package main

import (
	"context"
	"net/http"

	pb "github.com/vitalyo61/atlant/grpc"
)

type server struct {
	pb.UnimplementedProductServer
}

func (s *server) Fetch(ctx context.Context, req *pb.FetchRequest) (*pb.FetchResponce, error) {
	resp := new(pb.FetchResponce)

	r, err := http.Get(req.Url)
	if err != nil {
		resp.Error = "bad URL"
		return resp, nil
	}

	defer r.Body.Close()

	// scanner := bufio.NewScanner(file)
	// for scanner.Scan() {
	// 	scanner.Text()
	// }

	return resp, nil
}
func (s *server) List(context.Context, *pb.ListRequest) (*pb.ListResponce, error) {
	// return nil, status.Errorf(codes.Unimplemented, "method List not implemented")
	return nil, nil
}

func main() {
}
