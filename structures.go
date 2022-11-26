package main

import (
	"os"

	"gopkg.in/yaml.v3"
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
	Database string `default:"odoo" json:"database,omitempty"`
	Username string `default:"odoo" json:"username"`
	Password string `default:"odoo" json:"password"`
	Protocol string `default:"jsonrpc" json:"protocol,omitempty"`
	Schema   string `default:"http" json:"schema,omitempty"`
	Port     int    `default:"8069" json:"port,omitempty"`
	Jobcount int    `default:"1" json:"jobcount,omitempty"`
}

func GetConf(configFile string) map[string]Host {
	yamlFile, err := os.ReadFile(configFile)
	checkErr(err)
	data := make(map[string]Host)
	err = yaml.Unmarshal(yamlFile, data)
	checkErr(err)
	return data
}
