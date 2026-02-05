package validator

import (
	"errors"
	"fmt"
	"net/url"
	"time"

	"github.com/google/uuid"
	"github.com/logeshwarann-dev/news-api-rest/internal/model"
	"github.com/logeshwarann-dev/news-api-rest/internal/news"
)

func ValidateNewsRequest(newsRecord model.NewsRecord) (result news.Record, errs error) {

	if newsRecord.Author == "" {
		errs = errors.Join(errs, fmt.Errorf("author is empty: %s", newsRecord.Author))
	}
	if newsRecord.Title == "" {
		errs = errors.Join(errs, fmt.Errorf("title is empty: %s", newsRecord.Title))
	}
	if newsRecord.Summary == "" {
		errs = errors.Join(errs, fmt.Errorf("summary is empty: %s", newsRecord.Summary))
	}
	if newsRecord.Content == "" {
		errs = errors.Join(errs, fmt.Errorf("content is empty: %s", newsRecord.Content))
	}
	var createdTime time.Time
	var err error
	var urlValue *url.URL
	if createdTime, err = time.Parse(time.RFC3339, newsRecord.CreatedAt); err != nil {
		errs = errors.Join(errs, fmt.Errorf("invalid createdAt time: %s", newsRecord.CreatedAt))
	}
	if newsRecord.Source == "" {
		errs = errors.Join(errs, fmt.Errorf("source is empty: %s", newsRecord.Summary))
	}
	if urlValue, err = url.Parse(newsRecord.Source); err != nil {
		errs = errors.Join(errs, fmt.Errorf("invalid source url: %s", newsRecord.Source))
	}
	if len(newsRecord.Tags) == 0 {
		errs = errors.Join(errs, fmt.Errorf("tags are empty"))
	}
	result = news.Record{
		Author:    newsRecord.Author,
		Title:     newsRecord.Title,
		Summary:   newsRecord.Summary,
		Content:   newsRecord.Content,
		CreatedAt: createdTime,
		Source:    urlValue.String(),
		Tags:      newsRecord.Tags,
	}
	return result, errs
}

func ValidateNewsId(id string) (newsId uuid.UUID, err error) {
	if newsId, err = uuid.Parse(id); err != nil {
		return newsId, errors.Join(err, fmt.Errorf("unable to parse newsid: %s", id))
	}
	return newsId, nil
}
