package gojpegmarkers

import "testing"

func TestGetAllMarkers(t *testing.T) {
	t.Run("Get all markers of ./testdata/withoutRSTm.jpg", func(t *testing.T) {
		rawImg, err := ReadFile("./testdata/withoutRSTm.jpg")
		if err != nil {
			t.Fatal(err)
		}

		markers := GetAllMarkers(rawImg)

		if len(markers) != 30 {
			t.Errorf("image has 30 markers, not %d\n", len(markers))
		}
	})

	t.Run("Get all markers of ./testdata/withRSTm.jpg", func(t *testing.T) {
		rawImg, err := ReadFile("./testdata/withRSTm.jpg")
		if err != nil {
			t.Fatal(err)
		}

		markers := GetAllMarkers(rawImg)

		if len(markers) != 26 {
			t.Errorf("image has 26 markers, not %d\n", len(markers))
		}
	})
}

func TestHasRSTm(t *testing.T) {
	t.Run("RSTm flags does exist", func(t *testing.T) {
		rawImg, err := ReadFile("./testdata/withRSTm.jpg")
		if err != nil {
			t.Fatal(err)
		}

		hasRSTm := HasRSTm(rawImg)
		if !hasRSTm {
			t.Errorf("image should has RSTm")
		}
	})

	t.Run("RSTm flags does not exist", func(t *testing.T) {
		rawImg, err := ReadFile("./testdata/withoutRSTm.jpg")
		if err != nil {
			t.Fatal(err)
		}

		hasRSTm := HasRSTm(rawImg)
		if hasRSTm {
			t.Errorf("image should not has RSTm")
		}
	})
}

func TestGetWidthHeight(t *testing.T) {
	t.Run("Get width and height of ./testdata/withoutRSTm.jpg", func(t *testing.T) {
		rawImg, err := ReadFile("./testdata/withoutRSTm.jpg")
		if err != nil {
			t.Fatal(err)
		}

		width, height := GetWidthHeight(rawImg)
		if width != 783 || height != 522 {
			t.Errorf("image should has dimensions width = 783px and height = 522px")
		}
	})

	t.Run("Get width and height of ./testdata/withRSTm.jpg", func(t *testing.T) {
		rawImg, err := ReadFile("./testdata/withRSTm.jpg")
		if err != nil {
			t.Fatal(err)
		}

		width, height := GetWidthHeight(rawImg)
		if width != 50 || height != 50 {
			t.Errorf("image should has dimensions width = 50px and height = 50px")
		}
	})
}
