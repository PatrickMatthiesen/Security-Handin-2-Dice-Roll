package main

import (
	// "crypto/tls"
	"context"
	"crypto/tls"
	"crypto/x509"
	"encoding/pem"
	"flag"
	"fmt"
	"log"

	// "encoding/pem"
	"os"

	gRPC "github.com/PatrickMatthiesen/DiceRoll/proto"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
)

// Same principle as in client. Flags allows for user specific arguments/values
var clientsName = flag.String("name", "default", "Senders name")
var serverPort = flag.String("server", "5400", "Tcp server")

var server gRPC.DiceRollServiceClient  //the server
var ServerConn *grpc.ClientConn //the server connection

func main() {
	//parse flag/arguments
	flag.Parse()

	fmt.Println("--- CLIENT APP ---")

	//log to file instead of console
	//setLog()

	//connect to server and close the connection when program closes
	fmt.Println("--- join Server ---")
	ConnectToServer()
	defer ServerConn.Close()

	//start the biding
}

// connect to server
func ConnectToServer() {
	creds := credentials.NewTLS(getTLSConfig())

	log.Printf("client %s: Attempts to dial on port %s\n", *clientsName, *serverPort)
    conn, err := grpc.Dial("localhost:5400",
        grpc.WithTransportCredentials(creds),
		grpc.WithBlock(),
    )
	if err != nil {
        log.Fatalf("did not connect: %v", err)
    }
    defer conn.Close()

	log.Println("client: Dialing successful")

	// makes a client from the server connection and saves the connection
	// and prints rather or not the connection was is READY
	server = gRPC.NewDiceRollServiceClient(conn)
	ServerConn = conn
	log.Println("the connection is: ", conn.GetState().String())

	responce, err := server.CommitRoll(context.Background(), &gRPC.Commitment{Commitment: 1})
	if err != nil {
		log.Fatalf("could not greet: %v", err)
	}
	log.Println(responce)
}


func getTLSConfig() *tls.Config {
	// Read certificate files
	srvPemBytes, err := os.ReadFile("keys/bob.cert.pem")
	if err != nil {
		log.Fatalf("%v\n", err)
	}
	
	certPool := x509.NewCertPool()
	certPool.AppendCertsFromPEM(srvPemBytes)

	// Decode and parse certs
	srvPemBlock, _ := pem.Decode(srvPemBytes)
	clientCert, err := x509.ParseCertificate(srvPemBlock.Bytes)
	if err != nil {
		log.Fatalf("%v\n", err)
	}

	// Enforce client authentication and allow self-signed certs
	clientCert.BasicConstraintsValid = true
	clientCert.IsCA = true
	clientCert.KeyUsage = x509.KeyUsageCertSign
	certPool.AppendCertsFromPEM(srvPemBytes)

    return &tls.Config{
		Certificates: []tls.Certificate{clientCert},
        RootCAs:      certPool,
		ClientCAs:    certPool,
    }
}