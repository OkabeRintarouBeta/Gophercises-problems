package main

import (
	"database/sql"
	"encoding/csv"
	"flag"
	"os"
	"unicode"

	"github.com/okaberintaroubeta/phoneNumberNormalizer/psqldb"
)

type Identity struct {
	Name  string
	Phone string
}

func normalize(phone string) string {
	var ans string
	for _, ch := range phone {
		if unicode.IsDigit(ch) {
			ans = ans + string(ch)
		}
	}
	return ans
}

func main() {
	path := flag.String("path", "data.csv", "Path to the phone book")
	flag.Parse()

	db, err := psqldb.Connect()
	if err != nil {
		panic(err)
	}
	addEntries(*path, db)

	psqldb.ListRecord(db)

	psqldb.Update(db, "John", normalize("(1582341)234"))
	psqldb.DeleteRecord(db, "Sammie")
	psqldb.ListRecord(db)
	psqldb.CloseConnection(db)
}

func addEntries(path string, db *sql.DB) {
	var identities []Identity
	file, err := os.Open(path)
	if err != nil {
		panic(err)
	}
	reader := csv.NewReader(file)
	data, err := reader.ReadAll()
	if err != nil {
		panic(err)
	}

	for _, row := range data {
		if len(row) < 2 {
			continue
		}
		identities = append(identities, Identity{row[0], normalize(row[1])})
	}

	for _, identity := range identities {
		err = psqldb.InsertRecord(db, identity.Name, identity.Phone)
		if err != nil {
			panic(err)
		}
	}
}
