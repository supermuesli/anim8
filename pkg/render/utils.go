package render

import (
	"image/png"
	"bytes"

	"github.com/faiface/pixel"
	"golang.org/x/image/font"
	"github.com/golang/freetype/truetype"
)

func loadPicture(data []byte) (pixel.Picture, error) {
	img, err := png.Decode(bytes.NewReader(data))
	if err != nil {
		return nil, err
	}
	return pixel.PictureDataFromImage(img), nil
}

func loadTTF(data []byte, size float64) (font.Face, error) {
	font, err := truetype.Parse(data)
	if err != nil {
		return nil, err
	}

	return truetype.NewFace(font, &truetype.Options{
		Size:              size,
		GlyphCacheEntries: 1,
	}), nil
}
