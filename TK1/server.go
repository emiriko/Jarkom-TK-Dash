package main

import (
	"fmt"
	"net"
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

	// btw gw td coba cari ide di chat gpt wkwkwk, itu idea nya ada di file serverIdeaFromChatGPT.go
	// okay

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
	// if len(os.Args) != 2 {
	//     fmt.Fprintf(os.Stderr, "Usage: %s ip-addr\n", os.Args[0])
	//     os.Exit(1)
	// }
	// service := os.Args[1]
	// tcpAddr, err := net.ResolveTCPAddr("tcp4", service)

	// if(err != nil) {
	// 	// error handling here
	// }

}

func HandleConnection(connection net.Conn) {
	//This progrom handles the incoming request from client
	defer connection.Close()

	buffer := make([]byte, 1024)

	n, err := connection.Read(buffer)

	if err != nil {
		fmt.Println("Error, reading request: ", err.Error())
	}

	// request := string(buffer[:n])

	req := RequestDecoder(buffer[:n])

	res := HandleRequest(req)

	go ResponseEncoder(res)

}

func HandleRequest(req HttpRequest) HttpResponse {
	//This program handles the routing to each view handler.
	var res HttpResponse
	const paramater = strings.Split(req.Uri, "")
	// handle the request based on its URI and method
	switch req.Uri {
	case "/":
		if req.Method == "GET" {
			// handle GET request for the root URI
			res.Version = "HTTP/1.1"
			res.StatusCode = "200 OK"
			res.ContentType = "text/plain"
			res.ContentLanguage = "en-US"
			res.Data = "Hello, World!"
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
		res.Data = "Page not found"
	}

	return res

}

func RequestDecoder(bytestream []byte) HttpRequest {
	//Put the decoding program for HTTP Request Packet here

	// you can use the bufio package to read the message data from the connection
	// parse the message into its components
	// and then create an HTTPRequest structure from those components.
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

/**
contoh HandleRequest yg gw dpt td

func HandleRequest(req HttpRequest) HttpResponse {
    var res HttpResponse

    // handle the request based on its URI and method
    switch req.Uri {
    case "/":
        if req.Method == "GET" {
            // handle GET request for the root URI
            res.Version = "HTTP/1.1"
            res.StatusCode = "200 OK"
            res.ContentType = "text/plain"
            res.ContentLanguage = "en-US"
            res.Data = "Hello, World!"
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
        res.Data = "Page not found"
    }

    return res
}
**/
