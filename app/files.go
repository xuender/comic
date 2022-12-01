package app

import (
	"io/fs"
	"log"
	"os"
	"path/filepath"
	"sort"

	"github.com/h2non/filetype"
	"github.com/mholt/archiver/v4"
	"github.com/samber/lo"
)

const (
	_headSize = 262
	_wait     = "wait"
)

type Files struct {
	cache *Cache
	ends  []int
	pages []string
	index int
}

func NewFiles(cache *Cache) *Files {
	return &Files{
		ends:  []int{},
		pages: []string{},
		cache: cache,
	}
}

func (p *Files) Len() int {
	return len(p.pages)
}

func (p *Files) Load(paths []string) {
	paths = lo.Map(paths, func(path string, _ int) string {
		abs, _ := filepath.Abs(path)

		return abs
	})
	paths = lo.Filter(paths, func(path string, _ int) bool { return path != "" })

	cha := make(chan string)

	go func(input chan<- string) {
		for _, path := range paths {
			if abs, err := filepath.Abs(path); err == nil {
				p.ReadBooks(abs, input)
			}
		}
	}(cha)

	p.pages = append(p.pages, ReadPage(cha))

	go func() {
		for page := range cha {
			if page == "" {
				if length := len(p.pages) - 1; len(p.ends) == 0 || p.ends[len(p.ends)-1] < length {
					p.ends = append(p.ends, length)
				}

				continue
			}

			if p.pages[0] == _wait {
				p.pages[0] = page

				continue
			}

			p.pages = append(p.pages, page)
		}

		p.logs()
	}()
}

func (p *Files) logs() {
	for i, page := range p.pages {
		log.Println(i, page)
	}
}

func (p *Files) Get() string {
	if len(p.pages) == 0 {
		return _wait
	}

	return p.pages[p.index]
}

func (p *Files) Start() {
	max := -1

	for _, end := range p.ends {
		if p.index <= end {
			p.index = max + 1

			return
		}

		max = end
	}

	p.index = 0
}

func (p *Files) End() {
	for _, end := range p.ends {
		if p.index <= end {
			p.index = end

			return
		}
	}

	p.index = len(p.pages) - 1
}

func (p *Files) Down() {
	p.index++

	if p.index >= len(p.pages) {
		p.index = 0
	}
}

func (p *Files) Up() {
	p.index--

	if p.index < 0 {
		p.index = len(p.pages) - 1
	}
}

func (p *Files) ReadBooks(mainPath string, input chan<- string) {
	defer close(input)

	log.Println("ReadBooks", mainPath)

	stat, err := os.Stat(mainPath)

	if os.IsNotExist(err) {
		return
	}

	if stat.IsDir() {
		p.readDir(mainPath, input)

		return
	}

	if IsImage(mainPath) {
		if reader, err := os.Open(mainPath); err == nil {
			p.cache.Load(mainPath, reader)
			input <- mainPath
		}

		return
	}

	if IsArchive(mainPath) {
		p.readArchive(mainPath, input)

		return
	}

	log.Println("pass", mainPath)
}

func (p *Files) readDir(path string, input chan<- string) {
	files, err := os.ReadDir(path)
	if err != nil {
		return
	}

	sort.Slice(files, func(i, j int) bool {
		if files[i].IsDir() == files[j].IsDir() {
			return files[i].Name() < files[j].Name()
		}

		return files[i].IsDir()
	})

	for _, file := range files {
		name := filepath.Join(path, file.Name())

		if file.IsDir() {
			p.readDir(name, input)

			continue
		}

		if IsImage(name) {
			if reader, err := os.Open(name); err == nil {
				p.cache.Load(name, reader)
				input <- name
			}

			continue
		}

		if IsArchive(name) {
			p.readArchive(name, input)
		}
	}

	input <- ""
}

func (p *Files) readArchive(mainPath string, input chan<- string) {
	log.Println("read archive", mainPath)

	if fsys, err := archiver.FileSystem(mainPath); err == nil {
		paths := []string{}
		_ = fs.WalkDir(fsys, ".", func(path string, entry fs.DirEntry, err error) error {
			if err != nil {
				log.Println(err)

				return err
			}

			if entry.IsDir() {
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
				paths = append(paths, path)
			}

			return nil
		})

		sort.Strings(paths)

		for _, path := range paths {
			name := filepath.Join(mainPath, path)
			file, _ := fsys.Open(path)

			p.cache.Load(name, file)
			input <- name
		}
	}

	input <- ""
}
