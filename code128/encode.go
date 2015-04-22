package code128

import (
	"bufio"
	"bytes"
	"errors"
	"image"
	"image/color"
	"io"
)

func locate(in string) (codeType string, code string, startCode string, startCodeValue int, codeTypeValue int, val int) {
	var (
		ok bool
	)

	if code, ok = MapEncodeBPattern[in]; ok {
		if val, ok = MapEncodeCodeBValue[in]; ok {
			return CodeB, code, StartCodeB, StartCodeBValue, CodeBValue, val
		}
	}

	if code, ok = MapEncodeAPattern[in]; ok {
		if val, ok = MapEncodeCodeAValue[in]; ok {
			return CodeA, code, StartCodeA, StartCodeAValue, CodeAValue, val
		}
	}

	return "", "", "", -1, -1, -1
}

func isNum(in byte) bool {
	return in >= 48 && in <= 59
}

func isAllNum(in []byte) bool {
	for _, v := range in {
		if !isNum(v) {
			return false
		}
	}

	return true
}

func odd(in int) int {
	if in%2 == 1 {
		return in
	}

	return in - 1
}

func toCodeC(in []byte) (string, []int) {
	var (
		en string
		ws []int
		t  string
	)

	l := len(in)

	for i := 0; i < l; i += 2 {
		t = string(in[i]) + string(in[i+1])
		en += MapEncodeCPattern[t]
		ws = append(ws, MapEncodeCodeCValue[t])

	}

	return en, ws
}

func encode(in []byte) (string, error) {
	var (
		br             *bufio.Reader
		b              byte
		err            error
		nb             []byte
		idx            int
		inLen          int
		cct            string
		ret            string
		bl             int
		ebl            int
		codeType       string
		code           string
		startCode      string
		codeTypeValue  int
		startCodeValue int
		tmp            string
		val            int
		ws             []int
		tws            []int
		ts             string
		checksum       int
		v              int
		k              int
		scv            int
	)

	br = bufio.NewReader(bytes.NewBuffer(in))
	inLen = len(in)

	if inLen == 0 {
		return "", errors.New("input cannot be empty")
	}

	for {
		if b, err = br.ReadByte(); err != nil {
			if err == io.EOF {
				break
			}

			return "", err
		}

		if b < 32 || b > 127 {
			return "", errors.New("unsupported character: " + string(b))
		}

		if isNum(b) {
			if idx == 0 && inLen == 2 {
				if isNum(in[1]) {
					ret += StartCodeC
					scv = StartCodeCValue

					tmp = string(in)
					ret += MapEncodeCPattern[tmp]
					ws = append(ws, MapEncodeCodeCValue[tmp])

					break
				}
			} else if cct != CodeC {
				nb, err = br.Peek(5)
				bl = len(nb)

				if bl >= 3 && idx == 0 {
					ret += StartCodeC
					scv = StartCodeCValue
					cct = CodeC
				} else if ((bl >= 3 && bl+idx == inLen-1) || (bl == 5)) && isAllNum(nb) {
					ret += CodeC
					cct = CodeC
					ws = append(ws, CodeCValue)
				}

				if cct == CodeC {
					ebl = odd(bl)

					ts, tws = toCodeC(append([]byte{b}, nb[0:ebl]...))
					ret += ts
					ws = append(ws, tws...)

					br.Read(make([]byte, ebl))
					idx += ebl

					continue
				}
			} else if cct == CodeC {
				if nb, err = br.Peek(1); err == nil {
					if isNum(nb[0]) {
						tmp = string(b) + string(nb)
						ret += MapEncodeCPattern[tmp]
						ws = append(ws, MapEncodeCodeCValue[tmp])

						idx += 2
						br.ReadByte()
						continue
					}
				}
			}
		}

		tmp = string(b)
		codeType, code, startCode, startCodeValue, codeTypeValue, val = locate(tmp)
		if val == -1 {
			return "", errors.New("unsupport string: " + tmp)
		}

		if codeType != cct {
			if idx == 0 {
				ret += startCode
				scv = startCodeValue
			} else {
				ret += codeType
				ws = append(ws, codeTypeValue)
			}

			cct = codeType
		}

		ret += code
		ws = append(ws, val)
		idx++
	}

	checksum += scv
	for k, v = range ws {
		checksum += v * (k + 1)
	}

	checksum = checksum % 103
	ret += MapEncodeChecksumPattern[checksum] + Stop

	return ret, nil
}

func drawRect(img *image.NRGBA, r image.Rectangle, c color.Color) {
	for i := r.Min.X; i <= r.Max.X; i++ {
		for j := r.Min.Y; j <= r.Max.Y; j++ {
			img.Set(i, j, c)
		}
	}
}

func makeImg(b string, h int, qx int, qy int, u int) *image.NRGBA {
	var (
		width  int
		height int
		canvas *image.NRGBA
		k      int
		v      rune
	)

	if u <= 0 {
		u = 1
	}

	if qx <= 0 {
		qx = 5
	}

	if qy <= 0 {
		qy = 3
	}

	if h <= 0 {
		h = 100
	}

	width = qx*2 + len(b)*u
	height = qy*2 + h
	canvas = image.NewNRGBA(image.Rect(0, 0, width, height))

	drawRect(canvas, image.Rect(0, 0, width, height), color.RGBA{255, 255, 255, 255})

	for k, v = range b {
		if v == 49 {
			drawRect(canvas, image.Rect(qx+k*u, qy, qx+k*u+u, qy+h), color.RGBA{0, 0, 0, 255})
		} else {
			drawRect(canvas, image.Rect(qx+k*u, qy, qx+k*u+u, qy+h), color.RGBA{255, 255, 255, 255})
		}
	}

	return canvas
}

func Encode(in []byte, h int, qx int, qy int, u int) (*image.NRGBA, error) {
	var (
		bc  string
		err error
	)

	if bc, err = encode(in); err != nil {
		return nil, err
	}

	return makeImg(bc, h, qx, qy, u), nil
}
