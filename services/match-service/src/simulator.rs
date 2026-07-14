use std::{error::Error, fmt};

use crate::events::{
    EventMeta, KeyMoment, MatchFinishedEvent, MatchPlayer, MatchScheduledEvent,
    MatchTacticalContext, MatchTeam, PlayerBoxScore,
};

const SCHEMA_VERSION: u16 = 1;
const MIN_PLAYERS_PER_TEAM: usize = 5;

#[derive(Debug, Clone, PartialEq)]
pub struct MatchSimulationInput {
    pub game_id: String,
    pub match_id: String,
    pub simulated_date: String,
    pub home_team: MatchTeam,
    pub away_team: MatchTeam,
    pub home_tactics: Option<MatchTacticalContext>,
    pub away_tactics: Option<MatchTacticalContext>,
    pub players: Vec<PlayerSimulationInput>,
    pub seed: u64,
}

#[derive(Debug, Clone, PartialEq)]
pub struct PlayerSimulationInput {
    pub player_id: String,
    pub team_id: String,
    pub expected_minutes: Option<u8>,
    pub rating: u8,
    pub scoring: u8,
    pub rebounding: u8,
    pub playmaking: u8,
    pub defense: u8,
    pub stamina: u8,
    pub fatigue: u8,
    pub emotional_state: i8,
}

#[derive(Debug, Clone, PartialEq, Eq)]
pub enum SimulationError {
    MissingGameID,
    MissingMatchID,
    MissingSimulatedDate,
    NotEnoughPlayers {
        team_id: String,
        found: usize,
        minimum: usize,
    },
}

impl fmt::Display for SimulationError {
    fn fmt(&self, formatter: &mut fmt::Formatter<'_>) -> fmt::Result {
        match self {
            Self::MissingGameID => write!(formatter, "missing game id"),
            Self::MissingMatchID => write!(formatter, "missing match id"),
            Self::MissingSimulatedDate => write!(formatter, "missing simulated date"),
            Self::NotEnoughPlayers {
                team_id,
                found,
                minimum,
            } => write!(
                formatter,
                "team {team_id} has {found} players, minimum required is {minimum}"
            ),
        }
    }
}

impl Error for SimulationError {}

pub fn simulate_match(input: &MatchSimulationInput) -> Result<MatchFinishedEvent, SimulationError> {
    validate_input(input)?;

    let mut rng = DeterministicRng::new(input.seed);
    let home_players = team_players(&input.players, &input.home_team.team_id);
    let away_players = team_players(&input.players, &input.away_team.team_id);
    let home_rotation = build_rotation(&home_players, input.home_tactics.as_ref());
    let away_rotation = build_rotation(&away_players, input.away_tactics.as_ref());
    let home_strength = team_strength(
        &input.home_team,
        input.home_tactics.as_ref(),
        &home_rotation,
        true,
    );
    let away_strength = team_strength(
        &input.away_team,
        input.away_tactics.as_ref(),
        &away_rotation,
        false,
    );
    let pace = tactical_pace(&input.home_team, input.home_tactics.as_ref())
        + tactical_pace(&input.away_team, input.away_tactics.as_ref());
    let pace = pace / 2;
    let base_score = 82 + pace / 4;
    let home_noise = i16::from(rng.range_u8(0, 17)) - 8;
    let away_noise = i16::from(rng.range_u8(0, 17)) - 8;
    let mut home_score = score_from_strength(base_score, home_strength, away_strength, home_noise);
    let mut away_score = score_from_strength(base_score, away_strength, home_strength, away_noise);

    if home_score == away_score {
        if rng.bool() {
            home_score += 1;
        } else {
            away_score += 1;
        }
    }

    let home_box = build_box_score(&mut rng, &home_rotation, home_score, true);
    let away_box = build_box_score(&mut rng, &away_rotation, away_score, false);
    let mut box_score = Vec::with_capacity(home_box.len() + away_box.len());
    box_score.extend(home_box);
    box_score.extend(away_box);

    let winner_team_id = if home_score > away_score {
        input.home_team.team_id.clone()
    } else {
        input.away_team.team_id.clone()
    };

    Ok(MatchFinishedEvent {
        meta: EventMeta {
            event_id: format!("match-finished-{}", input.match_id),
            game_id: input.game_id.clone(),
            occurred_at: format!("{}T00:00:00Z", input.simulated_date),
            schema_version: SCHEMA_VERSION,
        },
        match_id: input.match_id.clone(),
        simulated_date: input.simulated_date.clone(),
        home_team: input.home_team.clone(),
        away_team: input.away_team.clone(),
        home_score,
        away_score,
        winner_team_id,
        seed: input.seed,
        box_score,
        key_moments: build_key_moments(
            &mut rng,
            input,
            home_score,
            away_score,
            &home_rotation,
            &away_rotation,
        ),
    })
}

