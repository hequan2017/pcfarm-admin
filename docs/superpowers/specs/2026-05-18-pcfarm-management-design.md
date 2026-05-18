# pcfarm 管理系统设计

## 目标

基于当前 gin-vue-admin 项目开发 `pcfarm` 管理系统，用于管理独立装机 VLAN 内的服务器资产、IP 分配、PXE 启动策略、远程电源控制，以及 Ubuntu Live 启动后的自动注册和在线状态。

第一版采用单节点部署：pcfarm 后端所在节点同时作为 DHCP、TFTP、HTTP/NFS 等 PXE 服务节点。实现上优先使用 `dnsmasq` 一体化管理 DHCP/TFTP，代码边界保留后续替换 Kea DHCP 或拆分 PXE Agent 的能力。

## 已确认需求

- pcfarm 直接管理 DHCP/TFTP/PXE/IP 分配等底层装机网络能力。
- 服务器处于独立装机网段或 VLAN，pcfarm 可以接管该网段 DHCP。
- 远程控制同时支持 IPMI 和 Redfish，按服务器配置选择驱动。
- 服务器每次通过 PXE 启动进入全新的 Ubuntu Live 环境，不依赖本地盘状态。
- Ubuntu Live 启动后运行 agent，自动上报身份、IP、硬件信息和在线状态。
- 资产长期身份以序列号或资产编号为主，PXE/DHCP 绑定以 PXE 网卡 MAC 为主。
- MVP 固定一套 Ubuntu Live 镜像，不做多镜像版本管理。
- IP 分配采用地址池自动分配，录入 PXE MAC 后长期绑定固定 IP。
- pcfarm 部署节点同时运行 PXE 相关服务。
- 同时支持 UEFI PXE 和 Legacy BIOS PXE。
- 每台服务器有启动策略：本地盘、Ubuntu Live、维护模式。

## 非目标

- 不做多 Ubuntu 镜像版本管理。
- 不做多 PXE 服务节点。
- 不做跨 VLAN DHCP Relay 管理。
- 不做高级任务调度。
- 不做自动重装到本地盘。
- 不把 pcfarm 第一版拆成独立 agent 架构。

## 架构方案

推荐方案为 `dnsmasq` 一体化 MVP：

- pcfarm 后端管理服务器资产、IP 池、启动策略和远控动作。
- 后端根据数据库期望状态生成 `dnsmasq`、TFTP、GRUB/iPXE/pxelinux 等配置。
- Ubuntu Live 镜像由部署阶段放置到固定目录，通过 HTTP/NFS 提供。
- IPMI 和 Redfish 通过统一 `PowerProvider` 抽象封装。
- PXE 服务控制通过 `PXEProvider` 抽象封装，第一版实现为本机文件配置和服务重载。

该方案符合 KISS 和 YAGNI：先解决单装机网段、固定镜像、服务器可控启动的核心闭环，同时通过 provider 接口避免业务层绑定具体底层实现。

## 后端模块

新增 `pcfarm` 业务模块，沿用现有 `Router -> API -> Service -> Model` 分层。

建议目录：

```text
server/model/pcfarm/
server/model/pcfarm/request/
server/model/pcfarm/response/
server/service/pcfarm/
server/api/v1/pcfarm/
server/router/pcfarm/
```

如果后续确定采用插件化交付，再整体迁移到：

```text
server/plugin/pcfarm/
web/src/plugin/pcfarm/
```

第一版建议先作为稳定业务模块落地，避免为了插件化增加额外复杂度。

### 核心模型

`ServerAsset`

- 服务器资产记录。
- 字段包括资产编号、序列号、PXE MAC、BMC 地址、远控协议、启动策略、在线状态、最后心跳时间。
- 序列号或资产编号用于长期身份，PXE MAC 用于 DHCP/PXE 绑定。

`IPPool`

- 装机网段地址池。
- 字段包括网段、起始 IP、结束 IP、网关、DNS、绑定网卡、启用状态。
- 第一版只允许启用一个装机网段，降低 DHCP 误配置风险。

`IPAllocation`

- PXE MAC 到 IP 的长期绑定。
- 字段包括服务器 ID、PXE MAC、IP、分配状态、分配时间、释放时间。
- 同一 MAC 和同一 IP 必须唯一。

`BootPolicy`

- 每台服务器当前启动策略。
- 枚举值：`local_disk`、`ubuntu_live`、`maintenance`。
- 策略变更触发 PXE 配置刷新。

`PowerCredential`

- BMC/IPMI/Redfish 连接配置。
- 密码必须加密存储，不直接返回前端。

`ProvisionEvent`

- 审计和诊断事件。
- 记录资产变更、IP 分配、PXE 配置刷新、远控动作、Agent 注册、心跳丢失等。

### Provider 抽象

`IPAllocator`

- 从启用地址池中为新 PXE MAC 分配固定 IP。
- 检查地址池范围、重复 MAC、重复 IP、保留地址。

`PXEProvider`

- 根据服务器和启动策略生成 DHCP host 绑定、TFTP 启动文件、PXE 菜单。
- 支持 UEFI 和 Legacy BIOS。
- 采用临时文件生成、配置校验、原子替换、服务重载流程。

`PowerProvider`

- 统一封装开机、关机、重启、设置下次 PXE 启动。
- 第一版实现 `IPMIProvider` 和 `RedfishProvider`。
- 操作失败只记录失败事件，不把业务状态更新为成功。

