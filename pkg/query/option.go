package query

var DefaultQueryLimit = 10

type Statement struct {
	Table  string
	Limit  int
	Offset int
	Where  []struct {
		Column   string
		Operator string
		Value    any
	}
	OrderBy struct {
		Column string
		Desc   bool
	}
}

func (s *Statement) IncrementOffset(amount int) {
	s.Offset += amount
}

func New(table string, options ...option) (*Statement, error) {
	statement := &Statement{
		Table:  table,
		Limit:  DefaultQueryLimit,
		Offset: 0,
		Where: []struct {
			Column   string
			Operator string
			Value    any
		}{},
		OrderBy: struct {
			Column string
			Desc   bool
		}{},
	}

	for _, option := range options {
		option(statement)
	}

	return statement, nil
}

type option func(*Statement)

func WithLimit(limit int) option {
	return func(statement *Statement) {
		statement.Limit = limit
	}
}

func WithOffset(offset int) option {
	return func(statement *Statement) {
		statement.Offset = offset
	}
}

func WithWhere(column, operator string, value any) option {
	return func(statement *Statement) {
		statement.Where = append(statement.Where, struct {
			Column   string
			Operator string
			Value    any
		}{
			Column:   column,
			Operator: operator,
			Value:    value,
		})
	}
}

func WithOrderBy(column string, desc bool) option {
	return func(statement *Statement) {
		statement.OrderBy.Column = column
		statement.OrderBy.Desc = desc
	}
}
