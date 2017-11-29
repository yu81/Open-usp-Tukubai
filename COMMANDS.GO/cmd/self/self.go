package main

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"strings"
)

const separator = " "

func main() {
	argParser := NewParser()
	args := os.Args[1:]
	input := getInput(argParser, args)

	scanner := bufio.NewScanner(input)
	firstLineFields, columnCount, selectedColumns := ProcessFirstLine(scanner, argParser, args)
	PrintSingle(firstLineFields, selectedColumns, columnCount)

	for {
		ProcessLine(scanner, selectedColumns, columnCount)
	}
}

func getInput(argParser Parser, args []string) io.Reader {
	if argParser.hasExistingFileName(args) {
		input, err := os.Open(args[len(args)-1])
		if err != nil {
			os.Exit(1)
		}
		return input
	} else {
		return os.Stdin
	}
}

func PrintSingle(fields []string, targetColumnIndices []int, maxLength int) {
	outputFields := make([]string, 0, maxLength)
	for i := 0; i < len(targetColumnIndices); i++ {
		outputFields = append(outputFields, fields[targetColumnIndices[i]-1])
	}
	fmt.Println(strings.Join(outputFields, separator))
}

func ProcessFirstLine(scanner *bufio.Scanner, argParser Parser, args []string) ([]string, int, []int) {
	if ok := scanner.Scan(); !ok {
		os.Exit(1)
	}
	firstLine := scanner.Text()
	firstLineFields := strings.Split(firstLine, separator)
	columnCount := len(firstLineFields)
	selectedColumns := argParser.Parse(args, columnCount)

	return firstLineFields, columnCount, selectedColumns
}

func ProcessLine(scanner *bufio.Scanner, selectedColumns []int, columnCount int) {
	if ok := scanner.Scan(); !ok {
		os.Exit(0)
	}
	line := scanner.Text()
	fields := strings.Split(line, separator)
	PrintSingle(fields, selectedColumns, columnCount)
}
