package diffstree

import "vikunjabot/internal/diffslog"

func LogFlowToDiffsTree(flow diffslog.LogFlow) RootNode {
	rootNode := RootNode{
		Projects: map[int64]*ProjectNode{},
	}

	// Projects registration
	for _, ev := range flow.Instances {
		for _, project := range ev.MessageData.GetProjects() {
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
	}

	// Tasks registration
	for _, ev := range flow.Instances {
		for _, task := range ev.MessageData.GetTasks() {
			taskNode := TaskNode{
				ID: task.ID,
				Instance: task,
				Status: Unchanged,
				Comments: []CommentNode{},
				Assignees: []AssigneeNode{},
				Attachments: []AttachmentNode{},
				Relations: []RelationNode{},
				ImReminders: []ImReminderNode{},
				ImOverdue: false,
				DiffDoer: nil,
			}

			if _, ok := rootNode.Projects[task.ProjectID]; ok {
				rootNode.Projects[task.ProjectID].Tasks[task.ID] = &taskNode
			} else {
				rootNode.Projects[task.ProjectID] = &ProjectNode{
					ID: task.ProjectID,
					Instance: nil,
					Status: Unchanged,
					Tasks: map[int64]*TaskNode{
						task.ID: &taskNode,
					},
					ImUsersShared: []ImUserSharedNode{},
					ImTeamsShared: []ImTeamSharedNode{},
					ImTasksOverdue: map[int64]struct{}{},
					DiffDoer: nil,
					Topology: map[ChangingTopology]struct{}{},
				}
			}
		}
	}
	
	// Instances
	for _, ev := range flow.Instances {
		switch ev.Instance {
		case diffslog.InstanceProject:
			project := rootNode.Projects[ev.ProjectID]
			project.DiffDoer = &ev.MessageData.Doer
			project.Status = ChangeTypeFromActionType(ev.Action)
			project.Topology[TopSelf] = struct{}{}

		case diffslog.InstanceTask:
			project := rootNode.Projects[ev.ProjectID]
			task := project.Tasks[ev.TaskID]
			task.DiffDoer = &ev.MessageData.Doer
			task.Status = ChangeTypeFromActionType(ev.Action)
			project.Topology[TopSelfTask] = struct{}{}

		case diffslog.InstanceComment:
			project := rootNode.Projects[ev.ProjectID]
			task := project.Tasks[ev.TaskID]
			task.Comments = append(task.Comments, CommentNode{
				Status: ChangeTypeFromActionType(ev.Action),
				Instance: ev.MessageData.Comment,
				DiffDoer: &ev.MessageData.Doer,
			})
			project.Topology[TopSelfTaskComment] = struct{}{}

		case diffslog.InstanceAssignee:
			project := rootNode.Projects[ev.ProjectID]
			task := project.Tasks[ev.TaskID]
			task.Assignees = append(task.Assignees, AssigneeNode{
				Status: ChangeTypeFromActionType(ev.Action),
				Instance: ev.MessageData.Assignee,
				DiffDoer: &ev.MessageData.Doer,
			})
			project.Topology[TopSelfTaskAssignee] = struct{}{}

		case diffslog.InstanceAttachment:
			project := rootNode.Projects[ev.ProjectID]
			task := project.Tasks[ev.TaskID]
			task.Attachments = append(task.Attachments, AttachmentNode{
				Status: ChangeTypeFromActionType(ev.Action),
				Instance: ev.MessageData.Attachment,
				DiffDoer: &ev.MessageData.Doer,
			})
			project.Topology[TopSelfTaskAttachment] = struct{}{}

		case diffslog.InstanceRelation:
			project := rootNode.Projects[ev.ProjectID]
			task := project.Tasks[ev.TaskID]
			task.Relations = append(task.Relations, RelationNode{
				Status: ChangeTypeFromActionType(ev.Action),
				Instance: ev.MessageData.Relation,
				DiffDoer: &ev.MessageData.Doer,
			})
			project.Topology[TopSelfTaskRelation] = struct{}{}
		}
	}

	// Events
	for _, ev := range flow.Events {
		project, ok := rootNode.Projects[ev.ProjectID]
		if !ok {
			project = &ProjectNode{
				ID: ev.ProjectID,
				Instance: nil,
				Status: Unchanged,
				Tasks: map[int64]*TaskNode{},
				ImUsersShared: []ImUserSharedNode{},
				ImTeamsShared: []ImTeamSharedNode{},
				ImTasksOverdue: map[int64]struct{}{},
				DiffDoer: nil,
				Topology: map[ChangingTopology]struct{}{},
			}
			rootNode.Projects[ev.ProjectID] = project
		}

		switch ev.Type {
		case diffslog.ImmediateProjectSharedUser:
			project.ImUsersShared = append(project.ImUsersShared, ImUserSharedNode{
				Instance: ev.MessageData.User,
			})
			project.Topology[TopImSelfUserShared] = struct{}{}

		case diffslog.ImmediateProjectSharedTeam:
			project.ImTeamsShared = append(project.ImTeamsShared, ImTeamSharedNode{
				Instance: ev.MessageData.Team,
			})
			project.Topology[TopImSelfTeamShared] = struct{}{}

		case diffslog.ImmediateTaskOverdue:
			project.ImTasksOverdue[ev.MessageData.Task.ID] = struct{}{}
			if task, ok := project.Tasks[ev.MessageData.Task.ID]; ok {
				task.ImOverdue = true
				project.Topology[TopImSelfTaskOverdue] = struct{}{}
			}
			project.Topology[TopImSelfTasksOverdue] = struct{}{}

		case diffslog.ImmediateTasksOverdue:
			for _, taskPayload := range ev.MessageData.Tasks {
				project.ImTasksOverdue[taskPayload.ID] = struct{}{}
				if task, ok := project.Tasks[taskPayload.ID]; ok {
					task.ImOverdue = true
					project.Topology[TopImSelfTaskOverdue] = struct{}{}
				}
			}
			project.Topology[TopImSelfTasksOverdue] = struct{}{}

		case diffslog.ImmediateTaskReminderFired:
			if task, ok := project.Tasks[ev.MessageData.Task.ID]; ok {
				task.ImReminders = append(task.ImReminders, ImReminderNode{
					Instance: ev.MessageData.Reminder,
				})
				project.Topology[TopImSelfTaskReminder] = struct{}{}
			}
		}
	}

	return rootNode
}
