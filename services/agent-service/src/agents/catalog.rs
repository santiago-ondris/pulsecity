//! Canonical seed data for individual agents, relationships, and rival GMs.

use super::{
    AgentRelationship, AgentRelationshipSeed, IndividualAgentState, IndividualAgentTemplate,
    RivalGMProfile, clamp, clamp_unit, map_from_pairs, string_map_from_pairs,
};

#[must_use]
pub fn canonical_relationship_seeds() -> Vec<AgentRelationshipSeed> {
    vec![
        relationship_seed(
            "head_coach",
            "head_analytics",
            -0.18,
            "Guerra fria entre ojo y dato",
            "stable",
        ),
        relationship_seed(
            "head_coach",
            "team_doctor",
            -0.22,
            "Disponibilidad vs salud del jugador",
            "stable",
        ),
        relationship_seed(
            "gm",
            "team_doctor",
            0.02,
            "Confianza medica pendiente de decisiones reales",
            "stable",
        ),
        relationship_seed(
            "head_coach",
            "player_development_director",
            -0.12,
            "Desarrollo de largo plazo vs necesidades inmediatas",
            "stable",
        ),
        relationship_seed(
            "scouting_director",
            "head_analytics",
            -0.20,
            "Evaluacion tradicional vs modelos",
            "stable",
        ),
        relationship_seed(
            "marketing_director",
            "gm",
            -0.08,
            "Jugador marketeable vs jugador optimo",
            "stable",
        ),
        relationship_seed(
            "cfo",
            "gm",
            -0.10,
            "Presupuesto vs calidad del roster",
            "stable",
        ),
        relationship_seed(
            "mayor",
            "owner",
            -0.14,
            "Agenda politica vs intereses de la franquicia",
            "stable",
        ),
        relationship_seed(
            "chamber_commerce_president",
            "mayor",
            0.04,
            "Agenda economica y agenda politica se observan de cerca",
            "stable",
        ),
        relationship_seed(
            "sports_psychologist",
            "head_coach",
            -0.16,
            "Bienestar del jugador vs disponibilidad inmediata",
            "stable",
        ),
        relationship_seed(
            "pr_director",
            "gm",
            -0.06,
            "Control narrativo vs decisiones del GM",
            "stable",
        ),
        relationship_seed(
            "press",
            "roster_collective",
            -0.12,
            "Cobertura vs privacidad y estado emocional de jugadores",
            "stable",
        ),
    ]
}

