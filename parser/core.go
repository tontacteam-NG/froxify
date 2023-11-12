package parser

import (
	"bufio"
	"net/http"
	"strings"
)

func ParseRequest(msg string) (*http.Request, error) {
	reqReader := strings.NewReader(msg)
	req, err := http.ReadRequest(bufio.NewReader(reqReader))
	if err != nil {
		return nil, err
	}
	return req, nil
}

func ParseFolder(data string) (int, []string, string) {
	arr_temp := strings.Split(data, "/")
	len := len(arr_temp)
	arrDir := arr_temp[0 : len-1]
	endArr := arr_temp[len-1]
	return len, arrDir, endArr
}
