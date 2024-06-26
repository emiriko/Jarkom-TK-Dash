package main

import (
	"encoding/json"
	"encoding/xml"
	"fmt"
	"net"
	"os"
	"strings"
)

const (
	SERVER_HOST = ""
	SERVER_PORT = "2150"
	SERVER_TYPE = "tcp"
	BUFFER_SIZE = 1024
	GROUP_NAME  = "Dash"
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

func main() {
	//The Program logic should go here.

	listen, err := net.Listen(SERVER_TYPE, SERVER_HOST+":"+SERVER_PORT)

	if err != nil {
		fmt.Println("Error listening connection: ", err.Error())
		return
	}

	defer listen.Close()

	for {
		conn, err := listen.Accept()
		if err != nil {
			fmt.Println("Error accepting connection: ", err.Error())
		}
		go HandleConnection(conn)
	}
}

func HandleConnection(connection net.Conn) {
	//This progrom handles the incoming request from client
	defer connection.Close()

	buffer := make([]byte, BUFFER_SIZE)

	n, err := connection.Read(buffer)

	if err != nil {
		fmt.Println("Error, reading request: ", err.Error())
	}

	req := RequestDecoder(buffer[:n])

	res := HandleRequest(req)

	result := ResponseEncoder(res)

	_, err = connection.Write([]byte(result))

	if err != nil {
		fmt.Println("Error message:", err.Error())
		os.Exit(1)
	}

	defer connection.Close()
}

func HandleRequest(req HttpRequest) HttpResponse {
	//This program handles the routing to each view handler.
	var res HttpResponse
	student := []Student{
		{Nama: "Sasha Nabila Fortuna", Npm: "2106632226"},
		{Nama: "Dianisa Wulandari", Npm: "2106702150"},
		{Nama: "Alvaro Austin", Npm: "2106752180"},
	}

	jsonData, err := json.Marshal(student)

	if err != nil {
		fmt.Println(err)
	}

	xmlData, err := xml.Marshal(student)

	if err != nil {
		fmt.Println(err)
	}

	// handle the request based on its URI and method
	switch req.Uri {
	case "/", fmt.Sprintf("/?name=%s", GROUP_NAME):
		if req.Method == "GET" {
			// handle GET request for the root URI
			res.Version = "HTTP/1.1"
			res.StatusCode = "200 OK"
			res.ContentType = "text/html"
			res.ContentLanguage = "id-ID"
			res.Data = fmt.Sprintf("<html><body><h1>Halo, kami dari %s</h1></body></html>", GROUP_NAME)
		} else {
			// handle unsupported method for the root URI
			res.Version = "HTTP/1.1"
			res.StatusCode = "405 Method Not Allowed"
			res.ContentType = "text/plain"
			res.ContentLanguage = "en-US"
			res.Data = "Method not allowed"
		}
	case "/greeting":
		if req.Method == "GET" {
			switch req.AcceptLanguange {
			case "en-US":
				res.Version = "HTTP/1.1"
				res.StatusCode = "200 OK"
				res.ContentType = "text/html"
				res.ContentLanguage = req.AcceptLanguange
				res.Data = fmt.Sprintf("<html><body><h1>Hello, we are from %s</h1></body></html>", GROUP_NAME)
			case "id-ID":
				res.Version = "HTTP/1.1"
				res.StatusCode = "200 OK"
				res.ContentType = "text/html"
				res.ContentLanguage = req.AcceptLanguange
				res.Data = fmt.Sprintf("<html><body><h1>Halo, kami dari %s</h1></body></html>", GROUP_NAME)
			default:
				res.Version = "HTTP/1.1"
				res.StatusCode = "200 OK"
				res.ContentType = "text/html"
				res.ContentLanguage = "en-US"
				res.Data = fmt.Sprintf("<html><body><h1>Hello, we are from %s</h1></body></html>", GROUP_NAME)
			}
		} else {
			// handle unsupported method for the root URI
			res.Version = "HTTP/1.1"
			res.StatusCode = "405 Method Not Allowed"
			res.ContentType = "text/plain"
			res.ContentLanguage = "en-US"
			res.Data = "Method not allowed"
		}
	case "/data":
		if req.Method == "GET" {
			switch req.Accept {
			case "application/json":
				res.Version = "HTTP/1.1"
				res.StatusCode = "200 OK"
				res.ContentType = req.Accept
				res.ContentLanguage = "en-US"
				res.Data = string(jsonData)
			case "application/xml":
				res.Version = "HTTP/1.1"
				res.StatusCode = "200 OK"
				res.ContentType = req.Accept
				res.ContentLanguage = "en-US"
				res.Data = string(xmlData)
			default:
				res.Version = "HTTP/1.1"
				res.StatusCode = "200 OK"
				res.ContentType = "application/json"
				res.ContentLanguage = "en-US"
				res.Data = string(jsonData)
			}
		} else {
			// handle unsupported method for the root URI
			res.Version = "HTTP/1.1"
			res.StatusCode = "405 Method Not Allowed"
			res.ContentType = "text/plain"
			res.ContentLanguage = "en-US"
			res.Data = "Method not allowed"
		}
	default:
		// handle unknown URI
		res.Version = "HTTP/1.1"
		res.StatusCode = "404 Not Found"
		res.ContentType = "text/plain"
		res.ContentLanguage = "en-US"
		res.Data = ""
	}

	return res

}

func RequestDecoder(bytestream []byte) HttpRequest {
	//Put the decoding program for HTTP Request Packet here
	var req HttpRequest

	// split the request into lines
	lines := strings.Split(string(bytestream), "\r\n")

	// parse the request line
	requestLine := strings.Split(lines[0], " ")
	req.Method = requestLine[0]
	req.Uri = requestLine[1]
	req.Version = requestLine[2]

	// parse headers
	for _, line := range lines[1:] {
		if line == "" {
			break // end of headers
		}
		parts := strings.SplitN(line, ":", 2)
		if len(parts) != 2 {
			continue // skip invalid header lines
		}
		key := strings.TrimSpace(parts[0])
		value := strings.TrimSpace(parts[1])
		switch key {
		case "Host":
			req.Host = value
		case "Accept":
			req.Accept = value
		case "Accept-Language":
			req.AcceptLanguange = value
		}
	}
	return req

}

func ResponseEncoder(res HttpResponse) []byte {
	//Put the encoding program for HTTP Response Struct here
	var result string

	// write the status line
	result += res.Version + " " + res.StatusCode + "\r\n"

	// write the headers
	result += "Content-Type: " + res.ContentType + "\r\n"
	result += "Content-Language: " + res.ContentLanguage + "\r\n"

	// write the response body
	result += "\r\n" + res.Data

	return []byte(result)

}
