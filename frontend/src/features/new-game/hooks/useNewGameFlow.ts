import { useCallback, useEffect, useRef, useState } from "react";

import type {
  AgentClientStates,
  ChatClientMessages,
  ChatMessageEvent,
  CityClientState,
  FinanceClientState,
  GameSetup,
  GameSummary,
  GuestSession,
  GuestUpgradeResult,
  MapClientState,
  NarrativeChoice,
  NarrativeEvent,
  RealtimeEvent,
  RelationshipClientStates,
  RosterClientStates,
  SeasonClientState,
  SeasonMatchSummary,
  TimeClientState,
  UserSession,
} from "../../../types";
import { clampPage, pageFromPath, relationshipId } from "../helpers";
import {
  flowPages,
  initialDraft,
  pagePaths,
  stageMeta,
  type CityManagementModeId,
  type FlowPage,
  type ScenarioId,
} from "../constants";
import type { NewGameDraft } from "../types";

const gatewayBaseUrl = "http://localhost:8080";
const socketBaseUrl = "ws://localhost:8080/ws";
const guestTokenStorageKey = "pulsecity_guest_token";
const sessionTokenStorageKey = "pulsecity_session_token";

const initialMapState: MapClientState = {
  game_id: "",
  stage: "idle",
  progress: 0,
  message: "Esperando la orden de fundacion.",
};

const initialTimeState: TimeClientState = {
  simulated_date: "2026-10-01",
  speed: 1,
  paused: true,
  days_processed: 0,
};

const initialSeasonState: SeasonClientState = {
  wins: 0,
  losses: 0,
  points_for: 0,
  points_against: 0,
};

const initialCityState: CityClientState = {
  fan_sentiment: 50,
  ticket_sales_index: 50,
  local_economy_index: 50,
  stadium_district_land_value: 100,
  win_streak: 0,
  loss_streak: 0,
};

const initialFinanceState: FinanceClientState = {
  simulated_date: "2026-10-01",
  source_event_id: "",
  source_subject: "",
  cap_base: 0,
  luxury_tax_line: 0,
  committed_salary: 0,
  cap_space: 0,
  luxury_tax_space: 0,
  roster_count: 0,
  status: "under_cap",
  near_luxury_tax: false,
  projected_tax_payment: 0,
};

const initialAgentStates: AgentClientStates = {};
const initialRosterStates: RosterClientStates = {};
const initialRelationshipStates: RelationshipClientStates = {};
const initialChatMessages: ChatClientMessages = {};

