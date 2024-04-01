package gojpegmarkers

func isRST(id int) bool {
	for i := 0; i <= 7; i++ {
		if id == RST0+i {
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
		n, marker := Scan(img[offset:])
		marker.Offset = offset
		markers = append(markers, marker)
		if n == 0 { // EOI
			break
		}
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
