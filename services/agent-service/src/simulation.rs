pub const DEFAULT_SIMULATED_DATE: &str = "2026-10-01";
pub const DEFAULT_SPEED: u8 = 1;

#[derive(Debug, Clone, PartialEq, Eq)]
pub struct SimulationState {
    pub game_id: String,
    pub current_simulated_date: String,
    pub speed: u8,
    pub paused: bool,
    pub session_active: bool,
    pub last_tick_processed_at: Option<String>,
}

impl SimulationState {
    #[must_use]
    pub fn new(game_id: impl Into<String>) -> Self {
        Self {
            game_id: game_id.into(),
            current_simulated_date: DEFAULT_SIMULATED_DATE.to_string(),
            speed: DEFAULT_SPEED,
            paused: true,
            session_active: false,
            last_tick_processed_at: None,
        }
    }

    #[must_use]
    pub fn can_advance(&self) -> bool {
        self.session_active && !self.paused
    }

    pub fn set_session_active(&mut self, active: bool) {
        self.session_active = active;
    }

    pub fn set_paused(&mut self, paused: bool) {
        self.paused = paused;
    }

    pub fn set_speed(&mut self, speed: u8) -> Result<(), InvalidSpeed> {
        if matches!(speed, 1 | 5 | 20) {
            self.speed = speed;
            return Ok(());
        }

        Err(InvalidSpeed(speed))
    }
}

#[derive(Debug, Clone, Copy, PartialEq, Eq)]
pub struct InvalidSpeed(pub u8);

impl std::fmt::Display for InvalidSpeed {
    fn fmt(&self, formatter: &mut std::fmt::Formatter<'_>) -> std::fmt::Result {
        write!(formatter, "invalid simulation speed {}", self.0)
    }
}

impl std::error::Error for InvalidSpeed {}

#[cfg(test)]
mod tests {
    use super::*;

    #[test]
    fn new_state_starts_paused_without_active_session() {
        let state = SimulationState::new("game-1");

        assert_eq!(state.game_id, "game-1");
        assert_eq!(state.current_simulated_date, DEFAULT_SIMULATED_DATE);
        assert_eq!(state.speed, DEFAULT_SPEED);
        assert!(!state.session_active);
        assert!(state.paused);
        assert!(!state.can_advance());
    }

    #[test]
    fn state_can_advance_only_with_active_session_and_without_pause() {
        let mut state = SimulationState::new("game-1");

        state.set_session_active(true);
        assert!(!state.can_advance());

        state.set_paused(false);
        assert!(state.can_advance());

        state.set_session_active(false);
        assert!(!state.can_advance());
    }

    #[test]
    fn speed_accepts_only_m2_values() {
        let mut state = SimulationState::new("game-1");

        assert!(state.set_speed(5).is_ok());
        assert_eq!(state.speed, 5);

        assert_eq!(state.set_speed(2), Err(InvalidSpeed(2)));
        assert_eq!(state.speed, 5);
    }
}
