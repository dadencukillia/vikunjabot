package diffslog

import "vikunjabot/internal/webhook"

// Log flow

type ActionType uint8
const (
	ActionUnknown ActionType = iota
	ActionCreate
	ActionUpdate
	ActionDelete
	ActionImmediate
)

func (a ActionType) MarshalJSON() ([]byte, error) {
	switch a {
	case ActionCreate: return []byte("\"create\""), nil
	case ActionUpdate: return []byte("\"update\""), nil
	case ActionDelete: return []byte("\"delete\""), nil
	case ActionImmediate: return []byte("\"immediate\""), nil
	default: return []byte("\"unknown\""), nil
	}
}

type InstanceType uint8
const (
	InstanceUnknown InstanceType = iota
	InstanceProject
	InstanceTask
	InstanceComment
	InstanceAssignee
	InstanceAttachment
	InstanceRelation
)

func (a InstanceType) MarshalJSON() ([]byte, error) {
	switch a {
	case InstanceProject: return []byte("\"project\""), nil
	case InstanceTask: return []byte("\"task\""), nil
	case InstanceComment: return []byte("\"comment\""), nil
	case InstanceAssignee: return []byte("\"assignee\""), nil
	case InstanceAttachment: return []byte("\"attachment\""), nil
	case InstanceRelation: return []byte("\"relation\""), nil
	default: return []byte("\"unknown\""), nil
	}
}

type LogEntry struct {
	Action ActionType `json:"action"`
	Instance InstanceType `json:"instance"`
	InstanceID string `json:"instance_id"`
	ProjectID int64 `json:"project_id,omitempty"`
	TaskID int64 `json:"task_id,omitempty"`
	MessageData *webhook.WebhookMessageData `json:"-"`
}

type LogFlow struct {
	Instances []LogEntry `json:"instances"`
	Events []ImmediateNode `json:"events"`
}

// Creation node

type CreationNodeType uint8
const (
	CreatedNoType CreationNodeType = iota
	CreatedTask
	CreatedAssignee
	CreatedAttachment
	CreatedComment
	CreatedRelation
)

type CreationNode struct {
	Type CreationNodeType
	ProjectID int64
	TaskID int64

	InstanceID string
	MessageData *webhook.WebhookMessageData
}

// Update node

type UpdationNodeType uint8
const (
	UpdatedNoType UpdationNodeType = iota
	UpdatedProject
	UpdatedTask
	UpdatedComment
)

type UpdationNode struct {
	Type UpdationNodeType
	ProjectID int64
	TaskID int64

	InstanceID string
	MessageData *webhook.WebhookMessageData
}

// Delete node

type DeletionNodeType uint8
const (
	DeletedNoType DeletionNodeType = iota
	DeletedProject
	DeletedAssignee
	DeletedAttachment
	DeletedTask
	DeletedComment
	DeletedRelation
)

type DeletionNode struct {
	Type DeletionNodeType
	ProjectID int64
	TaskID int64

	InstanceID string
	MessageData *webhook.WebhookMessageData
}

// Immediate node

type ImmediateNodeType uint8
const (
	ImmediateNoType ImmediateNodeType = iota
	ImmediateProjectSharedUser
	ImmediateProjectSharedTeam
	ImmediateTaskReminderFired
	ImmediateTaskOverdue
	ImmediateTasksOverdue
	ImmediateReminderFired
)

func (a ImmediateNodeType) MarshalJSON() ([]byte, error) {
	switch a {
	case ImmediateProjectSharedUser: return []byte("\"project_shared_user\""), nil
	case ImmediateProjectSharedTeam: return []byte("\"project_shared_team\""), nil
	case ImmediateTaskReminderFired: return []byte("\"task_reminder_fired\""), nil
	case ImmediateTaskOverdue: return []byte("\"task_overdue\""), nil
	case ImmediateTasksOverdue: return []byte("\"tasks_overdue\""), nil
	case ImmediateReminderFired: return []byte("\"reminder_fired\""), nil
	default: return []byte("\"unknown\""), nil
	}
}

type ImmediateNode struct {
	Type ImmediateNodeType `json:"type"`
	ProjectID int64 `json:"project_id"`

	MessageData *webhook.WebhookMessageData `json:"-"`
}
