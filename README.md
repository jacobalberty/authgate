# AuthGate

AuthGate is a Go toolkit for integrating NetBird’s Zero Trust VPN into proxies and gateways by verifying user access and group membership.

## Integrations

Available integrations—more to come. Add a new folder under `cmd/` with your tool and create its README under `docs/`, then link it below.

| Integration       | Directory        | Documentation                        |
|-------------------|------------------|--------------------------------------|
| Squid ACL Helper  | `cmd/squidgate`  | [Squid Integration](docs/squidgate.md) |

## Contributing

1. Fork the repo
2. Add a new `cmd/<integration>` folder with your helper
3. Write a `docs/<integration>.md` for it
4. Update this table to include your integration
5. Open a pull request

## License

BSD 3-Clause License © Jacob Alberty
