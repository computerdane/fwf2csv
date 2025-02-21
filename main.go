package main

import (
	_ "embed"
	"fmt"
	"log"
	"os"

	"github.com/computerdane/fwf2csv/lib"
	"github.com/computerdane/gears"
)

var Version string

//go:embed README
var readme string

func printReadme() {
	fmt.Println(readme)
	os.Exit(1)
}

func init() {
	gears.Add(&gears.Flag{
		Name:      "version",
		ValueType: "bool",
		Shorthand: "v",
	})
	gears.Add(&gears.Flag{
		Name:      "help",
		ValueType: "bool",
		Shorthand: "h",
	})
	gears.Add(&gears.Flag{
		Name:         "schema",
		ValueType:    "string",
		DefaultValue: "",
		Shorthand:    "s",
	})
	gears.Add(&gears.Flag{
		Name:         "schema-file",
		ValueType:    "string",
		DefaultValue: "",
		Shorthand:    "S",
	})
	gears.Add(&gears.Flag{
		Name:         "delimiter",
		ValueType:    "string",
		DefaultValue: ",",
		Shorthand:    "d",
	})
}

func main() {
	gears.Load()

	if gears.BoolValue("version") {
		if Version == "" {
			Version = "unknown version"
		}
		fmt.Printf("fwf2csv  %s\n", Version)
		os.Exit(0)
	}

	if gears.BoolValue("help") {
		printReadme()
	}

	args := gears.Positionals()

	if len(args) != 1 {
		printReadme()
	}

	schema := gears.StringValue("schema")
	schemaFile := gears.StringValue("schema-file")
	delimiter := gears.StringValue("delimiter")

	if delimiter == "" {
		log.Fatal("Delimiter must not be empty!")
	}

	if schema == "" && schemaFile == "" {
		log.Fatal("Must specify either --schema or --schema-file")
	}
	if schema != "" && schemaFile != "" {
		log.Fatal("The options --schema and --schema-file cannot be used together")
	}

	if schemaFile != "" {
		schemaBytes, err := os.ReadFile(schemaFile)
		if err != nil {
			log.Fatalf("Failed to read %s: %s", schemaFile, err)
		}
		schema = string(schemaBytes)
	}

	columns, err := lib.ReadSchema(schema)
	if err != nil {
		log.Fatal(err)
	}

	var input *os.File
	defer input.Close()

	file := args[0]
	if file == "-" {
		input = os.Stdin
	} else {
		input, err = os.Open(file)
		if err != nil {
			log.Fatalf("Failed to open %s: %s", file, err)
		}
	}

	if err := lib.Convert(columns, input, os.Stdout, delimiter); err != nil {
		log.Fatal(err)
	}
}
