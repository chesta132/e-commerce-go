package userlib

import (
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"net/http/cookiejar"
	"product-service/config"
	"time"

	"github.com/chesta132/e-commerce-go/shared/smodel"
	"github.com/chesta132/goreply/reply"
)

type GetUserError struct {
	error
	Status int
}

func GetUserData(cookies []*http.Cookie, endpoint string) (user *smodel.User, cookie string, err error) {
	url := config.USER_SERVICE_URL + endpoint

	jar, err := cookiejar.New(nil)
	if err != nil {
		return nil, "", GetUserError{error: err, Status: 500}
	}
	client := &http.Client{
		Jar:     jar,
		Timeout: 10 * time.Second,
	}

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, "", GetUserError{error: err, Status: 500}
	}

	for _, v := range cookies {
		req.AddCookie(v)
	}

	resp, err := client.Do(req)
	if err != nil {
		return nil, "", GetUserError{error: err, Status: 502}
	}
	defer resp.Body.Close()

	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, "", GetUserError{error: err, Status: 502}
	}

	var body reply.ReplyEnvelope
	err = json.Unmarshal(bodyBytes, &body)
	if err != nil {
		return nil, "", GetUserError{error: err, Status: 502}
	}

	dataBytes, err := json.Marshal(body.Data)
	if err != nil {
		return nil, "", GetUserError{error: err, Status: 502}
	}

	var errPayload reply.ErrorPayload
	if json.Unmarshal(dataBytes, &errPayload) == nil && errPayload.Message != "" {
		return nil, "", GetUserError{error: errors.New(errPayload.Message), Status: resp.StatusCode}
	}

	var u smodel.User
	if err := json.Unmarshal(dataBytes, &u); err != nil {
		return nil, "", GetUserError{error: errors.New("bad-request: user service doesn't sent valid data"), Status: 502}
	}

	return &u, resp.Header.Get("Set-Cookie"), nil
}
