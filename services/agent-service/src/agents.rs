use std::collections::BTreeMap;

mod catalog;
mod match_reactions;
mod medical_reactions;
mod salary_cap_reactions;
mod trade_reactions;

pub use catalog::{canonical_relationship_seeds, individual_agent_templates};
use catalog::{rival_gm_profile, rival_team_templates};
pub use match_reactions::{
    apply_match_finished, apply_match_to_player_agents, apply_match_to_relationships,
};
pub use medical_reactions::apply_gm_decision_to_relationships;
pub use salary_cap_reactions::apply_salary_cap_to_core_agents;
pub use trade_reactions::{
    RivalGMTradeEvaluation, apply_trade_accepted_to_player_agents, evaluate_trade_proposal,
};

use crate::events::{AgentRelationshipChangedEvent, AgentStateChangedEvent, RosterPatchEnvelope};

pub const OWN_TEAM_ID: &str = "pulsecity";
pub const CORE_AGENT_IDS: [&str; 5] = [
    "owner",
    "head_coach",
    "cfo",
    "scouting_director",
    "sports_psychologist",
];

pub const INDIVIDUAL_AGENT_COUNT: usize = 30;
pub const RIVAL_GM_COUNT: usize = 30;

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

#[derive(Debug, Clone, PartialEq)]
pub struct TeamRosterPlayer {
    pub player_id: String,
    pub game_id: String,
    pub full_name: String,
    pub position: String,
    pub overall_rating: u8,
    pub roster_status: String,
}

#[derive(Debug, Clone, PartialEq)]
pub struct PlayerAgentState {
    pub game_id: String,
    pub player_id: String,
    pub full_name: String,
    pub position: String,
    pub emotional_state: String,
    pub satisfaction: f64,
    pub loyalty: f64,
    pub ego: f64,
    pub competitive_drive: f64,
    pub city_connection: f64,
    pub last_match_id: Option<String>,
}

#[derive(Debug, Clone, PartialEq)]
pub struct MatchAgentReactions {
    pub core_agent_changes: Vec<AgentStateChange>,
    pub roster_patch: Option<RosterPatchEnvelope>,
    pub relationship_changes: Vec<AgentRelationshipChange>,
}

#[derive(Debug, Clone, PartialEq)]
pub struct AgentRelationship {
    pub game_id: String,
    pub agent_a_id: String,
    pub agent_b_id: String,
    pub trust: f64,
    pub last_event: String,
    pub trend: String,
    pub short_history: Vec<String>,
    pub last_source_event_id: Option<String>,
}

#[derive(Debug, Clone, PartialEq)]
pub struct AgentRelationshipChange {
    pub relationship: AgentRelationship,
    pub event: AgentRelationshipChangedEvent,
}

#[derive(Debug, Clone, PartialEq)]
pub struct AgentRelationshipSeed {
    pub agent_a_id: &'static str,
    pub agent_b_id: &'static str,
    pub trust: f64,
    pub last_event: &'static str,
    pub trend: &'static str,
    pub short_history: Vec<&'static str>,
}

#[derive(Debug, Clone, PartialEq)]
pub struct IndividualAgentTemplate {
    pub agent_id: &'static str,
    pub display_name: &'static str,
    pub category: &'static str,
    pub role: &'static str,
    pub domain: &'static str,
    pub emotional_state: &'static str,
    pub confidence: f64,
    pub satisfaction: f64,
    pub loyalty: f64,
    pub role_performance: f64,
    pub state: BTreeMap<String, f64>,
    pub agenda: BTreeMap<String, String>,
}

#[derive(Debug, Clone, PartialEq)]
pub struct IndividualAgentState {
    pub game_id: String,
    pub agent_id: String,
    pub display_name: String,
    pub category: String,
    pub role: String,
    pub domain: String,
    pub emotional_state: String,
    pub confidence: f64,
    pub satisfaction: f64,
    pub loyalty: f64,
    pub role_performance: f64,
    pub state: BTreeMap<String, f64>,
    pub agenda: BTreeMap<String, String>,
}

