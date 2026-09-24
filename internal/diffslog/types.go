package diffslog

import "vikunjabot/internal/webhook"

type ActionType uint8
const (
	ActionUnknown ActionType = iota
	ActionCreate
	ActionUpdate
	ActionDelete
	ActionImmediate
)

// Creation node

type CreationNodeType uint8
const (
	CreatedNoType CreationNodeType = iota
	CreatedAssignee
	CreatedTask
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

type ImmediateNode struct {
	Type ImmediateNodeType
	ProjectID int64

	MessageData *webhook.WebhookMessageData
}
