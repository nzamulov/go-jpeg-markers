package main

import (
	"fmt"
	_ "image/jpeg"
)

const (
	SOI  = 0xFFD8 // Start Of Image
	EOI  = 0xFFD9 // End Of Image
	DQT  = 0xFFDB // Define Quantization Table(s)
	DHT  = 0xFFC4 // Define Huffman Table(s)
	SOF0 = 0xFFC0 // Baseline DCT
	SOF2 = 0xFFC2 // Progressive DCT, Huffman coding
	COM  = 0xFFFE // Comment
	SOS  = 0xFFDA // Start Of Scan
	DRI  = 0xFFDD // Define Restart Interval
)

const (
	APP0 = 0xFFE0 + iota
	EXIF
	APP2
	APP3
	APP4
	APP5
	APP6
	APP7
	APP8
	APP9
	APP10
	APP11
	APP12
	APP13
	APP14
	APP15
)

const (
	RST0 = 0xFFD0 + iota
	RST1
	RST2
	RST3
	RST4
	RST5
	RST6
	RST7
)

type Marker struct {
	ID, Offset int
	Comment    string
}

func scan(b []byte) (int, Marker) {
	h := uint16(b[0])<<8 | uint16(b[1])

	switch h {
	case SOI:
		return 2, Marker{
			ID:      SOI,
			Comment: "0xFFD8: Start Of Image",
		}
	case APP0:
		return 2 + int(b[2])<<8 + int(b[3]), Marker{
			ID: APP0,
			// According to https://en.wikipedia.org/wiki/JPEG_File_Interchange_Format
			Comment: fmt.Sprintf("0xFFE0: JFIF ["+
				"Identifier:%s, "+
				"JFIF version:%d.%02d "+
				"Density units:%d, "+
				"Xdensity:%d, "+
				"Ydensity:%d, "+
				"Xthumbnail:%d, "+
				"Ythumbnail:%d]",
				string(b[4:9]),        // Identifier (5 bytes), 4A 46 49 46 00 = "JFIF" in ASCII, terminated by a null byte
				int(b[9]), int(b[10]), // First byte for major version, second byte for minor version (01 02 for 1.02)
				b[11],
				int(b[12])<<8+int(b[13]),
				int(b[14])<<8+int(b[15]),
				b[16],
				b[17],
				// TODO: thumbnail data???
			),
		}
	case APP2, APP3, APP4, APP5, APP6, APP7, APP8, APP9, APP10, APP11, APP12, APP13, APP14, APP15:
		return 2 + int(b[2])<<8 + int(b[3]), Marker{
			ID:      int(h),
			Comment: fmt.Sprintf("0x%X: APP%d", h, h-APP0),
		}
	case EXIF:
		return 2 + int(b[2])<<8 + int(b[3]), Marker{
			ID:      EXIF,
			Comment: "0xFFE1: EXIF",
		}
	case DQT:
		return 2 + int(b[2])<<8 + int(b[3]), Marker{
			ID:      DQT,
			Comment: "0xFFDB: Define Quantization Table(s)",
		}
	case DHT:
		return 2 + int(b[2])<<8 + int(b[3]), Marker{
			ID:      DHT,
			Comment: "0xFFC4: Define Huffman Table(s)",
		}
	case SOF0: // TODO: all SOFs???
		return 2 + int(b[2])<<8 + int(b[3]), Marker{
			ID:      SOF0,
			Comment: "0xFFC0: Start of Frame (Baseline DCT)",
		}
	case SOF2:
		return 2 + int(b[2])<<8 + int(b[3]), Marker{
			ID: SOF2,
			Comment: fmt.Sprintf("0xFFC2: Start Of Frame (Progressive DCT, Huffman coding) [P:%d, Y:%d, X:%d, Nf:%d]",
				b[4],
				int(b[5])<<8+int(b[6]),
				int(b[7])<<8+int(b[8]),
				b[9],
			),
		}
	case COM:
		return 2 + int(b[2])<<8 + int(b[3]), Marker{
			ID:      COM,
			Comment: "0xFFF3: Comment",
		}
	case SOS:
		header := 2 + int(b[2])<<8 + int(b[3])
		offset := header
		for {
			if b[offset] == 0xFF && b[offset+1] != 0x00 { // run to any other marker
				break
			}
			offset++
		}
		return offset, Marker{
			ID:      SOS,
			Comment: fmt.Sprintf("0xFFDA: Start Of Scan [Ns: %d] (%d bytes)", b[4], offset-header),
		}
	case DRI:
		return 2 + int(b[2])<<8 + int(b[3]), Marker{
			ID:      DRI,
			Comment: fmt.Sprintf("0xFFDD: Define Restart Interval [Ri: %d]", int(b[4])<<8+int(b[5])),
		}
	case RST0, RST1, RST2, RST3, RST4, RST5, RST6, RST7:
		header := 2
		offset := header
		for {
			if b[offset] == 0xFF && b[offset+1] != 0x00 { // run to any other marker
				break
			}
			offset++
		}
		return offset, Marker{
			ID:      int(h),
			Comment: fmt.Sprintf("0x%X: RST%d (%d bytes)", h, h-RST0, offset-header),
		}
	case EOI:
		return 0, Marker{
			ID:      EOI,
			Comment: "0xFFD9: End Of Image",
		}
	default:
		panic(fmt.Sprintf("unexpected marker: %x", h))
	}
	return 0, Marker{}
}

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
