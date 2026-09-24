package diffsummary

import (
	"strings"
	"vikunjabot/internal/diffstree"
	"vikunjabot/internal/utils"
)

var LANG_KEY_NOT_SET = "key not set"

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

	engine := NewReplacementsEngine()
	engine.RegisterDirect("SUMMARY_TITLE", a.localePack.GetOrDefault("TITLE", titleVariant, LANG_KEY_NOT_SET))
	engine.RegisterHandler("LANG", func(key string) (string, bool) {
		langKey := strings.SplitN(key, ":", 2)
		if len(langKey) != 2 {
			return "", false
		}

		return a.localePack.GetOrDefault(langKey[0], langKey[1], LANG_KEY_NOT_SET), true
	})

	return engine.ProcessString(`
^{SUMMARY_TITLE}

^{LANG:EMOJI:PROJECT} ^{LANG:NAME:PROJECT}
	`)
}
