package main

import (
	// "crypto/tls"
	"context"
	"flag"
	"fmt"
	"log"
	"net"

	// "net"

	"go.step.sm/crypto/tlsutil"

	// this has to be the same as the go.mod module,
	// followed by the path to the folder the proto file is in.
	gRPC "github.com/PatrickMatthiesen/DiceRoll/proto"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
)

type Server struct {
	gRPC.UnimplementedDiceRollServiceServer        // You need this line if you have a server
	name                             string // Not required but useful if you want to name your server
	port                             string // Not required but useful if your server needs to know what port it's listening to

	randomA int32
	commitment int32
}

// flags are used to get arguments from the terminal. Flags take a value, a default value and a description of the flag.
// to use a flag then just add it as an argument when running the program.
var port = flag.String("port", "5400", "Server port")           // set with "-port <port>" in terminal



func main() {
	flag.Parse()
	fmt.Println(".:Bob is waking up:.")
	log.Printf("Bob attempts to create listener on port %s\n", *port)

	c, err := tlsutil.NewServerCredentialsFromFile("keys/cointoss.pem", "keys/cointoss.key")
    if err != nil {
        log.Fatal("error creating server credentials: ", err)
    }

    opts := []grpc.ServerOption{
        grpc.Creds(credentials.NewTLS(c.TLSConfig())),
    }

    list, err := net.Listen("tcp", "localhost:5400")
    if err != nil {
        log.Fatalf("failed to listen: %v", err)
    }

	// cert, err := tls.LoadX509KeyPair("keys/cointoss.pem", "keys/cointoss.key")
    // if err != nil {
    //     log.Fatal(err)
    // }

	// makes gRPC server using the options
	// you can add options here if you want or remove the options part entirely
	// opts := []grpc.ServerOption {
	// 	grpc.Creds(credentials.NewServerTLSFromCert(&cert)),
	// }
	grpcServer := grpc.NewServer(opts...)

	server := &Server{
		port:           *port,
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
	return &gRPC.CommitmentResponse{ Random: 2}, nil
}