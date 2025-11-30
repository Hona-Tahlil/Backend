package dsl

import (
	"net/url"
	"strconv"
	"strings"
)

type ParsedQuery struct {
	Filters *Filters
	Sorts   *Sorts
	Page    int
	Limit   int
}
// ctx := *gin.Contex
// ctx.Request.URL.Query()
func ParseQuery(values url.Values) *ParsedQuery {

	filters := NewFilters()
	sorts := NewSorts()
	page := 1
	limit := 20

	for key, vals := range values {

		value := vals[0]

		// page
		if key == "page" {
			if p, err := strconv.Atoi(value); err == nil {
				page = p
			}
			continue
		}

		// limit
		if key == "limit" {
			if l, err := strconv.Atoi(value); err == nil {
				limit = l
			}
			continue
		}

		// sort
		if key == "sort" {
			if strings.HasPrefix(value, "-") {
				sorts.Desc(value[1:])
			} else {
				sorts.Asc(value)
			}
			continue
		}

		// IN: status=in:active,in_review
		if strings.HasPrefix(value, "in:") {
			items := strings.Split(value[3:], ",")
			anySlice := make([]any, len(items))
			for i, v := range items {
				anySlice[i] = v
			}
			filters.In(key, anySlice)
			continue
		}

		// LIKE: name~ali  OR  name~=ali
		if strings.Contains(value, "~") {
			parts := strings.Split(value, "~")
			filters.Like(key, parts[1])
			continue
		}

		// >= یا <=
		if strings.HasPrefix(value, ">=") {
			filters.Gte(key, value[2:])
			continue
		}
		if strings.HasPrefix(value, "<=") {
			filters.Lte(key, value[2:])
			continue
		}

		// > یا <
		if strings.HasPrefix(value, ">") {
			filters.Gt(key, value[1:])
			continue
		}
		if strings.HasPrefix(value, "<") {
			filters.Lt(key, value[1:])
			continue
		}

		// پیشفرض: EQ
		filters.Eq(key, value)
	}

	return &ParsedQuery{
		Filters: filters,
		Sorts:   sorts,
		Page:    page,
		Limit:   limit,
	}
}
