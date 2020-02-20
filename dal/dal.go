package dal

import (
	"context"
	"database/sql"
	"errors"

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
	connString := fmt.Sprintf("server=%s;port=%d;trusted_connection=yes;", host, port)

	var err error
	var db *sql.DB

	db, err = sql.Open("sqlserver", connString)
	if err != nil {
		log.Fatal("Error creating connection pool: ", err.Error())
	}
	ctx := context.Background()
	err = db.PingContext(ctx)
	if err != nil {
		log.Fatal(err.Error())
	}
	fmt.Printf("Connected!\n")

	res := &MsSQL{
		Host:     host,
		DataBase: db}

	return res, nil
}

// Create inserts new Receipt into DB
func (t *MsSQL) Create(current *app.Receipt) error {
	ctx := context.Background()
	var err error

	if t.DataBase == nil {
		err = errors.New("DB is null")
		log.Println("Null DB")
		return err
	}

	// Check if database is alive.
	err = t.DataBase.PingContext(ctx)
	if err != nil {
		log.Println("Ping error")
		return err
	}

	tsql := `USE storage INSERT INTO dbo.receipts values (
	@PostNum, 
	@PostAddr, 
	@OFD, 
	@Price, 
	@Currency, 
	@IsBankCard, 
	@IsFiscal, 
	@IsService, 
	@Time
	)`

	stmt, err := t.DataBase.Prepare(tsql)
	if err != nil {
		log.Println("Prepare error")
		return err
	}
	defer stmt.Close()

	row := stmt.QueryRowContext(
		ctx,
		sql.Named("PostNum", current.PostNum),
		sql.Named("PostAddr", current.PostAddr),
		sql.Named("OFD", current.OFD),
		sql.Named("Price", current.Price),
		sql.Named("Currency", current.Currency),
		sql.Named("IsBankCard", current.IsBankCard),
		sql.Named("IsFiscal", current.IsFiscal),
		sql.Named("IsService", current.IsService),
		sql.Named("Time", current.OperationTime),
	)

	var newID int64
	err = row.Scan(&newID)
	if err != nil {
		return nil
	}

	return nil
}
