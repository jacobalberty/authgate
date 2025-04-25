#!/bin/sh
set -eu
CAP_IF=${CAP_IFACE:-wg0}
# Remove old rules
iptables -t nat -D PREROUTING -i "$CAP_IF" -p tcp --dport 80  -j REDIRECT --to-port 3129 2>/dev/null || true
iptables -t nat -D PREROUTING -i "$CAP_IF" -p tcp --dport 443 -j REDIRECT --to-port 3130 2>/dev/null || true
# Insert new rules at the top
iptables -t nat -I PREROUTING 1 -i "$CAP_IF" -p tcp --dport 80  -j REDIRECT --to-port 3129
iptables -t nat -I PREROUTING 1 -i "$CAP_IF" -p tcp --dport 443 -j REDIRECT --to-port 3130
exec sleep infinity