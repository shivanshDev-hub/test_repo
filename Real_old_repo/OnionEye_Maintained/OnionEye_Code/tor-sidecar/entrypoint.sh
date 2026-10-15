#!/bin/sh
set -eu

# ── 1. Derive control-port auth hash from the CURRENT .env secret ─────────────
# Rotating TOR_CONTROL_PASSWORD in .env never requires rebuilding the sidecar:
# the Hash is recomputed and injected at every boot.
TOR_CONTROL_PASSWORD=${TOR_CONTROL_PASSWORD:-changeme_in_production}
echo "[entrypoint] Computing HashedControlPassword from env (rotates per restart)..."
hash=$(tor --hash-password "$TOR_CONTROL_PASSWORD" 2>/dev/null | tail -n1)
[ -n "$hash" ] || hash="16:31F07138218DC88460B39D4B4BCA44D85AB314C1DD8F62CF1CCFC13968"

# Randomized circuit lifetime (300–900s) per boot — makes Tor state non-static
# across restarts and reduces long-run fingerprint correlation.
dirtiness=$(awk 'BEGIN{srand(); print int(rand()*601)+300}')
echo "[entrypoint] MaxCircuitDirtiness randomized -> ${dirtiness}s"

# Render the runtime torrc from the template (rootfs is read-only, so the
# generated file lands on the /tmp tmpfs).
sed -e "s|__TOR_CONTROL_HASH__|$hash|" \
    -e "s|__MAX_CIRCUIT_DIRTINESS__|$dirtiness|" \
    /etc/tor/torrc.tpl > /tmp/torrc

# ── 2. Apply iptables kill switch (before Tor starts) ─────────────────────────
echo "[entrypoint] Applying iptables kill switch..."

# Ensure custom chain exists if referenced
iptables -N DOCKER-USER 2>/dev/null || true
ip6tables -N DOCKER-USER 2>/dev/null || true

# IPv4: allow loopback, established connections, Tor SOCKS, DNS, and internal subnets
iptables -A INPUT -m conntrack --ctstate ESTABLISHED,RELATED -j ACCEPT
iptables -A INPUT -i lo -j ACCEPT
iptables -A INPUT -p tcp --dport 9050 -j ACCEPT
iptables -A INPUT -p tcp --dport 9051 -j ACCEPT

iptables -A OUTPUT -m conntrack --ctstate ESTABLISHED,RELATED -j ACCEPT
iptables -A OUTPUT -o lo -j ACCEPT
iptables -A OUTPUT -p udp --dport 53 -j ACCEPT
iptables -A OUTPUT -p tcp --dport 9050 -j ACCEPT
iptables -A OUTPUT -d 172.16.0.0/12 -j ACCEPT

# IPv6: block everything except loopback and established
ip6tables -A INPUT -m conntrack --ctstate ESTABLISHED,RELATED -j ACCEPT
ip6tables -A INPUT -i lo -j ACCEPT
ip6tables -A OUTPUT -m conntrack --ctstate ESTABLISHED,RELATED -j ACCEPT
ip6tables -A OUTPUT -o lo -j ACCEPT

# Persist rules (if tool exists)
netfilter-persistent save 2>/dev/null || true

echo "[entrypoint] Kill switch active."

# ── 3. Start Tor daemon ───────────────────────────────────────────────────────
echo "[entrypoint] Starting Tor daemon..."
# DataDirectory must be owned by the user that runs Tor (root here: SETUID dropped).
mkdir -p /var/lib/tor
chown -R root:root /var/lib/tor
chmod 700 /var/lib/tor
tor -f /tmp/torrc &
TOR_PID=$!

# ── 4. Wait for SOCKS5 endpoint ───────────────────────────────────────────────
echo "[entrypoint] Waiting for SOCKS5 endpoint..."
until nc -z 127.0.0.1 9050 2>/dev/null; do
    sleep 1
done
echo "[entrypoint] SOCKS5 healthy on :9050"

# ── 5. Keep container alive ───────────────────────────────────────────────────
echo "[entrypoint] Tor sidecar ready."
wait "$TOR_PID"
