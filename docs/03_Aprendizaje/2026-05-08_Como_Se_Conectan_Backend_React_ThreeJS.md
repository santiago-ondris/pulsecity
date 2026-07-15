# 2026-05-08 — Cómo se conectan backend, React y Three.js en la generación del mapa

## La idea clave

React o Three.js **no generan el mapa**. El mapa se genera en el backend.

La separación correcta en PulseCity es:

- `map-service` decide cómo es el mapa
- `gateway` transporta la información al cliente
- React organiza estado y pantallas
- Three.js dibuja visualmente ese estado

En otras palabras:

- **backend** = lógica del mundo
- **React** = orquestación de UI y estado
- **Three.js** = render del mapa

---

## Flujo completo

Cuando el jugador crea una partida, el flujo esperado es:

1. El frontend muestra la pantalla de creación de partida
2. El jugador confirma la creación
3. React hace `POST /api/v1/games` al `gateway`
4. `gateway` publica `mapa.generacion_iniciada` en NATS
5. `map-service` escucha ese evento y genera el mapa
6. `map-service` publica eventos de progreso
7. `gateway` escucha esos eventos
8. `gateway` los traduce a `map.snapshot` y `map.patch`
9. El frontend recibe esos mensajes por WebSocket
10. React actualiza su estado local
11. Three.js vuelve a dibujar usando ese nuevo estado

---

## Qué hace React

React no necesita saber cómo funciona Perlin, Voronoi o cualquier algoritmo de generación procedural.

React se ocupa de:

- iniciar acciones del usuario
- abrir y mantener la conexión WebSocket
- guardar el estado actual del mapa
- aplicar `snapshot` y `patch`
- decidir qué componentes se muestran
- pasar los datos del mapa a la escena visual

Mentalmente, React cumple el rol de coordinador del frontend.

---

## Qué hace Three.js

Three.js toma datos ya generados y los convierte en imagen.

Por ejemplo, si recibe:

- ancho y alto de la grilla
- tipo de terreno por celda
- zona por celda
- ubicación del estadio

entonces puede:

- crear tiles o meshes
- asignar colores y materiales
- ubicar cámara y luces
- renderizar la ciudad o el terreno

Three.js no decide dónde va el estadio ni qué celdas son agua. Solo dibuja lo que el backend ya resolvió.

---

## Qué hace el backend

El backend sí es dueño del mundo.

En este caso:

- `map-service` genera el terreno
- calcula zonas o distritos
- decide dónde ubicar el estadio
- publica el resultado por etapas

Eso significa que el frontend no construye la verdad del mapa. Solo consume una verdad ya construida.

Este patrón es importante porque:

- evita lógica duplicada entre backend y frontend
- mantiene un solo lugar dueño del estado
- hace más simple sincronizar múltiples clientes

---

## Cómo pensar esto si venís de frontend tradicional

Una analogía útil es esta:

- En una app web tradicional, el backend devuelve JSON y React lo renderiza.
- En PulseCity, el backend devuelve un mundo incremental y React/Three.js lo renderizan.

La diferencia no está en el principio, sino en la complejidad del dato:

- antes: una lista, una tabla, un form
- acá: un mapa vivo que llega por snapshot y patches

Pero la idea base es la misma:

- backend produce datos
- frontend los consume
- UI los muestra

---

## Snapshot y patch en este contexto

Esto es importante para entender la integración:

- `map.snapshot` = estado completo base del mapa
- `map.patch` = cambio incremental sobre ese estado

El frontend ideal hace esto:

1. recibe snapshot
2. lo guarda como estado actual
3. recibe patches posteriores
4. actualiza el estado local
5. rerenderiza

Three.js no necesita saber si el dato vino por snapshot o patch. Solo recibe el estado actual ya resuelto por React.

---

## Ejemplo mental simple

Supongamos que llega este snapshot:

```json
{
  "type": "map.snapshot",
  "state": {
    "game_id": "abc",
    "map_data": {
      "width": 28,
      "height": 28,
      "cells": []
    }
  }
}
```

React haría algo como:

- guardar `state.map_data`
- pasar `map_data` a un componente de escena

Después llega un patch con el estadio:

```json
{
  "type": "map.patch",
  "patch": {
    "stadium": { "x": 14, "y": 13 }
  }
}
```

React actualiza el estado local, y Three.js redibuja la escena con el estadio visible.

---

## Frase corta para recordar

`map-service` crea el mundo, `gateway` lo transmite, React lo administra y Three.js lo dibuja.

---

## Por qué importaba tanto ordenar el contrato primero

Antes de construir el frontend real, necesitábamos dejar ordenado:

- qué mensajes salen del backend
- cómo se rehidrata una partida
- cómo se separa snapshot de patch
- dónde vive la persistencia del estado

Si eso está desordenado, el frontend se vuelve frágil.

Por eso el trabajo previo de:

- WebSocket
- snapshot/patch
- rehidratación
- persistencia

no fue accesorio: fue la base para que React y Three.js después entren con menos fricción.
