# MVP 自动化工作流端到端测试指南

本文档旨在提供一个完整的、端到端的测试流程，用于验证 `Intelligent-Pioneer` 后端自动化引擎的核心功能闭环。测试将完全通过 `curl` 命令模拟 Agent 的行为来完成。

## 🎯 测试目标

验证一个完整的“**触发 -> 诊断 -> 修复 -> 结束**”的自动化工作流能否成功执行。

## ✅ 预备条件

1.  **后端服务已启动:**
    
    ```bash
    make run
    ```
2.  **依赖服务已启动:**
    
    ```bash
    make docker-up
    ```
3.  **知识库条目已存在:** 确保 Elasticsearch 中已创建 `dns-flush-kb` 条目。如果不存在，请通过 Kibana Dev Tools 执行以下命令：
    ```json
    PUT pioneer-knowledge-base/_doc/dns-flush-kb
    {
      "diagnostics": [
        {
          "command": "ping -c 1 baidu.com"
        }
      ],
      "analysis_logic": "if exit_code == 0 return 'success'",
      "remediation": {
        "command": "echo 'DNS flushed successfully! (mock)'"
      }
    }
    ```

## 🧪 测试步骤

请按照以下顺序，在你的终端中逐步执行 `curl` 命令。建议**打开两个终端窗口**：一个用于运行后端服务并观察日志，另一个用于执行 `curl` 命令。

---

### **步骤 1: 注册 Agent**

我们首先需要一个合法的 Agent 身份。

**执行命令:**
```bash
curl -X POST http://localhost:8080/api/v1/agent/register \
-H "Content-Type: application/json" \
-d '{
  "hostname": "test-agent-01",
  "ip_address": "127.0.0.1",
  "os": "Test OS"
}'
```

**预期结果:**
你会收到一个包含 `agent_id` 的 JSON 响应。**请复制这个 `agent_id` 的值（不包括引号），我们将在后续所有步骤中使用它。**

**示例响应:**
```json
{
  "agent_id": "c1f7b8e2-a3d4-4b5c-8e9f-0a1b2c3d4e5f",
  "message": "Agent registered successfully."
}
```



**发送心跳包**

```json
curl -X POST http://localhost:8080/api/v1/agent/heartbeat \
-H "Content-Type: application/json" \
-d '{
  "agent_id": "242bfb75-0d19-4f51-91cb-8541156673c8"
}'
```



---

### **步骤 2: 触发自动化工作流**

现在，我们手动为刚刚注册的 Agent 触发一个知识库任务。

**执行命令 (请将 `<YOUR_AGENT_ID>` 替换为上一步复制的 ID):**
```bash
AGENT_ID="<YOUR_AGENT_ID>"

curl -X POST http://localhost:8080/api/v1/internal/trigger_kb \
-H "Content-Type: application/json" \
-d '{
  "agent_id": "'"$AGENT_ID"'",
  "kb_id": "dns-flush-kb"
}'

curl -X POST http://localhost:8080/api/v1/internal/trigger_kb \
-H "Content-Type: application/json" \
-d '{
  "agent_id": "242bfb75-0d19-4f51-91cb-8541156673c8",
  "kb_id": "dns-flush-kb"
}'
```

**预期结果:**
*   收到 `200 OK` 响应，表示任务已成功触发。
*   在后端服务的日志中，你会看到 `Starting KB workflow` 和 `Submitting new task to queue` 的日志。

---

### **步骤 3: Agent 获取“诊断”任务**

模拟 Agent 发起长轮询，获取它的第一个任务。

**执行命令 (请将 `<YOUR_AGENT_ID>` 替换为你的 ID):**
```bash
AGENT_ID="<YOUR_AGENT_ID>"

curl -i "http://localhost:8080/api/v1/agent/tasks?agent_id=$AGENT_ID"

curl -i "http://localhost:8080/api/v1/agent/tasks?agent_id=242bfb75-0d19-4f51-91cb-8541156673c8"
```

**预期结果:**

*   这个命令会**立即**返回，而不是等待30秒。
*   你会收到 `HTTP/1.1 200 OK` 的响应。
*   响应体是一个 JSON，包含了“诊断”任务的详细信息。**请复制响应中的 `ID` 字段的值（任务ID），下一步会用到。**

**示例响应:**
```json
{
    "ID": "f0e9d8c7-b6a5-4b4c-8a9b-1c2d3e4f5a6b",
    "AgentID": "c1f7b8e2-...",
    "Type": "diagnostic",
    "Command": "ping -c 1 baidu.com",
    "CreatedAt": "..."
}
```

---

### **步骤 4: Agent 上报“诊断”任务结果**

我们模拟 Agent 成功执行了诊断命令，并上报一个成功的结果。

