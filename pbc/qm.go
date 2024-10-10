package pbc

import (
	"bytes"
	"encoding/json"
	"io"
	"log/slog"
	"strings"
)

const (
	QTypeData    = "data"
	QTypeExpand  = "expand"
	QTypeFields  = "fields"
	QTypeFilters = "filters"
	QTypeHeaders = "headers"
	QTypePage    = "page"
	QTypeParams  = "params"
	QTypeSort    = "sort"
)

// QueryOption is a function type for configuring Query
type QueryOption func(*Query)

// Query represents a query with various options
type Query struct {
	Headers   map[string]string
	Params    []any
	Data      any
	DataBytes io.Reader
	Filters   string
	Sort      []string
	Fields    []string
	Page      int
	PerPage   int
	SkipTotal bool
	Expand    []string
}

// NewQuery creates a new Query with the given options
func NewQuery(opts ...QueryOption) *Query {
	q := &Query{
		Headers: make(map[string]string),
	}
	for _, opt := range opts {
		opt(q)
	}
	return q
}

// Headers represents
type Headers = map[string]string

// WithHeaders adds headers to the query
func WithHeaders(headers Headers) QueryOption {
	return func(q *Query) {
		for k, v := range headers {
			q.Headers[k] = v
		}
	}
}

// WithAuthorization adds Authorization header to request
func WithAuthorization(token string) QueryOption {
	return WithHeaders(Headers{"Authorization": token})
}

// WithParams adds parameters to the query
func WithParams(params ...any) QueryOption {
	return func(q *Query) {
		q.Params = append(q.Params, params...)
	}
}

// WithData adds data to the query
func WithData(data any) QueryOption {
	return func(q *Query) {
		q.Data = data
		q.DataBytes = jsonMarshal(data)
	}
}

// WithFilters adds filters to the query
func WithFilters(filters string) QueryOption {
	return func(q *Query) {
		q.Filters = filters
	}
}

// WithSort adds sort options to the query
func WithSort(sort ...string) QueryOption {
	return func(q *Query) {
		q.Sort = append(q.Sort, sort...)
	}
}

// WithFields adds fields to the query
func WithFields(fields ...string) QueryOption {
	return func(q *Query) {
		q.Fields = append(q.Fields, fields...)
	}
}

// WithPage adds pagination to the query
func WithPage(page, perPage int, skipTotal bool) QueryOption {
	return func(q *Query) {
		q.Page = page
		q.PerPage = perPage
		q.SkipTotal = skipTotal
	}
}

// WithExpand adds expand options to the query
func WithExpand(expand ...string) QueryOption {
	return func(q *Query) {
		q.Expand = append(q.Expand, expand...)
	}
}

// Helper functions

func QmListString(ss []string) string {
	return strings.Join(ss, ",")
}

func jsonMarshal(d any) io.Reader {
	b, err := json.Marshal(d)
	if err != nil {
		slog.Error("failed to marshal payload", "payload", d)
	}
	return bytes.NewReader(b)
}
