FROM debian:bookworm-slim
RUN apt-get update && apt-get install -y --no-install-recommends spin gcc libc6-dev python3 \
    && rm -rf /var/lib/apt/lists/*
WORKDIR /tp
ENTRYPOINT ["python3", "/tp/verificar_spin.py"]
