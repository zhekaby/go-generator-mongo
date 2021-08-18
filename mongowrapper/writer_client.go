package main

var writerClient = `{{ $tick := "` + "`" + `" }}
import (
	"context"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
	"sync"
)

var once sync.Once
var client *mongo.Client

func newClient(ctx context.Context, cs string) *mongo.Client {
	once.Do(func() {
		c, err := mongo.Connect(ctx, options.Client().ApplyURI(cs))
		if err != nil {
			panic(err)
		}

		if err = c.Ping(ctx, readpref.Primary()); err != nil {
			panic(err)
		}
		client = c
	})
	return client
}


`
