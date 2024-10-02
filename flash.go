package flash

import (
	"encoding/base64"
	"encoding/json"
	"net/http"
	"time"
)

type Severity string

const (
	SeverityNotice  Severity = "notice"
	SeverityInfo    Severity = "info"
	SeverityOk      Severity = "ok"
	SeverityWarning Severity = "warning"
	SeverityError   Severity = "error"
)

type Message struct {
	Title    string   `json:"title,omitempty"`
	Severity Severity `json:"severity,omitempty"`
	Body     string   `json:"body,omitempty"`
}

const (
	flashCookieName = "flash"
)

func Set(w http.ResponseWriter, messages ...Message) {
	data, err := json.Marshal(messages)
	if err != nil {
		return
	}
	http.SetCookie(w, &http.Cookie{
		Name:  flashCookieName,
		Value: encode(data),
	})
}

func Get(w http.ResponseWriter, r *http.Request) []Message {
	c, err := r.Cookie(flashCookieName)
	if err != nil {
		return nil
	}
	data, err := decode(c.Value)
	if err != nil {
		return nil
	}
	var messages []Message
	_ = json.Unmarshal(data, &messages)
	http.SetCookie(w, &http.Cookie{
		Name:    flashCookieName,
		MaxAge:  -1,
		Expires: time.Unix(1, 0),
	})
	return messages
}

func encode(src []byte) string {
	return base64.URLEncoding.EncodeToString(src)
}

func decode(s string) ([]byte, error) {
	return base64.URLEncoding.DecodeString(s)
}
