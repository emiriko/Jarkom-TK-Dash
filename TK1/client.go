package main

import (
	"bufio"
	"encoding/json"
	"encoding/xml"
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
	var req HttpRequest
	var res HttpResponse

	// create reader and read user input
	reader := bufio.NewReader(os.Stdin)

	req = HttpRequest{Version: "HTTP/1.1",
		Method: "GET"}

	fmt.Print("input the url: ")
	url, err := reader.ReadString('\n')
	if err != nil {
		fmt.Println("Error message:", err.Error())
		os.Exit(1)
	}

	urlSplit := strings.Split(url, "/")
	req.Uri = "/" + strings.TrimSpace(strings.Join(urlSplit[3:], "/"))
	req.Host = urlSplit[2]

	fmt.Print("input the data type: ")
	req.Accept, err = reader.ReadString('\n')

	fmt.Print("input the language: ")
	req.AcceptLanguange, err = reader.ReadString('\n')

	// create tcp connection
	tcpServer, err := net.ResolveTCPAddr(SERVER_TYPE, req.Host)
	if err != nil {
		fmt.Println("Resolve TCPAddr failed:", err.Error())
		os.Exit(1)
	}

	conn, err := net.DialTCP(SERVER_TYPE, nil, tcpServer)

	// make request
	res, _, req = Fetch(req, conn)

	// close connection and print response
	defer conn.Close()
	fmt.Println("Status Code: " + (strings.SplitN(res.StatusCode, " ", 2))[0])
	fmt.Println("Body: " + res.Data)
}

func Fetch(req HttpRequest, connection net.Conn) (HttpResponse, []Student, HttpRequest) {
	//This program handles the request-making to the server
	var res HttpResponse
	var student []Student

	// encode and send request to server
	request := RequestEncoder(req)
	_, err := connection.Write([]byte(request))

	if err != nil {
		fmt.Println("Error message:", err.Error())
		os.Exit(1)
	}

	// read server response
	buffer := make([]byte, BUFFER_SIZE)
	bufLen, err := connection.Read(buffer)
	if err != nil {
		fmt.Println("Error message:", err.Error())
		os.Exit(1)
	}

	// decode response
	res = ResponseDecoder(buffer[:bufLen])
	student = []Student{}

	if res.ContentType == "application/json" {
		// Unmarshal data into byte into &student struct to store parsed data
		err := json.Unmarshal([]byte(res.Data), &student)
		if err != nil {
			fmt.Println(err)
		}

	} else if res.ContentType == "application/xml" {
		err := xml.Unmarshal([]byte(res.Data), &student)
		if err != nil {
			fmt.Printf("error: %v", err)
		}
	}

	return res, student, req

}

func ResponseDecoder(bytestream []byte) HttpResponse {
	var res HttpResponse

	response := strings.Split(string(bytestream), "\r\n")
	statusLine := strings.SplitN(response[0], " ", 2)
	res.Version = statusLine[0]
	res.StatusCode = statusLine[1]
	res.ContentType = strings.Split(response[1], " ")[1]
	res.ContentLanguage = strings.Split(response[2], " ")[1]
	res.Data = response[4]

	return res

}

func RequestEncoder(req HttpRequest) []byte {
	var result string

	result = fmt.Sprintf("%s %s %s\r\nHost: %s\r\nAccept: %s\r\nAccept-Language: %s\r\n\r\n",
		req.Method, req.Uri, req.Version, req.Host, req.Accept, req.AcceptLanguange)

	return []byte(result)

}
