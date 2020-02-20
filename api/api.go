package api

import (
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/RadiumByte/XPath-and-MS-SQL/app"
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
	for i, data := range xmlquery.Find(message, "//item") {
		currentReceipt := app.NewReceipt()

		// Parse general package data:
		if postaddr := receipts.SelectElement("postaddr"); postaddr != nil {
			currentReceipt.PostAddr = postaddr.InnerText()
			// TODO: add logging here
		}

		if ofd := receipts.SelectElement("ofd"); ofd != nil {
			currentReceipt.OFD = ofd.InnerText()
			// TODO: add logging here
		}

		// Parse item-specific data:
		if postnum := data.SelectElement("postnum"); postnum != nil {
			numStr := postnum.InnerText()
			numInt, err := strconv.ParseInt(numStr, 10, 32)
			if err != nil {
				currentReceipt.PostNum = 0
			}
			currentReceipt.PostNum = numInt
			// TODO: add logging here
		}

		if price := data.SelectElement("price"); price != nil {
			priceStr := price.InnerText()
			priceInt, err := strconv.ParseInt(priceStr, 10, 32)
			if err != nil {
				currentReceipt.Price = 0
			}
			currentReceipt.Price = priceInt
			// TODO: add logging here
		}

		if currency := data.SelectElement("currency"); currency != nil {
			currencyStr := currency.InnerText()
			currencyInt, err := strconv.ParseInt(currencyStr, 10, 32)
			if err != nil {
				currentReceipt.Currency = 643
			}
			currentReceipt.Price = currencyInt
			// TODO: add logging here
		}

		if isbankcard := data.SelectElement("isbankcard"); isbankcard != nil {
			isbankcardStr := isbankcard.InnerText()

			isbankcardBool, err := strconv.ParseBool(isbankcardStr)
			if err != nil {
				currentReceipt.IsBankCard = false
			}
			currentReceipt.IsBankCard = isbankcardBool
			// TODO: add logging here
		}

		if isfiscal := data.SelectElement("isfiscal"); isfiscal != nil {
			isfiscalStr := isfiscal.InnerText()

			isfiscalBool, err := strconv.ParseBool(isfiscalStr)
			if err != nil {
				currentReceipt.IsFiscal = false
			}
			currentReceipt.IsFiscal = isfiscalBool
			// TODO: add logging here
		}

		if isservice := data.SelectElement("isservice"); isservice != nil {
			isserviceStr := isservice.InnerText()

			isserviceBool, err := strconv.ParseBool(isserviceStr)
			if err != nil {
				currentReceipt.IsService = true
			}
			currentReceipt.IsService = isserviceBool
			// TODO: add logging here
		}

		if datetime := data.SelectElement("time"); datetime != nil {
			datetimeStr := datetime.InnerText()

			datetimeObject, err := time.Parse(time.RFC3339, datetimeStr)
			if err != nil {
				currentReceipt.OperationTime = time.Now()
			}

			currentReceipt.OperationTime = datetimeObject
			// TODO: add logging here
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

	log.Info("Server is starting on port", port)
	errc <- fasthttp.ListenAndServe(port, router.Handler)
}

// NewWebServer constructs Web Server
func NewWebServer(application app.IncomeRegistration) *WebServer {
	res := &WebServer{}
	res.application = application

	return res
}
