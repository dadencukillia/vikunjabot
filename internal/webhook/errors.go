package webhook

import "errors"

var ErrWebInvalidHash = errors.New("invalid hash format")
var ErrWebCouldntParseHex = errors.New("couldn't parse hex")
var ErrWebInvalidEventRequest = errors.New("couldn't parse event body")
var ErrWebNoEventProject = errors.New("each event must be related to at least one project")
