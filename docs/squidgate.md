# squidgate

> Documentation for the `squidgate` integration (see `docs/squidgate.md`).

A Squid external‑ACL helper in Go that checks user access with NetBird’s Zero Trust VPN by verifying peer IP and group membership. It returns `OK` or `ERR` for ACL enforcement.

## Build

```bash
cd cmd/squidgate
go build -o squidgate
```

## Configuration

1. Create a JSON config file named `squidgate.json` in the current directory (defaults to this if no path is passed), or supply a custom file path:

   ```json
   {
     "root_endpoint": "https://api.netbird.io",
     "token": "YOUR_SERVICE_TOKEN"
   }
   ```

2. Ensure `squidgate` is executable and in your PATH.
3. In your `squid.conf`, configure the external helper, passing the config file as its sole argument:

   ```conf
   external_acl_type group_acl ttl=30 %DATA %SRC /usr/local/bin/squidgate /path/to/squidgate.json
   
   # Define VPN groups
   acl allow_internet external group_acl allow-internet
   acl admin          external group_acl admin

   # Define a blocklist for allow-internet users
   acl blocklist dstdomain "/etc/squid/blocklist.txt"
   
   # Access rules:
   # 1) Admins get full access
   http_access allow admin

   # 2) allow_internet users get all except blocklist
   http_access deny blocklist allow allow_internet
   http_access allow allow_internet

   # 3) Deny all others
   http_access deny all
   ```

## License

BSD 3-Clause License © Jacob AlbertyMIT © Jacob Alberty
