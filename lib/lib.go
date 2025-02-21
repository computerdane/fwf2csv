package lib

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"strconv"
	"strings"
)

type Column struct {
	N    int
	Name string
}

func ReadSchema(schema string) ([]Column, error) {
	lines := strings.Split(strings.TrimSpace(schema), "\n")
	columns := make([]Column, len(lines))

	for i, line := range lines {
		tokens := strings.SplitN(line, " ", 2)
		if len(tokens) != 2 {
			return nil, fmt.Errorf("Invalid column definition on line %d: %s", i+1, line)
		}

		n64, err := strconv.ParseInt(tokens[0], 10, 32)
		if err != nil {
			return nil, fmt.Errorf("Invalid column definition on line %d: %s is not an int", i+1, tokens[0])
		}

		n := int(n64)
		columns[i] = Column{N: n, Name: tokens[1]}
	}

	return columns, nil
}

func Convert(columns []Column, r io.Reader, w io.Writer, delimiter string) error {
	for i, column := range columns {
		if i > 0 {
			fmt.Fprint(w, delimiter)
		}
		fmt.Fprint(w, column.Name)
	}
	fmt.Fprintln(w)

	scanner := bufio.NewScanner(r)
	for scanner.Scan() {
		line := scanner.Text()
		a := 0
		for i, column := range columns {
			if i > 0 {
				if column.N <= columns[i-1].N {
					return fmt.Errorf("Invalid column definition at index %d: N should be greater than %d in the previous column", i, columns[i-1].N)
				}

				fmt.Fprint(w, delimiter)
			}
			b := min(column.N, len(line))
			value := line[a:b]
			trimmed := strings.TrimSpace(value)
			fmt.Fprint(w, trimmed)
			a = b
		}
		fmt.Fprintln(w)
	}

	return nil
}

func ConvertString(columns []Column, data string, delimiter string) (string, error) {
	r := strings.NewReader(data)
	var w bytes.Buffer

	if err := Convert(columns, r, &w, delimiter); err != nil {
		return "", err
	}

	return w.String(), nil
}
