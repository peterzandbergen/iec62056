package cache

import (

	// This project.
	"github.com/peterzandbergen/iec62056/model"

	// Stdlib
	"encoding/json"
	"errors"
	"time"

	// Vendor
	"github.com/syndtr/goleveldb/leveldb"
	"github.com/syndtr/goleveldb/leveldb/comparer"
	"github.com/syndtr/goleveldb/leveldb/opt"
)

var _ model.MeasurementRepo = &Cache{}

// type MeasurementRepo interface {
// 	Put(*Measurement) error
// 	Get(key []byte) ([]Measurement, error)
// 	GetN(n int) ([]Measurement, error)
// 	Delete(*Measurement) error
// }

var (
	ErrNotImplemented = errors.New("not implemented")
	ErrClosed         = errors.New("db closed")
)

func marshalMeasurement(m *model.Measurement) ([]byte, error) {
	// Encode to json.
	return json.Marshal(m)
}

func unmarshalMeasurement(b []byte) (*model.Measurement, error) {
	var m = &model.Measurement{}
	err := json.Unmarshal(b, m)
	if err != nil {
		return nil, err
	}
	return m, nil
}

func key(m *model.Measurement) []byte {
	return []byte(m.Time.String() + "|" + m.ManufacturerID + "|" + m.Identification)
}

func (c *Cache) Put(m *model.Measurement) error {
	if c.db == nil {
		return ErrClosed
	}
	k := key(m)
	v, err := marshalMeasurement(m)
	if err != nil {
		return err
	}
	return c.db.Put(k, v, nil)
}

func (c *Cache) Get(key []byte) (*model.Measurement, error) {
	if c.db == nil {
		return nil, ErrClosed
	}
	v, err := c.db.Get(key, nil)
	if err != nil {
		return nil, err
	}
	m, err := unmarshalMeasurement(v)
	if err != nil {
		return nil, err
	}
	return m, nil
}

func (c *Cache) GetN(n int) ([]*model.Measurement, error) {
	if c.db == nil {
		return nil, ErrClosed
	}
	// Get an iterator.
	it := c.db.NewIterator(nil, nil)
	if it.Error() != nil {
		return nil, it.Error()
	}
	defer it.Release()
	ms := make([]*model.Measurement, 0, n)
	for i := 0; i < n && it.Next(); i++ {
		if v, err := unmarshalMeasurement(it.Value()); err != nil {
		} else {
			ms = append(ms, v)
		}
	}
	return ms, nil
}

func (c *Cache) Delete(m *model.Measurement) error {
	if c.db == nil {
		return ErrClosed
	}
	k := key(m)
	err := c.db.Delete(k, nil)
	return err
}

// Cache type wraps the cache. The cache is used to store the messages until the can be sent to the cloud storage.
// This allows us to keep recording, even when network access is down.
// Cache implements the
type Cache struct {
	db       *leveldb.DB
	comparer comparer.Comparer
	options  *opt.Options
}

type CachedObject struct {
	Time   time.Time
	Object interface{}
}

func New(filename string) (*Cache, error) {
	db, err := leveldb.OpenFile(filename, nil)
	if err != nil {
		return nil, err
	}
	return &Cache{
		db: db,
	}, nil
}

func (c *Cache) Close() error {
	if c.db == nil {
		return nil
	}
	err := c.db.Close()
	if err != nil {
		return err
	}
	c.db = nil
	return nil
}
