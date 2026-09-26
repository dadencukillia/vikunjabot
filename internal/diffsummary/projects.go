package diffsummary

import (
	"fmt"
	"strings"
	"vikunjabot/internal/diffstree"
	"vikunjabot/internal/texts"
	"vikunjabot/internal/utils"
)

type ProjectFormatter struct {
	projectNode *diffstree.ProjectNode
	vikunjaHost string
	timeZone string
}

func NewProjectFormatter(projectNode *diffstree.ProjectNode, vikunjaHost string, timeZone string) ProjectFormatter {
	return ProjectFormatter{
		projectNode: projectNode,
		vikunjaHost: vikunjaHost,
		timeZone: timeZone,
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

func (a ProjectFormatter) GetDescriptionPureLine() string {
	return strings.TrimSpace(strings.ReplaceAll(
		utils.ExtractHTMLText(a.projectNode.Instance.Description),
		"\n", " ",
	))
}

func (a ProjectFormatter) GetTeamSharedName() string {
	if len(a.projectNode.ImTeamsShared) == 0 {
		return ""
	}

	return a.projectNode.ImTeamsShared[0].Instance.Name
}

func (a ProjectFormatter) GetUserSharedUsername() string {
	if len(a.projectNode.ImUsersShared) == 0 {
		return ""
	}

	return a.projectNode.ImUsersShared[0].Instance.Username
}

func (a ProjectFormatter) GetTaskID() int64 {
	for _, task := range a.projectNode.Tasks {
		return task.ID
	}

	return -1
}

func (a ProjectFormatter) BuildText(b texts.TextBuilder, t *texts.TextTools) string {
	if a.projectNode.Instance == nil {
		return ""
	}

	description := utils.CutString(a.GetDescriptionPureLine(), 100, true)

	return b.
		Line(true, t.Lang("TITLE", a.GetProjectTitleVariant())).
		EmptyLine(true).
		Line(true, t.Quote(t.Concat(
			t.Bold(t.Concat(
				t.Cond(a.projectNode.Status == diffstree.Updated,
					t.Lang("EMOJI", "updated"), t,
					t.Lang("STATUS_PROJECT", "updated"), t,
				),
				t.Cond(a.projectNode.Status == diffstree.Deleted,
					t.Lang("EMOJI", "deleted"), t,
					t.Lang("STATUS_PROJECT", "deleted"), t,
				),
				t.Cond(a.projectNode.Status == diffstree.Unchanged,
					t.Lang("EMOJI", "project"), t,
					t.Lang("SINGULAR", "project"), t, 
				),
			)),
			t.Anchor(
				t.EscHTML(a.projectNode.Instance.Title),
				utils.UrlConcat(a.vikunjaHost, "/projects/", a.projectNode.ID),
			),
		))).
		Line(len(a.projectNode.ImTasksOverdue) != 0, t.LangMap("TASKS_OVERDUE_COUNT", "", map[string]string{
			"COUNT": fmt.Sprint(len(a.projectNode.ImTasksOverdue)),
		})).
		Line(len(a.projectNode.ImTeamsShared) == 1, t.LangMap("PROJECT_SHARED_TEAM", "", map[string]string{
			"NAME": t.EscHTML(a.GetTeamSharedName()),
		})).
		Line(len(a.projectNode.ImTeamsShared) > 1, t.LangMap("PROJECT_SHARED_TEAMS", "", map[string]string{
			"COUNT": fmt.Sprint(len(a.projectNode.ImTeamsShared)),
		})).
		Line(len(a.projectNode.ImUsersShared) == 1, t.LangMap("PROJECT_SHARED_USER", "", map[string]string{
			"NAME": t.EscHTML(a.GetUserSharedUsername()),
		})).
		Line(len(a.projectNode.ImUsersShared) > 1, t.LangMap("PROJECT_SHARED_USERS", "", map[string]string{
			"COUNT": fmt.Sprint(len(a.projectNode.ImUsersShared)),
		})).
		Line(a.projectNode.Status != diffstree.Unchanged && a.projectNode.Status != diffstree.Deleted && len(description) != 0, t.Code(description)).
		EmptyLine(true).
		BuildScope(func(b texts.TextBuilder, t *texts.TextTools) string {
			if len(a.projectNode.Tasks) == 0 {
				return ""
			}

			b.Line(true, t.Lang("EMOJI", "task"), t, t.Bold(t.Lang("PLURAL", "task")))

			for _, task := range a.projectNode.Tasks {
				b.SubBuild(NewTaskFormatter(task, a.vikunjaHost, a.timeZone))
			}

			return b.String()
		}).
		EmptyLine(true).
		Line(true,
			"🌐 ",
			t.Anchor(t.Lang("ANCHOR", "website"), a.vikunjaHost),
			" | ", t.Anchor(t.Lang("ANCHOR", "project"), utils.UrlConcat(a.vikunjaHost, "/projects/", a.projectNode.ID)),
			t.Cond(len(a.projectNode.Tasks) > 0, " | ", t.Anchor(t.Lang("ANCHOR", "task"), utils.UrlConcat(a.vikunjaHost, "/tasks/", a.GetTaskID()))),
		).
		String()
}
