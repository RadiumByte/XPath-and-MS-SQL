package api

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/RadiumByte/XPath-and-MS-SQL/cmd/web/app"
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
	log.Info("API got new receipt...")

	currentReceipt := app.NewReceipt()

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

	if postaddr := channel.SelectElement("postaddr"); n != nil {
		currentReceipt.
			fmt.Printf(": %s\n", n.InnerText())
	}

	if n := channel.SelectElement("link"); n != nil {
		fmt.Printf("link: %s\n", n.InnerText())
	}
	for i, n := range xmlquery.Find(doc, "//item/title") {
		fmt.Printf("#%d %s\n", i, n.InnerText())
	}

	// TODO: parse payload with XPath

	// TODO: write data to the currentReceipt model

	log.Info("API got new receipt...")

	server.application.RegisterReceipt(currentReceipt)

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
