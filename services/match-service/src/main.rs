use anyhow::Context;
use match_service::{SERVICE_NAME, nats_url_from_env, runtime};
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

    info!(service = SERVICE_NAME, nats_url, "connected to nats");
    runtime::run(client.clone())
        .await
        .context("run match-service runtime")?;

    drop(client);
    Ok(())
}
