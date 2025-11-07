### **Intelligent-Pioneer 项目初始化指南**

请按照以下步骤操作。我将提供一个可以直接复制粘贴到终端执行的脚本来完成大部分工作。

#### **第一步：创建项目目录并进入**

```bash
mkdir intelligent-pioneer
cd intelligent-pioneer
```

#### **第二步：初始化Go模块 (Go Module)**

这是Go项目管理的基石。模块路径通常是你的代码托管地址。

```bash
# 将 <your-username> 替换为你的GitHub用户名或其他代码托管平台的用户名
go mod init github.com/<your-username>/intelligent-pioneer
```
例如：`go mod init github.com/my-awesome-org/intelligent-pioneer`

#### **第三步：创建推荐的项目结构和初始文件 (一键执行)**

下面是一个Shell脚本，它会自动创建我们之前讨论过的、适合大型项目的目录结构，并生成必要的初始文件。

**直接复制下面的所有内容，粘贴到你的终端里，然后按回车执行。**

```bash
#!/bin/bash

# --- 创建核心目录结构 ---
echo "Creating directory structure..."
mkdir -p cmd/intelligent-pioneer
mkdir -p internal/agent internal/api internal/config internal/core/engine internal/store
mkdir -p pkg/utils
mkdir -p api
mkdir -p configs

# --- 创建 .gitignore ---
echo "Creating .gitignore..."
cat <<EOL > .gitignore
# Binaries for programs and plugins
*.exe
*.exe~
*.dll
*.so
*.dylib

# Test binary, built with 'go test -c'
*.test

# Output of the go coverage tool, specifically when run with the -o flag
*.out

# IDE files
.idea/
.vscode/

# Environment files
.env

# Build artifacts
bin/
vendor/
EOL

# --- 创建 Makefile 用于简化常用命令 ---
echo "Creating Makefile..."
cat <<EOL > Makefile
.PHONY: run build test clean deps docker-up docker-down

# Go parameters
GOCMD=go
GOBUILD=$(GOCMD) build
GOCLEAN=$(GOCMD) clean
GOTEST=$(GOCMD) test
GOGET=$(GOCMD) get
BINARY_NAME=intelligent-pioneer
BINARY_PATH=./bin/$(BINARY_NAME)

run:
	@echo "Running the application..."
	@$(GOCMD) run ./cmd/intelligent-pioneer/main.go

build:
	@echo "Building the application..."
	@$(GOBUILD) -o $(BINARY_PATH) ./cmd/intelligent-pioneer/main.go

test:
	@echo "Running tests..."
	@$(GOTEST) -v ./...

clean:
	@echo "Cleaning up..."
	@$(GOCLEAN)
	@rm -f $(BINARY_PATH)

deps:
	@echo "Installing dependencies..."
	@$(GOCMD) mod tidy
	@$(GOCMD) mod vendor

docker-up:
	@echo "Starting Docker services..."
	@docker-compose up -d

docker-down:
	@echo "Stopping Docker services..."
	@docker-compose down
EOL

# --- 创建初始的 main.go 入口文件 ---
echo "Creating main.go..."
cat <<'EOL' > cmd/intelligent-pioneer/main.go
package main

import "fmt"

func main() {
	fmt.Println("🚀 Starting Intelligent-Pioneer... The journey begins!")

	// TODO: 1. Load configuration (Viper)
	// TODO: 2. Initialize logger (Zap/Logrus)
	// TODO: 3. Initialize database connections (PostgreSQL, Elasticsearch)
	// TODO: 4. Initialize message queue producer/consumer (Kafka)
	// TODO: 5. Initialize HTTP server (Gin) and register routes
	// TODO: 6. Start the server and wait for shutdown signal
}
EOL

# --- 创建初始的配置文件 ---
echo "Creating default config.yaml..."
cat <<EOL > configs/config.yaml
server:
  port: "8080"
  mode: "debug" # debug, release, test

database:
  postgres:
    host: "localhost"
    port: "5432"
    user: "pioneer_user"
    password: "pioneer_password"
    dbname: "pioneer_db"
    sslmode: "disable"

  elasticsearch:
    addresses:
      - "http://localhost:9200"

kafka:
  brokers:
    - "localhost:9092"
EOL

# --- 创建 Docker Compose 文件以启动依赖服务 ---
echo "Creating docker-compose.yml..."
cat <<EOL > docker-compose.yml
version: '3.8'

services:
  postgres:
    image: postgres:14-alpine
    container_name: pioneer_postgres
    environment:
      POSTGRES_USER: pioneer_user
      POSTGRES_PASSWORD: pioneer_password
      POSTGRES_DB: pioneer_db
    ports:
      - "5432:5432"
    volumes:
      - postgres_data:/var/lib/postgresql/data
    restart: unless-stopped

  elasticsearch:
    image: docker.elastic.co/elasticsearch/elasticsearch:8.5.3
    container_name: pioneer_elasticsearch
    environment:
      - discovery.type=single-node
      - "ES_JAVA_OPTS=-Xms512m -Xmx512m"
      - xpack.security.enabled=false # Disable security for easy development
    ports:
      - "9200:9200"
      - "9300:9300"
    volumes:
      - es_data:/usr/share/elasticsearch/data
    restart: unless-stopped

  zookeeper:
    image: confluentinc/cp-zookeeper:7.3.0
    container_name: pioneer_zookeeper
    environment:
      ZOOKEEPER_CLIENT_PORT: 2181
      ZOOKEEPER_TICK_TIME: 2000
    ports:
      - "2181:2181"

  kafka:
    image: confluentinc/cp-kafka:7.3.0
    container_name: pioneer_kafka
    depends_on:
      - zookeeper
    ports:
      - "9092:9092"
    environment:
      KAFKA_BROKER_ID: 1
      KAFKA_ZOOKEEPER_CONNECT: zookeeper:2181
      KAFKA_ADVERTISED_LISTENERS: PLAINTEXT://kafka:29092,PLAINTEXT_HOST://localhost:9092
      KAFKA_LISTENER_SECURITY_PROTOCOL_MAP: PLAINTEXT:PLAINTEXT,PLAINTEXT_HOST:PLAINTEXT
      KAFKA_INTER_BROKER_LISTENER_NAME: PLAINTEXT
      KAFKA_OFFSETS_TOPIC_REPLICATION_FACTOR: 1

volumes:
  postgres_data:
  es_data:
EOL

echo ""
echo "✅ Project 'Intelligent-Pioneer' initialized successfully!"
echo ""
```

