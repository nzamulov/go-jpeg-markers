package gojpegmetrics

func isRST(id int) bool {
	for i := 0; i <= 7; i++ {
		if id == RST0+i {
			return true
		}
	}
	return false
}

func Scan(img []byte) (markers []Marker) {
	offset := 0
	for {
		n, marker := scan(img[offset:])
		marker.Offset = offset
		markers = append(markers, marker)
		if n == 0 { // EOI
			break
		}
		offset += n
	}
	return
}

func CheckRSTm(img []byte) bool {
	markers := Scan(img)
	for _, marker := range markers {
		if isRST(marker.ID) {
			return true
		}
	}
	return false
}
