package psqldb

import (
	"database/sql"
	"fmt"
	"os"

	_ "github.com/lib/pq"
)

var (
	host     = "localhost"
	port     = 5432
	user     = "postgres"
	password = os.Getenv("DB_PASSWORD")
	dbname   = "postgres"
)

func Connect() (*sql.DB, error) {
	psqlInfo := fmt.Sprintf("host=%s port=%d user=%s "+
		"password=%s dbname=%s sslmode=disable",
		host, port, user, password, dbname)
	var err error
	db, err := sql.Open("postgres", psqlInfo)
	if err != nil {
		return nil, err
	}

	err = db.Ping()
	if err != nil {
		return nil, err
	}

	fmt.Println("Successfully connected!")

	createSqlStatement := `CREATE TABLE IF NOT EXISTS PhoneNumber (
		id SERIAL PRIMARY KEY,
		name TEXT NOT NULL,
		number TEXT NOT NULL
	);`
	_, err = db.Exec(createSqlStatement)

	return db, err
}

func CloseConnection(db *sql.DB) {
	defer db.Close()
}

func InsertRecord(db *sql.DB, name string, phone string) error {
	sqlStatement := `INSERT INTO PhoneNumber (name,number)
	VALUES ($1,$2)`
	_, err := db.Exec(sqlStatement, name, phone)
	if err != nil {
		return err
	}
	return nil
}

func DeleteRecord(db *sql.DB, name string) {
	sqlStatement := `DELETE FROM PhoneNumber WHERE name=$1;`
	_, err := db.Exec(sqlStatement, name)
	if err != nil {
		panic(err)
	} else {
		fmt.Printf("Record with name=%s deleted\n", name)
	}
}

func Update(db *sql.DB, name string, number string) {
	sqlStatement := `UPDATE PhoneNumber SET number= $1 WHERE name=$2;`
	_, err := db.Exec(sqlStatement, number, name)
	if err != nil {
		panic(err)
	} else {
		fmt.Printf("Record with name=%s updated\n", name)
	}
}

func ListRecord(db *sql.DB) error {
	sqlStatement := `SELECT * FROM PhoneNumber;`
	rows, err := db.Query(sqlStatement)
	if err != nil {
		return err
	}
	defer rows.Close()
	for rows.Next() {
		var id int
		var name string
		var number string
		err = rows.Scan(&id, &name, &number)
		if err != nil {
			panic(err)
		}
		fmt.Printf("%d: %s, %s\n", id, name, number)
	}
	err = rows.Err()
	if err != nil {
		panic(err)
	}
	return nil
}
