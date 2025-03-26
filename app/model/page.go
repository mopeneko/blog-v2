package model

import (
	"fmt"
	"net/http"
	"net/url"
	"os"
	"time"

	"github.com/gofiber/fiber/v3/client"
	"github.com/mopeneko/blog-v2/app/newt"
)

type Page struct {
	newt.BaseContent

	Title       string      `json:"title"`
	Slug        string      `json:"slug"`
	Content     string      `json:"content"`
	Thumbnail   *newt.Image `json:"thumbnail"`
	PublishedAt time.Time   `json:"published_at"`
	UpdatedAt   time.Time   `json:"updated_at"`
}

type PageClient struct {
	cc *client.Client
}

func NewPageClient() *PageClient {
	cc := client.New()

	baseURL := fmt.Sprintf("https://%s.cdn.newt.so/v1/", os.Getenv("NEWT_SPACE_UID"))
	if os.Getenv("ENV") == "development" {
		baseURL = fmt.Sprintf("https://%s.api.newt.so/v1/", os.Getenv("NEWT_SPACE_UID"))
	}

	cc.SetBaseURL(baseURL)
	cc.SetHeader("Authorization", "Bearer "+os.Getenv("NEWT_TOKEN"))

	return &PageClient{
		cc: cc,
	}
}

func (c *PageClient) FetchPages() ([]*Page, error) {
	// https://www.newt.so/docs/cdn-api-newt-api
	u, err := url.JoinPath(os.Getenv("NEWT_APP_UID"), os.Getenv("NEWT_PAGE_MODEL_UID"))
	if err != nil {
		return nil, fmt.Errorf("failed to join URL: %w", err)
	}

	resp, err := c.cc.Get(u, client.Config{
		Param: map[string]string{
			"order":  "-published_at",
			"select": "title,slug,thumbnail,tags,published_at,updated_at",
		},
	})
	if err != nil {
		return nil, fmt.Errorf("failed to send request: %w", err)
	}

	defer resp.Close()

	if resp.StatusCode() != http.StatusOK {
		return nil, fmt.Errorf("failed to get pages; status=%d", resp.StatusCode())
	}

	contents := new(newt.Contents[*Page])
	if err := resp.JSON(contents); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	return contents.Items, nil
}
func (c *PageClient) FetchPage(slug string) (*Page, error) {
	// https://www.newt.so/docs/cdn-api-newt-api
	u, err := url.JoinPath(os.Getenv("NEWT_APP_UID"), os.Getenv("NEWT_PAGE_MODEL_UID"))
	if err != nil {
		return nil, fmt.Errorf("failed to join URL: %w", err)
	}

	resp, err := c.cc.Get(u, client.Config{
		Param: map[string]string{
			"slug":  slug,
			"depth": "2",
		},
	})
	if err != nil {
		return nil, fmt.Errorf("failed to send request: %w", err)
	}

	defer resp.Close()

	if resp.StatusCode() != http.StatusOK {
		return nil, fmt.Errorf("failed to get pages; status=%d", resp.StatusCode())
	}

	page := new(newt.Contents[*Page])
	if err := resp.JSON(page); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	if len(page.Items) == 0 {
		return nil, fmt.Errorf("page not found")
	}

	return page.Items[0], nil
}
