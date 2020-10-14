package db

import (
	"context"
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

		t.Logf("%+v", result)
		prs := products[result.Name]
		assert.Equal(t, int(result.Count), len(prs))
		assert.Equal(t, result.Price.String(), prs[len(prs)-1])
	}

	productsBulk := []*ProductForUpdate{
		&ProductForUpdate{
			Name:  "a",
			Price: "10.40",
		},
		&ProductForUpdate{
			Name:  "b",
			Price: "20.50",
		},
		&ProductForUpdate{
			Name:  "c",
			Price: "30.20",
		},
	}

	err = ProductsUpdate(db, productsBulk)
	if !assert.NoError(t, err) {
		return
	}

	ctx, cancel = context.WithTimeout(context.Background(), db.Timeout)
	defer cancel()

	cur, err = db.Product.Find(ctx, bson.D{})
	if !assert.NoError(t, err) {
		return
	}

	defer cur.Close(ctx)

	var count int

	for cur.Next(ctx) {
		var result Product
		err := cur.Decode(&result)
		if !assert.NoError(t, err) {
			return
		}

		t.Logf("%+v", result)
		pr := productsBulk[count]
		prs := products[pr.Name]
		assert.Equal(t, int(result.Count), len(prs)+1)
		assert.Equal(t, result.Price.String(), pr.Price)

		count++
	}

}
