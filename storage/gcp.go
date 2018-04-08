package storage

import (
	"bytes"
	"fmt"
	"io"
	"time"

	"golang.org/x/net/context"

	gcs "cloud.google.com/go/storage"
	"github.com/golang/protobuf/proto"
	"github.com/nov1n/dumpertijdmachine/types"
)

type GCS struct {
	ctx       context.Context
	db        *gcs.Client
	bucket    string
	projectID string
}

func NewGCS(projectID, bucket string) (*GCS, error) {
	ctx := context.Background()
	db, err := gcs.NewClient(ctx)
	if err != nil {
		return nil, fmt.Errorf("create client: %s", err)
	}

	return &GCS{
		ctx:       ctx,
		bucket:    bucket,
		db:        db,
		projectID: projectID,
	}, nil
}

func (g *GCS) PutDay(date time.Time, day *types.Day) error {
	data, err := proto.Marshal(day)
	if err != nil {
		return fmt.Errorf("marshal day: %s", err)
	}

	bkt := g.db.Bucket(g.bucket)
	if err := bkt.Create(g.ctx, g.projectID, nil); err != nil {
		// pass
	}

	obj := bkt.Object(KeyForDate(date))
	w := obj.NewWriter(g.ctx)
	if _, err := fmt.Fprintf(w, "%s", data); err != nil {
		return fmt.Errorf("object write: %s", err)
	}
	if err := w.Close(); err != nil {
		return fmt.Errorf("object write close: %s", err)
	}

	return nil
}

func (g *GCS) GetDay(date time.Time) (*types.Day, error) {
	dateKey := KeyForDate(date)

	day := &types.Day{}

	bkt := g.db.Bucket(g.bucket)
	obj := bkt.Object(dateKey)
	r, err := obj.NewReader(g.ctx)
	if err != nil {
		return nil, fmt.Errorf("object read %s: %s", dateKey, err)
	}
	defer r.Close()

	data := &bytes.Buffer{}
	if _, err := io.Copy(data, r); err != nil {
		return nil, fmt.Errorf("object read copy: %s", err)
	}

	if data == nil {
		return nil, fmt.Errorf("day not found for date %s", date)
	}

	err = proto.Unmarshal(data.Bytes(), day)
	if err != nil {
		return nil, fmt.Errorf("unmarshal day: %s", err)
	}

	return day, nil
}

func (g *GCS) Close() error {
	return g.db.Close()
}
