package db

import (
	"context"
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/vitalyo61/atlant/config"
	"go.mongodb.org/mongo-driver/bson"
)

func TestDB(t *testing.T) {
	db, err := New(&config.DB{
		Address: "mongodb://localhost:30000",
		Name:    "atlant",
		Timeout: 5,
	})
	if !assert.NoError(t, err) {
		return
	}

	defer func() {
		err = db.Close()
		assert.NoError(t, err)
	}()

	assert.Equal(t, db.Product.Name(), "product")

	err = CreateIndexes(db.Product, db.Timeout)
	assert.NoError(t, err)

	products := map[string][]string{
		"a": []string{"10.10", "10.20", "10.30"},
		"b": []string{"20.10", "20.20"},
		"c": []string{"30.10"},
	}

	for p, prs := range products {
		for _, pr := range prs {
			err = ProductUpdate(db, p, pr)
			if !assert.NoError(t, err) {
				return
			}
		}
	}

	ctx, cancel := context.WithTimeout(context.Background(), db.Timeout)
	defer cancel()

	cur, err := db.Product.Find(ctx, bson.D{})
	if !assert.NoError(t, err) {
		return
	}

	defer cur.Close(ctx)

	for cur.Next(ctx) {
		var result Product
		err := cur.Decode(&result)
		if !assert.NoError(t, err) {
			return
		}

		prs := products[result.Name]
		assert.Equal(t, int(result.Count), len(prs))
		assert.Equal(t, result.Price.String(), prs[len(prs)-1])
	}

	productsBulk := make([]*ProductForUpdate, 0)

	for i := 97; i <= 122; i++ {
		productsBulk = append(productsBulk,
			&ProductForUpdate{
				Name:  fmt.Sprintf("%c", i),
				Price: fmt.Sprintf("%d.00", 123-i),
			})
	}

	err = ProductsUpdate(db, productsBulk)
	if !assert.NoError(t, err) {
		return
	}

	productsList, err := ProductsList(db, 4, 0, "name", 1)
	if !assert.NoError(t, err) {
		return
	}

	for i, p := range productsList {
		assert.Equal(t, p.Name, fmt.Sprintf("%c", i+97))
		assert.Equal(t, p.Count, int32(4-i))
		assert.Equal(t, p.Price.String(), fmt.Sprintf("%d.00", 26-i))
	}

	productsList, err = ProductsList(db, 4, 4, "name", 1)
	if !assert.NoError(t, err) {
		return
	}

	for i, p := range productsList {
		assert.Equal(t, p.Name, fmt.Sprintf("%c", i+97+4))
		assert.Equal(t, p.Price.String(), fmt.Sprintf("%d.00", 26-i-4))
	}

	productsList, err = ProductsList(db, 4, 0, "price", -1)
	if !assert.NoError(t, err) {
		return
	}

	for i, p := range productsList {
		assert.Equal(t, p.Price.String(), fmt.Sprintf("%d.00", 26-i))
	}

	productsList, err = ProductsList(db, 4, 0, "price", 1)
	if !assert.NoError(t, err) {
		return
	}

	for i, p := range productsList {
		assert.Equal(t, p.Price.String(), fmt.Sprintf("%d.00", i+1))
	}
}
