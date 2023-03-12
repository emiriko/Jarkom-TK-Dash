//package main

import (
	"bufio"
	"fmt"
	"net"
	"os"
	"strings"
)

type HttpRequest struct {
	Method          string
	Uri             string
	Version         string
	Host            string
	Accept          string
	AcceptLanguange string
}

type HttpResponse struct {
	Version         string
	StatusCode      string
	ContentType     string
	ContentLanguage string
	Data            string
}

type Student struct {
	Nama string
	Npm  string
}

const (
	SERVER_TYPE = "tcp"
	BUFFER_SIZE = 1024
)

func main() {
	//The Program logic should go here.
	reader := bufio.NewReader(os.Stdin)

	req := HttpRequest{Version: "HTTP/1.1",
						Method: "GET"}

	fmt.Print("input the url: ")
	url, err := reader.ReadString('\n')
	if err != nil {
		fmt.Println("Error message:", err.Error())
		os.Exit(1)
	}

	urlSplit := strings.Split(url, "/")
	req.Uri = "/" + strings.Join(urlSplit[3:], "/")
	req.Host = urlSplit[2]

	fmt.Print("input the data type: ")
	req.Accept, err = reader.ReadString('\n')
	
	fmt.Print("input the language: ")
	req.AcceptLanguange, err = reader.ReadString('\n')

	tcpServer, err := net.ResolveTCPAddr(SERVER_TYPE, req.Host+ req.Uri)
	if err != nil {
		fmt.Println("Resolve TCPAddr failed:", err.Error())
		os.Exit(1)
	}

	conn, err := net.DialTCP(SERVER_TYPE, nil, tcpServer)

	Fetch(req, conn)
	
	defer conn.Close()

}

func Fetch(req HttpRequest, connection net.Conn) (HttpResponse, []Student, HttpRequest) {
	//This program handles the request-making to the server
	var res HttpResponse
	var Student []Student

	string request = RequestEncoder(req)
	_, err = conn.Write([]byte(request))
	if err != nil {
		fmt.Println("Error message:", err.Error())
		os.Exit(1)
	}

	buffer := make([]byte, BUFFER_SIZE)
	bufLen, err := conn.Read(buffer)
	if err != nil {
		fmt.Println("Error message:", err.Error())
		os.Exit(1)
	}

	res = ResponseDecoder(buffer) 


	return res, Student, req

}

func ResponseDecoder(bytestream []byte) HttpResponse {
	var res HttpResponse

	response := bytes.Split(bytestream, []byte("\r\n"))
	statusLine := strings.SplitN(string(response[0]), " ", 2)
	res.Version = string(statusLine[0])
	res.StatusCode = string(statusLine[1])
	res.ContentType = strings.Split(string(response[1]), " ")[1]
	res.ContentLanguage = strings.Split(string(response[2]), " ")[1]

	if res.ContentType == "application/json" {

	} else if res.ContentType == "application/xml" {

	} else if res.ContentType == "text/html" {
		res.Data = strings.SplitN(string(response[3]), " ", 2)
	}

	return res

}

func RequestEncoder(req HttpRequest) []byte {
	var result string

	result = fmt.Sprintf("%s %s %s\r\nHost: %s\r\nAccept: %s\r\nAccept-Language: %s",
			HttpRequest.Method, HttpRequest.Uri, HttpRequest.Version, HttpRequest.Host, HttpRequest.Accept, HttpRequest.AcceptLanguange)

	return []byte(result)

}
