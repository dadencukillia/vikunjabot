package webhook

import "time"

// Vikunja primitives

type VikunjaRelationKind string
const (
	RelationKindUnknown VikunjaRelationKind = `unknown`
	RelationKindSubtask VikunjaRelationKind = `subtask`
	RelationKindParenttask VikunjaRelationKind = `parenttask`
	RelationKindRelated VikunjaRelationKind = `related`
	RelationKindDuplicateOf VikunjaRelationKind = `duplicateof`
	RelationKindDuplicates VikunjaRelationKind = `duplicates`
	RelationKindBlocking VikunjaRelationKind = `blocking`
	RelationKindBlocked VikunjaRelationKind = `blocked`
	RelationKindPreceeds VikunjaRelationKind = `precedes`
	RelationKindFollows VikunjaRelationKind = `follows`
	RelationKindCopiedFrom VikunjaRelationKind = `copiedfrom`
	RelationKindCopiedTo VikunjaRelationKind = `copiedto`
)

type VikunjaReminderRelation string
const (
	ReminderRelationDueDate VikunjaReminderRelation = `due_date`
	ReminderRelationStartDate VikunjaReminderRelation = `start_date`
	ReminderRelationEndDate VikunjaReminderRelation = `end_date`
)

type VikunjaUser struct {
	ID int64 `json:"id"`
	Name string `json:"name"`
	Username string `json:"username"`
	Email *string `json:"email,omitempty"`
	BotOwnerID *int64  `json:"bot_owner_id,omitempty"`
	Created time.Time `json:"created"`
	Updated time.Time `json:"updated"`
}

type VikunjaProject struct {
	ID int64 `json:"id"`
	Title string `json:"title"`
	Description string `json:"description"`
	Identifier string `json:"identifier"`
	HexColor string `json:"hex_color"`
	ParentProjectID int64 `json:"parent_project_id"`
	Owner VikunjaUser `json:"owner"`
	IsArchived bool `json:"is_archived"`
	Position float64 `json:"position"`
	Created time.Time `json:"created"`
	Updated time.Time `json:"updated"`
}

type VikunjaTask struct {
	Title string `json:"title"`
	Id int `json:"id"`
	Description string `json:"description"`
	DueDate string `json:"due_date"`
	Priority int `json:"priority"`
	Identifier string `json:"identifier"`
	ProjectID int `json:"project_id"`
}

type VikunjaTaskComment struct {
	ID int64 `json:"id"`
	Comment string `json:"comment"`
	Author VikunjaUser `json:"author"`
	Created time.Time `json:"created"`
	Updated time.Time `json:"updated"`
}

type VikunjaTaskAttachment struct {
	ID int64 `json:"id"`
	TaskID int64 `json:"task_id"`
	CreatedBy VikunjaUser `json:"created_by"`
	File VikunjaFile `json:"file"`
	Created time.Time `json:"created"`
}

type VikunjaFile struct {
	ID int64 `json:"id"`
	Name string `json:"name"`
	Mime string `json:"mime"`
	Size uint64 `json:"size"`
	Created time.Time `json:"created"`
}

type VikunjaTaskRelation struct {
	TaskID int64 `json:"task_id"`
	OtherTaskID int64 `json:"other_task_id"`
	RelationKind VikunjaRelationKind `json:"relation_kind"`
	CreatedBy VikunjaUser `json:"created_by"`
	Created time.Time `json:"created"`
}

type VikunjaTaskReminder struct {
	Reminder time.Time `json:"reminder"`
	RelativePeriod int64 `json:"relative_period"`
	RelativeTo VikunjaReminderRelation `json:"relative_to"`
}

type VikunjaTeam struct {
	ID int64 `json:"id"`
	Name string `json:"name"`
	Description string `json:"description"`
	ExternalID string `json:"external_id"`
	CreatedBy VikunjaUser `json:"created_by"`
	Members []VikunjaUser `json:"members"`
	Created time.Time `json:"created"`
	Updated time.Time `json:"updated"`
	IsPublic bool `json:"is_public"`
}

// Webhook

type WebhookMessage struct {
	EventName string `json:"event_name"`
	Time time.Time `json:"time"`
	Data WebhookMessageData `json:"data"`
}

type WebhookMessageData struct {
	Doer VikunjaUser `json:"doer"`
	Task *VikunjaTask `json:"task,omitempty"`
	Tasks []*VikunjaTask `json:"tasks,omitempty"`
	Assignee *VikunjaUser `json:"assignee,omitempty"`
	Comment *VikunjaTaskComment `json:"comment,omitempty"`
	Attachment *VikunjaTaskAttachment `json:"attachment,omitempty"`
	Relation *VikunjaTaskRelation `json:"relation,omitempty"`
	Reminder *VikunjaTaskReminder `json:"reminder,omitempty"`
	Project *VikunjaProject `json:"project,omitempty"`
	Projects map[int64]*VikunjaProject `json:"projects,omitempty"`
	User *VikunjaUser `json:"user,omitempty"`
	Team *VikunjaTeam `json:"team"`
}
