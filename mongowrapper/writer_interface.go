package main

//func (s *usersRepository) Project(ctx context.Context, pipeline mongo.Pipeline)  ([]*User, error) {
//	if cursor, err := s.c.Aggregate(ctx, pipeline); err != nil {
//		if err == mongo.ErrNoDocuments {
//			return nil, nil
//		}
//	}else {
//		for cursor.Next(ctx) {
//
//		}
//	}
//}
var writerIface = `{{ $tick := "` + "`" + `" }}
import (
	"context"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
	{{ if not .DbVar }}"net/url"{{end}}
	"os"
	"sync"
	"time"
)

type {{ .Typ }}Repository interface {
	Ping() error
	FindOne(ctx context.Context, findQuery bson.M) (*{{ .Typ }}, error)
	FindOneById(ctx context.Context, id string) (*{{ .Typ }}, error)
	FindMany(ctx context.Context, findQuery bson.M, sort bson.D, skip, limit int64) ([]*{{ .Typ }}, error)
	InsertOne(ctx context.Context, record *{{ .Typ }}) (InsertedID primitive.ObjectID, err error)
	InsertMany(ctx context.Context, records []*{{ .Typ }}) (InsertedID []primitive.ObjectID, err error)
	UpdateOne(ctx context.Context, findQuery, updateQuery bson.M) (matched bool, modified bool, err error)
	UpdateOneById(ctx context.Context, id string, updateQuery bson.M) (matched bool, modified bool, err error)
	UpdateOneFluent(ctx context.Context, findQuery bson.M, updater {{ .Typ }}Updater) (matched bool, modified bool, err error)
	UpdateOneByIdFluent(ctx context.Context, id string, updater {{ .Typ }}Updater) (matched bool, modified bool, err error)
	DeleteOne(ctx context.Context, findQuery bson.M) (isDeleted bool, err error)
	DeleteOneById(ctx context.Context, id string) (isDeleted bool, err error)
	DeleteMany(ctx context.Context, findQuery bson.M) (delete int64, err error)
	Watch(pipeline mongo.Pipeline, ch chan<- {{ .Typ }}ChangeEvent) error
{{range .Aggregations}}
	{{ if eq .Name $.Name }}AggregateTo{{ .Typ }}(ctx context.Context, pipeline mongo.Pipeline, limit int)  ([]*{{ .Typ }}, error){{ end }}
{{end}}
}

type {{ .Name }}Repository struct {
	client *mongo.Client
	ctx    context.Context
	c      *mongo.Collection
}

func New{{ .Typ }}RepositoryDefault(ctx context.Context) {{ .Typ }}Repository {
	{{if $.CsVar }}
	cs := os.Getenv("{{ $.CsVar }}")
	if cs == "" {
		cs = "{{ $.Cs }}"
	}
	client := newClient(ctx, cs)
	{{else}}
	client := newClient(ctx, "{{ $.Cs }}")
	{{end}}
	return &{{ .Name }}Repository{
		client: client,
		ctx:    ctx,
		c:      database.Collection("{{ .Name }}"),
	}
}

func New{{ .Typ }}Repository(ctx context.Context, cs string) {{ .Typ }}Repository {
	u, err := url.Parse(cs)
	if err != nil {
		panic(err)
	}
	client := newClient(ctx, u.String())
	return &{{ .Name }}Repository{
		client: client,
		ctx:    ctx,
		c:      database.Collection("{{ .Name }}"),
	}
}

func (s *{{ .Name }}Repository) Ping() error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	return s.client.Ping(ctx, readpref.Primary())
}

func (s *{{ .Name }}Repository) FindMany(ctx context.Context, findQuery bson.M, sort bson.D, skip, limit int64) ([]*{{ .Typ }}, error) {
	opts := &options.FindOptions{}
	opts.SetLimit(limit)
	opts.SetSkip(skip)
	opts.SetSort(sort)

	if cursor, err := s.c.Find(ctx, findQuery, opts); err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, nil
		}
		return nil, err
	} else {
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
		return nil, err
	}
	var r {{ .Typ }}
	if err := s.c.FindOne(ctx, bson.M{"_id": prim}).Decode(&r); err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, nil
		}
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

func (s *{{ .Name }}Repository) UpdateOne(ctx context.Context, findQuery, updateQuery bson.M) (matched bool, modified bool, err error) {
	if res, err := s.c.UpdateOne(ctx, findQuery, updateQuery); err != nil {
		return false, false, err
	} else {
		return res.MatchedCount > 0, res.ModifiedCount > 0, nil
	}
}

func (s *{{ .Name }}Repository) UpdateOneById(ctx context.Context, id string, updateQuery bson.M) (matched bool, modified bool, err error) {
	if bsonId, err := primitive.ObjectIDFromHex(id); err != nil {
		return false, false, err
	} else {
		if res, err := s.c.UpdateOne(ctx, bson.M{"_id": bsonId}, updateQuery); err != nil {
			return false, false, err
		} else {
			return res.MatchedCount > 0, res.ModifiedCount > 0, nil
		}
	}
}

func (s *{{ .Name }}Repository) UpdateOneFluent(ctx context.Context, findQuery bson.M, updater {{ .Typ }}Updater) (matched bool, modified bool, err error) {
	if res, err := s.c.UpdateOne(ctx, findQuery, updater.(*{{ .Name }}_updater).compile()); err != nil {
		return false, false, err
	} else {
		return res.MatchedCount > 0, res.ModifiedCount > 0, nil
	}
}

func (s *{{ .Name }}Repository) UpdateOneByIdFluent(ctx context.Context, id string, updater {{ .Typ }}Updater) (matched bool, modified bool, err error) {
	if bsonId, err := primitive.ObjectIDFromHex(id); err != nil {
		return false, false, err
	} else {
		if res, err := s.c.UpdateOne(ctx, bson.M{"_id": bsonId}, updater.(*{{ .Name }}_updater).compile()); err != nil {
			return false, false, err
		} else {
			return res.MatchedCount > 0, res.ModifiedCount > 0, nil
		}
	}
}

func (s *{{ .Name }}Repository) DeleteOne(ctx context.Context, findQuery bson.M) (isDeleted bool, err error) {
	res, err := s.c.DeleteOne(ctx, findQuery)
	if err != nil {
		return false, err
	}
	return res.DeletedCount > 0, nil
}

func (s *{{ .Name }}Repository) DeleteOneById(ctx context.Context, id string) (isDeleted bool, err error) {
	if bsonId, err := primitive.ObjectIDFromHex(id); err != nil {
		return false, err
	} else {
		res, err := s.c.DeleteOne(ctx, bson.M{"_id":bsonId})
		if err != nil {
			return false, err
		}
		return res.DeletedCount > 0, nil
	}
}

func (s *{{ .Name }}Repository) DeleteMany(ctx context.Context, findQuery bson.M) (delete int64, err error) {
	res, err := s.c.DeleteMany(ctx, findQuery)
	if err != nil {
		return 0, err
	}
	return res.DeletedCount, nil
}


func (s *{{ .Name }}Repository)Watch(pipeline mongo.Pipeline, ch chan<- {{ .Typ }}ChangeEvent) error {
	updateLookup := options.UpdateLookup
	opts1 := &options.ChangeStreamOptions{
		FullDocument: &updateLookup,
	}
	stream, err := s.c.Watch(s.ctx, pipeline, opts1)
	if err != nil {
		return err
	}
	var wg sync.WaitGroup
	go func() {
		wg.Add(1)
		defer wg.Done()
		for {
			select {
			case <-s.ctx.Done():
				return
			default:
				iterate{{ .Typ }}ChangeStream(s.ctx, stream, ch)
			}
		}
	}()
	wg.Wait()
	return nil
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
