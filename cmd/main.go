package main

import (
	"flag"
	"log"
	"net"
	"os"
	"os/signal"
	"syscall"

	"github.com/ilyakaznacheev/cleanenv"
	"github.com/vitalyo61/atlant/config"
	"github.com/vitalyo61/atlant/db"
	pb "github.com/vitalyo61/atlant/grpc"
	"github.com/vitalyo61/atlant/server"
	"google.golang.org/grpc"
)

func main() {
	var (
		cfg        *config.Config
		err        error
		flagConfig string
	)

	flag.StringVar(&flagConfig, "config", "", "the confg file(yaml)")

	envHelp := "Environment:"
	flag.Usage = cleanenv.Usage(new(config.Config), &envHelp, flag.Usage)

	flag.Parse()

	cfg, err = config.Get(flagConfig)
	if err != nil {
		log.Fatalf("failed to config: %v", err)
	}

	log.Printf("Loaded config: %+v\n", cfg)

	dbase, err := db.New(&cfg.DB)
	if err != nil {
		log.Fatalf("failed to database: %v", err)
	}
	defer dbase.Close()

	log.Println("Connected database")

	grpcS := server.NewServer(dbase)

	lis, err := net.Listen("tcp", ":"+cfg.GRPC.Port)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	s := grpc.NewServer()
	pb.RegisterProductServer(s, grpcS)
	log.Printf("Registered gRPC server on port %s\n", cfg.GRPC.Port)

	// interrapt application
	idleConnsClosed := make(chan struct{})
	go func() {
		sigint := make(chan os.Signal, 1)
		signal.Notify(sigint, os.Interrupt, syscall.SIGTERM)
		<-sigint

		s.GracefulStop()
		log.Println("gRPC server shutdown")
		close(idleConnsClosed)
	}()

	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
