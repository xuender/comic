package app

import (
	"image/color"
	"log"
	"net/url"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
)

type App struct {
	app    fyne.App
	main   fyne.Window
	center *fyne.Container
	border *fyne.Container
	img    *canvas.Image
	cache  *Cache
	files  *Files
}

func NewApp(
	cache *Cache,
	files *Files,
) *App {
	fyneApp := app.New()
	main := fyneApp.NewWindow("Comic")
	img := NoneImage()
	center := container.New(layout.NewCenterLayout(), img)
	scroll := container.NewScroll(center)

	app := &App{
		app:    fyneApp,
		main:   main,
		center: center,
		img:    img,
		cache:  cache,
		files:  files,
	}

	toolbar := widget.NewToolbar(
		widget.NewToolbarAction(theme.HomeIcon(), app.home),
		widget.NewToolbarAction(theme.NavigateBackIcon(), app.back),
		widget.NewToolbarAction(theme.NavigateNextIcon(), app.next),
		widget.NewToolbarSeparator(),
		widget.NewToolbarAction(theme.ViewFullScreenIcon(), app.full),
		widget.NewToolbarAction(theme.ViewRefreshIcon(), app.refresh),
		widget.NewToolbarAction(theme.ZoomInIcon(), app.zoomIn),
		widget.NewToolbarAction(theme.ZoomOutIcon(), app.zoomOut),
		widget.NewToolbarSpacer(),
		widget.NewToolbarAction(theme.HelpIcon(), app.help),
	)
	border := container.NewBorder(toolbar, nil, nil, nil, scroll)
	// nolint: gomnd
	main.Resize(fyne.NewSize(800, 600))
	main.SetContent(border)
	main.CenterOnScreen()

	app.border = border

	return app
}

