package diffslog

import (
	"fmt"
	"vikunjabot/internal/webhook"
)

func MapWebhookMessage(event *webhook.WebhookMessage) any {
	var createdType CreationNodeType = 0 
	var updatedType UpdationNodeType = 0
	var deletedType DeletionNodeType = 0
	var immediateType ImmediateNodeType = 0

	var instanceId string = ""

	switch event.EventName {
	case "task.created":
		createdType = CreatedTask
		instanceId = fmt.Sprintf("TASK_%d", event.Data.Task.ID)
	case "task.assignee.created":
		createdType = CreatedAssignee
		instanceId = fmt.Sprintf("ASSIGNEE_%d", event.Data.Assignee.ID)
	case "task.comment.created":
		createdType = CreatedComment
		instanceId = fmt.Sprintf("COMMENT_%d", event.Data.Comment.ID)
	case "task.relation.created":
		createdType = CreatedRelation
		firstTask := min(event.Data.Relation.TaskID, event.Data.Relation.OtherTaskID)
		secondTask := max(event.Data.Relation.TaskID, event.Data.Relation.OtherTaskID)
		instanceId = fmt.Sprintf("RELATION_%d_%d", firstTask, secondTask)
	case "task.attachment.created":
		createdType = CreatedAttachment
		instanceId = fmt.Sprintf("ATTACHMENT_%d", event.Data.Attachment.ID)

	case "project.updated":
		updatedType = UpdatedProject
		instanceId = fmt.Sprintf("PROJECT_%d", event.Data.Project.ID)
	case "task.updated":
		updatedType = UpdatedTask
		instanceId = fmt.Sprintf("TASK_%d", event.Data.Task.ID)
	case "task.comment.edited":
		updatedType = UpdatedComment
		instanceId = fmt.Sprintf("COMMENT_%d", event.Data.Comment.ID)

	case "project.deleted":
		deletedType = DeletedProject
		instanceId = fmt.Sprintf("PROJECT_%d", event.Data.Project.ID)
	case "task.deleted":
		deletedType = DeletedTask
		instanceId = fmt.Sprintf("TASK_%d", event.Data.Task.ID)
	case "task.assignee.deleted":
		deletedType = DeletedAssignee
		instanceId = fmt.Sprintf("ASSIGNEE_%d", event.Data.Assignee.ID)
	case "task.comment.deleted":
		deletedType = DeletedComment
		instanceId = fmt.Sprintf("COMMENT_%d", event.Data.Comment.ID)
	case "task.relation.deleted":
		deletedType = DeletedRelation
		firstTask := min(event.Data.Relation.TaskID, event.Data.Relation.OtherTaskID)
		secondTask := max(event.Data.Relation.TaskID, event.Data.Relation.OtherTaskID)
		instanceId = fmt.Sprintf("RELATION_%d_%d", firstTask, secondTask)
	case "task.attachment.deleted":
		deletedType = DeletedAttachment
		instanceId = fmt.Sprintf("ATTACHMENT_%d", event.Data.Attachment.ID)

	case "project.shared.team": immediateType = ImmediateProjectSharedTeam
	case "project.shared.user": immediateType = ImmediateProjectSharedUser
	case "task.overdue": immediateType = ImmediateTaskOverdue
	case "task.reminder.fired": immediateType = ImmediateReminderFired
	case "tasks.overdue": immediateType = ImmediateTasksOverdue

	default: return nil
	}

	if createdType != 0 {
		return CreationNode{
			Type: createdType,
			ProjectID: event.Data.GetProjects()[0].ID,
			TaskID: event.Data.GetTasks()[0].ID,
			InstanceID: instanceId,
			MessageData: &event.Data,
		}
	}

	if updatedType != 0 {
		return UpdationNode{
			Type: updatedType,
			ProjectID: event.Data.GetProjects()[0].ID,
			TaskID: event.Data.GetTasks()[0].ID,
			InstanceID: instanceId,
			MessageData: &event.Data,
		}
	}

	if deletedType != 0 {
		return DeletionNode{
			Type: deletedType,
			ProjectID: event.Data.GetProjects()[0].ID,
			TaskID: event.Data.GetTasks()[0].ID,
			InstanceID: instanceId,
			MessageData: &event.Data,
		}
	}

	if immediateType != 0 {
		return ImmediateNode{
			Type: immediateType,
			ProjectID: event.Data.GetProjects()[0].ID,
			MessageData: &event.Data,
		}
	}

	return nil
}
