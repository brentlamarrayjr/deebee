package mysql

import (
	"database/sql"
	"errors"
	"fmt"
	"net/url"
	"strings"
	"time"

	db "../../../deebee"
	"github.com/brentlrayjr/types"
	_ "github.com/go-sql-driver/MySQL" // Reason: Ease of use internally
)

type MySQL struct {
	dsn string
	db  *sql.DB
}

//Connect opens a connection to a DB provider
func New(dsn string) (*MySQL, error) {

	u, err := url.Parse(dsn)
	if err != nil {
		return nil, db.ErrInvalidDSN
	}

	scheme := u.Scheme
	if scheme != "mysql" {
		return nil, db.ErrInvalidDSN
	}

	password, found := u.User.Password()
	if !found {
		return nil, db.ErrInvalidDSN
	}

	db, err := sql.Open(scheme, fmt.Sprintf("%s:%s@tcp(%s)%s", u.User.Username(), password, u.Host, u.Path))
	if err != nil {
		return nil, err
	}

	return &MySQL{dsn: dsn, db: db}, nil
}

func (conn *MySQL) Query(statement *db.Statement) ([]*db.Row, error) {

	rows := make([]*db.Row, 0)

	switch statement.Manipulation() {

	case db.SELECT:

		qRows, err := conn.db.Query(statement.String(), statement.Values()...)
		if err != nil {
			return nil, err
		}

		columns, err := qRows.Columns()
		if err != nil {
			return nil, err
		}


		rowNum := 0

		for qRows.Next() {

			bytes := make([]*sql.RawBytes, len(columns))
			values := make([]interface{}, len(columns))

			for i := 0; i < len(columns); i++ {
				values[i] = new(sql.RawBytes)
			}
			err = qRows.Scan(values...)
			if err != nil {
				return nil, err
			}

			for i := 0; i < len(columns); i++ {
				bytes[i] = values[i].(*sql.RawBytes)
			}

			r := db.NewRow(rowNum, make(map[string]*sql.RawBytes))
			for i, column := range columns {
				r.AddColumn(column, bytes[i])
			}

			rows = append(rows, r)
			rowNum++

		}

		return rows, nil

	case db.INSERT, db.UPDATE, db.DELETE:

		result, err := conn.db.Exec(statement.String(), statement.Values()...)
		if err != nil {
			return nil, err
		}

		rowsAffected, err := result.RowsAffected()
		if err != nil {
			return nil, err
		}

		for rowNum := 0; rowNum < int(rowsAffected); rowNum++ {
			r := db.NewRow(rowNum, nil)
			rows = append(rows, r)
		}

		return rows, nil

	case db.CREATE:

		_, err := conn.db.Exec(statement.String())
		if err != nil {
			return nil, err
		}
		return nil, nil

	default:
		return nil, errors.New("UNSUPPORTED STATEMENT")

	}
}

func (conn *MySQL) Close() error {

	return conn.db.Close()
}
func (conn *MySQL) Ping() error {

	return conn.db.Ping()
}

func (conn *MySQL) TableCount() (int, error) {

	qRows, err := conn.db.Query("SELECT COUNT(*) FROM information_schema.tables WHERE table_schema = 'deebee'")
	if err != nil {
		return -1, err
	}

	var count int

	for qRows.Next() {

		err = qRows.Scan(&count)
		if err != nil {
			return -1, err
		}

	}

	return count, nil
}

func (conn *MySQL) HasTable(i interface{}) bool {

	return false
}

func (conn *MySQL) CreateTable(i interface{}, foreignKeys ...db.ForeignKey) (*db.Statement, error) {



	s, err := types.Structure(i)
	if err != nil {
		return nil, err
	}

	m, err := s.Map(true)
	if err != nil {
		return nil, err
	}

	fields, err := s.Fields()
	if err != nil {
		return nil, err
	}

	if _, found := m["id"]; !found {
		return nil, errors.New("ID field is required")
	}

	if _, found := m["dateTimeUpdated"]; !found {
		return nil, errors.New("DateTimeUpdated field is required")
	}

	if _, found := m["dateTimeCreated"]; !found {
		return nil, errors.New("DateTimeCreated field is required")
	}

	str := ""
	str += "id int PRIMARY KEY NOT NULL AUTO_INCREMENT, "

	for _, field := range fields {

		if name := field.Name(true); name == "id" || name == "dateTimeUpdated" || name == "dateTimeCreated" {
			continue
		}

		switch field.Value().(type) {

		case int:

			tag, err := field.Tag("db")
			if err != nil {
				fmt.Println(err.Error())
			}
			str += fmt.Sprintf("%s INT %s, ", field.Name(true), tag)
			break

		case float64:

			tag, err := field.Tag("db")
			if err != nil {
				fmt.Println(err.Error())
			}
			str += fmt.Sprintf("%s DOUBLE %s, ", field.Name(true), tag)
			break

		case string:

			tag, err := field.Tag("db")
			if err != nil {
				fmt.Println(err.Error())
			}
			str += fmt.Sprintf("%s VARCHAR(65000) %s, ", field.Name(true), tag)
			break

		case bool:

			tag, err := field.Tag("db")
			if err != nil {
				fmt.Println(err.Error())
			}
			fmt.Println("TAG: " + tag)
			str += fmt.Sprintf("%s TINYINT(1) %s, ", field.Name(true), tag)
			break

		case time.Time:

			tag, err := field.Tag("db")
			if err != nil {
				fmt.Println(err.Error())
			}
			str += fmt.Sprintf("%s DATETIME %s,", field.Name(true), tag)
			break

		default:

			return nil, errors.New("Unsupported type")

		}

	}

	str += "dateTimeUpdated DATETIME DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP, dateTimeCreated DATETIME DEFAULT CURRENT_TIMESTAMP"

	fkStrings := make([]string, len(foreignKeys))

	for i := range foreignKeys {
		fkStrings = append(fkStrings, foreignKeys[i].String())
	}

	FKs := strings.Join(fkStrings, ",")
	if len(foreignKeys) == 0 {
		FKs = ", " + FKs
	}

	stmt := db.NewStatement(db.CREATE, nil, nil)

	stmt.AddPart("TABLE", s.Name())
	stmt.AddPart("COLUMNS", str)
	stmt.AddPart("FOREIGN KEYS", FKs)

	return stmt, nil
}

