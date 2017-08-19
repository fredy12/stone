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
	fmt.Println("    xfs-quotactl [global options] command [command options] [arguments...]\n")
	fmt.Println("VERSION:")
	fmt.Printf("    %s\n\n", getVersion())
	fmt.Println("COMMANDS:")
	fmt.Println("    get     get quota info of path")
	fmt.Println("GLOBAL OPTIONS:")
	fmt.Println("   --help, -h     Show help")
}

func subCmdUsage(name string, cmdName string, options []map[string]string) {
	fmt.Println("USAGE:")
	if len(options) > 0 {
		fmt.Printf("    xfs-quotactl [options] %s %s\n\n", name, cmdName)
		fmt.Println("OPTIONS:")

	} else {
		fmt.Printf("    xfs-quotactl %s %s\n\n", name, cmdName)
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
			subCmdUsage("get", "", []map[string]string{
				{"path": "path to show quota"},
				{"verbose": "Optional show verbose logs"},
			})
		}

		f1.Parse(os.Args[2:])
		logrus.SetLevel(logrus.PanicLevel)
		path := os.Args[2]
		if len(os.Args) > 3 {
			if os.Args[3] == "verbose" {
				logrus.SetLevel(logrus.DebugLevel)
			}
		}

		projectId, err := getProjectID(path)
		assert(err)

		qc := XfsQuotaControl{}
		qc.quotas[path] = projectId
		quotas, err := qc.GetQuota(path)
		assert(err)

		type resp struct {
			ProjectId  uint32
			QuotaLimit uint64
		}
		r := &resp{
			ProjectId:  projectId,
			QuotaLimit: quotas.Size,
		}
		toJSON(r)
	}
}

func main() {
	Run()
}