export function useNewGameFlow() {
  const [draft, setDraft] = useState<NewGameDraft>(initialDraft);
  const [currentPage, setCurrentPage] = useState<FlowPage>("session");
  const [unlockedPage, setUnlockedPage] = useState<FlowPage>("home");
  const [gameId, setGameId] = useState("");
  const [selectedGameId, setSelectedGameId] = useState("");
  const [guestToken, setGuestToken] = useState("");
  const [userSession, setUserSession] = useState<UserSession | null>(null);
  const [games, setGames] = useState<GameSummary[]>([]);
  const [status, setStatus] = useState("Elegí cómo entrar a PulseCity.");
  const [socketStatus, setSocketStatus] = useState("Socket inactivo");
  const [mapState, setMapState] = useState<MapClientState>(initialMapState);
  const [timeState, setTimeState] = useState<TimeClientState>(initialTimeState);
  const [seasonState, setSeasonState] = useState<SeasonClientState>(initialSeasonState);
  const [recentResults, setRecentResults] = useState<SeasonMatchSummary[]>([]);
  const [cityState, setCityState] = useState<CityClientState>(initialCityState);
  const [financeState, setFinanceState] = useState<FinanceClientState>(initialFinanceState);
  const [agentStates, setAgentStates] = useState<AgentClientStates>(initialAgentStates);
  const [rosterStates, setRosterStates] = useState<RosterClientStates>(initialRosterStates);
  const [relationshipStates, setRelationshipStates] =
    useState<RelationshipClientStates>(initialRelationshipStates);
  const [chatMessages, setChatMessages] = useState<ChatClientMessages>(initialChatMessages);
  const [events, setEvents] = useState<RealtimeEvent[]>([]);
  const [narrativeInbox, setNarrativeInbox] = useState<NarrativeEvent[]>([]);
  const [activeNarrativeEvent, setActiveNarrativeEvent] = useState<NarrativeEvent | null>(null);
  const [ownerIntroResponse, setOwnerIntroResponse] = useState<NarrativeChoice | null>(null);
  const [submittingNarrativeChoice, setSubmittingNarrativeChoice] = useState(false);
  const [creatingGame, setCreatingGame] = useState(false);
  const [creatingGuestSession, setCreatingGuestSession] = useState(false);
  const [authenticatingUser, setAuthenticatingUser] = useState(false);
  const [restoringSession, setRestoringSession] = useState(true);
  const socketRef = useRef<WebSocket | null>(null);
  const activeGameIdRef = useRef("");

  const activeAuthKind: "none" | "guest" | "user" = userSession
    ? "user"
    : guestToken
      ? "guest"
      : "none";

  const syncPage = useCallback((nextPage: FlowPage, replace = false) => {
    if (replace) {
      window.history.replaceState({}, "", pagePaths[nextPage]);
    } else {
      window.history.pushState({}, "", pagePaths[nextPage]);
    }
    window.scrollTo({ top: 0, behavior: "smooth" });
    setCurrentPage(nextPage);
  }, []);

  useEffect(() => {
    activeGameIdRef.current = gameId;
  }, [gameId]);

  useEffect(() => {
    if (games.length === 0) {
      setSelectedGameId("");
      return;
    }

    if (!games.some((game) => game.game_id === selectedGameId)) {
      setSelectedGameId(games[0].game_id);
    }
  }, [games, selectedGameId]);

  useEffect(() => {
    void restoreAuthState();
  }, []);

  useEffect(() => {
    const requested = pageFromPath(window.location.pathname);
    const safePage = resolveAccessiblePage(requested, unlockedPage, activeAuthKind);
    if (safePage !== requested) {
      window.history.replaceState({}, "", pagePaths[safePage]);
    }
    setCurrentPage(safePage);

    const handlePopState = () => {
      const rawPage = pageFromPath(window.location.pathname);
      const nextPage = resolveAccessiblePage(rawPage, unlockedPage, activeAuthKind);
      if (nextPage !== rawPage) {
        window.history.replaceState({}, "", pagePaths[nextPage]);
      }
      setCurrentPage(nextPage);
    };

    window.addEventListener("popstate", handlePopState);
    return () => {
      window.removeEventListener("popstate", handlePopState);
    };
  }, [activeAuthKind, unlockedPage]);

  useEffect(() => {
    return () => {
      socketRef.current?.close();
    };
  }, []);

  useEffect(() => {
    if (!gameId || activeAuthKind === "none") {
      return;
    }

    void loadGameSetup(gameId);
  }, [gameId, activeAuthKind, userSession, guestToken]);

  const applyEvent = useCallback((payload: RealtimeEvent) => {
    const payloadGameId = eventGameId(payload);
    if (payloadGameId && payloadGameId !== activeGameIdRef.current) {
      return;
    }

    if (payload.type === "narrative.event") {
      if (payload.kind !== "owner_intro") {
        setNarrativeInbox((current) => prependNarrativeEvent(current, payload));
        setStatus(`${payload.emitter}: ${payload.title}`);
        return;
      }
      if (ownerIntroResponse) {
        return;
      }
      setActiveNarrativeEvent(payload);
      return;
    }

    if (payload.type === "narrative.response") {
      setOwnerIntroResponse(payload.choice);
      setActiveNarrativeEvent(null);
      setStatus(`Direccion inicial confirmada: ${payload.choice.label}.`);
      return;
    }

    if (payload.type === "chat.message") {
      setChatMessages((current) => appendChatMessage(current, payload));
      setStatus(`${agentLabel(payload.agent_id)} respondió en el chat.`);
      return;
    }

    if (payload.type === "map.snapshot") {
      setMapState(payload.state);
      setGameId(payload.state.game_id);
      return;
    }

    if (payload.type === "time.patch") {
      setTimeState((current) => ({
        simulated_date: payload.patch.simulated_date ?? current.simulated_date,
        speed: payload.patch.speed ?? current.speed,
        paused: payload.patch.paused ?? current.paused,
        days_processed: payload.patch.days_processed ?? current.days_processed,
      }));
      return;
    }

    if (payload.type === "season.patch") {
      setSeasonState((current) => ({
        wins: payload.patch.wins ?? current.wins,
        losses: payload.patch.losses ?? current.losses,
        points_for: payload.patch.points_for ?? current.points_for,
        points_against: payload.patch.points_against ?? current.points_against,
        last_result: payload.patch.last_result ?? current.last_result,
      }));
      if (payload.patch.last_result) {
        setRecentResults((current) => prependMatchResult(current, payload.patch.last_result!));
      }
      return;
    }

    if (payload.type === "finance.patch") {
      setFinanceState(payload.patch);
      return;
    }

    if (payload.type === "city.patch") {
      setCityState((current) => ({
        fan_sentiment: payload.patch.fan_sentiment ?? current.fan_sentiment,
        ticket_sales_index: payload.patch.ticket_sales_index ?? current.ticket_sales_index,
        local_economy_index: payload.patch.local_economy_index ?? current.local_economy_index,
        stadium_district_land_value:
          payload.patch.stadium_district_land_value ?? current.stadium_district_land_value,
        win_streak: payload.patch.win_streak ?? current.win_streak,
        loss_streak: payload.patch.loss_streak ?? current.loss_streak,
        last_match_id: payload.patch.last_match_id ?? current.last_match_id,
        reason: payload.patch.reason ?? current.reason,
      }));
      return;
    }

    if (payload.type === "agent.patch") {
      setAgentStates((current) => {
        const previous = current[payload.agent_id];
        return {
          ...current,
          [payload.agent_id]: {
            agent_id: payload.agent_id,
            mood: payload.patch.mood ?? previous?.mood ?? "calm",
            state: {
              ...(previous?.state ?? {}),
              ...(payload.patch.state ?? {}),
            },
            summary: payload.patch.summary ?? previous?.summary ?? "",
            simulated_date: payload.patch.simulated_date ?? previous?.simulated_date,
            source_event_id: payload.patch.source_event_id ?? previous?.source_event_id,
            source_subject: payload.patch.source_subject ?? previous?.source_subject,
          },
        };
      });
      return;
    }

    if (payload.type === "roster.patch") {
      setRosterStates((current) => {
        const next = { ...current };
        for (const player of payload.patch.players) {
          next[player.player_id] = {
            ...current[player.player_id],
            ...player,
            simulated_date: payload.patch.simulated_date,
            source_event_id: payload.patch.source_event_id,
            source_subject: payload.patch.source_subject,
          };
        }
        return next;
      });
      return;
    }

    if (payload.type === "relations.patch") {
      setRelationshipStates((current) => {
        const next = { ...current };
        for (const relationship of payload.patch.relationships) {
          const id = relationshipId(relationship.agent_a_id, relationship.agent_b_id);
          next[id] = {
            relationship_id: id,
            ...relationship,
            simulated_date: payload.patch.simulated_date,
            source_event_id: payload.patch.source_event_id,
            source_subject: payload.patch.source_subject,
          };
        }
        return next;
      });
      return;
    }

    setMapState((current) => ({
      ...current,
      game_id: payload.game_id,
      stage: payload.patch.stage ?? current.stage,
      progress: payload.patch.progress ?? current.progress,
      message: payload.patch.message ?? current.message,
      map_data: payload.patch.map_data ?? current.map_data,
      stadium: payload.patch.stadium ?? current.stadium,
    }));
  }, [ownerIntroResponse]);

  const connectSocket = useCallback((nextGameId: string) => {
    socketRef.current?.close();

    const params = new URLSearchParams();
    if (nextGameId) {
      params.set("game_id", nextGameId);
    }
    if (userSession?.session_token) {
      params.set("session_token", userSession.session_token);
    } else if (guestToken) {
      params.set("guest_token", guestToken);
    }

    const query = params.toString();
    const socket = new WebSocket(query ? `${socketBaseUrl}?${query}` : socketBaseUrl);
    socketRef.current = socket;

    socket.addEventListener("open", () => {
      setSocketStatus(nextGameId ? `Socket activo para ${nextGameId}` : "Socket activo");
    });

    socket.addEventListener("close", () => {
      setSocketStatus("Socket desconectado");
    });

    socket.addEventListener("error", () => {
      setSocketStatus("Error en WebSocket");
    });

    socket.addEventListener("message", (event) => {
      const payload = JSON.parse(event.data) as RealtimeEvent;
      if (!isEventForGame(payload, activeGameIdRef.current)) {
        return;
      }
      applyEvent(payload);
      setEvents((current) => [payload, ...current].slice(0, 12));
    });
  }, [applyEvent, guestToken, userSession]);

  function updateDraft<K extends keyof NewGameDraft>(key: K, value: NewGameDraft[K]) {
    setDraft((current) => ({
      ...current,
      [key]: value,
    }));
  }

  function startNewGame() {
    if (activeAuthKind === "none") {
      setStatus("Primero necesitás entrar como invitado o iniciar sesión.");
      return;
    }

    setUnlockedPage("identity");
    setStatus("Nueva partida iniciada. Primero definí la identidad de la franquicia.");
    syncPage("identity");
  }

  function goBack() {
    switch (currentPage) {
      case "scenario":
        syncPage("identity");
        return;
      case "management":
        syncPage("scenario");
        return;
      case "launch":
        syncPage("management");
        return;
      default:
        syncPage("home");
    }
  }

  function completeIdentityStep() {
    setUnlockedPage("scenario");
    setStatus("Identidad confirmada. Ahora elegí desde dónde arranca la historia.");
    syncPage("scenario");
  }

  function completeScenarioStep() {
    setUnlockedPage("management");
    setStatus("Escenario inicial confirmado. Ahora definí cómo gobernás la ciudad.");
    syncPage("management");
  }

  function completeManagementStep() {
    setUnlockedPage("launch");
    setStatus("Modo de gestión confirmado. Revisá todo antes de fundar el mundo.");
    syncPage("launch");
  }

  async function restoreAuthState() {
    setRestoringSession(true);
    const storedGuestToken = window.localStorage.getItem(guestTokenStorageKey) ?? "";
    const storedSessionToken = window.localStorage.getItem(sessionTokenStorageKey) ?? "";

    if (storedGuestToken) {
      setGuestToken(storedGuestToken);
    }

    if (storedSessionToken) {
      const restored = await restoreUserSession(storedSessionToken);
      if (restored) {
        syncPage("home", true);
        setRestoringSession(false);
        return;
      }
    }

    if (storedGuestToken) {
      await loadGamesForGuest(storedGuestToken, { silentOnSuccess: true });
      setStatus("Sesión invitada restaurada.");
      syncPage("home", true);
    } else {
      setStatus("Elegí cómo entrar a PulseCity.");
      syncPage("session", true);
    }

    setRestoringSession(false);
  }

  async function restoreUserSession(sessionToken: string) {
    try {
      const response = await fetch(`${gatewayBaseUrl}/api/v1/auth/session`, {
        headers: {
          "X-Session-Token": sessionToken,
        },
      });
      const payload = (await response.json()) as UserSession | { error?: string };
      if (!response.ok || !("session_token" in payload)) {
        clearUserSession({ preserveGuestFallback: true, resetNavigation: false });
        return false;
      }

      setUserSession(payload);
      await loadGamesForUser(payload.session_token, { silentOnSuccess: true });
      setStatus(`Sesión restaurada para ${payload.user.display_name}.`);
      return true;
    } catch {
      clearUserSession({ preserveGuestFallback: true, resetNavigation: false });
      return false;
    }
  }

  async function register(email: string, displayName: string, password: string) {
    if (authenticatingUser) {
      return;
    }

    const normalizedEmail = email.trim().toLowerCase();
    const normalizedDisplayName = displayName.trim();
    const registrationError = validateRegistrationInput(
      normalizedEmail,
      normalizedDisplayName,
      password,
    );
    if (registrationError) {
      setStatus(registrationError);
      return;
    }

    setAuthenticatingUser(true);
    setStatus("Creando cuenta...");

    try {
      const response = await fetch(`${gatewayBaseUrl}/api/v1/auth/register`, {
        method: "POST",
        headers: {
          "Content-Type": "application/json",
        },
        body: JSON.stringify({
          email: normalizedEmail,
          display_name: normalizedDisplayName,
          password,
        }),
      });
      const payload = (await response.json()) as UserSession | { error?: string };
      if (!response.ok || !("session_token" in payload)) {
        setStatus(("error" in payload && payload.error) || "No se pudo crear la cuenta.");
        return;
      }

      window.localStorage.setItem(sessionTokenStorageKey, payload.session_token);
      setUserSession(payload);
      const upgradeMessage = await upgradeGuestOwnership(payload);
      if (!upgradeMessage) {
        setStatus(`Cuenta creada. Sesión iniciada como ${payload.user.display_name}.`);
      }
      await loadGamesForUser(payload.session_token);
      syncPage("home", true);
      if (upgradeMessage) {
        setStatus(upgradeMessage);
      }
    } catch (error) {
      setStatus(error instanceof Error ? error.message : "Fallo de red al crear la cuenta.");
    } finally {
      setAuthenticatingUser(false);
    }
  }

  async function login(email: string, password: string) {
    if (authenticatingUser) {
      return;
    }

    const normalizedEmail = email.trim().toLowerCase();
    const loginError = validateLoginInput(normalizedEmail, password);
    if (loginError) {
      setStatus(loginError);
      return;
    }

    setAuthenticatingUser(true);
    setStatus("Iniciando sesión...");

    try {
      const response = await fetch(`${gatewayBaseUrl}/api/v1/auth/login`, {
        method: "POST",
        headers: {
          "Content-Type": "application/json",
        },
        body: JSON.stringify({
          email: normalizedEmail,
          password,
        }),
      });
      const payload = (await response.json()) as UserSession | { error?: string };
      if (!response.ok || !("session_token" in payload)) {
        setStatus(("error" in payload && payload.error) || "No se pudo iniciar sesión.");
        return;
      }

      window.localStorage.setItem(sessionTokenStorageKey, payload.session_token);
      setUserSession(payload);
      const upgradeMessage = await upgradeGuestOwnership(payload);
      if (!upgradeMessage) {
        setStatus(`Sesión iniciada como ${payload.user.display_name}.`);
      }
      await loadGamesForUser(payload.session_token);
      syncPage("home", true);
      if (upgradeMessage) {
        setStatus(upgradeMessage);
      }
    } catch (error) {
      setStatus(error instanceof Error ? error.message : "Fallo de red al iniciar sesión.");
    } finally {
      setAuthenticatingUser(false);
    }
  }

  async function createGuestSession() {
    if (creatingGuestSession) {
      return;
    }

    setCreatingGuestSession(true);
    setStatus("Abriendo sesión invitada...");

    try {
      const response = await fetch(`${gatewayBaseUrl}/api/v1/guest-sessions`, {
        method: "POST",
      });
      const payload = (await response.json()) as GuestSession | { error?: string };
      if (!response.ok || !("guest_token" in payload)) {
        setStatus(("error" in payload && payload.error) || "No se pudo crear la sesión invitada.");
        return;
      }

      window.localStorage.setItem(guestTokenStorageKey, payload.guest_token);
      setGuestToken(payload.guest_token);
      if (!userSession) {
        setGames([]);
        setStatus("Sesión invitada lista. Ya podés crear o continuar una partida.");
        setUnlockedPage("home");
        syncPage("home", true);
        await loadGamesForGuest(payload.guest_token);
      } else {
        setStatus("Sesión invitada guardada. La sesión activa sigue siendo la cuenta autenticada.");
      }
    } catch (error) {
      setStatus(
        error instanceof Error ? error.message : "Fallo de red al crear la sesión invitada.",
      );
    } finally {
      setCreatingGuestSession(false);
    }
  }

  async function createGame() {
    if (activeAuthKind === "none") {
      setStatus("Primero necesitás una sesión activa.");
      return;
    }

    setCreatingGame(true);
    setStatus(`Fundando ${draft.cityName} ${draft.franchiseName}...`);
    setEvents([]);
    setRecentResults([]);
    setNarrativeInbox([]);
    setMapState(initialMapState);
    setTimeState(initialTimeState);
    setSeasonState(initialSeasonState);
    setCityState(initialCityState);
    setFinanceState(initialFinanceState);
    setAgentStates(initialAgentStates);
    setRosterStates(initialRosterStates);
    setRelationshipStates(initialRelationshipStates);
    setChatMessages(initialChatMessages);
    setOwnerIntroResponse(null);
    setActiveNarrativeEvent(null);

    try {
      const response = await fetch(`${gatewayBaseUrl}/api/v1/games`, {
        method: "POST",
        headers: buildAuthHeaders({
          guestToken,
          sessionToken: userSession?.session_token,
          includeContentType: true,
        }),
        body: JSON.stringify({
          city_name: draft.cityName,
          franchise_name: draft.franchiseName,
          abbreviation: draft.abbreviation,
          primary_color: draft.primaryColor,
          secondary_color: draft.secondaryColor,
          accent_color: draft.accentColor,
          initial_scenario: draft.selectedScenario,
          city_management_mode: draft.cityManagementMode,
        }),
      });

      const payload = (await response.json()) as { game_id?: string; error?: string };
      if (!response.ok || !payload.game_id) {
        setStatus(payload.error ?? "No se pudo crear la partida.");
        return;
      }

      setGameId(payload.game_id);
      setUnlockedPage("ceremony");
      setStatus(`Partida creada para ${draft.cityName}. Ceremonia en curso.`);
      syncPage("ceremony");
      connectSocket(payload.game_id);
      await loadGames();
    } catch (error) {
      setStatus(
        error instanceof Error ? error.message : "Fallo de red al crear la partida.",
      );
    } finally {
      setCreatingGame(false);
    }
  }

  async function continueGame(nextGameId: string) {
    if (activeAuthKind === "none") {
      setStatus("La sesión actual ya no está disponible.");
      return;
    }

    setStatus(`Recuperando la partida ${nextGameId}...`);
    setUnlockedPage("ceremony");
    setGameId(nextGameId);
    syncPage("ceremony");
    connectSocket(nextGameId);
  }

  function continueSelectedGame() {
    if (!selectedGameId) {
      setStatus("Elegí una partida para continuar.");
      return;
    }

    void continueGame(selectedGameId);
  }

  async function switchToGuestSession() {
    if (!guestToken) {
      setStatus("No hay una sesión invitada guardada para activar.");
      return;
    }

    resetRuntimeToEntry();
    window.localStorage.removeItem(sessionTokenStorageKey);
    setUserSession(null);
    setStatus("Volviendo al perfil invitado...");
    syncPage("home", true);
    await loadGamesForGuest(guestToken);
  }

  async function logoutUser() {
    if (!userSession) {
      return;
    }

    resetRuntimeToEntry();
    window.localStorage.removeItem(sessionTokenStorageKey);
    setUserSession(null);

    if (guestToken) {
      setStatus("Sesión de cuenta cerrada. Volviendo al invitado guardado...");
      syncPage("home", true);
      await loadGamesForGuest(guestToken);
      return;
    }

    setStatus("Sesión cerrada. Elegí cómo entrar a PulseCity.");
  }

  function clearAllAccess() {
    resetRuntimeToEntry();
    clearGuestSessionStorage();
    window.localStorage.removeItem(sessionTokenStorageKey);
    setGuestToken("");
    setUserSession(null);
    setGames([]);
    setSelectedGameId("");
    setStatus("Sesiones locales limpiadas. Elegí cómo entrar a PulseCity.");
    syncPage("session", true);
  }

  function forgotPassword() {
    setStatus("Recuperación de contraseña todavía no implementada en esta etapa.");
  }

  async function upgradeGuestOwnership(nextSession: UserSession) {
    if (!guestToken) {
      return "";
    }

    try {
      const response = await fetch(`${gatewayBaseUrl}/api/v1/auth/upgrade-guest`, {
        method: "POST",
        headers: {
          "X-Guest-Token": guestToken,
          "X-Session-Token": nextSession.session_token,
        },
      });
      const payload = (await response.json()) as GuestUpgradeResult | { error?: string };
      if (!response.ok || !("migrated_games" in payload)) {
        setStatus(
          ("error" in payload && payload.error) ||
            "La cuenta se creó, pero no se pudieron migrar las partidas invitadas.",
        );
        return "";
      }

      clearGuestSessionStorage();
      setGuestToken("");
      return payload.migrated_games > 0
        ? `Sesión iniciada como ${nextSession.user.display_name}. Se migraron ${payload.migrated_games} partida(s) del invitado actual.`
        : `Sesión iniciada como ${nextSession.user.display_name}. No había partidas guest para migrar.`;
    } catch {
      setStatus("La cuenta quedó autenticada, pero falló la migración de partidas guest.");
      return "";
    }
  }

  async function loadGames() {
    if (userSession?.session_token) {
      return loadGamesForUser(userSession.session_token);
    }
    if (guestToken) {
      return loadGamesForGuest(guestToken);
    }

    return false;
  }

  async function loadGamesForGuest(
    nextGuestToken: string,
    options?: { silentOnSuccess?: boolean },
  ) {
    try {
      const response = await fetch(`${gatewayBaseUrl}/api/v1/games`, {
        headers: buildAuthHeaders({
          guestToken: nextGuestToken,
          includeContentType: false,
        }),
      });
      const payload = (await response.json()) as { games?: GameSummary[]; error?: string };
      if (!response.ok || !payload.games) {
        if (response.status === 401 && !userSession) {
          clearGuestSession({ resetNavigation: false });
          setStatus("La sesión invitada ya no es válida. Creá una nueva para continuar.");
        }
        return false;
      }

      setGames(payload.games);
      if (!options?.silentOnSuccess) {
        setStatus(
          payload.games.length === 0
            ? "Sesión invitada lista. Todavía no hay partidas asociadas."
            : `Sesión invitada lista. ${payload.games.length} partida(s) asociada(s).`,
        );
      }
      return true;
    } catch {
      return false;
    }
  }

  async function loadGamesForUser(
    sessionToken: string,
    options?: { silentOnSuccess?: boolean },
  ) {
    try {
      const response = await fetch(`${gatewayBaseUrl}/api/v1/games`, {
        headers: buildAuthHeaders({
          sessionToken,
          includeContentType: false,
        }),
      });
      const payload = (await response.json()) as { games?: GameSummary[]; error?: string };
      if (!response.ok || !payload.games) {
        if (response.status === 401) {
          clearUserSession({ preserveGuestFallback: true, resetNavigation: false });
          setStatus("La sesión autenticada dejó de ser válida.");
        }
        return false;
      }

      setGames(payload.games);
      if (!options?.silentOnSuccess) {
        setStatus(
          payload.games.length === 0
            ? "Sesión autenticada lista. Todavía no hay partidas asociadas."
            : `Sesión autenticada lista. ${payload.games.length} partida(s) asociada(s).`,
        );
      }
      return true;
    } catch {
      return false;
    }
  }

  async function loadGameSetup(nextGameId: string) {
    try {
      const response = await fetch(`${gatewayBaseUrl}/api/v1/games/${nextGameId}`, {
        headers: buildAuthHeaders({
          guestToken,
          sessionToken: userSession?.session_token,
          includeContentType: false,
        }),
      });
      const payload = (await response.json()) as GameSetup | { error?: string };
      if (!response.ok || !("game_id" in payload)) {
        if (response.status === 401) {
          if (userSession) {
            clearUserSession({ preserveGuestFallback: true, resetNavigation: true });
            setStatus("La sesión autenticada dejó de ser válida.");
          } else {
            clearGuestSession({ resetNavigation: true });
            setStatus("La sesión invitada dejó de ser válida.");
          }
        }
        return false;
      }

      applyGameSetup(payload);
      return true;
    } catch {
      return false;
    }
  }

  function applyGameSetup(setup: GameSetup) {
    setDraft({
      cityName: setup.city_name,
      franchiseName: setup.franchise_name,
      abbreviation: setup.abbreviation,
      primaryColor: setup.primary_color,
      secondaryColor: setup.secondary_color,
      accentColor: setup.accent_color,
      selectedScenario: isScenarioId(setup.initial_scenario)
        ? setup.initial_scenario
        : initialDraft.selectedScenario,
      cityManagementMode: isCityManagementModeId(setup.city_management_mode)
        ? setup.city_management_mode
        : initialDraft.cityManagementMode,
    });
    setOwnerIntroResponse(setup.owner_intro_response ?? null);
    if (setup.owner_intro_event && !setup.owner_intro_response) {
      setActiveNarrativeEvent(setup.owner_intro_event);
    } else if (setup.owner_intro_response) {
      setActiveNarrativeEvent(null);
    }
  }

  async function submitOwnerIntroChoice(choice: NarrativeChoice) {
    if (!gameId || activeAuthKind === "none" || submittingNarrativeChoice) {
      return;
    }

    setSubmittingNarrativeChoice(true);
    setStatus(`Confirmando direccion inicial: ${choice.label}...`);

    try {
      const response = await fetch(`${gatewayBaseUrl}/api/v1/games/${gameId}/owner-intro-response`, {
        method: "POST",
        headers: buildAuthHeaders({
          guestToken,
          sessionToken: userSession?.session_token,
          includeContentType: true,
        }),
        body: JSON.stringify({
          choice_id: choice.id,
        }),
      });
      const payload = (await response.json()) as NarrativeChoice | { error?: string };
      if (!response.ok || !("id" in payload)) {
        setStatus(("error" in payload && payload.error) || "No se pudo registrar la respuesta.");
        return;
      }

      setOwnerIntroResponse(payload);
      setActiveNarrativeEvent(null);
      setStatus(`Direccion inicial confirmada: ${payload.label}.`);
      setEvents((current) => [
        {
          type: "narrative.response" as const,
          subject: "narrativa.respuesta_gm",
          game_id: gameId,
          event_id: activeNarrativeEvent?.event_id ?? "owner-intro",
          choice: payload,
          emitter: "gm",
          timestamp: new Date().toISOString(),
        },
        ...current,
      ].slice(0, 12));
      await loadGames();
    } catch (error) {
      setStatus(
        error instanceof Error ? error.message : "Fallo de red al enviar la respuesta inicial.",
      );
    } finally {
      setSubmittingNarrativeChoice(false);
    }
  }

  async function updateTimeControl(next: { speed?: 1 | 5 | 20; paused?: boolean }) {
    if (!gameId || activeAuthKind === "none") {
      setStatus("Necesitás una partida activa para controlar el tiempo.");
      return;
    }

    setTimeState((current) => ({
      ...current,
      speed: next.speed ?? current.speed,
      paused: next.paused ?? current.paused,
    }));

    try {
      const response = await fetch(`${gatewayBaseUrl}/api/v1/games/${gameId}/time-control`, {
        method: "POST",
        headers: buildAuthHeaders({
          guestToken,
          sessionToken: userSession?.session_token,
          includeContentType: true,
        }),
        body: JSON.stringify(next),
      });
      const payload = (await response.json()) as { error?: string };
      if (!response.ok) {
        setStatus(payload.error ?? "No se pudo cambiar el tiempo.");
      }
    } catch (error) {
      setStatus(error instanceof Error ? error.message : "Fallo de red al cambiar el tiempo.");
    }
  }

  async function sendAgentChatMessage(agentId: string, message: string, conversationId?: string) {
    const trimmedAgentId = agentId.trim();
    const trimmedMessage = message.trim();
    if (!gameId || activeAuthKind === "none") {
      setStatus("Necesitás una partida activa para hablar con agentes.");
      return "";
    }
    if (!trimmedAgentId || !trimmedMessage) {
      setStatus("Elegí un agente y escribí un mensaje.");
      return "";
    }

    const nextConversationId = conversationId || `chat-local-${trimmedAgentId}`;
    setChatMessages((current) =>
      appendChatMessage(current, {
        type: "chat.message",
        subject: "agente.consulta_iniciada",
        game_id: gameId,
        conversation_id: nextConversationId,
        message_id: `local-${Date.now()}`,
        agent_id: trimmedAgentId,
        sender: "gm",
        body: trimmedMessage,
        created_at: new Date().toISOString(),
      }),
    );
    setStatus(`Consultando a ${agentLabel(trimmedAgentId)}...`);

    try {
      const response = await fetch(`${gatewayBaseUrl}/api/v1/games/${gameId}/agent-chat`, {
        method: "POST",
        headers: buildAuthHeaders({
          guestToken,
          sessionToken: userSession?.session_token,
          includeContentType: true,
        }),
        body: JSON.stringify({
          agent_id: trimmedAgentId,
          message: trimmedMessage,
          conversation_id: nextConversationId,
        }),
      });
      const payload = (await response.json()) as { conversation_id?: string; error?: string };
      if (!response.ok || !payload.conversation_id) {
        setStatus(payload.error ?? "No se pudo iniciar el chat.");
        return "";
      }

      return payload.conversation_id;
    } catch (error) {
      setStatus(error instanceof Error ? error.message : "Fallo de red al enviar el chat.");
      return "";
    }
  }

  function clearGuestSession(options?: { resetNavigation?: boolean }) {
    clearGuestSessionStorage();
    setGuestToken("");
    if (!userSession) {
      resetRuntimeToEntry();
      if (options?.resetNavigation !== false) {
        syncPage("session", true);
      }
    }
  }

  function clearUserSession(options?: { preserveGuestFallback?: boolean; resetNavigation?: boolean }) {
    window.localStorage.removeItem(sessionTokenStorageKey);
    setUserSession(null);
    if (options?.preserveGuestFallback && guestToken) {
      resetRuntimeToEntry();
      syncPage("home", true);
      void loadGamesForGuest(guestToken, { silentOnSuccess: true });
      return;
    }

    resetRuntimeToEntry();
    if (options?.resetNavigation !== false) {
      syncPage("session", true);
    }
  }

  function resetRuntimeToEntry() {
    socketRef.current?.close();
    setGames([]);
    setSelectedGameId("");
    setUnlockedPage("home");
    setGameId("");
    setMapState(initialMapState);
    setTimeState(initialTimeState);
    setSeasonState(initialSeasonState);
    setRecentResults([]);
    setCityState(initialCityState);
    setFinanceState(initialFinanceState);
    setAgentStates(initialAgentStates);
    setRosterStates(initialRosterStates);
    setRelationshipStates(initialRelationshipStates);
    setChatMessages(initialChatMessages);
    setEvents([]);
    setNarrativeInbox([]);
    setActiveNarrativeEvent(null);
    setOwnerIntroResponse(null);
    setSocketStatus("Socket inactivo");
  }

  const currentStage = stageMeta[mapState.stage] ?? stageMeta.idle;
  const selectedGame = games.find((game) => game.game_id === selectedGameId) ?? null;

  return {
    activeAuthKind,
    activeNarrativeEvent,
    authenticatingUser,
    creatingGame,
    creatingGuestSession,
    currentPage,
    currentStage,
    cityState,
    chatMessages,
    agentStates,
    draft,
    events,
    financeState,
    gameId,
    games,
    guestToken,
    mapState,
    narrativeInbox,
    ownerIntroResponse,
    recentResults,
    relationshipStates,
    rosterStates,
    restoringSession,
    selectedGame,
    selectedGameId,
    seasonState,
    socketStatus,
    status,
    submittingNarrativeChoice,
    timeState,
    userSession,
    updateDraft,
    continueSelectedGame,
    continueGame,
    createGame,
    createGuestSession,
    clearAllAccess,
    completeIdentityStep,
    completeManagementStep,
    completeScenarioStep,
    forgotPassword,
    goBack,
    login,
    logoutUser,
    register,
    setSelectedGameId,
    startNewGame,
    switchToGuestSession,
    submitOwnerIntroChoice,
    updateTimeControl,
    sendAgentChatMessage,
  };
}