impl From<MatchScheduledEvent> for MatchSimulationInput {
    fn from(event: MatchScheduledEvent) -> Self {
        Self {
            game_id: event.meta.game_id,
            match_id: event.match_id,
            simulated_date: event.simulated_date,
            home_team: event.home_team,
            away_team: event.away_team,
            home_tactics: event.home_tactics,
            away_tactics: event.away_tactics,
            players: event
                .players
                .into_iter()
                .map(PlayerSimulationInput::from)
                .collect(),
            seed: event.seed,
        }
    }
}

impl From<MatchPlayer> for PlayerSimulationInput {
    fn from(player: MatchPlayer) -> Self {
        Self {
            player_id: player.player_id,
            team_id: player.team_id,
            expected_minutes: player.expected_minutes,
            rating: player.rating,
            scoring: player.scoring,
            rebounding: player.rebounding,
            playmaking: player.playmaking,
            defense: player.defense,
            stamina: player.stamina,
            fatigue: player.fatigue,
            emotional_state: player.emotional_state,
        }
    }
}

fn validate_input(input: &MatchSimulationInput) -> Result<(), SimulationError> {
    if input.game_id.trim().is_empty() {
        return Err(SimulationError::MissingGameID);
    }
    if input.match_id.trim().is_empty() {
        return Err(SimulationError::MissingMatchID);
    }
    if input.simulated_date.trim().is_empty() {
        return Err(SimulationError::MissingSimulatedDate);
    }

    for team_id in [&input.home_team.team_id, &input.away_team.team_id] {
        let found = input
            .players
            .iter()
            .filter(|player| player.team_id == *team_id)
            .count();
        if found < MIN_PLAYERS_PER_TEAM {
            return Err(SimulationError::NotEnoughPlayers {
                team_id: team_id.clone(),
                found,
                minimum: MIN_PLAYERS_PER_TEAM,
            });
        }
    }

    Ok(())
}

fn team_players<'a>(
    players: &'a [PlayerSimulationInput],
    team_id: &str,
) -> Vec<&'a PlayerSimulationInput> {
    let mut selected: Vec<_> = players
        .iter()
        .filter(|player| player.team_id == team_id)
        .collect();
    selected.sort_by_key(|player| std::cmp::Reverse(player.rating));
    selected
}

#[derive(Debug, Clone, Copy)]
struct RotationPlayer<'a> {
    player: &'a PlayerSimulationInput,
    minutes: u8,
}

fn build_rotation<'a>(
    players: &[&'a PlayerSimulationInput],
    tactics: Option<&MatchTacticalContext>,
) -> Vec<RotationPlayer<'a>> {
    let profile = rotation_profile(tactics);
    let mut rotation = Vec::with_capacity(players.len().min(profile.len()));

    for (index, player) in players.iter().enumerate() {
        let minutes = player
            .expected_minutes
            .unwrap_or_else(|| profile.get(index).copied().unwrap_or(0))
            .min(42);
        if minutes > 0 {
            rotation.push(RotationPlayer {
                player: *player,
                minutes,
            });
        }
    }

    if rotation.is_empty() {
        for (index, player) in players.iter().take(MIN_PLAYERS_PER_TEAM).enumerate() {
            rotation.push(RotationPlayer {
                player: *player,
                minutes: profile.get(index).copied().unwrap_or(12).max(1),
            });
        }
    }

    rotation.sort_by_key(|slot| std::cmp::Reverse(slot.minutes));
    rotation
}

fn rotation_profile(tactics: Option<&MatchTacticalContext>) -> &'static [u8] {
    match tactics.map(|context| context.rotation_preference.as_str()) {
        Some("top_heavy") => &[38, 36, 34, 32, 30, 24, 18, 14, 8, 6],
        Some("deep") => &[29, 28, 27, 26, 25, 23, 22, 21, 20, 19],
        _ => &[34, 32, 30, 29, 28, 24, 22, 18, 13, 10],
    }
}

