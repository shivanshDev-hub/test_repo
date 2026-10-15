# syntax=docker/dockerfile:1
#
# OnionEye — Tor Sidecar (Team 2)
# Provides SOCKS5 egress for the crawler workers, with an iptables kill switch
# applied BEFORE the Tor daemon starts so no traffic can leak outside Tor.
#
# Requires NET_ADMIN capability + /dev/net/tun device (see docker-compose.yml).
# Entrypoint applies firewall rules as root, then Tor drops privileges.
# Base layer is digest-pinned for supply-chain integrity.

FROM alpine:3.19@sha256:6baf43584bcb78f2e5847d1de515f23499913ac9f12bdf834811a3145eb11ca1

# tor, iptables (IPv4 + IPv6), netcat for SOCKS health check
RUN apk --no-cache add tor iptables ip6tables netcat-openbsd

# Tor needs a data dir with permissive ownership (entrypoint runs as root for iptables)
RUN addgroup -S -g 10004 toruser && adduser -S -u 10004 -G toruser toruser \
    && mkdir -p /var/lib/tor /etc/tor \
    && chmod 700 /var/lib/tor \
    && chown root:root /var/lib/tor

# torrc is a TEMPLATE (placeholders) — the entrypoint renders the runtime file
# to /tmp/torrc so HashedControlPassword + MaxCircuitDirtiness rotate with env.
COPY torrc /etc/tor/torrc.tpl
COPY entrypoint.sh /entrypoint.sh
RUN chmod +x /entrypoint.sh

# The entrypoint must stay root to apply iptables, and Tor now runs as root
# because Docker dropped SETUID/SETGID. This container has NET_ADMIN only.
ENTRYPOINT ["/entrypoint.sh"]

LABEL org.opencontainers.image.title="onioneye-tor-sidecar" \
      org.opencontainers.image.description="Tor SOCKS5 sidecar with iptables kill switch (Team 2)" \
      org.opencontainers.image.source="https://github.com/onioneye/onioneye" \
      org.opencontainers.image.version="0.1.0"

EXPOSE 9050 9051
