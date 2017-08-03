package main

import (
	"github.com/kelseyhightower/confd/builtin/databases/consul"
	"github.com/kelseyhightower/confd/plugin"
)

func main() {
	plugin.Serve(&plugin.ServeOpts{
		DatabaseFunc: consul.Database,
	})
}