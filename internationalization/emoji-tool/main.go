package main

import (
	"fmt"
	"image"
	"image/gif"
	"io"
	"log"
	"os"

	"golang.org/x/image/font"
	"golang.org/x/image/font/opentype"
	"golang.org/x/image/math/fixed"
)

const AppleEmoji = "/System/Library/Fonts/Apple Color Emoji.ttc"

func main() {
	/*
		for _, i := range internationalization.Cultures {
			fmt.Printf("%s\n", i.Flag)
			build(i.Flag, i.CountryName+".png")
		}
	*/
	build("🇩🇪", "de.png")
}

func build(char, filename string) {
	const (
		width        = 400
		height       = 400
		startingDotX = 6
		startingDotY = 280
	)

	ttfR, err := os.Open(AppleEmoji)
	if err != nil {
		log.Fatal(err)
	}
	defer ttfR.Close()
	ttf, err := io.ReadAll(ttfR)
	fatal(err)

	f, err := opentype.ParseCollection(ttf)
	fatal(err)
	log.Printf("Collection has %d fonts", f.NumFonts())
	fnt, err := f.Font(1)
	fatal(err)

	face, err := opentype.NewFace(fnt, &opentype.FaceOptions{
		Size:    100,
		DPI:     72,
		Hinting: font.HintingNone,
	})
	fatal(err)

	dst := image.NewGray(image.Rect(0, 0, width, height))
	d := font.Drawer{
		Dst:  dst,
		Src:  image.White,
		Face: face,
		Dot:  fixed.P(startingDotX, startingDotY),
	}
	fmt.Printf("The dot is at %v\n", d.Dot)
	d.DrawString(char)

	fo, err := os.Create(filename)
	fatal(err)
	gif.Encode(fo, dst, nil)

	if err := gif.Encode(fo, dst, nil); err != nil {
		fo.Close()
		log.Fatal(err)
	}

	if err := fo.Close(); err != nil {
		log.Fatal(err)
	}
}

func fatal(err error) {
	if err != nil {
		log.Fatal(err)
	}
}
