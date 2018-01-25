package imagediff

import (
	"image"
	"image/png"
	"log"
	"os"
)

// Args are uint32, but coming from Color::RGBA, they contain just an uint16.
// Return is just a byte.
func absDiff(a, b uint32) byte {
	c, d := int32(a), int32(b)
	if c > d {
		return uint8((c - d) >> 8)
	} else {
		return uint8((d - c) >> 8)
	}
}

func Diff(a, b image.Image) image.Image {
	abox, bbox := a.Bounds(), b.Bounds()
	axlen, bxlen := abox.Max.X-abox.Min.X, bbox.Max.X-bbox.Min.X
	aylen, bylen := abox.Max.Y-abox.Min.Y, bbox.Max.Y-bbox.Min.Y
	if axlen != bxlen {
		log.Panicf("widths different: %d vs %d", axlen, bxlen)
	}
	if aylen != bylen {
		log.Panicf("widths different: %d vs %d", aylen, bylen)
	}

	zbox := image.Rectangle{image.Point{0, 0}, image.Point{axlen, aylen}}
	z := image.NewRGBA(zbox)

	for x := 0; x < axlen; x++ {
		for y := 0; y < aylen; y++ {
			ar, ag, ab, _ := a.At(x, y).RGBA()
			br, bg, bb, _ := b.At(x, y).RGBA()
			z.Pix[y*z.Stride+x*4+0] = absDiff(ar, br)
			z.Pix[y*z.Stride+x*4+1] = absDiff(ag, bg)
			z.Pix[y*z.Stride+x*4+2] = absDiff(ab, bb)
			z.Pix[y*z.Stride+x*4+3] = 255
		}
	}
	return z
}

func DiffFilenames(f1, f2, f3 string) {
	ar, err := os.Open(f1)
	if err != nil {
		log.Fatalf("Cannot open %q: %v", f1, err)
	}
	br, err := os.Open(f2)
	if err != nil {
		log.Fatalf("Cannot open %q: %v", f2, err)
	}

	ai, _, err := image.Decode(ar)
	if err != nil {
		log.Fatalf("Cannot decode %q: %v", f1, err)
	}
	bi, _, err := image.Decode(br)
	if err != nil {
		log.Fatalf("Cannot decode %q: %v", f2, err)
	}
	zi := Diff(ai, bi)

	zw, err := os.Create(f3)
	if err != nil {
		log.Fatalf("Cannot create %q: %v", f3, err)
	}
	err = png.Encode(zw, zi)
	if err != nil {
		log.Fatalf("Cannot encode %q: %v", f3, err)
	}
	zw.Close()
}
