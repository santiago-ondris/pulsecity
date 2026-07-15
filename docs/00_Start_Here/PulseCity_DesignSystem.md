# PulseCity — Design System & Identity Kit

> Documento de referencia para desarrollo frontend. Toda decisión visual de la app debe respetar este sistema. Aplica a Claude Code y cualquier agente que toque el front.

---

## 1. Filosofía visual

PulseCity es una app de gestión densa y de uso prolongado. El diseño prioriza:

- **Claridad sobre decoración** — cada elemento visual tiene un propósito funcional
- **Jerarquía legible** — el usuario sabe en todo momento qué es importante y qué es soporte
- **Profundidad real** — el dark mode tiene capas distinguibles, no todo aplastado en el mismo plano
- **Identidad propia** — no es un dashboard SaaS genérico, es una sala de control con personalidad editorial

Referentes visuales: Cities Skylines 2 (densidad y panels flotantes sobre mapa), NBA 2K franchise mode (profundidad de gestión), The Athletic (peso narrativo tipográfico).

### 1.1 Dirección visual de pantallas de flujo

Las pantallas de entrada, creación de partida, revisión y ceremonia no se diseñan como dashboards ni como formularios administrativos. Se diseñan como **pantallas de juego funcionales**:

- Fondo visual atmosférico de ciudad/arena/mapa, oscurecido para no competir con la UI.
- Topbar compacto con contexto mínimo: paso, estado o identidad activa.
- Bloque principal centrado o levemente elevado del centro geométrico.
- Una decisión dominante por pantalla.
- CTA principal cerca del contexto de decisión, nunca perdido al final de una página larga.
- Preview solo cuando ayuda a decidir; no se usa como decoración.
- Texto breve. La explicación larga se elimina o se mueve a documentación, no a la pantalla.

El usuario debe entender en menos de 3 segundos:

- dónde está
- qué está decidiendo
- qué botón avanza

### 1.2 Anti-dashboard

PulseCity puede tener pantallas densas de gestión, pero el flujo inicial no debe parecer un SaaS. Evitar:

- hero titles enormes que empujan la acción fuera del viewport
- grillas genéricas de cards sin jerarquía
- formularios reacomodados visualmente sin repensar la tarea
- sidebars decorativos
- paneles flotantes que dividen la atención sin necesidad
- fondos negros planos cuando la pantalla representa un momento de juego
- texto visible explicando cómo usar la interfaz

Si una pantalla pide una decisión, el layout debe organizarse alrededor de esa decisión.

---

## 2. Paleta de color

### 2.1 Fondos — siempre fijos, nunca se mezclan con colores de franquicia

| Token | Hex | Uso |
|-------|-----|-----|
| `bg-base` | `#0A0A0B` | Fondo de toda la app, detrás de todo |
| `bg-surface` | `#111113` | Panels principales, sidebars, modales |
| `bg-elevated` | `#1A1A1E` | Cards dentro de panels, dropdowns, tooltips |
| `bg-border` | `#252528` | Separadores de sección, dividers |

### 2.2 Texto

| Token | Hex | Uso |
|-------|-----|-----|
| `text-primary` | `#F0F0F0` | Títulos, valores importantes, texto principal |
| `text-muted` | `#666666` | Labels, metadata, texto de soporte |
| `text-disabled` | `#333333` | Elementos inactivos |

### 2.3 Acento primario — identidad PulseCity

| Token | Hex | Uso |
|-------|-----|-----|
| `accent-primary` | `#00C896` | Botones primarios, estados activos, highlights, identidad |
| `accent-primary-bg` | `#00C89616` | Fondo de badges y áreas con acento verde |
| `accent-primary-border` | `#00C89640` | Borde de badges y componentes con acento verde |

El verde `#00C896` cumple dos roles: **identidad visual** de PulseCity y **estado de éxito** (contrato firmado, trade completado, objetivo cumplido). Son el mismo color deliberadamente.

### 2.4 Estados semánticos

