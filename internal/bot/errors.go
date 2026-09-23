package bot

import "errors"

var ErrTelegramSendingError = errors.New("couldn't send an http request")
var ErrTelegramInvalidBodyResponse = errors.New("invalid body received from telegram api")
var ErrTelegramUnexpectedResponseFormat = errors.New("unexpected telegram response format")
var ErrTelegramInvalidRequestObject = errors.New("invalid request object")
