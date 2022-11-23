package github

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"time"
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
		retryAfter := resp.Header.Get("Retry-After")
		if retryAfter != "" {
			retryDuration, err := time.ParseDuration(retryAfter + "s")
			if err == nil {
				fmt.Printf("Retry after %v\n", retryDuration)
				time.Sleep(retryDuration)
				return c.Get(url)
			}
		}

		header := ""
		for name, values := range resp.Header {
			for _, v := range values {
				header += fmt.Sprintf("%s: %s\n", name, v)
			}
		}
		return nil, fmt.Errorf("StatusCode: %d, URL: %s\n%s", resp.StatusCode, url, header)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	return body, nil
}
