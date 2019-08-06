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


func usage() {
    fmt.Fprintf(os.Stderr, "usage: %v <protofile> <packetname> <gofile>\n", os.Args[0])
    fmt.Fprintf(os.Stderr, "\n")
}


func main() {
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

    code, err := Generate(parsed)
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

