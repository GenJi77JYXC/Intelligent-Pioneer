# 后端业务逻辑DO

1. **实现 handler_internal.go 中的 TriggerKB**：让它调用 engine.StartKBWorkflow。

2. **实现 handler_agent.go 中的 GetTasks**：让它调用 engine.TM.GetTaskForAgent。

3. **实现 handler_agent.go 中的 PostTaskResults**：让它解析请求体中的结果，并调用 engine.HandleTaskResult。

4. **（可选，但推荐）** 在 engine/model.go 中定义 Workflow 结构体，并在 store/postgres.go 的 migrateDatabase 中添加 &model.Workflow{}，然后在 StartKBWorkflow 和 HandleTaskResult 中实现对数据库的读写，以持久化追踪工作流的状态。

5. **设计一个离线检测机制（未来工作）**

   Heartbeat API 只能将 Agent 标记为 online。那么，如何将一个掉线的 Agent 标记为 offline 呢？

   这通常不是由 API 直接完成的，而是通过一个**后台定时任务**来实现。这里给出设计思路，你可以作为后续的迭代任务：

   1. **创建一个后台任务（Go Cron Job）:**
      - 使用类似 robfig/cron 的库，创建一个定时任务，比如每 5 分钟执行一次。
   2. **编写检测逻辑:**
      - 这个任务会执行一条 SQL 查询，查找所有 status = 'online' 并且 updated_at 在**5分钟之前**（即超过一个心跳周期没有更新）的 Agent。
      - SQL 示例: UPDATE agents SET status = 'offline' WHERE status = 'online' AND updated_at < NOW() - INTERVAL '5 minutes';
   3. **执行更新:**
      - 将查询到的这些“超时”的 Agent 状态批量更新为 offline。

   这个机制确保了即使 Agent 异常掉线（断电、断网），后端也能在一定延迟后准确地反映出它的离线状态。



**注册Agent**

```bash
 curl -X POST http://localhost:8080/api/v1/agent/register \
-H "Content-Type: application/json" \
-d '{
  "hostname": "my-dev-machine",
  "ip_address": "192.168.1.100",
  "os": "Ubuntu 22.04"
}'
```

**触发工作流**

```bash
curl -X POST http://localhost:8080/api/v1/internal/trigger_kb \
-H "Content-Type: application/json" \
-d '{
  "agent_id": "b180b1bc-c384-4dea-84b2-778c80a20a99",
  "kb_id": "dns-flush-kb"
}'
```



```bash
curl -X POST http://localhost:8080/api/v1/agent/heartbeat \
-H "Content-Type: application/json" \
-d '{
  "agent_id": "b180b1bc-c384-4dea-84b2-778c80a20a99"
}'

curl -X POST http://localhost:8080/api/v1/agent/heartbeat \
-H "Content-Type: application/json" \
-d '{
  "agent_id": "this-is-a-fake-id"
}'
```

**测试长轮询**

```bash
curl -v "http://localhost:8080/api/v1/agent/tasks?agent_id=b180b1bc-c384-4dea-84b2-778c80a20a99"


curl -X POST http://localhost:8080/api/v1/internal/trigger_kb \
-H "Content-Type: application/json" \
-d '{
  "agent_id": "b180b1bc-c384-4dea-84b2-778c80a20a99",
  "kb_id": "dns-flush-kb"
}'

```

