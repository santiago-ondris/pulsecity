#!/bin/bash
set -e

echo "==> Configurando entorno PulseCity..."

# Rust: agregar targets y componentes útiles
rustup component add clippy rustfmt
rustup target add wasm32-unknown-unknown

# Go: herramientas
go install github.com/air-verse/air@latest

# Agentes de código
echo "==> Instalando Claude Code y Codex..."
npm install -g @anthropic-ai/claude-code
npm install -g @openai/codex

# Levantar base de datos y NATS
echo "==> Levantando servicios (TimescaleDB + NATS)..."

echo "==> Entorno listo."