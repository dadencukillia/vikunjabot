package webhook

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"net/http"
)

func (a *WebhookServer) verifyHeader(r *http.Request, body []byte) error {
	signature := r.Header.Get("X-Vikunja-Signature")
	if len(signature) != 32 {
		return ErrWebInvalidHash
	}

	hexSign, err := hex.DecodeString(signature)
	if err != nil {
		return errors.Join(ErrWebCouldntParseHex, err)
	}

	mac := hmac.New(sha256.New, []byte(a.webhookSecret))
	mac.Write(body)
	expectedMAC := mac.Sum(nil)
	hmac.Equal(hexSign, expectedMAC)

	return nil
}
