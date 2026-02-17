/*
 * Copyright (C) 2025 NullSpook
 *
 * This file is part of psirng.
 *
 * psirng is free software: you can redistribute it and/or modify it under the
 * terms of the GNU Affero General Public License as published by the Free
 * Software Foundation, either version 3 of the License, or (at your option)
 * any later version.
 *
 * psirng is distributed in the hope that it will be useful, but WITHOUT ANY
 * WARRANTY; without even the implied warranty of MERCHANTABILITY or FITNESS
 * FOR A PARTICULAR PURPOSE.  See the GNU Affero General Public License for
 * more details.
 *
 * You should have received a copy of the GNU Affero General Public License
 * along with psirng.  If not, see <https://www.gnu.org/licenses/>.
 */

package main

import (
	"bytes"
	"crypto/sha256"
	"crypto/x509"
	"encoding/pem"
	"flag"
	"log"
	"net"
	"net/http"
	"os"
	psirngGrpc "psirng/grpc"
	"psirng/healthcheckers"
	"psirng/httpapi"
	"psirng/providers/rng"
	"psirng/services"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/health/grpc_health_v1"
)

var testingCertSha256 = []byte{
	0x31, 0x47, 0x1C, 0x16, 0x80, 0xAA, 0xE3, 0xE9,
	0xF7, 0x2C, 0x47, 0xC0, 0x58, 0xFF, 0x19, 0x3A,
	0xCF, 0x7B, 0x27, 0x8B, 0xA1, 0xFD, 0xF1, 0x8D,
	0x01, 0xFE, 0xFB, 0x30, 0x5A, 0xE3, 0x1A, 0x9E,
}

func main() {
	certFilePath := flag.String("cert", "", "TLS certificate file")
	keyFilePath := flag.String("key", "", "TLS key file")
	flag.Parse()

	rng, err := rng.NewQwqng()
	if err != nil {
		log.Fatalln(err)
	}

	rngService := services.NewRngService(rng)

	rngHealthChecker := healthcheckers.NewRngHealthChecker(rngService)

	rngController := httpapi.NewRngController(rngService)
	rngRouter := httpapi.NewRouter(rngController, rngHealthChecker)

	listener, err := net.Listen("tcp", ":50051")
	if err != nil {
		log.Fatalln(err)
	}

	var grpcServer *grpc.Server

	if *certFilePath != "" && *keyFilePath != "" {
		certFile, err := os.ReadFile(*certFilePath)
		if err != nil {
			log.Fatalln(err)
		}
		block, _ := pem.Decode(certFile)
		if block == nil {
			log.Fatalln("failed to decode the certificate")
		}
		cert, err := x509.ParseCertificate(block.Bytes)
		if err != nil {
			log.Fatalln(err)
		}
		certSha256 := sha256.Sum256(cert.Raw)
		if bytes.Equal(certSha256[:], testingCertSha256) {
			log.Print("WARNING: USING TESTING CERTIFICATE. DO NOT USE THIS IN PRODUCTION.\n\n")
		}

		creds, err := credentials.NewServerTLSFromFile(*certFilePath, *keyFilePath)
		if err != nil {
			log.Fatalln(err)
		}
		grpcServer = grpc.NewServer(grpc.Creds(creds), grpc.UnaryInterceptor(psirngGrpc.UnaryLoggingInterceptor()))
	} else {
		grpcServer = grpc.NewServer(grpc.UnaryInterceptor(psirngGrpc.UnaryLoggingInterceptor()))
	}

	rngServer := psirngGrpc.NewRngServerImpl(rngService)
	healthServer := psirngGrpc.NewHealthServerImpl(rngHealthChecker)

	psirngGrpc.RegisterRngServer(grpcServer, rngServer)
	grpc_health_v1.RegisterHealthServer(grpcServer, healthServer)

	if *certFilePath != "" && *keyFilePath != "" {
		log.Println("Starting HTTPS server on port 8443")
		go func() {
			err := http.ListenAndServeTLS(":8443", *certFilePath, *keyFilePath, rngRouter.MuxRouter)
			if err != nil {
				log.Fatalln(err)
			}
		}()
	} else {
		log.Println("TLS certificate and key not provided, skipping HTTPS server startup")
	}

	log.Println("Starting HTTP server on port 8080")
	go func() {
		err := http.ListenAndServe(":8080", rngRouter.MuxRouter)
		if err != nil {
			log.Fatalln(err)
		}
	}()

	log.Println("Starting gRPC server on port 50051")
	go func() {
		if err := grpcServer.Serve(listener); err != nil {
			log.Fatalln(err)
		}
	}()

	select {}
}
