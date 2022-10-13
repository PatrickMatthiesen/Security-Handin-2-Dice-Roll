package main

import (
	"context"
	"crypto/tls"
	"crypto/x509"
	"flag"
	"fmt"
	"log"
	"net"
	"os"

	// this has to be the same as the go.mod module,
	// followed by the path to the folder the proto file is in.
	gRPC "github.com/PatrickMatthiesen/DiceRoll/proto"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
)

type Server struct {
	gRPC.UnimplementedDiceRollServiceServer        // You need this line if you have a server
	name                                    string // Not required but useful if you want to name your server
	port                                    string // Not required but useful if your server needs to know what port it's listening to

	randomA    int32
	commitment int32
}

// flags are used to get arguments from the terminal. Flags take a value, a default value and a description of the flag.
// to use a flag then just add it as an argument when running the program.
var port = flag.String("port", "5400", "Server port") // set with "-port <port>" in terminal

func main() {
	flag.Parse()
	fmt.Println("--- Bob is waking up ---")
	log.Printf("Bob attempts to create listener on port %s\n", *port)

	// creds, _ := credentials.NewServerTLSFromFile("keys/server-cert.pem", "keys/server-key.pem")

	creds := credentials.NewTLS(getTLSConfig())

	grpcServer := grpc.NewServer(grpc.Creds(creds))

	list, err := net.Listen("tcp", "localhost:5400")
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	server := &Server{
		port:       *port,
		name:       "Bob",
		randomA:    0,
		commitment: 0,
	}

	gRPC.RegisterDiceRollServiceServer(grpcServer, server) //Registers the server to the gRPC server.

	log.Printf("Bob is Listening at %v\n", list.Addr())

	if err := grpcServer.Serve(list); err != nil {
		log.Fatalf("failed to serve %v", err)
	}
	// code here is unreachable because grpcServer.Serve occupies the current thread.
}

func (s *Server) CommitRoll(cxt context.Context, req *gRPC.Commitment) (*gRPC.CommitmentResponse, error) {
	log.Printf("Bob: Received commitment from Alice: %d\n", req.GetCommitment())
	s.commitment = req.GetCommitment()
	return &gRPC.CommitmentResponse{Random: 2}, nil
}

func getTLSConfig() *tls.Config {
	cert, err := tls.LoadX509KeyPair("keys/server_cert.pem", "keys/server_key.pem")
	if err != nil {
		log.Fatalf("failed to load key pair: %s", err)
	}

	ca := x509.NewCertPool()
	caFilePath := "keys/client_ca_cert.pem"
	caBytes, err := os.ReadFile(caFilePath)
	if err != nil {
		log.Fatalf("failed to read ca cert %q: %v", caFilePath, err)
	}
	if ok := ca.AppendCertsFromPEM(caBytes); !ok {
		log.Fatalf("failed to parse %q", caFilePath)
	}

	return &tls.Config{
		ClientAuth:   tls.RequireAndVerifyClientCert,
		Certificates: []tls.Certificate{cert},
		ClientCAs:    ca,
	}
}
