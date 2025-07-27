package store

import (
	"net/http"
	"strconv"
	"time"
)

type Pagination struct {
	Limit  int      `json:"limit" validate:"gte=0,lte=20"`
	Offset int      `json:"offset" validate:"gte=0"`
	Order  string   `json:"order" validate:"oneof=asc desc"`
	Search string   `json:"search" validate:"omitempty,max=100"`
	Since  *string  `json:"since" validate:"omitempty,datetime=2006-01-02T15:04:05Z07:00"`
	Tags   []string `json:"tags" validate:"dive,max=50"`
}

func (p *Pagination) GetPaginated(r *http.Request) (*Pagination, error) {
	query := r.URL.Query()

	limit := query.Get("limit")
	if limit != "" {
		var err error
		p.Limit, err = strconv.Atoi(limit)
		if err != nil {
			return nil, err
		}
	}
	offset := query.Get("offset")
	if offset != "" {
		var err error
		p.Offset, err = strconv.Atoi(offset)
		if err != nil {
			return nil, err
		}
	}
	order := query.Get("order")
	if order != "" {
		p.Order = order
	}

	search := query.Get("search")
	if search != "" {
		p.Search = search
	}
	since := query.Get("since")
	if since != "" {
		sin := p.formatTime(since)
		p.Since = &sin
	}
	tags := query["tags"]
	if len(tags) > 0 {
		p.Tags = tags
	}

	return p, nil
}

func (p *Pagination) formatTime(since string) string {
	t, err := time.Parse(time.RFC3339, since)
	if err != nil {
		return since
	}
	return t.Format(time.DateTime)
}