func (conn *MySQL) Migrate(i interface{}) error {

	return db.ErrMethodNotImplemented
}

func (conn *MySQL) Select(i interface{}, columns ...string) (*db.Statement, error) {

	for _, column := range columns {
		if column == "*"{
			//TODO: ERROR
			return nil, fmt.Errorf("")
		}
	}


	s, err := types.Structure(i)
	if err != nil {
		return nil, err
	}

	if len(columns) == 0 {

		names, err := s.Names(true)

		if err != nil {
			return nil, err
		}

		columns = append(columns, names...)

	}else {

		names, err := s.Names(true)

		if err != nil { return nil, err }

		for _, name := range names {

			found := false

			for _, column := range columns {
				if column == name {
					found = true
					break
				}
			}

			if !found {fmt.Errorf("")}
		}
	}

	stmt := db.NewStatement(db.SELECT, nil, nil)
	stmt.AddPart("COLUMNS", strings.Join(columns, ","))
	stmt.AddPart("TABLE", s.Name())

	return stmt, nil
}

func (conn *MySQL) Insert(i interface{}, columns ...string) (*db.Statement, error) {

	s, err := types.Structure(i)
	if err != nil {
		return nil, err
	}

	names := make([]string, 0)
	values := make([]interface{}, 0)


		m, err := s.Map(true)
		if err != nil {
			return nil, err
		}



		for name, value := range m {

			if name == "id" || name == "dateTimeUpdated" || name == "dateTimeCreated" {
				continue
			}


			found := false
			if  len(columns) > 0 {

				for _, column := range columns {

					if column == name {
						found = true
					}

				}
			}else{
				found = true
			}


			if !found {
				return nil, fmt.Errorf("")
			}

			names = append(names, name)
			values = append(values, value)
		}

		placeholders := strings.TrimSpace(strings.Repeat("?, ", len(values)))

		stmt := db.NewStatement(db.INSERT, nil, values)

		stmt.AddPart("TABLE", s.Name())
		stmt.AddPart("COLUMNS", strings.Join(names, ","))
		stmt.AddPart("VALUES", placeholders[:len(placeholders)-1])

		return stmt, nil

	}

func (conn *MySQL) Update(i interface{}, columns ...string) (*db.Statement, error) {

	s, err := types.Structure(i)
	if err != nil {
		return nil, err
	}

	m, err := s.Map(true)
	if err != nil {
		return nil, err
	}

	delete(m, "id")
	delete(m, "dateTimeUpdated")
	delete(m, "dateTimeCreated")

	if len(columns) > 0 {

		for name := range m {

			found := false

			for _, column := range columns {

				if column == name {
					found = true
				}
			}

			if !found {
				delete(m, name)
			}

		}

	}

	placeholders := ""

	values := make([]interface{}, 0)

	for name, value := range m {

		switch value.(type) {

		case int, string, bool, time.Time:
			placeholders += fmt.Sprintf("%s = ?, ", name)
			break
		default:
			return nil, db.ErrTypeNotSupported
		}

		values = append(values, value)

	}

	 stmt := db.NewStatement(db.UPDATE, nil, values)

	 stmt.AddPart("TABLE", s.Name())
	 stmt.AddPart("COLUMNS", placeholders[:len(placeholders)-2])

	 return stmt, nil

}

//Delete return a sql DELETE statement
func (conn *MySQL) Delete(i interface{}) (*db.Statement, error) {

	s, err := types.Structure(i)
	if err != nil {
		return nil, err
	}

	stmt := db.NewStatement(db.DELETE, nil, nil)
	stmt.AddPart("TABLE", s.Name())

	return stmt, nil
}
