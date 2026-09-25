package diffsummary

import (
	"vikunjabot/internal/diffstree"
	"vikunjabot/internal/texts"
)


type SummariesGenerator struct {
	localePack *texts.LocalePack
}

func NewSummariesGenerator(localePack *texts.LocalePack) SummariesGenerator {
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
	b, _ := texts.NewTextBuilder(a.localePack)

	return b.
		SubBuild(NewProjectFormatter(project)).
		String()
}
