package deebee

import "fmt"

type ForeignKey struct {
	table  string
	column string
}

func FK(table string, column string) *ForeignKey {

	return &ForeignKey{table: table, column: column}
}

func (fk *ForeignKey) String() string {

	return fmt.Sprintf("FOREIGN KEY(%s) REFERENCES %s(%s) ON DELETE CASCADE", fk.column, fk.table, fk.column)
}
