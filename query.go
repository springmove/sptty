package sptty

import (
	"strconv"
	"strings"

	"github.com/kataras/iris/v12"
	"gorm.io/gorm"
)

const (
	DefaultPageSize = 20
)

type IQuery interface {
	FromCtx(ctx iris.Context)
	loadDB(db *gorm.DB)
	ToQuery(paging bool) *gorm.DB
}

type QueryBase struct {
	IQuery

	db *gorm.DB

	Page     int64
	PageSize int64
	IDs      []string
}

func (s *QueryBase) loadDB(db *gorm.DB) {
	s.db = db
}

func (s *QueryBase) FromCtx(ctx iris.Context) {
	ids := ctx.URLParam("IDs")
	if ids != "" {
		s.IDs = strings.Split(ids, ",")
	}

	page, err := strconv.ParseInt(ctx.URLParam("Page"), 10, 32)
	if err != nil {
		page = 0
	}

	if page > 0 {
		page -= 1
	}

	s.Page = page

	pageSize, err := strconv.ParseInt(ctx.URLParam("PageSize"), 10, 32)
	if err != nil {
		pageSize = DefaultPageSize
	}

	s.PageSize = pageSize
}

func (s *QueryBase) ToQuery(paging bool) *gorm.DB {
	q := s.db

	q = q.Where("deleted = ?", false)

	if s.PageSize == 0 {
		s.PageSize = DefaultPageSize
	}

	if len(s.IDs) > 0 {
		q = q.Where("id in (?)", s.IDs)
	}

	if paging {
		q = q.Limit(int(s.PageSize)).Offset(int(s.Page * s.PageSize))
	}

	return q
}

func CreateQueryFromContext(query IQuery, db *gorm.DB, ctx ...iris.Context) IQuery {
	query.loadDB(db)

	if len(ctx) > 0 {
		query.FromCtx(ctx[0])
	}

	return query
}