fn team_strength(
    team: &MatchTeam,
    tactics: Option<&MatchTacticalContext>,
    rotation: &[RotationPlayer<'_>],
    is_home: bool,
) -> i16 {
    let total_minutes = rotation
        .iter()
        .map(|slot| u16::from(slot.minutes))
        .sum::<u16>()
        .max(1);
    let rotation_rating = rotation
        .iter()
        .map(|slot| {
            let player = slot.player;
            let player_rating = i32::from(player.rating) - i32::from(player.fatigue / 4)
                + i32::from(player.emotional_state / 4)
                + stamina_modifier(player.stamina);
            player_rating * i32::from(slot.minutes)
        })
        .sum::<i32>()
        / i32::from(total_minutes);
    let home_bonus = if is_home {
        i16::from(team.home_court_advantage)
    } else {
        0
    };
    let (offense_bonus, defense_bonus) = tactical_rating_bonus(tactics);

    (i16::from(team.rating)
        + i16::from(team.offense_rating)
        + offense_bonus
        + i16::from(team.defense_rating)
        + defense_bonus
        + rotation_rating as i16)
        / 4
        + home_bonus
}

fn stamina_modifier(stamina: u8) -> i32 {
    (i32::from(stamina) - 80) / 10
}

fn tactical_rating_bonus(tactics: Option<&MatchTacticalContext>) -> (i16, i16) {
    let flexibility_bonus = tactics
        .map(|context| i16::from(context.flexibility.saturating_sub(50)) / 25)
        .unwrap_or(0);

    match tactics.map(|context| context.system.as_str()) {
        Some("pace_and_space") => (2 + flexibility_bonus, -1),
        Some("defensive_grind") => (-1, 2 + flexibility_bonus),
        Some("balanced") | None => (flexibility_bonus, flexibility_bonus),
        _ => (0, 0),
    }
}

fn tactical_pace(team: &MatchTeam, tactics: Option<&MatchTacticalContext>) -> u16 {
    let adjustment = match tactics.map(|context| context.system.as_str()) {
        Some("pace_and_space") => 4,
        Some("defensive_grind") => -4,
        _ => 0,
    };

    (i16::from(team.pace) + adjustment).clamp(80, 110) as u16
}

fn score_from_strength(
    base_score: u16,
    own_strength: i16,
    opponent_strength: i16,
    noise: i16,
) -> u16 {
    let score =
        i16::try_from(base_score).unwrap_or(100) + (own_strength - opponent_strength) / 2 + noise;
    score.clamp(82, 138) as u16
}

fn build_box_score(
    rng: &mut DeterministicRng,
    rotation: &[RotationPlayer<'_>],
    team_score: u16,
    starters_first: bool,
) -> Vec<PlayerBoxScore> {
    let weights: Vec<u16> = rotation
        .iter()
        .map(|slot| {
            let minutes_weight = u16::from(slot.minutes).max(1);
            (u16::from(slot.player.scoring) + u16::from(rng.range_u8(0, 9))) * minutes_weight
        })
        .collect();
    let total_weight = weights.iter().copied().sum::<u16>().max(1);
    let mut remaining_points = team_score;
    let mut box_score = Vec::with_capacity(rotation.len());

    for (index, player) in rotation.iter().enumerate() {
        let is_last = index == rotation.len() - 1;
        let points = if is_last {
            remaining_points
        } else {
            let expected = (u32::from(team_score) * u32::from(weights[index])
                / u32::from(total_weight)) as u16;
            let variation = u16::from(rng.range_u8(0, 5));
            let points = expected.saturating_add(variation).min(remaining_points);
            remaining_points -= points;
            points
        };
        let minute_variation = i16::from(rng.range_u8(0, 5)) - 2;
        let minutes = (i16::from(player.minutes) + minute_variation).clamp(1, 42) as u8;
        let rebound_factor = u16::from(player.player.rebounding) / 12;
        let assist_factor = u16::from(player.player.playmaking) / 16;
        let defense_factor = u16::from(player.player.defense) / 28;

        box_score.push(PlayerBoxScore {
            player_id: player.player.player_id.clone(),
            team_id: player.player.team_id.clone(),
            minutes: if starters_first {
                minutes
            } else {
                minutes.saturating_sub(1)
            },
            points,
            rebounds: rebound_factor + u16::from(rng.range_u8(0, 5)),
            assists: assist_factor + u16::from(rng.range_u8(0, 4)),
            steals: defense_factor + u16::from(rng.range_u8(0, 2)),
            blocks: defense_factor / 2 + u16::from(rng.range_u8(0, 2)),
            turnovers: u16::from(rng.range_u8(0, 4)),
        });
    }

    box_score
}

fn build_key_moments(
    rng: &mut DeterministicRng,
    input: &MatchSimulationInput,
    home_score: u16,
    away_score: u16,
    home_players: &[RotationPlayer<'_>],
    away_players: &[RotationPlayer<'_>],
) -> Vec<KeyMoment> {
    let margin = home_score.abs_diff(away_score);
    let home_won = home_score > away_score;
    let winner = if home_won {
        &input.home_team
    } else {
        &input.away_team
    };
    let loser = if home_won {
        &input.away_team
    } else {
        &input.home_team
    };
    let winner_players = if home_won { home_players } else { away_players };
    let loser_players = if home_won { away_players } else { home_players };
    let winner_lead = winner_players[0].player;
    let loser_lead = loser_players[0].player;

    vec![
        KeyMoment {
            quarter: 1,
            clock: "06:42".to_string(),
            kind: "tempo_established".to_string(),
            description: format!(
                "{} marca el ritmo temprano y obliga a {} a ajustar la rotacion.",
                winner.name, loser.name
            ),
            team_id: winner.team_id.clone(),
            player_id: Some(winner_lead.player_id.clone()),
        },
        KeyMoment {
            quarter: 2,
            clock: "03:18".to_string(),
            kind: "bench_run".to_string(),
            description: format!(
                "La segunda unidad de {} sostiene el partido antes del descanso.",
                if rng.bool() {
                    &input.home_team.name
                } else {
                    &input.away_team.name
                }
            ),
            team_id: if rng.bool() {
                input.home_team.team_id.clone()
            } else {
                input.away_team.team_id.clone()
            },
            player_id: None,
        },
        KeyMoment {
            quarter: 4,
            clock: if margin <= 5 { "01:12" } else { "04:36" }.to_string(),
            kind: if margin <= 5 {
                "clutch_sequence"
            } else {
                "decisive_run"
            }
            .to_string(),
            description: if margin <= 5 {
                format!(
                    "{} cierra una posesion pesada para proteger una ventaja minima.",
                    winner_lead.player_id
                )
            } else {
                format!(
                    "{} rompe el partido con una racha que deja sin respuesta a {}.",
                    winner.name, loser.name
                )
            },
            team_id: winner.team_id.clone(),
            player_id: Some(winner_lead.player_id.clone()),
        },
        KeyMoment {
            quarter: 4,
            clock: "00:18".to_string(),
            kind: "final_response".to_string(),
            description: format!(
                "{} intenta una ultima reaccion, pero el margen queda en {} puntos.",
                loser_lead.player_id, margin
            ),
            team_id: loser.team_id.clone(),
            player_id: Some(loser_lead.player_id.clone()),
        },
    ]
}

#[derive(Debug, Clone, Copy)]
struct DeterministicRng {
    state: u64,
}

impl DeterministicRng {
    const fn new(seed: u64) -> Self {
        Self {
            state: seed ^ 0x9E37_79B9_7F4A_7C15,
        }
    }

    fn next_u64(&mut self) -> u64 {
        self.state = self
            .state
            .wrapping_mul(6_364_136_223_846_793_005)
            .wrapping_add(1_442_695_040_888_963_407);
        self.state
    }

    fn range_u8(&mut self, start: u8, end_exclusive: u8) -> u8 {
        let width = end_exclusive.saturating_sub(start).max(1);
        start + (self.next_u64() % u64::from(width)) as u8
    }

    fn bool(&mut self) -> bool {
        self.next_u64() & 1 == 1
    }
}

#[cfg(test)]
mod tests {
    use super::*;

    #[test]
    fn simulate_match_is_deterministic_for_same_input() {
        let input = sample_input();

        let first = simulate_match(&input).expect("valid input");
        let second = simulate_match(&input).expect("valid input");

        assert_eq!(first, second);
        assert_ne!(first.home_score, first.away_score);
        assert_eq!(first.box_score.len(), 20);
        assert_eq!(first.key_moments.len(), 4);
    }

    #[test]
    fn simulate_match_changes_when_seed_changes() {
        let mut first_input = sample_input();
        let mut second_input = sample_input();
        first_input.seed = 11;
        second_input.seed = 12;

        let first = simulate_match(&first_input).expect("valid input");
        let second = simulate_match(&second_input).expect("valid input");

        assert_ne!(
            (first.home_score, first.away_score),
            (second.home_score, second.away_score)
        );
    }

    #[test]
    fn simulate_match_requires_full_rotation() {
        let mut input = sample_input();
        input.players.retain(|player| player.team_id != "away");

        let err = simulate_match(&input).expect_err("away team should be invalid");

        assert_eq!(
            err,
            SimulationError::NotEnoughPlayers {
                team_id: "away".to_string(),
                found: 0,
                minimum: MIN_PLAYERS_PER_TEAM
            }
        );
    }

    #[test]
    fn box_score_points_match_team_scores() {
        let result = simulate_match(&sample_input()).expect("valid input");
        let home_points = result
            .box_score
            .iter()
            .filter(|line| line.team_id == "home")
            .map(|line| line.points)
            .sum::<u16>();
        let away_points = result
            .box_score
            .iter()
            .filter(|line| line.team_id == "away")
            .map(|line| line.points)
            .sum::<u16>();

        assert_eq!(home_points, result.home_score);
        assert_eq!(away_points, result.away_score);
    }

    #[test]
    fn explicit_minutes_drive_box_score_role() {
        let mut input = sample_input();
        set_expected_minutes(&mut input, "home-player-0", 8);
        set_expected_minutes(&mut input, "home-player-9", 38);

        let result = simulate_match(&input).expect("valid input");
        let former_lead = box_line(&result, "home-player-0");
        let promoted_guard = box_line(&result, "home-player-9");

        assert!(promoted_guard.minutes > former_lead.minutes);
        assert!(promoted_guard.points > former_lead.points);
    }

    #[test]
    fn rotation_changes_match_output() {
        let standard = simulate_match(&sample_input()).expect("valid input");
        let mut changed_input = sample_input();
        set_expected_minutes(&mut changed_input, "home-player-0", 4);
        set_expected_minutes(&mut changed_input, "home-player-9", 40);
        downgrade_player(&mut changed_input, "home-player-9", 55);

        let changed = simulate_match(&changed_input).expect("valid input");

        assert_ne!(
            (standard.home_score, standard.away_score),
            (changed.home_score, changed.away_score)
        );
    }

    fn sample_input() -> MatchSimulationInput {
        let home_team = MatchTeam {
            team_id: "home".to_string(),
            name: "PulseCity Lighthouses".to_string(),
            abbreviation: "LHT".to_string(),
            rating: 78,
            offense_rating: 79,
            defense_rating: 76,
            pace: 99,
            home_court_advantage: 3,
        };
        let away_team = MatchTeam {
            team_id: "away".to_string(),
            name: "Seattle Rainmakers".to_string(),
            abbreviation: "SEA".to_string(),
            rating: 77,
            offense_rating: 76,
            defense_rating: 78,
            pace: 97,
            home_court_advantage: 2,
        };
        let mut players = Vec::with_capacity(20);
        for index in 0..10 {
            players.push(sample_player("home", index));
            players.push(sample_player("away", index));
        }

        MatchSimulationInput {
            game_id: "game-1".to_string(),
            match_id: "match-1".to_string(),
            simulated_date: "2026-10-22".to_string(),
            home_team,
            away_team,
            home_tactics: Some(MatchTacticalContext {
                system: "balanced".to_string(),
                rotation_preference: "standard".to_string(),
                flexibility: 55,
            }),
            away_tactics: Some(MatchTacticalContext {
                system: "balanced".to_string(),
                rotation_preference: "standard".to_string(),
                flexibility: 55,
            }),
            players,
            seed: 42,
        }
    }

    fn sample_player(team_id: &str, index: u8) -> PlayerSimulationInput {
        PlayerSimulationInput {
            player_id: format!("{team_id}-player-{index}"),
            team_id: team_id.to_string(),
            expected_minutes: None,
            rating: 78u8.saturating_sub(index / 2),
            scoring: 76u8.saturating_sub(index / 3),
            rebounding: 70 + index % 8,
            playmaking: 72 + index % 7,
            defense: 74 + index % 6,
            stamina: 82,
            fatigue: index,
            emotional_state: 2,
        }
    }

    fn set_expected_minutes(input: &mut MatchSimulationInput, player_id: &str, minutes: u8) {
        let player = input
            .players
            .iter_mut()
            .find(|player| player.player_id == player_id)
            .expect("sample player exists");
        player.expected_minutes = Some(minutes);
    }

    fn downgrade_player(input: &mut MatchSimulationInput, player_id: &str, rating: u8) {
        let player = input
            .players
            .iter_mut()
            .find(|player| player.player_id == player_id)
            .expect("sample player exists");
        player.rating = rating;
        player.scoring = rating;
        player.rebounding = rating;
        player.playmaking = rating;
        player.defense = rating;
    }

    fn box_line<'a>(result: &'a MatchFinishedEvent, player_id: &str) -> &'a PlayerBoxScore {
        result
            .box_score
            .iter()
            .find(|line| line.player_id == player_id)
            .expect("box score line exists")
    }
}
