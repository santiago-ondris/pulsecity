//! Core scaffolding for PulseCity's match-service.

pub mod events;

pub const SERVICE_NAME: &str = "match-service";
pub const DEFAULT_NATS_URL: &str = "nats://127.0.0.1:4222";

#[must_use]
pub fn nats_url_from_env() -> String {
    std::env::var("NATS_URL").unwrap_or_else(|_| DEFAULT_NATS_URL.to_string())
}

#[cfg(test)]
mod tests {
    use super::*;

    #[test]
    fn service_name_identifies_match_service() {
        assert_eq!(SERVICE_NAME, "match-service");
    }
}
