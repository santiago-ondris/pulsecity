use std::{error::Error, fmt};

use crate::events::{
    EventMeta, KeyMoment, MatchFinishedEvent, MatchPlayer, MatchScheduledEvent, MatchTeam,
    PlayerBoxScore,
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
    pub players: Vec<PlayerSimulationInput>,
    pub seed: u64,
}

#[derive(Debug, Clone, PartialEq)]
pub struct PlayerSimulationInput {
    pub player_id: String,
    pub team_id: String,
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
    let home_strength = team_strength(&input.home_team, &home_players, true);
    let away_strength = team_strength(&input.away_team, &away_players, false);
    let pace = (u16::from(input.home_team.pace) + u16::from(input.away_team.pace)) / 2;
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

    let home_box = build_box_score(&mut rng, &home_players, home_score, true);
    let away_box = build_box_score(&mut rng, &away_players, away_score, false);
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
            &home_players,
            &away_players,
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

fn team_strength(team: &MatchTeam, players: &[&PlayerSimulationInput], is_home: bool) -> i16 {
    let rotation_rating = players
        .iter()
        .take(8)
        .map(|player| {
            i16::from(player.rating) - i16::from(player.fatigue / 4)
                + i16::from(player.emotional_state / 4)
        })
        .sum::<i16>()
        / 8;
    let home_bonus = if is_home {
        i16::from(team.home_court_advantage)
    } else {
        0
    };

    (i16::from(team.rating)
        + i16::from(team.offense_rating)
        + i16::from(team.defense_rating)
        + rotation_rating)
        / 4
        + home_bonus
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
    players: &[&PlayerSimulationInput],
    team_score: u16,
    starters_first: bool,
) -> Vec<PlayerBoxScore> {
    let rotation_size = players.len().min(10);
    let rotation = &players[..rotation_size];
    let weights: Vec<u16> = rotation
        .iter()
        .enumerate()
        .map(|(index, player)| {
            let role_bonus = if index < 5 { 26 } else { 12 };
            u16::from(player.scoring) + role_bonus + u16::from(rng.range_u8(0, 9))
        })
        .collect();
    let total_weight = weights.iter().copied().sum::<u16>().max(1);
    let mut remaining_points = team_score;
    let mut box_score = Vec::with_capacity(rotation_size);

    for (index, player) in rotation.iter().enumerate() {
        let is_last = index == rotation_size - 1;
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
        let minutes_base = if index < 5 { 29 } else { 15 };
        let minutes = minutes_base + rng.range_u8(0, 6);
        let rebound_factor = u16::from(player.rebounding) / 12;
        let assist_factor = u16::from(player.playmaking) / 16;
        let defense_factor = u16::from(player.defense) / 28;

        box_score.push(PlayerBoxScore {
            player_id: player.player_id.clone(),
            team_id: player.team_id.clone(),
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
    home_players: &[&PlayerSimulationInput],
    away_players: &[&PlayerSimulationInput],
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
    let winner_lead = &winner_players[0];
    let loser_lead = &loser_players[0];

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
            players,
            seed: 42,
        }
    }

    fn sample_player(team_id: &str, index: u8) -> PlayerSimulationInput {
        PlayerSimulationInput {
            player_id: format!("{team_id}-player-{index}"),
            team_id: team_id.to_string(),
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
}
