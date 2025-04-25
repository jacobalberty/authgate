# Squidgate + Docker — Intercept-Proxy Quick-Start

A complete guide to run **squidgate** with an intercepting Squid proxy for NetBird peers, providing three policy levels:

| Tag                | Client capabilities                                                     |
|--------------------|------------------------------------------------------------------------|
| **internet-block** | No web access (captive portal or quarantine)                           |
| **internet-filter**| Only allow-listed sites and block downloads with dangerous extensions  |
| **internet-open**  | Full internet access without restrictions                               |

`squidgate` queries NetBird to determine each peer’s tag, and Squid enforces policies based on that tag.

---

## 1  Repository Layout
```
infra/
├── compose.yaml
├── netns/
│   ├── Dockerfile             # builds netns container with iptables
│   └── intercept-rules.sh     # redirects HTTP/HTTPS to Squid
├── squid/
│   ├── conf/squid.conf        # Squid configuration
│   ├── policies/
│   │   ├── allowlist.txt      # domains allowed for "filter" tag
│   │   └── dangerous_ext.txt   # regex of file extensions to block
│   ├── ssl_db/                # persisted SSL cert DB for ssl-bump
│   ├── spool/                 # Squid cache_dir and swap space
│   └── log/                   # access.log, cache.log
└── squidgate/
    ├── squidgate             # compiled Go binary
    └── squidgate.json        # NetBird API token config
```
Adapt paths if you clone from `examples/intercept`; the structure must match your volume mounts.

---

## 2  netns Container: Intercept Rules
**intercept-rules.sh** (placed in `netns/`):
```bash
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
```
**Dockerfile** (`netns/Dockerfile`):
```dockerfile
FROM debian:12-slim
RUN apt-get update && apt-get install -y iptables && rm -rf /var/lib/apt/lists/*
COPY intercept-rules.sh /intercept-rules.sh
ENTRYPOINT ["/bin/sh","/intercept-rules.sh"]
```

---

## 3  Squid Configuration
Edit **`squid/conf/squid.conf`**:
```conf
### Intercept Listeners
http_port 3129 intercept
https_port 3130 intercept ssl-bump \
    cert=/etc/squid/ssl_db/ca.pem \
    generate-host-certificates=on \
    dynamic_cert_mem_cache_size=8MB \
    tls-dh=/etc/squid/conf/dhparam.pem

### SSL-Bump Policy (peek + splice)
acl step1 at_step SslBump1
ssl_bump peek step1
ssl_bump splice all    # change to "bump all" for full MITM

### External ACL with squidgate
external_acl_type group_acl ttl=30 %DATA %SRC /usr/local/bin/squidgate /etc/squidgate.json
acl block_internet   external group_acl internet-block
acl filter_internet  external group_acl internet-filter
acl allow_internet   external group_acl internet-open

### Domain and Download ACLs
acl allow_sites        dstdomain "/etc/squid/policies/allowlist.txt"
acl dangerous_download urlpath_regex -i "/etc/squid/policies/dangerous_ext.txt"

### Policy Enforcement
http_access deny  block_internet
http_access deny  filter_internet !allow_sites
http_access deny  filter_internet dangerous_download
http_access allow allow_internet
http_access deny all
```

---

## 4  compose.yaml
```yaml
version: "3.9"
services:
  # 1) namespace anchor with intercept rules
  netns:
    build: ./netns
    cap_add: [ NET_ADMIN ]

  # 2) Squid + squidgate binary
  squid:
    image: ghcr.io/jacobalberty/squid-docker:v6.12
    network_mode: "service:netns"
    restart: always
    cap_add: [ NET_BIND_SERVICE ]
    volumes:
      - ./squid/conf:/conf
      - ./squid/spool:/var/spool/squid
      - ./squid/ssl_db:/var/spool/squid/ssl_db
      - ./squid/log:/var/log/squid
      - ./squidgate/squidgate:/usr/local/bin/squidgate:ro
      - ./squidgate/squidgate.json:/etc/squidgate.json:ro
    depends_on: [ netns ]

  # 3) NetBird agent as exit node
  netbird:
    image: netbirdio/netbird:latest
    network_mode: "service:netns"
    cap_add: [ NET_ADMIN ]
    environment:
      - NB_TOKEN=YOUR_SERVICE_TOKEN
    depends_on: [ netns ]
```

---

## 5  Build the squidgate Binary
On your host with Go installed:
```bash
mkdir -p output
go build -o output/squidgate ./cmd/squidgate
mv output/squidgate squidgate/squidgate
```

---

## 6  Generate CA & DH Parameters
```bash
# SSL cert DB (MITM CA) persisted to ./squid/ssl_db
docker compose run --rm squid \
  /usr/lib/squid/security_file_certgen -s /var/spool/squid/ssl_db -M 8MB -c

# DH parameters for TLS
openssl dhparam -out ./squid/conf/dhparam.pem 2048
```

---

## 7  Configure NetBird as Exit Node
1. In the NetBird UI, go to **Peers**.
2. Click the **⋮ menu** next to the desired peer and select **Set Up Exit Node**.
3. Choose or create a **Distribution Group**, then add peers to it.

Those peers’ Internet traffic will now flow through the intercepting Squid proxy.

---

## 8  Launch and Test
```bash
docker compose up -d
```
- Tag a peer as **internet-filter** → only allow-listed domains load; dangerous downloads blocked.
- Tag as **internet-open** → full access without restrictions.

Reload Squid after policy changes:
```bash
docker compose exec squid squid -k reconfigure
```

---

## 9  Extensions
- Add file extensions in `squid/policies/dangerous_ext.txt` (`apk`, `dmg`, …).
- Switch to `ssl_bump bump all` for full HTTPS MITM after distributing `ca.pem`.
- Use `delay_pools` keyed on `group_acl` tags for bandwidth shaping.

Enjoy dynamic per-peer internet policies with Squidgate and Docker!

