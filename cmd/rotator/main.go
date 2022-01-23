package main

import (
	"flag"
)

var configFile string

func init() {
	flag.StringVar(&configFile, "config", "configs/config.dev.yaml", "Path to configuration file")
}

func main() {
	flag.Parse()
}
