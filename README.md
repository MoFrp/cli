# MoFrp

[README](README.md) | [中文文档](README_zh.md)

MoFrp is a frp-based intranet penetration client.

## About

This is the client build repository for MoFrp. Upstream project: [fatedier/frp](https://github.com/fatedier/frp).

Changes compared to the original frpc:

- Disabled Web UI to reduce binary size
- Adjusted log output format
- Added quickstart mode for fast startup via token
- AutoTLS support

Some implementations are referenced from [LoliaFrp](https://github.com/Lolia-FRP/lolia-frp), including AutoTLS and log output format.

## Download

Download the binary for your platform from the [Releases](https://github.com/mofrp/cli/releases) page.

## Usage

### quickstart mode

```bash
./frpc -t <token>
```

Get the token from MoFrp dashboard, format: `TunnelID:SecretKey`.

### Traditional mode

```bash
./frpc -c frpc.toml
```

See [frp documentation](https://gofrp.org/docs/) for configuration.

## Credits

- [fatedier/frp](https://github.com/fatedier/frp) - Original project
- [LoliaFrp](https://github.com/Lolia-FRP/lolia-frp) - AutoTLS and log format reference
