package github

import (
	"context"
	"fmt"
	"io"
	"net/http"
)

type Client struct {
	Token string
}

func (c Client) Get(url string) ([]byte, error) {
	client := &http.Client{}

	ctx := context.Background()
	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Add("User-Agent", "GitHub Actions Usage Calculator")
	req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", c.Token))

	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("StatusCode: %d, URL: %s", resp.StatusCode, url)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	return body, nil
}
