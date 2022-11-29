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
	app    fyne.App
	main   fyne.Window
	center *fyne.Container
	img    *canvas.Image
	cache  *Cache
	scroll *container.Scroll
	paths  []string
	index  int
}

func NewApp(
	cache *Cache,
) *App {
	fyneApp := app.New()
	main := fyneApp.NewWindow("Comic")
	img := NoneImage()
	center := container.New(layout.NewCenterLayout(), img)
	scroll := container.NewScroll(center)
	// nolint: gomnd
	main.Resize(fyne.NewSize(800, 600))
	main.SetContent(scroll)
	main.CenterOnScreen()

	return &App{
		app:    fyneApp,
		main:   main,
		center: center,
		scroll: scroll,
		img:    img,
		cache:  cache,
		index:  0,
	}
}

func (p *App) fullScreen() {
	p.main.SetFullScreen(!p.main.FullScreen())
}

func (p *App) original() {
	// 原始尺寸
	if p.img.Image == nil {
		return
	}

	size := fyne.NewSize(float32(p.img.Image.Bounds().Dx()), float32(p.img.Image.Bounds().Dy()))
	p.img.Resize(size)
	p.img.SetMinSize(size)
}

func (p *App) plus() {
	size := fyne.NewSize(
		p.img.Size().Width+p.img.Size().Width*0.1,
		p.img.Size().Height+p.img.Size().Height*0.1,
	)
	p.img.Resize(size)
	p.img.SetMinSize(size)
}

func (p *App) minus() {
	size := fyne.NewSize(
		p.img.Size().Width-p.img.Size().Width*0.1,
		p.img.Size().Height-p.img.Size().Height*0.1,
	)
	p.img.Resize(size)
	p.img.SetMinSize(size)
}

func (p *App) max() {
	size := ToSize(p.img.Size(), p.main.Canvas().Size())
	p.img.Resize(size)
	p.img.SetMinSize(size)
}

func (p *App) init() {
	funcs := map[fyne.KeyName]func(){
		// 全屏
		fyne.KeyF11: p.fullScreen,
		// 原始尺寸
		fyne.Key8:        p.original,
		fyne.KeyAsterisk: p.original,
		// 放大
		fyne.KeyPlus:  p.plus,
		fyne.KeyEqual: p.plus,
		// 缩小
		fyne.KeyMinus: p.minus,
		// 上页
		fyne.KeyPageDown: p.down,
		// 下页
		fyne.KeyPageUp: p.up,
		// 退出
		fyne.KeyQ:      p.app.Quit,
		fyne.KeyEscape: p.app.Quit,
		// 最大尺寸
		fyne.KeySlash: p.max,
	}

	p.main.Canvas().SetOnTypedKey(func(ke *fyne.KeyEvent) {
		if call, has := funcs[ke.Name]; has {
			call()
			log.Println(ke.Name)
		}
	})
}

func (p *App) Run(args []string) {
	defer p.cache.Close()

	// p.paths = args
	p.paths = []string{"doc/logo.png", "doc/maskable_icon.png"}
	// 加载缓存
	go p.cache.Load(p.paths)

	p.init()
	p.show()
	p.main.Show()
	p.app.Run()
}

func (p *App) show() {
	if len(p.paths) < 1 {
		p.img = NoneImage()
		// nolint: gomnd
		p.img.SetMinSize(fyne.NewSize(400, 400))
		p.img.FillMode = canvas.ImageFillStretch
		p.img.ScaleMode = canvas.ImageScaleFastest
		p.main.Canvas().SetContent(p.scroll)
		p.main.SetTitle("Comic")

		return
	}

	path := p.paths[p.index]
	p.img = p.cache.Image(path)
	// p.img = Image(path)
	if p.img.Image != nil {
		p.img.SetMinSize(fyne.NewSize(float32(p.img.Image.Bounds().Dx()), float32(p.img.Image.Bounds().Dy())))
	}

	p.center.RemoveAll()
	p.center.Add(p.img)

	p.main.Canvas().SetContent(p.scroll)
	p.main.SetTitle(path)
}

func Image(path string) *canvas.Image {
	image, _ := ReadImage(path)
	log.Println(image.Bounds())
	img := canvas.NewImageFromImage(image)
	img.FillMode = canvas.ImageFillStretch
	img.ScaleMode = canvas.ImageScaleFastest

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
