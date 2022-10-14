package main

import (
	"bytes"
	"context"
	"crypto/sha1"
	"flag"
	"fmt"
	"log"
	"math/rand"
	"net"

	"github.com/PatrickMatthiesen/DiceRoll/util"

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

	commitment []byte
	randomB    int64
}

// flags are used to get arguments from the terminal. Flags take a value, a default value and a description of the flag.
// to use a flag then just add it as an argument when running the program.
var port = flag.String("port", "5400", "Server port") // set with "-port <port>" in terminal

func main() {
	flag.Parse()
	fmt.Println("--- Bob is waking up ---")
	log.Printf("Bob attempts to create listener on port %s\n", *port)

	creds := credentials.NewTLS(util.GetServerTLSConfig())

	grpcServer := grpc.NewServer(grpc.Creds(creds))

	list, err := net.Listen("tcp", "localhost:5400")
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	server := &Server{
		port:       *port,
		name:       "Bob",
		randomB:    rand.Int63(),
		commitment: nil,
	}

	gRPC.RegisterDiceRollServiceServer(grpcServer, server) //Registers the server to the gRPC server.

	log.Printf("Bob is Listening at %v\n", list.Addr())

	if err := grpcServer.Serve(list); err != nil {
		log.Fatalf("failed to serve %v", err)
	}
	// code here is unreachable because grpcServer.Serve occupies the current thread.
}

func (s *Server) CommitRoll(cxt context.Context, req *gRPC.Commitment) (*gRPC.CommitmentResponse, error) {
	log.Println()
	log.Println("Commitment from Alice: ", req.GetCommitment())
	s.commitment = req.GetCommitment()

	randomB := rand.Int63()
	log.Println("Bob random: ", randomB)
	s.randomB = randomB

	return &gRPC.CommitmentResponse{Random: randomB}, nil
}

func (s *Server) ValidateRoll(cxt context.Context, req *gRPC.RollValidation) (*gRPC.RollValidationResponse, error) {
	randomA := req.GetRandom()
	log.Println("Alice random:", randomA)

	sum := sha1.New().Sum([]byte(fmt.Sprint(randomA)))

	log.Println("Alice commitment:", s.commitment)
	log.Println("Bob commitment:", sum)

	if !bytes.Equal(s.commitment, sum) {
		log.Println("Roll is invalid")
		return &gRPC.RollValidationResponse{Valid: false}, nil
	}

	log.Println("Bob random cal:", s.randomB)
	roll := util.CalculateDiceRoll(randomA, s.randomB)

	log.Println("Bob: Roll is valid, sending result to Alice")
	log.Println("Bob: Roll is", roll)

	return &gRPC.RollValidationResponse{Valid: true, Roll: roll}, nil
}
