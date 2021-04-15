package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"time"

	geojson "github.com/paulmach/go.geojson"
	gj "github.com/venicegeo/geojson-go/geojson"
)

func main() {
	log.SetFlags(log.LstdFlags | log.Lmicroseconds)
	if len(os.Args) < 2 {
		log.Fatal("no filename given")
	}
	f, err := os.Open(os.Args[1])
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()
	collection := new(geojson.FeatureCollection)
	dec := json.NewDecoder(f)
	for {
		log.Println("loading")
		now := time.Now()
		if err := dec.Decode(collection); err != nil {
			log.Fatal(err)
		}
		log.Printf("loaded (%s)", time.Since(now))
		for i, feature := range collection.Features {
			if feature.Geometry.IsPolygon() {
				bb, err := gj.NewBoundingBox(feature.Geometry.Polygon)
				if err != nil {
					log.Fatal(err)
				}
				c := bb.Centroid()
				pt := c.Coordinates
				fmt.Printf("%.6f,%.6f\n", pt[0], pt[1])
			} else {
				log.Printf("feature %d is a %s\n", i, feature.Type)
			}
		}
		if !dec.More() {
			break
		}
	}
}
