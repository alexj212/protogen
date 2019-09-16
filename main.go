package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"os"
)

//
// go run protogen /home/alexj/projects/rooms/submodules/protocol-definitions/centralserver/csserver.proto Packet /home/alexj/projects/rooms/proto/cs/mapping.go
//

var DATE string
var LATEST_COMMIT string
var BUILD_NUMBER string
var BUILT_ON_IP string
var BUILT_ON_OS string

func usage() {
	fmt.Fprintf(os.Stderr, "usage: %v <protofile> <packetname> <gofile>\n", os.Args[0])
	fmt.Fprintf(os.Stderr, "\n")
	fmt.Fprintf(os.Stderr, "DATE          : %v\n", DATE)
	fmt.Fprintf(os.Stderr, "LATEST_COMMIT : %v\n", LATEST_COMMIT)
	fmt.Fprintf(os.Stderr, "BUILD_NUMBER  : %v\n", BUILD_NUMBER)
	fmt.Fprintf(os.Stderr, "BUILT_ON_IP   : %v\n", BUILT_ON_IP)
	fmt.Fprintf(os.Stderr, "BUILT_ON_OS   : %v\n", BUILT_ON_OS)
	fmt.Fprintf(os.Stderr, "\n")
}

func main() {
	var formatCode bool

	flag.BoolVar(&formatCode, "format", false, "Format code after generating")
	flag.Parse()

	if len(os.Args) != 4 {
		usage()
		os.Exit(1)
		return
	}

	protoFile := os.Args[1]
	enumName := os.Args[2]
	goFile := os.Args[3]

	if !fileExists(protoFile) {
		fmt.Fprintf(os.Stderr, "Protobuf file %v does not exist\n", protoFile)
		os.Exit(1)
		return
	}

	parsed, err := Parse(protoFile, enumName)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error parsing: %v\n", err)
		os.Exit(1)
		return
	}

	code, err := Generate(parsed, false)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error generating code: %v\n", err)
		os.Exit(1)
		return
	}

	f, err := os.Create(goFile)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error creating blob file:%v\n", err)
		os.Exit(1)
	}
	defer f.Close()

	err = ioutil.WriteFile(goFile, code, os.ModePerm)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error writing blob file error: %v\n", err)
		os.Exit(1)
	}

	fmt.Fprintf(os.Stdout, "File Written: %v\n", goFile)
}

// fileExists checks if a file exists and is not a directory before we
// try using it to prevent further errors.
func fileExists(filename string) bool {
	info, err := os.Stat(filename)
	if os.IsNotExist(err) {
		return false
	}
	return !info.IsDir()
}
