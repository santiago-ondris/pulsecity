use serde::{Deserialize, Serialize};

pub const SUBJECT_MAP_GENERATION_STARTED: &str = "mapa.generacion_iniciada";
pub const SUBJECT_MAP_TERRAIN_READY: &str = "mapa.terreno_listo";
pub const SUBJECT_MAP_ZONES_CALCULATED: &str = "mapa.zonas_calculadas";
pub const SUBJECT_MAP_STADIUM_LOCATED: &str = "mapa.estadio_ubicado";
pub const SUBJECT_MAP_GENERATION_COMPLETE: &str = "mapa.generacion_completa";

#[derive(Debug, Deserialize)]
pub struct MapGenerationRequest {
    pub game_id: String,
    pub city_name: Option<String>,
}

#[derive(Debug, Serialize)]
pub struct MapGenerationProgress {
    pub game_id: String,
    pub stage: &'static str,
    pub progress: u8,
    pub message: String,
}
