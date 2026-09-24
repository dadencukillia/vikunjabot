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

// Cleans log and removes redundant nodes for better UX.
// An Event Compactor is a method designed to optimize log
// by compressing, merging, or filtering a stream of events
func (a *DiffsLog) Squash() {
	actions := map[string]ActionType{}
	data := map[string]*webhook.WebhookMessageData{}
	deletedProjects := map[int64]struct{}{}
	deletedTasks := map[int64]struct{}{}

	// Deciding the operation type
	// since there could be only one creation, updation or deletion node
	// for one instance (referred by instance ID)

	for _, cNode := range a.creations {
		actions[cNode.InstanceID] = ActionCreate
		data[cNode.InstanceID] = cNode.MessageData
	}

	for _, uNode := range a.updations {
		if _, ok := actions[uNode.InstanceID]; !ok {
			actions[uNode.InstanceID] = ActionUpdate
		}
		data[uNode.InstanceID] = uNode.MessageData
	}

	for _, dNode := range a.deletions {
		if _, ok := actions[dNode.InstanceID]; ok {
			delete(actions, dNode.InstanceID)
			delete(data, dNode.InstanceID)
		} else {
			actions[dNode.InstanceID] = ActionDelete
			data[dNode.InstanceID] = dNode.MessageData
		}

		switch dNode.Type {
		case DeletedTask:
			deletedTasks[dNode.TaskID] = struct{}{}
		case DeletedProject:
			deletedProjects[dNode.ProjectID] = struct{}{}
		}
	}

	// Removing extra nodes that isn't matched with our determined operation types

	{
		// Logically creation nodes are able to change their data if they were updated in the same batch
		// or they could be removed, if the instance has been removed in the outcome (or were linked to removed instance)

		newCreationNodes := []CreationNode{}

		for _, cNode := range a.creations {
			if action, ok := actions[cNode.InstanceID]; !ok || action != ActionCreate {
				continue
			}

			if _, ok := deletedProjects[cNode.ProjectID]; ok && cNode.ProjectID != -1 {
				continue
			}
			if _, ok := deletedTasks[cNode.TaskID]; ok && cNode.TaskID != -1 {
				continue
			}

			cNode.MessageData = data[cNode.InstanceID]
			newCreationNodes = append(newCreationNodes, cNode)
		}

		a.creations = newCreationNodes
	}

	{
		// The only ways update nodes could be changed is removing. Here's cases:
		// 1. Update nodes that affect on instance created in the batch are removed
		// 2. Update nodes referred on instance that is removed in outcome (or were linked to removed instance)
		// 3. Update nodes that affect on instance that has another update node are removed to leave the latest one

		newUpdationNodes := []UpdationNode{}

		for _, uNode := range a.updations {
			if action, ok := actions[uNode.InstanceID]; !ok || action != ActionUpdate {
				continue
			}

			if nodeData, ok := data[uNode.InstanceID]; !ok || nodeData != uNode.MessageData {
				continue
			}

			if _, ok := deletedProjects[uNode.ProjectID]; ok && uNode.ProjectID != -1 {
				continue
			}
			if _, ok := deletedTasks[uNode.TaskID]; ok && uNode.TaskID != -1 {
				continue
			}

			newUpdationNodes = append(newUpdationNodes, uNode)
		}

		a.updations = newUpdationNodes
	}

	{
		// The only way a deletion nodes array could be modified is
		// when there was a deletion node for an instance that was created in the same batch

		newDeletionNodes := []DeletionNode{}

		for _, dNode := range a.deletions {
			if action, ok := actions[dNode.InstanceID]; !ok || action != ActionDelete {
				continue
			}

			if _, ok := deletedProjects[dNode.ProjectID]; ok && dNode.ProjectID != -1 {
				continue
			}
			if _, ok := deletedTasks[dNode.TaskID]; ok && dNode.TaskID != -1 {
				continue
			}

			newDeletionNodes = append(newDeletionNodes, dNode)
		}

		a.deletions = newDeletionNodes
	}

	{
		// Immediate nodes could be removed if they are linked
		// to project or task that is removed in outcome

		newImmediateNodes := []ImmediateNode{}
		
		for _, iNode := range a.immediates {
			switch iNode.Type {
			// Task and project checking
			case ImmediateTaskOverdue, ImmediateTaskReminderFired, ImmediateTasksOverdue:
				taskId := iNode.MessageData.GetTaskID()
				if _, ok := deletedTasks[taskId]; ok && taskId != -1 {
					continue
				}
				if _, ok := deletedProjects[iNode.ProjectID]; ok && iNode.ProjectID != -1 {
					continue
				}

			// Project checking
			case ImmediateProjectSharedUser, ImmediateProjectSharedTeam:
				if _, ok := deletedProjects[iNode.ProjectID]; ok && iNode.ProjectID != -1 {
					continue
				}
			}

			newImmediateNodes = append(newImmediateNodes, iNode)
		}

		a.immediates = newImmediateNodes
	}
}
