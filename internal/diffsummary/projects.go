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

func (a ProjectFormatter) GetProjectTitleVariant() string {
	if a.projectNode.TopCheck().Any(
		diffstree.TopSelf,
		diffstree.TopImSelfTeamShared,
		diffstree.TopImSelfUserShared,
	).Fits() {
		return "project"

	} else if a.projectNode.TopCheck().Any(
		diffstree.TopSelfTask,
		diffstree.TopSelfTaskRelation,
		diffstree.TopSelfTaskAttachment,
		diffstree.TopSelfTaskAssignee,
		diffstree.TopSelfTaskComment,
	).Fits() {
		return "task"

	} else if a.projectNode.TopCheck().Any(
		diffstree.TopImSelfTasksOverdue,
		diffstree.TopImSelfTaskOverdue,
		diffstree.TopImSelfTaskReminder,
	).Fits() {
		return "reminders"
	}

	return ""
}

func (a ProjectFormatter) BuildText(b texts.TextBuilder, t *texts.TextTools) string {
	return b.
		Line(true, t.Lang("TITLE", a.GetProjectTitleVariant())).
		String()
}
