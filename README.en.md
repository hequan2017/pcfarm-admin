[简体中文](README.md) | [English](README.en.md)

# pcfarm-admin

A bare-metal provisioning system for dedicated install network segments/VLANs: server asset inventory, IP pool allocation, PXE boot policies, IPMI/Redfish remote power control, and Ubuntu Live Agent registration with heartbeats.

![Go](https://img.shields.io/badge/Go-1.24-00ADD8?logo=go&logoColor=white)
![Vue](https://img.shields.io/badge/Vue-3.5-4FC08D?logo=vue.js&logoColor=white)
![License](https://img.shields.io/badge/License-Apache--2.0-green)

## 📖 Introduction

Data-center batch provisioning is usually scattered across separate tools: spreadsheets for inventory, Excel for IPs, and hand-typed IPMI commands for power. `pcfarm-admin` brings this chain into one Gin + Vue admin console: centered on server assets, it allocates long-term fixed IPs from the install network pool by PXE MAC, pushes one of three boot policies per asset (local disk / Ubuntu Live / maintenance), and performs remote power actions through the BMC address with an encrypted stored password. Once the Live system boots, the Agent registers itself with the server and keeps reporting heartbeats, closing the loop of "register → allocate → provision → manage".

It suits ops teams that maintain dedicated install network segments/VLANs, and doubles as the bare-metal onboarding step in front of a GPU compute cluster (see the same author's [docker-gpu-manage](https://github.com/hequan2017/docker-gpu-manage)).

## ✨ Features

- 🖥️ Server asset management: asset code, serial number, PXE MAC, BMC address/account, encrypted password storage, power protocol (IPMI/Redfish), boot policy, status, agent version, hardware summary, last heartbeat time
- 🌐 IP address pools: pool management (name, CIDR, start/end IP, gateway, DNS, bound interface, enable toggle) with automatic long-term fixed IP allocation by PXE MAC and traceable allocation/release records
- 🥾 PXE boot policies: local disk (local_disk), Ubuntu Live, and maintenance modes; a local dnsmasq provider renders host bindings and boot menus, with Refresh / Status APIs
- 🔌 Remote power control: IPMI and Redfish providers supporting power on, power off, restart, and one-time PXE boot
- 🤖 Ubuntu Live Agent: registers the asset by serial number or PXE MAC after the Live system boots and keeps sending heartbeats, writing provision events
- 🧭 Frontend pages: server list/detail, IP pool management, PXE settings
- 📝 Provision event log: registration, allocation, and other key actions are recorded for troubleshooting and audit

## 🛠 Tech Stack

| Layer | Technologies |
|---|---|
| Backend | Go 1.24, Gin 1.10, GORM 1.31, Casbin v3, JWT, Swagger (swaggo) |
| Frontend | Vue 3.5, Vite 8, Element Plus 2.13, Pinia, Axios |
| PXE approach | MVP uses a single-node `dnsmasq` control plane, replaceable later with Kea DHCP or a standalone PXE Agent |

## 🚀 Quick Start

### Requirements

- Go 1.24+, Node.js and npm (works from Windows PowerShell, Linux, or macOS shells)

### Start the backend

```bash
cd server
go run .
```

It reads `server/config.yaml` by default: backend API at `http://127.0.0.1:8888`, Swagger at `http://127.0.0.1:8888/swagger/index.html`. To initialize the database, first check `system.db-type`, the database connection, and `disable-auto-migrate` in that file.

### Start the frontend

```bash
cd web
npm install
npm run serve
```

The dev server runs at `http://127.0.0.1:8080`, with `/api` proxied to `http://127.0.0.1:8888`.

### Build & Test

```bash
# Frontend build, output to web/dist/
cd web && npm run build

# Backend pcfarm-focused tests
cd server
go test ./model/pcfarm ./service/pcfarm ./api/v1/pcfarm ./router/pcfarm ./initialize -count=1
```

## 📁 Directory Structure

```text
├─ server/
│  ├─ model|service|api|router/pcfarm/   # assets, IP pool, PXE, Agent, power control
│  └─ initialize/
├─ web/src/view/pcfarm/                  # server / ipPool / pxe pages
├─ deploy/                               # docker / docker-compose / kubernetes
├─ docs/  ├─ aiDoc/                      # project docs and AI collaboration doc layer
└─ AGENT.MD
```

## ⚠️ Current Implementation Boundary

- The pcfarm module's base models, services, APIs, routes, and frontend pages are implemented.
- The PXE provider and IPMI/Redfish providers are MVP-scoped: real `dnsmasq` config writing, service reload, and hardware power command execution need to be completed for your deployment environment.
- Menu permissions follow the underlying admin framework; configure the pcfarm page entries in the admin menu management UI.

## 🔗 Related Projects

Sibling projects in the same GPU/compute management direction:

- [docker-gpu-manage](https://github.com/hequan2017/docker-gpu-manage) — TianQi compute management platform (GPU container scheduling, HAMi VRAM splitting, K8s management)
- [kapigpu](https://github.com/hequan2017/kapigpu) — KaPi GPU management platform (lightweight GVA edition)
- [tianqi](https://github.com/hequan2017/tianqi) — TianQi GPU Manager
- [DockerGPU](https://github.com/hequan2017/DockerGPU) — the previous-generation GPU rental management system

## 📄 License

Released under the [Apache-2.0](./LICENSE) license.
