FROM golang:1.24-trixie AS build

WORKDIR /app
COPY . .

RUN chmod +x docker/install_deps.sh \
    && docker/install_deps.sh

RUN go build -o /usr/local/bin/psirng ./cmd/psirng

FROM debian:trixie-slim

COPY --from=build /app/cert.pem /app/cert.pem
COPY --from=build /app/key.pem /app/key.pem
COPY --from=build /app/3rdparty /app/3rdparty
COPY --from=build /app/docker/install_deps.sh /app/docker/install_deps.sh
COPY --from=build /usr/local/bin/psirng /usr/local/bin/psirng

# Use testing certificate if not provided
RUN mkdir -p /etc/psirng \
    && [ ! -f /etc/psirng/cert.pem ] && [ ! -f /etc/psirng/key.pem ] && cp /app/cert.pem /etc/psirng/cert.pem && cp /app/key.pem /etc/psirng/key.pem

RUN /app/docker/install_deps.sh \
    && apt purge -y --auto-remove g++ cmake \
    && apt -y install curl \
    && apt clean \
    && rm -rf /var/lib/apt/lists/* /tmp/* /var/tmp/*

EXPOSE 8443
EXPOSE 8080
EXPOSE 50051

HEALTHCHECK CMD ["curl", "-f", "http://localhost:8080/healthz"]

CMD ["psirng", "-cert=/etc/psirng/cert.pem", "-key=/etc/psirng/key.pem"]