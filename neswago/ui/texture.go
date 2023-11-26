package ui

import (
	"image"
	"path"
	"strings"

	"github.com/go-gl/gl/v2.1/gl"
)

const textureSize = 4096
const textureDim = textureSize / 256
const textureCount = textureDim * textureDim

type _PathIndexPair struct {
	Path  string
	Index int
}

type Texture struct {
	texture uint32
	lookup  []_PathIndexPair
	reverse [textureCount]string
	access  [textureCount]int
	counter int
}

func NewTexture() *Texture {
	texture := createTexture()
	gl.BindTexture(gl.TEXTURE_2D, texture)
	gl.TexImage2D(
		gl.TEXTURE_2D, 0, gl.RGBA,
		textureSize, textureSize,
		0, gl.RGBA, gl.UNSIGNED_BYTE, nil)
	gl.BindTexture(gl.TEXTURE_2D, 0)
	t := Texture{}
	t.texture = texture
	return &t
}

func (t *Texture) Bind() {
	gl.BindTexture(gl.TEXTURE_2D, t.texture)
}

func (t *Texture) Unbind() {
	gl.BindTexture(gl.TEXTURE_2D, 0)
}

func (t *Texture) Lookup(path string) (x, y, dx, dy float32) {
	for i := 0; i < len(t.lookup); i++ {
		if t.lookup[i].Path == path {
			index := t.lookup[i].Index
			return t.coord(index)
		}
	}
	return t.coord(t.load(path))
}

func (t *Texture) addLookpath(path string, index int) {
	t.lookup = append(t.lookup, _PathIndexPair{
		path, index,
	})
}

func (t *Texture) delLookpath(path string) {
	for i := 0; i < len(t.lookup); i++ {
		if t.lookup[i].Path == path {
			lastIdx := len(t.lookup) - 1
			t.lookup[i] = t.lookup[lastIdx]
			t.lookup = t.lookup[:lastIdx]
			return
		}
	}
}

func (t *Texture) mark(index int) {
	t.counter++
	t.access[index] = t.counter
}

func (t *Texture) lru() int {
	minIndex := 0
	minValue := t.counter + 1
	for i, n := range t.access {
		if n < minValue {
			minIndex = i
			minValue = n
		}
	}
	return minIndex
}

func (t *Texture) coord(index int) (x, y, dx, dy float32) {
	x = float32(index%textureDim) / textureDim
	y = float32(index/textureDim) / textureDim
	dx = 1.0 / textureDim
	dy = dx * 240 / 256
	return
}

func (t *Texture) load(path string) int {
	index := t.lru()
	t.delLookpath(t.reverse[index])
	t.mark(index)
	t.addLookpath(path, index)
	t.reverse[index] = path
	x := int32((index % textureDim) * 256)
	y := int32((index / textureDim) * 256)
	im := copyImage(t.loadThumbnail(path))
	size := im.Rect.Size()
	gl.TexSubImage2D(
		gl.TEXTURE_2D, 0, x, y, int32(size.X), int32(size.Y),
		gl.RGBA, gl.UNSIGNED_BYTE, gl.Ptr(im.Pix))
	return index
}

func (t *Texture) loadThumbnail(romPath string) image.Image {
	_, name := path.Split(romPath)
	name = strings.TrimSuffix(name, ".nes")
	name = strings.Replace(name, "_", " ", -1)
	name = strings.Title(name)
	return CreateGenericThumbnail(name)
}
