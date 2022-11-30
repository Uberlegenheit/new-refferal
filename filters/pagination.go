package filters

import (
	"github.com/wedancedalot/squirrel"
)

const (
	MaxLimitSize = 1000
	DefaultLimit = 100
)

type (
	Pagination struct {
		Limit uint64 `schema:"limit"`
		Page  uint64 `schema:"page"`
	}
	PeriodInfoRequest struct {
		Start string `schema:"start"`
		End   string `schema:"end"`
	}
)

func (p *Pagination) SetFilter(q squirrel.SelectBuilder) squirrel.SelectBuilder {
	if p.Limit != 0 {
		q = q.Limit(p.Limit)
	}
	if p.Page > 1 {
		offset := (p.Page - 1) * p.Limit
		q = q.Offset(offset)
	}
	return q
}

func (p *Pagination) Validate() {
	if p.Page == 0 {
		p.Page = 1
	}

	if p.Limit > MaxLimitSize {
		p.Limit = MaxLimitSize
	}
	if p.Limit == 0 {
		p.Limit = DefaultLimit
	}
}

func (p *Pagination) Offset() uint64 {
	return p.Limit * (p.Page - 1)
}
