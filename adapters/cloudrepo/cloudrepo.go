// Package cloudrepo implements a repo that stores measurements at
// a remote service.
package cloudrepo

import (
	"fmt"
	"net/url"
	"time"

	"github.com/peterzandbergen/iec62056/model"
)

// JSONTime implements json formatting using the UTC function.
type JSONTime time.Time

// MarshalJSON implements the marshal interface.
func (t JSONTime) MarshalJSON() ([]byte, error) {
	//do your serializing here
	stamp := fmt.Sprintf("\"%s\"", time.Time(t).UTC())
	return []byte(stamp), nil
}

// Measurement resource with json format tags.
type Measurement struct {
	JSONTime       time.Time `json:"time"`
	ManufacturerID string    `json:"manufacturerID"`
	Identification string    `json:"identification"`
	Readings       []DataSet `json:"datasets"`
}

// DataSet resource with json format tags.
type DataSet struct {
	Address string `json:"address"`
	Value   string `json:"value"`
	Unit    string `json:"unit"`
}

// CloudRepo implements a repo somewhere on the internet.
type CloudRepo struct {
	EndPoint *url.URL
	// TODO: Add credentials.
}

func (c *CloudRepo) Put(*model.Measurement) error {
	return nil
}
func (c *CloudRepo) Get(key []byte) (*model.Measurement, error) {
	return nil, nil
}
func (c *CloudRepo) GetN(n int) ([]*model.Measurement, error) {
	return nil, nil
}
func (c *CloudRepo) Delete(*model.Measurement) error {
	return nil
}