**执行命令 (请将 `<YOUR_AGENT_ID>` 和 `<YOUR_TASK_ID>` 替换为真实的值):**
```bash
AGENT_ID="<YOUR_AGENT_ID>"
TASK_ID="<YOUR_TASK_ID>" # 从上一步获取的任务ID

curl -X POST http://localhost:8080/api/v1/agent/tasks/results \
-H "Content-Type: application/json" \
-d '{
  "task_id": "'"$TASK_ID"'",
  "agent_id": "'"$AGENT_ID"'",
  "success": true,
  "output": "PING baidu.com ...",
  "error": "",
  "exit_code": 0
}'

curl -X POST http://localhost:8080/api/v1/agent/tasks/results \
-H "Content-Type: application/json" \
-d '{
  "task_id": "'"70bf6085-ecb8-434c-8c5b-4f972fe3c597"'",
  "agent_id": "'"242bfb75-0d19-4f51-91cb-8541156673c8"'",
  "success": true,
  "output": "PING baidu.com ...",
  "error": "",
  "exit_code": 0
}'
```

**预期结果:**
*   收到 `200 OK` 响应，表示结果已收到。
*   在后端日志中，你会看到 `Received task result from agent`，紧接着是 `Handling task result in engine` 和 `Remediation task submitted` 的日志。

---

### **步骤 5: Agent 获取“修复”任务**

由于上一步诊断成功，引擎应该已经下发了新的“修复”任务。我们再次模拟 Agent 拉取任务。

**执行命令 (请将 `<YOUR_AGENT_ID>` 替换为你的 ID):**
```bash
AGENT_ID="<YOUR_AGENT_ID>"

curl -i "http://localhost:8080/api/v1/agent/tasks?agent_id=$AGENT_ID"

curl -i "http://localhost:8080/api/v1/agent/tasks?agent_id=242bfb75-0d19-4f51-91cb-8541156673c8"
```

**预期结果:**
*   这个命令同样会**立即**返回。
*   响应体是一个 JSON，包含了“修复”任务的详细信息。**请再次复制 `ID` 字段的值。**

**示例响应:**
```json
{
    "ID": "a9b8c7d6-e5f4-4c3d-8b2a-1b2c3d4e5f6g",
    "AgentID": "c1f7b8e2-...",
    "Type": "remediation",
    "Command": "echo 'DNS flushed successfully! (mock)'",
    "CreatedAt": "..."
}
```

---

### **步骤 6: Agent 上报“修复”任务结果**

我们模拟 Agent 成功执行了修复命令，并上报结果。

**执行命令 (请将 `<YOUR_AGENT_ID>` 和 `<YOUR_TASK_ID>` 替换为真实的值):**
```bash
AGENT_ID="<YOUR_AGENT_ID>"
TASK_ID="<YOUR_TASK_ID>" # 从上一步获取的任务ID b0d9006e-02ec-4391-b0c6-775b43b95a7e

curl -X POST http://localhost:8080/api/v1/agent/tasks/results \
-H "Content-Type: application/json" \
-d '{
  "task_id": "'"$TASK_ID"'",
  "agent_id": "'"$AGENT_ID"'",
  "success": true,
  "output": "DNS flushed successfully! (mock)",
  "error": "",
  "exit_code": 0
}'

curl -X POST http://localhost:8080/api/v1/agent/tasks/results \
-H "Content-Type: application/json" \
-d '{
  "task_id": "'"31718fbc-d451-44ef-a6de-c9d6920b49ce"'",
  "agent_id": "'"242bfb75-0d19-4f51-91cb-8541156673c8"'",
  "success": true,
  "output": "DNS flushed successfully! (mock)",
  "error": "",
  "exit_code": 0
}'
```

**预期结果:**
*   收到 `200 OK` 响应。
*   在后端日志中，你会看到 `Handling task result...` 和 `No remediation step found. Workflow completed.` 的日志，因为修复任务之后没有更多步骤了。

---

### **步骤 7: 验证工作流已结束**

现在工作流已经完成，Agent 应该再也获取不到新任务了。我们来验证一下。

**执行命令 (请将 `<YOUR_AGENT_ID>` 替换为你的 ID):**
```bash
AGENT_ID="<YOUR_AGENT_ID>"

curl -i "http://localhost:8080/api/v1/agent/tasks?agent_id=$AGENT_ID"

curl -i "http://localhost:8080/api/v1/agent/tasks?agent_id=242bfb75-0d19-4f51-91cb-8541156673c8"
```

**预期结果:**
*   这次，你的终端会**卡住**。
*   等待大约 **30 秒** 后，命令会结束。
*   你会看到 `HTTP/1.1 204 No Content` 的响应头，并且**没有响应体**。
*   在后端日志中，你会看到 `Polling timeout, no tasks for agent` 的日志。

---

**🎉 恭喜！** 如果你成功地完成了以上所有步骤，那就证明你的 `Intelligent-Pioneer` 后端自动化引擎的核心逻辑已经完全打通，并且工作正常！