use std::collections::BTreeMap;

use crate::events::{
    AgentRelationshipChangedEvent, AgentStateChangedEvent, EventMeta, GMDecisionRegisteredEvent,
    MatchFinishedEvent, PlayerBoxScore, PlayerEmotionalPatch, RosterPatchEnvelope,
    RosterStatePatch, SUBJECT_GM_DECISION_REGISTERED, SUBJECT_MATCH_FINISHED, SUBJECT_ROSTER_PATCH,
    SalaryCapCalculatedEvent,
};

pub const OWN_TEAM_ID: &str = "pulsecity";
pub const CORE_AGENT_IDS: [&str; 5] = [
    "owner",
    "head_coach",
    "cfo",
    "scouting_director",
    "sports_psychologist",
];

pub const INDIVIDUAL_AGENT_COUNT: usize = 30;

const MIN_STATE_VALUE: f64 = -1.0;
const MAX_STATE_VALUE: f64 = 1.0;
const SCHEMA_VERSION: u16 = 1;

#[derive(Debug, Clone, PartialEq)]
pub struct CoreAgentState {
    pub game_id: String,
    pub agent_id: String,
    pub mood: String,
    pub state: BTreeMap<String, f64>,
    pub last_match_id: Option<String>,
}

#[derive(Debug, Clone, PartialEq)]
pub struct AgentStateChange {
    pub state: CoreAgentState,
    pub event: AgentStateChangedEvent,
}

#[derive(Debug, Clone, PartialEq)]
pub struct TeamRosterPlayer {
    pub player_id: String,
    pub game_id: String,
    pub full_name: String,
    pub position: String,
    pub overall_rating: u8,
    pub roster_status: String,
}

#[derive(Debug, Clone, PartialEq)]
pub struct PlayerAgentState {
    pub game_id: String,
    pub player_id: String,
    pub full_name: String,
    pub position: String,
    pub emotional_state: String,
    pub satisfaction: f64,
    pub loyalty: f64,
    pub ego: f64,
    pub competitive_drive: f64,
    pub city_connection: f64,
    pub last_match_id: Option<String>,
}

#[derive(Debug, Clone, PartialEq)]
pub struct MatchAgentReactions {
    pub core_agent_changes: Vec<AgentStateChange>,
    pub roster_patch: Option<RosterPatchEnvelope>,
    pub relationship_changes: Vec<AgentRelationshipChange>,
}

#[derive(Debug, Clone, PartialEq)]
pub struct AgentRelationship {
    pub game_id: String,
    pub agent_a_id: String,
    pub agent_b_id: String,
    pub trust: f64,
    pub last_event: String,
    pub trend: String,
    pub short_history: Vec<String>,
    pub last_source_event_id: Option<String>,
}

#[derive(Debug, Clone, PartialEq)]
pub struct AgentRelationshipChange {
    pub relationship: AgentRelationship,
    pub event: AgentRelationshipChangedEvent,
}

#[derive(Debug, Clone, PartialEq)]
pub struct AgentRelationshipSeed {
    pub agent_a_id: &'static str,
    pub agent_b_id: &'static str,
    pub trust: f64,
    pub last_event: &'static str,
    pub trend: &'static str,
    pub short_history: Vec<&'static str>,
}

#[derive(Debug, Clone, PartialEq)]
pub struct IndividualAgentTemplate {
    pub agent_id: &'static str,
    pub display_name: &'static str,
    pub category: &'static str,
    pub role: &'static str,
    pub domain: &'static str,
    pub emotional_state: &'static str,
    pub confidence: f64,
    pub satisfaction: f64,
    pub loyalty: f64,
    pub role_performance: f64,
    pub state: BTreeMap<String, f64>,
    pub agenda: BTreeMap<String, String>,
}

#[derive(Debug, Clone, PartialEq)]
pub struct IndividualAgentState {
    pub game_id: String,
    pub agent_id: String,
    pub display_name: String,
    pub category: String,
    pub role: String,
    pub domain: String,
    pub emotional_state: String,
    pub confidence: f64,
    pub satisfaction: f64,
    pub loyalty: f64,
    pub role_performance: f64,
    pub state: BTreeMap<String, f64>,
    pub agenda: BTreeMap<String, String>,
}

#[must_use]
pub fn default_individual_agent_states(game_id: &str) -> Vec<IndividualAgentState> {
    individual_agent_templates()
        .into_iter()
        .map(|template| template.into_state(game_id))
        .collect()
}

#[must_use]
pub fn default_player_agent_state(player: &TeamRosterPlayer) -> PlayerAgentState {
    let rating_factor = (f64::from(player.overall_rating).clamp(60.0, 90.0) - 60.0) / 30.0;
    let ego = match player.position.as_str() {
        "PG" | "SG" | "SF" => 0.42 + rating_factor * 0.24,
        _ => 0.34 + rating_factor * 0.20,
    };

    PlayerAgentState {
        game_id: player.game_id.clone(),
        player_id: player.player_id.clone(),
        full_name: player.full_name.clone(),
        position: player.position.clone(),
        emotional_state: "steady".to_string(),
        satisfaction: 0.04,
        loyalty: 0.62,
        ego: clamp_unit(ego),
        competitive_drive: clamp_unit(0.58 + rating_factor * 0.22),
        city_connection: 0.35,
        last_match_id: None,
    }
}