| Estado | Token | Hex | Uso típico |
|--------|-------|-----|------------|
| Urgencia / peligro | `semantic-urgent` | `#FF6B2B` | Trade deadline, lesión grave, cap en rojo |
| Warning suave | `semantic-warning` | `#FFAA00` | Moral bajando, cap ajustado, agente insatisfecho |
| Info neutral | `semantic-info` | `#7B8CDE` | Eventos de ciudad, notas de prensa, recordatorios |
| Negativo / pérdida | `semantic-negative` | `#E05555` | Partido perdido, jugador que se va, aprobación cayendo |

Cada estado semántico tiene variantes de fondo y borde siguiendo el patrón del acento:
- Fondo: `{hex}16` (10% opacidad)
- Borde: `{hex}40` (25% opacidad)

### 2.5 Fondos visuales y overlays

Para pantallas de flujo (`LandingPage`, `IdentityPage`, `ScenarioPage`, `ManagementPage`, `LaunchPage`, `CeremonyPage`) se permite usar un fondo visual full-screen. El fondo cumple rol de atmósfera, no de contenido principal.

Reglas:

- El fondo debe ocupar todo el viewport.
- Siempre debe tener overlay oscuro suficiente para legibilidad.
- El contenido interactivo se mantiene sobre panels translúcidos con borde sutil.
- No dejar bandas superiores o laterales del `body`/wrapper visibles.
- No usar el fondo para comunicar datos críticos.

Los colores de franquicia pueden aparecer como acentos locales en pantallas de identidad/revisión, pero no reemplazan la paleta base de PulseCity.

### 2.6 Elevación y bordes

Las capas se distinguen por **color + borde sutil**. Nunca solo color, nunca borde sin diferencia de tono.

```css
/* Borde universal para cards y componentes */
border: 1px solid rgba(255, 255, 255, 0.031); /* #ffffff08 */
```

Regla: cada capa tiene su propio tono de fondo + este borde. No hay sombras (`box-shadow`) salvo en modales sobre el mapa de ciudad.

---

## 3. Tipografía

### 3.1 Familias

| Rol | Familia | Fuente | Uso |
|-----|---------|--------|-----|
| Display / títulos | **Bona Nova SC** | Google Fonts | Nombres de jugadores, stats protagonistas, titulares de eventos, citas de agentes |
| UI / datos / cuerpo | **DM Sans** | Google Fonts | Tablas, labels, badges, navegación, texto de soporte, diálogos |

```html
<!-- Import en el proyecto -->
<link href="https://fonts.googleapis.com/css2?family=Bona+Nova+SC:ital,wght@0,400;0,700;1,400&family=DM+Sans:wght@400;500&display=swap" rel="stylesheet">
```

### 3.2 Escala tipográfica

| Token | Familia | Tamaño | Peso | Uso |
|-------|---------|--------|------|-----|
| `display-xl` | Bona Nova SC | 48px | 700 | Stats protagonistas (PPG, puntos) |
| `display-lg` | Bona Nova SC | 36px | 700 | Nombre de jugador, titular principal |
| `display-md` | Bona Nova SC | 24px | 700 | Subtítulos de sección, eventos de ciudad |
| `display-sm` | Bona Nova SC | 18px | 700 | Titulares secundarios |
| `display-italic` | Bona Nova SC | 18px | 400 italic | Citas de agentes, momentos narrativos |
| `ui-lg` | DM Sans | 14px | 500 | Labels de sección, navegación activa |
| `ui-md` | DM Sans | 13px | 400 | Texto de tabla, contenido de cards |
| `ui-sm` | DM Sans | 11px | 400 | Metadata, fechas, eyebrows |
| `ui-xs` | DM Sans | 10px | 500 | Badges, tags, indicadores |

### 3.3 Reglas tipográficas

- Bona Nova SC **solo** para títulos y momentos de peso narrativo. Nunca en tablas ni navegación.
- DM Sans para todo lo funcional. Es la tipografía que el usuario "no ve" — simplemente lee.
- Números en stats grandes van en Bona Nova SC. Números en tablas van en DM Sans.
- Eyebrows (texto sobre título) siempre en DM Sans `ui-sm`, uppercase, `letter-spacing: 0.1em`, color `text-muted`.

---

## 4. Forma y espaciado

### 4.1 Border radius

