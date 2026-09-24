package diffslog

import (
	"vikunjabot/internal/webhook"
)

type DiffsLog struct {
	creations []CreationNode
	updations []UpdationNode
	deletions []DeletionNode
	immediates []ImmediateNode
}

func NewDiffsLog() *DiffsLog {
	return &DiffsLog{}
}

func (a *DiffsLog) AddEvent(event *webhook.WebhookMessage) {
	mapped := MapWebhookMessage(event)
	if mapped != nil {
		switch node := mapped.(type) {
		case CreationNode: 
			a.creations = append(a.creations, node)
		case UpdationNode:
			a.updations = append(a.updations, node)
		case DeletionNode:
			a.deletions = append(a.deletions, node)
		case ImmediateNode:
			a.immediates = append(a.immediates, node)
		}
	}
}

func (a *DiffsLog) GetLogFlow() LogFlow {
	logFlow := LogFlow{}

	for _, cNode := range a.creations {
		instanceType := InstanceUnknown

		switch cNode.Type {
		case CreatedAssignee: instanceType = InstanceAssignee
		case CreatedTask: instanceType = InstanceTask
		case CreatedAttachment: instanceType = InstanceAttachment
		case CreatedComment: instanceType = InstanceComment
		case CreatedRelation: instanceType = InstanceRelation
		default: continue
		}

		logFlow.Instances = append(logFlow.Instances, LogEntry{
			Action: ActionCreate,
			MessageData: cNode.MessageData,
			ProjectID: cNode.ProjectID,
			InstanceID: cNode.InstanceID,
			TaskID: cNode.TaskID,
			Instance: instanceType,
		})
	}

	for _, uNode := range a.updations {
		instanceType := InstanceUnknown

		switch uNode.Type {
		case UpdatedTask: instanceType = InstanceTask
		case UpdatedComment: instanceType = InstanceComment
		case UpdatedProject: instanceType = InstanceProject
		default: continue
		}

		logFlow.Instances = append(logFlow.Instances, LogEntry{
			Action: ActionUpdate,
			MessageData: uNode.MessageData,
			ProjectID: uNode.ProjectID,
			InstanceID: uNode.InstanceID,
			TaskID: uNode.TaskID,
			Instance: instanceType,
		})
	}

	for _, dNode := range a.deletions {
		instanceType := InstanceUnknown

		switch dNode.Type {
		case DeletedProject: instanceType = InstanceProject
		case DeletedAssignee: instanceType = InstanceAssignee
		case DeletedAttachment: instanceType = InstanceAttachment
		case DeletedTask: instanceType = InstanceTask
		case DeletedComment: instanceType = InstanceComment
		case DeletedRelation: instanceType = InstanceRelation
		default: continue
		}

		logFlow.Instances = append(logFlow.Instances, LogEntry{
			Action: ActionDelete,
			MessageData: dNode.MessageData,
			ProjectID: dNode.ProjectID,
			InstanceID: dNode.InstanceID,
			TaskID: dNode.TaskID,
			Instance: instanceType,
		})
	}

	logFlow.Events = a.immediates

	return logFlow
}
