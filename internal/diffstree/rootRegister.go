package diffstree

import (
	"vikunjabot/internal/diffslog"
	"vikunjabot/internal/webhook"
)

func rootRegister(rootNode *RootNode, flow *diffslog.LogFlow) {
	collectedProjects := make([]*webhook.VikunjaProject, 0, len(flow.Instances) + len(flow.Events))
	collectedTasksInstances := make([]*webhook.VikunjaTask, 0, len(flow.Instances)) // Higher priority
	collectedTasksEvents := make([]*webhook.VikunjaTask, 0, len(flow.Events)) // Lower priority

	// Collecting projects and tasks from instances

	for _, ev := range flow.Instances {
		for _, project := range ev.MessageData.GetProjects() {
			collectedProjects = append(collectedProjects, project)
		}
		for _, task := range ev.MessageData.GetTasks() {
			collectedTasksInstances = append(collectedTasksInstances, task)
		}
	}

	// Collecting projects and tasks from events

	for _, ev := range flow.Events {
		for _, project := range ev.MessageData.GetProjects() {
			collectedProjects = append(collectedProjects, project)
		}
		for _, task := range ev.MessageData.GetTasks() {
			collectedTasksEvents = append(collectedTasksEvents, task)
		}
	}

	// Adding projects into our tree

	for _, project := range collectedProjects {
		if p, ok := rootNode.Projects[project.ID]; !ok || (p.Instance == nil && project != nil) {
			rootNode.Projects[project.ID] = &ProjectNode{
				ID: project.ID,
				Instance: project,
				Status: Unchanged,
				Tasks: map[int64]*TaskNode{},
				ImUsersShared: []ImUserSharedNode{},
				ImTeamsShared: []ImTeamSharedNode{},
				ImTasksOverdue: map[int64]struct{}{},
				DiffDoer: nil,
				Topology: map[ChangingTopology]struct{}{},
			}
		}
	}

	// Adding tasks into our tree

	taskNodeTemplate := TaskNode{
		// Requires to fill ID and Instance fields
		Status: Unchanged,
		Comments: []CommentNode{},
		Assignees: []AssigneeNode{},
		Attachments: []AttachmentNode{},
		Relations: []RelationNode{},
		ImReminders: []ImReminderNode{},
		ImOverdue: false,
		DiffDoer: nil,
	}

	for _, taskList := range [][]*webhook.VikunjaTask{
		collectedTasksEvents,
		collectedTasksInstances,
	} {
		for _, task := range taskList {
			if project, ok := rootNode.Projects[task.ProjectID]; ok {
				taskNodeCopy := taskNodeTemplate
				taskNodeCopy.ID = task.ID
				taskNodeCopy.Instance = task
				project.Tasks[task.ID] = &taskNodeCopy
			}
		}
	}
}
