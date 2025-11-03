package engine

import (
	"encoding/json"
	"errors"
	"github.com/GenJi77JYXC/intelligent-pioneer/internal/logger"
	"github.com/GenJi77JYXC/intelligent-pioneer/internal/model"
	"github.com/GenJi77JYXC/intelligent-pioneer/internal/store"
	"github.com/google/uuid"
	"time"
)

// StartKBWorkflow 是启动知识库工作流的入口
func StartKBWorkflow(agentID, kbID string) (string, error) {
	logger.L.Infow("Starting KB workflow", "agent_id", agentID, "kb_id", kbID)

	// 1. 创建并存储工作流状态到数据库 (PostgreSQL)
	workflow := &model.Workflow{
		ID:        uuid.NewString(), // 生成工作流唯一ID
		KBID:      kbID,
		AgentID:   agentID,
		Status:    "pending",
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
	if err := store.DB.Create(workflow).Error; err != nil {
		logger.L.Errorw("Failed to create workflow record", "error", err)
		return "", err
	}
	// 2. 从 Elasticsearch 中获取知识库条目
	// 假设你有一个函数来获取KB条目
	kbItem, err := getKBItemFromES(kbID)
	if err != nil {
		logger.L.Errorw("Failed to get KB item from Elasticsearch", "kb_id", kbID, "error", err)
		return "", err
	}
	if len(kbItem.Diagnostics) == 0 {
		return "", errors.New("KB item has no diagnostic steps")
	}

	// 3. 提交第一个诊断任务
	firstDiagnosticStep := kbItem.Diagnostics[0]
	task := &Task{
		ID:         uuid.NewString(),
		AgentID:    agentID,
		WorkflowID: workflow.ID,
		Type:       "diagnostic",
		Command:    firstDiagnosticStep["command"], // 假设诊断步骤的格式是 {"command": "..."}
		CreatedAt:  time.Now(),
	}

	// 4. 更新工作流状态为 "diagnosing"
	updateData := map[string]interface{}{
		"status":          "diagnosing",
		"current_task_id": task.ID,
	}
	if err := store.DB.Model(&model.Workflow{}).Where("id = ?", workflow.ID).Updates(updateData).Error; err != nil {
		logger.L.Errorw("Failed to update workflow status to diagnosing", "error", err)
		return "", err
	}

	TM.SubmitTask(task)

	return workflow.ID, nil
}

// HandleTaskResult 是处理 Agent 返回结果的入口
func HandleTaskResult(result *TaskResult) {
	logger.L.Infow("Handling task result", "task_id", result.TaskID, "success", result.Success)

	// 1. 根据 result.TaskID 找到对应的工作流 (workflow)
	// 我们假设一个 Agent 的一个任务只属于一个工作流
	var workflow model.Workflow
	dbResult := store.DB.Where("agent_id = ? AND current_task_id = ?", result.AgentID, result.TaskID).First(&workflow)
	if dbResult.Error != nil {
		logger.L.Errorw("Cannot find workflow for this task result", "task_id", result.TaskID, "agent_id", result.AgentID, "error", dbResult.Error)
		return
	}

	// 2. 从ES中再次获取KB条目

	kbItem, err := getKBItemFromES(workflow.KBID)
	if err != nil {
		logger.L.Errorw("Cannot find KB item for task result", "kb_id", workflow.KBID, "task_id", result.TaskID, "error", err)
		// 更新工作流状态为 "failed"
		updateWorkflowStatus(workflow.ID, "failed")
		return
	}

	// 3. 执行分析逻辑 (analysis_logic)
	switch workflow.Status {
	case "diagnosing":
		handleDiagnosingResult(result, &workflow, kbItem)
	case "remediating":
		handleRemediatingResult(result, &workflow, kbItem)
	default:
		logger.L.Warnw("Received task result for a workflow in an unexpected state", "workflow_id", workflow.ID, "status", workflow.Status)
	}
}

// getKBItemFromES 是一个示例函数，用于从ES获取KB条目
func getKBItemFromES(kbID string) (*KnowledgeBaseItem, error) {
	// 实际项目中，索引名应该来自配置
	res, err := store.ESClient.Get("pioneer-knowledge-base", kbID)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	if res.IsError() {
		return nil, errors.New("document not found or other error")
	}

	var response struct {
		Source *KnowledgeBaseItem `json:"_source"`
	}
	if err := json.NewDecoder(res.Body).Decode(&response); err != nil {
		return nil, err
	}

	return response.Source, nil
}

// handleDiagnosingResult 处理诊断任务的结果
func handleDiagnosingResult(result *TaskResult, workflow *model.Workflow, kbItem *KnowledgeBaseItem) {
	// 分析逻辑 (MVP: 仅判断 success)
	if result.Success {
		logger.L.Info("Diagnostic step succeeded. Proceeding to remediation.")

		if kbItem.Remediation == nil || kbItem.Remediation["command"] == "" {
			logger.L.Info("No remediation step. Workflow completed.", "workflow_id", workflow.ID)
			updateWorkflowStatus(workflow.ID, "completed")
			return
		}

		// 构造并提交修复任务
		remediationTask := &Task{
			ID:         uuid.NewString(),
			AgentID:    workflow.AgentID,
			WorkflowID: workflow.ID,
			Type:       "remediation",
			Command:    kbItem.Remediation["command"],
			CreatedAt:  time.Now(),
		}
		// 更新工作流状态为 "remediating"
		updateData := map[string]interface{}{
			"status":          "remediating",
			"current_task_id": remediationTask.ID,
		}
		if err := store.DB.Model(workflow).Updates(updateData).Error; err != nil {
			logger.L.Errorw("Failed to update workflow to remediating", "error", err)
			return
		}

		TM.SubmitTask(remediationTask)

	} else {
		logger.L.Errorw("Diagnostic step failed", "workflow_id", workflow.ID, "output", result.Output)
		// 更新工作流状态为 "failed"
		updateWorkflowStatus(workflow.ID, "failed")
	}
}

// handleRemediatingResult 处理修复任务的结果
func handleRemediatingResult(result *TaskResult, workflow *model.Workflow, kbItem *KnowledgeBaseItem) {
	if result.Success {
		logger.L.Info("Remediation step succeeded. Workflow completed.", "workflow_id", workflow.ID)
		// 更新工作流状态为 "completed"
		updateWorkflowStatus(workflow.ID, "completed")
	} else {
		// 更新工作流状态为 "failed"
		logger.L.Errorw("Remediation step failed", "workflow_id", workflow.ID, "output", result.Output)
		updateWorkflowStatus(workflow.ID, "failed")
	}
}

// updateWorkflowStatus 是一个辅助函数，用于更新工作流状态
func updateWorkflowStatus(workflowID, status string) {
	updateData := map[string]interface{}{"status": status}
	if err := store.DB.Model(&model.Workflow{}).Where("id = ?", workflowID).Updates(updateData).Error; err != nil {
		logger.L.Errorw("Failed to update workflow status", "workflow_id", workflowID, "status", status, "error", err)
	}
}
