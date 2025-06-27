package gojpegmarkers

import (
	"encoding/binary"
	"fmt"
	_ "image/jpeg"
	"math"
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
	JPG = 0xFFC8 // Reserved for JPEG extensions
	DAC = 0xFFCC // Define arithmetic coding conditioning(s)
	DNL = 0xFFDC // Define number of lines
	DHP = 0xFFDE // Define hierarchical progression
	EXP = 0xFFDF // Expand reference component(s)
)

// Start of Frame header specifies the source image characteristics: the components of the frame, and the
// sampling factors for each component, and specifies the destinations from which the quantized tables
// to be used with each component are retrieved.
//
// Start of Frame: non-differential, Huffman coding.
const (
	SOF0 = 0xFFC0 + iota // Baseline DCT
	SOF1                 // Extended sequential DCT
	SOF2                 // Progressive DCT, Huffman coding
	SOF3                 // Lossless (sequential)
)

// Start of Frame: differential, Huffman coding.
const (
	SOF4 = DHT + iota // SOF4 = DHT
	SOF5              // Differential sequential DCT
	SOF6              // Differential progressive DCT
	SOF7              // Differential lossless (sequential)
)

// Start of Frame: non-differential, arithmetic coding.
const (
	SOF8  = JPG + iota // SOF8 = JPG
	SOF9               // Extended sequential DCT
	SOF10              // Progressive DCT
	SOF11              // Lossless (sequential)
)

// Start of Frame: differential, arithmetic coding.
const (
	SOF12 = DAC + iota // SOF12 = DAC
	SOF13              // Differential sequential DCT
	SOF14              // Differential progressive DCT
	SOF15              // Differential lossless (sequential)
)

// SOFComments goes from CCITT Rec. T.81 (1992 E)
var SOFComments = map[int]string{
	SOF0: "Baseline DCT",
	SOF1: "Extended sequential DCT",
	SOF2: "Progressive DCT, Huffman coding",
	SOF3: "Lossless (sequential)",
	// SOF4 = DHT
	SOF5: "Differential sequential DCT",
	SOF6: "Differential progressive DCT",
	SOF7: "Differential lossless (sequential)",
	// SOF8 = JPG
	SOF9:  "Extended sequential DCT",
	SOF10: "Progressive DCT",
	SOF11: "Lossless (sequential)",
	// SOF12 = DAC
	SOF13: "Differential sequential DCT",
	SOF14: "Differential progressive DCT",
	SOF15: "Differential lossless (sequential)",
}

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

func getOffsetMaybeWithLen(b []byte, skipLen bool) int {
	// if 'b' passed as broken data (less than marker length itself), just return the length of b
	if len(b) < 2 {
		return len(b)
	}
	if skipLen {
		return 2
	}
	// if 'b' passed as broken data (with marker but with broken length), just return the length of b
	if len(b) < 4 {
		return len(b)
	}
	//fmt.Println(2 + int(b[2])<<8 + int(b[3]))
	//fmt.Printf("%s\n", b[4:16])
	return 2 + int(b[2])<<8 + int(b[3])
}

