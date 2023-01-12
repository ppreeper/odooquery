package main

import "github.com/ppreeper/odoojrpc"

func odooConnect(host Host) (*odoojrpc.Odoo, error) {
	oc := odoojrpc.Odoo{
		Hostname: host.Hostname,
		Port:     host.Port,
		Database: host.Database,
		Username: host.Username,
		Password: host.Password,
		Schema:   host.Schema,
	}
	err := oc.Login()
	if err != nil {
		return nil, err
	}

	return &oc, nil
}
