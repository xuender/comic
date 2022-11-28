package app

import (
	"log"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/theme"
)

type App struct {
	app     fyne.App
	main    fyne.Window
	content *fyne.Container
	img     *canvas.Image
	cache   *Cache
	scroll  *container.Scroll
	paths   []string
	index   int
}

func NewApp(
	cache *Cache,
) *App {
	fyneApp := app.New()
	main := fyneApp.NewWindow("Comic")
	img := NoneImage()
	content := container.New(layout.NewCenterLayout(), img)
	scroll := container.NewScroll(content)

	main.Resize(fyne.NewSize(500, 500))
	scroll.Resize(fyne.NewSize(600, 600))
	main.SetContent(scroll)
	main.CenterOnScreen()

	return &App{
		app:     fyneApp,
		main:    main,
		content: content,
		scroll:  scroll,
		img:     img,
		cache:   cache,
		index:   0,
	}
}

func (p *App) init() {
	p.main.Canvas().SetOnTypedKey(func(ke *fyne.KeyEvent) {
		switch ke.Name {
		case fyne.KeyF11:
			// 全屏
			p.main.SetFullScreen(!p.main.FullScreen())
		case fyne.Key8, fyne.KeyAsterisk:
			// 原始尺寸
			if p.img.Image == nil {
				return
			}

			p.img.FillMode = canvas.ImageFillContain
			p.img.Resize(fyne.NewSize(float32(p.img.Image.Bounds().Dx()), float32(p.img.Image.Bounds().Dy())))
			p.img.FillMode = canvas.ImageFillOriginal
			log.Println("*")
		case fyne.KeyPlus, fyne.KeyEqual:
			// 放大
			size := p.img.Size()
			w := size.Width + size.Width*0.1
			h := size.Height + size.Width*0.1
			p.img.FillMode = canvas.ImageFillContain
			newSize := fyne.NewSize(float32(w), float32(h))
			p.img.Resize(newSize)
			p.img.FillMode = canvas.ImageFillOriginal
			// p.content.Resize(newSize)
			old := p.content.Position()
			log.Println(p.img.Size(), p.content.Size(), old)
			// p.content.RemoveAll()
			// p.content.Add(p.img)
			// p.scroll.Content = p.content
			// pos := fyne.NewPos((p.scroll.Size().Width-newSize.Width)/2, (p.scroll.Size().Height-newSize.Height)/2)
			// pos := fyne.NewPos(0, 0)
			// p.content.Move(pos)
			// log.Println("+", size, p.scroll.Size(), pos, old)
		case fyne.KeyMinus:
			// 缩小
			size := p.img.Size()
			w := size.Width - size.Width*0.1
			h := size.Height - size.Width*0.1
			p.img.FillMode = canvas.ImageFillContain
			p.img.Resize(fyne.NewSize(float32(w), float32(h)))
			p.img.FillMode = canvas.ImageFillOriginal
			log.Println("-")
		case fyne.KeyPageDown:
			// 下页
			p.down()
		case fyne.KeyPageUp:
			// 上页
			p.up()
		case fyne.KeyQ, fyne.KeyEscape:
			// 退出
			p.app.Quit()
		}
	})
}

func (p *App) Run(args []string) {
	defer p.cache.Close()

	// p.paths = args
	p.paths = []string{"doc/logo.png", "doc/maskable_icon.png"}

	p.init()
	p.show()
	p.main.Show()
	p.app.Run()
}

func (p *App) show() {
	if len(p.paths) < 1 {
		p.img = NoneImage()
		p.content.RemoveAll()
		p.content.Add(p.img)
		p.main.Canvas().SetContent(p.content)
		p.main.SetTitle("Comic")

		return
	}

	path := p.paths[p.index]
	p.img = Image(path)
	if p.img.Image != nil {
		p.img.Resize(fyne.NewSize(float32(p.img.Image.Bounds().Dx()), float32(p.img.Image.Bounds().Dy())))
	}

	p.content.RemoveAll()
	p.content.Add(p.img)
	p.scroll.Content = p.content
	p.main.Canvas().SetContent(p.scroll)
	p.main.SetTitle(path)
}

func Image(path string) *canvas.Image {
	image, _ := ReadImage(path)
	log.Println(image.Bounds())
	img := canvas.NewImageFromImage(image)
	// img := canvas.NewImageFromFile(path)
	// img.FillMode = canvas.ImageFillContain
	img.FillMode = canvas.ImageFillOriginal

	return img
}

func (p *App) down() {
	p.index++

	if p.index >= len(p.paths) {
		p.index = 0
	}

	p.show()
}

func (p *App) up() {
	p.index--

	if p.index < 0 {
		p.index = len(p.paths) - 1
	}

	p.show()
}

func NoneImage() *canvas.Image {
	img := canvas.NewImageFromResource(theme.FyneLogo())
	img.FillMode = canvas.ImageFillContain

	return img
}
