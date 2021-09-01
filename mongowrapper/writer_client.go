package main

var writerClient = `{{ $tick := "` + "`" + `" }}
import (
	"fmt"
	"context"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
	"net/url"
	"sync"
)

var once sync.Once
var client *mongo.Client
var database *mongo.Database

func newClient(ctx context.Context, cs string) *mongo.Client {
	once.Do(func() {
		fmt.Printf("cs: %s\n", cs)
		u, err := url.Parse(cs)
		if err != nil {
			panic(err)
		}
		c, err := mongo.Connect(ctx, options.Client().ApplyURI(cs))
		if err != nil {
			panic(err)
		}

		if err = c.Ping(ctx, readpref.Primary()); err != nil {
			panic(err)
		}
		client = c
		{{ if .DbVar }}
		dbName := os.Getenv("{{ .DbVar }}");
		if dbName == "" {
			panic("{{ .DbVar }} passed but empty")
		} 
		database = client.Database(dbName)
		{{else}}
		database = client.Database(u.Path[1:])
		{{end}}
	})
	return client
}


`
