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
	db, err := sql.Open("mssql", connString)
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
	var query_str = "USE storage INSERT INTO dbo.receipts values ("
	query_str += current.PostNum + ", "
	query_str += "'" + current.PostAddr + "', "
	query_str += "'" + current.OFD + "', "
	query_str += current.Price + ", "
	query_str += current.Currency + ", "
	query_str += current.IsBankCard + ", "
	query_str += current.IsFiscal + ", "
	query_str += current.IsService + ", "

	query_str += "'" + current.OperationTime + "')"

	log.Println(query_str)

	_, err := t.DataBase.Query(query_str)
	if err != nil {
		log.Println("Query problem")
		return err
	}

	return nil
}
