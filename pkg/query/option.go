package query

import "fmt"

var DefaultQueryLimit = 10

type SqlMethod string

const (
	Select SqlMethod = "SELECT"
	Insert SqlMethod = "INSERT"
	Update SqlMethod = "UPDATE"
	Delete SqlMethod = "DELETE"
)

const AllColumns = "*"

var ErrEmptyMethodProvided = fmt.Errorf("empty method provided")
var ErrEmptyTableProvided = fmt.Errorf("empty table provided")

type StatementProcessor func(*Statement) error

type Statement struct {
	Method  SqlMethod
	Table   string
	Columns []string
	Values  []any
	Set     []struct {
		Column string
		Value  any
	}
	Where []struct {
		Column   string
		Operator string
		Value    any
	}
	OrderBy struct {
		Column string
		Desc   bool
	}
	Limit  int
	Offset int

	PreProcess  StatementProcessor
	PostProcess StatementProcessor
}

func (s *Statement) IncrementOffset(amount int) {
	s.Offset += amount
}

func New(table string, method SqlMethod, options ...option) (*Statement, error) {
	statement := &Statement{
		Table:   table,
		Method:  method,
		Columns: []string{AllColumns},
		Where: []struct {
			Column   string
			Operator string
			Value    any
		}{},
		OrderBy: struct {
			Column string
			Desc   bool
		}{},
		Limit:  DefaultQueryLimit,
		Offset: 0,
	}

	for _, option := range options {
		option(statement)
	}
	if statement.Method == "" {
		return nil, ErrEmptyMethodProvided
	}
	if statement.Table == "" {
		return nil, ErrEmptyTableProvided
	}

	return statement, nil
}

type option func(*Statement)

func WithMethod(method SqlMethod) option {
	return func(statement *Statement) {
		statement.Method = method
	}
}

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

func WithPreProcess(processor StatementProcessor) option {
	return func(statement *Statement) {
		statement.PreProcess = processor
	}
}

func WithPostProcess(processor StatementProcessor) option {
	return func(statement *Statement) {
		statement.PostProcess = processor
	}
}