#[derive(Debug, Clone, PartialEq)]
pub struct RivalGMProfile {
    pub game_id: String,
    pub rival_team_id: String,
    pub gm_agent_id: String,
    pub display_name: String,
    pub team_name: String,
    pub negotiation_style: String,
    pub urgency_current: f64,
    pub build_philosophy: String,
    pub roster_needs: Vec<String>,
    pub relationship_trust: f64,
    pub relationship_history: Vec<String>,
    pub last_interaction_event_id: Option<String>,
}

#[must_use]
pub fn default_individual_agent_states(game_id: &str) -> Vec<IndividualAgentState> {
    individual_agent_templates()
        .into_iter()
        .map(|template| template.into_state(game_id))
        .collect()
}

#[must_use]
pub fn default_rival_gms(game_id: &str) -> Vec<RivalGMProfile> {
    rival_team_templates()
        .into_iter()
        .enumerate()
        .map(|(index, team)| rival_gm_profile(game_id, index, team))
        .collect()
}

#[must_use]
pub fn default_player_agent_state(player: &TeamRosterPlayer) -> PlayerAgentState {
    let rating_factor = (f64::from(player.overall_rating).clamp(60.0, 90.0) - 60.0) / 30.0;
    let ego = match player.position.as_str() {
        "PG" | "SG" | "SF" => 0.42 + rating_factor * 0.24,
        _ => 0.34 + rating_factor * 0.20,
    };

    PlayerAgentState {
        game_id: player.game_id.clone(),
        player_id: player.player_id.clone(),
        full_name: player.full_name.clone(),
        position: player.position.clone(),
        emotional_state: "steady".to_string(),
        satisfaction: 0.04,
        loyalty: 0.62,
        ego: clamp_unit(ego),
        competitive_drive: clamp_unit(0.58 + rating_factor * 0.22),
        city_connection: 0.35,
        last_match_id: None,
    }
}

#[must_use]
pub fn default_agent_relationships(game_id: &str) -> Vec<AgentRelationship> {
    canonical_relationship_seeds()
        .into_iter()
        .map(|seed| seed.into_relationship(game_id))
        .collect()
}

