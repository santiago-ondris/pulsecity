use std::collections::BTreeMap;

use crate::events::{
    AgentStateChangedEvent, EventMeta, MatchFinishedEvent, SUBJECT_MATCH_FINISHED,
};

pub const OWN_TEAM_ID: &str = "pulsecity";
pub const CORE_AGENT_IDS: [&str; 5] = [
    "owner",
    "head_coach",
    "cfo",
    "scouting_director",
    "sports_psychologist",
];

const MIN_STATE_VALUE: f64 = -1.0;
const MAX_STATE_VALUE: f64 = 1.0;
const SCHEMA_VERSION: u16 = 1;

#[derive(Debug, Clone, PartialEq)]
pub struct CoreAgentState {
    pub game_id: String,
    pub agent_id: String,
    pub mood: String,
    pub state: BTreeMap<String, f64>,
    pub last_match_id: Option<String>,
}

#[derive(Debug, Clone, PartialEq)]
pub struct AgentStateChange {
    pub state: CoreAgentState,
    pub event: AgentStateChangedEvent,
}

#[must_use]
pub fn default_core_agent_states(game_id: &str) -> Vec<CoreAgentState> {
    CORE_AGENT_IDS
        .iter()
        .map(|agent_id| default_core_agent_state(game_id, agent_id))
        .collect()
}

#[must_use]
pub fn default_core_agent_state(game_id: &str, agent_id: &str) -> CoreAgentState {
    let state = match agent_id {
        "owner" => map_from_pairs(&[
            ("sporting_trust", 0.0),
            ("business_trust", 0.0),
            ("patience_remaining", 0.75),
            ("satisfaction", 0.0),
        ]),
        "head_coach" => map_from_pairs(&[
            ("gm_trust", 0.0),
            ("roster_satisfaction", 0.0),
            ("results_pressure", 0.25),
            ("locker_room_relationship", 0.0),
        ]),
        "cfo" => map_from_pairs(&[
            ("financial_trust", 0.0),
            ("budget_alert", 0.15),
            ("financial_conservatism", 0.55),
        ]),
        "scouting_director" => map_from_pairs(&[
            ("criteria_trust", 0.0),
            ("motivation", 0.3),
            ("perceived_precision", 0.0),
        ]),
        "sports_psychologist" => map_from_pairs(&[
            ("locker_room_climate", 0.0),
            ("emotional_alert", 0.2),
            ("player_trust", 0.0),
        ]),
        _ => BTreeMap::new(),
    };

    CoreAgentState {
        game_id: game_id.to_string(),
        agent_id: agent_id.to_string(),
        mood: "calm".to_string(),
        state,
        last_match_id: None,
    }
}

#[must_use]
pub fn apply_match_finished(
    current_states: Vec<CoreAgentState>,
    event: &MatchFinishedEvent,
    occurred_at: String,
) -> Vec<AgentStateChange> {
    let context = MatchContext::from_event(event);

    current_states
        .into_iter()
        .map(|state| apply_match_to_agent(state, event, &context, &occurred_at))
        .collect()
}

fn apply_match_to_agent(
    mut state: CoreAgentState,
    event: &MatchFinishedEvent,
    context: &MatchContext,
    occurred_at: &str,
) -> AgentStateChange {
    match state.agent_id.as_str() {
        "owner" => apply_owner(&mut state, context),
        "head_coach" => apply_head_coach(&mut state, context),
        "cfo" => apply_cfo(&mut state, context),
        "scouting_director" => apply_scouting_director(&mut state, context),
        "sports_psychologist" => apply_sports_psychologist(&mut state, context),
        _ => {}
    }

    state.last_match_id = Some(event.match_id.clone());
    let summary = summarize_agent_change(&state.agent_id, context);
    let event = AgentStateChangedEvent {
        meta: EventMeta {
            event_id: format!(
                "agent-state-{}-{}-{}",
                event.meta.game_id, event.match_id, state.agent_id
            ),
            game_id: event.meta.game_id.clone(),
            occurred_at: occurred_at.to_string(),
            schema_version: SCHEMA_VERSION,
        },
        simulated_date: event.simulated_date.clone(),
        agent_id: state.agent_id.clone(),
        source_event_id: event.meta.event_id.clone(),
        source_subject: SUBJECT_MATCH_FINISHED.to_string(),
        mood: state.mood.clone(),
        state: state.state.clone(),
        summary,
    };

    AgentStateChange { state, event }
}

fn apply_owner(state: &mut CoreAgentState, context: &MatchContext) {
    let result_delta = if context.won { 1.0 } else { -1.0 };
    adjust(&mut state.state, "sporting_trust", 0.06 * result_delta);
    adjust(
        &mut state.state,
        "business_trust",
        if context.home_game {
            0.03 * result_delta
        } else {
            0.015 * result_delta
        },
    );
    adjust(
        &mut state.state,
        "patience_remaining",
        if context.won { 0.02 } else { -0.05 },
    );
    adjust(&mut state.state, "satisfaction", 0.05 * result_delta);

    if context.blowout {
        adjust(&mut state.state, "sporting_trust", 0.03 * result_delta);
        adjust(&mut state.state, "satisfaction", 0.02 * result_delta);
    }

    let patience = metric(&state.state, "patience_remaining");
    state.mood = if !context.won && patience < 0.35 {
        "frustrated"
    } else if !context.won {
        "concerned"
    } else if context.blowout {
        "excited"
    } else {
        "calm"
    }
    .to_string();
}

