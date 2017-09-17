package mysql

import (
	"fmt"
	"testing"
	"time"
	"../../../deebee"

	"github.com/stretchr/testify/require"
)

func TestMethodNew(t *testing.T) {

	type Struct struct {
		ID              int    `db:""`
		Name            string `db:""`
		Active          bool   `db:"NOT NULL"`
		DateTimeUpdated time.Time
		DateTimeCreated time.Time
	}

	db, err := New("mysql://root:toor@localhost:3306/db")
	require.NoError(t, err, "MySQL DB constructor should not throw an error")
	defer db.Close()

	err = db.Ping()
	require.NoError(t, err, "Ping should not throw an error")

}

func TestMethodCreateTable(t *testing.T) {

	t.SkipNow()

	type Struct struct {
		ID              int    `db:""`
		Name            string `db:""`
		Active          bool   `db:"NOT NULL"`
		DateTimeUpdated time.Time
		DateTimeCreated time.Time
	}

	db, err := New("mysql://root:toor@localhost:3306/db")
	require.NoError(t, err, "MySQL DB constructor should not throw an error")
	defer db.Close()

	err = db.Ping()
	require.NoError(t, err, "Ping should not throw an error")

	stmt, err := db.CreateTable(&Struct{})
	require.NoError(t, err, "CreateTable function should not throw an error")

	_, err = db.Query(stmt)
	require.NoError(t, err, "Query function should not throw an error")

}

func TestMethodInsert(t *testing.T) {

	type Struct struct {
		ID              int    `db:""`
		Name            string `db:""`
		Active          bool   `db:"NOT NULL"`
		DateTimeUpdated time.Time
		DateTimeCreated time.Time
	}

	db, err := New("mysql://root:toor@localhost:3306/db")
	require.NoError(t, err, "MySQL DB constructor should not throw an error")
	defer db.Close()

	err = db.Ping()
	require.NoError(t, err, "Ping should not throw an error")

	stmt, err := db.Insert(&Struct{Name: "db", Active: true})
	require.NoError(t, err, "Select function should not throw an error")
	fmt.Println("INSERT: " + stmt.String())

	_, err = db.Query(stmt)
	require.NoError(t, err, "Query function should not throw an error")

}

func TestMethodUpdate(t *testing.T) {

	type Struct struct {
		ID              int    `db:""`
		Name            string `db:""`
		Active          bool   `db:"NOT NULL"`
		DateTimeUpdated time.Time
		DateTimeCreated time.Time
	}

	db, err := New("mysql://root:toor@localhost:3306/db")
	require.NoError(t, err, "MySQL DB constructor should not throw an error")
	defer db.Close()

	err = db.Ping()
	require.NoError(t, err, "Ping should not throw an error")

	stmt, err := db.Update(&Struct{Name: "Brent", Active: true})
	require.NoError(t, err, "Select function should not throw an error")
	fmt.Println("UPDATE: " + stmt.String())

	//_, err = db.Query(stmt)
	//require.NoError(t, err, "Query function should not throw an error")

}

func TestMethodQuery(t *testing.T) {

	type Struct struct {
		ID              int    `db:""`
		Name            string `db:""`
		Active          bool   `db:"NOT NULL"`
		DateTimeUpdated time.Time
		DateTimeCreated time.Time
	}

	db, err := New("mysql://root:toor@localhost:3306/db")
	require.NoError(t, err, "MySQL DB constructor should not throw an error")
	defer db.Close()

	err = db.Ping()
	require.NoError(t, err, "Ping should not throw an error")

	stmt, err := db.Select(&Struct{})
	require.NoError(t, err, "Select function should not throw an error")

	rows, err := db.Query(stmt)
	require.NoError(t, err, "Query function should not throw an error")
	require.Truef(t, len(rows) > 0, "Query method should return rows but lenght = %d", len(rows))
	fmt.Printf("ROW COUNT: %d \n", len(rows))
	fmt.Printf("QUERY: %s \n", stmt.String())
	fmt.Println()

}

