package main

import (
	"fmt"
	"os"
	"unicode"

	_ "github.com/lib/pq"
	phonedb "github.com/okaberintaroubeta/phoneNumberNormalizer/db"
)

var (
	host     = "localhost"
	port     = 5432
	user     = "postgres"
	password = os.Getenv("DB_PASSWORD")
	dbname   = "phonedb"
)

func main() {

	psqlInfo := fmt.Sprintf("host=%s port=%d user=%s "+
		"password=%s dbname=%s sslmode=disable",
		host, port, user, password, dbname)
	must(phonedb.Reset("postgres", psqlInfo, dbname))

	psqlInfo = fmt.Sprintf("%s dbname=%s", psqlInfo, dbname)

	must(phonedb.Migrate("postgres", psqlInfo))

	db, err := phonedb.Open("postgres", psqlInfo)
	must(err)
	defer db.CloseConnection()

	err = db.Seed()
	must(err)

	phones, err := db.ListRecord()
	must(err)
	for _, p := range phones {
		fmt.Printf("Working on... %+v\n", p)
		number := normalize(p.Number)
		if number != p.Number {
			fmt.Println("Updating or removing...", number)
			existing, err := db.FindPhone(number)
			must(err)
			if existing != nil {
				must(db.DeleteRecord(p.Id))
			} else {
				p.Number = number
				must(db.UpdateRecord(&p))
			}
		} else {
			fmt.Println("No changes required")
		}
	}

}

func must(err error) {
	if err != nil {
		panic(err)
	}
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
