package templates

import (
	"text/template"
)

const queryTemplate = `SELECT {{.ColumnsString}} FROM {{.Table}}
WHERE {{.WhereString}}
ORDER BY {{.OrderByString}}
LIMIT {{.Limit}} OFFSET {{.Offset}}`

const countTemplate = `SELECT COUNT(*) FROM {{.Table}}
WHERE {{.Where}}`

func QueryOptions() *template.Template {
	return template.Must(template.New("query").Parse(queryTemplate))
}

func CountOptions() *template.Template {
	return template.Must(template.New("count").Parse(countTemplate))
}
