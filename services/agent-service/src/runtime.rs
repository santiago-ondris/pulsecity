use std::time::{Duration, SystemTime, UNIX_EPOCH};

use anyhow::{Context, Result};
use async_nats::Client;
use futures_util::StreamExt;
use tokio::time::{MissedTickBehavior, interval};
use tracing::{debug, error, info, warn};

use crate::{
    events::{
        DayAdvancedEvent, EventMeta, PauseChangedEvent, SUBJECT_TIME_DAY_ADVANCED,
        SUBJECT_TIME_PAUSE_CHANGED, SUBJECT_TIME_SESSION_ENDED, SUBJECT_TIME_SESSION_STARTED,
        SUBJECT_TIME_SPEED_CHANGED, SessionEndedEvent, SessionStartedEvent, SpeedChangedEvent,
    },
    persistence::Store,
    simulation::{SimulationAccumulator, SimulationState, advance_simulated_date},
};

const TICK_INTERVAL: Duration = Duration::from_millis(100);
const SCHEMA_VERSION: u16 = 1;

pub async fn run(client: Client, store: Store, mut state: SimulationState) -> Result<()> {
    let mut session_started = client
        .subscribe(SUBJECT_TIME_SESSION_STARTED)
        .await
        .context("subscribe tiempo.sesion_iniciada")?;
    let mut session_ended = client
        .subscribe(SUBJECT_TIME_SESSION_ENDED)
        .await
        .context("subscribe tiempo.sesion_terminada")?;
    let mut speed_changed = client
        .subscribe(SUBJECT_TIME_SPEED_CHANGED)
        .await
        .context("subscribe tiempo.velocidad_cambiada")?;
    let mut pause_changed = client
        .subscribe(SUBJECT_TIME_PAUSE_CHANGED)
        .await
        .context("subscribe tiempo.pausa_activada")?;

    let mut ticker = interval(TICK_INTERVAL);
    ticker.set_missed_tick_behavior(MissedTickBehavior::Delay);
    let mut accumulator = SimulationAccumulator::new();

    info!(
        game_id = %state.game_id,
        "simulation loop started"
    );

    loop {
        tokio::select! {
            _ = tokio::signal::ctrl_c() => {
                info!("shutdown signal received");
                break;
            }
            _ = ticker.tick() => {
                process_tick(&client, &store, &mut state, &mut accumulator).await?;
            }
            Some(message) = session_started.next() => {
                let event = match decode_event::<SessionStartedEvent>(SUBJECT_TIME_SESSION_STARTED, &message.payload) {
                    Some(event) => event,
                    None => continue,
                };
                if event.meta.game_id != state.game_id {
                    continue;
                }

                state.set_session_active(true);
                store.save_simulation_state(&state).await?;
                info!(game_id = %state.game_id, session_id = %event.session_id, "simulation session activated");
            }
            Some(message) = session_ended.next() => {
                let event = match decode_event::<SessionEndedEvent>(SUBJECT_TIME_SESSION_ENDED, &message.payload) {
                    Some(event) => event,
                    None => continue,
                };
                if event.meta.game_id != state.game_id {
                    continue;
                }

                state.set_session_active(false);
                accumulator.reset();
                store.save_simulation_state(&state).await?;
                info!(game_id = %state.game_id, session_id = %event.session_id, "simulation session deactivated");
            }
            Some(message) = speed_changed.next() => {
                let event = match decode_event::<SpeedChangedEvent>(SUBJECT_TIME_SPEED_CHANGED, &message.payload) {
                    Some(event) => event,
                    None => continue,
                };
                if event.meta.game_id != state.game_id {
                    continue;
                }

                if let Err(err) = state.set_speed(event.speed) {
                    warn!(game_id = %state.game_id, speed = event.speed, error = %err, "ignoring invalid simulation speed");
                    continue;
                }

                store.save_simulation_state(&state).await?;
                info!(game_id = %state.game_id, speed = state.speed, "simulation speed changed");
            }
            Some(message) = pause_changed.next() => {
                let event = match decode_event::<PauseChangedEvent>(SUBJECT_TIME_PAUSE_CHANGED, &message.payload) {
                    Some(event) => event,
                    None => continue,
                };
                if event.meta.game_id != state.game_id {
                    continue;
                }

                state.set_paused(event.paused);
                if state.paused {
                    accumulator.reset();
                }
                store.save_simulation_state(&state).await?;
                info!(game_id = %state.game_id, paused = state.paused, "simulation pause changed");
            }
        }
    }

    Ok(())
}

