# Barcode

Generates 1D-Barcode in golang.

## Usage

```go
ts := "HI3456HI"

img, err := Encode([]byte(ts), 0, 0, 0, 2)

if err != nil {
	log.Fatal(err)
}

if file, err := os.OpenFile(ts+".png", os.O_RDWR|os.O_CREATE, 0666); err != nil {
	log.Fatal(err)
} else {
	png.Encode(file, img)
}
```

## Restrict

Only one supports type is [code128](http://en.wikipedia.org/wiki/Code_128), and **Excludes** FNCx codes.