#[must_use]
pub fn individual_agent_templates() -> Vec<IndividualAgentTemplate> {
    vec![
        template(
            "owner",
            "Owner",
            "basketball_ops",
            "Owner",
            "ownership, franchise vision, pressure, spending",
            "watchful",
            0.05,
            0.05,
            0.80,
            0.72,
            &[
                ("sporting_trust", 0.0),
                ("business_trust", 0.0),
                ("patience_remaining", 0.75),
                ("disposition_spending", 0.58),
            ],
            &[
                (
                    "primary_goal",
                    "make the expansion franchise credible without losing business discipline",
                ),
                ("risk_posture", "measured ambition"),
            ],
        ),
        template(
            "president_basketball_ops",
            "President of Basketball Operations",
            "basketball_ops",
            "President of Basketball Operations",
            "front office alignment, roster philosophy, owner buffer",
            "measured",
            0.02,
            0.04,
            0.72,
            0.70,
            &[
                ("owner_access", 0.70),
                ("gm_alignment", 0.05),
                ("leadership_control", 0.48),
                ("job_security", 0.68),
            ],
            &[
                ("primary_goal", "build a coherent basketball identity"),
                ("risk_posture", "protect organizational credibility"),
            ],
        ),
        template(
            "assistant_gm_cap",
            "Assistant GM, Cap Strategy",
            "basketball_ops",
            "Assistant General Manager",
            "contracts, cap sheets, trade mechanics",
            "focused",
            0.03,
            0.02,
            0.70,
            0.74,
            &[
                ("technical_level", 0.76),
                ("negotiation_skill", 0.68),
                ("gm_ambition", 0.55),
                ("cap_discipline", 0.82),
            ],
            &[
                ("specialization", "contracts and salary cap"),
                ("primary_goal", "keep optionality for future moves"),
            ],
        ),
        template(
            "assistant_gm_personnel",
            "Assistant GM, Player Personnel",
            "basketball_ops",
            "Assistant General Manager",
            "personnel evaluation, market reads, player agents",
            "curious",
            0.02,
            0.03,
            0.68,
            0.71,
            &[
                ("technical_level", 0.72),
                ("negotiation_skill", 0.70),
                ("gm_ambition", 0.62),
                ("contact_network", 0.66),
            ],
            &[
                ("specialization", "scouting and personnel"),
                ("primary_goal", "find undervalued players before the market"),
            ],
        ),
        template(
            "assistant_gm_operations",
            "Assistant GM, Operations",
            "basketball_ops",
            "Assistant General Manager",
            "front office operations, staff coordination, internal process",
            "organized",
            0.02,
            0.03,
            0.69,
            0.70,
            &[
                ("technical_level", 0.68),
                ("process_discipline", 0.74),
                ("gm_ambition", 0.48),
                ("internal_coordination", 0.72),
            ],
            &[
                ("specialization", "basketball operations"),
                (
                    "primary_goal",
                    "make the front office operate like one coherent room",
                ),
            ],
        ),
        template(
            "scouting_director",
            "Director de Scouting",
            "basketball_ops",
            "Director de Scouting",
            "draft, scouting reports, talent projection",
            "motivated",
            0.03,
            0.05,
            0.74,
            0.73,
            &[
                ("historical_precision", 0.50),
                ("coverage_capacity", 0.62),
                ("criteria_trust", 0.0),
                ("workload", 0.35),
            ],
            &[
                (
                    "evaluation_bias",
                    "values basketball IQ and positional size",
                ),
                (
                    "primary_goal",
                    "prove the expansion pipeline can find real talent",
                ),
            ],
        ),
        template(
            "player_personnel_director",
            "Director de Player Personnel",
            "basketball_ops",
            "Director de Player Personnel",
            "roster balance, player agents, trade opportunities",
            "alert",
            0.02,
            0.03,
            0.70,
            0.72,
            &[
                ("market_reading", 0.68),
                ("trade_management", 0.64),
                ("agent_relationships", 0.66),
                ("scouting_tension", 0.18),
            ],
            &[
                (
                    "primary_goal",
                    "turn roster gaps into practical opportunities",
                ),
                ("risk_posture", "opportunistic but not reckless"),
            ],
        ),
        template(
            "head_analytics",
            "Head of Analytics",
            "basketball_ops",
            "Head of Analytics",
            "models, lineup data, projections",
            "analytical",
            0.01,
            0.02,
            0.66,
            0.75,
            &[
                ("model_sophistication", 0.72),
                ("communication_capacity", 0.58),
                ("data_coverage", 0.64),
                ("traditional_scouting_tension", 0.28),
            ],
            &[
                ("primary_goal", "make decisions falsifiable and measurable"),
                ("risk_posture", "trust the signal over reputation"),
            ],
        ),
        template(
            "head_coach",
            "Head Coach",
            "basketball_ops",
            "Head Coach",
            "rotation, tactics, locker room leadership",
            "calm",
            0.0,
            0.0,
            0.70,
            0.70,
            &[
                ("gm_trust", 0.0),
                ("roster_satisfaction", 0.0),
                ("results_pressure", 0.25),
                ("locker_room_relationship", 0.0),
            ],
            &[
                ("system", "balanced pace and space"),
                (
                    "primary_goal",
                    "turn an expansion roster into a serious team identity",
                ),
            ],
        ),
        template(
            "assistant_coach_offense",
            "Assistant Coach, Offense",
            "basketball_ops",
            "Assistant Coach",
            "offensive sets, shot profile, player communication",
            "engaged",
            0.02,
            0.03,
            0.72,
            0.69,
            &[
                ("technical_level", 0.70),
                ("player_communication", 0.68),
                ("head_coach_loyalty", 0.78),
                ("head_coach_ambition", 0.42),
            ],
            &[
                ("specialization", "offense"),
                (
                    "primary_goal",
                    "create efficient shots with limited star power",
                ),
            ],
        ),
        template(
            "assistant_coach_defense",
            "Assistant Coach, Defense",
            "basketball_ops",
            "Assistant Coach",
            "defensive system, matchup prep, effort standards",
            "demanding",
            0.02,
            0.02,
            0.72,
            0.70,
            &[
                ("technical_level", 0.72),
                ("player_communication", 0.62),
                ("head_coach_loyalty", 0.74),
                ("defensive_standards", 0.80),
            ],
            &[
                ("specialization", "defense"),
                ("primary_goal", "make effort non-negotiable"),
            ],
        ),
        template(
            "player_development_director",
            "Director de Player Development",
            "basketball_ops",
            "Director de Player Development",
            "young players, skill plans, long-term growth",
            "patient",
            0.03,
            0.04,
            0.76,
            0.73,
            &[
                ("development_history", 0.58),
                ("young_player_trust", 0.65),
                ("coach_alignment", 0.52),
                ("patience", 0.82),
            ],
            &[
                (
                    "methodology",
                    "individual plans with slow compounding gains",
                ),
                (
                    "primary_goal",
                    "make the roster better by February than it is in October",
                ),
            ],
        ),
        template(
            "team_doctor",
            "Medico del Equipo",
            "basketball_ops",
            "Medico del Equipo",
            "injury diagnosis, return protocol, health risk",
            "careful",
            0.04,
            0.04,
            0.78,
            0.76,
            &[
                ("diagnostic_level", 0.78),
                ("return_protocol_conservatism", 0.72),
                ("coach_tension", 0.24),
                ("player_trust", 0.62),
            ],
            &[
                ("specialization", "sports medicine and orthopedics"),
                (
                    "primary_goal",
                    "avoid preventable injuries even under competitive pressure",
                ),
            ],
        ),
        template(
            "strength_conditioning_coach",
            "Fisioterapeuta / Strength & Conditioning Coach",
            "basketball_ops",
            "Strength & Conditioning Coach",
            "load management, prevention, conditioning",
            "steady",
            0.03,
            0.04,
            0.76,
            0.74,
            &[
                ("prevention_methodology", 0.72),
                ("personalization", 0.66),
                ("coach_tension", 0.22),
                ("player_trust", 0.64),
            ],
            &[
                (
                    "primary_goal",
                    "keep bodies available without flattening intensity",
                ),
                ("risk_posture", "prevention first"),
            ],
        ),
        template(
            "sports_psychologist",
            "Sports Psychologist",
            "basketball_ops",
            "Sports Psychologist",
            "emotional climate, burnout, player trust",
            "attentive",
            0.04,
            0.04,
            0.78,
            0.75,
            &[
                ("locker_room_climate", 0.0),
                ("emotional_alert", 0.2),
                ("player_trust", 0.0),
                ("confidentiality_tension", 0.35),
            ],
            &[
                ("methodology", "trust-first emotional diagnostics"),
                (
                    "primary_goal",
                    "catch pressure before it becomes performance collapse",
                ),
            ],
        ),
        template(
            "video_coordinator",
            "Video Coordinator",
            "basketball_ops",
            "Video Coordinator",
            "film, scouting clips, opponent tendencies",
            "busy",
            0.02,
            0.03,
            0.72,
            0.68,
            &[
                ("production_speed", 0.70),
                ("analysis_quality", 0.62),
                ("opponent_coverage", 0.58),
                ("analytics_relationship", 0.48),
            ],
            &[
                (
                    "primary_goal",
                    "turn raw film into useful decisions quickly",
                ),
                ("risk_posture", "detail-oriented"),
            ],
        ),
        template(
            "international_scout",
            "International Scout",
            "basketball_ops",
            "International Scout",
            "international scouting, cultural fit, overseas markets",
            "observant",
            0.01,
            0.03,
            0.68,
            0.69,
            &[
                ("geographic_coverage", 0.62),
                ("international_network", 0.66),
                ("evaluation_adaptability", 0.64),
                ("director_relationship", 0.52),
            ],
            &[
                ("coverage", "Europe, Latin America, emerging markets"),
                ("primary_goal", "make global talent feel less unknown"),
            ],
        ),
        template(
            "ceo_business_ops",
            "CEO / President of Business Operations",
            "business_ops",
            "CEO / President of Business Operations",
            "business strategy, fan experience, city relationships",
            "composed",
            0.03,
            0.04,
            0.74,
            0.73,
            &[
                ("business_vision", 0.70),
                ("growth_orientation", 0.58),
                ("basketball_ops_relationship", 0.48),
                ("city_relationship", 0.55),
            ],
            &[
                (
                    "primary_goal",
                    "make the franchise feel like a civic institution and a healthy business",
                ),
                ("risk_posture", "growth with operational discipline"),
            ],
        ),
        template(
            "cfo",
            "CFO",
            "business_ops",
            "CFO",
            "budget, salary cap, financial risk",
            "calm",
            0.02,
            0.03,
            0.76,
            0.74,
            &[
                ("financial_trust", 0.0),
                ("budget_alert", 0.15),
                ("financial_conservatism", 0.55),
                ("cap_sophistication", 0.72),
            ],
            &[
                (
                    "primary_goal",
                    "protect flexibility and avoid expensive traps",
                ),
                ("risk_posture", "skeptical until numbers justify risk"),
            ],
        ),
        template(
            "marketing_director",
            "Director de Marketing & Brand",
            "business_ops",
            "Director de Marketing & Brand",
            "brand, fanbase, campaigns, marketability",
            "energetic",
            0.02,
            0.04,
            0.70,
            0.72,
            &[
                ("campaign_creativity", 0.72),
                ("fanbase_reading", 0.64),
                ("digital_capacity", 0.70),
                ("media_vs_sport_orientation", 0.42),
            ],
            &[
                (
                    "primary_goal",
                    "give the expansion team a recognizable public identity",
                ),
                (
                    "risk_posture",
                    "bold campaigns if the team gives her a story",
                ),
            ],
        ),
        template(
            "ticket_sales_director",
            "Director de Ticket Sales",
            "business_ops",
            "Director de Ticket Sales",
            "attendance, pricing, season tickets",
            "practical",
            0.02,
            0.03,
            0.70,
            0.71,
            &[
                ("bad_season_sales_capacity", 0.64),
                ("pricing_aggression", 0.48),
                ("loyalty_programs", 0.58),
                ("result_sensitivity", 0.68),
            ],
            &[
                (
                    "primary_goal",
                    "fill the building before wins are guaranteed",
                ),
                ("risk_posture", "protect long-term fan habits"),
            ],
        ),
        template(
            "partnerships_director",
            "Director de Corporate Partnerships & Sponsors",
            "business_ops",
            "Director de Corporate Partnerships & Sponsors",
            "sponsors, corporate relationships, activations",
            "polished",
            0.02,
            0.04,
            0.72,
            0.72,
            &[
                ("corporate_network", 0.66),
                ("commercial_negotiation", 0.68),
                ("retention_rate", 0.55),
                ("image_sensitivity", 0.70),
            ],
            &[
                (
                    "primary_goal",
                    "turn local excitement into durable sponsor trust",
                ),
                ("risk_posture", "image-conscious"),
            ],
        ),
        template(
            "pr_director",
            "Director de PR & Communications",
            "business_ops",
            "Director de PR & Communications",
            "media strategy, crisis response, public narrative",
            "vigilant",
            0.02,
            0.03,
            0.70,
            0.73,
            &[
                ("crisis_management", 0.74),
                ("proactive_narrative", 0.58),
                ("press_network", 0.66),
                ("control_need", 0.62),
            ],
            &[
                (
                    "primary_goal",
                    "keep the story coherent before the press defines it",
                ),
                ("risk_posture", "transparent only when prepared"),
            ],
        ),
        template(
            "arena_operations_director",
            "Director de Arena Operations",
            "business_ops",
            "Director de Arena Operations",
            "arena logistics, events, maintenance, fan experience",
            "operational",
            0.02,
            0.03,
            0.72,
            0.74,
            &[
                ("operational_efficiency", 0.72),
                ("alternate_event_capacity", 0.58),
                ("maintenance_management", 0.64),
                ("fan_experience", 0.66),
            ],
            &[
                ("primary_goal", "make every event feel professionally run"),
                ("risk_posture", "maintenance before spectacle"),
            ],
        ),
        template(
            "legal_counsel",
            "Legal Counsel",
            "business_ops",
            "Legal Counsel",
            "contracts, legal risk, regulatory questions",
            "precise",
            0.03,
            0.03,
            0.78,
            0.76,
            &[
                ("response_speed", 0.66),
                ("legal_network", 0.62),
                ("contract_risk_detection", 0.76),
                ("cfo_relationship", 0.64),
            ],
            &[
                (
                    "specialization",
                    "sports contracts and cap-adjacent legal risk",
                ),
                ("primary_goal", "make fast deals defensible"),
            ],
        ),
        template(
            "mayor",
            "Alcalde",
            "city",
            "Alcalde",
            "politics, permits, city agenda",
            "calculating",
            0.0,
            0.03,
            0.62,
            0.70,
            &[
                ("current_popularity", 0.55),
                ("concession_tolerance", 0.42),
                ("electorate_pressure", 0.50),
                ("franchise_view", 0.48),
            ],
            &[
                (
                    "political_agenda",
                    "visible growth without looking captured by private money",
                ),
                (
                    "primary_goal",
                    "make the franchise help the city more than it costs politically",
                ),
            ],
        ),
        template(
            "police_chief",
            "Jefe de Policia",
            "city",
            "Jefe de Policia",
            "stadium security, logistics, public safety",
            "firm",
            0.01,
            0.03,
            0.66,
            0.71,
            &[
                ("resource_priority", 0.52),
                ("operational_capacity", 0.64),
                ("effectiveness", 0.68),
                ("political_sensitivity", 0.58),
            ],
            &[
                (
                    "primary_goal",
                    "keep game days orderly without draining the rest of the city",
                ),
                ("risk_posture", "resource-aware"),
            ],
        ),
        template(
            "chamber_commerce_president",
            "Presidente de la Camara de Comercio",
            "city",
            "Presidente de la Camara de Comercio",
            "local business, sponsors, district economy",
            "opportunistic",
            0.02,
            0.04,
            0.66,
            0.70,
            &[
                ("franchise_alignment", 0.56),
                ("business_network", 0.70),
                ("mayor_relationship", 0.50),
                ("owner_relationship", 0.48),
            ],
            &[
                (
                    "primary_goal",
                    "make arena traffic become local business growth",
                ),
                ("risk_posture", "pro-growth"),
            ],
        ),
        template(
            "urbanism_director",
            "Director de Urbanismo",
            "city",
            "Director de Urbanismo",
            "permits, planning, zoning process",
            "technical",
            0.01,
            0.03,
            0.68,
            0.72,
            &[
                ("process_efficiency", 0.62),
                ("technical_knowledge", 0.76),
                ("political_alignment", 0.54),
                ("regulation_tension", 0.46),
            ],
            &[
                (
                    "primary_goal",
                    "move projects without creating regulatory debt",
                ),
                ("risk_posture", "process-conscious"),
            ],
        ),
        template(
            "press",
            "La Prensa",
            "press",
            "Agente colectivo",
            "coverage, public sentiment, dominant narrative",
            "watching",
            0.0,
            0.0,
            0.50,
            0.68,
            &[
                ("general_sentiment", 0.0),
                ("coverage_intensity", 0.45),
                ("fanbase_impact", 0.62),
                ("sponsor_impact", 0.50),
            ],
            &[
                (
                    "dominant_narrative",
                    "the expansion experiment is interesting but unproven",
                ),
                (
                    "primary_goal",
                    "turn ambiguity into a readable public story",
                ),
            ],
        ),
    ]
}

