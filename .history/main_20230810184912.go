package main

import (
	"fmt"
	"log"
	"os"

	"github.com/alexj212/gox"
	"github.com/alexj212/gox/utilx"
	"github.com/droundy/goopt"
)

//
// go run protogen /home/alexj/projects/rooms/submodules/protocol-definitions/centralserver/csserver.proto Packet /home/alexj/projects/rooms/proto/cs/mapping.go
//

var (
	// BuildDate date string of when build was performed filled in by -X compile flag
	BuildDate string

	// GitRepo string of the git repo url when build was performed filled in by -X compile flag
	GitRepo string

	// BuiltBy date string of who performed build filled in by -X compile flag
	BuiltBy string

	// CommitDate date string of when commit of the build was performed filled in by -X compile flag
	CommitDate string

	// Branch string of branch in the git repo filled in by -X compile flag
	Branch string

	// LatestCommit date string of when build was performed filled in by -X compile flag
	LatestCommit string

	// Version string of build filled in by -X compile flag
	Version string
)

var (
	formatCode       = goopt.Flag([]string{"--format"}, nil, "Format code after generating", "")
	forceOverwrite   = goopt.Flag([]string{"--overwrite", "-o"}, nil, "Format code after generating", "")
	fieldPrefix      = goopt.String([]string{"--fieldPrefix"}, "", "Field Prefix to use - defaults to enumName")
	protoFile        = goopt.String([]string{"--proto"}, "", "protobuf file")
	jsonFile         = goopt.String([]string{"--json"}, "", "json mappings file")
	templateFileName = goopt.String([]string{"--template"}, "", "Template file to use")
	enumName         = goopt.String([]string{"--enum"}, "", "Enum Name")
)

func init() {
	// Setup goopts
	goopt.Description = func() string {
		return fmt.Sprintf("ProtoGen")
	}
	goopt.Version = fmt.Sprintf("v%v - GitCommit: %v - BuildDate: %v", Version, LatestCommit, BuildDate)
	goopt.Summary = `
` //Parse options
	goopt.Parse(nil)

} // init

func main() {
	exitApp, err := gox.HandleHistory()
	if err != nil {
		fmt.Printf("Error handling history: %v\n", err)
		os.Exit(1)
	}
	if exitApp {
		fmt.Printf("exitApp: %v\n", exitApp)
		os.Exit(1)
	}

	if len(goopt.Args) != 1 {
		fmt.Printf("must specify arguments\n\n")
		fmt.Printf(goopt.Usage())
		os.Exit(1)
	}

	outputFile := goopt.Args[0]

	codeTemplate := messageMapperTemplate

	if *templateFileName != "" {
		templateBytes, err := os.ReadFile(*templateFileName)
		if err != nil {
			log.Printf("Failed to load template file: %s, err: %v", *templateFileName, err)
			os.Exit(1)
			return
		}
		codeTemplate = string(templateBytes)
		log.Printf("Using template file: %s", *templateFileName)
	}

	if *protoFile == "" && *jsonFile == "" {
		fmt.Fprintf(os.Stderr, "No event list loaded, specify --proto <file> or --json <file> \n")
		os.Exit(1)
		return
	}

	var parsed *MessageMapper

	if *protoFile != "" {
		if !fileExists(*protoFile) {
			fmt.Fprintf(os.Stderr, "Protobuf file %v does not exist\n", protoFile)
			os.Exit(1)
			return
		}

		if *fieldPrefix == "" {
			fieldPrefix = enumName
		}

		parsed, err = ParseProtoBuf(*protoFile, *enumName, *fieldPrefix)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error parsing: %v\n", err)
			os.Exit(1)
			return
		}
	}

	if *jsonFile != "" {
		parsed, err = utilx.LoadJson(*jsonFile, parsed)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error parsing: %v\n", err)
			os.Exit(1)
			return
		}
	}

	if parsed == nil {
		fmt.Fprintf(os.Stderr, "No event list loaded, specify --proto <file> \n")
		os.Exit(1)
		return
	}

	code, err := Generate(codeTemplate, parsed, false)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error generating code: %v\n", err)
		os.Exit(1)
		return
	}

	exists := false

	if _, err := os.Stat(outputFile); err == nil {
		exists = true
	}

	if exists && !*forceOverwrite {
		fmt.Fprintf(os.Stderr, "Output file exists: %s\n", outputFile)
		os.Exit(1)
		return
	}

	f, err := os.Create(outputFile)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error creating blob file:%v\n", err)
		os.Exit(1)
	}
	defer f.Close()

	err = os.WriteFile(outputFile, code, os.ModePerm)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error writing blob file error: %v\n", err)
		os.Exit(1)
	}

	fmt.Fprintf(os.Stdout, "File Written: %v\n", outputFile)
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
