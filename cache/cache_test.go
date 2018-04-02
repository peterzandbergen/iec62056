package cache

import (
	"testing"
	"time"

	"github.com/peterzandbergen/iec62056/model"
)

const testDB = "./.testdb"

func TestOpenDB(t *testing.T) {
	var c *Cache

	c, err := Open(testDB)
	if err != nil {
		t.Fatalf("Error opening database: %s", err.Error())
	}
	defer c.Close()
}

var putM = &model.Measurement{
	Time:           time.Date(2018, 3, 21, 12, 30, 0, 0, time.Local),
	ManufacturerID: "T01",
	Identification: "sddsadfada",
	Readings: []model.DataSet{
		{
			Address: "1.2.3",
			Value:   "500",
			Unit:    "",
		},
	},
}

func TestPutMeasurement(t *testing.T) {
	var c *Cache

	c, err := Open(testDB)
	if err != nil {
		t.Fatalf("Error opening database: %s", err.Error())
	}
	defer c.Close()

	err = c.Put(putM)
	if err != nil {
		t.Errorf("error saving m: %s", err.Error())
	}
}

var getM = &model.Measurement{
	Time:           time.Date(2018, 3, 21, 12, 30, 0, 0, time.Local),
	ManufacturerID: "T02",
	Identification: "SSSSS",
	Readings: []model.DataSet{
		{
			Address: "1.2.3",
			Value:   "500",
			Unit:    "",
		},
	},
}

func TestGetMeasurement(t *testing.T) {
	var c *Cache

	c, err := Open(testDB)
	if err != nil {
		t.Fatalf("Error opening database: %s", err.Error())
	}
	defer c.Close()

	err = c.Put(getM)
	if err != nil {
		t.Errorf("error saving m: %s", err.Error())
	}

	m, err := c.Get(key(getM))
	if err != nil {
		t.Fatalf("error get: %s", err.Error())
	}
	t.Logf("Found m with id: %s", m.Identification)
}

// "Mon Jan 2 15:04:05 -0700 MST 2006"
// "2006 Jan "
var getMultM = []*model.Measurement{
	{
		Time:           time.Date(2018, 3, 21, 12, 30, 0, 0, time.Local),
		ManufacturerID: "T03",
		Identification: "SSSSS",
		Readings: []model.DataSet{
			{
				Address: "1.2.3",
				Value:   "500",
				Unit:    "",
			},
		},
	},
	{
		Time:           time.Date(2018, 3, 21, 12, 35, 0, 0, time.Local),
		ManufacturerID: "T03",
		Identification: "SSSSS",
		Readings: []model.DataSet{
			{
				Address: "1.2.3",
				Value:   "500",
				Unit:    "",
			},
		},
	},
	{
		Time:           time.Date(2018, 3, 21, 12, 40, 0, 0, time.Local),
		ManufacturerID: "T03",
		Identification: "SSSSS",
		Readings: []model.DataSet{
			{
				Address: "1.2.3",
				Value:   "500",
				Unit:    "",
			},
		},
	},
}

func TestGetNMeasurement(t *testing.T) {
	var c *Cache

	c, err := Open(testDB)
	if err != nil {
		t.Fatalf("Error opening database: %s", err.Error())
	}
	defer c.Close()

	for _, m := range getMultM {
		c.Put(m)
	}

	n := 1000
	ms, err := c.GetN(n)
	if err != nil {
		t.Fatalf("error get: %s", err.Error())
	}
	if len(ms) < len(getMultM) {
		t.Errorf("expected at least %d messages, received %d", len(getMultM), len(ms))
	}
	t.Logf("Found %d messages", len(ms))
}
