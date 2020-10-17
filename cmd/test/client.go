package main

import (
	"context"
	"log"
	"os"
	"time"

	pb "github.com/vitalyo61/atlant/grpc"
	"google.golang.org/grpc"
)

func main() {
	address := os.Getenv("GPRS_SERVER")
	if address == "" {
		address = "localhost:50051"
	}

	log.Println("Address:", address)

	csvUrl := os.Getenv("CSV_URL")
	if csvUrl == "" {
		panic("haven't csv url")
	}

	conn, err := grpc.Dial(address, grpc.WithInsecure())
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}

	defer conn.Close()
	c := pb.NewProductClient(conn)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	_, err = c.Fetch(ctx, &pb.FetchRequest{Url: "http://" + csvUrl + "/product.csv"})
	if err != nil {
		log.Fatalf("Fetch error: %s/n", err)
	}

	reqL := &pb.ListRequest{
		Sorting: &pb.Sorting{
			Type: pb.Sorting_DESC,
		},
	}
	resp, err := c.List(ctx, reqL)
	if err != nil {
		log.Fatalf("List error %s/n", err)
	}

	for _, p := range resp.Products {
		log.Printf("%+v\n", p)
	}
}
