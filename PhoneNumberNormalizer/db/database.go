package phonedb

import (
	"database/sql"
	"fmt"

	_ "github.com/lib/pq"
)

type DB struct {
	db *sql.DB
}

type Phone struct {
	Id     int
	Number string
}

func Open(driverName, dataSource string) (*DB, error) {

	db, err := sql.Open(driverName, dataSource)
	if err != nil {
		return nil, err
	}

	return &DB{db}, nil
}

func (db *DB) CloseConnection() error {
	return db.db.Close()
}

func (db *DB) Seed() error {
	data := []string{
		"1234567890",
		"123 456 7891",
		"(123) 456 7892",
		"(123) 456-7893",
		"123-456-7894",
		"123-456-7890",
		"1234567892",
		"(123)456-7892",
	}
	for _, number := range data {
		if _, err := insertRecord(db, number); err != nil {
			return err
		}
	}
	return nil
}

func insertRecord(db *DB, number string) (int, error) {
	sqlStatement := `INSERT INTO PhoneNumber (number) VALUES ($1) RETURNING id`
	var id int
	err := db.db.QueryRow(sqlStatement, number).Scan(&id)
	if err != nil {
		return -1, err
	}
	return id, nil
}

func (db *DB) DeleteRecord(id int) error {
	sqlStatement := `DELETE FROM PhoneNumber WHERE id=$1`
	_, err := db.db.Exec(sqlStatement, id)
	if err != nil {
		return err
	} else {
		fmt.Printf("Record with id=%d deleted\n", id)
		return nil
	}
}

func (db *DB) UpdateRecord(phone *Phone) error {
	sqlStatement := `UPDATE PhoneNumber SET number= $1 WHERE id=$2;`
	_, err := db.db.Exec(sqlStatement, phone.Number, phone.Id)
	if err != nil {
		return err
	} else {
		fmt.Printf("Record with id=%d updated\n", phone.Id)
		return nil
	}
}

func (db *DB) ListRecord() ([]Phone, error) {
	sqlStatement := `SELECT * FROM PhoneNumber;`
	var phones []Phone
	rows, err := db.db.Query(sqlStatement)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	for rows.Next() {
		var id int
		var number string
		err = rows.Scan(&id, &number)
		if err != nil {
			return nil, err
		}
		phones = append(phones, Phone{id, number})
	}
	err = rows.Err()
	if err != nil {
		return nil, err
	}
	return phones, nil
}

func (db *DB) FindPhone(number string) (*Phone, error) {
	var p Phone
	row := db.db.QueryRow("SELECT * FROM PhoneNumber WHERE number=$1", number)

	err := row.Scan(&p.Id, &p.Number)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		} else {
			return nil, err
		}
	}
	return &p, nil
}

func Migrate(driverName, dataSource string) error {
	db, err := sql.Open(driverName, dataSource)
	if err != nil {
		return err
	}
	err = createPhoneNumberTable(db)
	if err != nil {
		return err
	}
	return db.Close()
}

func createPhoneNumberTable(db *sql.DB) error {
	createSqlStatement := `CREATE TABLE IF NOT EXISTS PhoneNumber (
		id SERIAL,
		number VARCHAR(255)
	);`
	_, err := db.Exec(createSqlStatement)
	return err
}

func Reset(driverName, dataSource, dbName string) error {
	db, err := sql.Open(driverName, dataSource)
	if err != nil {
		return err
	}
	err = resetDB(db, dbName)
	if err != nil {
		return err
	}
	return db.Close()
}

func resetDB(db *sql.DB, name string) error {
	_, err := db.Exec("DROP DATABASE IF EXISTS " + name)
	if err != nil {
		return err
	}
	return createDB(db, name)
}
func createDB(db *sql.DB, name string) error {
	_, err := db.Exec("CREATE DATABASE " + name)
	if err != nil {
		return err
	}
	return nil
}

func getPhone(db *sql.DB, id int) (string, error) {
	var number string
	row := db.QueryRow("SELECT * FROM PhoneNumber WHERE id=$1", id)
	err := row.Scan(&id, &number)
	if err != nil {
		return "", err
	}
	return number, nil
}
