# pcfarm 管理系统

## 状态

Active

## 需求来源

用户希望基于当前 gin-vue-admin 项目开发 pcfarm 管理系统，用于管理服务器，使服务器通过 PXE 引导进入全新的 Ubuntu Live 系统，并支持远程启动、关闭、进入 PXE 和 IP 分配。

## 已确认范围

- pcfarm 直接管理 DHCP/TFTP/PXE/IP 分配。
- 装机网络为独立 VLAN，pcfarm 可以接管该网段 DHCP。
- 远控同时支持 IPMI 和 Redfish。
- 服务器每次通过 PXE 启动进入全新的 Ubuntu Live 环境。
- Ubuntu Live 启动后运行 agent 自动注册和上报状态。
- 资产身份使用序列号或资产编号，PXE/DHCP 使用 PXE MAC 绑定。
- MVP 固定一套 Ubuntu Live 镜像。
- IP 从地址池自动分配，并与 PXE MAC 长期绑定。
- pcfarm 部署节点同时运行 DHCP、TFTP、HTTP/NFS 等 PXE 相关服务。
- 支持 UEFI 和 Legacy BIOS PXE。
- 每台服务器支持本地盘、Ubuntu Live、维护模式三种启动策略。

## 设计文档

- `docs/superpowers/specs/2026-05-18-pcfarm-management-design.md`

## 当前决策

采用 `dnsmasq` 一体化单节点 MVP，后端通过 `PXEProvider`、`PowerProvider`、`IPAllocator` 保留后续替换 Kea DHCP 或拆分 PXE Agent 的边界。
