//! Emotional and relational reactions triggered by completed matches.

use super::{
    AgentRelationship, AgentRelationshipChange, AgentStateChange, CoreAgentState, OWN_TEAM_ID,
    PlayerAgentState, SCHEMA_VERSION, adjust, adjusted, apply_relationship_delta, clamp,
    clamp_unit, metric, relationship_key,
};
use crate::events::{
    AgentRelationshipChangedEvent, AgentStateChangedEvent, EventMeta, MatchFinishedEvent,
    PlayerBoxScore, PlayerEmotionalPatch, RosterPatchEnvelope, RosterStatePatch,
    SUBJECT_MATCH_FINISHED, SUBJECT_ROSTER_PATCH,
};

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

#[must_use]
pub fn apply_match_to_player_agents(
    current_states: Vec<PlayerAgentState>,
    event: &MatchFinishedEvent,
) -> (Vec<PlayerAgentState>, Option<RosterPatchEnvelope>) {
    if current_states.is_empty() {
        return (Vec::new(), None);
    }

    let context = MatchContext::from_event(event);
    let mut next_states = Vec::with_capacity(current_states.len());
    let mut patches = Vec::with_capacity(current_states.len());

    for mut state in current_states {
        let line = event
            .box_score
            .iter()
            .find(|line| line.player_id == state.player_id && line.team_id == OWN_TEAM_ID);
        let Some(line) = line else {
            next_states.push(state);
            continue;
        };

        apply_box_score_to_player(&mut state, line, &context, &event.match_id);
        patches.push(player_patch(&state, line, &context));
        next_states.push(state);
    }

    let roster_patch = if patches.is_empty() {
        None
    } else {
        Some(RosterPatchEnvelope {
            event_type: SUBJECT_ROSTER_PATCH.to_string(),
            subject: SUBJECT_ROSTER_PATCH.to_string(),
            game_id: event.meta.game_id.clone(),
            patch: RosterStatePatch {
                simulated_date: event.simulated_date.clone(),
                source_event_id: event.meta.event_id.clone(),
                source_subject: SUBJECT_MATCH_FINISHED.to_string(),
                players: patches,
            },
        })
    };

    (next_states, roster_patch)
}

