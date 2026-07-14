use serde::{Deserialize, Serialize};

pub const SUBJECT_MATCH_SCHEDULED: &str = "partido.programado";
pub const SUBJECT_MATCH_STARTING: &str = "partido.iniciando";
pub const SUBJECT_MATCH_FINISHED: &str = "partido.terminado";

#[derive(Debug, Clone, PartialEq, Eq, Serialize, Deserialize)]
pub struct EventMeta {
    pub event_id: String,
    pub game_id: String,
    pub occurred_at: String,
    pub schema_version: u16,
}

#[derive(Debug, Clone, PartialEq, Serialize, Deserialize)]
pub struct MatchScheduledEvent {
    #[serde(flatten)]
    pub meta: EventMeta,
    pub match_id: String,
    pub simulated_date: String,
    pub home_team: MatchTeam,
    pub away_team: MatchTeam,
    #[serde(default, skip_serializing_if = "Option::is_none")]
    pub home_tactics: Option<MatchTacticalContext>,
    #[serde(default, skip_serializing_if = "Option::is_none")]
    pub away_tactics: Option<MatchTacticalContext>,
    pub players: Vec<MatchPlayer>,
    pub seed: u64,
}

#[derive(Debug, Clone, PartialEq, Eq, Serialize, Deserialize)]
pub struct MatchStartingEvent {
    #[serde(flatten)]
    pub meta: EventMeta,
    pub match_id: String,
    pub simulated_date: String,
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
    pub box_score: Vec<PlayerBoxScore>,
    pub key_moments: Vec<KeyMoment>,
}

#[derive(Debug, Clone, PartialEq, Serialize, Deserialize)]
pub struct MatchTeam {
    pub team_id: String,
    pub name: String,
    pub abbreviation: String,
    pub rating: u8,
    pub offense_rating: u8,
    pub defense_rating: u8,
    pub pace: u8,
    pub home_court_advantage: i8,
}

#[derive(Debug, Clone, PartialEq, Eq, Serialize, Deserialize)]
pub struct MatchTacticalContext {
    pub system: String,
    pub rotation_preference: String,
    pub flexibility: u8,
}

#[derive(Debug, Clone, PartialEq, Eq, Serialize, Deserialize)]
pub struct MatchPlayer {
    pub player_id: String,
    pub team_id: String,
    #[serde(default, skip_serializing_if = "Option::is_none")]
    pub expected_minutes: Option<u8>,
    pub rating: u8,
    pub scoring: u8,
    pub rebounding: u8,
    pub playmaking: u8,
    pub defense: u8,
    pub stamina: u8,
    pub fatigue: u8,
    pub emotional_state: i8,
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
pub struct KeyMoment {
    pub quarter: u8,
    pub clock: String,
    pub kind: String,
    pub description: String,
    pub team_id: String,
    #[serde(skip_serializing_if = "Option::is_none")]
    pub player_id: Option<String>,
}

#[cfg(test)]
mod tests {
    use super::*;

    #[test]
    fn subjects_follow_m2_contracts() {
        assert_eq!(SUBJECT_MATCH_SCHEDULED, "partido.programado");
        assert_eq!(SUBJECT_MATCH_STARTING, "partido.iniciando");
        assert_eq!(SUBJECT_MATCH_FINISHED, "partido.terminado");
    }

    #[test]
    fn match_scheduled_event_accepts_m2_payload_without_tactics() {
        let payload = r#"{
            "event_id": "event-1",
            "game_id": "game-1",
            "occurred_at": "2026-10-22T00:00:00Z",
            "schema_version": 1,
            "match_id": "match-1",
            "simulated_date": "2026-10-22",
            "home_team": {
                "team_id": "home",
                "name": "Home",
                "abbreviation": "HOM",
                "rating": 78,
                "offense_rating": 79,
                "defense_rating": 76,
                "pace": 99,
                "home_court_advantage": 3
            },
            "away_team": {
                "team_id": "away",
                "name": "Away",
                "abbreviation": "AWY",
                "rating": 77,
                "offense_rating": 76,
                "defense_rating": 78,
                "pace": 97,
                "home_court_advantage": 2
            },
            "players": [],
            "seed": 42
        }"#;

        let event: MatchScheduledEvent = serde_json::from_str(payload).expect("valid m2 payload");

        assert_eq!(event.home_tactics, None);
        assert_eq!(event.away_tactics, None);
    }
}
