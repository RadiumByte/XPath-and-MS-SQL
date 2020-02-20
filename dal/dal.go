package dal

import (
	"database/sql"

	"XPath-and-MS-SQL/app"
	"fmt"
	"log"

	_ "github.com/denisenkom/go-mssqldb"
)

// MsSQL represents data for connection to Data base
type MsSQL struct {
	Host     string
	DataBase *sql.DB
}

var (
	sqlversion string
)

// NewMsSQL constructs object of MsSQL
func NewMsSQL(host string, port int) (*MsSQL, error) {
	// Create connection string
	connString := fmt.Sprintf("server=%s;port=%d;trusted_connection=yes;",
		host, port)

	// Create connection pool
	db, err := sql.Open("sqlserver", connString)
	if err != nil {
		log.Fatal("Error creating connection pool: " + err.Error())
	}
	log.Printf("Connected!\n")

	rows, err := db.Query("select @@version")
	if err != nil {
		log.Fatal(err)
	}
	for rows.Next() {
		err := rows.Scan(&sqlversion)
		if err != nil {
			log.Fatal(err)
		}
		log.Println(sqlversion)
	}

	res := &MsSQL{
		Host:     host,
		DataBase: db}

	return res, nil
}

// Create inserts new Receipt into DB
func (t *MsSQL) Create(current *app.Receipt) error {
	/*
		var target Receipt
		target.Price = current.Price
		target.Post = current.Post

		if current.IsBankCard {
			target.IsBankCard = 1
		} else {
			target.IsBankCard = -1
		}

		target.IsProcessed = -1

		err := t.DataBase.Insert(&target)

		if err != nil {
			return nil, err
		}
	*/
	return nil
}
