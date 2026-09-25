package diffsummary

import (
	"strings"
	"vikunjabot/internal/diffstree"
	"vikunjabot/internal/utils"
)


type SummariesGenerator struct {
	localePack *LocalePack
}

func NewSummariesGenerator(localePack *LocalePack) SummariesGenerator {
	return SummariesGenerator{
		localePack: localePack,
	}
}

func (a SummariesGenerator) GenerateHTMLSummaries(root *diffstree.RootNode) (res []string) {
	for _, project := range root.Projects {
		res = append(res, a.GenerateProjectHTMLSummary(project))
	}

	return
}

func (a SummariesGenerator) GenerateProjectHTMLSummary(project *diffstree.ProjectNode) string {
	titleVariant := ""

	if utils.NewContainChecker(project.Topology).Any(
		diffstree.TopSelf,
		diffstree.TopImSelfTeamShared,
		diffstree.TopImSelfUserShared,
	).Fits() {
		titleVariant = "project"

	} else if utils.NewContainChecker(project.Topology).Any(
		diffstree.TopSelfTask,
		diffstree.TopSelfTaskRelation,
		diffstree.TopSelfTaskAttachment,
		diffstree.TopSelfTaskAssignee,
		diffstree.TopSelfTaskComment,
	).Fits() {
		titleVariant = "task"

	} else if utils.NewContainChecker(project.Topology).Any(
		diffstree.TopImSelfTasksOverdue,
		diffstree.TopImSelfTaskOverdue,
		diffstree.TopImSelfTaskReminder,
	).Fits() {
		titleVariant = "reminders"

	}

	b, t := NewTextBuilder(a.localePack)

	return b.
		Line(true, t.Lang("TITLE", titleVariant)).
		EmptyLine(true).
		String()
}
