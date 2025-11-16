package userlib

import (
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"product-service/config"

	"github.com/chesta132/e-commerce-go/shared"
	"github.com/chesta132/goreply/reply"
)

func GetUserDataWithAuth() (*shared.User, error) {
	resp, err := http.Get(config.USER_SERVICE_URL + "/auth/admin")
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var body reply.ReplyEnvelope
	err = json.Unmarshal(bodyBytes, &body)
	if err != nil {
		return nil, err
	}

	if err, ok := body.Data.(reply.ErrorPayload); ok {
		return nil, errors.New(err.Message)
	}

	user, ok := body.Data.(*shared.User)
	if !ok {
		return nil, errors.New("bad-request: user service doesn't sent a valid data")
	}
	return user, nil
}
