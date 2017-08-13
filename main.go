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
	UseRootDisk  bool
	SocketPath   string
)

func init() {
	log.SetPrefix(os.Args[0] + " | ")

	flag.BoolVar(&PrintVersion, "v", false, "Show the plugin version information")
	flag.StringVar(&SocketPath, "s", "/run/docker/plugins", "Path to the plugin socket")
	flag.BoolVar(&UseRootDisk, "r", false, "Use the root path / or /boot or /var/lib/docker to produce volumes")
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
		fmt.Printf("ZaneCloud Stone Docker Volume plugin: %s\n", VERSION)
		return
	}

	plugin := stone_plugin.NewPluginAPI(SocketPath, UseRootDisk)

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
