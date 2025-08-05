package utils

import (
	"bufio"
	"fmt"
	"strconv"
	"strings"

	"github.com/fatih/color"
)

func ReadAndSanitizeInput(msg string, reader *bufio.Reader) (string, error) {
	color.Blue(msg)
	input, err := reader.ReadString('\n')
	if err != nil {
		return "", err
	}
	return strings.TrimRight(input, "\r\n"), nil
}

func ReadIntList(msg string, reader *bufio.Reader) ([]int, error) {
	input, err := ReadAndSanitizeInput(msg, reader)
	if err != nil {
		return nil, err
	}

	var res []int

	for str := range strings.SplitSeq(input, " ") {
		num, err := strconv.Atoi(strings.TrimSpace(str))
		if err == nil {
			res = append(res, num)
		}
	}

	return res, nil
}

func PrintListInRows(list []string) {
	for i, name := range list {
		fmt.Print(color.BlueString("%d. %s\t", i+1, name))
		if (i+1)%5 == 0 {
			fmt.Println("")
		}
	}
	fmt.Println("")
}
