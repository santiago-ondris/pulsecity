pub const DEFAULT_SIMULATED_DATE: &str = "2026-10-01";
pub const DEFAULT_SPEED: u8 = 1;
pub const REAL_MILLIS_PER_SIMULATED_DAY_X1: u128 = 1_600;

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
pub struct TickAdvance {
    pub days_processed: u16,
}

#[derive(Debug, Clone, Copy, PartialEq, Eq)]
pub struct SimulationAccumulator {
    elapsed_millis: u128,
}

impl SimulationAccumulator {
    #[must_use]
    pub const fn new() -> Self {
        Self { elapsed_millis: 0 }
    }

    pub fn reset(&mut self) {
        self.elapsed_millis = 0;
    }

    #[must_use]
    pub const fn elapsed_millis(&self) -> u128 {
        self.elapsed_millis
    }

    pub fn tick(&mut self, state: &SimulationState, elapsed_millis: u128) -> TickAdvance {
        if !state.can_advance() {
            self.reset();
            return TickAdvance { days_processed: 0 };
        }

        self.elapsed_millis += elapsed_millis;
        let millis_per_day = REAL_MILLIS_PER_SIMULATED_DAY_X1 / u128::from(state.speed);
        let days_processed = self.elapsed_millis / millis_per_day;
        self.elapsed_millis %= millis_per_day;

        TickAdvance {
            days_processed: days_processed.min(u128::from(u16::MAX)) as u16,
        }
    }
}

impl Default for SimulationAccumulator {
    fn default() -> Self {
        Self::new()
    }
}

pub fn advance_simulated_date(current: &str, days: u16) -> Result<String, InvalidSimulatedDate> {
    let (year, month, day) = parse_date(current)?;
    let day_number = days_from_civil(year, month, day);
    let (next_year, next_month, next_day) = civil_from_days(day_number + i64::from(days));

    Ok(format!("{next_year:04}-{next_month:02}-{next_day:02}"))
}

#[derive(Debug, Clone, Copy, PartialEq, Eq)]
pub struct InvalidSpeed(pub u8);

impl std::fmt::Display for InvalidSpeed {
    fn fmt(&self, formatter: &mut std::fmt::Formatter<'_>) -> std::fmt::Result {
        write!(formatter, "invalid simulation speed {}", self.0)
    }
}

impl std::error::Error for InvalidSpeed {}

#[derive(Debug, Clone, PartialEq, Eq)]
pub struct InvalidSimulatedDate(String);

impl std::fmt::Display for InvalidSimulatedDate {
    fn fmt(&self, formatter: &mut std::fmt::Formatter<'_>) -> std::fmt::Result {
        write!(formatter, "invalid simulated date {}", self.0)
    }
}

impl std::error::Error for InvalidSimulatedDate {}

fn parse_date(input: &str) -> Result<(i32, u32, u32), InvalidSimulatedDate> {
    let mut parts = input.split('-');
    let year = parts
        .next()
        .and_then(|part| part.parse::<i32>().ok())
        .ok_or_else(|| InvalidSimulatedDate(input.to_string()))?;
    let month = parts
        .next()
        .and_then(|part| part.parse::<u32>().ok())
        .ok_or_else(|| InvalidSimulatedDate(input.to_string()))?;
    let day = parts
        .next()
        .and_then(|part| part.parse::<u32>().ok())
        .ok_or_else(|| InvalidSimulatedDate(input.to_string()))?;

    if parts.next().is_some() || !(1..=12).contains(&month) || !(1..=31).contains(&day) {
        return Err(InvalidSimulatedDate(input.to_string()));
    }

    Ok((year, month, day))
}

fn days_from_civil(year: i32, month: u32, day: u32) -> i64 {
    let year = year - i32::from(month <= 2);
    let era = if year >= 0 { year } else { year - 399 } / 400;
    let year_of_era = year - era * 400;
    let month = month as i32;
    let day = day as i32;
    let day_of_year = (153 * (month + if month > 2 { -3 } else { 9 }) + 2) / 5 + day - 1;
    let day_of_era = year_of_era * 365 + year_of_era / 4 - year_of_era / 100 + day_of_year;

    i64::from(era * 146_097 + day_of_era - 719_468)
}

fn civil_from_days(days: i64) -> (i32, u32, u32) {
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

    #[test]
    fn accumulator_processes_days_by_speed() {
        let mut state = SimulationState::new("game-1");
        state.set_session_active(true);
        state.set_paused(false);

        let mut accumulator = SimulationAccumulator::new();
        assert_eq!(accumulator.tick(&state, 1_500).days_processed, 0);
        assert_eq!(accumulator.elapsed_millis(), 1_500);
        assert_eq!(accumulator.tick(&state, 100).days_processed, 1);
        assert_eq!(accumulator.elapsed_millis(), 0);

        state.set_speed(5).expect("x5 is valid");
        assert_eq!(accumulator.tick(&state, 320).days_processed, 1);

        state.set_speed(20).expect("x20 is valid");
        assert_eq!(accumulator.tick(&state, 240).days_processed, 3);
    }

    #[test]
    fn accumulator_resets_when_time_cannot_advance() {
        let mut state = SimulationState::new("game-1");
        state.set_session_active(true);
        state.set_paused(false);

        let mut accumulator = SimulationAccumulator::new();
        assert_eq!(accumulator.tick(&state, 800).days_processed, 0);

        state.set_paused(true);
        assert_eq!(accumulator.tick(&state, 100).days_processed, 0);
        assert_eq!(accumulator.elapsed_millis(), 0);
    }

    #[test]
    fn simulated_date_advances_across_months_and_leap_years() {
        assert_eq!(
            advance_simulated_date("2026-10-01", 1).expect("valid date"),
            "2026-10-02"
        );
        assert_eq!(
            advance_simulated_date("2026-10-31", 1).expect("valid date"),
            "2026-11-01"
        );
        assert_eq!(
            advance_simulated_date("2028-02-28", 1).expect("valid date"),
            "2028-02-29"
        );
    }
}
