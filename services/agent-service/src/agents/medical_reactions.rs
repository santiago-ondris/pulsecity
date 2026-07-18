//! Relationship reactions triggered by the GM's medical decisions.

use super::{
    AgentRelationship, AgentRelationshipChange, SCHEMA_VERSION, apply_relationship_delta,
    relationship_key,
};
use crate::events::{
    AgentRelationshipChangedEvent, EventMeta, GMDecisionRegisteredEvent,
    SUBJECT_GM_DECISION_REGISTERED,
};

#[must_use]
pub fn apply_gm_decision_to_relationships(
    current_relationships: Vec<AgentRelationship>,
    event: &GMDecisionRegisteredEvent,
    occurred_at: String,
) -> Vec<AgentRelationshipChange> {
    if event.kind != "medical_decision" {
        return Vec::new();
    }

    let Some(choice_id) = event.payload.get("choice_id") else {
        return Vec::new();
    };
    let mut changes = Vec::new();
    for mut relationship in current_relationships {
        let Some((delta, reason)) =
            relationship_delta_for_medical_decision(&relationship, choice_id)
        else {
            continue;
        };
        apply_relationship_delta(&mut relationship, delta, reason, &event.meta.event_id);
        changes.push(AgentRelationshipChange {
            event: relationship_changed_event_from_decision(&relationship, event, &occurred_at),
            relationship,
        });
    }

    changes
}

fn relationship_delta_for_medical_decision(
    relationship: &AgentRelationship,
    choice_id: &str,
) -> Option<(f64, &'static str)> {
    let key = relationship_key(&relationship.agent_a_id, &relationship.agent_b_id);
    match (key.as_str(), choice_id) {
        ("gm:team_doctor", "rest") => Some((
            0.045,
            "El GM respeta el protocolo medico y fortalece la confianza del staff de salud.",
        )),
        ("gm:team_doctor", "reduce_minutes") => Some((
            0.025,
            "El GM acepta bajar carga y el Medico interpreta la decision como prudente.",
        )),
        ("gm:team_doctor", "ignore_doctor") => Some((
            -0.055,
            "El GM ignora la recomendacion medica y erosiona la confianza del staff de salud.",
        )),
        ("gm:team_doctor", "force_return") => Some((
            -0.09,
            "El GM fuerza una alta anticipada y abre una fractura seria con el Medico.",
        )),
        ("head_coach:team_doctor", "force_return") => Some((
            -0.035,
            "La urgencia competitiva vuelve a tensionar al Coach con el Medico.",
        )),
        ("head_coach:team_doctor", "rest") => Some((
            0.018,
            "Coach y Medico quedan alineados alrededor del protocolo de recuperacion.",
        )),
        ("head_coach:team_doctor", "reduce_minutes") => Some((
            0.012,
            "La reduccion de carga crea un compromiso aceptable entre competencia y salud.",
        )),
        _ => None,
    }
}

fn relationship_changed_event_from_decision(
    relationship: &AgentRelationship,
    event: &GMDecisionRegisteredEvent,
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
        source_subject: SUBJECT_GM_DECISION_REGISTERED.to_string(),
    }
}
