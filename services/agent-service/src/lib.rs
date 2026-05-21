//! Core scaffolding for PulseCity's agent-service.

pub mod events;
pub mod persistence;
pub mod simulation;

pub const SERVICE_NAME: &str = "agent-service";
pub const DEFAULT_NATS_URL: &str = "nats://127.0.0.1:4222";
pub const DEFAULT_DATABASE_URL: &str =
    "postgres://pulsecity:pulsecity@localhost:5433/pulsecity_dev?sslmode=disable";

#[must_use]
pub fn nats_url_from_env() -> String {
    std::env::var("NATS_URL").unwrap_or_else(|_| DEFAULT_NATS_URL.to_string())
}

#[must_use]
pub fn database_url_from_env() -> String {
    std::env::var("DATABASE_URL").unwrap_or_else(|_| DEFAULT_DATABASE_URL.to_string())
}

#[must_use]
pub fn game_id_from_env() -> String {
    std::env::var("GAME_ID").unwrap_or_else(|_| "local-dev".to_string())
}

#[cfg(test)]
mod tests {
    use super::*;

    #[test]
    fn service_name_identifies_agent_service() {
        assert_eq!(SERVICE_NAME, "agent-service");
    }
}
