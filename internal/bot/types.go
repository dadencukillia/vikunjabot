package bot

// Primitives

type ParseModeEnum string
const (
	HTMLParseMode ParseModeEnum = "HTML"
	MarkdownParseMode ParseModeEnum = "MarkdownV2"
	PlainTextParseMode ParseModeEnum = ""
)

type TelegramUser struct {
	ID int64 `json:"id"`
	FirstName string `json:"first_name"`
	Username string `json:"username"`
}

type TelegramChat struct {
	ID int64 `json:"id"`
	FirstName string `json:"first_name"`
	LastName string `json:"last_name"`
	Username string `json:"username"`
	Type string `json:"type"`
}

// Response template

type TelegramResponse[T any] struct {
	Ok bool `json:"ok"`
	Result *T `json:"result,omitempty"`
	ErrorCode int64 `json:"error_code,omitempty"`
	ErrorDescription string `json:"description,omitempty"`
}

// Request/Response models

type GetMeResponse TelegramUser

type SendMessageRequest struct {
	ChatID int64 `json:"chat_id"`
	Text string `json:"text"`
	ParseMode *string `json:"parse_mode,omitempty"`
	DisableWebPagePreview bool `json:"disable_web_page_preview"`
	DisableNotification bool `json:"disable_notification"`
	ReplyToMessageID *int64 `json:"reply_to_message_id,omitempty"`
}

type SendMessageResponse struct {
	MessageID int64 `json:"message_id"`
	Date string `json:"date"`
	Text string `json:"text"`
	From TelegramUser `json:"from"`
	Chat TelegramChat `json:"chat"`
}