function clearGuestSessionStorage() {
  window.localStorage.removeItem(guestTokenStorageKey);
}

function prependMatchResult(current: SeasonMatchSummary[], next: SeasonMatchSummary) {
  return [next, ...current.filter((result) => result.match_id !== next.match_id)].slice(0, 8);
}

function prependNarrativeEvent(current: NarrativeEvent[], next: NarrativeEvent) {
  return [next, ...current.filter((event) => event.event_id !== next.event_id)].slice(0, 8);
}

function appendChatMessage(current: ChatClientMessages, next: ChatMessageEvent) {
  const messages = current[next.conversation_id] ?? [];
  if (messages.some((message) => message.message_id === next.message_id)) {
    return current;
  }

  return {
    ...current,
    [next.conversation_id]: [...messages, next].slice(-30),
  };
}

function agentLabel(agentId: string) {
  return agentId.replaceAll("_", " ");
}

function buildAuthHeaders({
  guestToken,
  sessionToken,
  includeContentType,
}: {
  guestToken?: string;
  sessionToken?: string;
  includeContentType: boolean;
}) {
  const headers: Record<string, string> = {};
  if (includeContentType) {
    headers["Content-Type"] = "application/json";
  }
  if (sessionToken) {
    headers["X-Session-Token"] = sessionToken;
  } else if (guestToken) {
    headers["X-Guest-Token"] = guestToken;
  }

  return headers;
}

