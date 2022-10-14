package util

import (
	"crypto/tls"
	"crypto/x509"
	"log"
	"os"
)

func GetServerTLSConfig() *tls.Config {
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

func GetClientTLSConfig() *tls.Config {
	cert, err := tls.LoadX509KeyPair("keys/client_cert.pem", "keys/client_key.pem")
	if err != nil {
		log.Fatalf("failed to load client cert: %v", err)
	}

	ca := x509.NewCertPool()
	caFilePath := "keys/ca_cert.pem"
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

func CalculateDiceRoll(randA int64, randB int64) int64 {
	return ((randA ^ randB) % 6) + 1
}