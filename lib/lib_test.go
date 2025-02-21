package lib

import (
	"bytes"
	"fmt"
	"os"
	"testing"
)

func getTestColumns(t *testing.T) []Column {
	schemaFile := "../examples/schema"
	schemaBytes, err := os.ReadFile(schemaFile)
	if err != nil {
		t.Fatalf("Failed to read %s: %s", schemaFile, err)
	}
	schema := string(schemaBytes)

	columns, err := ReadSchema(schema)
	if err != nil {
		t.Fatal(err)
	}

	if len(columns) != 5 {
		t.Fatalf("len(columns) is %d; expected 5", len(columns))
	}

	return columns
}

func TestReadSchema(t *testing.T) {
	columns := getTestColumns(t)

	if columns[0].N != 10 ||
		columns[0].Name != "first_name" ||
		columns[1].N != 20 ||
		columns[1].Name != "last_name" ||
		columns[2].N != 24 ||
		columns[2].Name != "year" ||
		columns[3].N != 26 ||
		columns[3].Name != "month" ||
		columns[4].N != 28 ||
		columns[4].Name != "day" {
		t.Error("columns does not match expected values")
	}
}

func TestConvert(t *testing.T) {
	columns := getTestColumns(t)

	file := "../examples/data"
	data, err := os.Open(file)
	if err != nil {
		t.Fatalf("Failed to open %s: %s", file, err)
	}

	var w bytes.Buffer
	if err := Convert(columns, data, &w); err != nil {
		t.Fatal(err)
	}

	output := w.String()

	csvFile := "../examples/data.csv"
	csvBytes, err := os.ReadFile(csvFile)
	if err != nil {
		t.Fatalf("Failed to read %s: %s", csvFile, err)
	}
	csv := string(csvBytes)

	if output != csv {
		fmt.Println(output)
		fmt.Println(csv)
		t.Error("Converted output and expected CSV do not match.")
	}
}
