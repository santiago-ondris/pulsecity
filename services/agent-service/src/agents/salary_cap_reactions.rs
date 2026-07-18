//! Core agent reactions triggered by salary cap calculations.

use super::{AgentStateChange, CoreAgentState, SCHEMA_VERSION, adjust};
use crate::events::{AgentStateChangedEvent, EventMeta, SalaryCapCalculatedEvent};

#[must_use]
pub fn apply_salary_cap_to_core_agents(
    current_states: Vec<CoreAgentState>,
    event: &SalaryCapCalculatedEvent,
    occurred_at: String,
) -> Vec<AgentStateChange> {
    let mut changes = Vec::new();
    for mut state in current_states {
        match state.agent_id.as_str() {
            "cfo" => {
                if event.status == "luxury_tax" {
                    adjust(&mut state.state, "budget_alert", 0.18);
                    adjust(&mut state.state, "financial_trust", -0.04);
                    state.mood = "concerned".to_string();
                } else if event.near_luxury_tax {
                    adjust(&mut state.state, "budget_alert", 0.08);
                    state.mood = "watchful".to_string();
                } else {
                    adjust(&mut state.state, "budget_alert", -0.02);
                    state.mood = "calm".to_string();
                }
            }
            "owner" => {
                if event.status == "luxury_tax" {
                    adjust(&mut state.state, "patience_remaining", -0.04);
                    adjust(&mut state.state, "business_trust", -0.03);
                    state.mood = "concerned".to_string();
                } else {
                    continue;
                }
            }
            _ => continue,
        }

        changes.push(AgentStateChange {
            event: AgentStateChangedEvent {
                meta: EventMeta {
                    event_id: format!("agent-state-{}-{}", event.meta.event_id, state.agent_id),
                    game_id: event.meta.game_id.clone(),
                    occurred_at: occurred_at.clone(),
                    schema_version: SCHEMA_VERSION,
                },
                simulated_date: event.simulated_date.clone(),
                agent_id: state.agent_id.clone(),
                mood: state.mood.clone(),
                state: state.state.clone(),
                summary: salary_cap_summary(&state.agent_id, event),
                source_event_id: event.meta.event_id.clone(),
                source_subject: "salary_cap.calculado".to_string(),
            },
            state,
        });
    }

    changes
}

fn salary_cap_summary(agent_id: &str, event: &SalaryCapCalculatedEvent) -> String {
    match (agent_id, event.status.as_str()) {
        ("cfo", "luxury_tax") => "El CFO alerta que la nomina entro en luxury tax.".to_string(),
        ("cfo", _) if event.near_luxury_tax => {
            "El CFO marca que la franquicia esta cerca de la linea de luxury tax.".to_string()
        }
        ("owner", "luxury_tax") => {
            "El Owner reduce paciencia por el costo proyectado de la nomina.".to_string()
        }
        _ => "La situacion financiera se mantiene controlada.".to_string(),
    }
}
