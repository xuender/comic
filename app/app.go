package app

import (
	"image/color"
	"log"
	"net/url"
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
)

type Mode int

const (
	ModeFit Mode = iota
	ModeFull
	ModeByWidth
	ModeByHeight
)

var modes = map[Mode]string{
	ModeFit:      "*",
	ModeFull:     "/",
	ModeByWidth:  "W",
	ModeByHeight: "H",
}

type App struct {
	app    fyne.App
	main   fyne.Window
	center *fyne.Container
	border *fyne.Container
	img    *canvas.Image
	cache  *Cache
	files  *Files
	help   dialog.Dialog
	mode   Mode
	radio  *widget.RadioGroup
	scroll *container.Scroll
	show   func()
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
		scroll: scroll,
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
		widget.NewToolbarAction(theme.HelpIcon(), app.showHelp),
	)

	options := map[string]func(){
		"*": app.refresh,
		"/": app.full,
		"W": app.width,
		"H": app.height,
	}
	radio := widget.NewRadioGroup([]string{"*", "/", "W", "H"}, func(option string) {
		if call, has := options[option]; has {
			call()
		}

		log.Println(option)
	})
	radio.Horizontal = true
	border := container.NewBorder(toolbar, radio, nil, nil, scroll)
	// nolint: gomnd
	main.Resize(fyne.NewSize(800, 600))
	main.SetContent(border)
	main.CenterOnScreen()

	app.border = border
	app.radio = radio
	app.help = app.createHelp()
	app.show = Debounced(app.Show, time.Millisecond*300)

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
		fyne.KeyQ:        p.app.Quit,
		fyne.KeySlash:    p.full,
		fyne.KeyM:        p.full,
		fyne.KeyHome:     p.home,
		fyne.KeyEnd:      p.end,
		fyne.KeyW:        p.width,
		fyne.KeyH:        p.height,
		fyne.KeyF1:       p.showHelp,
		fyne.KeyF10:      p.showHelp,
		fyne.KeyEscape:   p.hideHelp,
		fyne.KeyUp:       p.top,
		fyne.KeyDown:     p.bottom,
	}

	p.main.Canvas().SetOnTypedKey(func(ke *fyne.KeyEvent) {
		if call, has := funcs[ke.Name]; has {
			call()
			log.Println(ke.Name)
		}
	})
	p.SetMode(ModeFit)
}

func (p *App) Run(args []string) {
	defer p.cache.Close()
	// p.files.Load([]string{"doc/logo.png", "doc/maskable_icon.png"})
	// p.files.Load([]string{"doc/a.zip", "doc/logo.png", "doc/maskable_icon.png"})
	// p.files.Load([]string{"doc"})
	// p.files.Load([]string{"doc/a.zip"})
	p.files.Load(args)

	p.init()
	p.Show()
	p.main.Show()
	p.app.Run()
}

func (p *App) Show() {
	log.Println("Show")
	if p.files.Len() < 1 {
		p.img = NoneImage()
		// nolint: gomnd
		p.img.SetMinSize(fyne.NewSize(400, 400))
		p.img.FillMode = canvas.ImageFillStretch
		p.img.ScaleMode = canvas.ImageScaleFastest
		p.reset()
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
	p.top()

	p.reset()

	p.main.Canvas().SetContent(p.border)
	p.main.SetTitle(path)
}

func (p *App) createHelp() dialog.Dialog {
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
		canvas.NewText("To Top", white), canvas.NewText("Up", grey), canvas.NewText("", grey),
		canvas.NewText("To Bottom", white), canvas.NewText("Down", grey), canvas.NewText("", grey),
		canvas.NewText("Help", white), canvas.NewText("F1", grey), canvas.NewText("F10", grey),
		canvas.NewText("Quit", white), canvas.NewText("Q", grey), canvas.NewText("", grey),
		canvas.NewText("", white), canvas.NewText("", grey), widget.NewHyperlink("xuender/comic", url),
	)

	return dialog.NewCustom("Help", "Close", grid, p.main)
}

func (p *App) reset() {
	switch p.mode {
	case ModeFull:
		p.full()
	case ModeByWidth:
		p.width()
	case ModeByHeight:
		p.height()
	default:
	}
}

func (p *App) showHelp() {
	p.help.Show()
}

func (p *App) hideHelp() {
	if p.main.FullScreen() {
		p.fullScreen()
	}

	p.help.Hide()
}

func (p *App) top() {
	log.Println("top")
	p.scroll.Scrolled(&fyne.ScrollEvent{Scrolled: fyne.Delta{DX: 0, DY: p.scroll.Offset.Y}})
}

func (p *App) bottom() {
	log.Println("bottom")
	p.scroll.ScrollToBottom()
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
	if p.main.FullScreen() {
		p.border.Objects[1].Hide()
		p.border.Objects[2].Hide()

		return
	}

	p.border.Objects[1].Show()
	p.border.Objects[2].Show()
}

func (p *App) refresh() {
	p.SetMode(ModeFit)
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

func (p *App) SetMode(mode Mode) {
	p.mode = mode
	log.Println("mode", modes[p.mode])
	p.radio.SetSelected(modes[p.mode])
}

func (p *App) full() {
	p.SetMode(ModeFull)
	size := ToSize(p.img.Size(), p.main.Canvas().Size())
	p.img.Resize(size)
	p.img.SetMinSize(size)

	log.Println("full")
}

func (p *App) width() {
	p.SetMode(ModeByWidth)
	width := p.main.Canvas().Size().Width
	height := width / p.img.Size().Width * p.img.Size().Height
	size := fyne.NewSize(width, height)
	p.img.Resize(size)
	p.img.SetMinSize(size)
}

func (p *App) height() {
	p.SetMode(ModeByHeight)
	height := p.main.Canvas().Size().Height
	width := height / p.img.Size().Height * p.img.Size().Width
	size := fyne.NewSize(width, height)
	p.img.Resize(size)
	p.img.SetMinSize(size)
}
