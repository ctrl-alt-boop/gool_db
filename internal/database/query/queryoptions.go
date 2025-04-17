package query

type QueryOptions struct {
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