impl IndividualAgentTemplate {
    pub(super) fn into_state(self, game_id: &str) -> IndividualAgentState {
        IndividualAgentState {
            game_id: game_id.to_string(),
            agent_id: self.agent_id.to_string(),
            display_name: self.display_name.to_string(),
            category: self.category.to_string(),
            role: self.role.to_string(),
            domain: self.domain.to_string(),
            emotional_state: self.emotional_state.to_string(),
            confidence: self.confidence,
            satisfaction: self.satisfaction,
            loyalty: self.loyalty,
            role_performance: self.role_performance,
            state: self.state,
            agenda: self.agenda,
        }
    }
}

impl AgentRelationshipSeed {
    pub(super) fn into_relationship(self, game_id: &str) -> AgentRelationship {
        AgentRelationship {
            game_id: game_id.to_string(),
            agent_a_id: self.agent_a_id.to_string(),
            agent_b_id: self.agent_b_id.to_string(),
            trust: clamp(self.trust),
            last_event: self.last_event.to_string(),
            trend: self.trend.to_string(),
            short_history: self.short_history.into_iter().map(str::to_string).collect(),
            last_source_event_id: None,
        }
    }
}

fn template(
    agent_id: &'static str,
    display_name: &'static str,
    category: &'static str,
    role: &'static str,
    domain: &'static str,
    emotional_state: &'static str,
    confidence: f64,
    satisfaction: f64,
    loyalty: f64,
    role_performance: f64,
    state: &[(&str, f64)],
    agenda: &[(&str, &str)],
) -> IndividualAgentTemplate {
    IndividualAgentTemplate {
        agent_id,
        display_name,
        category,
        role,
        domain,
        emotional_state,
        confidence,
        satisfaction,
        loyalty,
        role_performance,
        state: map_from_pairs(state),
        agenda: string_map_from_pairs(agenda),
    }
}

