mod events;

use std::time::Duration;

use anyhow::Context;
use async_nats::Client;
use events::{
    GridPoint, MapCell, MapData, MapGenerationProgress, MapGenerationRequest,
    SUBJECT_MAP_GENERATION_COMPLETE, SUBJECT_MAP_GENERATION_STARTED, SUBJECT_MAP_STADIUM_LOCATED,
    SUBJECT_MAP_TERRAIN_READY, SUBJECT_MAP_ZONES_CALCULATED,
};
use tokio::signal;
use tokio_stream::StreamExt;
use tracing::{error, info};

#[tokio::main]
async fn main() -> anyhow::Result<()> {
    tracing_subscriber::fmt()
        .with_env_filter(
            tracing_subscriber::EnvFilter::try_from_default_env().unwrap_or_else(|_| "info".into()),
        )
        .init();

    let nats_url =
        std::env::var("NATS_URL").unwrap_or_else(|_| "nats://127.0.0.1:4222".to_string());
    let client = async_nats::connect(&nats_url)
        .await
        .with_context(|| format!("connect nats at {nats_url}"))?;

    let mut subscription = client
        .subscribe(SUBJECT_MAP_GENERATION_STARTED.to_string())
        .await?;
    info!("map-service listening for {SUBJECT_MAP_GENERATION_STARTED}");

    loop {
        tokio::select! {
            maybe_message = subscription.next() => {
                let Some(message) = maybe_message else {
                    break;
                };

                match serde_json::from_slice::<MapGenerationRequest>(&message.payload) {
                    Ok(request) => {
                        if let Err(err) = process_generation(&client, request).await {
                            error!("process map generation: {err:#}");
                        }
                    }
                    Err(err) => error!("decode map generation request: {err}"),
                }
            }
            _ = signal::ctrl_c() => {
                info!("map-service shutdown signal received");
                break;
            }
        }
    }

    Ok(())
}

async fn process_generation(client: &Client, request: MapGenerationRequest) -> anyhow::Result<()> {
    let city_name = request
        .city_name
        .clone()
        .unwrap_or_else(|| "Nueva PulseCity".to_string());
    let terrain_map = generate_terrain_map(20, 20);
    let zoned_map = apply_zones(&terrain_map);
    let stadium = locate_stadium(&zoned_map);

    info!(game_id = %request.game_id, city = %city_name, "starting simple grid map generation");

    publish_progress(
        client,
        SUBJECT_MAP_TERRAIN_READY,
        MapGenerationProgress {
            game_id: request.game_id.clone(),
            stage: "terrain",
            progress: 25,
            message: format!("Terreno base generado para {city_name}"),
            map_data: Some(terrain_map.clone()),
            stadium: None,
        },
    )
    .await?;

    tokio::time::sleep(Duration::from_millis(350)).await;

    publish_progress(
        client,
        SUBJECT_MAP_ZONES_CALCULATED,
        MapGenerationProgress {
            game_id: request.game_id.clone(),
            stage: "zoning",
            progress: 55,
            message: "Zonas base calculadas sobre la grilla".to_string(),
            map_data: Some(zoned_map.clone()),
            stadium: None,
        },
    )
    .await?;

    tokio::time::sleep(Duration::from_millis(350)).await;

    publish_progress(
        client,
        SUBJECT_MAP_STADIUM_LOCATED,
        MapGenerationProgress {
            game_id: request.game_id.clone(),
            stage: "stadium",
            progress: 80,
            message: "Estadio ubicado en el distrito central".to_string(),
            map_data: Some(zoned_map.clone()),
            stadium: Some(stadium.clone()),
        },
    )
    .await?;

    tokio::time::sleep(Duration::from_millis(350)).await;

    publish_progress(
        client,
        SUBJECT_MAP_GENERATION_COMPLETE,
        MapGenerationProgress {
            game_id: request.game_id,
            stage: "complete",
            progress: 100,
            message: "Generacion de mapa completada".to_string(),
            map_data: Some(zoned_map),
            stadium: Some(stadium),
        },
    )
    .await?;

    Ok(())
}

async fn publish_progress(
    client: &Client,
    subject: &str,
    progress: MapGenerationProgress,
) -> anyhow::Result<()> {
    let payload = serde_json::to_vec(&progress)?;
    client.publish(subject.to_string(), payload.into()).await?;
    info!(game_id = %progress.game_id, subject, progress = progress.progress, "published map event");
    Ok(())
}

fn generate_terrain_map(width: usize, height: usize) -> MapData {
    let mut cells = Vec::with_capacity(height);

    for y in 0..height {
        let mut row = Vec::with_capacity(width);
        for x in 0..width {
            let terrain = if x < 3 || y < 2 {
                "water"
            } else if y > height - 4 || x > width - 3 {
                "hill"
            } else if (x + y) % 7 == 0 {
                "forest"
            } else {
                "plain"
            };

            row.push(MapCell {
                terrain,
                zone: None,
            });
        }
        cells.push(row);
    }

    MapData {
        width,
        height,
        cells,
    }
}

fn apply_zones(map: &MapData) -> MapData {
    let mut zoned = map.clone();

    for (y, row) in zoned.cells.iter_mut().enumerate() {
        for (x, cell) in row.iter_mut().enumerate() {
            if cell.terrain == "water" || cell.terrain == "hill" {
                continue;
            }

            cell.zone = Some(if x < zoned.width / 3 {
                "residential"
            } else if y > zoned.height / 2 && x > zoned.width / 2 {
                "industrial"
            } else if (zoned.width / 3..=(zoned.width * 2 / 3)).contains(&x) {
                "commercial"
            } else {
                "park"
            });
        }
    }

    zoned
}

fn locate_stadium(map: &MapData) -> GridPoint {
    let target_x = map.width / 2;
    let target_y = map.height / 2;

    for radius in 0..map.width.max(map.height) {
        let min_x = target_x.saturating_sub(radius);
        let max_x = (target_x + radius).min(map.width - 1);
        let min_y = target_y.saturating_sub(radius);
        let max_y = (target_y + radius).min(map.height - 1);

        for y in min_y..=max_y {
            for x in min_x..=max_x {
                let cell = &map.cells[y][x];
                if cell.terrain != "water" && cell.terrain != "hill" {
                    return GridPoint { x, y };
                }
            }
        }
    }

    GridPoint { x: 0, y: 0 }
}
