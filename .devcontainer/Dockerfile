FROM mcr.microsoft.com/devcontainers/base:ubuntu-24.04

RUN apt-get update && apt-get install -y \
    curl \
    git \
    build-essential \
    pkg-config \
    libssl-dev \
    protobuf-compiler \
    && rm -rf /var/lib/apt/lists/*

# NATS CLI
RUN NATS_VERSION=$(curl -s https://api.github.com/repos/nats-io/natscli/releases/latest | grep tag_name | cut -d'"' -f4) && \
    curl -fsSL "https://github.com/nats-io/natscli/releases/download/${NATS_VERSION}/nats-${NATS_VERSION#v}-linux-amd64.zip" -o nats.zip && \
    unzip nats.zip && \
    mv nats-*/nats /usr/local/bin/ && \
    rm -rf nats.zip nats-*