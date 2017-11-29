package main

import (
	"os"
	"regexp"
	"strconv"
	"strings"
)

const (
	columnPattern = "^((NF|[0-9]+))|([0-9]|NF)/([0-9]|NF)$"
)

var columnPatternRegex = regexp.MustCompile(columnPattern)

type Parser struct {
}

func NewParser() Parser {
	return Parser{}
}

func (p Parser) Parse(args []string, fieldCount int) []int {
	length := len(args)
	if length < 1 {
		return []int{}
	}
	if p.hasExistingFileName(args) {
		length -= 1
	}
	result := make([]int, 0, length)
	for _, v := range args {
		if !p.isColumn(v) {
			return []int{}
		}
		// NF(Number of field) を含むケースは、NFをフィールド数の値に置き換える。
		if strings.Contains(v, "NF") {
			v = strings.Replace(v, "NF", strconv.Itoa(fieldCount), 2)
		}
		// 数値のみのケース
		if v, err := strconv.Atoi(v); err == nil {
			result = append(result, v)
			continue
		}
		// "/" を含むケースは、"/" の左側の値を始点、右側を終点(閉区間)とした連続する整数値を与える
		// 例: 2/NF -> 2,3,4,5,6 (NF == 6 の場合)
		if strings.Contains(v, "/") {
			slashedResult := p.slashedColumnToNumbers(v)
			if len(slashedResult) == 0 {
				return []int{}
			}
			result = append(result, slashedResult...)
			continue
		}
	}

	return result
}

func (p Parser) slashedColumnToNumbers(v string) []int {
	split := strings.Split(v, "/")
	if len(split) != 2 {
		return []int{}
	}
	start, err := strconv.Atoi(split[0])
	if err != nil {
		return []int{}
	}
	end, err := strconv.Atoi(split[1])
	if err != nil {
		return []int{}
	}
	length := end - start + 1
	result := make([]int, 0, length)
	for i := start; i <= end; i++ {
		result = append(result, i)
	}
	return result
}

func (p Parser) isColumn(arg string) bool {
	return columnPatternRegex.Copy().MatchString(arg)
}

func (p Parser) hasExistingFileName(args []string) bool {
	fileNameCandidate := args[len(args)-1]
	_, err := os.Stat(fileNameCandidate)
	if err != nil {
		return false
	}
	return true
}
