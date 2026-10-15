#!/bin/sh
# OnionEye — generate htpasswd for nginx basic auth at container start.
# Credentials come from the environment (NGINX_AUTH_USER / NGINX_AUTH_PASS,
# sourced from .env) so no credential file is ever committed to the repo.
set -eu

if [ "${NGINX_AUTH_USER:-}" ] && [ "${NGINX_AUTH_PASS:-}" ]; then
    # BusyBox mkpasswd — nginx supports $6$ (SHA-512) hashes in auth_basic_user_file.
    hash=$(mkpasswd -m sha512 "$NGINX_AUTH_PASS")
    printf '%s:%s\n' "$NGINX_AUTH_USER" "$hash" > /etc/nginx/htpasswd
    chmod 640 /etc/nginx/htpasswd    # while still owned by root
    chown nginx:nginx /etc/nginx/htpasswd
    echo "[onioneye] htpasswd generated for user: $NGINX_AUTH_USER"
else
    echo "[onioneye] WARNING: NGINX_AUTH_USER/NGINX_AUTH_PASS unset — auth_basic disabled"
fi