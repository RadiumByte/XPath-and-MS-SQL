package dal

import (
	"database/sql"

	"github.com/RadiumByte/XPath-and-MS-SQL/app"
	"github.com/denisenkom/go-mssqldb"
)

// MsSQL represents data for connection to Data base
type MsSQL struct {
	Host     string
	DataBase *sql.DB
}

// NewMsSQL constructs object of MsSQL
func NewMsSQL(host string, port int) (*MsSQL, error) {
	// Create connection string
	connString := fmt.Sprintf("server=%s;port=%d;trusted_connection=yes;",
		host, port)

	// Create connection pool
	db, err = sql.Open("sqlserver", connString)
	if err != nil {
		log.Fatal("Error creating connection pool: " + err.Error())
	}
	log.Printf("Connected!\n")

	res := &MsSQL{
		Host:     host,
		DataBase: db}

	return res, nil
}
