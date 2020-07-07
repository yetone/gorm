package postgres

import (
	"database/sql/driver"

	"encoding/json"
	"errors"
	"fmt"

	"github.com/jackc/pgtype"
	_ "github.com/jackc/pgx/v4"
)

type Hstore map[string]*string

// Value get value of Hstore
func (h Hstore) Value() (driver.Value, error) {
	hstore := pgtype.Hstore{Map: map[string]pgtype.Text{}}
	if len(h) == 0 {
		return nil, nil
	}

	for key, value := range h {
		s := pgtype.Text{Status: pgtype.Null}
		if value != nil {
			s.String = *value
			s.Status = pgtype.Present
		}
		hstore.Map[key] = s
	}
	return hstore.Value()
}

// Scan scan value into Hstore
func (h *Hstore) Scan(value interface{}) error {
	hstore := pgtype.Hstore{}

	if err := hstore.Scan(value); err != nil {
		return err
	}

	if len(hstore.Map) == 0 {
		return nil
	}

	*h = Hstore{}
	for k := range hstore.Map {
		if hstore.Map[k].Status == pgtype.Present {
			s := hstore.Map[k].String
			(*h)[k] = &s
		} else {
			(*h)[k] = nil
		}
	}

	return nil
}

// Jsonb Postgresql's JSONB data type
type Jsonb struct {
	json.RawMessage
}

// Value get value of Jsonb
func (j Jsonb) Value() (driver.Value, error) {
	if len(j.RawMessage) == 0 {
		return nil, nil
	}
	return j.MarshalJSON()
}

// Scan scan value into Jsonb
func (j *Jsonb) Scan(value interface{}) error {
	bytes, ok := value.([]byte)
	if !ok {
		return errors.New(fmt.Sprint("Failed to unmarshal JSONB value:", value))
	}

	return json.Unmarshal(bytes, j)
}
