package app

import (
	"github.com/powerman/structlog"
)

var log = structlog.New()

// IncomeRegistration is an interface for accepting income Receipts from Web Server
type IncomeRegistration interface {
	RegisterReceipt(*Receipt)
}

// DataAccessLayer is an interface for DAL usage from Application
type DataAccessLayer interface {
	Create(*Receipt) error
}

// Application is responsible for all logics and communicates with other layers
type Application struct {
	DB   DataAccessLayer
	errc chan<- error
}

// RegisterReceipt sends Receipt to DAL for saving/registration
func (app *Application) RegisterReceipt(currentData *Receipt) {
	err := app.DB.Create(currentData)

	if err != nil {
		app.errc <- err
		return
	}

	log.Info("New receipt added to MS SQL DBMS...")
}

// NewApplication constructs Application
func NewApplication(db DataAccessLayer, errchannel chan<- error) *Application {
	res := &Application{}

	res.DB = db
	res.errc = errchannel

	return res
}