| Token | Valor | Uso |
|-------|-------|-----|
| `radius-panel` | `8px` | Cards, panels, modales, contenedores principales |
| `radius-component` | `4px` | Badges, inputs, botones, mini-componentes internos |
| `radius-pill` | `999px` | Solo para pills de estado muy específicos |

### 4.2 Densidad

PulseCity es una app de **densidad media-alta**. La información se agrupa en bloques compactos con separadores claros. No hay mucho aire entre elementos pero la jerarquía visual evita la sensación de asfixia.

Espaciado base: múltiplos de `4px`.

| Token | Valor | Uso |
|-------|-------|-----|
| `space-xs` | `4px` | Gap entre elementos muy relacionados |
| `space-sm` | `8px` | Padding interno de badges y mini-componentes |
| `space-md` | `12px` | Gap entre elementos en una card |
| `space-lg` | `16px` | Padding interno de cards |
| `space-xl` | `24px` | Gap entre cards o secciones |
| `space-2xl` | `32px` | Separación entre secciones mayores |

---

## 5. Iconografía

**Set oficial: Phosphor Icons**

Se usa en sus distintos pesos según el contexto:

| Peso | Uso |
|------|-----|
| `thin` | Tablas densas, metadata, contextos de mucha información |
| `regular` | Uso general en UI, navegación, labels |
| `bold` | Botones de acción, estados importantes, CTAs |
| `fill` | Estado activo/seleccionado de un ícono navegable |

Regla: nunca mezclar pesos distintos en el mismo componente. Un panel elige su peso y lo mantiene consistente.

```bash
# Instalación
npm install @phosphor-icons/react
```

```jsx
// Uso en React
import { Basketball, Trophy, Buildings } from '@phosphor-icons/react'

// Thin en tabla
<Basketball size={16} weight="thin" />

// Bold en botón de acción
<Trophy size={20} weight="bold" />
```

---

## 6. Componentes base

### 6.1 Badge de estado

```jsx
// Estructura universal de badge semántico
<span style={{
  display: 'inline-flex',
  alignItems: 'center',
  gap: '5px',
  fontSize: '10px',
  fontWeight: 500,
  fontFamily: 'DM Sans',
  padding: '3px 8px',
  borderRadius: '4px',
  background: `${color}16`,
  color: color,
  border: `1px solid ${color}40`,
}}>
  <span style={{ width: 6, height: 6, borderRadius: '50%', background: color }} />
  {label}
</span>
```

### 6.2 Card estándar

```css
.card {
  background: #1A1A1E;          /* bg-elevated */
  border: 1px solid #ffffff08;  /* borde sutil */
  border-radius: 8px;           /* radius-panel */
  padding: 16px;                /* space-lg */
}
```

### 6.3 Eyebrow + título

```jsx
// Patrón estándar de encabezado de sección
<div>
  <p style={{
    fontFamily: 'DM Sans',
    fontSize: '10px',
    fontWeight: 500,
    color: '#666',
    letterSpacing: '0.1em',
    textTransform: 'uppercase',
    marginBottom: '6px',
  }}>
    {eyebrow}
  </p>
  <h2 style={{
    fontFamily: 'Bona Nova SC',
    fontSize: '36px',
    fontWeight: 700,
    color: '#F0F0F0',
    lineHeight: 1.05,
  }}>
    {title}
  </h2>
</div>
```

---

## 7. Patrones de pantalla

### 7.1 Menú principal / landing

El menú principal es una pantalla de juego, no una landing de marketing.

Reglas:

- Fondo atmosférico full-screen.
- Acciones principales apiladas verticalmente en el centro.
- No mostrar biblioteca, login o formularios como contenido permanente de la landing.
- `Nueva partida`, `Cargar partida`, `Jugar como invitado` deben leerse como decisiones hermanas.
- Cuenta y biblioteca se abren como pantalla/panel posterior, no compiten con el primer impacto.

### 7.2 Onboarding de nueva partida

Cada paso resuelve una sola pregunta:

| Pantalla | Pregunta |
|---|---|
| Identidad | ¿Cómo se llama y se ve la franquicia? |
| Escenario | ¿Desde dónde empieza la historia? |
| Gestión | ¿Qué relación de poder tenés con la ciudad? |
| Revisión | ¿Confirmás esta fundación? |
| Ceremonia | ¿Cómo nace el mundo en tiempo real? |

