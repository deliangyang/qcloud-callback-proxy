package main

import (
	"flag"
	"net/http"

	"github.com/deliangyang/qcloud-callback-proxy/internal"
)

func main() {

	configFile := flag.String("config", "configs/proxy.toml", "configuration file")
	flag.Parse()

	if err := internal.Parse(*configFile); err != nil {
		internal.Logger().WithError(err).Fatal("parse config fail")
	}
	conf := internal.GetConfig()

	http.HandleFunc(conf.URI, internal.Dispatch)
	if err := http.ListenAndServe(conf.Port, nil); err != nil {
		internal.Logger().WithError(err).Fatal("create server fail")
	}

}
