package diffstree

import (
	"vikunjabot/internal/diffslog"
)

func LogFlowToDiffsTree(flow diffslog.LogFlow) RootNode {
	rootNode := RootNode{
		Projects: map[int64]*ProjectNode{},
	}

	rootRegister(&rootNode, &flow)
	
	// Instances

	for _, ev := range flow.Instances {
		project := rootNode.Projects[ev.ProjectID]

		switch ev.Instance {
		case diffslog.InstanceProject:
			project.DiffDoer = &ev.MessageData.Doer
			project.Status = ChangeTypeFromActionType(ev.Action)
			project.Topology[TopSelf] = struct{}{}

		case diffslog.InstanceTask:
			task := project.Tasks[ev.TaskID]
			task.DiffDoer = &ev.MessageData.Doer
			task.Status = ChangeTypeFromActionType(ev.Action)
			project.Topology[TopSelfTask] = struct{}{}

		case diffslog.InstanceComment:
			task := project.Tasks[ev.TaskID]
			task.Comments = append(task.Comments, CommentNode{
				Status: ChangeTypeFromActionType(ev.Action),
				Instance: ev.MessageData.Comment,
				DiffDoer: &ev.MessageData.Doer,
			})
			project.Topology[TopSelfTaskComment] = struct{}{}

		case diffslog.InstanceAssignee:
			task := project.Tasks[ev.TaskID]
			task.Assignees = append(task.Assignees, AssigneeNode{
				Status: ChangeTypeFromActionType(ev.Action),
				Instance: ev.MessageData.Assignee,
				DiffDoer: &ev.MessageData.Doer,
			})
			project.Topology[TopSelfTaskAssignee] = struct{}{}

		case diffslog.InstanceAttachment:
			task := project.Tasks[ev.TaskID]
			task.Attachments = append(task.Attachments, AttachmentNode{
				Status: ChangeTypeFromActionType(ev.Action),
				Instance: ev.MessageData.Attachment,
				DiffDoer: &ev.MessageData.Doer,
			})
			project.Topology[TopSelfTaskAttachment] = struct{}{}

		case diffslog.InstanceRelation:
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
		project := rootNode.Projects[ev.ProjectID]

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
