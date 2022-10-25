package main

import (
	// "crypto/tls"
	"context"
	"crypto/sha256"
	"flag"
	"fmt"
	"log"
	"math/rand"

	// "encoding/pem"

	gRPC "github.com/PatrickMatthiesen/DiceRoll/proto"
	"github.com/PatrickMatthiesen/DiceRoll/util"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
)

// Same principle as in client. Flags allows for user specific arguments/values
var clientsName = flag.String("name", "default", "Senders name")
var serverPort = flag.String("server", "5400", "Tcp server")

var server gRPC.DiceRollServiceClient //the server

func main() {
	//parse flag/arguments
	flag.Parse()

	fmt.Println("--- CLIENT APP ---")

	//log to file instead of console
	//setLog()

	//connect to server and close the connection when program closes
	ConnectToServer()

	//start the biding
}

// connect to server
func ConnectToServer() {
	creds := credentials.NewTLS(util.GetClientTLSConfig())

	log.Printf("client %s: Attempts to dial on port %s\n", *clientsName, *serverPort)
	conn, err := grpc.Dial("localhost:5400",
		grpc.WithTransportCredentials(creds),
		grpc.WithBlock(),
	)
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer conn.Close()

	server = gRPC.NewDiceRollServiceClient(conn)
	log.Println("the connection is: ", conn.GetState().String())

	for {
		rollTheDice()
	}
}

func rollTheDice() {
	//commit to a roll
	randomA, randomB := commitRoll()
	//roll the dice
	roll := ValidateRoll(randomA, randomB)

	if roll == 0 {
		log.Println("Roll was invalid")
		return
	}

	log.Println("Roll was", roll)
}

func commitRoll() (int64, int64) {
	randomA := rand.Int63()
	commitment := sha256.New().Sum([]byte(fmt.Sprint(randomA)))

	responce, err := server.CommitRoll(context.Background(), &gRPC.Commitment{Commitment: commitment[:]})
	if err != nil {
		log.Fatalf("could not commit roll: %v", err)
	}
	log.Println("Alice random:", randomA)
	log.Println("Bob random:", responce.Random)

	return randomA, responce.Random
}

func ValidateRoll(randomA int64, randomB int64) int64 {
	responce, err := server.ValidateRoll(context.Background(), &gRPC.RollValidation{Random: randomA})
	if err != nil {
		log.Fatalln("could not validate roll: ", err)
	}
	if !responce.Valid {
		return 0
	}

	return util.CalculateDiceRoll(randomA, randomB)
}

