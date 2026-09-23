package bot

import (
	"bytes"
	"errors"
	"io"
	"net/http"

	"github.com/goccy/go-json"
)

type Bot struct {
	telegramToken string
	chatId int64
	client http.Client
}

func NewBot(telegramToken string, chatId int64) *Bot {
	return &Bot{
		telegramToken: telegramToken,
		chatId: chatId,
		client: http.Client{},
	}
}

func (a *Bot) GetMe() (TelegramResponse[GetMeResponse], error) {
	resp, err := a.client.Post(
		wrapApiUrl(a.telegramToken, "getMe"),
		"", nil,
	)
	if err != nil {
		return TelegramResponse[GetMeResponse]{}, errors.Join(ErrTelegramSendingError, err)
	}

	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return TelegramResponse[GetMeResponse]{}, errors.Join(ErrTelegramInvalidBodyResponse, err)
	}

	var ret TelegramResponse[GetMeResponse]
	err = json.Unmarshal(body, &ret)
	if err != nil {
		return TelegramResponse[GetMeResponse]{}, errors.Join(ErrTelegramUnexpectedResponseFormat, err)
	}

	return ret, nil
}

func (a *Bot) SendTextMessage(text string, parseMode ParseModeEnum, webPreview bool, notify bool) (TelegramResponse[SendMessageResponse], error) {
	request := SendMessageRequest{
		ChatID: a.chatId,
		Text: text,
		ReplyToMessageID: nil,
		DisableWebPagePreview: !webPreview,
		DisableNotification: !notify,
	}
	if parseMode != PlainTextParseMode {
		request.ParseMode = (*string)(&parseMode)
	}

	requestJson, err := json.Marshal(request)
	if err != nil {
		return TelegramResponse[SendMessageResponse]{}, errors.Join(ErrTelegramInvalidRequestObject, err)
	}

	resp, err := a.client.Post(
		wrapApiUrl(a.telegramToken, "sendMessage"),
		"application/json",
		bytes.NewReader(requestJson),
	)
	if err != nil {
		return TelegramResponse[SendMessageResponse]{}, errors.Join(ErrTelegramSendingError, err)
	}

	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return TelegramResponse[SendMessageResponse]{}, errors.Join(ErrTelegramInvalidBodyResponse, err)
	}

	var ret TelegramResponse[SendMessageResponse]
	err = json.Unmarshal(body, &ret)
	if err != nil {
		return TelegramResponse[SendMessageResponse]{}, errors.Join(ErrTelegramUnexpectedResponseFormat, err)
	}

	return ret, nil
}
