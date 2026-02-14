package metadev

import (
	"github.com/MetaDiv-AI/metaorm"
)

type Empty struct{}

type RequestPathId struct {
	ID uint `uri:"id" json:"-"`
}

type RequestListing struct {
	metaorm.PaginationQueryImpl
	metaorm.SortingQueryImpl
	Keyword string `form:"keyword"`
}

func (r *RequestListing) PageQuery() metaorm.PaginationQuery {
	return metaorm.Paginate(r.Page, r.Size)
}

func (r *RequestListing) SortingQuery() metaorm.SortingQuery {
	return metaorm.Sort(r.Field, r.Asc)
}

// SimilarKeyword build a query for similar keyword
// if the field name is prefixed with "*", it will be treated as a decrypted field
func (r *RequestListing) SimilarKeyword(fields ...string) metaorm.Query {
	if len(fields) == 0 {
		return nil
	}

	qb := metaorm.NewQueryBuilder()

	queries := make([]metaorm.Query, 0)
	for _, field := range fields {
		queries = append(queries, qb.Field(field).Similar(r.Keyword))
	}
	return qb.Or(queries...)
}
