package deebee

import (
	"fmt"
	"strings"
	"strconv"
)

//Statement represent a Statement built using provided Statement builders
type Statement struct {
	manipulation Manipulation
	parts      *Parts
	values       []interface{}
}

func NewStatement(manipulation Manipulation, parts map[string][]string, values []interface{}) *Statement {

	if parts == nil {
		parts = make(map[string][]string)
	}

	if values == nil {
		values = make([]interface{}, 0)
	}

	return &Statement{manipulation, NewParts(), values}
}

func (stmt *Statement) Manipulation() Manipulation {

	return stmt.manipulation
}

func (stmt *Statement) AddPart(name string, part string) {

	stmt.parts.Add(name, part)
}

func (stmt *Statement) Values() []interface{} {

	return stmt.values
}

//Where adds or appends to WHERE clause in Statement
func (stmt *Statement) Where(conditions ...*Condition) {

	conds := make([]string, 0)

	for _, cond := range conditions {

		conds = append(conds, cond.String())
		stmt.values = append(stmt.values, cond.values...)

	}

		stmt.parts.Add("WHERE", strings.Join(conds, " AND "))
	}

//Having adds or appends to WHERE clause in Statement
func (stmt *Statement) Having(conditions ...*Condition) {

	conds := make([]string, 0)

	for _, cond := range conditions {

		conds = append(conds, cond.String())
		stmt.values = append(stmt.values, cond.values...)

	}

	stmt.parts.Add("HAVING", strings.Join(conds, " AND "))

}

//GroupBy appends GROUP BY clause to Statement
func (stmt *Statement) GroupBy(column string, order string) {

	if order != "ASC" &&  order != "DESC" {
		return
	}

	stmt.parts.Add("GROUP BY", column + " " + order)
}

//OrderBy appends ORDER BY clause to Statement
func (stmt *Statement) OrderBy(order string, fields ...string) {

	if order != "ASC" &&  order != "DESC" {
		return
	}

	stmt.parts.Add("ORDER BY", strings.Join(fields, ",") + " " + order)
}

//Limit ...
func (stmt *Statement) Limit(offset int, limit int) {

	if !stmt.parts.Has("LIMIT") {
		stmt.parts.Add("LIMIT", strconv.Itoa(offset) + ", " + strconv.Itoa(limit))
	}
}

func (stmt *Statement) String() string {

	statement := ""

	switch stmt.manipulation {

	case CREATE:
		statement += fmt.Sprintf("CREATE TABLE IF NOT EXISTS %s (%s) ", stmt.parts.Get("TABLE")[0], stmt.parts.Get("COLUMNS")[0])
		break
	case SELECT:
		statement += fmt.Sprintf("SELECT %s FROM %s ", stmt.parts.Get("COLUMNS")[0], stmt.parts.Get("TABLE")[0])
		break
	case INSERT:
		statement += fmt.Sprintf("INSERT INTO %s(%s) VALUES (%s) ", stmt.parts.Get("TABLE")[0], stmt.parts.Get("COLUMNS")[0], stmt.parts.Get("VALUES")[0])
		break
	case UPDATE:
		statement += fmt.Sprintf("UPDATE %s SET %s ", stmt.parts.Get("TABLE")[0], stmt.parts.Get("COLUMNS")[0])
		break
	case DELETE:
		statement += fmt.Sprintf("DELETE FROM %s ", stmt.parts.Get("TABLE")[0])
		break

	}

	if stmt.parts.Has("WHERE") {

		statement += "WHERE "

		for i, where := range stmt.parts.Get("WHERE") {

			if i < stmt.parts.Len("WHERE")-1 {
				statement += "(" + where + ") OR "
				continue
			}

			statement += "(" + where + ")"

		}
	}

	if stmt.parts.Has("GROUP BY") {
		statement += "GROUP BY " + strings.Join(stmt.parts.Get("GROUP BY"), ",") + " "
	}

	if stmt.parts.Has("HAVING") {

		statement += "HAVING "

		for i, having := range stmt.parts.Get("HAVING") {

			if i < stmt.parts.Len("HAVING")-1 {
				statement += "(" + having + ") OR "
				continue
			}

			statement += "(" + having + ")"

		}

		statement += " "
	}

	if stmt.parts.Has("ORDER BY") {
		statement += "ORDER BY " + stmt.parts.Get("ORDER BY")[0] + " "
	}

	if stmt.parts.Has("LIMIT") {
		statement += "LIMIT " + stmt.parts.Get("LIMIT")[0] + " "
	}


	return statement
}
