package main

import (
	// "crypto/tls"
	"context"
	"crypto/tls"
	"crypto/x509"
	"flag"
	"fmt"
	"log"
	"os"

	gRPC "github.com/PatrickMatthiesen/DiceRoll/proto"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	// "google.golang.org/grpc/credentials/insecure"
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
	rootPEM, err := os.ReadFile("keys/cointoss.pem")
	if err != nil {
		log.Fatalf("failed to read cert.pem: %s", err)
	}
	
	certPool := x509.NewCertPool()
	ok := certPool.AppendCertsFromPEM(rootPEM)
	if !ok {
		panic("failed to parse root certificate")
	}

	cert, err := tls.LoadX509KeyPair("keys/cointoss.pem", "keys/cointoss.key")
	if err != nil {
		log.Fatal(err)
	}

	cfg := &tls.Config{
		RootCAs: certPool,
		Certificates: []tls.Certificate{cert},
	}

	log.Printf("client %s: Attempts to dial on port %s\n", *clientsName, *serverPort)
    conn, err := grpc.Dial("localhost:5400",
        grpc.WithTransportCredentials(credentials.NewTLS(cfg)),
		grpc.WithBlock(),
    )
	if err != nil {
        log.Fatalf("did not connect: %v", err)
    }
    defer conn.Close()

	// var opts []grpc.DialOption
	// opts = append(opts, grpc.WithBlock(), grpc.WithTransportCredentials(credentials.NewTLS(cfg)))
	// //dial the server, with the flag "server", to get a connection to it
	// conn, err := grpc.Dial(fmt.Sprintf(":%s", *serverPort), opts...)
	// if err != nil {
	// 	log.Printf("Fail to Dial : %v", err)
	// 	return
	// }

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
