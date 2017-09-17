package deebee

type operator string

var ASC operator = "ASC"
var DESC operator = "DESC"

var GREATER_THAN operator = "%s > %s"
var GREATER_THAN_EQUAL operator = "%s >= %s"
var LESS_THAN operator = "%s < %s"
var LESS_THAN_EQUAL operator = "%s <= %s"
var EQUAL operator = "%s = %s"
var NOT_EQUAL operator = "NOT (%s) = %s"
var IN operator = "IN (%s)"
var NOT_IN operator = "NOT IN (%s)"

func (opr operator) String() string {
	return string(opr)
}