function eventGameId(payload: RealtimeEvent) {
  if (payload.type === "map.snapshot") {
    return payload.state.game_id;
  }

  return payload.game_id;
}

function isEventForGame(payload: RealtimeEvent, gameId: string) {
  const payloadGameId = eventGameId(payload);
  if (!payloadGameId || !gameId) {
    return true;
  }

  return payloadGameId === gameId;
}

function isScenarioId(value: string): value is ScenarioId {
  return ["rebuild", "contention", "decline", "expansion"].includes(value);
}

function isCityManagementModeId(value: string): value is CityManagementModeId {
  return ["owner_influence", "dual_figure"].includes(value);
}

function resolveAccessiblePage(
  requested: FlowPage,
  unlockedPage: FlowPage,
  activeAuthKind: "none" | "guest" | "user",
) {
  if (activeAuthKind === "none") {
    return "session";
  }

  if (requested === "session") {
    return "home";
  }

  const maxUnlockedPage = flowPages.includes(unlockedPage) ? unlockedPage : "home";
  return clampPage(requested, maxUnlockedPage);
}

function validateRegistrationInput(email: string, displayName: string, password: string) {
  if (!email) {
    return "Necesitás un email para crear la cuenta.";
  }
  if (!/^[^\s@]+@[^\s@]+\.[^\s@]+$/.test(email)) {
    return "El formato del email no es válido.";
  }
  if (!displayName) {
    return "Necesitás un nombre visible para la cuenta.";
  }
  if (displayName.length < 3) {
    return "El nombre visible debe tener al menos 3 caracteres.";
  }
  if (displayName.length > 40) {
    return "El nombre visible no puede superar los 40 caracteres.";
  }
  if (!password) {
    return "Necesitás una contraseña para crear la cuenta.";
  }
  if (password.length < 8) {
    return "La contraseña debe tener al menos 8 caracteres.";
  }
  if (password.length > 72) {
    return "La contraseña no puede superar los 72 caracteres.";
  }

  return "";
}

function validateLoginInput(email: string, password: string) {
  if (!email) {
    return "Necesitás un email para iniciar sesión.";
  }
  if (!/^[^\s@]+@[^\s@]+\.[^\s@]+$/.test(email)) {
    return "El formato del email no es válido.";
  }
  if (!password) {
    return "Necesitás una contraseña para iniciar sesión.";
  }

  return "";
}