#[must_use]
pub fn apply_match_to_relationships(
    current_relationships: Vec<AgentRelationship>,
    event: &MatchFinishedEvent,
    occurred_at: String,
) -> Vec<AgentRelationshipChange> {
    let context = MatchContext::from_event(event);
    let mut changes = Vec::new();

    for mut relationship in current_relationships {
        let Some((delta, reason)) = relationship_delta_for_match(&relationship, &context) else {
            continue;
        };
        apply_relationship_delta(&mut relationship, delta, reason, &event.meta.event_id);

        changes.push(AgentRelationshipChange {
            event: relationship_changed_event(&relationship, event, &occurred_at),
            relationship,
        });
    }

    changes
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

fn apply_box_score_to_player(
    state: &mut PlayerAgentState,
    line: &PlayerBoxScore,
    context: &MatchContext,
    match_id: &str,
) {
    let performance = player_performance_index(line);
    let result_delta = if context.won { 1.0 } else { -1.0 };
    let role_pressure = if line.minutes >= 28 { 0.018 } else { -0.01 };

    state.satisfaction = clamp(adjusted(
        state.satisfaction,
        0.035 * result_delta + performance * 0.035,
    ));
    state.loyalty = clamp_unit(state.loyalty + 0.012 * result_delta + performance * 0.01);
    state.ego = clamp_unit(state.ego + performance * 0.035 + role_pressure);
    state.competitive_drive =
        clamp_unit(state.competitive_drive + if context.won { 0.008 } else { 0.02 });
    state.city_connection = clamp_unit(
        state.city_connection
            + if context.home_game && context.won {
                0.018
            } else {
                0.004
            },
    );
    state.last_match_id = Some(match_id.to_string());

    state.emotional_state = if !context.won && performance < -0.35 {
        "frustrated"
    } else if !context.won {
        "restless"
    } else if performance > 0.40 {
        "confident"
    } else {
        "steady"
    }
    .to_string();
}

fn relationship_delta_for_match(
    relationship: &AgentRelationship,
    context: &MatchContext,
) -> Option<(f64, &'static str)> {
    let key = relationship_key(&relationship.agent_a_id, &relationship.agent_b_id);
    match key.as_str() {
        "head_analytics:head_coach" => {
            if context.won {
                Some((
                    0.015,
                    "El resultado reduce la friccion entre modelo y rotacion.",
                ))
            } else if context.blowout {
                Some((
                    -0.035,
                    "La derrota amplia reabre la tension entre datos y decisiones de cancha.",
                ))
            } else {
                Some((
                    -0.012,
                    "La derrota deja mas preguntas sobre la lectura tactica.",
                ))
            }
        }
        "head_coach:sports_psychologist" => {
            if context.won {
                Some((
                    0.012,
                    "La victoria baja la tension sobre el manejo emocional del vestuario.",
                ))
            } else {
                Some((
                    -0.025,
                    "La derrota aumenta la tension entre bienestar y exigencia competitiva.",
                ))
            }
        }
        "mayor:owner" => {
            if context.home_game && context.won {
                Some((
                    0.012,
                    "Una victoria local mejora la lectura civica del proyecto.",
                ))
            } else if context.home_game {
                Some((
                    -0.018,
                    "Una derrota local enfria el valor politico del proyecto.",
                ))
            } else {
                None
            }
        }
        "gm:pr_director" => {
            if context.blowout && !context.won {
                Some((
                    -0.025,
                    "La derrota amplia complica la narrativa publica del GM.",
                ))
            } else if context.won {
                Some((
                    0.01,
                    "La victoria hace mas defendible la direccion publica del GM.",
                ))
            } else {
                None
            }
        }
        "press:roster_collective" => {
            if context.blowout && !context.won {
                Some((
                    -0.03,
                    "La cobertura se endurece sobre el estado emocional del roster.",
                ))
            } else if context.won {
                Some((
                    0.015,
                    "La cobertura positiva reduce presion sobre el vestuario.",
                ))
            } else {
                Some((
                    -0.01,
                    "La derrota sostiene una cobertura mas incomoda para el roster.",
                ))
            }
        }
        _ => None,
    }
}

fn relationship_changed_event(
    relationship: &AgentRelationship,
    event: &MatchFinishedEvent,
    occurred_at: &str,
) -> AgentRelationshipChangedEvent {
    AgentRelationshipChangedEvent {
        meta: EventMeta {
            event_id: format!(
                "agent-relationship-{}-{}",
                event.meta.event_id,
                relationship_key(&relationship.agent_a_id, &relationship.agent_b_id)
            ),
            game_id: event.meta.game_id.clone(),
            occurred_at: occurred_at.to_string(),
            schema_version: SCHEMA_VERSION,
        },
        simulated_date: event.simulated_date.clone(),
        agent_a_id: relationship.agent_a_id.clone(),
        agent_b_id: relationship.agent_b_id.clone(),
        trust: relationship.trust,
        trend: relationship.trend.clone(),
        last_event: relationship.last_event.clone(),
        short_history: relationship.short_history.clone(),
        source_event_id: event.meta.event_id.clone(),
        source_subject: SUBJECT_MATCH_FINISHED.to_string(),
    }
}

fn player_patch(
    state: &PlayerAgentState,
    line: &PlayerBoxScore,
    context: &MatchContext,
) -> PlayerEmotionalPatch {
    PlayerEmotionalPatch {
        player_id: state.player_id.clone(),
        full_name: state.full_name.clone(),
        position: state.position.clone(),
        emotional_state: state.emotional_state.clone(),
        satisfaction: state.satisfaction,
        loyalty: state.loyalty,
        ego: state.ego,
        competitive_drive: state.competitive_drive,
        city_connection: state.city_connection,
        summary: summarize_player_change(state, line, context),
    }
}

fn summarize_player_change(
    state: &PlayerAgentState,
    line: &PlayerBoxScore,
    context: &MatchContext,
) -> String {
    let result = if context.won { "victoria" } else { "derrota" };
    let role = if line.minutes >= 28 {
        "rol alto"
    } else if line.minutes >= 16 {
        "rotacion estable"
    } else {
        "minutos limitados"
    };

    format!(
        "{} procesa la {} con {} y {} puntos.",
        state.full_name, result, role, line.points
    )
}

fn player_performance_index(line: &PlayerBoxScore) -> f64 {
    let production = f64::from(line.points) * 0.04
        + f64::from(line.rebounds) * 0.025
        + f64::from(line.assists) * 0.03
        + f64::from(line.steals + line.blocks) * 0.04
        - f64::from(line.turnovers) * 0.04;
    let minutes_expectation = f64::from(line.minutes) * 0.018;

    (production - minutes_expectation).clamp(-1.0, 1.0)
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
