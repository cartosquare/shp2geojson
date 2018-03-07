package shp2geojson

import (
	"fmt"
	"reflect"

	"github.com/jonas-p/go-shp"
	"github.com/paulmach/go.geojson"
)

// Convert shapefile to geojson
// currently not support Multi Geometry
func Convert(input string) (string, error) {
	shape, err := shp.Open(input)
	if err != nil {
		return "", err
	}
	defer shape.Close()

	// fields from the attribute table (DBF)
	fields := shape.Fields()
	fc := geojson.NewFeatureCollection()
	// loop through all features in the shapefile
	for shape.Next() {
		n, s := shape.Shape()
		var feature *geojson.Feature

		switch s.(type) {
		case *shp.Point:
			pt, _ := s.(*shp.Point)
			feature = geojson.NewPointFeature([]float64{pt.X, pt.Y})
		case *shp.PolyLine:
			line, _ := s.(*shp.PolyLine)
			if line.NumParts != 1 {
				fmt.Println("Warning: more than 1 part in polyline!")
			}
			coordinates := make([][]float64, line.NumPoints)
			for i := range coordinates {
				coordinates[i] = make([]float64, 2)
			}

			for idx, pt := range line.Points {
				coordinates[idx][0] = pt.X
				coordinates[idx][1] = pt.Y
			}
			feature = geojson.NewLineStringFeature(coordinates)
		case *shp.Polygon:
			polygon, _ := s.(*shp.Polygon)
			coordinates := make([][][]float64, polygon.NumParts)
			var i int32
			for i = 0; i < polygon.NumParts; i = i + 1 {
				var startIndex, endIndex int32
				startIndex = polygon.Parts[i]
				if i == polygon.NumParts-1 {
					endIndex = int32(len(polygon.Points))
				} else {
					endIndex = polygon.Parts[i+1]
				}

				coordinates[i] = make([][]float64, endIndex-startIndex)
				cnt := 0
				for j := startIndex; j < endIndex; j = j + 1 {
					coordinates[i][cnt] = make([]float64, 2)
					coordinates[i][cnt][0] = polygon.Points[j].X
					coordinates[i][cnt][1] = polygon.Points[j].Y
					cnt = cnt + 1
				}
			}
			feature = geojson.NewPolygonFeature(coordinates)
		default:
			fmt.Println("Not support geometry type", reflect.TypeOf(s).Elem())
			continue
		}

		for k, f := range fields {
			val := shape.ReadAttribute(n, k)
			feature.Properties[f.String()] = val
		}
		fc.AddFeature(feature)
	}

	rawJSON, err := fc.MarshalJSON()
	if err != nil {
		return "", err
	}
	return string(rawJSON), nil
}
