use std::collections::BTreeMap;

use serde::{Deserialize, Serialize};

pub const SUBJECT_TIME_SESSION_STARTED: &str = "tiempo.sesion_iniciada";
pub const SUBJECT_TIME_SESSION_ENDED: &str = "tiempo.sesion_terminada";
pub const SUBJECT_TIME_SPEED_CHANGED: &str = "tiempo.velocidad_cambiada";
pub const SUBJECT_TIME_PAUSE_CHANGED: &str = "tiempo.pausa_activada";
pub const SUBJECT_TIME_DAY_ADVANCED: &str = "tiempo.dia_avanzado";
pub const SUBJECT_MAP_GENERATION_STARTED: &str = "mapa.generacion_iniciada";
pub const SUBJECT_MATCH_FINISHED: &str = "partido.terminado";
pub const SUBJECT_AGENT_STATE_CHANGED: &str = "agente.estado_cambio";
pub const SUBJECT_AGENT_CRITICAL_EVENT: &str = "agente.evento_critico";

#[derive(Debug, Clone, PartialEq, Eq, Serialize, Deserialize)]
pub struct EventMeta {
    pub event_id: String,
    pub game_id: String,
    pub occurred_at: String,
    pub schema_version: u16,
}

#[derive(Debug, Clone, PartialEq, Eq, Serialize, Deserialize)]
pub struct SessionStartedEvent {
    #[serde(flatten)]
    pub meta: EventMeta,
    pub session_id: String,
    pub client_id: String,
}

#[derive(Debug, Clone, PartialEq, Eq, Serialize, Deserialize)]
pub struct SessionEndedEvent {
    #[serde(flatten)]
    pub meta: EventMeta,
    pub session_id: String,
    pub reason: String,
}

#[derive(Debug, Clone, PartialEq, Eq, Serialize, Deserialize)]
pub struct SpeedChangedEvent {
    #[serde(flatten)]
    pub meta: EventMeta,
    pub speed: u8,
}

#[derive(Debug, Clone, PartialEq, Eq, Serialize, Deserialize)]
pub struct PauseChangedEvent {
    #[serde(flatten)]
    pub meta: EventMeta,
    pub paused: bool,
}

#[derive(Debug, Clone, PartialEq, Eq, Serialize, Deserialize)]
pub struct DayAdvancedEvent {
    #[serde(flatten)]
    pub meta: EventMeta,
    pub simulated_date: String,
    pub speed: u8,
    pub days_processed: u16,
}

#[derive(Debug, Clone, PartialEq, Eq, Serialize, Deserialize)]
pub struct MapGenerationStartedEvent {
    pub game_id: String,
    pub city_name: String,
    #[serde(default)]
    pub seed: Option<u64>,
}

#[derive(Debug, Clone, PartialEq, Serialize, Deserialize)]
pub struct MatchFinishedEvent {
    #[serde(flatten)]
    pub meta: EventMeta,
    pub match_id: String,
    pub simulated_date: String,
    pub home_team: MatchTeam,
    pub away_team: MatchTeam,
    pub home_score: u16,
    pub away_score: u16,
    pub winner_team_id: String,
    pub seed: u64,
}

#[derive(Debug, Clone, PartialEq, Eq, Serialize, Deserialize)]
pub struct MatchTeam {
    pub team_id: String,
    pub name: String,
    pub abbreviation: String,
}

#[derive(Debug, Clone, PartialEq, Serialize, Deserialize)]
pub struct AgentStateChangedEvent {
    #[serde(flatten)]
    pub meta: EventMeta,
    pub simulated_date: String,
    pub agent_id: String,
    pub source_event_id: String,
    pub source_subject: String,
    pub mood: String,
    pub state: BTreeMap<String, f64>,
    pub summary: String,
}

#[derive(Debug, Clone, PartialEq, Eq, Serialize, Deserialize)]
pub struct AgentCriticalEvent {
    #[serde(flatten)]
    pub meta: EventMeta,
    pub simulated_date: String,
    pub agent_id: String,
    pub severity: String,
    pub source_event_id: String,
    pub source_subject: String,
    pub title: String,
    pub summary: String,
}

#[cfg(test)]
mod tests {
    use super::*;

    #[test]
    fn subjects_follow_m2_contracts() {
        assert_eq!(SUBJECT_TIME_DAY_ADVANCED, "tiempo.dia_avanzado");
        assert_eq!(SUBJECT_MATCH_FINISHED, "partido.terminado");
        assert_eq!(SUBJECT_AGENT_STATE_CHANGED, "agente.estado_cambio");
    }
}