func TestMethodWhere(t *testing.T) {

	type Struct struct {
		ID              int    `db:""`
		Name            string `db:""`
		Active          bool   `db:"NOT NULL"`
		DateTimeUpdated time.Time
		DateTimeCreated time.Time
	}

	db, err := New("mysql://root:toor@localhost:3306/db")
	require.NoError(t, err, "MySQL DB constructor should not throw an error")
	defer db.Close()

	err = db.Ping()
	require.NoError(t, err, "Ping should not throw an error")

	stmt, err := db.Select(&Struct{})
	require.NoError(t, err, "Select function should not throw an error")

	stmt.Where(deebee.GreaterThan("id", 2))
	stmt.Where(deebee.EqualTo("active", true), deebee.LessThan("dateTimeUpdated", time.Now().String()))

	require.Equal(t, "SELECT id,name,active,dateTimeUpdated,dateTimeCreated FROM Struct WHERE (id > ?) OR (active = ? AND dateTimeUpdated < ?)", stmt.String(), "Statement string should match expected value.")

	rows, err := db.Query(stmt)
	require.NoError(t, err, "Query function should not throw an error")
	require.Truef(t, len(rows) > 0, "Query method should return rows but lenght = %d", len(rows))

	fmt.Printf("ROW COUNT: %d \n", len(rows))
	fmt.Printf("QUERY: %s \n", stmt.String())
	fmt.Println()

}

func TestMethodGroupBy(t *testing.T) {

	t.SkipNow()

	type Struct struct {
		ID              int    `db:""`
		Name            string `db:""`
		Active          bool   `db:"NOT NULL"`
		DateTimeUpdated time.Time
		DateTimeCreated time.Time
	}

	db, err := New("mysql://root:toor@localhost:3306/db")
	require.NoError(t, err, "MySQL DB constructor should not throw an error")
	defer db.Close()

	err = db.Ping()
	require.NoError(t, err, "Ping should not throw an error")

	stmt, err := db.Select(&Struct{})
	require.NoError(t, err, "Select function should not throw an error")

	stmt.GroupBy("dateTimeUpdated", "ASC")

	require.Equal(t, "SELECT id,name,active,dateTimeUpdated,dateTimeCreated FROM Struct GROUP BY dateTimeUpdated ASC " , stmt.String(), "Statement string should match expected value.")

	rows, err := db.Query(stmt)
	require.NoError(t, err, "Query function should not throw an error")
	require.Truef(t, len(rows) > 0, "Query method should return rows but length = %d", len(rows))

	fmt.Printf("ROW COUNT: %d \n", len(rows))
	fmt.Printf("QUERY: %s \n", stmt.String())
	fmt.Println()

}

func TestMethodHaving(t *testing.T) {

	type Struct struct {
		ID              int    `db:""`
		Name            string `db:""`
		Active          bool   `db:"NOT NULL"`
		DateTimeUpdated time.Time
		DateTimeCreated time.Time
	}

	db, err := New("mysql://root:toor@localhost:3306/db")
	require.NoError(t, err, "MySQL DB constructor should not throw an error")
	defer db.Close()

	err = db.Ping()
	require.NoError(t, err, "Ping should not throw an error")

	stmt, err := db.Select(&Struct{})
	require.NoError(t, err, "Select function should not throw an error")

	stmt.Having(deebee.GreaterThan("id", 2))
	stmt.Having(deebee.EqualTo("active", true), deebee.LessThan("dateTimeUpdated", time.Now().String()))

	require.Equal(t, "SELECT id,name,active,dateTimeUpdated,dateTimeCreated FROM Struct HAVING (id > ?) OR (active = ? AND dateTimeUpdated < ?) ", stmt.String(), "Statement string should match expected value.")

	rows, err := db.Query(stmt)
	require.NoError(t, err, "Query function should not throw an error")
	require.Truef(t, len(rows) > 0, "Query method should return rows but length = %d", len(rows))

	fmt.Printf("ROW COUNT: %d \n", len(rows))
	fmt.Printf("QUERY: %s \n", stmt.String())
	fmt.Println()

}

