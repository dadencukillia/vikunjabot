package diffsummary

import (
	"strings"
	"vikunjabot/internal/diffstree"
	"vikunjabot/internal/texts"
	"vikunjabot/internal/utils"
)

type TaskFormatter struct {
	taskNode *diffstree.TaskNode
	vikunjaHost string
}

func NewTaskFormatter(taskNode *diffstree.TaskNode, vikunjaHost string) TaskFormatter {
	return TaskFormatter{
		taskNode: taskNode,
		vikunjaHost: vikunjaHost,
	}
}

func (a TaskFormatter) GetDescriptionPureLine() string {
	return strings.TrimSpace(strings.ReplaceAll(
		utils.ExtractHTMLText(a.taskNode.Instance.Description),
		"\n", " ",
	))
}

func (a TaskFormatter) BuildText(b texts.TextBuilder, t *texts.TextTools) string {
	description := utils.CutString(a.GetDescriptionPureLine(), 200, true)
	statusBody := a.BuildStatusBody(b, t)

	return b.
		Line(true, t.Quote(t.Concat(
			t.Concat(
				t.Cond(a.taskNode.Status == diffstree.Created, t.Lang("EMOJI", "created"), t),
				t.Cond(a.taskNode.Status == diffstree.Updated, t.Lang("EMOJI", "updated"), t),
				t.Cond(a.taskNode.Status == diffstree.Deleted, t.Lang("EMOJI", "deleted"), t),
				t.Anchor(t.EscHTML(a.taskNode.Instance.Title), utils.UrlConcat(a.vikunjaHost, "/tasks/", a.taskNode.ID)),
			),
			t.Cond(a.taskNode.Status != diffstree.Deleted, t.Concat(
				t.Cond(len(description) != 0 && a.taskNode.Status != diffstree.Deleted && a.taskNode.Status != diffstree.Unchanged, "\n", t.Code(t.EscHTML(description))),
				t.Cond(len(statusBody) != 0, "\n\n", strings.TrimRight(statusBody, "\n")),
			)),
		))).
		String()
}

func (a TaskFormatter) BuildStatusBody(b texts.TextBuilder, t *texts.TextTools) string {
	b.Line(a.taskNode.ImOverdue, t.Bold(t.Lang("TASK_OVERDUE", "")))

	for _, attachment := range a.taskNode.Attachments {
		b.Line(true, 
			t.Lang("EMOJI", "file"),
			t.Cond(attachment.Status == diffstree.Deleted, t.Lang("EMOJI", "deleted")), 
			t, t.Bold(attachment.Instance.File.Name),
		)
	}

	for _, comment := range a.taskNode.Comments {
		commentText := utils.CutString(utils.ExtractHTMLText(comment.Instance.Comment), 60, true)
		b.Line(true,
			t.Lang("EMOJI", "comment"),
			t.Cond(comment.Status == diffstree.Deleted, t.Lang("EMOJI", "deleted")),
			t, t.Bold(comment.Instance.Author.Username),
			t.Cond(comment.Status != diffstree.Deleted, ": ", t.EscHTML(commentText)),
		)
	}

	for _, assign := range a.taskNode.Assignees {
		b.Line(assign.Status == diffstree.Created, "+ ", t.Bold(assign.Instance.Username))
		b.Line(assign.Status == diffstree.Deleted, "- ", t.Bold(assign.Instance.Username))
	}

	addRelations := []string{}
	minusRelations := []string{}
	for _, relation := range a.taskNode.Relations {
		switch relation.Status {
		case diffstree.Created:
			addRelations = append(addRelations, t.Lang("TASK_RELATION", string(relation.Instance.RelationKind)))
		case diffstree.Deleted:
			minusRelations = append(minusRelations, t.Lang("TASK_RELATION", string(relation.Instance.RelationKind)))
		}
	}
	b.Line(len(addRelations) > 0, t.Lang("EMOJI", "created"), t, t.Bold(t.Lang("TASK_RELATIONS", "") + ": "), strings.Join(addRelations, ", "))
	b.Line(len(minusRelations) > 0, t.Lang("EMOJI", "deleted"), t, t.Bold(t.Lang("TASK_RELATIONS", "") + ": "), strings.Join(minusRelations, ", "))

	return b.String()
}
