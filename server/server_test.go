package server

import (
	"context"
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/vitalyo61/atlant/config"
	"github.com/vitalyo61/atlant/db"
	pb "github.com/vitalyo61/atlant/grpc"
)

func TestServer(t *testing.T) {
	address := os.Getenv("MONGO")
	if address == "" {
		address = "mongodb://localhost:27017"
	}

	dbase, err := db.New(&config.DB{
		Address: address,
		Name:    "atlant",
		Timeout: 5,
	})
	if !assert.NoError(t, err) {
		return
	}

	defer func() {
		err = dbase.Close()
		assert.NoError(t, err)
	}()

	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		for i := 97; i < 123; i++ {
			fmt.Fprintf(w, "%c;%d.00\n", i, 123-i)
		}
	}))
	defer ts.Close()

	s := NewServer(dbase)

	ctx := context.Background()
	reqF := &pb.FetchRequest{Url: "ts.URL"}
	_, err = s.Fetch(ctx, reqF)
	assert.Error(t, err)

	reqF.Url = ts.URL
	_, err = s.Fetch(ctx, reqF)
	assert.NoError(t, err)

	reqL := &pb.ListRequest{
		Sorting: &pb.Sorting{
			Type: pb.Sorting_DESC,
		},
	}
	resp, err := s.List(ctx, reqL)
	assert.NoError(t, err)

	for i, p := range resp.Products {
		assert.Equal(t, fmt.Sprintf("%c", 122-i), p.Name)
	}
}
