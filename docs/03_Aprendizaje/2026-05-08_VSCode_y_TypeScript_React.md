# 2026-05-08 — VSCode y falsos errores de TypeScript en el frontend

## El síntoma

A veces VSCode marca errores como:

- `Cannot find module 'react'`
- `JSX element implicitly has type 'any'`

aunque el proyecto compile correctamente por terminal.

---

## Qué significa normalmente

En este contexto, eso no suele indicar que el frontend esté roto.

Si desde `frontend/` pasan comandos como:

```bash
./node_modules/.bin/tsc -p tsconfig.app.json --noEmit
npm run build
```

entonces el problema casi siempre es del **TypeScript server del editor**, no del proyecto.

---

## Qué quedó configurado en este repo

Se agregó configuración local de VSCode en:

- [.vscode/settings.json](/workspace/.vscode/settings.json)

Objetivo:

- forzar el uso de la versión local de TypeScript del frontend
- reducir falsos errores por desincronización del editor
- evitar ruido de watchers sobre `node_modules`, `dist` y `target`

También quedó:

- [.vscode/extensions.json](/workspace/.vscode/extensions.json)

con extensiones recomendadas para Go, Rust y frontend.

---

## Qué hacer si vuelven a aparecer errores en rojo

En este orden:

1. Ejecutar en `frontend/`:

```bash
npm install
```

2. En VSCode:

- `TypeScript: Select TypeScript Version`
- elegir `Use Workspace Version`

3. Luego:

- `TypeScript: Restart TS Server`

4. Si sigue mal:

- `Developer: Reload Window`

---

## Cómo distinguir error real de error del editor

### Si esto pasa:

```bash
cd frontend
./node_modules/.bin/tsc -p tsconfig.app.json --noEmit
```

sin errores, entonces:

- el proyecto está bien tipado
- el error es del editor o de su caché

### Si falla ese comando:

entonces sí hay un problema real del proyecto y hay que corregir código o configuración.

---

## Regla práctica

Para este repo, si VSCode muestra algo raro en el frontend:

- primero confiar en `tsc` y `npm run build`
- después reiniciar el TS server
- recién ahí asumir que el código está roto
