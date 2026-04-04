# MoFrp

[README](README.md) | [中文文档](README_zh.md)

MoFrp 是基于 [frp](https://github.com/fatedier/frp) 的内网穿透服务客户端。

## 说明

本仓库为 MoFrp 的客户端构建仓库，上游项目为 [fatedier/frp](https://github.com/fatedier/frp)。

相比原版 frpc，这里做了一些调整：

- 禁用了 Web UI，减小二进制体积
- 调整了日志输出格式
- 支持 quickstart 模式，通过 token 快速启动
- AutoTLS 功能

部分实现参考了 [LoliaFrp](https://github.com/Lolia-FRP/lolia-frp)，包括 AutoTLS 和日志输出格式。

## 下载

前往 [Releases](https://github.com/mofrp/cli/releases) 页面下载对应平台的二进制文件。

## 使用

### 快速启动 模式

```bash
./frpc -t <token>
```

token 从 MoFrp 控制面板获取，格式为 `隧道ID:密钥`。

### 传统模式

```bash
./frpc -c frpc.toml
```

配置文件参考 [frp 官方文档](https://gofrp.org/zh-cn/docs/)。

## 致谢

- [fatedier/frp](https://github.com/fatedier/frp) - 原始项目
- [LoliaFrp](https://github.com/Lolia-FRP/lolia-frp) - AutoTLS 和日志格式参考
