mod events;

use std::time::Duration;

use anyhow::Context;
use async_nats::Client;
use events::{
    MapGenerationProgress, MapGenerationRequest, SUBJECT_MAP_GENERATION_COMPLETE,
    SUBJECT_MAP_GENERATION_STARTED, SUBJECT_MAP_STADIUM_LOCATED, SUBJECT_MAP_TERRAIN_READY,
    SUBJECT_MAP_ZONES_CALCULATED,
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

    info!(game_id = %request.game_id, city = %city_name, "starting mock map generation");

    publish_progress(
        client,
        SUBJECT_MAP_TERRAIN_READY,
        MapGenerationProgress {
            game_id: request.game_id.clone(),
            stage: "terrain",
            progress: 25,
            message: format!("Terreno base generado para {city_name}"),
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
            message: "Zonas Voronoi calculadas".to_string(),
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