fn apply_head_coach(state: &mut CoreAgentState, context: &MatchContext) {
    let result_delta = if context.won { 1.0 } else { -1.0 };
    adjust(&mut state.state, "gm_trust", 0.03 * result_delta);
    adjust(&mut state.state, "roster_satisfaction", 0.05 * result_delta);
    adjust(
        &mut state.state,
        "results_pressure",
        if context.won { -0.04 } else { 0.06 },
    );
    adjust(
        &mut state.state,
        "locker_room_relationship",
        0.035 * result_delta,
    );

    if context.close_game {
        adjust(&mut state.state, "results_pressure", 0.015);
    }
    if context.blowout {
        adjust(&mut state.state, "results_pressure", -0.02 * result_delta);
    }

    let pressure = metric(&state.state, "results_pressure");
    state.mood = if pressure > 0.65 {
        "pressured"
    } else if !context.won {
        "frustrated"
    } else {
        "calm"
    }
    .to_string();
}

fn apply_cfo(state: &mut CoreAgentState, context: &MatchContext) {
    let result_delta = if context.won { 1.0 } else { -1.0 };
    let home_multiplier = if context.home_game { 1.0 } else { 0.6 };
    adjust(
        &mut state.state,
        "financial_trust",
        0.02 * result_delta * home_multiplier,
    );
    adjust(
        &mut state.state,
        "budget_alert",
        if context.won {
            -0.01 * home_multiplier
        } else {
            0.025 * home_multiplier
        },
    );
    adjust(
        &mut state.state,
        "financial_conservatism",
        if context.won { -0.01 } else { 0.02 },
    );

    state.mood = if metric(&state.state, "budget_alert") > 0.55 {
        "concerned"
    } else {
        "calm"
    }
    .to_string();
}

fn apply_scouting_director(state: &mut CoreAgentState, context: &MatchContext) {
    let result_delta = if context.won { 1.0 } else { -1.0 };
    adjust(&mut state.state, "criteria_trust", 0.025 * result_delta);
    adjust(
        &mut state.state,
        "motivation",
        if context.close_game { 0.025 } else { 0.015 },
    );
    adjust(
        &mut state.state,
        "perceived_precision",
        if context.won { 0.02 } else { -0.015 },
    );

    if !context.won && context.blowout {
        adjust(&mut state.state, "criteria_trust", -0.02);
        adjust(&mut state.state, "perceived_precision", -0.015);
    }

    state.mood = if !context.won && context.blowout {
        "concerned"
    } else if context.won {
        "excited"
    } else {
        "calm"
    }
    .to_string();
}

fn apply_sports_psychologist(state: &mut CoreAgentState, context: &MatchContext) {
    let result_delta = if context.won { 1.0 } else { -1.0 };
    adjust(
        &mut state.state,
        "locker_room_climate",
        0.055 * result_delta,
    );
    adjust(
        &mut state.state,
        "emotional_alert",
        if context.won { -0.03 } else { 0.055 },
    );
    adjust(&mut state.state, "player_trust", 0.03 * result_delta);

    if context.close_game {
        adjust(&mut state.state, "emotional_alert", 0.015);
    }

    state.mood = if metric(&state.state, "emotional_alert") > 0.55 {
        "concerned"
    } else if context.won {
        "calm"
    } else {
        "pressured"
    }
    .to_string();
}

fn summarize_agent_change(agent_id: &str, context: &MatchContext) -> String {
    let result = if context.won { "victoria" } else { "derrota" };
    let venue = if context.home_game {
        "en casa"
    } else {
        "como visitante"
    };
    let intensity = if context.blowout {
        "amplia"
    } else if context.close_game {
        "cerrada"
    } else {
        "normal"
    };

    match agent_id {
        "owner" => format!("El owner ajusta su confianza tras una {result} {intensity} {venue}."),
        "head_coach" => format!("El head coach recalibra presion y satisfaccion tras la {result}."),
        "cfo" => format!("El CFO actualiza su lectura financiera despues de la {result} {venue}."),
        "scouting_director" => {
            format!(
                "El scouting director reevalua su criterio despues de una {result} {intensity}."
            )
        }
        "sports_psychologist" => {
            format!("La sports psychologist ajusta el clima emocional tras la {result}.")
        }
        _ => format!("El agente reacciona al resultado del partido."),
    }
}

fn adjust(state: &mut BTreeMap<String, f64>, key: &str, delta: f64) {
    let current = metric(state, key);
    state.insert(key.to_string(), clamp(current + delta));
}

fn metric(state: &BTreeMap<String, f64>, key: &str) -> f64 {
    state.get(key).copied().unwrap_or_default()
}

fn clamp(value: f64) -> f64 {
    value.clamp(MIN_STATE_VALUE, MAX_STATE_VALUE)
}

