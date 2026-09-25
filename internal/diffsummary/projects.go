package diffsummary

import (
	"vikunjabot/internal/diffstree"
	"vikunjabot/internal/texts"
)

type ProjectFormatter struct {
	projectNode *diffstree.ProjectNode
}

func NewProjectFormatter(projectNode *diffstree.ProjectNode) ProjectFormatter {
	return ProjectFormatter{
		projectNode: projectNode,
	}
}

func (a ProjectFormatter) ProjectEventsPresent() bool {
	return a.projectNode.TopCheck().Any(
		diffstree.TopSelf,
		diffstree.TopImSelfTeamShared,
		diffstree.TopImSelfUserShared,
	).Fits()
}

func (a ProjectFormatter) TaskEventsPresent() bool {
	return a.projectNode.TopCheck().Any(
		diffstree.TopSelfTask,
		diffstree.TopSelfTaskRelation,
		diffstree.TopSelfTaskAttachment,
		diffstree.TopSelfTaskAssignee,
		diffstree.TopSelfTaskComment,
	).Fits()
}

func (a ProjectFormatter) RemindEventsPresent() bool {
	return a.projectNode.TopCheck().Any(
		diffstree.TopImSelfTasksOverdue,
		diffstree.TopImSelfTaskOverdue,
		diffstree.TopImSelfTaskReminder,
	).Fits()
}

func (a ProjectFormatter) GetProjectTitleVariant() string {
	if a.ProjectEventsPresent() {
		return "project"
	} else if a.TaskEventsPresent() {
		return "task"
	} else if a.RemindEventsPresent() {
		return "reminders"
	}

	return ""
}

func (a ProjectFormatter) BuildText(b texts.TextBuilder, t *texts.TextTools) string {
	return b.
		Line(true, t.Lang("TITLE", a.GetProjectTitleVariant())).
		Line(true, t.Concat(t.Lang("EMOJI", "project"), t, t.Lang("NAME", "project"), t, a.projectNode.Instance.Title)).
		EmptyLine(true).
		Line(len(a.projectNode.ImTasksOverdue) != 0, t.QuoteLine(t.Lang("TASKS_OVERDUE_COUNT", ""))).
		String()
}
