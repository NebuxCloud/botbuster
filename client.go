package botbuster

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
)

// Client interface
type Client interface {
	Verify(string) (bool, error)
}

// Client implementation
type botbusterClient struct {
	url    string
	client *http.Client
}

type response struct {
	Success bool `json:"success"`
}

func NewClient(url string) Client {
	return &botbusterClient{
		url:    url,
		client: &http.Client{},
	}
}

func (c *botbusterClient) Verify(solution string) (bool, error) {
	url := fmt.Sprintf("%s/v1/verify", c.url)

	req, err := http.NewRequest(
		http.MethodPost,
		url,
		strings.NewReader(solution),
	)

	if err != nil {
		return false, err
	}

	req.Header.Set("Content-Type", "text/plain; charset=utf-8")

	res, err := c.client.Do(req)
	if err != nil {
		return false, err
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK &&
		res.StatusCode != http.StatusForbidden {
		return false, nil
	}

	var vr response
	if err := json.NewDecoder(res.Body).Decode(&vr); err != nil {
		return false, err
	}

	return vr.Success, nil
}

// Dummy client implementation
type DummyClient struct {
}

func (c *DummyClient) Verify(token string) (bool, error) {
	return true, nil
}