fn map_from_pairs(pairs: &[(&str, f64)]) -> BTreeMap<String, f64> {
    pairs
        .iter()
        .map(|(key, value)| ((*key).to_string(), *value))
        .collect()
}

#[derive(Debug, Clone, Copy, PartialEq, Eq)]
struct MatchContext {
    won: bool,
    home_game: bool,
    close_game: bool,
    blowout: bool,
}

impl MatchContext {
    fn from_event(event: &MatchFinishedEvent) -> Self {
        let own_score = if event.home_team.team_id == OWN_TEAM_ID {
            event.home_score
        } else {
            event.away_score
        };
        let opponent_score = if event.home_team.team_id == OWN_TEAM_ID {
            event.away_score
        } else {
            event.home_score
        };
        let margin = own_score.abs_diff(opponent_score);

        Self {
            won: event.winner_team_id == OWN_TEAM_ID,
            home_game: event.home_team.team_id == OWN_TEAM_ID,
            close_game: margin <= 5,
            blowout: margin >= 15,
        }
    }
}

#[cfg(test)]
mod tests {
    use super::*;
    use crate::events::{EventMeta, MatchTeam};

    #[test]
    fn defaults_create_the_five_m2_core_agents() {
        let states = default_core_agent_states("game-1");

        assert_eq!(states.len(), 5);
        assert!(states.iter().any(|state| state.agent_id == "owner"));
        assert!(states.iter().any(|state| state.agent_id == "head_coach"));
        assert!(states.iter().any(|state| state.agent_id == "cfo"));
        assert!(
            states
                .iter()
                .any(|state| state.agent_id == "scouting_director")
        );
        assert!(
            states
                .iter()
                .any(|state| state.agent_id == "sports_psychologist")
        );
    }

    #[test]
    fn win_improves_owner_and_reduces_coach_pressure() {
        let event = sample_match(true, 112, 101);
        let changes = apply_match_finished(
            default_core_agent_states("game-1"),
            &event,
            "2026-05-24T00:00:00Z".to_string(),
        );

        let owner = state_for(&changes, "owner");
        assert!(owner.state["sporting_trust"] > 0.0);
        assert!(owner.state["patience_remaining"] > 0.75);
        assert_eq!(owner.mood, "calm");

        let coach = state_for(&changes, "head_coach");
        assert!(coach.state["results_pressure"] < 0.25);
    }

    #[test]
    fn blowout_loss_increases_pressure_and_emotional_alert() {
        let event = sample_match(false, 88, 111);
        let changes = apply_match_finished(
            default_core_agent_states("game-1"),
            &event,
            "2026-05-24T00:00:00Z".to_string(),
        );

        let coach = state_for(&changes, "head_coach");
        assert!(coach.state["results_pressure"] > 0.25);

        let psychologist = state_for(&changes, "sports_psychologist");
        assert!(psychologist.state["emotional_alert"] > 0.2);
    }

    #[test]
    fn state_change_events_are_deterministic_per_match_and_agent() {
        let event = sample_match(true, 112, 101);
        let changes = apply_match_finished(
            default_core_agent_states("game-1"),
            &event,
            "2026-05-24T00:00:00Z".to_string(),
        );

        let coach = changes
            .iter()
            .find(|change| change.state.agent_id == "head_coach")
            .expect("coach change exists");

        assert_eq!(
            coach.event.meta.event_id,
            "agent-state-game-1-match-1-head_coach"
        );
        assert_eq!(coach.event.source_subject, "partido.terminado");
        assert_eq!(coach.event.source_event_id, "match-finished-match-1");
    }

    fn state_for<'a>(changes: &'a [AgentStateChange], agent_id: &str) -> &'a CoreAgentState {
        &changes
            .iter()
            .find(|change| change.state.agent_id == agent_id)
            .expect("agent state exists")
            .state
    }

    fn sample_match(own_home: bool, own_score: u16, opponent_score: u16) -> MatchFinishedEvent {
        let own = MatchTeam {
            team_id: OWN_TEAM_ID.to_string(),
            name: "PulseCity".to_string(),
            abbreviation: "PUL".to_string(),
        };
        let opponent = MatchTeam {
            team_id: "rival-1".to_string(),
            name: "Rival".to_string(),
            abbreviation: "RIV".to_string(),
        };
        let (home_team, away_team, home_score, away_score) = if own_home {
            (own, opponent, own_score, opponent_score)
        } else {
            (opponent, own, opponent_score, own_score)
        };
        let winner_team_id = if own_score > opponent_score {
            OWN_TEAM_ID
        } else {
            "rival-1"
        };

        MatchFinishedEvent {
            meta: EventMeta {
                event_id: "match-finished-match-1".to_string(),
                game_id: "game-1".to_string(),
                occurred_at: "2026-05-24T00:00:00Z".to_string(),
                schema_version: 1,
            },
            match_id: "match-1".to_string(),
            simulated_date: "2026-10-22".to_string(),
            home_team,
            away_team,
            home_score,
            away_score,
            winner_team_id: winner_team_id.to_string(),
            seed: 123,
        }
    }
}
