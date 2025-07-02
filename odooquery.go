package main

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"strings"

	"github.com/charmbracelet/fang"
	"github.com/ppreeper/odoorpc"
	"github.com/ppreeper/odoorpc/odoojrpc"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

type QueryDef struct {
	Model  string
	Filter string
	Offset int
	Limit  int
	Fields string
	Count  bool
}

// Conf config structure
type Host struct {
	Hostname string `default:"localhost" json:"hostname"`
	Database string `default:"odoo" json:"database"`
	Username string `default:"odoo" json:"username"`
	Password string `default:"odoo" json:"password"`
	Protocol string `default:"jsonrpc" json:"protocol,omitempty"`
	Schema   string `default:"https" json:"schema,omitempty"`
	Port     int    `default:"443" json:"port,omitempty"`
}

func NewHost() *Host {
	return &Host{
		Hostname: "localhost",
		Database: "odoo",
		Username: "odoo",
		Password: "odoo",
		Protocol: "jsonrpc",
		Schema:   "https",
		Port:     443,
	}
}

func getServerConfigByName(name string) (*Host, error) {
	if name == "" {
		return nil, fmt.Errorf("server name cannot be empty")
	}

	server := NewHost()
	server.Hostname = viper.GetString(name + ".hostname")
	server.Database = viper.GetString(name + ".database")
	server.Username = viper.GetString(name + ".username")
	server.Password = viper.GetString(name + ".password")
	protocol := viper.GetString(name + ".protocol")
	if protocol != "" {
		server.Protocol = protocol
	}
	schema := viper.GetString(name + ".schema")
	if schema != "" {
		server.Schema = schema
	}
	port := viper.GetInt(name + ".port")
	if port == 0 {
		if server.Schema == "https" {
			server.Port = 443
		} else {
			server.Port = 8069
		}
	} else {
		server.Port = port
	}

	return server, nil
}

func getOdooClient(server *Host) odoorpc.Odoo {
	odooClient := odoojrpc.NewOdoo().
		WithHostname(server.Hostname).
		WithPort(server.Port).
		WithSchema(server.Schema).
		WithDatabase(server.Database).
		WithUsername(server.Username).
		WithPassword(server.Password)
	return odooClient
}

func main() {
	// Config File
	userConfigDir, err := os.UserConfigDir()
	checkErr(err)
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath(userConfigDir + "/odooquery")
	if err := viper.ReadInConfig(); err != nil {
		fmt.Println("Error reading config file:", err)
		return
	}

	rootCmd := &cobra.Command{
		Use:  "odooquery <system> <model>",
		Args: cobra.ExactArgs(2),
		Run: func(cmd *cobra.Command, args []string) {
			system := args[0]
			model := args[1]

			// Extract flags
			limit, _ := cmd.Flags().GetInt("limit")
			offset, _ := cmd.Flags().GetInt("offset")
			fields, _ := cmd.Flags().GetString("fields")
			filter, _ := cmd.Flags().GetString("filter")
			count, _ := cmd.Flags().GetBool("count")

			q := QueryDef{
				Model:  model,
				Offset: offset,
				Limit:  limit,
				Fields: fields,
				Filter: filter,
				Count:  count,
			}

			server, err := getServerConfigByName(system)
			fatalErr(err)
			oc := getOdooClient(server)
			err = oc.Login()
			fatalErr(err, "login failed")

			getRecords(oc, q)
		},
	}
	rootCmd.Flags().IntP("offset", "o", 0, "offset records, 0 for no offset")
	rootCmd.Flags().IntP("limit", "l", 0, "limit records, 0 for no limit")
	rootCmd.Flags().StringP("fields", "f", "", "fields field1,field2,...fieldN")
	rootCmd.Flags().StringP("filter", "F", "", "filter \"[('field', 'op', value), ...]\"")
	rootCmd.Flags().BoolP("count", "c", false, "count records")

	if err := fang.Execute(context.Background(), rootCmd); err != nil {
		fmt.Println("Error executing command:", err)
		os.Exit(1)
	}
}

func getRecords(oc odoorpc.Odoo, q QueryDef) {
	umdl := strings.ReplaceAll(q.Model, "_", ".")

	filtp, err := parseFilter(q.Filter)
	checkErr(err)

	if q.Count {
		count, err := oc.Count(umdl, filtp)
		fatalErr(err, "query error", "check if model exists")
		fmt.Println("records:", count)
		return
	}

	fields := parseFields(q.Fields)

	rr, err := oc.SearchRead(umdl, q.Offset, q.Limit, fields, filtp)
	fatalErr(err, "search read error")

	jsonStr, err := json.MarshalIndent(rr, "", "  ")
	fatalErr(err, "record marshalling error")
	fmt.Println(string(jsonStr))
}

func checkErr(err error, msg ...string) {
	if err != nil {
		if len(msg) == 0 {
			fmt.Printf("error: %v\n", err.Error())
		}
		fmt.Printf("%v\n", strings.Join(msg, " "))
	}
}

func fatalErr(err error, msg ...string) {
	if err != nil {
		if len(msg) == 0 {
			fmt.Printf("error: %v\n", err.Error())
		}
		fmt.Printf("%v\n", strings.Join(msg, " "))
		os.Exit(1)
	}
}
