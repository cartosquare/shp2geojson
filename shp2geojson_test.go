package shp2geojson

import (
	"io/ioutil"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestPoint(t *testing.T) {
	data := "./data/point.shp"
	jsonStr, err := Convert(data)
	assert.Equal(t, err, nil)

	err = ioutil.WriteFile("./fixture/point.geojson", []byte(jsonStr), 0644)
	assert.Equal(t, err, nil)
}

func TestLine(t *testing.T) {
	data := "./data/line.shp"
	jsonStr, err := Convert(data)
	assert.Equal(t, err, nil)

	err = ioutil.WriteFile("./fixture/line.geojson", []byte(jsonStr), 0644)
	assert.Equal(t, err, nil)
}

func TestPolygon(t *testing.T) {
	data := "./data/polygon.shp"
	jsonStr, err := Convert(data)
	assert.Equal(t, err, nil)

	err = ioutil.WriteFile("./fixture/polygon.geojson", []byte(jsonStr), 0644)
	assert.Equal(t, err, nil)
}
