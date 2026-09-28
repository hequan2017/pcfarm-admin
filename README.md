[简体中文](README.md) | [English](README.en.md)

# pcfarm-admin

面向独立装机网段/VLAN 的裸金属装机管理系统：服务器资产台账、IP 地址池分配、PXE 启动策略、IPMI/Redfish 远程电源控制，以及 Ubuntu Live Agent 注册与心跳。

![Go](https://img.shields.io/badge/Go-1.24-00ADD8?logo=go&logoColor=white)
![Vue](https://img.shields.io/badge/Vue-3.5-4FC08D?logo=vue.js&logoColor=white)
![License](https://img.shields.io/badge/License-Apache--2.0-green)

## 📖 项目介绍

机房批量装机通常分散在好几套工具里：资产表格管台账、Excel 管 IP、手工敲 IPMI 命令开关机。`pcfarm-admin` 把这条链路收进一个基于 Gin + Vue 管理框架的后台：以服务器资产为中心，围绕它的 PXE MAC 自动分配装机网段的长期固定 IP，按资产下发三种启动策略（本地盘 / Ubuntu Live / 维护模式），通过 BMC 地址与加密存储的密码执行远程电源操作；Live 系统启动后，Agent 自动注册到管理端并持续上报心跳，形成"登记 → 分配 → 装机 → 纳管"的闭环。

它适合维护独立装机网段/VLAN 的运维团队，也可作为 GPU 算力集群的裸金属上架前置环节（参见同作者的 [docker-gpu-manage](https://github.com/hequan2017/docker-gpu-manage)）。

## ✨ 功能特性

- 🖥️ 服务器资产管理：资产编号、序列号、PXE MAC、BMC 地址/账号、密码密文存储、远控协议（IPMI/Redfish）、启动策略、状态、Agent 版本、硬件摘要、最后心跳时间
- 🌐 IP 地址池：地址池（名称、CIDR、起止 IP、网关、DNS、绑定网卡、启用开关）管理，按 PXE MAC 自动分配长期固定 IP，分配/释放记录可追溯
- 🥾 PXE 启动策略：本地盘（local_disk）、Ubuntu Live、维护模式（maintenance）三种；本地 dnsmasq provider 渲染主机绑定与启动菜单，提供 Refresh / Status 接口
- 🔌 远程电源控制：IPMI 与 Redfish 双 provider，支持开机、关机、重启、下次 PXE 启动
- 🤖 Ubuntu Live Agent：Live 系统启动后按序列号或 PXE MAC 注册资产并上报心跳，写入装机事件（ProvisionEvent）
- 🧭 前端页面：服务器列表/详情、IP 池管理、PXE 设置
- 📝 装机事件流水：注册、分配等关键动作均有事件记录，便于排障追溯

## 🛠 技术栈

| 层 | 技术 |
|---|---|
| 后端 | Go 1.24、Gin 1.10、GORM 1.31、Casbin v3、JWT、Swagger（swaggo） |
| 前端 | Vue 3.5、Vite 8、Element Plus 2.13、Pinia、Axios |
| PXE 方案 | MVP 采用单节点 `dnsmasq` 控制面，后续可替换为 Kea DHCP 或独立 PXE Agent |

## 🚀 快速开始

### 环境要求

- Go 1.24+、Node.js 与 npm（Windows PowerShell、Linux、macOS shell 均可运行）

### 启动后端

```bash
cd server
go run .
```

默认读取 `server/config.yaml`：后端 API `http://127.0.0.1:8888`，Swagger `http://127.0.0.1:8888/swagger/index.html`。如需初始化数据库，请先确认其中 `system.db-type`、数据库连接与 `disable-auto-migrate` 配置。

### 启动前端

```bash
cd web
npm install
npm run serve
```

前端开发地址 `http://127.0.0.1:8080`，`/api` 代理到 `http://127.0.0.1:8888`。

### 构建与测试

```bash
# 前端构建，产物输出到 web/dist/
cd web && npm run build

# 后端 pcfarm 聚焦测试
cd server
go test ./model/pcfarm ./service/pcfarm ./api/v1/pcfarm ./router/pcfarm ./initialize -count=1
```

## 📁 目录结构

```text
├─ server/
│  ├─ model|service|api|router/pcfarm/   # 资产、IP 池、PXE、Agent、电源控制
│  └─ initialize/
├─ web/src/view/pcfarm/                  # server / ipPool / pxe 页面
├─ deploy/                               # docker / docker-compose / kubernetes
├─ docs/  ├─ aiDoc/                      # 项目文档与 AI 协作文档层
└─ AGENT.MD
```

## ⚠️ 当前实现边界

- pcfarm 管理模块的基础模型、服务、API、路由和前端页面已实现。
- PXE provider 与 IPMI/Redfish provider 目前为 MVP 边界实现：真实 `dnsmasq` 配置写入、服务重载、硬件远控命令执行需结合部署环境继续补齐。
- 菜单权限沿用管理框架机制，需在后台菜单管理中配置 pcfarm 页面入口。

## 🔗 相关项目

同一 GPU/算力管理方向的姊妹项目：

- [docker-gpu-manage](https://github.com/hequan2017/docker-gpu-manage) — 天启算力管理平台（GPU 容器实例调度、HAMi 显存切分、K8s 管理）
- [kapigpu](https://github.com/hequan2017/kapigpu) — 卡皮巴拉 GPU 管理平台（GVA 轻量版）
- [tianqi](https://github.com/hequan2017/tianqi) — TianQi GPU Manager
- [DockerGPU](https://github.com/hequan2017/DockerGPU) — 前代 GPU 算力租用管理系统

## 📄 License

本项目基于 [Apache-2.0](./LICENSE) 协议开源。
