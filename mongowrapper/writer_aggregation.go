package main

var writerAggregation = `
import (
	"context"
	"go.mongodb.org/mongo-driver/mongo"
)

{{range .Aggregations}}
func (s *{{ .Name }}Repository) AggregateTo{{ .Typ }}(ctx context.Context, pipeline mongo.Pipeline, limit int)  ([]*{{ .Typ }}, error) {
	if cursor, err := s.c.Aggregate(ctx, pipeline); err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, nil
		} else {
			return nil, err
		}
	}else {
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
{{end}}
`
