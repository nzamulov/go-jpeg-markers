package gojpegmarkers

func isRST(id int) bool {
	for i := 0; i <= 7; i++ {
		if id == RST0+i {
			return true
		}
	}
	return false
}

func isSOF(id int) bool {
	for i := 0; i <= 15; i++ {
		if id == SOF0+i {
			return true
		}
	}
	return false
}

func GetAllMarkers(img []byte) (markers []Marker) {
	offset := 0
	for {
		// broken data was passed and offset goes far away than bound
		if offset >= len(img) {
			break
		}
		n, marker := scan(img[offset:])
		if n == 0 { // EOI
			break
		}
		marker.Offset = offset
		markers = append(markers, marker)
		offset += n
	}
	return
}

func HasRSTm(img []byte) bool {
	markers := GetAllMarkers(img)
	for _, marker := range markers {
		if isRST(marker.ID) {
			return true
		}
	}
	return false
}

func GetWidthHeight(img []byte) (int, int) {
	markers := GetAllMarkers(img)
	for _, marker := range markers {
		if isSOF(marker.ID) && len(marker.AdditionalInfo) > 0 {
			return marker.AdditionalInfo[WIDTH].(int), marker.AdditionalInfo[HEIGHT].(int)
		}
	}
	return 0, 0
}