#[must_use]
pub fn default_agent_relationships(game_id: &str) -> Vec<AgentRelationship> {
    canonical_relationship_seeds()
        .into_iter()
        .map(|seed| seed.into_relationship(game_id))
        .collect()
}

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
pub fn relationship_key(agent_a_id: &str, agent_b_id: &str) -> String {
    if agent_a_id <= agent_b_id {
        format!("{agent_a_id}:{agent_b_id}")
    } else {
        format!("{agent_b_id}:{agent_a_id}")
    }
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
    fn into_state(self, game_id: &str) -> IndividualAgentState {
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
    fn into_relationship(self, game_id: &str) -> AgentRelationship {
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

#[must_use]
pub fn default_core_agent_states(game_id: &str) -> Vec<CoreAgentState> {
    CORE_AGENT_IDS
        .iter()
        .map(|agent_id| default_core_agent_state(game_id, agent_id))
        .collect()
}

#[must_use]
pub fn default_core_agent_state(game_id: &str, agent_id: &str) -> CoreAgentState {
    let state = match agent_id {
        "owner" => map_from_pairs(&[
            ("sporting_trust", 0.0),
            ("business_trust", 0.0),
            ("patience_remaining", 0.75),
            ("satisfaction", 0.0),
        ]),
        "head_coach" => map_from_pairs(&[
            ("gm_trust", 0.0),
            ("roster_satisfaction", 0.0),
            ("results_pressure", 0.25),
            ("locker_room_relationship", 0.0),
        ]),
        "cfo" => map_from_pairs(&[
            ("financial_trust", 0.0),
            ("budget_alert", 0.15),
            ("financial_conservatism", 0.55),
        ]),
        "scouting_director" => map_from_pairs(&[
            ("criteria_trust", 0.0),
            ("motivation", 0.3),
            ("perceived_precision", 0.0),
        ]),
        "sports_psychologist" => map_from_pairs(&[
            ("locker_room_climate", 0.0),
            ("emotional_alert", 0.2),
            ("player_trust", 0.0),
        ]),
        _ => BTreeMap::new(),
    };

    CoreAgentState {
        game_id: game_id.to_string(),
        agent_id: agent_id.to_string(),
        mood: "calm".to_string(),
        state,
        last_match_id: None,
    }
}

#[must_use]
pub fn apply_match_finished(
    current_states: Vec<CoreAgentState>,
    event: &MatchFinishedEvent,
    occurred_at: String,
) -> Vec<AgentStateChange> {
    let context = MatchContext::from_event(event);

    current_states
        .into_iter()
        .map(|state| apply_match_to_agent(state, event, &context, &occurred_at))
        .collect()
}

#[must_use]
pub fn apply_salary_cap_to_core_agents(
    current_states: Vec<CoreAgentState>,
    event: &SalaryCapCalculatedEvent,
    occurred_at: String,
) -> Vec<AgentStateChange> {
    let mut changes = Vec::new();
    for mut state in current_states {
        match state.agent_id.as_str() {
            "cfo" => {
                if event.status == "luxury_tax" {
                    adjust(&mut state.state, "budget_alert", 0.18);
                    adjust(&mut state.state, "financial_trust", -0.04);
                    state.mood = "concerned".to_string();
                } else if event.near_luxury_tax {
                    adjust(&mut state.state, "budget_alert", 0.08);
                    state.mood = "watchful".to_string();
                } else {
                    adjust(&mut state.state, "budget_alert", -0.02);
                    state.mood = "calm".to_string();
                }
            }
            "owner" => {
                if event.status == "luxury_tax" {
                    adjust(&mut state.state, "patience_remaining", -0.04);
                    adjust(&mut state.state, "business_trust", -0.03);
                    state.mood = "concerned".to_string();
                } else {
                    continue;
                }
            }
            _ => continue,
        }

        changes.push(AgentStateChange {
            event: AgentStateChangedEvent {
                meta: EventMeta {
                    event_id: format!("agent-state-{}-{}", event.meta.event_id, state.agent_id),
                    game_id: event.meta.game_id.clone(),
                    occurred_at: occurred_at.clone(),
                    schema_version: SCHEMA_VERSION,
                },
                simulated_date: event.simulated_date.clone(),
                agent_id: state.agent_id.clone(),
                mood: state.mood.clone(),
                state: state.state.clone(),
                summary: salary_cap_summary(&state.agent_id, event),
                source_event_id: event.meta.event_id.clone(),
                source_subject: "salary_cap.calculado".to_string(),
            },
            state,
        });
    }

    changes
}

fn salary_cap_summary(agent_id: &str, event: &SalaryCapCalculatedEvent) -> String {
    match (agent_id, event.status.as_str()) {
        ("cfo", "luxury_tax") => "El CFO alerta que la nomina entro en luxury tax.".to_string(),
        ("cfo", _) if event.near_luxury_tax => {
            "El CFO marca que la franquicia esta cerca de la linea de luxury tax.".to_string()
        }
        ("owner", "luxury_tax") => {
            "El Owner reduce paciencia por el costo proyectado de la nomina.".to_string()
        }
        _ => "La situacion financiera se mantiene controlada.".to_string(),
    }
}

#[must_use]
pub fn apply_match_to_player_agents(
    current_states: Vec<PlayerAgentState>,
    event: &MatchFinishedEvent,
) -> (Vec<PlayerAgentState>, Option<RosterPatchEnvelope>) {
    if current_states.is_empty() {
        return (Vec::new(), None);
    }

    let context = MatchContext::from_event(event);
    let mut next_states = Vec::with_capacity(current_states.len());
    let mut patches = Vec::with_capacity(current_states.len());

    for mut state in current_states {
        let line = event
            .box_score
            .iter()
            .find(|line| line.player_id == state.player_id && line.team_id == OWN_TEAM_ID);
        let Some(line) = line else {
            next_states.push(state);
            continue;
        };

        apply_box_score_to_player(&mut state, line, &context, &event.match_id);
        patches.push(player_patch(&state, line, &context));
        next_states.push(state);
    }

    let roster_patch = if patches.is_empty() {
        None
    } else {
        Some(RosterPatchEnvelope {
            event_type: SUBJECT_ROSTER_PATCH.to_string(),
            subject: SUBJECT_ROSTER_PATCH.to_string(),
            game_id: event.meta.game_id.clone(),
            patch: RosterStatePatch {
                simulated_date: event.simulated_date.clone(),
                source_event_id: event.meta.event_id.clone(),
                source_subject: SUBJECT_MATCH_FINISHED.to_string(),
                players: patches,
            },
        })
    };

    (next_states, roster_patch)
}

#[must_use]
pub fn apply_match_to_relationships(
    current_relationships: Vec<AgentRelationship>,
    event: &MatchFinishedEvent,
    occurred_at: String,
) -> Vec<AgentRelationshipChange> {
    let context = MatchContext::from_event(event);
    let mut changes = Vec::new();

    for mut relationship in current_relationships {
        let Some((delta, reason)) = relationship_delta_for_match(&relationship, &context) else {
            continue;
        };
        apply_relationship_delta(&mut relationship, delta, reason, &event.meta.event_id);

        changes.push(AgentRelationshipChange {
            event: relationship_changed_event(&relationship, event, &occurred_at),
            relationship,
        });
    }

    changes
}

#[must_use]
pub fn apply_gm_decision_to_relationships(
    current_relationships: Vec<AgentRelationship>,
    event: &GMDecisionRegisteredEvent,
    occurred_at: String,
) -> Vec<AgentRelationshipChange> {
    if event.kind != "medical_decision" {
        return Vec::new();
    }

    let Some(choice_id) = event.payload.get("choice_id") else {
        return Vec::new();
    };
    let mut changes = Vec::new();
    for mut relationship in current_relationships {
        let Some((delta, reason)) =
            relationship_delta_for_medical_decision(&relationship, choice_id)
        else {
            continue;
        };
        apply_relationship_delta(&mut relationship, delta, reason, &event.meta.event_id);
        changes.push(AgentRelationshipChange {
            event: relationship_changed_event_from_decision(&relationship, event, &occurred_at),
            relationship,
        });
    }

    changes
}

fn apply_match_to_agent(
    mut state: CoreAgentState,
    event: &MatchFinishedEvent,
    context: &MatchContext,
    occurred_at: &str,
) -> AgentStateChange {
    match state.agent_id.as_str() {
        "owner" => apply_owner(&mut state, context),
        "head_coach" => apply_head_coach(&mut state, context),
        "cfo" => apply_cfo(&mut state, context),
        "scouting_director" => apply_scouting_director(&mut state, context),
        "sports_psychologist" => apply_sports_psychologist(&mut state, context),
        _ => {}
    }

    state.last_match_id = Some(event.match_id.clone());
    let summary = summarize_agent_change(&state.agent_id, context);
    let event = AgentStateChangedEvent {
        meta: EventMeta {
            event_id: format!(
                "agent-state-{}-{}-{}",
                event.meta.game_id, event.match_id, state.agent_id
            ),
            game_id: event.meta.game_id.clone(),
            occurred_at: occurred_at.to_string(),
            schema_version: SCHEMA_VERSION,
        },
        simulated_date: event.simulated_date.clone(),
        agent_id: state.agent_id.clone(),
        source_event_id: event.meta.event_id.clone(),
        source_subject: SUBJECT_MATCH_FINISHED.to_string(),
        mood: state.mood.clone(),
        state: state.state.clone(),
        summary,
    };

    AgentStateChange { state, event }
}

fn apply_owner(state: &mut CoreAgentState, context: &MatchContext) {
    let result_delta = if context.won { 1.0 } else { -1.0 };
    adjust(&mut state.state, "sporting_trust", 0.06 * result_delta);
    adjust(
        &mut state.state,
        "business_trust",
        if context.home_game {
            0.03 * result_delta
        } else {
            0.015 * result_delta
        },
    );
    adjust(
        &mut state.state,
        "patience_remaining",
        if context.won { 0.02 } else { -0.05 },
    );
    adjust(&mut state.state, "satisfaction", 0.05 * result_delta);

    if context.blowout {
        adjust(&mut state.state, "sporting_trust", 0.03 * result_delta);
        adjust(&mut state.state, "satisfaction", 0.02 * result_delta);
    }

    let patience = metric(&state.state, "patience_remaining");
    state.mood = if !context.won && patience < 0.35 {
        "frustrated"
    } else if !context.won {
        "concerned"
    } else if context.blowout {
        "excited"
    } else {
        "calm"
    }
    .to_string();
}

fn apply_head_coach(state: &mut CoreAgentState, context: &MatchContext) {
    let result_delta = if context.won { 1.0 } else { -1.0 };
    adjust(&mut state.state, "gm_trust", 0.03 * result_delta);
    adjust(&mut state.state, "roster_satisfaction", 0.05 * result_delta);
    adjust(
        &mut state.state,
        "results_pressure",
        if context.won { -0.04 } else { 0.06 },
    );
    adjust(
        &mut state.state,
        "locker_room_relationship",
        0.035 * result_delta,
    );

    if context.close_game {
        adjust(&mut state.state, "results_pressure", 0.015);
    }
    if context.blowout {
        adjust(&mut state.state, "results_pressure", -0.02 * result_delta);
    }

    let pressure = metric(&state.state, "results_pressure");
    state.mood = if pressure > 0.65 {
        "pressured"
    } else if !context.won {
        "frustrated"
    } else {
        "calm"
    }
    .to_string();
}

fn apply_cfo(state: &mut CoreAgentState, context: &MatchContext) {
    let result_delta = if context.won { 1.0 } else { -1.0 };
    let home_multiplier = if context.home_game { 1.0 } else { 0.6 };
    adjust(
        &mut state.state,
        "financial_trust",
        0.02 * result_delta * home_multiplier,
    );
    adjust(
        &mut state.state,
        "budget_alert",
        if context.won {
            -0.01 * home_multiplier
        } else {
            0.025 * home_multiplier
        },
    );
    adjust(
        &mut state.state,
        "financial_conservatism",
        if context.won { -0.01 } else { 0.02 },
    );

    state.mood = if metric(&state.state, "budget_alert") > 0.55 {
        "concerned"
    } else {
        "calm"
    }
    .to_string();
}

fn apply_scouting_director(state: &mut CoreAgentState, context: &MatchContext) {
    let result_delta = if context.won { 1.0 } else { -1.0 };
    adjust(&mut state.state, "criteria_trust", 0.025 * result_delta);
    adjust(
        &mut state.state,
        "motivation",
        if context.close_game { 0.025 } else { 0.015 },
    );
    adjust(
        &mut state.state,
        "perceived_precision",
        if context.won { 0.02 } else { -0.015 },
    );

    if !context.won && context.blowout {
        adjust(&mut state.state, "criteria_trust", -0.02);
        adjust(&mut state.state, "perceived_precision", -0.015);
    }

    state.mood = if !context.won && context.blowout {
        "concerned"
    } else if context.won {
        "excited"
    } else {
        "calm"
    }
    .to_string();
}

fn apply_sports_psychologist(state: &mut CoreAgentState, context: &MatchContext) {
    let result_delta = if context.won { 1.0 } else { -1.0 };
    adjust(
        &mut state.state,
        "locker_room_climate",
        0.055 * result_delta,
    );
    adjust(
        &mut state.state,
        "emotional_alert",
        if context.won { -0.03 } else { 0.055 },
    );
    adjust(&mut state.state, "player_trust", 0.03 * result_delta);

    if context.close_game {
        adjust(&mut state.state, "emotional_alert", 0.015);
    }

    state.mood = if metric(&state.state, "emotional_alert") > 0.55 {
        "concerned"
    } else if context.won {
        "calm"
    } else {
        "pressured"
    }
    .to_string();
}

fn apply_box_score_to_player(
    state: &mut PlayerAgentState,
    line: &PlayerBoxScore,
    context: &MatchContext,
    match_id: &str,
) {
    let performance = player_performance_index(line);
    let result_delta = if context.won { 1.0 } else { -1.0 };
    let role_pressure = if line.minutes >= 28 { 0.018 } else { -0.01 };

    state.satisfaction = clamp(adjusted(
        state.satisfaction,
        0.035 * result_delta + performance * 0.035,
    ));
    state.loyalty = clamp_unit(state.loyalty + 0.012 * result_delta + performance * 0.01);
    state.ego = clamp_unit(state.ego + performance * 0.035 + role_pressure);
    state.competitive_drive =
        clamp_unit(state.competitive_drive + if context.won { 0.008 } else { 0.02 });
    state.city_connection = clamp_unit(
        state.city_connection
            + if context.home_game && context.won {
                0.018
            } else {
                0.004
            },
    );
    state.last_match_id = Some(match_id.to_string());

    state.emotional_state = if !context.won && performance < -0.35 {
        "frustrated"
    } else if !context.won {
        "restless"
    } else if performance > 0.40 {
        "confident"
    } else {
        "steady"
    }
    .to_string();
}

fn relationship_delta_for_match(
    relationship: &AgentRelationship,
    context: &MatchContext,
) -> Option<(f64, &'static str)> {
    let key = relationship_key(&relationship.agent_a_id, &relationship.agent_b_id);
    match key.as_str() {
        "head_analytics:head_coach" => {
            if context.won {
                Some((
                    0.015,
                    "El resultado reduce la friccion entre modelo y rotacion.",
                ))
            } else if context.blowout {
                Some((
                    -0.035,
                    "La derrota amplia reabre la tension entre datos y decisiones de cancha.",
                ))
            } else {
                Some((
                    -0.012,
                    "La derrota deja mas preguntas sobre la lectura tactica.",
                ))
            }
        }
        "head_coach:sports_psychologist" => {
            if context.won {
                Some((
                    0.012,
                    "La victoria baja la tension sobre el manejo emocional del vestuario.",
                ))
            } else {
                Some((
                    -0.025,
                    "La derrota aumenta la tension entre bienestar y exigencia competitiva.",
                ))
            }
        }
        "mayor:owner" => {
            if context.home_game && context.won {
                Some((
                    0.012,
                    "Una victoria local mejora la lectura civica del proyecto.",
                ))
            } else if context.home_game {
                Some((
                    -0.018,
                    "Una derrota local enfria el valor politico del proyecto.",
                ))
            } else {
                None
            }
        }
        "gm:pr_director" => {
            if context.blowout && !context.won {
                Some((
                    -0.025,
                    "La derrota amplia complica la narrativa publica del GM.",
                ))
            } else if context.won {
                Some((
                    0.01,
                    "La victoria hace mas defendible la direccion publica del GM.",
                ))
            } else {
                None
            }
        }
        "press:roster_collective" => {
            if context.blowout && !context.won {
                Some((
                    -0.03,
                    "La cobertura se endurece sobre el estado emocional del roster.",
                ))
            } else if context.won {
                Some((
                    0.015,
                    "La cobertura positiva reduce presion sobre el vestuario.",
                ))
            } else {
                Some((
                    -0.01,
                    "La derrota sostiene una cobertura mas incomoda para el roster.",
                ))
            }
        }
        _ => None,
    }
}

fn relationship_delta_for_medical_decision(
    relationship: &AgentRelationship,
    choice_id: &str,
) -> Option<(f64, &'static str)> {
    let key = relationship_key(&relationship.agent_a_id, &relationship.agent_b_id);
    match (key.as_str(), choice_id) {
        ("gm:team_doctor", "rest") => Some((
            0.045,
            "El GM respeta el protocolo medico y fortalece la confianza del staff de salud.",
        )),
        ("gm:team_doctor", "reduce_minutes") => Some((
            0.025,
            "El GM acepta bajar carga y el Medico interpreta la decision como prudente.",
        )),
        ("gm:team_doctor", "ignore_doctor") => Some((
            -0.055,
            "El GM ignora la recomendacion medica y erosiona la confianza del staff de salud.",
        )),
        ("gm:team_doctor", "force_return") => Some((
            -0.09,
            "El GM fuerza una alta anticipada y abre una fractura seria con el Medico.",
        )),
        ("head_coach:team_doctor", "force_return") => Some((
            -0.035,
            "La urgencia competitiva vuelve a tensionar al Coach con el Medico.",
        )),
        ("head_coach:team_doctor", "rest") => Some((
            0.018,
            "Coach y Medico quedan alineados alrededor del protocolo de recuperacion.",
        )),
        ("head_coach:team_doctor", "reduce_minutes") => Some((
            0.012,
            "La reduccion de carga crea un compromiso aceptable entre competencia y salud.",
        )),
        _ => None,
    }
}

fn apply_relationship_delta(
    relationship: &mut AgentRelationship,
    delta: f64,
    reason: &str,
    source_event_id: &str,
) {
    relationship.trust = clamp(relationship.trust + delta);
    relationship.trend = if delta > 0.0 {
        "improving"
    } else if delta < 0.0 {
        "deteriorating"
    } else {
        "stable"
    }
    .to_string();
    relationship.last_event = reason.to_string();
    relationship.last_source_event_id = Some(source_event_id.to_string());
    relationship.short_history.push(reason.to_string());
    if relationship.short_history.len() > 5 {
        relationship.short_history.remove(0);
    }
}

fn relationship_changed_event(
    relationship: &AgentRelationship,
    event: &MatchFinishedEvent,
    occurred_at: &str,
) -> AgentRelationshipChangedEvent {
    AgentRelationshipChangedEvent {
        meta: EventMeta {
            event_id: format!(
                "agent-relationship-{}-{}",
                event.meta.event_id,
                relationship_key(&relationship.agent_a_id, &relationship.agent_b_id)
            ),
            game_id: event.meta.game_id.clone(),
            occurred_at: occurred_at.to_string(),
            schema_version: SCHEMA_VERSION,
        },
        simulated_date: event.simulated_date.clone(),
        agent_a_id: relationship.agent_a_id.clone(),
        agent_b_id: relationship.agent_b_id.clone(),
        trust: relationship.trust,
        trend: relationship.trend.clone(),
        last_event: relationship.last_event.clone(),
        short_history: relationship.short_history.clone(),
        source_event_id: event.meta.event_id.clone(),
        source_subject: SUBJECT_MATCH_FINISHED.to_string(),
    }
}

fn relationship_changed_event_from_decision(
    relationship: &AgentRelationship,
    event: &GMDecisionRegisteredEvent,
    occurred_at: &str,
) -> AgentRelationshipChangedEvent {
    AgentRelationshipChangedEvent {
        meta: EventMeta {
            event_id: format!(
                "agent-relationship-{}-{}",
                event.meta.event_id,
                relationship_key(&relationship.agent_a_id, &relationship.agent_b_id)
            ),
            game_id: event.meta.game_id.clone(),
            occurred_at: occurred_at.to_string(),
            schema_version: SCHEMA_VERSION,
        },
        simulated_date: event.simulated_date.clone(),
        agent_a_id: relationship.agent_a_id.clone(),
        agent_b_id: relationship.agent_b_id.clone(),
        trust: relationship.trust,
        trend: relationship.trend.clone(),
        last_event: relationship.last_event.clone(),
        short_history: relationship.short_history.clone(),
        source_event_id: event.meta.event_id.clone(),
        source_subject: SUBJECT_GM_DECISION_REGISTERED.to_string(),
    }
}

fn player_patch(
    state: &PlayerAgentState,
    line: &PlayerBoxScore,
    context: &MatchContext,
) -> PlayerEmotionalPatch {
    PlayerEmotionalPatch {
        player_id: state.player_id.clone(),
        emotional_state: state.emotional_state.clone(),
        satisfaction: state.satisfaction,
        loyalty: state.loyalty,
        ego: state.ego,
        competitive_drive: state.competitive_drive,
        city_connection: state.city_connection,
        summary: summarize_player_change(state, line, context),
    }
}

fn summarize_player_change(
    state: &PlayerAgentState,
    line: &PlayerBoxScore,
    context: &MatchContext,
) -> String {
    let result = if context.won { "victoria" } else { "derrota" };
    let role = if line.minutes >= 28 {
        "rol alto"
    } else if line.minutes >= 16 {
        "rotacion estable"
    } else {
        "minutos limitados"
    };

    format!(
        "{} procesa la {} con {} y {} puntos.",
        state.full_name, result, role, line.points
    )
}

fn player_performance_index(line: &PlayerBoxScore) -> f64 {
    let production = f64::from(line.points) * 0.04
        + f64::from(line.rebounds) * 0.025
        + f64::from(line.assists) * 0.03
        + f64::from(line.steals + line.blocks) * 0.04
        - f64::from(line.turnovers) * 0.04;
    let minutes_expectation = f64::from(line.minutes) * 0.018;

    (production - minutes_expectation).clamp(-1.0, 1.0)
}

fn summarize_agent_change(agent_id: &str, context: &MatchContext) -> String {
    let result = if context.won { "victoria" } else { "derrota" };
    let venue = if context.home_game {
        "en casa"
    } else {
        "como visitante"
    };
    let intensity = if context.blowout {
        "amplia"
    } else if context.close_game {
        "cerrada"
    } else {
        "normal"
    };

    match agent_id {
        "owner" => format!("El owner ajusta su confianza tras una {result} {intensity} {venue}."),
        "head_coach" => format!("El head coach recalibra presion y satisfaccion tras la {result}."),
        "cfo" => format!("El CFO actualiza su lectura financiera despues de la {result} {venue}."),
        "scouting_director" => {
            format!(
                "El scouting director reevalua su criterio despues de una {result} {intensity}."
            )
        }
        "sports_psychologist" => {
            format!("La sports psychologist ajusta el clima emocional tras la {result}.")
        }
        _ => format!("El agente reacciona al resultado del partido."),
    }
}

fn adjust(state: &mut BTreeMap<String, f64>, key: &str, delta: f64) {
    let current = metric(state, key);
    state.insert(key.to_string(), adjusted(current, delta));
}

fn adjusted(current: f64, delta: f64) -> f64 {
    clamp(current + delta)
}

fn metric(state: &BTreeMap<String, f64>, key: &str) -> f64 {
    state.get(key).copied().unwrap_or_default()
}

fn clamp(value: f64) -> f64 {
    value.clamp(MIN_STATE_VALUE, MAX_STATE_VALUE)
}

fn clamp_unit(value: f64) -> f64 {
    value.clamp(0.0, 1.0)
}

fn map_from_pairs(pairs: &[(&str, f64)]) -> BTreeMap<String, f64> {
    pairs
        .iter()
        .map(|(key, value)| ((*key).to_string(), *value))
        .collect()
}

fn string_map_from_pairs(pairs: &[(&str, &str)]) -> BTreeMap<String, String> {
    pairs
        .iter()
        .map(|(key, value)| ((*key).to_string(), (*value).to_string()))
        .collect()
}

#[derive(Debug, Clone, Copy, PartialEq, Eq)]
struct MatchContext {
    won: bool,
    home_game: bool,
    close_game: bool,
    blowout: bool,
}

impl MatchContext {
    fn from_event(event: &MatchFinishedEvent) -> Self {
        let own_score = if event.home_team.team_id == OWN_TEAM_ID {
            event.home_score
        } else {
            event.away_score
        };
        let opponent_score = if event.home_team.team_id == OWN_TEAM_ID {
            event.away_score
        } else {
            event.home_score
        };
        let margin = own_score.abs_diff(opponent_score);

        Self {
            won: event.winner_team_id == OWN_TEAM_ID,
            home_game: event.home_team.team_id == OWN_TEAM_ID,
            close_game: margin <= 5,
            blowout: margin >= 15,
        }
    }
}

#[cfg(test)]
mod tests {
    use super::*;
    use crate::events::{EventMeta, MatchTeam};

    #[test]
    fn defaults_create_the_five_m2_core_agents() {
        let states = default_core_agent_states("game-1");

        assert_eq!(states.len(), 5);
        assert!(states.iter().any(|state| state.agent_id == "owner"));
        assert!(states.iter().any(|state| state.agent_id == "head_coach"));
        assert!(states.iter().any(|state| state.agent_id == "cfo"));
        assert!(
            states
                .iter()
                .any(|state| state.agent_id == "scouting_director")
        );
        assert!(
            states
                .iter()
                .any(|state| state.agent_id == "sports_psychologist")
        );
    }

    #[test]
    fn individual_catalog_seeds_thirty_canon_agents() {
        let agents = default_individual_agent_states("game-1");

        assert_eq!(agents.len(), INDIVIDUAL_AGENT_COUNT);
        assert!(agents.iter().any(|agent| agent.agent_id == "owner"));
        assert!(agents.iter().any(|agent| agent.agent_id == "head_coach"));
        assert!(agents.iter().any(|agent| agent.agent_id == "cfo"));
        assert!(
            agents
                .iter()
                .any(|agent| agent.agent_id == "scouting_director")
        );
        assert!(
            agents
                .iter()
                .any(|agent| agent.agent_id == "sports_psychologist")
        );
        assert!(agents.iter().any(|agent| agent.agent_id == "mayor"));
        assert!(agents.iter().any(|agent| agent.agent_id == "press"));
    }

    #[test]
    fn individual_catalog_has_unique_ids_and_universal_state() {
        let agents = default_individual_agent_states("game-1");
        let mut ids = std::collections::BTreeSet::new();

        for agent in agents {
            assert!(ids.insert(agent.agent_id));
            assert!((-1.0..=1.0).contains(&agent.confidence));
            assert!((-1.0..=1.0).contains(&agent.satisfaction));
            assert!((0.0..=1.0).contains(&agent.loyalty));
            assert!((0.0..=1.0).contains(&agent.role_performance));
            assert!(!agent.domain.is_empty());
            assert!(!agent.agenda.is_empty());
            assert!(!agent.state.is_empty());
        }
    }

    #[test]
    fn default_player_agent_state_uses_team_player_id() {
        let player = sample_roster_player("game-1-player-01", 82, "PG");
        let state = default_player_agent_state(&player);

        assert_eq!(state.player_id, "game-1-player-01");
        assert_eq!(state.game_id, "game-1");
        assert_eq!(state.full_name, "Mateo Cross");
        assert!((0.0..=1.0).contains(&state.ego));
        assert!((0.0..=1.0).contains(&state.competitive_drive));
    }

    #[test]
    fn player_agents_react_to_box_score_with_roster_patch() {
        let mut event = sample_match(true, 112, 101);
        event.box_score = vec![PlayerBoxScore {
            player_id: "game-1-player-01".to_string(),
            team_id: OWN_TEAM_ID.to_string(),
            minutes: 32,
            points: 26,
            rebounds: 6,
            assists: 8,
            steals: 1,
            blocks: 0,
            turnovers: 2,
        }];
        let initial =
            default_player_agent_state(&sample_roster_player("game-1-player-01", 82, "PG"));

        let (states, patch) = apply_match_to_player_agents(vec![initial], &event);

        assert_eq!(states.len(), 1);
        assert_eq!(states[0].last_match_id.as_deref(), Some("match-1"));
        assert!(states[0].satisfaction > 0.04);

        let patch = patch.expect("roster patch exists");
        assert_eq!(patch.event_type, SUBJECT_ROSTER_PATCH);
        assert_eq!(patch.patch.players.len(), 1);
        assert_eq!(patch.patch.players[0].player_id, "game-1-player-01");
    }

    #[test]
    fn relationship_catalog_seeds_canon_tensions() {
        let relationships = default_agent_relationships("game-1");

        assert!(
            relationships
                .iter()
                .any(|relationship| relationship.agent_a_id == "head_coach"
                    && relationship.agent_b_id == "team_doctor")
        );
        assert!(relationships.iter().any(
            |relationship| relationship.agent_a_id == "cfo" && relationship.agent_b_id == "gm"
        ));
        assert!(
            relationships
                .iter()
                .any(|relationship| relationship.agent_a_id == "press"
                    && relationship.agent_b_id == "roster_collective")
        );
        assert!(
            relationships
                .iter()
                .any(|relationship| relationship.agent_a_id == "gm"
                    && relationship.agent_b_id == "team_doctor")
        );
    }

    #[test]
    fn medical_decision_moves_doctor_gm_relationship() {
        let relationships = default_agent_relationships("game-1");
        let event = GMDecisionRegisteredEvent {
            meta: EventMeta {
                event_id: "decision-medical-injury-1".to_string(),
                game_id: "game-1".to_string(),
                occurred_at: "2026-10-29T00:00:00Z".to_string(),
                schema_version: SCHEMA_VERSION,
            },
            decision_id: "medical-injury-1".to_string(),
            kind: "medical_decision".to_string(),
            payload: string_map_from_pairs(&[
                ("choice_id", "force_return"),
                ("injury_id", "injury-1"),
                ("player_id", "player-1"),
            ]),
            simulated_date: "2026-10-29".to_string(),
            agents_affected: vec!["team_doctor".to_string(), "head_coach".to_string()],
            source_event_id: Some("injury-1".to_string()),
            source_subject: Some("jugador.lesionado".to_string()),
        };

        let changes = apply_gm_decision_to_relationships(
            relationships,
            &event,
            "2026-10-29T00:00:01Z".to_string(),
        );
        let doctor_gm = changes
            .iter()
            .find(|change| {
                relationship_key(
                    &change.relationship.agent_a_id,
                    &change.relationship.agent_b_id,
                ) == "gm:team_doctor"
            })
            .expect("doctor gm relationship moves");

        assert_eq!(doctor_gm.relationship.trend, "deteriorating");
        assert_eq!(
            doctor_gm.event.source_subject,
            SUBJECT_GM_DECISION_REGISTERED
        );
    }

    #[test]
    fn luxury_tax_moves_cfo_and_owner_state() {
        let states = vec![
            default_core_agent_state("game-1", "owner"),
            default_core_agent_state("game-1", "cfo"),
        ];
        let event = SalaryCapCalculatedEvent {
            meta: EventMeta {
                event_id: "salary-cap-game-1".to_string(),
                game_id: "game-1".to_string(),
                occurred_at: "2026-10-22T00:00:00Z".to_string(),
                schema_version: SCHEMA_VERSION,
            },
            simulated_date: "2026-10-22".to_string(),
            cap_base: 141_000_000,
            luxury_tax_line: 171_000_000,
            committed_salary: 180_000_000,
            cap_space: -39_000_000,
            luxury_tax_space: -9_000_000,
            roster_count: 15,
            status: "luxury_tax".to_string(),
            near_luxury_tax: true,
            projected_tax_payment: 18_000_000,
        };

        let changes =
            apply_salary_cap_to_core_agents(states, &event, "2026-10-22T00:00:01Z".to_string());
        let cfo = state_for(&changes, "cfo");
        let owner = state_for(&changes, "owner");

        assert!(cfo.state["budget_alert"] > 0.15);
        assert!(owner.state["patience_remaining"] < 0.75);
    }

    #[test]
    fn match_result_moves_relevant_relationships() {
        let event = sample_match(false, 88, 111);
        let relationships = default_agent_relationships("game-1");

        let changes =
            apply_match_to_relationships(relationships, &event, "2026-05-25T00:00:00Z".to_string());

        let coach_analytics = changes
            .iter()
            .find(|change| {
                relationship_key(
                    &change.relationship.agent_a_id,
                    &change.relationship.agent_b_id,
                ) == "head_analytics:head_coach"
            })
            .expect("coach analytics relationship moves");
        assert_eq!(coach_analytics.relationship.trend, "deteriorating");
        assert_eq!(coach_analytics.event.source_subject, SUBJECT_MATCH_FINISHED);
    }

    #[test]
    fn win_improves_owner_and_reduces_coach_pressure() {
        let event = sample_match(true, 112, 101);
        let changes = apply_match_finished(
            default_core_agent_states("game-1"),
            &event,
            "2026-05-24T00:00:00Z".to_string(),
        );

        let owner = state_for(&changes, "owner");
        assert!(owner.state["sporting_trust"] > 0.0);
        assert!(owner.state["patience_remaining"] > 0.75);
        assert_eq!(owner.mood, "calm");

        let coach = state_for(&changes, "head_coach");
        assert!(coach.state["results_pressure"] < 0.25);
    }

    #[test]
    fn blowout_loss_increases_pressure_and_emotional_alert() {
        let event = sample_match(false, 88, 111);
        let changes = apply_match_finished(
            default_core_agent_states("game-1"),
            &event,
            "2026-05-24T00:00:00Z".to_string(),
        );

        let coach = state_for(&changes, "head_coach");
        assert!(coach.state["results_pressure"] > 0.25);

        let psychologist = state_for(&changes, "sports_psychologist");
        assert!(psychologist.state["emotional_alert"] > 0.2);
    }

    #[test]
    fn state_change_events_are_deterministic_per_match_and_agent() {
        let event = sample_match(true, 112, 101);
        let changes = apply_match_finished(
            default_core_agent_states("game-1"),
            &event,
            "2026-05-24T00:00:00Z".to_string(),
        );

        let coach = changes
            .iter()
            .find(|change| change.state.agent_id == "head_coach")
            .expect("coach change exists");

        assert_eq!(
            coach.event.meta.event_id,
            "agent-state-game-1-match-1-head_coach"
        );
        assert_eq!(coach.event.source_subject, "partido.terminado");
        assert_eq!(coach.event.source_event_id, "match-finished-match-1");
    }

    fn state_for<'a>(changes: &'a [AgentStateChange], agent_id: &str) -> &'a CoreAgentState {
        &changes
            .iter()
            .find(|change| change.state.agent_id == agent_id)
            .expect("agent state exists")
            .state
    }

    fn sample_match(own_home: bool, own_score: u16, opponent_score: u16) -> MatchFinishedEvent {
        let own = MatchTeam {
            team_id: OWN_TEAM_ID.to_string(),
            name: "PulseCity".to_string(),
            abbreviation: "PUL".to_string(),
        };
        let opponent = MatchTeam {
            team_id: "rival-1".to_string(),
            name: "Rival".to_string(),
            abbreviation: "RIV".to_string(),
        };
        let (home_team, away_team, home_score, away_score) = if own_home {
            (own, opponent, own_score, opponent_score)
        } else {
            (opponent, own, opponent_score, own_score)
        };
        let winner_team_id = if own_score > opponent_score {
            OWN_TEAM_ID
        } else {
            "rival-1"
        };

        MatchFinishedEvent {
            meta: EventMeta {
                event_id: "match-finished-match-1".to_string(),
                game_id: "game-1".to_string(),
                occurred_at: "2026-05-24T00:00:00Z".to_string(),
                schema_version: 1,
            },
            match_id: "match-1".to_string(),
            simulated_date: "2026-10-22".to_string(),
            home_team,
            away_team,
            home_score,
            away_score,
            winner_team_id: winner_team_id.to_string(),
            seed: 123,
            box_score: Vec::new(),
        }
    }

    fn sample_roster_player(
        player_id: &str,
        overall_rating: u8,
        position: &str,
    ) -> TeamRosterPlayer {
        TeamRosterPlayer {
            player_id: player_id.to_string(),
            game_id: "game-1".to_string(),
            full_name: "Mateo Cross".to_string(),
            position: position.to_string(),
            overall_rating,
            roster_status: "active".to_string(),
        }
    }
}
