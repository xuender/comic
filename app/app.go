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
	app      fyne.App
	main     fyne.Window
	center   *fyne.Container
	border   *fyne.Container
	img      *canvas.Image
	cache    *Cache
	files    *Files
	help     dialog.Dialog
	mode     Mode
	radio    *widget.RadioGroup
	scroll   *container.Scroll
	show     func()
	commands []Command
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

	commands := []Command{
		{Help: "First Page", Call: app.home, Icon: theme.MediaSkipPreviousIcon(), Key1: fyne.KeyHome},
		{Help: "Page Up", Call: app.back, Icon: theme.NavigateBackIcon(), Key1: fyne.KeyPageUp},
		{Help: "Page Down", Call: app.next, Icon: theme.NavigateNextIcon(), Key1: fyne.KeyPageDown, Key2: fyne.KeySpace},
		{Help: "Last Page", Call: app.end, Icon: theme.MediaSkipNextIcon(), Key1: fyne.KeyEnd},
		{Help: "Separator"},
		{Help: "Full Screen", Call: app.fullScreen, Icon: theme.ViewFullScreenIcon(), Key1: fyne.KeyF11},
		{Help: "Max", Call: app.full, Icon: theme.ViewRestoreIcon(), Key1: fyne.KeySlash, Key2: fyne.KeyM},
		{Help: "Width", Call: app.width, Icon: theme.MoreHorizontalIcon(), Key1: fyne.KeyW},
		{Help: "Height", Call: app.height, Icon: theme.MoreVerticalIcon(), Key1: fyne.KeyH},
		{Help: "Refresh", Call: app.refresh, Icon: theme.ViewRefreshIcon(), Key1: fyne.KeyAsterisk, Key2: fyne.Key8},
		{Help: "Zoom In", Call: app.zoomIn, Icon: theme.ZoomInIcon(), Key1: fyne.KeyPlus, Key2: fyne.KeyEqual},
		{Help: "Zoom Out", Call: app.zoomOut, Icon: theme.ZoomOutIcon(), Key1: fyne.KeyMinus},
		{Help: "Close", Call: app.hideHelp, Key1: fyne.KeyEscape},
		{Help: "To Top", Call: app.top, Key1: fyne.KeyUp},
		{Help: "To Bottom", Call: app.bottom, Key1: fyne.KeyDown},
		{Help: "Spacer"},
		{Help: "Help", Call: app.showHelp, Icon: theme.HelpIcon(), Key1: fyne.KeyF1, Key2: fyne.KeyF10},
		{Help: "Quit", Call: fyneApp.Quit, Key1: fyne.KeyQ},
	}

	items := []widget.ToolbarItem{}

	for _, command := range commands {
		if command.Help == "Separator" {
			items = append(items, widget.NewToolbarSeparator())

			continue
		}

		if command.Help == "Spacer" {
			items = append(items, widget.NewToolbarSpacer())

			continue
		}

		if command.Icon == nil {
			continue
		}

		items = append(items, widget.NewToolbarAction(command.Icon, command.Call))
	}

	toolbar := widget.NewToolbar(items...)
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
	app.commands = commands
	app.help = app.createHelp()
	app.show = Debounced(app.Show, time.Millisecond*300)

	return app
}

func (p *App) init() {
	p.main.Canvas().SetOnTypedKey(func(ke *fyne.KeyEvent) {
		for _, command := range p.commands {
			if ke.Name == command.Key1 || ke.Name == command.Key2 {
				command.Call()
				log.Println(ke.Name)
			}
		}
	})
	p.SetMode(ModeFit)
}

func (p *App) Run(args []string) {
	defer p.cache.Close()
	// p.files.Load([]string{"doc/logo.png", "doc/maskable_icon.png"})
	// p.files.Load([]string{"doc/a.zip", "doc/logo.png", "doc/maskable_icon.png"})
	p.files.Load([]string{"doc"})
	// p.files.Load([]string{"doc/a.zip"})
	// p.files.Load(args)

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
	four := 4

	objects := []fyne.CanvasObject{
		canvas.NewText("Function", white),
		canvas.NewText("Key1", white),
		canvas.NewText("Key2", white),
		canvas.NewText("", white),
	}

	for _, command := range p.commands {
		if command.Key1 == "" {
			continue
		}

		objects = append(objects,
			canvas.NewText(command.Help, white),
			canvas.NewText(string(command.Key1), grey),
			canvas.NewText(string(command.Key2), grey),
			widget.NewIcon(command.Icon),
		)
	}

	objects = append(objects, canvas.NewText("", white), canvas.NewText("", white), canvas.NewText("", white),
		widget.NewHyperlink("xuender/comic", url))

	grid := container.New(layout.NewGridLayout(four), objects...)

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
