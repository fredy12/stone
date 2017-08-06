// Copyright (c) 2017-2018, ZANECLOUD CORPORATION. All rights reserved.

package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"runtime"
	"runtime/debug"

	"github.com/zanecloud/stone/stone_plugin"
)

var (
	PrintVersion bool
	ListenAddr   string
	VolumesPath  string
	SocketPath   string

	Version string
)

func init() {
	log.SetPrefix(os.Args[0] + " | ")

	flag.BoolVar(&PrintVersion, "v", false, "Show the plugin version information")
	flag.StringVar(&ListenAddr, "l", "localhost:3476", "Server listen address")
	flag.StringVar(&VolumesPath, "d", "/var/lib/stone/volumes", "Path where to store the volumes")
	flag.StringVar(&SocketPath, "s", "/run/docker/plugins", "Path to the plugin socket")
}

func assert(err error) {
	if err != nil {
		log.Panicln("Error:", err)
	}
}

func exit() {
	if err := recover(); err != nil {
		if _, ok := err.(runtime.Error); ok {
			log.Println(err)
		}
		if os.Getenv("STONE_DEBUG") != "" {
			log.Printf("%s", debug.Stack())
		}
		os.Exit(1)
	}
	os.Exit(0)
}

func main() {
	flag.Parse()
	defer exit()

	if PrintVersion {
		fmt.Printf("ZaneCloud Stone Docker Volume plugin: %s\n", Version)
		return
	}

	plugin := stone_plugin.NewPluginAPI(SocketPath)

	log.Println("Serving plugin API at", SocketPath)
	p := plugin.Serve()
L:
	for {
		select {
		case <-p:
			break L
		}
	}
	assert(plugin.Error())
	log.Println("Successfully terminated")
}
