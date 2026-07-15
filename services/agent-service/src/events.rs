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
pub const SUBJECT_AGENT_RELATIONSHIP_CHANGED: &str = "agente.relacion_cambio";
pub const SUBJECT_AGENT_CRITICAL_EVENT: &str = "agente.evento_critico";
pub const SUBJECT_ROSTER_PATCH: &str = "roster.patch";
pub const SUBJECT_GM_DECISION_REGISTERED: &str = "decision.gm_registrada";
pub const SUBJECT_SALARY_CAP_CALCULATED: &str = "salary_cap.calculado";
pub const SUBJECT_TRADE_PROPOSED: &str = "trade.propuesta_enviada";
pub const SUBJECT_TRADE_REJECTED: &str = "trade.rechazada";
pub const SUBJECT_TRADE_COUNTERED: &str = "trade.contraoferta";
pub const SUBJECT_TRADE_ACCEPTED: &str = "trade.aceptada";

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
    #[serde(default)]
    pub box_score: Vec<PlayerBoxScore>,
}

#[derive(Debug, Clone, PartialEq, Eq, Serialize, Deserialize)]
pub struct PlayerBoxScore {
    pub player_id: String,
    pub team_id: String,
    pub minutes: u8,
    pub points: u16,
    pub rebounds: u16,
    pub assists: u16,
    pub steals: u16,
    pub blocks: u16,
    pub turnovers: u16,
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

#[derive(Debug, Clone, PartialEq, Serialize, Deserialize)]
pub struct AgentRelationshipChangedEvent {
    #[serde(flatten)]
    pub meta: EventMeta,
    pub simulated_date: String,
    pub agent_a_id: String,
    pub agent_b_id: String,
    pub trust: f64,
    pub trend: String,
    pub last_event: String,
    pub short_history: Vec<String>,
    pub source_event_id: String,
    pub source_subject: String,
}

#[derive(Debug, Clone, PartialEq, Serialize, Deserialize)]
pub struct RosterPatchEnvelope {
    #[serde(rename = "type")]
    pub event_type: String,
    pub subject: String,
    pub game_id: String,
    pub patch: RosterStatePatch,
}

#[derive(Debug, Clone, PartialEq, Serialize, Deserialize)]
pub struct RosterStatePatch {
    pub simulated_date: String,
    pub source_event_id: String,
    pub source_subject: String,
    pub players: Vec<PlayerEmotionalPatch>,
}

#[derive(Debug, Clone, PartialEq, Serialize, Deserialize)]
pub struct PlayerEmotionalPatch {
    pub player_id: String,
    pub emotional_state: String,
    pub satisfaction: f64,
    pub loyalty: f64,
    pub ego: f64,
    pub competitive_drive: f64,
    pub city_connection: f64,
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

#[derive(Debug, Clone, PartialEq, Serialize, Deserialize)]
pub struct GMDecisionRegisteredEvent {
    #[serde(flatten)]
    pub meta: EventMeta,
    pub decision_id: String,
    pub kind: String,
    pub payload: BTreeMap<String, String>,
    pub simulated_date: String,
    pub agents_affected: Vec<String>,
    #[serde(default)]
    pub source_event_id: Option<String>,
    #[serde(default)]
    pub source_subject: Option<String>,
}

#[derive(Debug, Clone, PartialEq, Eq, Serialize, Deserialize)]
pub struct SalaryCapCalculatedEvent {
    #[serde(flatten)]
    pub meta: EventMeta,
    pub simulated_date: String,
    pub cap_base: i64,
    pub luxury_tax_line: i64,
    pub committed_salary: i64,
    pub cap_space: i64,
    pub luxury_tax_space: i64,
    pub roster_count: u8,
    pub status: String,
    pub near_luxury_tax: bool,
    pub projected_tax_payment: i64,
}

#[derive(Debug, Clone, PartialEq, Eq, Serialize, Deserialize)]
pub struct TradeProposedEvent {
    #[serde(flatten)]
    pub meta: EventMeta,
    pub proposal_id: String,
    pub simulated_date: String,
    pub rival_team_id: String,
    pub offered_player_id: String,
    pub offered_player_name: String,
    pub offered_salary: i64,
    pub requested_position: String,
    pub incoming_salary: i64,
    pub cap_space_after: i64,
}

#[derive(Debug, Clone, PartialEq, Eq, Serialize, Deserialize)]
pub struct TradeRejectedEvent {
    #[serde(flatten)]
    pub meta: EventMeta,
    pub proposal_id: String,
    pub simulated_date: String,
    pub rival_team_id: String,
    pub reason: String,
    pub detail: String,
}

#[derive(Debug, Clone, PartialEq, Eq, Serialize, Deserialize)]
pub struct TradeCounteredEvent {
    #[serde(flatten)]
    pub meta: EventMeta,
    pub proposal_id: String,
    pub simulated_date: String,
    pub rival_team_id: String,
    pub requested_position: String,
    pub additional_asset_required: String,
    pub detail: String,
}

#[derive(Debug, Clone, PartialEq, Eq, Serialize, Deserialize)]
pub struct TradeAcceptedEvent {
    #[serde(flatten)]
    pub meta: EventMeta,
    pub proposal_id: String,
    pub simulated_date: String,
    pub rival_team_id: String,
    pub outgoing_player_id: String,
    pub outgoing_player_name: String,
    pub incoming_player_id: String,
    pub incoming_player_name: String,
    pub incoming_position: String,
    pub incoming_rating: u8,
    pub incoming_salary: i64,
    #[serde(default)]
    pub accepted_additional_asset: Option<String>,
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
