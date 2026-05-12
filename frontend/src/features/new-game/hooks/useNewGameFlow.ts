import { useCallback, useEffect, useRef, useState } from "react";

import type {
  GameSetup,
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
  const [status, setStatus] = useState("Pantalla inicial lista.");
  const [socketStatus, setSocketStatus] = useState("Socket inactivo");
  const [mapState, setMapState] = useState<MapClientState>(initialMapState);
  const [events, setEvents] = useState<RealtimeEvent[]>([]);
  const [activeNarrativeEvent, setActiveNarrativeEvent] = useState<NarrativeEvent | null>(null);
  const [ownerIntroResponse, setOwnerIntroResponse] = useState<NarrativeChoice | null>(null);
  const [submittingNarrativeChoice, setSubmittingNarrativeChoice] = useState(false);
  const [creatingGame, setCreatingGame] = useState(false);
  const socketRef = useRef<WebSocket | null>(null);

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
    const requested = pageFromPath(window.location.pathname);
    const safePage = clampPage(requested, unlockedPage);
    if (safePage !== requested) {
      window.history.replaceState({}, "", pagePaths[safePage]);
    }
    setCurrentPage(safePage);

    const handlePopState = () => {
      const rawPage = pageFromPath(window.location.pathname);
      const nextPage = clampPage(rawPage, unlockedPage);
      if (nextPage !== rawPage) {
        window.history.replaceState({}, "", pagePaths[nextPage]);
      }
      setCurrentPage(nextPage);
    };

    window.addEventListener("popstate", handlePopState);
    return () => {
      window.removeEventListener("popstate", handlePopState);
    };
  }, [unlockedPage]);

  useEffect(() => {
    return () => {
      socketRef.current?.close();
    };
  }, []);

  useEffect(() => {
    if (!gameId) {
      return;
    }

    void loadGameSetup(gameId);
  }, [gameId]);

  const applyEvent = useCallback((payload: RealtimeEvent) => {
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

  const connectSocket = useCallback((nextGameId: string) => {
    socketRef.current?.close();

    const query = nextGameId ? `?game_id=${encodeURIComponent(nextGameId)}` : "";
    const socket = new WebSocket(`${socketBaseUrl}${query}`);
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

  async function createGame() {
    setCreatingGame(true);
    setStatus(`Fundando ${draft.cityName} ${draft.franchiseName}...`);
    setEvents([]);
    setMapState(initialMapState);
    setOwnerIntroResponse(null);
    setActiveNarrativeEvent(null);

    try {
      const response = await fetch(`${gatewayBaseUrl}/api/v1/games`, {
        method: "POST",
        headers: {
          "Content-Type": "application/json",
        },
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
    } catch (error) {
      setStatus(
        error instanceof Error ? error.message : "Fallo de red al crear la partida.",
      );
    } finally {
      setCreatingGame(false);
    }
  }

  async function loadGameSetup(nextGameId: string) {
    try {
      const response = await fetch(`${gatewayBaseUrl}/api/v1/games/${nextGameId}`);
      const payload = (await response.json()) as GameSetup | { error?: string };
      if (!response.ok || !("game_id" in payload)) {
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
    if (!gameId || submittingNarrativeChoice) {
      return;
    }

    setSubmittingNarrativeChoice(true);
    setStatus(`Confirmando direccion inicial: ${choice.label}...`);

    try {
      const response = await fetch(`${gatewayBaseUrl}/api/v1/games/${gameId}/owner-intro-response`, {
        method: "POST",
        headers: {
          "Content-Type": "application/json",
        },
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
    } catch (error) {
      setStatus(
        error instanceof Error ? error.message : "Fallo de red al enviar la respuesta inicial.",
      );
    } finally {
      setSubmittingNarrativeChoice(false);
    }
  }

  const currentStage = stageMeta[mapState.stage] ?? stageMeta.idle;

  return {
    activeNarrativeEvent,
    creatingGame,
    currentPage,
    currentStage,
    draft,
    events,
    gameId,
    mapState,
    ownerIntroResponse,
    socketStatus,
    status,
    submittingNarrativeChoice,
    updateDraft,
    createGame,
    completeIdentityStep,
    completeManagementStep,
    completeScenarioStep,
    goBack,
    startNewGame,
    submitOwnerIntroChoice,
  };
}

function isScenarioId(value: string): value is ScenarioId {
  return ["rebuild", "contention", "decline", "expansion"].includes(value);
}

function isCityManagementModeId(value: string): value is CityManagementModeId {
  return ["owner_influence", "dual_figure"].includes(value);
}
