package sqlutils

import (
	"strings"
)

type QueryBuilder struct {
	query  strings.Builder
	isFull bool
}

func (q *QueryBuilder) SetFull() *QueryBuilder {
	q.isFull = true
	return q
}

func (q *QueryBuilder) InsertInto(table string, values ...Column) *QueryBuilder {
	q.query.WriteString("INSERT INTO " + table + " ")
	if len(values) == 0 {
		return q
	}
	q.query.WriteString("(")
	q.writeSplitColumns(values)
	q.query.WriteString(") ")
	return q
}

func (q *QueryBuilder) DeleteFrom(table string) *QueryBuilder {
	q.query.WriteString("DELETE FROM " + table + " ")
	return q
}

func (q *QueryBuilder) Returning(columns ...Column) *QueryBuilder {
	q.query.WriteString("RETURNING ")
	q.writeSplitColumns(columns)
	return q
}

func (q *QueryBuilder) Select(columns ...Column) *QueryBuilder {
	q.query.WriteString("SELECT ")
	if q.isFull {
		q.query.WriteString(Full(columns...))
	} else {
		q.writeSplitColumns(columns)
	}
	q.query.WriteString(" ")
	return q
}

func (q *QueryBuilder) SelectStatement(s *SelectStatement) *QueryBuilder {
	q.query.WriteString("SELECT " + s.String() + " ")
	return q
}

func (q *QueryBuilder) From(table string) *QueryBuilder {
	q.query.WriteString("FROM " + table + " ")
	return q
}

func (q *QueryBuilder) WhereColumnEqual(c Column, raw string) *QueryBuilder {
	column := q.string(c)
	q.query.WriteString("WHERE " + column + " = " + raw + " ")
	return q
}

func (q *QueryBuilder) InnerJoin(table string) *QueryBuilder {
	q.query.WriteString("INNER JOIN " + table + " ")
	return q
}

func (q *QueryBuilder) LeftJoin(table string) *QueryBuilder {
	q.query.WriteString("LEFT JOIN " + table + " ")
	return q
}

func (q *QueryBuilder) OnColumnsEquals(c1, c2 Column) *QueryBuilder {
	column1, column2 := q.string(c1), q.string(c2)
	q.query.WriteString("ON " + column1 + " = " + column2 + " ")
	return q
}

func (q *QueryBuilder) And() *QueryBuilder {
	q.query.WriteString("AND ")
	return q
}

func (q *QueryBuilder) Or() *QueryBuilder {
	q.query.WriteString("OR ")
	return q
}

func (q *QueryBuilder) Not() *QueryBuilder {
	q.query.WriteString("NOT ")
	return q
}

func (q *QueryBuilder) ColumnsEquals(c1, c2 Column) *QueryBuilder {
	column1, column2 := q.string(c1), q.string(c2)
	q.query.WriteString(column1 + " = " + column2 + " ")
	return q
}

func (q *QueryBuilder) ColumnEquals(c Column, raw string) *QueryBuilder {
	column := q.string(c)
	q.query.WriteString(column + " = " + raw + " ")
	return q
}

func (q *QueryBuilder) ColumnsNotEquals(c1, c2 Column) *QueryBuilder {
	column1, column2 := q.string(c1), q.string(c2)
	q.query.WriteString(column1 + " != " + column2 + " ")
	return q
}

func (q *QueryBuilder) ColumnNotEquals(c Column, raw string) *QueryBuilder {
	column := q.string(c)
	q.query.WriteString(column + " != " + raw + " ")
	return q
}

func (q *QueryBuilder) GroupBy(columns ...Column) *QueryBuilder {
	q.query.WriteString("GROUP BY ")
	if q.isFull {
		q.query.WriteString(Full(columns...))
	} else {
		q.writeSplitColumns(columns)
	}
	return q
}

func (q *QueryBuilder) Update(table string) *QueryBuilder {
	q.query.WriteString("UPDATE " + table + " SET ")
	return q
}

func (q *QueryBuilder) SetColumn(c Column, raw string) *QueryBuilder {
	column := q.string(c)
	q.query.WriteString(column + " = " + raw + " ")
	return q
}

func (q *QueryBuilder) WriteQuery(s string) *QueryBuilder {
	q.query.WriteString(s)
	q.query.WriteRune(' ')
	return q
}

func (q *QueryBuilder) Build() string {
	return q.query.String()
}

func (q *QueryBuilder) string(c Column) string {
	column := c.String()
	if q.isFull {
		column = Full(c)
	}
	return column
}

func (q *QueryBuilder) writeSplitColumns(columns []Column) {
	lastIndex := len(columns) - 1
	for i, v := range columns {
		q.query.WriteString(v.String())
		if i != lastIndex {
			q.query.WriteRune(',')
		}
	}
}