Reglas:

- Mantener topbar compacto.
- No usar `step-card + grid` como patrón por defecto.
- Mostrar una vista/resumen del estado seleccionado.
- Las opciones deben ser botones grandes, legibles y comparables.
- El estado activo debe ser evidente por borde, fondo y acento.
- El contenido debe caber en desktop sin quedar tirado abajo.
- En mobile, el contenido crítico no puede quedar cortado.

### 7.3 Pantallas operativas

Cuando la pantalla muestra estado vivo del backend, como la ceremonia del mapa:

- El dato vivo principal ocupa el área dominante.
- Métricas, pipeline y eventos quedan como soporte.
- No ocultar estado técnico relevante si ayuda a entender el proceso.
- Mantener `snapshot`/`patch` y deltas como modelo mental: el frontend observa cambios, no inventa estado.

---

## 8. Animaciones

Filosofía: **el movimiento tiene propósito, nunca es decorativo**.

| Tipo | Duración | Easing | Uso |
|------|----------|--------|-----|
| Transición de estado | `150ms` | `ease-out` | Cambio de valor, aparición de badge, apertura de panel |
| Feedback de acción | `100ms` | `ease-out` | Scale `0.97` en botones al hacer click |
| Highlight de fila | `150ms` | `ease-out` | Flash momentáneo en fila afectada por una acción |
| Evento narrativo | `250ms` | `ease-out` | Entrada de panel de evento importante (slide desde abajo) |

Reglas:
- Navegación entre secciones: instantánea o `100ms` máximo
- Hover states: `100ms`
- Sin rebotes (`spring`, `bounce`), sin easing dramático
- Sin animaciones en scroll
- Los eventos narrativos (trade completado, lesión, evento de ciudad) son el único momento donde la animación puede tener carácter — entran desde abajo con el color semántico correspondiente

---

## 9. Checklist de aceptación visual

Antes de dar por buena una pantalla nueva:

- ¿La acción principal se entiende en 3 segundos?
- ¿Hay una sola decisión dominante?
- ¿El CTA está cerca del contexto de decisión?
- ¿El fondo acompaña sin competir?
- ¿La pantalla parece parte de un juego, no de un dashboard genérico?
- ¿No hay bandas vacías del wrapper/body visibles?
- ¿Los textos largos caben sin desbordar cards o botones?
- ¿Desktop no empuja el bloque principal demasiado abajo?
- ¿Mobile no corta contenido crítico?
- ¿La UI conserva la paleta PulseCity aunque use acentos de franquicia?

---

## 10. CSS Variables — referencia rápida

```css
:root {
  /* Fondos */
  --bg-base: #0A0A0B;
  --bg-surface: #111113;
  --bg-elevated: #1A1A1E;
  --bg-border: #252528;

  /* Texto */
  --text-primary: #F0F0F0;
  --text-muted: #666666;
  --text-disabled: #333333;

  /* Acento primario */
  --accent-primary: #00C896;
  --accent-primary-bg: #00C89616;
  --accent-primary-border: #00C89640;

  /* Semánticos */
  --semantic-urgent: #FF6B2B;
  --semantic-warning: #FFAA00;
  --semantic-info: #7B8CDE;
  --semantic-negative: #E05555;

  /* Borde universal */
  --border-subtle: rgba(255, 255, 255, 0.031);

  /* Radio */
  --radius-panel: 8px;
  --radius-component: 4px;

  /* Tipografía */
  --font-display: 'Bona Nova SC', serif;
  --font-ui: 'DM Sans', sans-serif;
}
```

---

## 11. Lo que NO es PulseCity

- No usa gradientes decorativos genéricos como reemplazo de dirección visual.
- No usa sombras (`box-shadow`) salvo en modales flotantes sobre el mapa
- No usa más de un color de acento por componente
- No usa Bona Nova SC en tablas, navegación ni texto de soporte
- No usa animaciones sin propósito funcional
- No convierte colores de franquicia en theme global de la app
- No usa formularios genéricos como experiencia final de juego
- No usa cards anidadas como estructura principal de página

---

*PulseCity Design System v1.1 — actualizado con dirección visual del onboarding el 2026-05-20*
