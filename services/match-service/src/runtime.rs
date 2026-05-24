use anyhow::{Context, Result};
use async_nats::Client;
use tokio_stream::StreamExt;
use tracing::{error, info};

use crate::{
    events::{
        EventMeta, MatchScheduledEvent, MatchStartingEvent, SUBJECT_MATCH_FINISHED,
        SUBJECT_MATCH_SCHEDULED, SUBJECT_MATCH_STARTING,
    },
    simulator::{MatchSimulationInput, simulate_match},
};

const SCHEMA_VERSION: u16 = 1;

pub async fn run(client: Client) -> Result<()> {
    let mut scheduled = client
        .subscribe(SUBJECT_MATCH_SCHEDULED)
        .await
        .context("subscribe partido.programado")?;

    info!("match-service listening for {SUBJECT_MATCH_SCHEDULED}");

    loop {
        tokio::select! {
            maybe_message = scheduled.next() => {
                let Some(message) = maybe_message else {
                    break;
                };
                if let Err(err) = handle_scheduled_match(&client, &message.payload).await {
                    error!("handle partido.programado: {err:#}");
                }
            }
            _ = tokio::signal::ctrl_c() => {
                info!("match-service shutdown signal received");
                break;
            }
        }
    }

    Ok(())
}

async fn handle_scheduled_match(client: &Client, payload: &[u8]) -> Result<()> {
    let event: MatchScheduledEvent =
        serde_json::from_slice(payload).context("decode partido.programado")?;
    let input = MatchSimulationInput::from(event.clone());

    let starting = MatchStartingEvent {
        meta: EventMeta {
            event_id: format!("match-starting-{}", event.match_id),
            game_id: event.meta.game_id,
            occurred_at: format!("{}T00:00:00Z", event.simulated_date),
            schema_version: SCHEMA_VERSION,
        },
        match_id: event.match_id,
        simulated_date: event.simulated_date,
    };
    publish_json(client, SUBJECT_MATCH_STARTING, &starting)
        .await
        .context("publish partido.iniciando")?;

    let finished = simulate_match(&input).context("simulate match")?;
    publish_json(client, SUBJECT_MATCH_FINISHED, &finished)
        .await
        .context("publish partido.terminado")?;

    Ok(())
}

async fn publish_json<T>(client: &Client, subject: &str, payload: &T) -> Result<()>
where
    T: serde::Serialize,
{
    let encoded = serde_json::to_vec(payload).context("encode nats payload")?;
    client
        .publish(subject.to_string(), encoded.into())
        .await
        .with_context(|| format!("publish {subject}"))?;
    Ok(())
}
