package main

import (
	"fmt"
	"github.com/antoligy/offthegrid"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"os"
	"time"
)

func main() {
	var err error

	if len(os.Args) < 2 {
		fmt.Fprintf(os.Stderr, "usage: %s <path/to/config.yaml>\n", os.Args[0])
		os.Exit(1)
	}

	cfgdata, err := ioutil.ReadFile(os.Args[1])
	if err != nil {
		fmt.Fprintf(os.Stderr, "error: can't read config file: %s\n", err.Error())
		os.Exit(2)
	}

	cfg := offthegrid.Config{}
	err = yaml.Unmarshal(cfgdata, &cfg)
	if err != nil {
		fmt.Fprintf(os.Stderr, "error: can't read config file: %s\n", err.Error())
		os.Exit(2)
	}

	server, err := offthegrid.NewServer(cfg)
	if err != nil {
		fmt.Fprintf(os.Stderr, "error: can't create server: %s\n", err.Error())
		os.Exit(3)
	}

	fmt.Fprintf(os.Stdout, "offthegrid/%d starting up\n", offthegrid.VERSION)
	cstarted, cerr := server.Run()
	for {
		select {
		case _ = <-cstarted:
			fmt.Fprintf(os.Stdout, "listening on '%s'\n", cfg.ListenSocket)
		case err := <-cerr:
			fmt.Fprintf(os.Stderr, "error: %s\n", err.Error())
		case _ = <-time.After(100 * time.Millisecond):
		}
	}
}
