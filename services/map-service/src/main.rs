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
    let seed = city_seed(&city_name);
    let terrain_map = generate_terrain_map(28, 28, seed);
    let zoned_map = apply_zones(&terrain_map, seed);
    let stadium = locate_stadium(&zoned_map);

    info!(game_id = %request.game_id, city = %city_name, seed, "starting procedural-lite map generation");

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
            message: "Zonas base calculadas sobre densidad y accesibilidad".to_string(),
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

fn city_seed(city_name: &str) -> f64 {
    let hash = city_name.bytes().fold(0u64, |acc, byte| {
        acc.wrapping_mul(109).wrapping_add(byte as u64 + 17)
    });

    (hash % 10_000) as f64 / 10_000.0
}

fn generate_terrain_map(width: usize, height: usize, seed: f64) -> MapData {
    let mut cells = Vec::with_capacity(height);
    let perlin = Perlin2D::new((seed * 1_000_000.0) as u64 + 1);

    for y in 0..height {
        let mut row = Vec::with_capacity(width);
        for x in 0..width {
            let terrain = classify_terrain(sample_elevation(x, y, width, height, seed, &perlin));

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

fn apply_zones(map: &MapData, seed: f64) -> MapData {
    let mut zoned = map.clone();
    let districts = district_seeds(zoned.width, zoned.height, seed);

    for (y, row) in zoned.cells.iter_mut().enumerate() {
        for (x, cell) in row.iter_mut().enumerate() {
            if !is_buildable(cell.terrain) {
                continue;
            }

            let district = nearest_district(x as f64, y as f64, &districts);
            let moisture = moisture_noise(x, y, zoned.width, zoned.height);
            let center_bias = 1.0
                - normalized_distance(
                    x as f64,
                    y as f64,
                    zoned.width as f64 * 0.52,
                    zoned.height as f64 * 0.48,
                );

            cell.zone = Some(if cell.terrain == "forest" && moisture > 0.56 {
                "park"
            } else if district.kind == "industrial" && moisture < 0.66 {
                "industrial"
            } else if district.kind == "commercial" || center_bias > 0.82 {
                "commercial"
            } else if district.kind == "park" || moisture > 0.67 {
                "park"
            } else if district.kind == "residential" || center_bias > 0.42 {
                "residential"
            } else {
                "residential"
            });
        }
    }

    zoned
}

fn locate_stadium(map: &MapData) -> GridPoint {
    let target_x = map.width as f64 * 0.54;
    let target_y = map.height as f64 * 0.50;
    let mut best = GridPoint {
        x: map.width / 2,
        y: map.height / 2,
    };
    let mut best_score = f64::MIN;

    for y in 0..map.height {
        for x in 0..map.width {
            let cell = &map.cells[y][x];
            if !is_buildable(cell.terrain) {
                continue;
            }

            let center_bias = 1.0 - normalized_distance(x as f64, y as f64, target_x, target_y);
            let zone_bonus = match cell.zone {
                Some("commercial") => 0.45,
                Some("residential") => 0.24,
                Some("park") => 0.12,
                Some("industrial") => -0.25,
                _ => 0.0,
            };
            let terrain_bonus = match cell.terrain {
                "plain" => 0.2,
                "forest" => -0.05,
                _ => 0.0,
            };
            let score = center_bias + zone_bonus + terrain_bonus;

            if score > best_score {
                best_score = score;
                best = GridPoint { x, y };
            }
        }
    }

    best
}

fn sample_elevation(
    x: usize,
    y: usize,
    width: usize,
    height: usize,
    seed: f64,
    perlin: &Perlin2D,
) -> f64 {
    let nx = x as f64 / width as f64 - 0.5;
    let ny = y as f64 / height as f64 - 0.5;
    let edge_falloff = ((nx * nx + ny * ny).sqrt() * 1.55).clamp(0.0, 1.0);

    let coarse = perlin.noise(nx * 2.8 + seed * 2.1, ny * 2.8 - seed * 1.7);
    let detail = perlin.noise(nx * 6.6 - seed * 1.1, ny * 6.1 + seed * 1.4);
    let ridges = perlin.noise(nx * 11.0 + 3.0, ny * 11.0 - 2.0);
    let ridge_shape = (0.5 - (ridges - 0.5).abs()) * 0.34;
    let coast_bias = (0.28 - (nx + 0.16).abs()).max(0.0) * 0.42;

    0.46 + coarse * 0.34 + detail * 0.18 + ridge_shape + coast_bias - edge_falloff * 0.72
}

fn classify_terrain(elevation: f64) -> &'static str {
    if elevation < 0.24 {
        "water"
    } else if elevation > 0.86 {
        "hill"
    } else if elevation > 0.58 {
        "forest"
    } else {
        "plain"
    }
}

fn moisture_noise(x: usize, y: usize, width: usize, height: usize) -> f64 {
    let nx = x as f64 / width as f64;
    let ny = y as f64 / height as f64;

    0.5 + ((nx * 8.0).sin() * 0.18) + ((ny * 10.0).cos() * 0.12) + (((nx - ny) * 6.0).sin() * 0.1)
}

fn normalized_distance(x: f64, y: f64, target_x: f64, target_y: f64) -> f64 {
    let dx = x - target_x;
    let dy = y - target_y;
    let distance = (dx * dx + dy * dy).sqrt();
    let max_distance = (target_x.powi(2) + target_y.powi(2)).sqrt().max(1.0);
    (distance / max_distance).clamp(0.0, 1.0)
}

fn is_buildable(terrain: &str) -> bool {
    !matches!(terrain, "water" | "hill")
}

#[derive(Clone, Copy)]
struct DistrictSeed {
    x: f64,
    y: f64,
    kind: &'static str,
}

fn district_seeds(width: usize, height: usize, seed: f64) -> Vec<DistrictSeed> {
    let columns = 3usize;
    let rows = 3usize;
    let cell_width = width as f64 / columns as f64;
    let cell_height = height as f64 / rows as f64;
    let district_kinds = [
        "residential",
        "commercial",
        "industrial",
        "park",
        "residential",
        "commercial",
        "park",
        "industrial",
        "residential",
    ];

    let mut districts = Vec::with_capacity(columns * rows);

    for row in 0..rows {
        for col in 0..columns {
            let idx = row * columns + col;
            let jitter_x = seeded_unit(seed, idx as u64 * 2 + 11);
            let jitter_y = seeded_unit(seed, idx as u64 * 2 + 29);
            let kind_index = ((seeded_unit(seed, idx as u64 * 7 + 53) * district_kinds.len() as f64)
                .floor() as usize)
                % district_kinds.len();

            let x = (col as f64 + 0.2 + jitter_x * 0.6) * cell_width;
            let y = (row as f64 + 0.2 + jitter_y * 0.6) * cell_height;

            districts.push(DistrictSeed {
                x,
                y,
                kind: district_kinds[(idx + kind_index) % district_kinds.len()],
            });
        }
    }

    districts
}

fn nearest_district(x: f64, y: f64, districts: &[DistrictSeed]) -> DistrictSeed {
    let mut best = districts[0];
    let mut best_distance = f64::MAX;

    for district in districts {
        let dx = x - district.x;
        let dy = y - district.y;
        let distance = dx * dx + dy * dy;
        if distance < best_distance {
            best_distance = distance;
            best = *district;
        }
    }

    best
}

fn lerp(a: f64, b: f64, t: f64) -> f64 {
    a + (b - a) * t
}

fn seeded_unit(seed: f64, salt: u64) -> f64 {
    let seed_int = (seed * 1_000_000.0) as u64;
    let mixed = seed_int
        .wrapping_mul(0x9E3779B185EBCA87)
        .wrapping_add(salt.wrapping_mul(0xC2B2AE3D27D4EB4F));
    let value = ((mixed as f64) * 0.0000000000000002).sin() * 43758.5453123;
    value.fract().abs()
}

struct Perlin2D {
    perm: [u8; 512],
}

impl Perlin2D {
    fn new(seed: u64) -> Self {
        let mut values = [0u8; 256];
        for (i, value) in values.iter_mut().enumerate() {
            *value = i as u8;
        }

        let mut state = seed.max(1);
        for i in (1..256).rev() {
            state = state.wrapping_mul(6364136223846793005).wrapping_add(1);
            let j = (state % (i as u64 + 1)) as usize;
            values.swap(i, j);
        }

        let mut perm = [0u8; 512];
        for i in 0..512 {
            perm[i] = values[i & 255];
        }

        Self { perm }
    }

    fn noise(&self, x: f64, y: f64) -> f64 {
        let xi = (x.floor() as i32 & 255) as usize;
        let yi = (y.floor() as i32 & 255) as usize;
        let xf = x - x.floor();
        let yf = y - y.floor();

        let u = fade(xf);
        let v = fade(yf);

        let aa = self.perm[self.perm[xi] as usize + yi] as usize;
        let ab = self.perm[self.perm[xi] as usize + yi + 1] as usize;
        let ba = self.perm[self.perm[xi + 1] as usize + yi] as usize;
        let bb = self.perm[self.perm[xi + 1] as usize + yi + 1] as usize;

        let x1 = lerp(grad2(aa, xf, yf), grad2(ba, xf - 1.0, yf), u);
        let x2 = lerp(grad2(ab, xf, yf - 1.0), grad2(bb, xf - 1.0, yf - 1.0), u);
        let value = lerp(x1, x2, v);

        (value + 1.0) * 0.5
    }
}

fn fade(t: f64) -> f64 {
    t * t * t * (t * (t * 6.0 - 15.0) + 10.0)
}

fn grad2(hash: usize, x: f64, y: f64) -> f64 {
    match hash & 7 {
        0 => x + y,
        1 => -x + y,
        2 => x - y,
        3 => -x - y,
        4 => x,
        5 => -x,
        6 => y,
        _ => -y,
    }
}