fn relationship_seed(
    agent_a_id: &'static str,
    agent_b_id: &'static str,
    trust: f64,
    last_event: &'static str,
    trend: &'static str,
) -> AgentRelationshipSeed {
    AgentRelationshipSeed {
        agent_a_id,
        agent_b_id,
        trust,
        last_event,
        trend,
        short_history: vec![last_event],
    }
}

#[derive(Debug, Clone, Copy, PartialEq, Eq)]
pub(super) struct RivalTeamTemplate {
    team_id: &'static str,
    team_name: &'static str,
    gm_name: &'static str,
}

pub(super) fn rival_team_templates() -> Vec<RivalTeamTemplate> {
    vec![
        rival_team("atl", "Atlanta Hawks", "Darius Bell"),
        rival_team("bos", "Boston Celtics", "Elliot Walsh"),
        rival_team("bkn", "Brooklyn Nets", "Miles Abram"),
        rival_team("cha", "Charlotte Hornets", "Noah Whitaker"),
        rival_team("chi", "Chicago Bulls", "Grant Mercer"),
        rival_team("cle", "Cleveland Cavaliers", "Isaiah Cole"),
        rival_team("dal", "Dallas Mavericks", "Owen Strickland"),
        rival_team("den", "Denver Nuggets", "Calvin Price"),
        rival_team("det", "Detroit Pistons", "Marcus Vale"),
        rival_team("gsw", "Golden State Warriors", "Adrian Cross"),
        rival_team("hou", "Houston Rockets", "Victor Hayes"),
        rival_team("ind", "Indiana Pacers", "Simon Reed"),
        rival_team("lac", "LA Clippers", "Julian Pierce"),
        rival_team("lal", "Los Angeles Lakers", "Wesley Kane"),
        rival_team("mem", "Memphis Grizzlies", "Malcolm Brooks"),
        rival_team("mia", "Miami Heat", "Rafael Soto"),
        rival_team("mil", "Milwaukee Bucks", "Leonard Frost"),
        rival_team("min", "Minnesota Timberwolves", "Theo Ramsey"),
        rival_team("nop", "New Orleans Pelicans", "Bennett Shaw"),
        rival_team("nyk", "New York Knicks", "Dominic Hale"),
        rival_team("okc", "Oklahoma City Thunder", "Preston Lane"),
        rival_team("orl", "Orlando Magic", "Felix Ward"),
        rival_team("phi", "Philadelphia 76ers", "Jonah Bishop"),
        rival_team("phx", "Phoenix Suns", "Gideon Marks"),
        rival_team("por", "Portland Trail Blazers", "Evan Holt"),
        rival_team("sac", "Sacramento Kings", "Trevor Nash"),
        rival_team("sas", "San Antonio Spurs", "Nolan Voss"),
        rival_team("tor", "Toronto Raptors", "Arthur Quinn"),
        rival_team("uta", "Utah Jazz", "Silas Finch"),
        rival_team("was", "Washington Wizards", "Graham Lowell"),
    ]
}

