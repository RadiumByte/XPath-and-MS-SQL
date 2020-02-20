package api

import (
	"strings"

	"XPath-and-MS-SQL/app"

	"github.com/antchfx/xmlquery"
	"github.com/buaazp/fasthttprouter"
	"github.com/powerman/structlog"
	"github.com/valyala/fasthttp"
)

var log = structlog.New()

// WebServer accepts POST requests with payload of XML docs of Receipts
// Then it parses them with XPath and pushes data to Application
type WebServer struct {
	application app.IncomeRegistration
}

// ParseXML uses XPath to parse Receipt data and pushes it to Application
func (server *WebServer) ParseXML(ctx *fasthttp.RequestCtx) {
	log.Info("API got new receipt. Parsing started...")

	payLoad := string(ctx.PostBody())

	// Create parse tree
	message, err := xmlquery.Parse(strings.NewReader(payLoad))
	if err != nil {
		// Not a XML inside payload
		ctx.SetStatusCode(fasthttp.StatusBadRequest)
		return
	}

	// Example of XML receipt
	/*
		`<?xml version="1.0" encoding="UTF-8" ?>
		<receipts>
		  <postaddr>Kosmonavtov 21A</postaddr>
		  <ofd>ofd.yandex.ru</ofd>
		  <item>
		    <postnum>4</postnum>
			<price>245</price>
			<currency>643</currency>
			<isbankcard>0</isbankcard>
			<isfiscal>1</isfiscal>
			<isservice>0</isservice>
			<time>2020-02-20T15:10:15.371Z</time>
		  </item>
		  <item>
		    <postnum>2</postnum>
			<price>100</price>
			<currency>643</currency>
			<isbankcard>1</isbankcard>
			<isfiscal>1</isfiscal>
			<isservice>0</isservice>
			<time>2020-02-20T17:20:5.123Z</time>
		  </item>
		</receipts>`
	*/

	// Select receipts block
	receipts := xmlquery.FindOne(message, "//receipts")

	// Find all receipts inside payload
	for _, data := range xmlquery.Find(message, "//item") {
		currentReceipt := app.NewReceipt()

		// Parse general package data:
		if postaddr := receipts.SelectElement("postaddr"); postaddr != nil {
			currentReceipt.PostAddr = postaddr.InnerText()
		}

		if ofd := receipts.SelectElement("ofd"); ofd != nil {
			currentReceipt.OFD = ofd.InnerText()
		}

		// Parse item-specific data:
		if postnum := data.SelectElement("postnum"); postnum != nil {
			currentReceipt.PostNum = postnum.InnerText()
		}

		if price := data.SelectElement("price"); price != nil {
			currentReceipt.Price = price.InnerText()
		}

		if currency := data.SelectElement("currency"); currency != nil {
			currentReceipt.Currency = currency.InnerText()
		}

		if isbankcard := data.SelectElement("isbankcard"); isbankcard != nil {
			currentReceipt.IsBankCard = isbankcard.InnerText()
		}

		if isfiscal := data.SelectElement("isfiscal"); isfiscal != nil {
			currentReceipt.IsFiscal = isfiscal.InnerText()
		}

		if isservice := data.SelectElement("isservice"); isservice != nil {
			currentReceipt.IsService = isservice.InnerText()
		}

		if datetime := data.SelectElement("time"); datetime != nil {
			currentReceipt.OperationTime = datetime.InnerText()
		}

		// Send receipt to Application
		server.application.RegisterReceipt(currentReceipt)
	}

	ctx.SetStatusCode(fasthttp.StatusOK)
}

// Start initializes Web Server, starts application and begins serving
func (server *WebServer) Start(errc chan<- error) {
	router := fasthttprouter.New()
	router.PUT("/receipt_xml", server.ParseXML)

	port := ":6633"

	/*
		currentReceipt := app.NewReceipt()
		currentReceipt.Currency = "643"
		currentReceipt.IsBankCard = "0"
		currentReceipt.IsFiscal = "1"
		currentReceipt.IsService = "0"
		currentReceipt.OFD = "yandex.ofd.ru"
		currentReceipt.PostNum = "28"
		currentReceipt.PostAddr = "Rostov"
		currentReceipt.Price = "754"
		currentReceipt.OperationTime = "20120618 10:34:09"
		server.application.RegisterReceipt(currentReceipt)
	*/

	log.Info("Server is starting on port", port)
	errc <- fasthttp.ListenAndServe(port, router.Handler)
}

// NewWebServer constructs Web Server
func NewWebServer(application app.IncomeRegistration) *WebServer {
	res := &WebServer{}
	res.application = application

	return res
}
