package github

import (
	"net/http"

	"github.com/google/go-github/v71/github"
)

func Ptr[T any](v T) *T {
	return github.Ptr(v)
}

func ParseWebHook(messageType string, payload []byte) (interface{}, error) {
	return github.ParseWebHook(messageType, payload)
}

func WebHookType(r *http.Request) string {
	return github.WebHookType(r)
}

func ValidatePayload(r *http.Request, secretToken []byte) (payload []byte, err error) {
	return github.ValidatePayload(r, secretToken)
}
