package deebee

import (
	"fmt"
	"reflect"
	"strconv"
	"strings"
	"time"

	"database/sql"

	"github.com/brentlrayjr/types"
)

//Row represents
type Row struct {
	id      int
	columns map[string]*sql.RawBytes
}

func NewRow(id int, columns map[string]*sql.RawBytes) *Row {


	return &Row{id, columns}
}

//Column returns a value of type sql.RawBytes at specified column name
func (row *Row) AddColumn(column string, data *sql.RawBytes)  {


	row.columns[column] = data
}

//Column returns a value of type sql.RawBytes at specified column name
func (row *Row) AddAllColumns(columns map[string]*sql.RawBytes)  {

	for column, data := range columns {

		row.AddColumn(column, data)
	}
}

//Column returns a value of type sql.RawBytes at specified column name
func (row *Row) RemoveColumn(column string)  {

	delete(row.columns, column)
}

//Column returns a value of type sql.RawBytes at specified column name
func (row *Row) Column(column string) *sql.RawBytes {
	if v, b := row.columns[column]; b {
		return v
	}
	return nil
}

//Columns returns a map of strings and sql.RawBytes representing column names and values
func (row *Row) Columns() map[string]*sql.RawBytes {
	return row.columns
}

//ColumnNames returns a slice of strings representing column names
func (row *Row) ColumnNames() []string {

	names := make([]string, 0, len(row.columns))
	for name := range row.columns {
		names = append(names, name)
	}

	return names
}

//ColumnValues returns a slice of sql.RawBytes representing column data
func (row *Row) ColumnValues() []*sql.RawBytes {

	values := make([]*sql.RawBytes, 0, len(row.columns))

	for _, value := range row.columns {
		values = append(values, value)
	}

	return values
}

//Populate takes a struct or struct pointer and populates the fields with values from Row
func (row *Row) Populate(i interface{}) error {

	if reflect.TypeOf(i).Kind() != reflect.Ptr {
		return ErrTypeNotSupported
	}

	s, err := types.Structure(i)
	if err != nil {
		return err
	}

	for name, value := range row.Columns() {

		var val interface{}
		switch name {

		case "id":

			field, err := s.FieldByName("ID")
			if err != nil {
				return err
			}

			val, err = strconv.Atoi(string(*value))
			if err != nil {
				return err
			}

			err = field.Set(val)
			if err != nil {
				return err
			}

			continue

		case "dateTimeUpdated", "dateTimeCreated":

			field, err := s.FieldByName(strings.Title(name))
			if err != nil {
				return fmt.Errorf("Field: %s; Value: %v; Error: %s", field.Name(true), value, err.Error())
			}

			val, err = time.Parse("2006-01-02 15:04:05", string(*value))
			if err != nil {
				return err
			}

			err = field.Set(val)
			if err != nil {
				return err
			}

			continue

		default:

			field, err := s.FieldByName(strings.Title(name))
			if err != nil {
				return fmt.Errorf("Field: %s; Value: %v; Error: %s", field.Name(true), value, err.Error())
			}

			switch field.Value().(type) {

			case int:
				val, err = strconv.Atoi(string(*value))
				if err != nil {
					return err
				}

				break

			case float64:
				val, err = strconv.ParseFloat(string(*value), 64)
				if err != nil {
					return err
				}
				break
			case string:
				val = string(*value)
				break
			case bool:
				val, err = strconv.ParseBool(string(*value))
				if err != nil {
					return err
				}
				break
			case time.Time:

				val, err = time.Parse(time.RFC822, string(*value))
				if err != nil {
					return err
				}
				break

			default:
				return fmt.Errorf("Type: %s; %s", reflect.TypeOf(field.Value()), ErrTypeNotSupported.Error())

			}

			err = field.Set(val)
			if err != nil {
				return err
			}

			continue

		}

	}

	return nil
}
