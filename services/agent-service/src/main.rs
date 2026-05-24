use agent_service::{
    SERVICE_NAME, database_url_from_env, game_id_from_env, nats_url_from_env, persistence::Store,
    runtime,
};
use anyhow::Context;
use tracing::info;

#[tokio::main]
async fn main() -> anyhow::Result<()> {
    tracing_subscriber::fmt()
        .with_env_filter(
            tracing_subscriber::EnvFilter::try_from_default_env().unwrap_or_else(|_| "info".into()),
        )
        .init();

    let nats_url = nats_url_from_env();
    let client = async_nats::connect(&nats_url)
        .await
        .with_context(|| format!("connect nats at {nats_url}"))?;

    let database_url = database_url_from_env();
    let store = Store::connect(&database_url).await?;
    store.ensure_schema().await?;

    let game_id = game_id_from_env();
    let state = store.load_or_initialize_simulation_state(&game_id).await?;

    info!(service = SERVICE_NAME, nats_url, "connected to nats");
    info!(
        service = SERVICE_NAME,
        game_id = %state.game_id,
        simulated_date = %state.current_simulated_date,
        speed = state.speed,
        paused = state.paused,
        session_active = state.session_active,
        "simulation state loaded"
    );

    runtime::run(client.clone(), store, state)
        .await
        .context("run simulation loop")?;

    drop(client);
    Ok(())
}
