package basiq

import (
	"errors"
	"fmt"
	"net/http"
)

type Error struct {
	HttpCode      int
	Type          string `json:"type"`
	CorrelationId string `json:"correlationId"`
	Data          []struct {
		Code   string `json:"code"`
		Detail string `json:"detail"`
		Source struct {
			Parameter string `json:"parameter"`
			Pointer   string `json:"pointer"`
		} `json:"source"`
		Title string `json:"title"`
		Type  string `json:"type"`
	} `json:"data"`
}

// --------------------------------------------------------------------------------------------------------------------

func (e *Error) Error() string {
	if len(e.Data) == 0 {
		return "Unknown error"
	}
	return fmt.Sprintf("%d: %s:%s", e.HttpCode, e.Data[0].Title, e.Data[0].Detail)
}

func IsUnauthorizedErr(err error) bool {
	var e *Error
	if !errors.As(err, &e) {
		return false
	}

	return e.HttpCode == http.StatusUnauthorized
}
