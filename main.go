package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"os"
	"strings"

	"github.com/ppreeper/odoojrpc"
)

func main() {
	// Config File
	userConfigDir, err := os.UserConfigDir()
	checkErr(err)
	// Hoster()

	HostMap := GetConf(userConfigDir + "/odooquery/config.yml")

	var o odoojrpc.Odoo
	var q QueryDef

	var host string
	var database string

	flag.StringVar(&host, "host", "localhost", "odoo host")
	flag.StringVar(&database, "d", "", "odoo database")

	flag.StringVar(&q.Model, "model", "", "model")
	flag.StringVar(&q.Filter, "filter", "", "filter")
	flag.IntVar(&q.Offset, "offset", 0, "offset")
	flag.IntVar(&q.Limit, "limit", 0, "limit")
	flag.StringVar(&q.Fields, "fields", "", "fields")
	flag.BoolVar(&q.Count, "count", false, "count records")

	flag.Parse()
	if _, ok := HostMap[host]; !ok {
		fmt.Println("no host specified")
		return
	}
	server := HostMap[host]

	fmt.Println("database", database)
	fmt.Println("server.Database", server.Database)

	if server.Database == "" && database == "" {
		fmt.Println("no database specified")
		return
	}
	if database == "" {
		database = server.Database
	}

	if q.Model == "" {
		fmt.Println("no model specified")
		return
	}

	o.Hostname = server.Hostname
	o.Port = server.Port
	o.Schema = server.Schema
	o.Username = server.Username
	o.Password = server.Password
	o.Database = database

	err = o.Login()
	checkErr(err)
	getRecords(&o, &q)
}

func getRecords(o *odoojrpc.Odoo, q *QueryDef) {
	umdl := strings.Replace(q.Model, "_", ".", -1)

	fields := parseFields(q.Fields)
	if q.Count {
		fields = []string{"id"}
	}

	filtp, err := parseFilter(q.Filter)
	checkErr(err)

	rr := o.SearchRead(umdl, filtp, q.Offset, q.Limit, fields)
	if q.Count {
		fmt.Println("records:", len(rr))
	} else {
		jsonStr, err := json.MarshalIndent(rr, "", "  ")
		checkErr(err)
		fmt.Println(string(jsonStr))
	}
}