func TestMethodOrderBy(t *testing.T) {

	type Struct struct {
		ID              int    `db:""`
		Name            string `db:""`
		Active          bool   `db:"NOT NULL"`
		DateTimeUpdated time.Time
		DateTimeCreated time.Time
	}

	db, err := New("mysql://root:toor@localhost:3306/db")
	require.NoError(t, err, "MySQL DB constructor should not throw an error")
	defer db.Close()

	err = db.Ping()
	require.NoError(t, err, "Ping should not throw an error")

	stmt, err := db.Select(&Struct{})
	require.NoError(t, err, "Select function should not throw an error")

	stmt.OrderBy("ASC", "dateTimeUpdated", "active")

	require.Equal(t, "SELECT id,name,active,dateTimeUpdated,dateTimeCreated FROM Struct ORDER BY dateTimeUpdated,active ASC ", stmt.String(), "Statement string should match expected value.")

	rows, err := db.Query(stmt)
	require.NoError(t, err, "Query function should not throw an error")
	require.Truef(t, len(rows) > 0, "Query method should return rows but length = %d", len(rows))

	fmt.Printf("ROW COUNT: %d \n", len(rows))
	fmt.Printf("QUERY: %s \n", stmt.String())
	fmt.Println()

}

func TestMethodLimit(t *testing.T) {

	type Struct struct {
		ID              int    `db:""`
		Name            string `db:""`
		Active          bool   `db:"NOT NULL"`
		DateTimeUpdated time.Time
		DateTimeCreated time.Time
	}

	db, err := New("mysql://root:toor@localhost:3306/db")
	require.NoError(t, err, "MySQL DB constructor should not throw an error")
	defer db.Close()

	err = db.Ping()
	require.NoError(t, err, "Ping should not throw an error")

	stmt, err := db.Select(&Struct{})
	require.NoError(t, err, "Select function should not throw an error")

	stmt.Limit(0, 3)

	require.Equal(t, "SELECT id,name,active,dateTimeUpdated,dateTimeCreated FROM Struct LIMIT 0, 3 ", stmt.String(), "Statement string should match expected value.")

	rows, err := db.Query(stmt)
	require.NoError(t, err, "Query function should not throw an error")
	require.Truef(t, len(rows) > 0, "Query method should return rows but length = %d", len(rows))

	fmt.Printf("ROW COUNT: %d \n", len(rows))
	fmt.Printf("QUERY: %s \n", stmt.String())
	fmt.Println()

}

func TestRowMethodPopulate(t *testing.T) {

	type Struct struct {
		ID              int    `db:""`
		Name            string `db:""`
		Active          bool   `db:"NOT NULL"`
		DateTimeUpdated time.Time
		DateTimeCreated time.Time
	}

	db, err := New("mysql://root:toor@localhost:3306/db")
	require.NoError(t, err, "MySQL DB constructor should not throw an error")
	defer db.Close()

	err = db.Ping()
	require.NoError(t, err, "Ping should not throw an error")

	stmt, err := db.Select(&Struct{})
	require.NoError(t, err, "Select function should not throw an error")

	rows, err := db.Query(stmt)
	require.NoError(t, err, "Query function should not throw an error")
	require.Truef(t, len(rows) > 0, "Query method should return rows but lenght = %d", len(rows))

	fmt.Printf("ROW COUNT: %d \n", len(rows))
	fmt.Printf("QUERY: %s \n", stmt.String())
	fmt.Println()

	for _, row := range rows {

		s := new(Struct)
		err = row.Populate(s)
		require.NoError(t, err, "Populate function should not throw an error")
		fmt.Printf("POPULATED STRUCT: %+v \n", s)

	}

}
