package app

import (
	"bytes"
	"encoding/gob"
	"log"

	"fyne.io/fyne/v2/canvas"
	"github.com/xujiajun/nutsdb"
)

type Cache struct {
	db *nutsdb.DB
}

func NewCache() *Cache {
	options := nutsdb.DefaultOptions
	options.SyncEnable = false

	db, err := nutsdb.Open(
		options,
		nutsdb.WithDir("/tmp/comic"),
	)
	if err != nil {
		log.Fatal(err)
	}

	return &Cache{db: db}
}

func (p *Cache) Close() {
	p.db.Close()
}

func (p *Cache) Get(path string) []byte {
	var res []byte

	key := []byte(path)

	p.db.View(func(tx *nutsdb.Tx) error {
		if entry, err := tx.Get("", key); err != nil {
			return err
		} else {
			res = entry.Value

			return nil
		}
	})

	return res
}

func (p *Cache) Put(path string, data []byte) {
	key := []byte(path)

	p.db.Update(func(tx *nutsdb.Tx) error {
		return tx.Put("", key, data, 0)
	})
}

func (p *Cache) Image(path string) *canvas.Image {
	res := &canvas.Image{}

	if data := p.Get(path); len(data) > 0 {
		buf := bytes.NewBuffer(data)
		decoder := gob.NewDecoder(buf)

		_ = decoder.Decode(res)

		log.Println("cache data", len(data))

		return res
	}

	// TODO 需要读取图像缓存
	// canvas.NewImageFromImage(src)
	img := canvas.NewImageFromFile(path)
	img.FillMode = canvas.ImageFillContain
	buf := &bytes.Buffer{}
	encode := gob.NewEncoder(buf)

	_ = encode.Encode(img)

	p.Put(path, buf.Bytes())

	return img
}

func Encode(img *canvas.Image) []byte {
	buf := &bytes.Buffer{}
	encode := gob.NewEncoder(buf)

	_ = encode.Encode(img)
	return buf.Bytes()
}

func Decode(data []byte) *canvas.Image {
	buf := bytes.NewBuffer(data)
	decoder := gob.NewDecoder(buf)
	res := &canvas.Image{}
	_ = decoder.Decode(res)

	return res
}