fn rival_team(
    team_id: &'static str,
    team_name: &'static str,
    gm_name: &'static str,
) -> RivalTeamTemplate {
    RivalTeamTemplate {
        team_id,
        team_name,
        gm_name,
    }
}

pub(super) fn rival_gm_profile(
    game_id: &str,
    index: usize,
    team: RivalTeamTemplate,
) -> RivalGMProfile {
    let styles = [
        "patient_value_hunter",
        "aggressive_star_chaser",
        "defensive_conservative",
        "asset_accumulator",
        "cap_flexible_operator",
        "win_now_pressure",
    ];
    let philosophies = [
        "draft_and_develop",
        "star_driven_contention",
        "defense_first_identity",
        "pace_and_spacing",
        "financial_optionalidad",
        "veteran_stability",
    ];
    let position_groups = [
        ["PG", "bench_shooting"],
        ["SG", "point_of_attack_defense"],
        ["SF", "wing_depth"],
        ["PF", "frontcourt_size"],
        ["C", "rim_protection"],
        ["secondary_creation", "salary_relief"],
    ];

    let style_index = deterministic_bucket(game_id, team.team_id, "style", styles.len());
    let philosophy_index =
        deterministic_bucket(game_id, team.team_id, "philosophy", philosophies.len());
    let needs_index = deterministic_bucket(game_id, team.team_id, "needs", position_groups.len());
    let urgency_bucket = deterministic_bucket(game_id, team.team_id, "urgency", 41);
    let trust_bucket = deterministic_bucket(game_id, team.team_id, "trust", 31);

    RivalGMProfile {
        game_id: game_id.to_string(),
        rival_team_id: team.team_id.to_string(),
        gm_agent_id: format!("rival_gm_{}", team.team_id),
        display_name: team.gm_name.to_string(),
        team_name: team.team_name.to_string(),
        negotiation_style: styles[style_index].to_string(),
        urgency_current: clamp_unit(0.25 + (urgency_bucket as f64 * 0.0125)),
        build_philosophy: philosophies[philosophy_index].to_string(),
        roster_needs: position_groups[needs_index]
            .iter()
            .map(|need| (*need).to_string())
            .collect(),
        relationship_trust: clamp(-0.15 + (trust_bucket as f64 * 0.01)),
        relationship_history: vec![
            "Sin historial directo con el GM de PulseCity".to_string(),
            format!("Perfil inicial sembrado como rival #{:02}", index + 1),
        ],
        last_interaction_event_id: None,
    }
}

fn deterministic_bucket(game_id: &str, team_id: &str, salt: &str, modulo: usize) -> usize {
    let mut hash = 14_695_981_039_346_656_037_u64;
    for byte in game_id.bytes().chain(team_id.bytes()).chain(salt.bytes()) {
        hash ^= u64::from(byte);
        hash = hash.wrapping_mul(1_099_511_628_211);
    }

    (hash as usize) % modulo
}