#[must_use]
pub fn relationship_key(agent_a_id: &str, agent_b_id: &str) -> String {
    if agent_a_id <= agent_b_id {
        format!("{agent_a_id}:{agent_b_id}")
    } else {
        format!("{agent_b_id}:{agent_a_id}")
    }
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

fn apply_relationship_delta(
    relationship: &mut AgentRelationship,
    delta: f64,
    reason: &str,
    source_event_id: &str,
) {
    relationship.trust = clamp(relationship.trust + delta);
    relationship.trend = if delta > 0.0 {
        "improving"
    } else if delta < 0.0 {
        "deteriorating"
    } else {
        "stable"
    }
    .to_string();
    relationship.last_event = reason.to_string();
    relationship.last_source_event_id = Some(source_event_id.to_string());
    relationship.short_history.push(reason.to_string());
    if relationship.short_history.len() > 5 {
        relationship.short_history.remove(0);
    }
}

fn adjust(state: &mut BTreeMap<String, f64>, key: &str, delta: f64) {
    let current = metric(state, key);
    state.insert(key.to_string(), adjusted(current, delta));
}

fn adjusted(current: f64, delta: f64) -> f64 {
    clamp(current + delta)
}

fn metric(state: &BTreeMap<String, f64>, key: &str) -> f64 {
    state.get(key).copied().unwrap_or_default()
}

fn clamp(value: f64) -> f64 {
    value.clamp(MIN_STATE_VALUE, MAX_STATE_VALUE)
}

fn clamp_unit(value: f64) -> f64 {
    value.clamp(0.0, 1.0)
}

fn map_from_pairs(pairs: &[(&str, f64)]) -> BTreeMap<String, f64> {
    pairs
        .iter()
        .map(|(key, value)| ((*key).to_string(), *value))
        .collect()
}

fn string_map_from_pairs(pairs: &[(&str, &str)]) -> BTreeMap<String, String> {
    pairs
        .iter()
        .map(|(key, value)| ((*key).to_string(), (*value).to_string()))
        .collect()
}

#[cfg(test)]
mod tests {
    use super::*;
    use crate::events::{
        EventMeta, GMDecisionRegisteredEvent, MatchFinishedEvent, MatchTeam, PlayerBoxScore,
        SUBJECT_GM_DECISION_REGISTERED, SUBJECT_MATCH_FINISHED, SUBJECT_ROSTER_PATCH,
        SalaryCapCalculatedEvent, TradeAcceptedEvent, TradeProposedEvent,
    };

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
    fn individual_catalog_seeds_thirty_canon_agents() {
        let agents = default_individual_agent_states("game-1");

        assert_eq!(agents.len(), INDIVIDUAL_AGENT_COUNT);
        assert!(agents.iter().any(|agent| agent.agent_id == "owner"));
        assert!(agents.iter().any(|agent| agent.agent_id == "head_coach"));
        assert!(agents.iter().any(|agent| agent.agent_id == "cfo"));
        assert!(
            agents
                .iter()
                .any(|agent| agent.agent_id == "scouting_director")
        );
        assert!(
            agents
                .iter()
                .any(|agent| agent.agent_id == "sports_psychologist")
        );
        assert!(agents.iter().any(|agent| agent.agent_id == "mayor"));
        assert!(agents.iter().any(|agent| agent.agent_id == "press"));
    }

    #[test]
    fn individual_catalog_has_unique_ids_and_universal_state() {
        let agents = default_individual_agent_states("game-1");
        let mut ids = std::collections::BTreeSet::new();

        for agent in agents {
            assert!(ids.insert(agent.agent_id));
            assert!((-1.0..=1.0).contains(&agent.confidence));
            assert!((-1.0..=1.0).contains(&agent.satisfaction));
            assert!((0.0..=1.0).contains(&agent.loyalty));
            assert!((0.0..=1.0).contains(&agent.role_performance));
            assert!(!agent.domain.is_empty());
            assert!(!agent.agenda.is_empty());
            assert!(!agent.state.is_empty());
        }
    }

    #[test]
    fn rival_gm_catalog_seeds_thirty_profiles() {
        let rival_gms = default_rival_gms("game-1");

        assert_eq!(rival_gms.len(), RIVAL_GM_COUNT);
        assert!(rival_gms.iter().any(|gm| gm.rival_team_id == "bos"));
        assert!(rival_gms.iter().any(|gm| gm.gm_agent_id == "rival_gm_lal"));
        assert!(
            rival_gms
                .iter()
                .all(|gm| !gm.roster_needs.is_empty() && !gm.relationship_history.is_empty())
        );
    }

    #[test]
    fn rival_gm_catalog_is_deterministic_per_game() {
        let first = default_rival_gms("game-1");
        let second = default_rival_gms("game-1");
        let other_game = default_rival_gms("game-2");

        assert_eq!(first, second);
        assert_ne!(first[0].negotiation_style, "");
        assert!(
            first
                .iter()
                .all(|gm| (0.0..=1.0).contains(&gm.urgency_current))
        );
        assert!(
            first
                .iter()
                .all(|gm| (-1.0..=1.0).contains(&gm.relationship_trust))
        );
        assert_ne!(
            first
                .iter()
                .map(|gm| (&gm.negotiation_style, &gm.roster_needs))
                .collect::<Vec<_>>(),
            other_game
                .iter()
                .map(|gm| (&gm.negotiation_style, &gm.roster_needs))
                .collect::<Vec<_>>()
        );
    }

    #[test]
    fn rival_gm_rejects_trade_when_needs_do_not_fit() {
        let rival_gm = RivalGMProfile {
            game_id: "game-1".to_string(),
            rival_team_id: "bos".to_string(),
            gm_agent_id: "rival_gm_bos".to_string(),
            display_name: "Elliot Walsh".to_string(),
            team_name: "Boston Celtics".to_string(),
            negotiation_style: "patient_value_hunter".to_string(),
            urgency_current: 0.25,
            build_philosophy: "draft_and_develop".to_string(),
            roster_needs: vec!["C".to_string()],
            relationship_trust: -0.1,
            relationship_history: vec!["Sin historial".to_string()],
            last_interaction_event_id: None,
        };

        let evaluation = evaluate_trade_proposal(
            &rival_gm,
            &sample_trade_proposal("PG", 12_000_000),
            "now".to_string(),
        );

        match evaluation {
            RivalGMTradeEvaluation::Rejected(event) => {
                assert_eq!(event.reason, "rival_needs_mismatch");
                assert_eq!(event.proposal_id, "proposal-1");
            }
            RivalGMTradeEvaluation::Countered(_) => {
                panic!("expected rejected trade evaluation")
            }
        }
    }

    #[test]
    fn rival_gm_counters_trade_when_needs_fit() {
        let rival_gm = RivalGMProfile {
            game_id: "game-1".to_string(),
            rival_team_id: "bos".to_string(),
            gm_agent_id: "rival_gm_bos".to_string(),
            display_name: "Elliot Walsh".to_string(),
            team_name: "Boston Celtics".to_string(),
            negotiation_style: "asset_accumulator".to_string(),
            urgency_current: 0.5,
            build_philosophy: "draft_and_develop".to_string(),
            roster_needs: vec!["PG".to_string()],
            relationship_trust: 0.05,
            relationship_history: vec!["Sin historial".to_string()],
            last_interaction_event_id: None,
        };

        let evaluation = evaluate_trade_proposal(
            &rival_gm,
            &sample_trade_proposal("PG", 12_000_000),
            "now".to_string(),
        );

        match evaluation {
            RivalGMTradeEvaluation::Countered(event) => {
                assert_eq!(event.additional_asset_required, "second_round_pick");
                assert_eq!(event.rival_team_id, "bos");
            }
            RivalGMTradeEvaluation::Rejected(_) => {
                panic!("expected countered trade evaluation")
            }
        }
    }

    #[test]
    fn trade_acceptance_updates_outgoing_and_incoming_player_agents() {
        let outgoing =
            default_player_agent_state(&sample_roster_player("game-1-player-01", 82, "PG"));
        let event = sample_trade_accepted();

        let (states, patch) = apply_trade_accepted_to_player_agents(vec![outgoing], &event);

        let outgoing = states
            .iter()
            .find(|state| state.player_id == "game-1-player-01")
            .expect("outgoing player state exists");
        let incoming = states
            .iter()
            .find(|state| state.player_id == "trade-1-incoming")
            .expect("incoming player state exists");

        assert_eq!(outgoing.emotional_state, "traded");
        assert_eq!(incoming.emotional_state, "arriving");
        assert_eq!(incoming.full_name, "Jalen Warren");
        assert_eq!(patch.patch.players.len(), 2);
        assert_eq!(patch.patch.source_subject, "trade.aceptada");
    }

    #[test]
    fn default_player_agent_state_uses_team_player_id() {
        let player = sample_roster_player("game-1-player-01", 82, "PG");
        let state = default_player_agent_state(&player);

        assert_eq!(state.player_id, "game-1-player-01");
        assert_eq!(state.game_id, "game-1");
        assert_eq!(state.full_name, "Mateo Cross");
        assert!((0.0..=1.0).contains(&state.ego));
        assert!((0.0..=1.0).contains(&state.competitive_drive));
    }

    #[test]
    fn player_agents_react_to_box_score_with_roster_patch() {
        let mut event = sample_match(true, 112, 101);
        event.box_score = vec![PlayerBoxScore {
            player_id: "game-1-player-01".to_string(),
            team_id: OWN_TEAM_ID.to_string(),
            minutes: 32,
            points: 26,
            rebounds: 6,
            assists: 8,
            steals: 1,
            blocks: 0,
            turnovers: 2,
        }];
        let initial =
            default_player_agent_state(&sample_roster_player("game-1-player-01", 82, "PG"));

        let (states, patch) = apply_match_to_player_agents(vec![initial], &event);

        assert_eq!(states.len(), 1);
        assert_eq!(states[0].last_match_id.as_deref(), Some("match-1"));
        assert!(states[0].satisfaction > 0.04);

        let patch = patch.expect("roster patch exists");
        assert_eq!(patch.event_type, SUBJECT_ROSTER_PATCH);
        assert_eq!(patch.patch.players.len(), 1);
        assert_eq!(patch.patch.players[0].player_id, "game-1-player-01");
        assert_eq!(patch.patch.players[0].full_name, "Mateo Cross");
        assert_eq!(patch.patch.players[0].position, "PG");
    }

    #[test]
    fn relationship_catalog_seeds_canon_tensions() {
        let relationships = default_agent_relationships("game-1");

        assert!(
            relationships
                .iter()
                .any(|relationship| relationship.agent_a_id == "head_coach"
                    && relationship.agent_b_id == "team_doctor")
        );
        assert!(relationships.iter().any(
            |relationship| relationship.agent_a_id == "cfo" && relationship.agent_b_id == "gm"
        ));
        assert!(
            relationships
                .iter()
                .any(|relationship| relationship.agent_a_id == "press"
                    && relationship.agent_b_id == "roster_collective")
        );
        assert!(
            relationships
                .iter()
                .any(|relationship| relationship.agent_a_id == "gm"
                    && relationship.agent_b_id == "team_doctor")
        );
    }

    #[test]
    fn medical_decision_moves_doctor_gm_relationship() {
        let relationships = default_agent_relationships("game-1");
        let event = GMDecisionRegisteredEvent {
            meta: EventMeta {
                event_id: "decision-medical-injury-1".to_string(),
                game_id: "game-1".to_string(),
                occurred_at: "2026-10-29T00:00:00Z".to_string(),
                schema_version: SCHEMA_VERSION,
            },
            decision_id: "medical-injury-1".to_string(),
            kind: "medical_decision".to_string(),
            payload: string_map_from_pairs(&[
                ("choice_id", "force_return"),
                ("injury_id", "injury-1"),
                ("player_id", "player-1"),
            ]),
            simulated_date: "2026-10-29".to_string(),
            agents_affected: vec!["team_doctor".to_string(), "head_coach".to_string()],
            source_event_id: Some("injury-1".to_string()),
            source_subject: Some("jugador.lesionado".to_string()),
        };

        let changes = apply_gm_decision_to_relationships(
            relationships,
            &event,
            "2026-10-29T00:00:01Z".to_string(),
        );
        let doctor_gm = changes
            .iter()
            .find(|change| {
                relationship_key(
                    &change.relationship.agent_a_id,
                    &change.relationship.agent_b_id,
                ) == "gm:team_doctor"
            })
            .expect("doctor gm relationship moves");

        assert_eq!(doctor_gm.relationship.trend, "deteriorating");
        assert_eq!(
            doctor_gm.event.source_subject,
            SUBJECT_GM_DECISION_REGISTERED
        );
    }

    #[test]
    fn luxury_tax_moves_cfo_and_owner_state() {
        let states = vec![
            default_core_agent_state("game-1", "owner"),
            default_core_agent_state("game-1", "cfo"),
        ];
        let event = SalaryCapCalculatedEvent {
            meta: EventMeta {
                event_id: "salary-cap-game-1".to_string(),
                game_id: "game-1".to_string(),
                occurred_at: "2026-10-22T00:00:00Z".to_string(),
                schema_version: SCHEMA_VERSION,
            },
            simulated_date: "2026-10-22".to_string(),
            cap_base: 141_000_000,
            luxury_tax_line: 171_000_000,
            committed_salary: 180_000_000,
            cap_space: -39_000_000,
            luxury_tax_space: -9_000_000,
            roster_count: 15,
            status: "luxury_tax".to_string(),
            near_luxury_tax: true,
            projected_tax_payment: 18_000_000,
        };

        let changes =
            apply_salary_cap_to_core_agents(states, &event, "2026-10-22T00:00:01Z".to_string());
        let cfo = state_for(&changes, "cfo");
        let owner = state_for(&changes, "owner");

        assert!(cfo.state["budget_alert"] > 0.15);
        assert!(owner.state["patience_remaining"] < 0.75);
    }

    #[test]
    fn match_result_moves_relevant_relationships() {
        let event = sample_match(false, 88, 111);
        let relationships = default_agent_relationships("game-1");

        let changes =
            apply_match_to_relationships(relationships, &event, "2026-05-25T00:00:00Z".to_string());

        let coach_analytics = changes
            .iter()
            .find(|change| {
                relationship_key(
                    &change.relationship.agent_a_id,
                    &change.relationship.agent_b_id,
                ) == "head_analytics:head_coach"
            })
            .expect("coach analytics relationship moves");
        assert_eq!(coach_analytics.relationship.trend, "deteriorating");
        assert_eq!(coach_analytics.event.source_subject, SUBJECT_MATCH_FINISHED);
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
            box_score: Vec::new(),
        }
    }

    fn sample_roster_player(
        player_id: &str,
        overall_rating: u8,
        position: &str,
    ) -> TeamRosterPlayer {
        TeamRosterPlayer {
            player_id: player_id.to_string(),
            game_id: "game-1".to_string(),
            full_name: "Mateo Cross".to_string(),
            position: position.to_string(),
            overall_rating,
            roster_status: "active".to_string(),
        }
    }

    fn sample_trade_proposal(requested_position: &str, incoming_salary: i64) -> TradeProposedEvent {
        TradeProposedEvent {
            meta: EventMeta {
                event_id: "trade-proposed-proposal-1".to_string(),
                game_id: "game-1".to_string(),
                occurred_at: "2026-11-01T00:00:00Z".to_string(),
                schema_version: 1,
            },
            proposal_id: "proposal-1".to_string(),
            simulated_date: "2026-11-01".to_string(),
            rival_team_id: "bos".to_string(),
            offered_player_id: "player-1".to_string(),
            offered_player_name: "Mateo Cross".to_string(),
            offered_salary: 10_000_000,
            requested_position: requested_position.to_string(),
            incoming_salary,
            cap_space_after: -12_000_000,
        }
    }

    fn sample_trade_accepted() -> TradeAcceptedEvent {
        TradeAcceptedEvent {
            meta: EventMeta {
                event_id: "trade-accepted-trade-1".to_string(),
                game_id: "game-1".to_string(),
                occurred_at: "2026-11-01T00:00:00Z".to_string(),
                schema_version: 1,
            },
            proposal_id: "trade-1".to_string(),
            simulated_date: "2026-11-01".to_string(),
            rival_team_id: "bos".to_string(),
            outgoing_player_id: "game-1-player-01".to_string(),
            outgoing_player_name: "Mateo Cross".to_string(),
            incoming_player_id: "trade-1-incoming".to_string(),
            incoming_player_name: "Jalen Warren".to_string(),
            incoming_position: "PG".to_string(),
            incoming_rating: 80,
            incoming_salary: 12_000_000,
            accepted_additional_asset: Some("second_round_pick".to_string()),
        }
    }
}
