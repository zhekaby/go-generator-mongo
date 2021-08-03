package main

var writerIface = `{{ $tick := "` + "`" + `" }}
import (
	"context"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"net/url"
	"os"
	"sync"
)

type {{ .Typ }}Repository interface {
	FindOne(ctx context.Context, findQuery bson.M) (*{{ .Typ }}, error)
	FindOneById(ctx context.Context, id string) (*User, error)
	FindMany(ctx context.Context, findQuery bson.M, skip, limit int64) ([]*{{ .Typ }}, error)
	InsertOne(ctx context.Context, record *{{ .Typ }}) (InsertedID primitive.ObjectID, err error)
	InsertMany(ctx context.Context, records []*{{ .Typ }}) (InsertedID []primitive.ObjectID, err error)
	DeleteOne(ctx context.Context, findQuery bson.M) (delete int64, err error)
	DeleteMany(ctx context.Context, findQuery bson.M) (delete int64, err error)
	Watch(pipeline mongo.Pipeline) (<-chan {{ .Typ }}ChangeEvent, error)
}

type {{ .Name }}Repository struct {
	client *mongo.Client
	ctx    context.Context
	c      *mongo.Collection
}

func New{{ .Typ }}RepositoryDefault(ctx context.Context) {{ .Typ }}Repository {
	cs := os.Getenv("{{ $.CsVar }}")
	if cs == "" {
		cs = "{{ $.Cs }}"
	}

	return New{{ .Typ }}Repository(ctx, cs)
}


func New{{ .Typ }}Repository(ctx context.Context, cs string) {{ .Typ }}Repository {
	u, err := url.Parse(cs)
	if err != nil {
		panic(err)
	}

	client, err := mongo.Connect(ctx, options.Client().ApplyURI(cs))
	if err != nil {
		panic(err)
	}
	database := client.Database(u.Path[1:])
	return &{{ .Name }}Repository{
		client: client,
		ctx:    ctx,
		c:      database.Collection("{{ .Name }}"),
	}
}

func (s *{{ .Name }}Repository) FindMany(ctx context.Context, findQuery bson.M, skip, limit int64) ([]*{{ .Typ }}, error) {
	opts := &options.FindOptions{}
	opts.SetLimit(limit)
	opts.SetSkip(skip)

	cursor, err := s.c.Find(ctx, findQuery, opts)
	if err != nil {
		return nil, err
	}
	records := make([]*{{ .Typ }}, 0, limit)
	for cursor.Next(ctx) {
		t := {{ .Typ }}{}
		err := cursor.Decode(&t)
		if err != nil {
			return records, err
		}
		records = append(records, &t)
	}

	return records, nil
}

func (s *{{ .Name }}Repository) FindOne(ctx context.Context, findQuery bson.M) (*{{ .Typ }}, error) {
	var r {{ .Typ }}
	if err := s.c.FindOne(ctx, findQuery).Decode(&r); err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, nil
		}
		return nil, err
	}
	return &r, nil
}

func (s *{{ .Name }}Repository) FindOneById(ctx context.Context, id string) (*{{ .Typ }}, error) {
	prim, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, nil
		}
		return nil, err
	}
	var r {{ .Typ }}
	if err := s.c.FindOne(ctx, bson.M{"_id": prim}).Decode(&r); err != nil {
		return nil, err
	}
	return &r, nil
}

func (s *{{ .Name }}Repository) InsertOne(ctx context.Context, record *{{ .Typ }}) (InsertedID primitive.ObjectID, err error) {
	res, err := s.c.InsertOne(ctx, record)
	if err != nil {
		return primitive.ObjectID{}, err
	}
	return res.InsertedID.(primitive.ObjectID), err
}

func (s *{{ .Name }}Repository) InsertMany(ctx context.Context, records []*{{ .Typ }}) (InsertedID []primitive.ObjectID, err error) {
	data := make([]interface{}, len(records))
	for i := range records {
		data[i] = records[i]
	}
	res, err := s.c.InsertMany(ctx, data)
	if err != nil {
		return []primitive.ObjectID{}, err
	}
	ids := make([]primitive.ObjectID, len(res.InsertedIDs))
	for i := range res.InsertedIDs {
		ids[i] = res.InsertedIDs[i].(primitive.ObjectID)
	}
	return ids, err
}

func (s *{{ .Name }}Repository) DeleteOne(ctx context.Context, findQuery bson.M) (delete int64, err error) {
	res, err := s.c.DeleteOne(ctx, findQuery)
	if err != nil {
		return 0, err
	}
	return res.DeletedCount, nil
}

func (s *{{ .Name }}Repository) DeleteMany(ctx context.Context, findQuery bson.M) (delete int64, err error) {
	res, err := s.c.DeleteMany(ctx, findQuery)
	if err != nil {
		return 0, err
	}
	return res.DeletedCount, nil
}


func (s *{{ .Name }}Repository) Watch(pipeline mongo.Pipeline) (<-chan {{ .Typ }}ChangeEvent, error) {
	updateLookup := options.UpdateLookup
	opts1 := &options.ChangeStreamOptions{
		FullDocument: &updateLookup,
	}
	stream, err := s.c.Watch(s.ctx, pipeline, opts1)
	if err != nil {
		return nil, err
	}
	var wg sync.WaitGroup
	ch := make(chan {{ .Typ }}ChangeEvent)
	go func() {
		wg.Add(1)
		defer wg.Done()
		for {
			select {
			case <-s.ctx.Done():
				close(ch)
			default:
				iterate{{ .Typ }}ChangeStream(s.ctx, stream, ch)
			}
		}
	}()
	wg.Wait()
	return ch, nil
}

func iterate{{ .Typ }}ChangeStream(ctx context.Context, stream *mongo.ChangeStream, ch chan<- {{ .Typ }}ChangeEvent) {
	for stream.Next(ctx) {
		var data {{ .Typ }}ChangeEvent
		if err := stream.Decode(&data); err != nil {
			continue
		}
		ch <- data
	}
}

type {{ .Typ }}ChangeEvent struct {
	ID struct {
		Data string {{ $tick }}bson:"_data"{{ $tick }}
	} {{ $tick }}bson:"_id"{{ $tick }}
	OperationType string              {{ $tick }}bson:"operationType"{{ $tick }}
	ClusterTime   primitive.Timestamp {{ $tick }}bson:"clusterTime"{{ $tick }}
	FullDocument  *{{ .Typ }}         {{ $tick }}bson:"fullDocument"{{ $tick }}
	DocumentKey   struct {
		ID primitive.ObjectID {{ $tick }}bson:"_id"{{ $tick }}
	} {{ $tick }}bson:"documentKey"{{ $tick }}
	Ns struct {
		Db   string {{ $tick }}bson:"db"{{ $tick }}
		Coll string {{ $tick }}bson:"coll"{{ $tick }}
	} {{ $tick }}bson:"ns"{{ $tick }}
}
`
