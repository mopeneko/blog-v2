package model

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"

	"github.com/mopeneko/blog-v2/app/newt"
)

type Article struct {
	newt.BaseContent

	Title     string     `json:"title"`
	Slug      string     `json:"slug"`
	Thumbnail newt.Image `json:"thumbnail"`
	Content   string     `json:"content"`
	Tags      []Tag      `json:"tags"`
}

func FetchArticles() ([]*Article, error) {
	// https://www.newt.so/docs/cdn-api-newt-api
	base_url := fmt.Sprintf("https://%s.cdn.newt.so/v1", os.Getenv("NEWT_SPACE_UID"))
	if os.Getenv("ENV") == "development" {
		base_url = fmt.Sprintf("https://%s.api.newt.so/v1", os.Getenv("NEWT_SPACE_UID"))
	}

	path := fmt.Sprintf("/%s/%s", os.Getenv("NEWT_APP_UID"), os.Getenv("NEWT_MODEL_UID"))

	req, err := http.NewRequest("GET", base_url+path, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Authorization", "Bearer "+os.Getenv("NEWT_TOKEN"))

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to send request: %w", err)
	}

	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("failed to get articles; status=%d", resp.StatusCode)
	}

	contents := new(newt.Contents[*Article])
	if err := json.NewDecoder(resp.Body).Decode(contents); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	return contents.Items, nil
}
