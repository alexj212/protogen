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

// DATE info from build
var DATE string

// LATEST_COMMIT info from build
var LATEST_COMMIT string

// BUILD_NUMBER info from build
var BUILD_NUMBER string

// BUILT_ON_IP info from build
var BUILT_ON_IP string

// BUILT_ON_OS info from build
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
	var fieldPrefix string

	flag.BoolVar(&formatCode, "format", false, "Format code after generating")
	flag.StringVar(&fieldPrefix, "fieldPrefix", "", "Field Prefix to use - defaults to enumName")
	flag.Parse()

	if len(flag.Args()) != 3 {
		usage()
		os.Exit(1)
		return
	}

	protoFile := flag.Args()[0]
	enumName := flag.Args()[1]
	goFile := flag.Args()[2]

	if !fileExists(protoFile) {
		fmt.Fprintf(os.Stderr, "Protobuf file %v does not exist\n", protoFile)
		os.Exit(1)
		return
	}

	if fieldPrefix == "" {
		fieldPrefix = enumName
	}

	parsed, err := Parse(protoFile, enumName, fieldPrefix)
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
