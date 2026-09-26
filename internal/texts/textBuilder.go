package texts

import (
	"fmt"
	"html"
	"strconv"
	"strings"
)

var LANG_KEY_NOT_SET = "key not set"

type TextFragmentBuilder interface {
	BuildText(b TextBuilder, t *TextTools) string
}

type TextBuilder struct {
	localePack *LocalePack
	textTools *TextTools
	currentState strings.Builder
}

func NewTextBuilder(locale *LocalePack) (b TextBuilder, t *TextTools) {
	tools := &TextTools{
		localePack: locale,
	}

	return TextBuilder{
		localePack: locale,
		textTools: tools,
		currentState: strings.Builder{},
	}, tools
}

func (a *TextBuilder) BuildScope(scope func(b TextBuilder, t *TextTools) string) *TextBuilder {
	builder := TextBuilder{
		localePack: a.localePack,
		textTools: a.textTools,
		currentState: strings.Builder{},
	}
	scopeResult := scope(builder, a.textTools)
	a.currentState.WriteString(scopeResult)

	return a
}

func (a *TextBuilder) Line(cond bool, parts... any) *TextBuilder {
	if (cond) {
		fmt.Fprint(&a.currentState, parts...)
		a.currentState.WriteByte('\n')
	}

	return a
}

func (a *TextBuilder) DelimLine(cond bool) *TextBuilder {
	if (cond) {
		a.currentState.WriteString("- - -\n")
	}

	return a
}

func (a *TextBuilder) EmptyLine(cond bool) *TextBuilder {
	if (cond) {
		a.currentState.WriteString("\n")
	}

	return a
}

func (a *TextBuilder) SubBuild(obj TextFragmentBuilder) *TextBuilder {
	return a.BuildScope(obj.BuildText)
}

func (a *TextBuilder) String() string {
	return a.currentState.String()
}

// Tools

type TextTools struct {
	localePack *LocalePack
}

func (a TextTools) String() string {
	return " "
}

func (a TextTools) Lang(key string, variant string) string {
	return a.localePack.GetOrDefault(key, variant, LANG_KEY_NOT_SET)
}

func (a TextTools) LangMap(key string, variant string, dict map[string]string) string {
	return a.localePack.GetWithDictOrDefault(key, variant, func(key string) (string, bool) {
		v, ok := dict[key]
		return v, ok
	}, LANG_KEY_NOT_SET)
}

func (a TextTools) ExtLang() *LocalePack {
	return a.localePack
}

func (a TextTools) EscHTML(text string) string {
	return html.EscapeString(text)
}

func (a TextTools) Cond(cond bool, parts... any) string {
	if cond {
		return fmt.Sprint(parts...)
	}

	return ""
}

func (a TextTools) Concat(parts... any) string {
	return fmt.Sprint(parts...)
}

func (a TextTools) Code(text string) string {
	return "<code>" + text + "</code>"
}

func (a TextTools) MonoBlock(text string, title string) string {
	return "<pre language=" + strconv.Quote(title) + ">" + text + "</pre>"
}

func (a TextTools) Italic(text string) string {
	return "<i>" + text + "</i>"
}

func (a TextTools) Bold(text string) string {
	return "<b>" + text + "</b>"
}

func (a TextTools) Strike(text string) string {
	return "<s>" + text + "</s>"
}

func (a TextTools) Underline(text string) string {
	return "<u>" + text + "</u>"
}

func (a TextTools) Quote(text string) string {
	return "<blockquote>" + text + "</blockquote>"
}

func (a TextTools) QuotePage(text string) string {
	return "<blockquote expandable>" + text + "</blockquote>"
}

func (a TextTools) Anchor(text string, url string) string {
	return "<a href=" + strconv.Quote(url) + ">" + text + "</a>"
}

