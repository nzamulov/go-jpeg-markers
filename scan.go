package gojpegmetrics

import (
	"fmt"
	_ "image/jpeg"
)

// A JPEG image consists of a sequence of segments, each beginning with a marker,
// each of which begins with a 0xFF byte, followed by a byte indicating what kind
// of marker it is. Some markers consist of just those two bytes; others are
// followed by two bytes (high then low), indicating the length of marker-specific
// payload data that follows. (The length includes the two bytes for the length,
// but not the two bytes for the marker.) Some markers are followed by entropy-coded data;
// the length of such a marker does not include the entropy-coded data. Note that consecutive
// 0xFF bytes are used as fill bytes for padding purposes, although this fill byte padding
// should only ever take place for markers immediately following entropy-coded scan data.

// Common JPEG markers.
const (
	SOI = 0xFFD8 // Start Of Image
	EOI = 0xFFD9 // End Of Image
	DQT = 0xFFDB // Define Quantization Table(s)
	DHT = 0xFFC4 // Define Huffman Table(s)
	COM = 0xFFFE // Comment
	SOS = 0xFFDA // Start Of Scan
	DRI = 0xFFDD // Define Restart Interval
)

// TODO: frame header scheme???
// Start of Frame header specifies the source image characteristics: the components of the frame, and the
// sampling factors for each component, and specifies the destinations from which the quantized tables
// to be used with each component are retrieved.
//
// Start of Frame: non-differential, Huffman coding
const (
	SOF0 = 0xFFC0 + iota // Baseline DCT
	SOF1                 // Extended sequential DCT
	SOF2                 // Progressive DCT, Huffman coding
	SOF3                 // Lossless (sequential)
)

// Start of Frame: differential, Huffman coding
const (
	SOF4 = DHT + iota // SOF4 = DHT
	SOF5              // Differential sequential DCT
	SOF6              // Differential progressive DCT
	SOF7              // Differential lossless (sequential)
)

// Start of Frame: non-differential, arithmetic coding
const (
	SOF8  = 0xFFC8 + iota // Reserved for JPEG extensions
	SOF9                  // Extended sequential DCT
	SOF10                 // Progressive DCT
	SOF11                 // Lossless (sequential)
)

// Start of Frame: differential, arithmetic coding
const (
	SOF12 = 0xFFCC + iota // TODO: ???
	SOF13                 // Differential sequential DCT
	SOF14                 // Differential progressive DCT
	SOF15                 // Differential lossless (sequential)
)

// Application-specific markers.
//
// For example, an Exif JPEG file uses an APP1 (EXIF) marker to store metadata,
// laid out in a structure based closely on TIFF.
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

// Restart markers.
//
// Inserted every r macroblocks, where r is the restart interval set by a DRI marker.
// Not used if there was no DRI marker. The low three bits of the marker code cycle
// in value from 0 to 7.
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
				"Ythumbnail:%d"+
				"]",
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
	case SOF0:
		return 2 + int(b[2])<<8 + int(b[3]), Marker{
			ID:      SOF0,
			Comment: "0xFFC0: Start of Frame (Baseline DCT)",
		}
	case SOF1, SOF3, SOF5, SOF6, SOF7, SOF8, SOF9, SOF10, SOF11, SOF12, SOF13, SOF14, SOF15:
		return 2 + int(b[2])<<8 + int(b[3]), Marker{
			ID:      int(h),
			Comment: fmt.Sprintf("0x%X: SOF%d", h, h-SOF0),
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

// API

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
