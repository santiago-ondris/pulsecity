#!/bin/bash
set -e

echo "==> Configurando entorno PulseCity..."

# Rust: agregar targets y componentes útiles
rustup component add clippy rustfmt
rustup target add wasm32-unknown-unknown

# Go: herramientas
go install github.com/air-verse/air@latest

echo "==> Entorno listo."