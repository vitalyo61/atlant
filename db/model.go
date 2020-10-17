package db

import (
	"context"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func CreateIndexes(collection *mongo.Collection, timeout time.Duration) error {
	indexView := collection.Indexes()
	opts := options.ListIndexes().SetMaxTime(timeout)
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	cursor, err := indexView.List(ctx, opts)
	if err != nil {
		return err
	}
	defer cursor.Close(ctx)

	haveInds := make(map[string]struct{})

	for cursor.Next(ctx) {
		var result bson.M
		if err := cursor.Decode(&result); err != nil {
			return err
		}
		haveInds[result["name"].(string)] = struct{}{}
	}
	if err := cursor.Err(); err != nil {
		return err
	}

	models := make([]mongo.IndexModel, 0)

	for _, ind := range []string{"price", "count", "date"} {
		if _, ok := haveInds[ind]; ok {
			continue
		}

		name := ind
		models = append(models, mongo.IndexModel{
			Keys: bson.D{{ind, 1}},
			Options: &options.IndexOptions{
				Name: &name,
			},
		})
	}

	if len(models) > 0 {
		ctx, cancel := context.WithTimeout(context.Background(), timeout)
		defer cancel()
		if _, err := collection.Indexes().CreateMany(ctx, models,
			options.CreateIndexes().SetMaxTime(time.Duration(len(models))*timeout)); err != nil {
			return err
		}
	}

	return nil
}

type Product struct {
	Name      string               `bson:"_id"`
	Price     primitive.Decimal128 `bson:"price"`
	Count     int32                `bson:"count"`
	Timestamp time.Time            `bson:"date"`
}

type ProductForUpdate struct {
	Name, Price string
}

func ProductsUpdate(db *Database, products []*ProductForUpdate) error {
	models := make([]mongo.WriteModel, 0)
	var (
		err error
		d   primitive.Decimal128
	)

	for _, p := range products {
		if d, err = primitive.ParseDecimal128(p.Price); err != nil {
			log.Printf("Error parse price for '%s': %s\n", p.Name, err)
			continue
		}

		models = append(models,
			mongo.NewUpdateOneModel().
				SetFilter(bson.M{"_id": p.Name}).
				SetUpdate(bson.M{
					"$currentDate": bson.M{"date": bson.M{"$type": "timestamp"}},
					"$set":         bson.M{"price": d},
					"$inc":         bson.M{"count": 1},
				}).SetUpsert(true))
	}

	opts := options.BulkWrite().SetOrdered(false)
	ctx, cancel := context.WithTimeout(context.Background(), db.Timeout)
	defer cancel()

	_, err = db.Product.BulkWrite(ctx, models, opts)

	return err
}

func ProductUpdate(db *Database, name, price string) error {
	d, err := primitive.ParseDecimal128(price)
	if err != nil {
		return err
	}

	opts := options.Update().SetUpsert(true)
	filter := bson.M{"_id": name}
	update := bson.M{
		"$currentDate": bson.M{"date": bson.M{"$type": "timestamp"}},
		"$set":         bson.M{"price": d},
		"$inc":         bson.M{"count": 1},
	}

	ctx, cancel := context.WithTimeout(context.Background(), db.Timeout)
	defer cancel()

	_, err = db.Product.UpdateOne(ctx, filter, update, opts)

	return err
}

func ProductsList(db *Database, limit, skip int64, sort string, order int) ([]*Product, error) {
	if sort == "name" {
		sort = "_id"
	}
	opts := options.Find().
		SetLimit(limit).
		SetSkip(skip).
		SetSort(bson.D{{sort, order}})

	products := make([]*Product, 0)

	ctx, cancelF := context.WithTimeout(context.Background(), db.Timeout)
	defer cancelF()
	cursor, err := db.Product.Find(ctx, bson.D{}, opts)
	if err != nil {
		return products, err
	}

	ctx, cancelC := context.WithTimeout(context.Background(), db.Timeout)
	defer cancelC()
	err = cursor.All(ctx, &products)
	return products, err
}
