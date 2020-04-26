package parser

import (
	"bytes"
	"errors"
	. "github.com/saichler/utils/golang"
	"strconv"
	"strings"
)

type Query struct {
	query      string
	schemaName string
	sortBy     string
	descending bool
	limit      int
	page       int
	matchCase  bool
	tables     []string
	column     []string
	where      *Expression
}

type parsed struct {
	select_     []string
	from_       []string
	where_      string
	sortby_     string
	descending_ string
	ascending_  string
	limit_      string
	page_       string
	matchcase_  string
}

const (
	Select     = "select"
	From       = "from"
	Where      = "where"
	SortBy     = "sort-by"
	Descending = "descending"
	Ascending  = "ascending"
	Limit      = "limit"
	Page       = "page"
	MatchCase  = "match-case"
)

var words = []string{Select, From, Where, SortBy, Descending, Ascending, Limit, Page, MatchCase}

func (query *Query) Where() *Expression {
	return query.where
}

func (query *Query) Tables() []string {
	return query.tables
}

func (query *Query) Columns() []string {
	return query.column
}

func (query *Query) SchemaName() string {
	return query.schemaName
}

func (query *Query) SortBy() string {
	return query.sortBy
}

func (query *Query) Descending() bool {
	return query.descending
}

func (query *Query) Limit() int {
	return query.limit
}

func (query *Query) Page() int {
	return query.page
}

func (query *Query) MatchCase() bool {
	return query.matchCase
}

func NewQuery(query string) (*Query, error) {
	cwql := &Query{}
	cwql.query = query
	e := cwql.init()
	return cwql, e
}

func TrimAndLowerNoKeys(sql string) string {
	buff := bytes.Buffer{}
	sql = strings.TrimSpace(sql)
	keyOpen := false
	for _, c := range sql {
		if c == '[' {
			keyOpen = true
		} else if c == ']' {
			keyOpen = false
		}
		if !keyOpen {
			buff.WriteString(strings.ToLower(string(c)))
		} else {
			buff.WriteString(string(c))
		}
	}
	return buff.String()
}

func (query *Query) split() *parsed {
	sql := TrimAndLowerNoKeys(query.query)
	data := &parsed{}
	data.select_ = getSplitTag(sql, Select)
	data.from_ = getSplitTag(sql, From)
	data.where_ = getTag(sql, Where)
	data.descending_ = getBoolTag(sql, Descending)
	data.ascending_ = getBoolTag(sql, Ascending)
	data.limit_ = getTag(sql, Limit)
	data.page_ = getTag(sql, Page)
	data.sortby_ = getTag(sql, SortBy)
	data.matchcase_ = getBoolTag(sql, MatchCase)
	return data
}

func getBoolTag(str, tag string) string {
	index := strings.Index(str, tag)
	if index != -1 {
		return "true"
	}
	return "false"
}

func getTag(str, tag string) string {
	index := strings.Index(str, tag)
	if index == -1 {
		return ""
	}
	index += len(tag)
	index2 := len(str)
	for _, t := range words {
		if t != tag {
			index3 := strings.Index(str, t)
			if index3 > index && index3 < index2 {
				index2 = index3
			}
		}
	}
	return strings.TrimSpace(str[index:index2])
}

func getSplitTag(str, tag string) []string {
	result := make([]string, 0)
	data := getTag(str, tag)
	if data == "" {
		return result
	}
	split := strings.Split(data, ",")
	for _, t := range split {
		result = append(result, t)
	}
	return result
}

func (query *Query) init() error {
	p := query.split()
	query.column = make([]string, 0)
	for _, col := range p.select_ {
		query.column = append(query.column, strings.TrimSpace(col))
	}
	query.tables = make([]string, 0)
	for _, tbl := range p.from_ {
		index := strings.Index(tbl, ".")
		if index != -1 {
			query.schemaName = tbl[0:index]
			tbl = tbl[index+1:]
		}
		query.tables = append(query.tables, strings.TrimSpace(tbl))
	}
	if p.where_ != "" {
		where, e := parseExpression(p.where_)
		if e != nil {
			return e
		}
		query.where = where
	}
	if p.limit_ != "" {
		limit, e := strconv.Atoi(p.limit_)
		if e != nil {
			Error("Invalid limit:" + p.limit_ + ", setting limity to 10")
			limit = 10
		}
		if limit >= 1000 {
			msg := "Invalid limit: Limit is limited up to 1000 elements"
			Error(msg)
			return errors.New(msg)
		}
		query.limit = limit
	}
	if p.page_ != "" {
		page, e := strconv.Atoi(p.page_)
		if e != nil {
			Error("Invalid page:" + p.page_)
			return e
		}
		query.page = page
	}
	query.sortBy = p.sortby_
	if p.descending_ == "true" {
		query.descending = true
	}
	if p.ascending_ == "true" {
		query.descending = false
	}
	if p.matchcase_ != "" {
		query.matchCase = true
	}
	return nil
}
