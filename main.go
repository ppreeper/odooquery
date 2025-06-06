package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/ppreeper/odoorpc"
	"github.com/ppreeper/odoorpc/odoojrpc"
)

type application struct {
	errorLog *log.Logger
	infoLog  *log.Logger
	oc       odoorpc.Odoo
}

func main() {
	app := &application{
		errorLog: log.New(os.Stderr, "ERROR\t", log.Ldate|log.Ltime|log.Lshortfile),
		infoLog:  log.New(os.Stdout, "INFO\t", log.Ldate|log.Ltime),
	}

	// Config File
	userConfigDir, err := os.UserConfigDir()
	checkErr(err)

	var configFile, system string
	flag.StringVar(&configFile, "c", userConfigDir+"/odooquery/config.yml", "odoo connection config file")
	flag.StringVar(&system, "system", "prod", "odoo host specified in config.yml")

	var cServer Host
	flag.StringVar(&cServer.Hostname, "host", "", "odoo host")
	flag.StringVar(&cServer.Database, "d", "", "odoo database")
	flag.StringVar(&cServer.Username, "U", "", "odoo username")
	flag.StringVar(&cServer.Password, "P", "", "odoo password")
	cServer.Protocol = "jsonrpc"
	flag.StringVar(&cServer.Schema, "S", "http", "odoo url schema [http|https]")
	flag.IntVar(&cServer.Port, "port", 8069, "odoo port")
	cServer.Jobcount = 1

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

	// Server connection profile
	server := HostMap[system]
	if server.Hostname == "" && cServer.Hostname == "" {
		fmt.Println("no host specified")
		return
	}
	if cServer.Hostname != "" {
		server.Hostname = cServer.Hostname
	}

	// 	Database
	if server.Database == "" && cServer.Database == "" {
		fmt.Println("no database specified")
		return
	}
	if cServer.Database != "" {
		server.Database = cServer.Database
	}

	// 	Username
	if cServer.Username != "" {
		server.Username = cServer.Username
	}
	// 	Password
	if cServer.Password != "" {
		server.Password = cServer.Password
	}
	// 	Protocol
	if server.Protocol == "" && cServer.Protocol != "" {
		server.Protocol = cServer.Protocol
	}
	// 	Schema
	if server.Schema == "" && cServer.Schema != "" {
		server.Schema = cServer.Schema
	}
	// 	Port
	if server.Port == 0 {
		server.Port = cServer.Port
	}
	// 	Jobcount
	if server.Jobcount == 0 {
		server.Jobcount = cServer.Jobcount
	}

	if q.Model == "" {
		fmt.Println("no model specified")
		return
	}

	app.oc = odoojrpc.NewOdoo().
		WithHostname(server.Hostname).
		WithPort(server.Port).
		WithDatabase(server.Database).
		WithUsername(server.Username).
		WithPassword(server.Password).
		WithSchema(server.Schema)

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

	rr, err := app.oc.SearchRead(umdl, q.Offset, q.Limit, fields, filtp)
	fatalErr(err)
	if q.Count {
		fmt.Println("records:", len(rr))
	} else {
		jsonStr, err := json.MarshalIndent(rr, "", "  ")
		checkErr(err)
		fmt.Println(string(jsonStr))
	}
}
