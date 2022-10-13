package main

import (
	// "crypto/tls"
	"context"
	"flag"
	"fmt"
	"log"
	"crypto/tls"
	"crypto/x509"
	"encoding/pem"
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
	// rootPEM, err := os.ReadFile("keys/cointoss.pem")
	// if err != nil {
	// 	log.Fatalf("failed to read cert.pem: %s", err)
	// }
	
	// certPool := x509.NewCertPool()
	// ok := certPool.AppendCertsFromPEM(rootPEM)
	// if !ok {
	// 	panic("failed to parse root certificate")
	// }

	// cert, err := tls.LoadX509KeyPair("keys/cointoss.pem", "keys/cointoss.key")
	// if err != nil {
	// 	log.Fatal(err)
	// }

	// cfg := &tls.Config{
	// 	RootCAs: certPool,
	// 	Certificates: []tls.Certificate{cert},
	// }

	// creds, _ := credentials.NewClientTLSFromFile("keys/client-cert.pem", "example.org")

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
    certPool := x509.NewCertPool()
    certs := []tls.Certificate{}

	// Read certificate files
	srvPemBytes, err := os.ReadFile(fmt.Sprintf("keys/%v.cert.pem", "bob"))
	if err != nil {
		log.Fatalf("%v\n", err)
	}

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

	// Load server certificates (essentially the same as the client certs)
	srvCert, err := tls.LoadX509KeyPair(fmt.Sprintf("keys/%v.cert.pem", "bob"), fmt.Sprintf("keys/%v.key.pem", "bob"))
	if err != nil {
		log.Fatalf("%v\n", err)
	}
	certs = append(certs, srvCert)

    return &tls.Config{
        Certificates: certs, // Server certs
        ClientAuth:   tls.RequireAndVerifyClientCert,
        ClientCAs:    certPool,
        RootCAs:      certPool,
    }
}