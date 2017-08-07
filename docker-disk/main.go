package docker_disk

import (
	"encoding/json"
	"flag"
	"fmt"
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

	if len(os.Args) < 2 {
		usage()
		os.Exit(2)
	}

	switch os.Args[1] {
	case "list":
		f1 := flag.NewFlagSet(os.Args[1], flag.ExitOnError)
		f1.Usage = func() {
			subCmdUsage("list", "", []map[string]string{})
		}
		f1.Parse(os.Args[2:])
		diskInfo, err := Collect()
		assert(err)
		toJSON(diskInfo)
	}
}

func main() {
	Run()
}
