package diffstree

import (
	"vikunjabot/internal/diffslog"
	"vikunjabot/internal/utils"
	"vikunjabot/internal/webhook"
)

type ChangeType uint8
const (
	Unchanged ChangeType = iota
	Created
	Updated
	Deleted
)

func ChangeTypeFromActionType(action diffslog.ActionType) ChangeType {
	switch action {
	case diffslog.ActionCreate: return Created
	case diffslog.ActionUpdate: return Updated
	case diffslog.ActionDelete: return Deleted
	default: return Unchanged
	}
}

type ChangingTopology uint8
const (
	TopSelf ChangingTopology = iota
	TopImSelfUserShared
	TopImSelfTeamShared
	TopImSelfTasksOverdue
	TopSelfTask
	TopSelfTaskComment
	TopSelfTaskAssignee
	TopSelfTaskAttachment
	TopSelfTaskRelation
	TopImSelfTaskReminder
	TopImSelfTaskOverdue
)

type RootNode struct {
	Projects map[int64]*ProjectNode
}

type ProjectNode struct {
	ID int64
	Instance *webhook.VikunjaProject
	Status ChangeType
	Tasks map[int64]*TaskNode
	ImUsersShared []ImUserSharedNode
	ImTeamsShared []ImTeamSharedNode
	ImTasksOverdue map[int64]struct{}
	DiffDoer *webhook.VikunjaUser
	Topology map[ChangingTopology]struct{} // what were changed at all
}

func (node *ProjectNode) TopCheck() *utils.ContainChecker[ChangingTopology, struct{}] {
	return utils.NewContainChecker(node.Topology)
}

type ImUserSharedNode struct {
	Instance *webhook.VikunjaUser
}

type ImTeamSharedNode struct {
	Instance *webhook.VikunjaTeam
}

type TaskNode struct {
	ID int64
	Instance *webhook.VikunjaTask
	Status ChangeType
	Comments []CommentNode
	Assignees []AssigneeNode
	Attachments []AttachmentNode
	Relations []RelationNode
	ImReminders []ImReminderNode
	ImOverdue bool
	DiffDoer *webhook.VikunjaUser
}

type CommentNode struct {
	Status ChangeType
	Instance *webhook.VikunjaTaskComment
	DiffDoer *webhook.VikunjaUser
}

type AssigneeNode struct {
	Status ChangeType
	Instance *webhook.VikunjaUser
	DiffDoer *webhook.VikunjaUser
}

type AttachmentNode struct {
	Status ChangeType
	Instance *webhook.VikunjaTaskAttachment
	DiffDoer *webhook.VikunjaUser
}

type RelationNode struct {
	Status ChangeType
	Instance *webhook.VikunjaTaskRelation
	DiffDoer *webhook.VikunjaUser
}

type ImReminderNode struct {
	Instance *webhook.VikunjaTaskReminder
}
