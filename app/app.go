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
	ModeActual Mode = iota
	ModeWindow
	ModeWidth
	ModeHeight
)

// nolint: gochecknoglobals
var _modes = map[Mode]string{
	ModeActual: "Actual",
	ModeWindow: "Window",
	ModeWidth:  "Width",
	ModeHeight: "Height",
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

	app.createCommands()
	app.createRadio()
	border := container.NewBorder(app.createToolbar(), app.radio, nil, nil, scroll)
	// nolint: gomnd
	main.Resize(fyne.NewSize(800, 600))
	main.SetContent(border)
	main.CenterOnScreen()

	app.border = border
	app.help = app.createHelp()
	// nolint: gomnd
	app.show = Debounced(app.Show, time.Millisecond*200)

	return app
}

func (p *App) createCommands() {
	p.commands = []Command{
		{Help: "First Page", Call: p.first, Icon: theme.MediaSkipPreviousIcon(), Key1: fyne.KeyHome},
		{Help: "Previous Page", Call: p.previous, Icon: theme.NavigateBackIcon(), Key1: fyne.KeyPageUp},
		{Help: "Next Page", Call: p.next, Icon: theme.NavigateNextIcon(), Key1: fyne.KeyPageDown, Key2: fyne.KeySpace},
		{Help: "Last Page", Call: p.last, Icon: theme.MediaSkipNextIcon(), Key1: fyne.KeyEnd},
		{Help: "Separator"},
		{Help: "Full screen", Call: p.fullScreen, Icon: theme.ViewFullScreenIcon(), Key1: fyne.KeyF11},
		{Help: "Actual Size", Call: p.modeActual, Icon: theme.ViewRefreshIcon(), Key1: fyne.KeyAsterisk, Key2: fyne.Key8},
		{Help: "Fit to window", Call: p.modeWindow, Icon: theme.ViewRestoreIcon(), Key1: fyne.KeySlash, Key2: fyne.KeyM},
		{Help: "Fit to width", Call: p.modeWidth, Icon: theme.MoreHorizontalIcon(), Key1: fyne.KeyW},
		{Help: "Fit to height", Call: p.modeHeight, Icon: theme.MoreVerticalIcon(), Key1: fyne.KeyH},
		{Help: "Zoom In", Call: p.zoomIn, Icon: theme.ZoomInIcon(), Key1: fyne.KeyPlus, Key2: fyne.KeyEqual},
		{Help: "Zoom Out", Call: p.zoomOut, Icon: theme.ZoomOutIcon(), Key1: fyne.KeyMinus},
		{Help: "Close", Call: p.close, Key1: fyne.KeyEscape},
		{Help: "To Top", Call: p.top, Key1: fyne.KeyUp},
		{Help: "To Bottom", Call: p.bottom, Key1: fyne.KeyDown},
		{Help: "Spacer"},
		{Help: "This Help", Call: p.showHelp, Icon: theme.HelpIcon(), Key1: fyne.KeyF1, Key2: fyne.KeyF10},
		{Help: "Quit", Call: p.app.Quit, Key1: fyne.KeyQ},
	}
}

func (p *App) createRadio() {
	options := map[string]func(){
		"Actual": p.modeActual,
		"Window": p.modeWindow,
		"Width":  p.modeWidth,
		"Height": p.modeHeight,
	}
	p.radio = widget.NewRadioGroup([]string{"Actual", "Window", "Width", "Height"}, func(option string) {
		if call, has := options[option]; has {
			call()
		}
	})
	p.radio.Horizontal = true
}

func (p *App) createToolbar() *widget.Toolbar {
	items := []widget.ToolbarItem{}

	for _, command := range p.commands {
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

	return widget.NewToolbar(items...)
}

func (p *App) Run(args []string) {
	defer p.cache.Close()
	// p.files.Load([]string{"doc/logo.png", "doc/maskable_icon.png"})
	// p.files.Load([]string{"doc/a.zip", "doc/logo.png", "doc/maskable_icon.png"})
	p.files.Load([]string{"doc"})
	// p.files.Load([]string{"doc/a.zip"})
	// p.files.Load(args)

	p.setKey()
	p.Show()
	p.main.Show()
	p.app.Run()
}

func (p *App) setKey() {
	p.main.Canvas().SetOnTypedKey(func(event *fyne.KeyEvent) {
		if event.Name == "" {
			return
		}

		log.Println(event.Name)

		for _, command := range p.commands {
			if event.Name == command.Key1 || event.Name == command.Key2 {
				command.Call()
			}
		}
	})
	p.SetMode(ModeActual)
}

func (p *App) Show() {
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
	if p.img.Image == nil {
		p.img.SetMinSize(p.main.Canvas().Size())
	} else {
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
		canvas.NewText("", white),
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
	case ModeWindow:
		p.modeWindow()
	case ModeWidth:
		p.modeWidth()
	case ModeHeight:
		p.modeHeight()
	case ModeActual:
	default:
	}
}

func (p *App) showHelp() {
	p.help.Show()
}

func (p *App) close() {
	if p.main.FullScreen() {
		p.fullScreen()
	}

	p.help.Hide()
}

func (p *App) top() {
	p.scroll.Scrolled(&fyne.ScrollEvent{Scrolled: fyne.Delta{DX: 0, DY: p.scroll.Offset.Y}})
}

func (p *App) bottom() {
	p.scroll.ScrollToBottom()
}

func (p *App) next() {
	p.files.Down()
	p.show()
}

func (p *App) first() {
	p.files.Start()
	p.show()
}

func (p *App) last() {
	p.files.End()
	p.show()
}

func (p *App) previous() {
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

func (p *App) modeActual() {
	p.SetMode(ModeActual)

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
	p.radio.SetSelected(_modes[p.mode])
}

func (p *App) modeWindow() {
	p.SetMode(ModeWindow)
	size := ToSize(p.img.Size(), p.main.Canvas().Size())
	p.img.Resize(size)
	p.img.SetMinSize(size)
}

func (p *App) modeWidth() {
	p.SetMode(ModeWidth)
	width := p.main.Canvas().Size().Width
	height := width / p.img.Size().Width * p.img.Size().Height
	size := fyne.NewSize(width, height)
	p.img.Resize(size)
	p.img.SetMinSize(size)
}

func (p *App) modeHeight() {
	p.SetMode(ModeHeight)
	height := p.main.Canvas().Size().Height
	width := height / p.img.Size().Height * p.img.Size().Width
	size := fyne.NewSize(width, height)
	p.img.Resize(size)
	p.img.SetMinSize(size)
}