---

### **第四步：检查并开始你的项目**

执行完上面的脚本后，你的项目目录结构看起来会是这样：

```
intelligent-pioneer/
├── api/
├── cmd/
│   └── intelligent-pioneer/
│       └── main.go
├── configs/
│   └── config.yaml
├── internal/
│   ├── agent/
│   ├── api/
│   ├── config/
│   ├── core/
│   │   └── engine/
│   └── store/
├── pkg/
│   └── utils/
├── .gitignore
├── docker-compose.yml
├── go.mod
└── Makefile
```

#### **你的下一步 (Next Steps)**

1.  **启动依赖服务:**
    打开终端，确保你已经安装了Docker和Docker Compose，然后运行：
    
    ```bash
    make docker-up
    ```
    这将会在后台启动PostgreSQL, Elasticsearch, Kafka和Zookeeper。
    
2.  **安装初始依赖:**
    虽然我们还没写代码，但可以先整理一下`go.mod`文件：
    ```bash
    make deps
    ```

3.  **运行你的应用:**
    现在，你可以运行初始的`main.go`了：
    ```bash
    make run
    ```
    你应该会看到输出：`🚀 Starting Intelligent-Pioneer... The journey begins!`

4.  **开始编码！**
    你现在拥有了一个非常专业的项目起点。可以按照我们之前制定的 **Phase 1 任务清单** 开始了：
    
    *   **任务 1.1.3 (配置管理):** 在 `internal/config` 中，使用 [Viper](https://github.com/spf13/viper) 读取 `configs/config.yaml`。
    *   **任务 1.1.4 (日志系统):** 集成 [Zap](https://github.com/uber-go/zap) 并创建一个全局的Logger。
    *   **任务 1.1.2 (Web框架):** `go get github.com/gin-gonic/gin`，然后在 `main.go` 中启动一个简单的Gin服务器。

**恭喜！你的 `Intelligent-Pioneer` 项目已经正式启航。Happy coding!**