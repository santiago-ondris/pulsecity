import { useCallback, useEffect, useRef, useState } from "react";

import type {
  GameSetup,
  GameSummary,
  GuestSession,
  MapClientState,
  NarrativeChoice,
  NarrativeEvent,
  RealtimeEvent,
} from "../../../types";
import { clampPage, pageFromPath } from "../helpers";
import {
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

const initialMapState: MapClientState = {
  game_id: "",
  stage: "idle",
  progress: 0,
  message: "Esperando la orden de fundacion.",
};

export function useNewGameFlow() {
  const [draft, setDraft] = useState<NewGameDraft>(initialDraft);
  const [currentPage, setCurrentPage] = useState<FlowPage>("home");
  const [unlockedPage, setUnlockedPage] = useState<FlowPage>("home");
  const [gameId, setGameId] = useState("");
  const [guestToken, setGuestToken] = useState("");
  const [guestGames, setGuestGames] = useState<GameSummary[]>([]);
  const [status, setStatus] = useState("Elegí cómo entrar a PulseCity.");
  const [socketStatus, setSocketStatus] = useState("Socket inactivo");
  const [mapState, setMapState] = useState<MapClientState>(initialMapState);
  const [events, setEvents] = useState<RealtimeEvent[]>([]);
  const [activeNarrativeEvent, setActiveNarrativeEvent] = useState<NarrativeEvent | null>(null);
  const [ownerIntroResponse, setOwnerIntroResponse] = useState<NarrativeChoice | null>(null);
  const [submittingNarrativeChoice, setSubmittingNarrativeChoice] = useState(false);
  const [creatingGame, setCreatingGame] = useState(false);
  const [creatingGuestSession, setCreatingGuestSession] = useState(false);
  const socketRef = useRef<WebSocket | null>(null);
  const activeGameIdRef = useRef("");

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
    const storedGuestToken = window.localStorage.getItem(guestTokenStorageKey) ?? "";
    if (storedGuestToken) {
      setGuestToken(storedGuestToken);
      void loadGuestGames(storedGuestToken);
    }
  }, []);

  useEffect(() => {
    const requested = pageFromPath(window.location.pathname);
    const maxUnlockedPage = guestToken ? unlockedPage : "home";
    const safePage = clampPage(requested, maxUnlockedPage);
    if (safePage !== requested) {
      window.history.replaceState({}, "", pagePaths[safePage]);
    }
    setCurrentPage(safePage);

    const handlePopState = () => {
      const rawPage = pageFromPath(window.location.pathname);
      const nextPage = clampPage(rawPage, guestToken ? unlockedPage : "home");
      if (nextPage !== rawPage) {
        window.history.replaceState({}, "", pagePaths[nextPage]);
      }
      setCurrentPage(nextPage);
    };

    window.addEventListener("popstate", handlePopState);
    return () => {
      window.removeEventListener("popstate", handlePopState);
    };
  }, [guestToken, unlockedPage]);

  useEffect(() => {
    return () => {
      socketRef.current?.close();
    };
  }, []);

  useEffect(() => {
    if (!gameId || !guestToken) {
      return;
    }

    void loadGameSetup(gameId, guestToken);
  }, [gameId, guestToken]);

  const applyEvent = useCallback((payload: RealtimeEvent) => {
    const payloadGameId = eventGameId(payload);
    if (payloadGameId && payloadGameId !== activeGameIdRef.current) {
      return;
    }

    if (payload.type === "narrative.event") {
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

    if (payload.type === "map.snapshot") {
      setMapState(payload.state);
      setGameId(payload.state.game_id);
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

  const connectSocket = useCallback((nextGameId: string, nextGuestToken: string) => {
    socketRef.current?.close();

    const params = new URLSearchParams();
    if (nextGameId) {
      params.set("game_id", nextGameId);
    }
    if (nextGuestToken) {
      params.set("guest_token", nextGuestToken);
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
  }, [applyEvent]);

  function updateDraft<K extends keyof NewGameDraft>(key: K, value: NewGameDraft[K]) {
    setDraft((current) => ({
      ...current,
      [key]: value,
    }));
  }

  function startNewGame() {
    if (!guestToken) {
      setStatus("Primero necesitás entrar como invitado.");
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
      setGuestGames([]);
      setStatus("Sesión invitada lista. Ya podés crear o continuar una partida.");
      setUnlockedPage("home");
      syncPage("home", true);
      await loadGuestGames(payload.guest_token);
    } catch (error) {
      setStatus(
        error instanceof Error ? error.message : "Fallo de red al crear la sesión invitada.",
      );
    } finally {
      setCreatingGuestSession(false);
    }
  }

  async function createGame() {
    if (!guestToken) {
      setStatus("Primero necesitás una sesión invitada.");
      return;
    }

    setCreatingGame(true);
    setStatus(`Fundando ${draft.cityName} ${draft.franchiseName}...`);
    setEvents([]);
    setMapState(initialMapState);
    setOwnerIntroResponse(null);
    setActiveNarrativeEvent(null);

    try {
      const response = await fetch(`${gatewayBaseUrl}/api/v1/games`, {
        method: "POST",
        headers: buildGuestHeaders(guestToken),
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
      connectSocket(payload.game_id, guestToken);
      await loadGuestGames(guestToken);
    } catch (error) {
      setStatus(
        error instanceof Error ? error.message : "Fallo de red al crear la partida.",
      );
    } finally {
      setCreatingGame(false);
    }
  }

  async function continueGame(nextGameId: string) {
    if (!guestToken) {
      setStatus("La sesión invitada ya no está disponible.");
      return;
    }

    setStatus(`Recuperando la partida ${nextGameId}...`);
    setUnlockedPage("ceremony");
    setGameId(nextGameId);
    syncPage("ceremony");
    connectSocket(nextGameId, guestToken);
  }

  async function loadGuestGames(nextGuestToken: string) {
    try {
      const response = await fetch(`${gatewayBaseUrl}/api/v1/games`, {
        headers: buildGuestHeaders(nextGuestToken),
      });
      const payload = (await response.json()) as { games?: GameSummary[]; error?: string };
      if (!response.ok || !payload.games) {
        if (response.status === 401) {
          clearGuestSession();
          setStatus("La sesión invitada ya no es válida. Creá una nueva para continuar.");
        }
        return false;
      }

      setGuestGames(payload.games);
      if (payload.games.length === 0) {
        setStatus("Sesión invitada lista. Todavía no hay partidas asociadas.");
      } else {
        setStatus(`Sesión invitada lista. ${payload.games.length} partida(s) asociada(s).`);
      }
      return true;
    } catch {
      return false;
    }
  }

  async function loadGameSetup(nextGameId: string, nextGuestToken: string) {
    try {
      const response = await fetch(`${gatewayBaseUrl}/api/v1/games/${nextGameId}`, {
        headers: buildGuestHeaders(nextGuestToken),
      });
      const payload = (await response.json()) as GameSetup | { error?: string };
      if (!response.ok || !("game_id" in payload)) {
        if (response.status === 401) {
          clearGuestSession();
          setStatus("La sesión invitada dejó de ser válida.");
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
    if (!gameId || !guestToken || submittingNarrativeChoice) {
      return;
    }

    setSubmittingNarrativeChoice(true);
    setStatus(`Confirmando direccion inicial: ${choice.label}...`);

    try {
      const response = await fetch(`${gatewayBaseUrl}/api/v1/games/${gameId}/owner-intro-response`, {
        method: "POST",
        headers: buildGuestHeaders(guestToken),
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
      await loadGuestGames(guestToken);
    } catch (error) {
      setStatus(
        error instanceof Error ? error.message : "Fallo de red al enviar la respuesta inicial.",
      );
    } finally {
      setSubmittingNarrativeChoice(false);
    }
  }

  function clearGuestSession() {
    window.localStorage.removeItem(guestTokenStorageKey);
    setGuestToken("");
    setGuestGames([]);
    setUnlockedPage("home");
    setGameId("");
  }

  const currentStage = stageMeta[mapState.stage] ?? stageMeta.idle;

  return {
    activeNarrativeEvent,
    creatingGame,
    creatingGuestSession,
    currentPage,
    currentStage,
    draft,
    events,
    gameId,
    guestGames,
    guestReady: guestToken.length > 0,
    guestToken,
    mapState,
    ownerIntroResponse,
    socketStatus,
    status,
    submittingNarrativeChoice,
    updateDraft,
    continueGame,
    createGame,
    createGuestSession,
    completeIdentityStep,
    completeManagementStep,
    completeScenarioStep,
    goBack,
    startNewGame,
    submitOwnerIntroChoice,
  };
}

function buildGuestHeaders(guestToken: string) {
  return {
    "Content-Type": "application/json",
    "X-Guest-Token": guestToken,
  };
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
