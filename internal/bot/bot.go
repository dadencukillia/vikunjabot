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

func (a *Bot) sendApiReq(endpoint string, reqBody []byte, responseObj any) error {
	resp, err := a.client.Post(
		wrapApiUrl(a.telegramToken, endpoint),
		"application/json", bytes.NewReader(reqBody),
	)
	if err != nil {
		return errors.Join(ErrTelegramSendingError, err)
	}


	defer resp.Body.Close()

	resBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return errors.Join(ErrTelegramInvalidBodyResponse, err)
	}

	err = json.Unmarshal(resBody, responseObj)
	if err != nil {
		return errors.Join(ErrTelegramUnexpectedResponseFormat, err)
	}

	return nil
}

func (a *Bot) GetMe() (TelegramResponse[GetMeResponse], error) {
	var ret TelegramResponse[GetMeResponse]
	err := a.sendApiReq("getMe", []byte{}, &ret)
	if err != nil {
		return TelegramResponse[GetMeResponse]{}, err
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

	var ret TelegramResponse[SendMessageResponse]
	err = a.sendApiReq("sendMessage", requestJson, &ret)
	if err != nil {
		return TelegramResponse[SendMessageResponse]{}, err
	}

	return ret, nil
}
