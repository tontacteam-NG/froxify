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
