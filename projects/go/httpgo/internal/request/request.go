package request

import (
	"io"
	"strings"
)

type Request struct {
	RequestLine RequestLine
}

type RequestLine struct {
	HttpVersion   string
	RequestTarget string
	Method        string
}

func RequestFromReader(reader io.Reader) (*Request, error) {

	var req Request
	var str = make([]byte, 1024)

	_, err := reader.Read(str)
	if err != nil {
		return nil, err
	}

	s := strings.Split(string(str), "\r\n")

	req.RequestLine.Method = s[0]
  req.

	return nil, nil
}
