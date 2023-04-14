package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/ppreeper/odoojrpc"
)

type application struct {
	errorLog *log.Logger
	infoLog  *log.Logger
	oc       *odoojrpc.Odoo
}

func main() {
	app := &application{
		errorLog: log.New(os.Stderr, "ERROR\t", log.Ldate|log.Ltime|log.Lshortfile),
		infoLog:  log.New(os.Stdout, "INFO\t", log.Ldate|log.Ltime),
	}

	// Config File
	userConfigDir, err := os.UserConfigDir()
	checkErr(err)

	var configFile, system, database string
	flag.StringVar(&configFile, "c", userConfigDir+"/odooquery/config.yml", "odoo connection config file")
	flag.StringVar(&system, "system", "prod", "odoo host specified in config.yml")
	flag.StringVar(&database, "d", "", "odoo database")

	var q QueryDef
	flag.StringVar(&q.Model, "model", "", "model")
	flag.StringVar(&q.Filter, "filter", "", "filter")
	flag.IntVar(&q.Offset, "offset", 0, "offset")
	flag.IntVar(&q.Limit, "limit", 0, "limit")
	flag.StringVar(&q.Fields, "fields", "", "fields")
	flag.BoolVar(&q.Count, "count", false, "count records")

	flag.Parse()

	// get config file
	HostMap := GetConf(configFile)

	if _, ok := HostMap[system]; !ok {
		fmt.Println("no host specified")
		return
	}
	server := HostMap[system]

	if server.Database == "" && database == "" {
		fmt.Println("no database specified")
		return
	}
	if database != "" {
		server.Database = database
	}

	if q.Model == "" {
		fmt.Println("no model specified")
		return
	}

	oc, err := odooConnect(server)
	if err != nil {
		app.errorLog.Println("error:", err)
		fatalErr(err)
	}
	app.oc = oc

	app.getRecords(&q)
}

func (app *application) getRecords(q *QueryDef) {
	umdl := strings.Replace(q.Model, "_", ".", -1)

	fields := parseFields(q.Fields)
	if q.Count {
		fields = []string{"id"}
	}

	filtp, err := parseFilter(q.Filter)
	checkErr(err)

	rr, err := app.oc.SearchRead(umdl, filtp, q.Offset, q.Limit, fields)
	fatalErr(err)
	if q.Count {
		fmt.Println("records:", len(rr))
	} else {
		jsonStr, err := json.MarshalIndent(rr, "", "  ")
		checkErr(err)
		fmt.Println(string(jsonStr))
	}
}
