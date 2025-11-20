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

func GetUserData(cookies []*http.Cookie, endpoint string) (user *smodel.User, cookie string, err error) {
	url := config.USER_SERVICE_URL + endpoint

	jar, err := cookiejar.New(nil)
	if err != nil {
		return nil, "", err
	}
	client := &http.Client{
		Jar:     jar,
		Timeout: 10 * time.Second,
	}

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, "", err
	}

	for _, v := range cookies {
		req.AddCookie(v)
	}

	resp, err := client.Do(req)
	if err != nil {
		return nil, "", err
	}
	defer resp.Body.Close()

	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, "", err
	}

	var body reply.ReplyEnvelope
	err = json.Unmarshal(bodyBytes, &body)
	if err != nil {
		return nil, "", err
	}

	dataBytes, err := json.Marshal(body.Data)
	if err != nil {
		return nil, "", err
	}

	var errPayload reply.ErrorPayload
	if json.Unmarshal(dataBytes, &errPayload) == nil && errPayload.Message != "" {
		return nil, "", errors.New(errPayload.Message)
	}

	var u smodel.User
	if err := json.Unmarshal(dataBytes, &u); err != nil {
		return nil, "", errors.New("bad-request: user service doesn't sent valid data")
	}

	return &u, resp.Header.Get("Set-Cookie"), nil
}
