# Intelligent-Pioneer (智能先锋)

`Intelligent-Pioneer` 是一个基于 **双引擎驱动** 的现代化终端安全运维智能体项目。它旨在通过融合 **数据驱动的AI预测** 与 **经验驱动的知识库决策**，将终端IT运维从被动响应的“救火”模式，革新为无人值守的“预防与自愈”模式。

## ✨ 项目愿景

在数字化时代，终端设备是企业运营的神经末梢。本项目致力于打造一个能够自主感知、分析、决策并行动的“智能运维先锋”，它不仅能像气象员一样精准预测系统风险，更能像经验丰富的工程师一样，高效地自动化诊断并修复常见故障，最终实现：

*   **降低运维成本 (Reduce Costs):** 将运维人员从重复、繁琐的工作中解放出来。
*   **提升系统稳定性 (Enhance Stability):** 变被动响应为主动预防，保障业务连续性。
*   **沉淀专家经验 (Accumulate Knowledge):** 将运维经验固化为可执行的知识库，构建可成长的智能运维体系。

## 🚀 核心特性

- **双引擎驱动:**
    - **🧠 AI 预测引擎:** 基于时序数据和机器学习模型，提前预警内存泄漏、性能瓶颈等潜在风险。
    - **📚 知识库引擎:** 将标准操作流程 (SOP) 结构化，实现对网络、客户端等已知问题的自动化诊断与修复。
- **轻量级终端智能体 (Agent):** 使用 Go 语言开发，性能卓越，资源占用低，跨平台兼容。
- **现代化技术栈:** 后端采用 Go (Gin)，前端采用 Vue3 (Vite)，数据库采用 PostgreSQL + Elasticsearch，消息队列采用 Kafka，实现高内聚、低耦合的微服务架构。
- **容器化部署:** 所有依赖服务均通过 Docker Compose 一键启动，实现开发环境的快速搭建与隔离。
- **可视化管理:** 提供友好的 Web 界面，用于监控终端状态、管理知识库、追踪自动化任务流。

## 🛠️ 技术栈

| 分类       | 技术                                       |
| :--------- | :----------------------------------------- |
| **后端**   | Go, Gin                                    |
| **前端**   | Vue.js 3, Vite, TypeScript, Element Plus |
| **数据库** | PostgreSQL, Elasticsearch                  |
| **消息队列** | Kafka                                      |
| **容器化** | Docker, Docker Compose                     |
| **DevOps** | Makefile                                   |

## 🏁 快速开始

在开始之前，请确保你的本地环境已经安装了 [Go (1.18+)](https://golang.org/), [Docker](https://www.docker.com/), [Docker Compose](https://docs.docker.com/compose/) 和 `make`。

### 1. 克隆项目

```bash
git clone https://github.com/GenJi77JYXC/intelligent-pioneer.git
cd intelligent-pioneer
```

### 2. 启动依赖服务

此命令将通过 Docker Compose 在后台启动 PostgreSQL, Elasticsearch 和 Kafka。

```bash
make docker-up
```
第一次启动会需要一些时间来拉取镜像。你可以通过 `docker ps` 命令检查所有容器是否都处于 `Up` 状态。

### 3. 安装项目依赖

此命令会下载 Go 模块并将其放入 `vendor` 目录。

```bash
make deps
```

### 4. 运行后端服务

```bash
make run
```
如果一切顺利，你将在终端看到 `🚀 Starting Intelligent-Pioneer...` 的日志，并且后端服务将运行在 `http://localhost:8080`。

### 5. 运行前端（待开发）

*当前为后端先行阶段，前端开发将在后续迭代中进行。*

```bash
# 进入前端代码目录 (未来)
# cd frontend

# 安装 npm 依赖
# npm install

# 启动开发服务器
# npm run dev
```

## 📂 项目结构说明

```
.
├── cmd/                # 程序主入口
├── configs/            # 配置文件
├── internal/           # 项目内部私有代码 (核心业务逻辑)
│   ├── agent/          # 终端智能体相关逻辑
│   ├── api/            # Gin 路由和 Handler
│   ├── config/         # 配置加载
│   ├── core/           # 核心引擎 (AI, KB, Decision)
│   └── store/          # 数据库交互层
├── pkg/                # 可被外部项目引用的公共代码
├── api/                # API 定义文件 (如 protobuf, openapi)
├── web/                # (未来) 前端Vue代码
├── .gitignore          # Git 忽略文件
├── docker-compose.yml  # Docker 服务编排
├── go.mod              # Go 模块管理
├── Makefile            # 项目自动化构建脚本
└── README.md           # 就是你正在看的这个文件
```

## 🗺️ 开发路线图 (Roadmap)

- **[✓] Phase 0: 项目初始化**
    - [x] 确定技术栈与架构
    - [x] 初始化项目结构
    - [x] 配置开发环境 (Docker, Makefile)

- **[▶️] Phase 1: 基础平台与知识库MVP (进行中)**
    - [ ] Agent: 实现心跳、注册、指令执行能力
    - [ ] Backend: 实现面向Agent的API, 搭建数据管道
    - [ ] KB Engine: 开发知识库核心工作流 (匹配->诊断->修复)
    - [ ] Frontend: 实现Agent列表和事件流展示页
    - **目标:** 跑通第一个端到端的自动化修复案例。

- **[ ] Phase 2: 完善知识库生态与告警集成**
    - [ ] 提供可视化的知识库管理界面
    - [ ] 对接日志、监控系统，实现自动触发
    - [ ] 扩充知识库至50+条目

- **[ ] Phase 3: AI预测引擎的引入**
    - [ ] 训练并部署第一个AI预测模型 (如内存泄漏)
    - [ ] 在前端展示预测性风险
    - [ ] 开发决策与仲裁引擎

- **[ ] Phase 4: 系统成熟化与优化**
    - [ ] AI与KB双引擎协同工作
    - [ ] 性能优化与稳定性加固
    - [ ] 完善文档与测试用例

## 🤝 贡献

我们欢迎任何形式的贡献！如果你有好的想法或建议，请随时提交。

## 📄 开源许可
本项目采用 [MIT License](./LICENSE) 开源许可。