package main

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
)

type Headers struct {
	IPAddr  string `json:"address"`
	ASN     string `json:"asn"`
	Country string `json:"country"`
}

func handler(ctx context.Context, req events.APIGatewayV2HTTPRequest) (events.APIGatewayV2HTTPResponse, error) {
	var responseBody string
	var contentType string

	if req.Headers["cloudfront-viewer-address"] == "" {
		return events.APIGatewayV2HTTPResponse{
			StatusCode: http.StatusForbidden,
			Headers: map[string]string{
				"Content-Type": "text/plain",
			},
			Body: "Forbidden",
		}, nil
	}

	var ipAddr string
	ipAddrSlice := strings.Split(req.Headers["cloudfront-viewer-address"], ":")
	if (len(ipAddrSlice) - 1) > 2 {
		ipAddr = strings.Join(ipAddrSlice[:len(ipAddrSlice)-1], ":")
	} else {
		ipAddr = ipAddrSlice[0]
	}

	h := Headers{
		IPAddr:  ipAddr,
		ASN:     req.Headers["cloudfront-viewer-asn"],
		Country: req.Headers["cloudfront-viewer-country"],
	}
	data, err := json.Marshal(h)
	if err != nil {
		data, _ = json.Marshal(Headers{})
	}

	if req.RawPath == "/whoami" {
		responseBody = h.IPAddr
		contentType = "text/plain"
	} else if req.RawPath == "/whoami/json" {
		responseBody = string(data)
		contentType = "application/json"
	}

	fmt.Println(string(data))

	return events.APIGatewayV2HTTPResponse{
		StatusCode: http.StatusOK,
		Headers: map[string]string{
			"Content-Type": contentType,
		},
		Body: responseBody,
	}, nil
}

func main() {
	lambda.Start(handler)
}
