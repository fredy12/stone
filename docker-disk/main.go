package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"github.com/Sirupsen/logrus"
	"log"
	"os"
	"runtime"
)

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
		os.Exit(1)
	}
	os.Exit(0)
}

func getVersion() string {
	return "v1.0.0"
}

func usage() {
	fmt.Println("USAGE:")
	fmt.Println("    docker-disk [global options] command [command options] [arguments...]\n")
	fmt.Println("VERSION:")
	fmt.Printf("    %s\n\n", getVersion())
	fmt.Println("COMMANDS:")
	fmt.Println("    list     List all disks")
	fmt.Println("GLOBAL OPTIONS:")
	fmt.Println("   --help, -h     Show help")
}

func subCmdUsage(name string, cmdName string, options []map[string]string) {
	fmt.Println("USAGE:")
	if len(options) > 0 {
		fmt.Printf("    docker-disk [options] %s %s\n\n", name, cmdName)
		fmt.Println("OPTIONS:")

	} else {
		fmt.Printf("    docker-disk %s %s\n\n", name, cmdName)
	}
	for _, option := range options {
		for optName, optDesp := range option {
			fmt.Printf("    %s    %s\n", optName, optDesp)
		}
	}
}

func toJSON(v interface{}) {
	s, err := json.Marshal(v)
	assert(err)
	fmt.Printf(string(s))
}

func Run() {
	defer exit()

	flag.Usage = usage
	flag.Parse()

	if len(os.Args) < 3 {
		usage()
		os.Exit(2)
	}

	log.SetPrefix(os.Args[0] + os.Args[1] + " | ")
	switch os.Args[1] {
	case "list":
		f1 := flag.NewFlagSet(os.Args[1], flag.ExitOnError)
		f1.Usage = func() {
			subCmdUsage("list", "", []map[string]string{
				{"fstype": "Show which fstype disks: all/xfs/ext4"},
				{"rootdisk": "Show root disks"},
				{"verbose": "Show verbose logs"},
			})
		}

		useRootDisk := false
		f1.Parse(os.Args[2:])
		logrus.SetLevel(logrus.PanicLevel)
		fsType := os.Args[2]
		if fsType != "xfs" && fsType != "ext4" {
			fsType = ""
		}
		for _, arg := range os.Args[2:] {
			if arg == "verbose" {
				logrus.SetLevel(logrus.DebugLevel)
			}
			if arg == "rootdisk" {
				useRootDisk = true
			}
		}

		diskInfo, err := Collect(fsType, useRootDisk)
		assert(err)
		toJSON(diskInfo)
	}
}

func main() {
	Run()
}
