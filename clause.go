package deebee

type clause string

var WHERE clause = ""
var HAVING clause = ""
var ORDER_BY clause = ""
var LIKE clause = ""
var GROUP_BY clause = ""
var LIMIT clause = ""
var OR clause = "OR"

func (c *clause) String() string {
	return string(*c)
}