async fn process_tick(
    client: &Client,
    store: &Store,
    state: &mut SimulationState,
    accumulator: &mut SimulationAccumulator,
) -> Result<()> {
    let advance = accumulator.tick(state, TICK_INTERVAL.as_millis());
    if advance.days_processed == 0 {
        return Ok(());
    }

    let processed_date =
        advance_simulated_date(&state.current_simulated_date, advance.days_processed)
            .with_context(|| {
                format!(
                    "advance simulated date from {}",
                    state.current_simulated_date
                )
            })?;

    state.current_simulated_date.clone_from(&processed_date);
    state.last_tick_processed_at = Some(now_rfc3339());
    store.save_simulation_state(state).await?;

    let event = DayAdvancedEvent {
        meta: EventMeta {
            event_id: format!("time-day-{}-{processed_date}", state.game_id),
            game_id: state.game_id.clone(),
            occurred_at: state
                .last_tick_processed_at
                .clone()
                .unwrap_or_else(now_rfc3339),
            schema_version: SCHEMA_VERSION,
        },
        simulated_date: processed_date,
        speed: state.speed,
        days_processed: advance.days_processed,
    };

    let payload = serde_json::to_vec(&event).context("encode tiempo.dia_avanzado")?;
    client
        .publish(SUBJECT_TIME_DAY_ADVANCED, payload.into())
        .await
        .context("publish tiempo.dia_avanzado")?;

    debug!(
        game_id = %state.game_id,
        simulated_date = %event.simulated_date,
        days_processed = event.days_processed,
        "simulation day advanced"
    );

    Ok(())
}

fn decode_event<T>(subject: &str, payload: &[u8]) -> Option<T>
where
    T: serde::de::DeserializeOwned,
{
    match serde_json::from_slice(payload) {
        Ok(event) => Some(event),
        Err(err) => {
            error!(subject, error = %err, "failed to decode nats event");
            None
        }
    }
}

fn now_rfc3339() -> String {
    let elapsed = SystemTime::now()
        .duration_since(UNIX_EPOCH)
        .unwrap_or_else(|_| Duration::from_secs(0));
    let seconds = elapsed.as_secs() as i64;
    let nanos = elapsed.subsec_nanos();
    let days = seconds.div_euclid(86_400);
    let second_of_day = seconds.rem_euclid(86_400);
    let (year, month, day) = civil_from_unix_days(days);
    let hour = second_of_day / 3_600;
    let minute = (second_of_day % 3_600) / 60;
    let second = second_of_day % 60;

    format!("{year:04}-{month:02}-{day:02}T{hour:02}:{minute:02}:{second:02}.{nanos:09}Z")
}

fn civil_from_unix_days(days: i64) -> (i32, u32, u32) {
    let days = days + 719_468;
    let era = if days >= 0 { days } else { days - 146_096 } / 146_097;
    let day_of_era = days - era * 146_097;
    let year_of_era =
        (day_of_era - day_of_era / 1_460 + day_of_era / 36_524 - day_of_era / 146_096) / 365;
    let year = year_of_era + era * 400;
    let day_of_year = day_of_era - (365 * year_of_era + year_of_era / 4 - year_of_era / 100);
    let month_prime = (5 * day_of_year + 2) / 153;
    let day = day_of_year - (153 * month_prime + 2) / 5 + 1;
    let month = month_prime + if month_prime < 10 { 3 } else { -9 };
    let year = year + i64::from(month <= 2);

    (year as i32, month as u32, day as u32)
}

#[cfg(test)]
mod tests {
    use super::*;

    #[test]
    fn rfc3339_timestamp_has_utc_suffix() {
        let timestamp = now_rfc3339();

        assert!(timestamp.ends_with('Z'));
        assert!(timestamp.contains('T'));
    }
}