func (p *App) init() {
	funcs := map[fyne.KeyName]func(){
		fyne.KeyF11:      p.fullScreen,
		fyne.Key8:        p.refresh,
		fyne.KeyAsterisk: p.refresh,
		fyne.KeyPlus:     p.zoomIn,
		fyne.KeyEqual:    p.zoomIn,
		fyne.KeyMinus:    p.zoomOut,
		fyne.KeyPageDown: p.next,
		fyne.KeySpace:    p.next,
		fyne.KeyPageUp:   p.back,
		fyne.KeyE:        p.app.Quit,
		fyne.KeyEscape:   p.app.Quit,
		fyne.KeySlash:    p.full,
		fyne.KeyM:        p.full,
		fyne.KeyHome:     p.home,
		fyne.KeyEnd:      p.end,
		fyne.KeyW:        p.width,
		fyne.KeyH:        p.height,
		fyne.KeyF1:       p.help,
		fyne.KeyF10:      p.help,
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
	// p.files.Load([]string{"doc/logo.png", "doc/maskable_icon.png"})
	// p.files.Load([]string{"doc/a.zip", "doc/logo.png", "doc/maskable_icon.png"})
	// p.files.Load([]string{"doc"})
	p.files.Load(args)

	p.init()
	p.show()
	p.main.Show()
	p.app.Run()
}

func (p *App) show() {
	if p.files.Len() < 1 {
		p.img = NoneImage()
		// nolint: gomnd
		p.img.SetMinSize(fyne.NewSize(400, 400))
		p.img.FillMode = canvas.ImageFillStretch
		p.img.ScaleMode = canvas.ImageScaleFastest
		p.main.Canvas().SetContent(p.border)
		p.main.SetTitle("Comic")

		return
	}

	path := p.files.Get()
	p.img = p.cache.Image(path)
	// p.img = Image(path)
	if p.img.Image != nil {
		p.img.SetMinSize(fyne.NewSize(float32(p.img.Image.Bounds().Dx()), float32(p.img.Image.Bounds().Dy())))
	}

	p.center.RemoveAll()
	p.center.Add(p.img)

	p.main.Canvas().SetContent(p.border)
	p.main.SetTitle(path)
}

func (p *App) help() {
	white := color.White
	grey := color.Gray16{0x8888}
	url, _ := url.Parse("https://github.com/xuender/comic")
	three := 3
	grid := container.New(layout.NewGridLayout(three),
		canvas.NewText("Function", white), canvas.NewText("Key1", white), canvas.NewText("Key2", white),
		canvas.NewText("Full Screen", white), canvas.NewText("F11", grey), canvas.NewText("", grey),
		canvas.NewText("Refresh", white), canvas.NewText("*", grey), canvas.NewText("8", grey),
		canvas.NewText("Zoom In", white), canvas.NewText("+", grey), canvas.NewText("=", grey),
		canvas.NewText("Zoom Out", white), canvas.NewText("-", grey), canvas.NewText("", grey),
		canvas.NewText("Page Down", white), canvas.NewText("PageDown", grey), canvas.NewText("Space", grey),
		canvas.NewText("Page Up", white), canvas.NewText("PageUp", grey), canvas.NewText("", grey),
		canvas.NewText("Max", white), canvas.NewText("/", grey), canvas.NewText("M", grey),
		canvas.NewText("Width", white), canvas.NewText("W", grey), canvas.NewText("", grey),
		canvas.NewText("Height", white), canvas.NewText("H", grey), canvas.NewText("", grey),
		canvas.NewText("First Page", white), canvas.NewText("Home", grey), canvas.NewText("", grey),
		canvas.NewText("Last Page", white), canvas.NewText("End", grey), canvas.NewText("", grey),
		canvas.NewText("Help", white), canvas.NewText("F1", grey), canvas.NewText("F10", grey),
		canvas.NewText("Exit", white), canvas.NewText("Esc", grey), canvas.NewText("E", grey),
		canvas.NewText("", white), canvas.NewText("", grey),
		widget.NewHyperlink("xuender/comic", url),
	)

	dialog.ShowCustom("Help", "Close", grid, p.main)
}

func (p *App) next() {
	p.files.Down()
	p.show()
}

func (p *App) home() {
	p.files.Start()
	p.show()
}

func (p *App) end() {
	p.files.End()
	p.show()
}

func (p *App) back() {
	p.files.Up()
	p.show()
}

func NoneImage() *canvas.Image {
	img := canvas.NewImageFromResource(theme.FyneLogo())
	img.FillMode = canvas.ImageFillContain

	return img
}

func (p *App) fullScreen() {
	p.main.SetFullScreen(!p.main.FullScreen())
}

func (p *App) refresh() {
	// 原始尺寸
	if p.img.Image == nil {
		return
	}

	size := fyne.NewSize(float32(p.img.Image.Bounds().Dx()), float32(p.img.Image.Bounds().Dy()))
	p.img.Resize(size)
	p.img.SetMinSize(size)
}

func (p *App) zoomIn() {
	size := fyne.NewSize(
		p.img.Size().Width+p.img.Size().Width*0.1,
		p.img.Size().Height+p.img.Size().Height*0.1,
	)
	p.img.Resize(size)
	p.img.SetMinSize(size)
}

func (p *App) zoomOut() {
	size := fyne.NewSize(
		p.img.Size().Width-p.img.Size().Width*0.1,
		p.img.Size().Height-p.img.Size().Height*0.1,
	)
	p.img.Resize(size)
	p.img.SetMinSize(size)
}

func (p *App) full() {
	size := ToSize(p.img.Size(), p.main.Canvas().Size())
	p.img.Resize(size)
	p.img.SetMinSize(size)
}

func (p *App) width() {
	width := p.main.Canvas().Size().Width
	height := width / p.img.Size().Width * p.img.Size().Height
	size := fyne.NewSize(width, height)
	p.img.Resize(size)
	p.img.SetMinSize(size)
}

func (p *App) height() {
	height := p.main.Canvas().Size().Height
	width := height / p.img.Size().Height * p.img.Size().Width
	size := fyne.NewSize(width, height)
	p.img.Resize(size)
	p.img.SetMinSize(size)
}
