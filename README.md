# pcfarm-admin

`pcfarm-admin` 是基于当前 Gin + Vue 管理框架扩展的服务器管理系统，面向独立装机网段/VLAN 的服务器资产、IP 分配、PXE 启动、远程电源控制和 Ubuntu Live Agent 注册场景。

## 功能说明

- 服务器资产管理：维护资产编号、序列号、PXE MAC、BMC 地址、远控协议和在线状态。
- IP 地址池：从装机网段地址池为服务器 PXE MAC 自动分配长期固定 IP。
- PXE 启动策略：支持本地盘、Ubuntu Live、维护模式三种启动策略。
- 远程电源控制：预留 IPMI 和 Redfish provider，支持开机、关机、重启、下次 PXE 启动。
- Ubuntu Live Agent：Live 系统启动后可注册到管理端并上报心跳。
- 前端页面：服务器列表、服务器详情、IP 池管理、PXE 设置。

## 技术栈

- 后端：Go、Gin、GORM、JWT、Casbin、Swagger。
- 前端：Vue 3、Vite、Pinia、Element Plus。
- MVP PXE 方案：单节点 `dnsmasq` 控制面，后续可替换为 Kea DHCP 或独立 PXE Agent。

## 环境要求

- Go 1.24 或兼容版本。
- Node.js 和 npm。
- Windows PowerShell、Linux shell 或 macOS shell 均可运行开发命令。

## 启动后端

```powershell
cd server
go run .
```

默认配置读取 `server/config.yaml`，当前常用开发端口为：

- 后端 API：`http://127.0.0.1:8888`
- Swagger：`http://127.0.0.1:8888/swagger/index.html`

如需初始化数据库，请先确认 `server/config.yaml` 中的 `system.db-type`、数据库连接和 `disable-auto-migrate` 配置。

## 启动前端

```powershell
cd web
npm install
npm run serve
```

默认开发地址：

- 前端：`http://127.0.0.1:8080`
- 前端代理：`/api` 转发到 `http://127.0.0.1:8888`

## 构建前端

```powershell
cd web
npm run build
```

构建产物输出到 `web/dist/`。

## 验证命令

后端 pcfarm 聚焦测试：

```powershell
cd server
go test ./model/pcfarm ./service/pcfarm ./api/v1/pcfarm ./router/pcfarm ./initialize -count=1
```

前端构建验证：

```powershell
cd web
npm run build
```

## 当前实现边界

- 当前版本已实现 pcfarm 管理模块的基础模型、服务、API、路由和前端页面。
- PXE provider 和 IPMI/Redfish provider 当前为 MVP 边界实现，真实 `dnsmasq` 配置写入、服务重载、硬件远控命令执行需结合部署环境继续补齐。
- 菜单权限仍沿用原管理框架机制，需要在后台菜单管理中配置 pcfarm 页面入口。
