package app

import (
	"io"
	"log"
	"os"

	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/theme"
	"github.com/xujiajun/nutsdb"
)

type Cache struct {
	db *nutsdb.DB
}

func NewCache() *Cache {
	options := nutsdb.DefaultOptions
	options.SyncEnable = false

	cache, err := nutsdb.Open(
		options,
		nutsdb.WithDir("/tmp/comic"),
	)
	if err != nil {
		log.Fatal(err)
	}

	return &Cache{db: cache}
}

func (p *Cache) Close() {
	p.db.Close()
}

func (p *Cache) Get(path string) []byte {
	var res []byte

	key := []byte(path)

	_ = p.db.View(func(tx *nutsdb.Tx) error {
		entry, err := tx.Get("", key)
		res = entry.Value

		return err
	})

	return res
}

func (p *Cache) Put(path string, data []byte) {
	key := []byte(path)

	_ = p.db.Update(func(tx *nutsdb.Tx) error {
		return tx.Put("", key, data, 0)
	})
}

func (p *Cache) Error() *canvas.Image {
	img := canvas.NewImageFromResource(theme.ErrorIcon())
	img.FillMode = canvas.ImageFillStretch

	return img
}

func (p *Cache) Image(path string) *canvas.Image {
	data := p.Get(path)

	if len(data) == 0 {
		if data, _ = os.ReadFile(path); len(data) > 0 {
			p.Put(path, data)
			log.Println("cache put", path, len(data))
		} else {
			return p.Error()
		}
	}

	log.Println("cache read", path, len(data))

	image, err := ToImage(data)
	if err != nil {
		log.Println(path, data, err)

		return p.Error()
	}

	img := canvas.NewImageFromImage(image)
	img.FillMode = canvas.ImageFillStretch
	img.ScaleMode = canvas.ImageScaleFastest

	return img
}

func (p *Cache) Load(path string, reader io.ReadCloser) {
	isOld := true
	_ = p.db.View(func(tx *nutsdb.Tx) error {
		if _, err := tx.Get("", []byte(path)); err == nil {
			isOld = false
		}

		return nil
	})

	if isOld {
		return
	}

	data, err := io.ReadAll(reader)
	if err != nil {
		return
	}

	defer reader.Close()

	_ = p.db.Update(func(tx *nutsdb.Tx) error {
		return tx.Put("", []byte(path), data, 0)
	})
}

func (p *Cache) Loads(paths []string) {
	loads := make([]string, 0, len(paths))

	_ = p.db.View(func(tx *nutsdb.Tx) error {
		for _, path := range paths {
			if _, err := tx.Get("", []byte(path)); err != nil {
				loads = append(loads, path)
			}
		}

		return nil
	})

	_ = p.db.Update(func(tx *nutsdb.Tx) error {
		for _, path := range loads {
			data, _ := os.ReadFile(path)
			_ = tx.Put("", []byte(path), data, 0)
			log.Println("load:", path, len(data))
		}

		return nil
	})
}