//! Rival GM evaluation and emotional reactions to accepted trades.

use super::{
    PlayerAgentState, RivalGMProfile, SCHEMA_VERSION, TeamRosterPlayer, adjusted, clamp_unit,
    default_player_agent_state,
};
use crate::events::{
    EventMeta, PlayerEmotionalPatch, RosterPatchEnvelope, RosterStatePatch, SUBJECT_ROSTER_PATCH,
    TradeAcceptedEvent, TradeCounteredEvent, TradeProposedEvent, TradeRejectedEvent,
};

#[derive(Debug, Clone, PartialEq, Eq)]
pub enum RivalGMTradeEvaluation {
    Rejected(TradeRejectedEvent),
    Countered(TradeCounteredEvent),
}

#[must_use]
pub fn evaluate_trade_proposal(
    rival_gm: &RivalGMProfile,
    event: &TradeProposedEvent,
    occurred_at: String,
) -> RivalGMTradeEvaluation {
    let need_fit = rival_gm
        .roster_needs
        .iter()
        .any(|need| need.eq_ignore_ascii_case(&event.requested_position));
    let salary_pressure = event.incoming_salary.saturating_sub(event.offered_salary);
    let leverage_score = rival_gm.urgency_current + rival_gm.relationship_trust;
    let style = rival_gm.negotiation_style.as_str();

    if !need_fit && leverage_score < 0.58 {
        return RivalGMTradeEvaluation::Rejected(TradeRejectedEvent {
            meta: trade_meta(
                format!("trade-rejected-{}", event.proposal_id),
                event,
                occurred_at,
            ),
            proposal_id: event.proposal_id.clone(),
            simulated_date: event.simulated_date.clone(),
            rival_team_id: event.rival_team_id.clone(),
            reason: "rival_needs_mismatch".to_string(),
            detail: format!(
                "{} no ve encaje claro con sus necesidades actuales.",
                rival_gm.display_name
            ),
        });
    }

    if salary_pressure > 8_000_000 && style != "win_now_pressure" {
        return RivalGMTradeEvaluation::Rejected(TradeRejectedEvent {
            meta: trade_meta(
                format!("trade-rejected-{}", event.proposal_id),
                event,
                occurred_at,
            ),
            proposal_id: event.proposal_id.clone(),
            simulated_date: event.simulated_date.clone(),
            rival_team_id: event.rival_team_id.clone(),
            reason: "salary_value_gap".to_string(),
            detail: format!(
                "{} rechaza absorber tanta diferencia salarial sin mas valor.",
                rival_gm.display_name
            ),
        });
    }

    let additional_asset_required = match style {
        "asset_accumulator" => "second_round_pick",
        "aggressive_star_chaser" => "rotation_player",
        "cap_flexible_operator" => "salary_relief",
        "defensive_conservative" => "defensive_wing",
        "win_now_pressure" => "veteran_depth",
        _ => "future_second",
    };

    RivalGMTradeEvaluation::Countered(TradeCounteredEvent {
        meta: trade_meta(
            format!("trade-countered-{}", event.proposal_id),
            event,
            occurred_at,
        ),
        proposal_id: event.proposal_id.clone(),
        simulated_date: event.simulated_date.clone(),
        rival_team_id: event.rival_team_id.clone(),
        requested_position: event.requested_position.clone(),
        additional_asset_required: additional_asset_required.to_string(),
        detail: format!(
            "{} no acepta el paquete inicial, pero deja abierta una contraoferta.",
            rival_gm.display_name
        ),
    })
}

#[must_use]
pub fn apply_trade_accepted_to_player_agents(
    current_states: Vec<PlayerAgentState>,
    event: &TradeAcceptedEvent,
) -> (Vec<PlayerAgentState>, RosterPatchEnvelope) {
    let mut updated_states = Vec::with_capacity(current_states.len() + 1);
    let mut incoming_found = false;
    let mut patches = Vec::with_capacity(2);

    for mut state in current_states {
        if state.player_id == event.outgoing_player_id {
            state.emotional_state = "traded".to_string();
            state.satisfaction = adjusted(state.satisfaction, -0.08);
            state.loyalty = clamp_unit(state.loyalty - 0.12);
            state.city_connection = clamp_unit(state.city_connection - 0.18);
            patches.push(PlayerEmotionalPatch {
                player_id: state.player_id.clone(),
                full_name: state.full_name.clone(),
                position: state.position.clone(),
                emotional_state: state.emotional_state.clone(),
                satisfaction: state.satisfaction,
                loyalty: state.loyalty,
                ego: state.ego,
                competitive_drive: state.competitive_drive,
                city_connection: state.city_connection,
                summary: format!("{} procesa su salida via trade.", state.full_name),
            });
        }
        if state.player_id == event.incoming_player_id {
            incoming_found = true;
            state.emotional_state = "arriving".to_string();
            state.satisfaction = adjusted(state.satisfaction, 0.03);
            patches.push(PlayerEmotionalPatch {
                player_id: state.player_id.clone(),
                full_name: state.full_name.clone(),
                position: state.position.clone(),
                emotional_state: state.emotional_state.clone(),
                satisfaction: state.satisfaction,
                loyalty: state.loyalty,
                ego: state.ego,
                competitive_drive: state.competitive_drive,
                city_connection: state.city_connection,
                summary: format!("{} llega a PulseCity via trade.", state.full_name),
            });
        }
        updated_states.push(state);
    }

    if !incoming_found {
        let player = TeamRosterPlayer {
            player_id: event.incoming_player_id.clone(),
            game_id: event.meta.game_id.clone(),
            full_name: event.incoming_player_name.clone(),
            position: event.incoming_position.clone(),
            overall_rating: event.incoming_rating,
            roster_status: "active".to_string(),
        };
        let mut state = default_player_agent_state(&player);
        state.emotional_state = "arriving".to_string();
        state.city_connection = 0.18;
        patches.push(PlayerEmotionalPatch {
            player_id: state.player_id.clone(),
            full_name: state.full_name.clone(),
            position: state.position.clone(),
            emotional_state: state.emotional_state.clone(),
            satisfaction: state.satisfaction,
            loyalty: state.loyalty,
            ego: state.ego,
            competitive_drive: state.competitive_drive,
            city_connection: state.city_connection,
            summary: format!("{} llega a PulseCity via trade.", state.full_name),
        });
        updated_states.push(state);
    }

    (
        updated_states,
        RosterPatchEnvelope {
            event_type: SUBJECT_ROSTER_PATCH.to_string(),
            subject: SUBJECT_ROSTER_PATCH.to_string(),
            game_id: event.meta.game_id.clone(),
            patch: RosterStatePatch {
                simulated_date: event.simulated_date.clone(),
                source_event_id: event.meta.event_id.clone(),
                source_subject: "trade.aceptada".to_string(),
                players: patches,
            },
        },
    )
}

fn trade_meta(event_id: String, event: &TradeProposedEvent, occurred_at: String) -> EventMeta {
    EventMeta {
        event_id,
        game_id: event.meta.game_id.clone(),
        occurred_at,
        schema_version: SCHEMA_VERSION,
    }
}
