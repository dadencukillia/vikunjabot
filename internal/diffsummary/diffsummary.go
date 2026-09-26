package diffsummary

import (
	"vikunjabot/internal/diffstree"
	"vikunjabot/internal/texts"
)


type SummariesGenerator struct {
	localePack *texts.LocalePack
	vikunjaHost string
}

func NewSummariesGenerator(localePack *texts.LocalePack, vikunjaHost string) SummariesGenerator {
	return SummariesGenerator{
		localePack: localePack,
		vikunjaHost: vikunjaHost,
	}
}

func (a SummariesGenerator) GenerateHTMLSummaries(root *diffstree.RootNode) (res []string) {
	for _, project := range root.Projects {
		res = append(res, a.GenerateProjectHTMLSummary(project))
	}

	return
}

func (a SummariesGenerator) GenerateProjectHTMLSummary(project *diffstree.ProjectNode) string {
	b, _ := texts.NewTextBuilder(a.localePack)

	return b.
		SubBuild(NewProjectFormatter(project, a.vikunjaHost)).
		String()
}