## 核心流程

### 资产录入

1. 管理员录入服务器资产编号、序列号、PXE MAC、BMC 地址、远控协议和凭据。
2. 系统校验 PXE MAC 唯一性。
3. 系统从 IP 池自动分配一个长期固定 IP。
4. 系统写入 `ServerAsset`、`IPAllocation` 和初始化 `BootPolicy`。
5. 系统记录资产创建和 IP 分配事件。

### 设置 Ubuntu Live 启动

1. 管理员将服务器启动策略切换为 `ubuntu_live`。
2. 后端更新策略并调用 `PXEProvider` 刷新 PXE 配置。
3. 管理员选择立即重启或下次启动生效。
4. 如选择立即执行，后端调用 `PowerProvider` 设置下次 PXE 启动并重启。
5. 服务器从 PXE 启动进入 Ubuntu Live。

### Live Agent 注册

1. Ubuntu Live 启动后自动运行 pcfarm agent。
2. agent 调用注册接口，上报序列号、PXE MAC、IP、硬件信息和 agent 版本。
3. 后端按序列号和 PXE MAC 匹配资产。
4. 匹配成功后更新在线状态和最后心跳时间。
5. agent 周期性发送心跳，心跳超时后后端标记为心跳丢失。

### 启动策略

`local_disk`

- PXE 菜单返回本地盘启动，或不给该 MAC 下发 Ubuntu Live 配置。

`ubuntu_live`

- PXE 菜单引导固定 Ubuntu Live 镜像。

`maintenance`

- 预留维护模式入口，可先指向同一套 Ubuntu Live 镜像并传入不同内核参数。
- 第一版只保留策略枚举和配置生成入口，不扩展复杂维护功能。

## 前端页面

### 服务器资产

- 列表字段：资产编号、序列号、PXE MAC、固定 IP、远控协议、启动策略、在线状态、最后心跳时间。
- 操作：新增、编辑、删除、分配 IP、释放 IP、切换启动策略、开机、关机、重启、设置下次 PXE。
- 支持批量开机、关机、重启、进入 PXE。

### 服务器详情

- 展示资产信息、IP 分配、BMC 配置摘要、Agent 状态、硬件信息、事件时间线。
- 密码字段只允许重置，不显示明文。

### IP 池管理

- 展示装机网段、起止 IP、网关、DNS、绑定网卡、已分配数量、可用数量。
- 第一版限制一个启用地址池。

### PXE 设置

- 展示 Ubuntu Live 镜像路径、TFTP 根目录、HTTP/NFS 地址、dnsmasq 服务状态。
- 提供刷新配置、校验配置、查看最近错误入口。

## API 边界

API 层负责 HTTP 参数绑定、校验、调用 Service 和统一响应；Service 不依赖 `gin.Context`。

接口建议按以下资源组织：

- `/pcfarm/server/*`
- `/pcfarm/ipPool/*`
- `/pcfarm/bootPolicy/*`
- `/pcfarm/power/*`
- `/pcfarm/pxe/*`
- `/pcfarm/agent/*`
- `/pcfarm/event/*`

对外接口必须补充 Swagger 注释，并保持项目统一响应结构 `{ code, data, msg }`。

## 安全与风险控制

- DHCP 只能绑定明确配置的装机 VLAN 网卡。
- 不允许在未配置装机网卡时启动 DHCP 管理能力。
- 配置刷新必须先校验再替换，失败时保留旧配置。
- IP 分配必须检查范围、重复 MAC、重复 IP 和保留地址。
- 远控动作必须记录审计事件。
- Live Agent 注册必须使用 token 或签名校验。
- BMC 密码加密存储，接口不返回明文。
- 批量远控动作需要二次确认，避免误关机。

## 测试策略

后端单元测试：

- IP 池分配、耗尽、重复 MAC、重复 IP。
- 启动策略枚举和状态变更。
- PXE 配置生成的输入输出。
- PowerProvider 失败时事件记录和状态不误更新。
- Agent 注册匹配逻辑。

后端集成测试：

- 资产录入后自动分配 IP。
- 切换启动策略后触发 PXE 配置刷新。
- Agent 注册后服务器状态变为在线。

前端验证：

- 服务器列表筛选、分页、批量操作。
- 详情页事件时间线。
- IP 池剩余容量显示。
- 密码字段不回显。

部署验证：

- dnsmasq 配置校验通过。
- UEFI PXE 能进入 Ubuntu Live。
- Legacy BIOS PXE 能进入 Ubuntu Live。
- IPMI 和 Redfish 至少各验证一台服务器。

## 设计原则应用

- KISS：MVP 采用单节点 `dnsmasq`，不引入多节点控制面。
- YAGNI：固定一套 Ubuntu Live 镜像，不提前实现多镜像版本管理。
- DRY：远控、PXE、IP 分配通过 provider 和 allocator 统一入口，避免业务流程重复拼装底层命令。
- SOLID：Service 只处理业务语义，Provider 只处理外部系统适配，API 只处理 HTTP 边界。

## 自检结果

- 没有保留未决占位符。
- MVP 范围与已确认需求一致。
- 架构没有要求第一版实现多节点或多镜像。
- 启动策略、身份绑定、IP 分配和远控协议均已明确。
