package deebee

import "fmt"

type Condition struct {
	opr    operator
	field  string
	values []interface{}
}

func (e *Condition) String() string {

	switch e.opr {

	case EQUAL:
		return fmt.Sprintf(EQUAL.String(), e.field, "?")

	case LESS_THAN:
		return fmt.Sprintf(LESS_THAN.String(), e.field, "?")

	case LESS_THAN_EQUAL:
		return fmt.Sprintf(LESS_THAN_EQUAL.String(), e.field, "?")

	case GREATER_THAN:
		return fmt.Sprintf(GREATER_THAN.String(), e.field, "?")

	case GREATER_THAN_EQUAL:
		return fmt.Sprintf(GREATER_THAN_EQUAL.String(), e.field, "?")

	default:
		return ""

	}

}

//EqualTo function
func EqualTo(field string, value interface{}) (cond *Condition) {

	cond = &Condition{}

	cond.opr = EQUAL
	cond.field = field
	cond.values = append(cond.values, value)
	return cond
}

func LessThan(field string, value interface{}) (cond *Condition) {

	cond = &Condition{}

	cond.opr = LESS_THAN
	cond.field = field
	cond.values = append(cond.values, value)
	return cond
}

func GreaterThan(field string, value interface{}) (cond *Condition) {

	cond = &Condition{}

	cond.opr = GREATER_THAN
	cond.field = field
	cond.values = append(cond.values, value)
	return cond
}
