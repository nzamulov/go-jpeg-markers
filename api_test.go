package gojpegmarkers

import "testing"

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
		rawImg, err := ReadFile("./testdata/image.jpg")
		if err != nil {
			t.Fatal(err)
		}

		hasRSTm := HasRSTm(rawImg)
		if hasRSTm {
			t.Errorf("image should not has RSTm")
		}
	})
}

func TestGetImageXY(t *testing.T) {
	t.Run("Get X and Y of testdata/image.jpg", func(t *testing.T) {
		rawImg, err := ReadFile("./testdata/image.jpg")
		if err != nil {
			t.Fatal(err)
		}

		x, y := GetImageXY(rawImg)
		if x != 783 || y != 522 {
			t.Errorf("image should has dimensions x = 783px and y = 522px")
		}
	})

	t.Run("Get X and Y of testdata/withRSTm.jpg", func(t *testing.T) {
		rawImg, err := ReadFile("./testdata/withRSTm.jpg")
		if err != nil {
			t.Fatal(err)
		}

		x, y := GetImageXY(rawImg)
		if x != 50 || y != 50 {
			t.Errorf("image should has dimensions x = 50px and y = 50px")
		}
	})
}
