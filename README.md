# go-jpeg-markers [![CI][ci-img]][ci]

Simply native Go tool for fetch all markers from JPEG image.
Only Go, no more...

## What are "markers" in JPEG?

According to [CCITT T.81 THE INTERNATIONAL (09/92)](https://www.w3.org/Graphics/JPEG/itu-t81.pdf) markers serve to
identify the various structural parts of the compressed data formats. Most markers start marker segments containing a
related group of parameters, some markers stand alone. All markers are assigned two-byte codes: an 0xFF byte followed 
by a byte which os not equal to 0x00 of 0xFF.

## Motivation

Working a lot of times with images, in most cases with JPEG, there is no instrument was found through the Internet that
will allow you to find or to list all JPEG markers spending a minimum time. This instrument will allow you to inspect or
to check needed JPEG markers in extremely short time. For example, very helpful for work with GPU CUDA-workers and for 
checking restart markers (RSTm) into the JPEG file. Also, helpful for checking any markers in source code at runtime.

Inspired by:
https://github.com/MavEtJu/jpeg-markers

## How to use as a lib?

```bash
go get -u github.com/nzamulov/go-jpeg-markers
```

```go
package main

import (
	"bufio"
	"fmt"
	"io"
	"os"

	"github.com/nzamulov/go-jpeg-markers"
)

func main() {
	file, err := os.Open("./image.jpg")
	if err != nil {
		panic(err)
	}
	defer file.Close()

	img, err := io.ReadAll(bufio.NewReader(file))
	if err != nil {
		panic(err)
	}

	var (
		markers     = gojpegmetrics.Scan(img)
		isRSTmExist = gojpegmetrics.CheckRSTm(img)
	)

	fmt.Printf("markers len: %d, RSTm: %t\n", len(markers), isRSTmExist)
}
```

## How to use as a CLI?

```bash
make scan <path_to_file_in_fs | link_from_internet>
```

File from FS:
```bash
make scan testdata/image.jpg
```
```bash
2024/03/29 12:41:15.652918 [go-jpeg-markers] offset:      0 - 0xFFD8: Start Of Image
2024/03/29 12:41:15.653262 [go-jpeg-markers] offset:      2 - 0xFFE0: JFIF [Identifier:JFIF, JFIF version:1.01 Density units:1, Xdensity:72, Ydensity:72, Xthumbnail:0, Ythumbnail:0]
2024/03/29 12:41:15.653266 [go-jpeg-markers] offset:     14 - 0xFFED: APP13
2024/03/29 12:41:15.653268 [go-jpeg-markers] offset:   223c - 0xFFE1: EXIF
2024/03/29 12:41:15.653270 [go-jpeg-markers] offset:   2262 - 0xFFE2: APP2
2024/03/29 12:41:15.653272 [go-jpeg-markers] offset:   2ebc - 0xFFE1: EXIF
2024/03/29 12:41:15.653274 [go-jpeg-markers] offset:   448f - 0xFFDB: Define Quantization Table(s)
2024/03/29 12:41:15.653277 [go-jpeg-markers] offset:   44d4 - 0xFFDB: Define Quantization Table(s)
2024/03/29 12:41:15.653279 [go-jpeg-markers] offset:   4519 - 0xFFC2: Start Of Frame (SOF2) (Progressive DCT, Huffman coding) [P:8, Y:522, X:783, Nf:3]
2024/03/29 12:41:15.653281 [go-jpeg-markers] offset:   452c - 0xFFC4: Define Huffman Table(s)
2024/03/29 12:41:15.653283 [go-jpeg-markers] offset:   454b - 0xFFC4: Define Huffman Table(s)
2024/03/29 12:41:15.653284 [go-jpeg-markers] offset:   4565 - 0xFFDA: Start Of Scan [Ns: 3] (8184 bytes)
2024/03/29 12:41:15.653286 [go-jpeg-markers] offset:   656b - 0xFFC4: Define Huffman Table(s)
2024/03/29 12:41:15.653288 [go-jpeg-markers] offset:   659f - 0xFFDA: Start Of Scan [Ns: 1] (21274 bytes)
2024/03/29 12:41:15.653290 [go-jpeg-markers] offset:   b8c3 - 0xFFC4: Define Huffman Table(s)
2024/03/29 12:41:15.653291 [go-jpeg-markers] offset:   b8f8 - 0xFFDA: Start Of Scan [Ns: 1] (1333 bytes)
2024/03/29 12:41:15.653293 [go-jpeg-markers] offset:   be37 - 0xFFC4: Define Huffman Table(s)
2024/03/29 12:41:15.653304 [go-jpeg-markers] offset:   be6b - 0xFFDA: Start Of Scan [Ns: 1] (1282 bytes)
2024/03/29 12:41:15.653306 [go-jpeg-markers] offset:   c377 - 0xFFC4: Define Huffman Table(s)
2024/03/29 12:41:15.653308 [go-jpeg-markers] offset:   c3cc - 0xFFDA: Start Of Scan [Ns: 1] (147491 bytes)
2024/03/29 12:41:15.653310 [go-jpeg-markers] offset:  303f9 - 0xFFC4: Define Huffman Table(s)
2024/03/29 12:41:15.653311 [go-jpeg-markers] offset:  30422 - 0xFFDA: Start Of Scan [Ns: 1] (44414 bytes)
2024/03/29 12:41:15.653313 [go-jpeg-markers] offset:  3b1aa - 0xFFDA: Start Of Scan [Ns: 3] (2426 bytes)
2024/03/29 12:41:15.653316 [go-jpeg-markers] offset:  3bb32 - 0xFFC4: Define Huffman Table(s)
2024/03/29 12:41:15.653317 [go-jpeg-markers] offset:  3bb57 - 0xFFDA: Start Of Scan [Ns: 1] (575 bytes)
2024/03/29 12:41:15.653319 [go-jpeg-markers] offset:  3bda0 - 0xFFC4: Define Huffman Table(s)
2024/03/29 12:41:15.653321 [go-jpeg-markers] offset:  3bdc5 - 0xFFDA: Start Of Scan [Ns: 1] (577 bytes)
2024/03/29 12:41:15.653324 [go-jpeg-markers] offset:  3c010 - 0xFFC4: Define Huffman Table(s)
2024/03/29 12:41:15.653326 [go-jpeg-markers] offset:  3c039 - 0xFFDA: Start Of Scan [Ns: 1] (48374 bytes)
2024/03/29 12:41:15.653328 [go-jpeg-markers] offset:  47d39 - 0xFFD9: End Of Image
```

File from link:
```bash
make scan https://sun9-80.userapi.com/V936_Yc4ULoWeaJ0_KRy_QfObYjex7mhPo1Qgg/i6DWh5it8lk.jpg
```
```bash
2024/03/29 13:16:44.789233 [go-jpeg-markers] offset:      0 - 0xFFD8: Start Of Image
2024/03/29 13:16:44.789690 [go-jpeg-markers] offset:      2 - 0xFFDB: Define Quantization Table(s)
2024/03/29 13:16:44.789693 [go-jpeg-markers] offset:     47 - 0xFFDB: Define Quantization Table(s)
2024/03/29 13:16:44.789695 [go-jpeg-markers] offset:     8c - 0xFFC0: Start Of Frame (SOF0) (Baseline DCT) [P:8, Y:50, X:50, Nf:3]
2024/03/29 13:16:44.789697 [go-jpeg-markers] offset:     9f - 0xFFC4: Define Huffman Table(s)
2024/03/29 13:16:44.789699 [go-jpeg-markers] offset:     bc - 0xFFC4: Define Huffman Table(s)
2024/03/29 13:16:44.789700 [go-jpeg-markers] offset:     ec - 0xFFC4: Define Huffman Table(s)
2024/03/29 13:16:44.789702 [go-jpeg-markers] offset:    103 - 0xFFC4: Define Huffman Table(s)
2024/03/29 13:16:44.789704 [go-jpeg-markers] offset:    124 - 0xFFDD: Define Restart Interval [Ri: 1]
2024/03/29 13:16:44.789706 [go-jpeg-markers] offset:    12a - 0xFFDA: Start Of Scan [Ns: 3] (6 bytes)
2024/03/29 13:16:44.789708 [go-jpeg-markers] offset:    13e - 0xFFD0: RST0 (240 bytes)
2024/03/29 13:16:44.789740 [go-jpeg-markers] offset:    230 - 0xFFD1: RST1 (79 bytes)
2024/03/29 13:16:44.789747 [go-jpeg-markers] offset:    281 - 0xFFD2: RST2 (6 bytes)
2024/03/29 13:16:44.789749 [go-jpeg-markers] offset:    289 - 0xFFD3: RST3 (113 bytes)
2024/03/29 13:16:44.789752 [go-jpeg-markers] offset:    2fc - 0xFFD4: RST4 (256 bytes)
2024/03/29 13:16:44.789753 [go-jpeg-markers] offset:    3fe - 0xFFD5: RST5 (108 bytes)
2024/03/29 13:16:44.789755 [go-jpeg-markers] offset:    46c - 0xFFD6: RST6 (6 bytes)
2024/03/29 13:16:44.789757 [go-jpeg-markers] offset:    474 - 0xFFD7: RST7 (67 bytes)
2024/03/29 13:16:44.789758 [go-jpeg-markers] offset:    4b9 - 0xFFD0: RST0 (177 bytes)
2024/03/29 13:16:44.789760 [go-jpeg-markers] offset:    56c - 0xFFD1: RST1 (99 bytes)
2024/03/29 13:16:44.789762 [go-jpeg-markers] offset:    5d1 - 0xFFD2: RST2 (6 bytes)
2024/03/29 13:16:44.789763 [go-jpeg-markers] offset:    5d9 - 0xFFD3: RST3 (6 bytes)
2024/03/29 13:16:44.789765 [go-jpeg-markers] offset:    5e1 - 0xFFD4: RST4 (36 bytes)
2024/03/29 13:16:44.789768 [go-jpeg-markers] offset:    607 - 0xFFD5: RST5 (11 bytes)
2024/03/29 13:16:44.789770 [go-jpeg-markers] offset:    614 - 0xFFD6: RST6 (6 bytes)
2024/03/29 13:16:44.789771 [go-jpeg-markers] offset:    61c - 0xFFD9: End Of Image
```

## License
[Mozilla Public License Version 2.0](./LICENSE)

[ci-img]: https://github.com/nzamulov/go-jpeg-markers/workflows/CI/badge.svg
[ci]: https://github.com/nzamulov/go-jpeg-markers/actions