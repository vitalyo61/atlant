package server

import (
	"bufio"
	"context"
	"fmt"
	"log"
	"net/http"
	"strings"

	"github.com/vitalyo61/atlant/db"
	pb "github.com/vitalyo61/atlant/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type Server struct {
	pb.UnimplementedProductServer
	db *db.Database
}

func NewServer(db *db.Database) *Server {
	return &Server{db: db}
}

func (s *Server) Fetch(ctx context.Context, req *pb.FetchRequest) (*pb.FetchResponce, error) {
	resp := new(pb.FetchResponce)

	r, err := http.Get(req.Url)
	if err != nil {
		return resp, status.Errorf(codes.InvalidArgument, err.Error())
	}

	defer r.Body.Close()

	var (
		pr           string
		prData       []string
		productsBulk = make([]*db.ProductForUpdate, 0)
	)

	scanner := bufio.NewScanner(r.Body)
	for scanner.Scan() {
		pr = scanner.Text()
		prData = strings.Split(pr, ";")
		if len(prData) != 2 {
			return resp, status.Errorf(codes.DataLoss, fmt.Sprintf("bad data: %q", pr))
		}

		productsBulk = append(productsBulk,
			&db.ProductForUpdate{
				Name:  prData[0],
				Price: prData[1],
			})
	}

	if err = scanner.Err(); err != nil {
		return resp, status.Errorf(codes.DataLoss, err.Error())
	}

	if err = db.ProductsUpdate(s.db, productsBulk); err != nil {
		return resp, status.Errorf(codes.Internal, err.Error())
	}

	return resp, nil
}
func (s *Server) List(ctx context.Context, req *pb.ListRequest) (*pb.ListResponce, error) {
	resp := new(pb.ListResponce)

	var (
		limit, skip int64
	)
	sort := "_id"
	order := 1

	if paging := req.GetPaging(); paging != nil {
		limit = paging.GetLimit()
		skip = paging.GetSkip()
	}

	if sorting := req.GetSorting(); sorting != nil {
		switch sorting.GetField() {
		case pb.Sorting_PRICE:
			sort = "price"
		case pb.Sorting_COUNT:
			sort = "count"
		case pb.Sorting_TIMESTAMP:
			sort = "date"
		}
		switch sorting.GetType() {
		case pb.Sorting_DESC:
			order = -1
		}
	}

	log.Println(limit, skip, sort, order)

	productsList, err := db.ProductsList(s.db, limit, skip, sort, order)
	if err != nil {
		return resp, status.Errorf(codes.Internal, err.Error())
	}

	resp.Products = make([]*pb.ListResponce_Product, len(productsList))
	for i, p := range productsList {
		resp.Products[i] = &pb.ListResponce_Product{
			Name:       p.Name,
			Price:      p.Price.String(),
			Count:      p.Count,
			LastUpdate: timestamppb.New(p.Timestamp),
		}
	}

	return resp, nil
}
