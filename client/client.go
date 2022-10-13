package main

import (
	// "crypto/tls"
	"context"
	"crypto/tls"
	"crypto/x509"
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
	cert, err := tls.LoadX509KeyPair("client_cert.pem", "client_key.pem")
	if err != nil {
		log.Fatalf("failed to load client cert: %v", err)
	}

	ca := x509.NewCertPool()
	caFilePath := "ca_cert.pem"
	caBytes, err := os.ReadFile(caFilePath)
	if err != nil {
		log.Fatalf("failed to read ca cert %q: %v", caFilePath, err)
	}
	if ok := ca.AppendCertsFromPEM(caBytes); !ok {
		log.Fatalf("failed to parse %q", caFilePath)
	}

	return &tls.Config{
		ServerName:   "x.test.example.com",
		Certificates: []tls.Certificate{cert},
		RootCAs:      ca,
	}
}