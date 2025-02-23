package comdirect

import (
	"fmt"
	"log/slog"
	"net/http"
)

func handleRequestError(
	res *http.Response,
) error {
	statusErr := fmt.Errorf("auth failed with status code %d", res.StatusCode)
	var body []byte
	if _, err := res.Body.Read(body); err != nil {
		return statusErr
	}

	slog.Debug(fmt.Sprintf("%v", body))
	return statusErr
}
