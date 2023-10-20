package sqlutils

import (
	"fmt"
	"strings"
)

type SelectStatement struct {
	query   strings.Builder
	isFull  bool
	started bool
}

func (s *SelectStatement) SetFull() *SelectStatement {
	s.isFull = true
	return s
}

func (s *SelectStatement) SelectColumns(columns ...Column) *SelectStatement {
	for _, c := range columns {
		if s.started {
			s.query.WriteRune(',')
		}
		if s.isFull {
			s.query.WriteString(Full(c))
		} else {
			s.query.WriteString(c.String())
		}
		s.started = true
	}
	return s
}

func (s *SelectStatement) SelectArrayAggs(columns ...Column) *SelectStatement {
	for _, c := range columns {
		if s.started {
			s.query.WriteRune(',')
		}
		column := c.String()
		if s.isFull {
			column = Full(c)
		}
		s.query.WriteString("array_agg(" + column + ")")
		s.started = true
	}
	return s
}

func (s *SelectStatement) Raw(query string, vars ...any) *SelectStatement {
	s.query.WriteString(fmt.Sprintf(","+query, vars...))
	return s
}

func (s *SelectStatement) String() string {
	return s.query.String()
}
