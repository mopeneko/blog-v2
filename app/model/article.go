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

type Article struct {
	newt.BaseContent

	Title       string      `json:"title"`
	Slug        string      `json:"slug"`
	Thumbnail   *newt.Image `json:"thumbnail"`
	Content     string      `json:"content"`
	Tags        []*Tag      `json:"tags"`
	PublishedAt time.Time   `json:"published_at"`
	UpdatedAt   time.Time   `json:"updated_at"`
	Product     *Product    `json:"product"`
}

type ArticleClient struct {
	cc *client.Client
}

func NewArticleClient() *ArticleClient {
	cc := client.New()

	baseURL := fmt.Sprintf("https://%s.cdn.newt.so/v1/", os.Getenv("NEWT_SPACE_UID"))
	if os.Getenv("ENV") == "development" {
		baseURL = fmt.Sprintf("https://%s.api.newt.so/v1/", os.Getenv("NEWT_SPACE_UID"))
	}

	cc.SetBaseURL(baseURL)
	cc.SetHeader("Authorization", "Bearer "+os.Getenv("NEWT_TOKEN"))

	return &ArticleClient{
		cc: cc,
	}
}

func (c *ArticleClient) FetchArticles() ([]*Article, error) {
	// https://www.newt.so/docs/cdn-api-newt-api
	u, err := url.JoinPath(os.Getenv("NEWT_APP_UID"), os.Getenv("NEWT_MODEL_UID"))
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
		return nil, fmt.Errorf("failed to get articles; status=%d", resp.StatusCode())
	}

	contents := new(newt.Contents[*Article])
	if err := resp.JSON(contents); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	return contents.Items, nil
}
func (c *ArticleClient) FetchArticle(slug string) (*Article, error) {
	// https://www.newt.so/docs/cdn-api-newt-api
	u, err := url.JoinPath(os.Getenv("NEWT_APP_UID"), os.Getenv("NEWT_MODEL_UID"))
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
		return nil, fmt.Errorf("failed to get articles; status=%d", resp.StatusCode())
	}

	article := new(newt.Contents[*Article])
	if err := resp.JSON(article); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	if len(article.Items) == 0 {
		return nil, fmt.Errorf("article not found")
	}

	return article.Items[0], nil
}