func scan(b []byte) (int, Marker) {
	if len(b) <= 1 {
		return len(b), Marker{Comment: "broken marker"}
	}

	h := uint16(b[0])<<8 | uint16(b[1])

	switch h {
	case SOI:
		return getOffsetMaybeWithLen(b, true), Marker{
			ID:      SOI,
			Comment: "0xFFD8: Start Of Image",
		}
	case APP0:
		if len(b) < 18 {
			return len(b), Marker{Comment: "broken marker"}
		}
		return getOffsetMaybeWithLen(b, false), Marker{
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
		comment := fmt.Sprintf("0x%X: APP%d", h, h-APP0)

		if h == APP2 {
			l := int(b[2])<<8 + int(b[3])
			fmt.Printf("LEN: %d\n", l)
			comment += fmt.Sprintf(" [ALGO: %s]", string(b[4:16]))
			fmt.Printf("%s\n", string(b[16:16+l]))

			// #region: Header (128 bytes)
			{
				// Profile size
				fmt.Println(int(b[16]), int(b[17]))
				fmt.Println(int(b[18]), int(b[19]))
				bits32 := binary.LittleEndian.Uint32(b[16:20])
				fmt.Println("Profile size = ", bits32)
				// CMM type
				fmt.Println(string(b[20:25]))
				// ...
				// Profile version
				fmt.Println(int(b[25])<<8 + int(b[26]))
				// b[27] - must be 0
				// b[28] - must be 0
				// Device class
				fmt.Println(string(b[29:34]))
				// Canonical input space
				fmt.Println(string(b[34:38]))
				// Canonical output space
				fmt.Println(string(b[38:42]))

				// Date
				year := int(b[42])<<8 + int(b[43])
				fmt.Println(year)
				month := int(b[44])<<8 + int(b[45])
				fmt.Println(month)
				day := int(b[46])<<8 + int(b[47])
				fmt.Println(day)
				hour := int(b[48])<<8 + int(b[49])
				fmt.Println(hour)
				minute := int(b[50])<<8 + int(b[51])
				fmt.Println(minute)
				second := int(b[52])<<8 + int(b[53])
				fmt.Println(second)

				// Profile file signature
				fmt.Println(string(b[54:58]))

				//  Primary platform target for the profile
				fmt.Println(string(b[58:62]))

				// Flags
				fmt.Println(string(b[62:66]))

				// Device manufacturer of the device for which this profile is created
				fmt.Println(string(b[66:70]))

				// Device model of the device for which this profile is created
				fmt.Println(string(b[70:74]))

				// Device attributes unique to the particular device setup such as media type
				fmt.Println(string(b[74:82]))

				// Rendering Intent
				fmt.Println(string(b[82:86]))

				// The XYZ values of the illuminant of the profile connection space. This must correspond to D50. It is explained in more detail in Annex A.1 'Profile Connection Spaces'.
				bits32 = binary.BigEndian.Uint32(b[86:90])
				fmt.Println(math.Float32frombits(bits32))

				bits32 = binary.BigEndian.Uint32(b[90:94])
				fmt.Println(math.Float32frombits(bits32))

				bits32 = binary.BigEndian.Uint32(b[94:98])
				fmt.Println(math.Float32frombits(bits32))

				// Identifies the creator of the profile
				fmt.Println(string(b[98:102]))

				// 44 bytes reserved for future expansion
				fmt.Println(b[102:146])
			}

			// Tags count
			bits32 := binary.BigEndian.Uint32(b[146:150])
			fmt.Println(bits32)

			fmt.Println(string(b[150:154]))
			fmt.Println(b[154:158], string(b[154:158]), binary.BigEndian.Uint32(b[154:158]))
			fmt.Println(b[158:162], string(b[158:162]), binary.BigEndian.Uint32(b[158:162]))
			fmt.Println(string(b[16+252 : 16+353]))

			fmt.Println(string(b[162:166]))
			fmt.Println(b[166:170], string(b[166:170]), binary.BigEndian.Uint32(b[166:170]))
			fmt.Println(b[174:178], string(b[174:178]), binary.BigEndian.Uint32(b[174:178]))

			fmt.Println(string(b[148 : 18+l-148]))
		}

		return getOffsetMaybeWithLen(b, false), Marker{
			ID:      int(h),
			Comment: comment,
		}
	case EXIF:
		return getOffsetMaybeWithLen(b, false), Marker{
			ID:      EXIF,
			Comment: "0xFFE1: EXIF",
		}
	case DQT:
		return getOffsetMaybeWithLen(b, false), Marker{
			ID:      DQT,
			Comment: "0xFFDB: Define Quantization Table(s)",
		}
	case DHT:
		return getOffsetMaybeWithLen(b, false), Marker{
			ID:      DHT,
			Comment: "0xFFC4: Define Huffman Table(s)",
		}
	case DNL:
		if len(b) < 6 {
			return len(b), Marker{Comment: "broken marker"}
		}
		return getOffsetMaybeWithLen(b, false), Marker{
			ID:      DNL,
			Comment: fmt.Sprintf("0xFFDC: Define number of lines [NL: %d]", int(b[4])<<8+int(b[5])),
		}
	case DHP:
		if len(b) < 10 {
			return len(b), Marker{Comment: "broken marker"}
		}
		return getOffsetMaybeWithLen(b, false), Marker{
			ID: DHP,
			Comment: fmt.Sprintf("0xFFDE: Define hierarchical progression [P:%d, Y:%d, X:%d, Nf:%d]",
				b[4],
				int(b[5])<<8+int(b[6]),
				int(b[7])<<8+int(b[8]),
				b[9],
			),
		}
	case EXP:
		if len(b) < 6 {
			return len(b), Marker{Comment: "broken marker"}
		}
		return getOffsetMaybeWithLen(b, false), Marker{
			ID:      EXP,
			Comment: fmt.Sprintf("0xFFDF: Expand reference component(s) [Eh:%d, Ev:%d]", b[4], b[5]),
		}
	case JPG:
		return getOffsetMaybeWithLen(b, false), Marker{
			ID:      JPG,
			Comment: "0xFFC8: Reserved for JPEG extensions",
		}
	case DAC:
		if len(b) < 7 {
			return len(b), Marker{Comment: "broken marker"}
		}
		return getOffsetMaybeWithLen(b, false), Marker{
			ID:      DAC,
			Comment: fmt.Sprintf("0xFFCC: Define arithmetic coding conditioning(s) [Tc:%d, Tb:%d, Cs:%d]", b[4], b[5], b[6]),
		}
	case SOF0, SOF1, SOF2, SOF3, SOF5, SOF6, SOF7, SOF9, SOF10, SOF11, SOF13, SOF14, SOF15:
		if len(b) < 10 {
			return len(b), Marker{Comment: "broken marker"}
		}
		return getOffsetMaybeWithLen(b, false), Marker{
			ID: SOF2,
			Comment: fmt.Sprintf("0x%X: Start Of Frame (SOF%d) (%s) [P:%d, Y:%d, X:%d, Nf:%d]",
				h,
				h-SOF0,
				SOFComments[int(h)],
				b[4],
				int(b[5])<<8+int(b[6]),
				int(b[7])<<8+int(b[8]),
				b[9],
			),
		}
	case COM:
		return getOffsetMaybeWithLen(b, false), Marker{
			ID:      COM,
			Comment: "0xFFF3: Comment",
		}
	case SOS:
		if len(b) < 5 {
			return len(b), Marker{Comment: "broken marker"}
		}
		header := getOffsetMaybeWithLen(b, false)
		offset := header
		for {
			if offset+1 >= len(b) || (b[offset] == 0xFF && b[offset+1] != 0x00) { // run to any other marker
				break
			}
			offset++
		}
		return offset, Marker{
			ID:      SOS,
			Comment: fmt.Sprintf("0xFFDA: Start Of Scan [Ns: %d] (%d bytes)", b[4], offset-header),
		}
	case DRI:
		if len(b) < 6 {
			return len(b), Marker{Comment: "broken marker"}
		}
		return getOffsetMaybeWithLen(b, false), Marker{
			ID:      DRI,
			Comment: fmt.Sprintf("0xFFDD: Define Restart Interval [Ri: %d]", int(b[4])<<8+int(b[5])),
		}
	case RST0, RST1, RST2, RST3, RST4, RST5, RST6, RST7:
		header := getOffsetMaybeWithLen(b, true)
		offset := header
		for {
			if offset+1 >= len(b) || (b[offset] == 0xFF && b[offset+1] != 0x00) { // run to any other marker
				break
			}
			offset++
		}
		return offset, Marker{
			ID:      int(h),
			Comment: fmt.Sprintf("0x%X: RST%d (%d bytes)", h, h-RST0, offset-header),
		}
	case EOI:
		return getOffsetMaybeWithLen(b, true), Marker{
			ID:      EOI,
			Comment: "0xFFD9: End Of Image",
		}
	default:
		header := getOffsetMaybeWithLen(b, true)
		offset := header
		for {
			if offset+1 >= len(b) || (b[offset] == 0xFF && b[offset+1] != 0x00) { // run to any other marker
				break
			}
			offset++
		}
		return offset, Marker{
			ID:      int(h),
			Comment: "unexpected marker",
		}
	}
}
