package app

import (
	"io/fs"
	"log"
	"os"
	"path/filepath"
	"sort"

	"github.com/h2non/filetype"
	"github.com/mholt/archiver/v4"
)

const _headSize = 262

type Book struct {
	path string
	subs []string
}

type Files struct {
	cache *Cache
	books []*Book
	index int
	book  int
}

func NewFiles(cache *Cache) *Files {
	return &Files{
		books: []*Book{},
		cache: cache,
	}
}

func (p *Files) Len() int {
	res := 0

	for _, book := range p.books {
		if len(book.subs) > 0 {
			res += len(book.subs)
		} else {
			res++
		}
	}

	return res
}

func (p *Files) Load(paths []string) {
	for _, path := range paths {
		if abs, err := filepath.Abs(path); err == nil {
			if books := p.ReadBooks(abs); books != nil {
				p.books = append(p.books, books...)
			}
		}
	}

	p.logs()
}

func (p *Files) logs() {
	for i, book := range p.books {
		log.Println(i, book.path)

		for f, sub := range book.subs {
			log.Println(i, f, sub)
		}
	}
}

func (p *Files) Get() string {
	if len(p.books) == 0 {
		return ""
	}

	book := p.books[p.book]
	if len(book.subs) > 0 {
		return p.books[p.book].subs[p.index]
	}

	log.Println(book.path)

	return book.path
}

func (p *Files) Start() {
	p.index = 0
}

func (p *Files) End() {
	if length := len(p.books[p.book].subs); length > 0 {
		p.index = length - 1
	}
}

func (p *Files) Down() {
	p.index++

	if p.index >= len(p.books[p.book].subs) {
		p.book++

		if p.book >= len(p.books) {
			p.book = 0
		}

		p.index = 0
	}
}

func (p *Files) Up() {
	p.index--

	if p.index < 0 {
		p.book--

		if p.book < 0 {
			p.book = len(p.books) - 1
		}

		p.index = len(p.books[p.book].subs) - 1
	}
}

func (p *Files) readDir(path string) []*Book {
	files, err := os.ReadDir(path)
	if err != nil {
		return nil
	}

	res := []*Book{}
	book := &Book{path: path, subs: []string{}}

	for _, file := range files {
		name := filepath.Join(path, file.Name())

		if file.IsDir() {
			res = append(res, p.readDir(name)...)

			continue
		}

		if IsImage(name) {
			if reader, err := os.Open(name); err == nil {
				go p.cache.Load(name, reader)
				book.subs = append(book.subs, name)
			}

			continue
		}

		if IsArchive(name) {
			res = append(res, p.readArchive(name)...)
		}
	}

	if len(book.subs) > 0 {
		sort.Strings(book.subs)
		res = append([]*Book{book}, res...)
	}

	return res
}

func (p *Files) readArchive(mainPath string) []*Book {
	log.Println("read achive", mainPath)
	if fsys, err := archiver.FileSystem(mainPath); err == nil {
		res := []string{}

		_ = fs.WalkDir(fsys, ".", func(path string, d fs.DirEntry, err error) error {
			if err != nil {
				log.Println(err)

				return err
			}

			if d.IsDir() {
				return nil
			}

			file, err := fsys.Open(path)
			if err != nil {
				return nil
			}

			log.Println(path)

			head := make([]byte, _headSize)
			_, _ = file.Read(head)
			file.Close()

			if filetype.IsImage(head) {
				name := filepath.Join(mainPath, path)
				file, _ := fsys.Open(path)

				go p.cache.Load(name, file)

				res = append(res, name)
			}

			return nil
		})

		sort.Strings(res)
		log.Println(res)

		return []*Book{{path: mainPath, subs: res}}
	}

	return nil
}

func (p *Files) ReadBooks(mainPath string) []*Book {
	log.Println("ReadBooks", mainPath)
	stat, err := os.Stat(mainPath)
	if os.IsNotExist(err) {
		return nil
	}

	if stat.IsDir() {
		return p.readDir(mainPath)
	}

	if IsImage(mainPath) {
		return []*Book{{path: mainPath}}
	}

	if IsArchive(mainPath) {
		return p.readArchive(mainPath)
	}

	log.Println("pass", mainPath)

	return nil
}
