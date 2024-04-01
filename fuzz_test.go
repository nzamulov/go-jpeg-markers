package gojpegmarkers

import "testing"

func FuzzGetAllMarkers(f *testing.F) {
	f.Fuzz(func(t *testing.T, b []byte) {
		GetAllMarkers(b)
	})
}
