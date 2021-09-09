package tetris

import (
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/gotk3/gotk3/cairo"
	"github.com/gotk3/gotk3/gdk"
	"github.com/gotk3/gotk3/glib"
	"github.com/gotk3/gotk3/gtk"
)

const (
	SIGNAL_ACTIVATE        = "activate"
	SIGNAL_DRAW            = "draw"
	SIGNAL_KEY_PRESS_EVENT = "key-press-event"

	KEY_LEFT  uint = 65361
	KEY_UP    uint = 65362
	KEY_RIGHT uint = 65363
	KEY_DOWN  uint = 65364

	ACTION_QUIT    = "app.quit"
	ACTION_PAUSE   = "win.pause"
	ACTION_RESUME  = "win.resume"
	ACTION_NEWGAME = "win.start"

	ACTION_ROTATE = "win.rotate"
	ACTION_LEFT   = "win.left"
	ACTION_RIGHT  = "win.right"
	ACTION_DOWN   = "win.down"

	LABEL_PAUSE     = "Pause"
	LABEL_RESUME    = "Resume"
	LABEL_STARTGAME = "Start Game"

	LABEL_SCORE = "SCORE"

	UNIT_SIZE = 32
	SPAN_SIZE = UNIT_SIZE - 2
)

var (
	RGB_COLOR_GRAY  = Rgb{231 / 255.0, 231 / 255.0, 231 / 255.0}
	RGB_COLOR_BLUE  = Rgb{168 / 255.0, 202 / 255.0, 1}
	RGB_COLOR_GREEN = Rgb{132 / 255.0, 212 / 255.0, 129 / 255.0}

	centralDa *gtk.DrawingArea
	rightDa   *gtk.DrawingArea

	stateLabel *gtk.Label
	scoreValue *gtk.Label
	levelValue *gtk.Label
)

type Rgb [3]float64

func GUI() {
	const appID = "com.github.cloudecho.tetris"
	application, err := gtk.ApplicationNew(appID, glib.APPLICATION_FLAGS_NONE)
	if err != nil {
		log.Fatal("Could not create application:", err)
	}

	// Initialize game
	game := NewGame()
	go showGame(game)

	application.Connect(SIGNAL_ACTIVATE, func() {
		win := newWindow(application, game)

		aQuit := glib.SimpleActionNew("quit", nil)
		aQuit.Connect(SIGNAL_ACTIVATE, func() {
			application.Quit()
		})
		application.AddAction(aQuit)

		win.ShowAll()
	})

	os.Exit(application.Run(os.Args))
}

func showGame(g *Game) {
	for {
		select {
		case pos := <-g.chanPos:
			showCurrentShape(pos, g)
		case <-g.chanNexts:
			showNextShape(g)
		case state := <-g.chanState:
			switch state {
			case SATE_GAMEOVER:
				stateLabel.SetLabel("GAME OVER")
			case STATE_GAMING:
				stateLabel.SetLabel("")
			case STATE_PAUSED:
				stateLabel.SetLabel("PAUSED")
			case STATE_ZERO:
				// reset gui
				resetGui()
			}
		case level := <-g.chanLevel:
			levelValue.SetMarkup(markup("#000", UNIT_SIZE, strconv.Itoa(int(level))))
		case score := <-g.chanScore:
			scoreValue.SetMarkup(markup("#000", UNIT_SIZE, strconv.FormatUint(score, 10)))
		case area := <-g.chanRedraw:
			redrawArea(area, centralDa, g)
		case row := <-g.chanHiligh:
			drawHiligh(row, centralDa)
		}
	}
}

func resetGui() {
	fillBackgroud(centralDa, ROW, COL)
	centralDa.QueueDraw()
}

// Redraw area (top~otop rows)
func redrawArea(area Point, da *gtk.DrawingArea, g *Game) {
	m := &g.model
	da.Connect(SIGNAL_DRAW, func(da *gtk.DrawingArea, cr *cairo.Context) {
		for i := area.top; i <= area.otop; i++ { // top
			for j := 0; j < COL; j++ { // left
				rgb := rgb(m[i][j])
				cr.SetSourceRGB(rgb[0], rgb[1], rgb[2])
				cr.Rectangle(float64(j*UNIT_SIZE), float64(i*UNIT_SIZE), SPAN_SIZE, SPAN_SIZE)
				cr.Fill()
			}
		}
	})
	da.QueueDraw()
}

func drawHiligh(row int, da *gtk.DrawingArea) {
	for i := 0; i < 3; i++ {
		_drawHiligh(row, da, RGB_COLOR_GRAY)
		time.Sleep(time.Millisecond * 100)
		_drawHiligh(row, da, RGB_COLOR_GREEN)
		time.Sleep(time.Millisecond * 100)
	}
}

func _drawHiligh(row int, da *gtk.DrawingArea, color Rgb) {
	y := row * UNIT_SIZE
	da.Connect(SIGNAL_DRAW, func(da *gtk.DrawingArea, cr *cairo.Context) {
		cr.SetSourceRGB(color[0], color[1], color[2])
		for j := 0; j < COL; j++ { // left
			cr.Rectangle(float64(j*UNIT_SIZE), float64(y), SPAN_SIZE, SPAN_SIZE)
		}
		cr.Fill()
	})
	da.QueueDraw()
}

func rgb(v uint8) Rgb {
	if v > 0 {
		return RGB_COLOR_BLUE
	}
	return RGB_COLOR_GRAY
}

func showCurrentShape(pos Point, g *Game) {
	shape := g.currShape

	// Hide the old shape
	opos := Point{left: pos.oleft, top: pos.otop}
	drawShape(opos, shape, RGB_COLOR_GRAY, true, centralDa, g)

	// Show the current shape
	drawShape(pos, shape, RGB_COLOR_BLUE, false, centralDa, g)
}

func showNextShape(g *Game) {
	shape := g.nextShape

	// Hide the old shape
	pos := Point{left: 0, top: 0, otop: -2}
	drawShape(pos, shape, RGB_COLOR_GRAY, true, rightDa, g)

	// Show the current shape
	drawShape(pos, shape, RGB_COLOR_BLUE, false, rightDa, g)
}

func drawShape(pos Point, shape *Shape, rgb Rgb, erase bool, da *gtk.DrawingArea, g *Game) {
	if pos.top < 0 {
		return
	}

	// Show the current shape
	da.Connect(SIGNAL_DRAW, func(da *gtk.DrawingArea, cr *cairo.Context) {
		cr.SetSourceRGB(rgb[0], rgb[1], rgb[2])
		for i := 0; i < SHAPE_SIZE; i++ { // left
			for j := 0; j < SHAPE_SIZE; j++ { // top
				if !checkOutOfBounds(pos.left+i, pos.top+j) &&
					(erase && (pos.otop <= -2 || g.model[pos.top+j][pos.left+i] == 0) ||
						!erase && shape.data[j][i] > 0) {
					cr.Rectangle(
						float64(pos.left+i)*UNIT_SIZE,
						float64(pos.top+j)*UNIT_SIZE,
						SPAN_SIZE,
						SPAN_SIZE)
				}
			}
		}
		cr.Fill()
	})

	da.QueueDraw()
}

func newWindow(application *gtk.Application, g *Game) *gtk.ApplicationWindow {
	win, err := gtk.ApplicationWindowNew(application)
	if err != nil {
		log.Fatal("Unable to create window:", err)
	}

	win.SetTitle("TETRIS")
	initTitleBar(win, g)

	// Centrol & Right panels
	box, _ := gtk.BoxNew(gtk.ORIENTATION_HORIZONTAL, 10)
	initCentralPanel(box)
	initRightPanel(box)
	addMovingButtonActions(win, g)

	// Assemble the window
	win.Add(box)
	win.SetPosition(gtk.WIN_POS_MOUSE)
	win.SetDefaultSize(500, 600)
	return win
}

func initTitleBar(win *gtk.ApplicationWindow, g *Game) {
	// Create a header bar
	header, err := gtk.HeaderBarNew()
	if err != nil {
		log.Fatal("Could not create header bar:", err)
	}
	header.SetShowCloseButton(false)
	header.SetTitle("TETRIS")
	header.SetSubtitle("github.com/cloudecho/tetris")

	// Set up the menu model for the button
	menu := glib.MenuNew()
	if menu == nil {
		log.Fatal("Could not create menu (nil)")
	}

	// Actions with the prefix 'app' reference actions on the application
	// Actions with the prefix 'win' reference actions on the current window (specific to ApplicationWindow)
	// Other prefixes can be added to widgets via InsertActionGroup
	menu.Append(LABEL_STARTGAME, ACTION_NEWGAME)
	menu.Append("Quit", ACTION_QUIT)

	// Create a new menu button
	mbtn, err := gtk.MenuButtonNew()
	if err != nil {
		log.Fatal("Could not create menu button:", err)
	}

	mbtn.SetMenuModel(&menu.MenuModel)

	// Add the menu button to the header
	header.PackStart(mbtn)

	// Title buttons
	btnPause := btnPause()
	btnStart := btnStart()

	// Add title buttons to the end
	buttonBox, _ := gtk.ButtonBoxNew(gtk.ORIENTATION_HORIZONTAL)
	buttonBox.Add(btnStart)
	buttonBox.Add(btnPause)
	header.PackEnd(buttonBox)

	addTitleButtonActions(win, btnPause, g)
	win.SetTitlebar(header)
}

func btnPause() *gtk.Button {
	btn, _ := gtk.ButtonNew()
	btn.SetActionName(ACTION_PAUSE)
	btn.SetLabel(LABEL_PAUSE)
	return btn
}

func btnStart() *gtk.Button {
	btn, _ := gtk.ButtonNew()
	btn.SetActionName(ACTION_NEWGAME)
	btn.SetLabel("Start")
	btn.SetTooltipText(LABEL_STARTGAME)
	return btn
}

func addTitleButtonActions(win *gtk.ApplicationWindow, btnPause *gtk.Button, g *Game) {
	addActionTo(win, simpleActionName4Win(ACTION_NEWGAME), func() {
		g.start()
	})

	addActionTo(win, simpleActionName4Win(ACTION_PAUSE), func() {
		btnPause.SetLabel(LABEL_RESUME)
		btnPause.SetActionName(ACTION_RESUME)
		g.pause()
	})

	addActionTo(win, simpleActionName4Win(ACTION_RESUME), func() {
		btnPause.SetLabel(LABEL_PAUSE)
		btnPause.SetActionName(ACTION_PAUSE)
		g.resume()
	})
}

// Create an action in the win action group
func addActionTo(
	win *gtk.ApplicationWindow,
	actionName string,
	activateFunc func()) {
	a := glib.SimpleActionNew(actionName, nil)
	a.Connect(SIGNAL_ACTIVATE, activateFunc)
	win.AddAction(a)
}

func initCentralPanel(parent *gtk.Box) {
	da, _ := gtk.DrawingAreaNew()
	fillBackgroud(da, ROW, COL)
	da.SetSizeRequest(COL*UNIT_SIZE, (ROW+1)*UNIT_SIZE)
	centralDa = da

	parent.PackStart(da, true, true, 10)
}

func fillBackgroud(da *gtk.DrawingArea, row, col int) {
	da.Connect(SIGNAL_DRAW, func(da *gtk.DrawingArea, cr *cairo.Context) {
		cr.SetSourceRGB(RGB_COLOR_GRAY[0], RGB_COLOR_GRAY[1], RGB_COLOR_GRAY[2])
		for i := 0; i < row; i++ {
			for j := 0; j < col; j++ {
				cr.Rectangle(float64(j*UNIT_SIZE), float64(i*UNIT_SIZE), SPAN_SIZE, SPAN_SIZE)
			}
		}
		cr.Fill()
	})
}

func initRightPanel(parent *gtk.Box) {
	initValueLabels()
	btnRotate, btnLeft, btnRight, btnDown := initMovingButtons()

	da, _ := gtk.DrawingAreaNew()
	fillBackgroud(da, SHAPE_SIZE, SHAPE_SIZE)
	da.SetSizeRequest(SHAPE_SIZE*UNIT_SIZE, SHAPE_SIZE*UNIT_SIZE)
	da.SetMarginTop(UNIT_SIZE)
	rightDa = da

	stateLabel, _ = gtk.LabelNew("")
	scoreLabel, _ := gtk.LabelNew("")
	levelLabel, _ := gtk.LabelNew("")
	separator1, _ := gtk.LabelNew("")
	separator2, _ := gtk.LabelNew("")
	separator3, _ := gtk.LabelNew("")
	separator4, _ := gtk.LabelNew("")

	scoreLabel.SetMarkup(markup("#000", UNIT_SIZE, "SCORE"))
	scoreValue.SetMarkup(markup("#000", UNIT_SIZE, "0"))
	levelLabel.SetMarkup(markup("#000", UNIT_SIZE, "LEVEL"))
	levelValue.SetMarkup(markup("#000", UNIT_SIZE, "0"))
	separator1.SetMarkup(markup("#000", 4*UNIT_SIZE, " "))
	separator2.SetMarkup(markup("#000", UNIT_SIZE, " "))
	separator3.SetMarkup(markup("#000", UNIT_SIZE, " "))
	separator4.SetMarkup(markup("#000", UNIT_SIZE/2, " "))

	grid, _ := gtk.GridNew()
	grid.Attach(da, 0, 0, 3, 1)
	grid.Attach(separator1, 0, 1, 3, 1)
	grid.Attach(scoreLabel, 0, 2, 3, 1)
	grid.Attach(scoreValue, 0, 3, 3, 1)
	grid.Attach(separator2, 0, 4, 3, 1)
	grid.Attach(levelLabel, 0, 5, 3, 1)
	grid.Attach(levelValue, 0, 6, 3, 1)
	grid.Attach(separator3, 0, 7, 3, 1)
	grid.Attach(btnRotate, 1, 8, 1, 1)
	grid.Attach(btnLeft, 0, 9, 1, 1)
	grid.Attach(btnRight, 2, 9, 1, 1)
	grid.Attach(btnDown, 1, 10, 1, 1)
	grid.Attach(separator4, 0, 11, 3, 1)
	grid.Attach(stateLabel, 0, 12, 3, 1)

	parent.PackEnd(grid, true, true, 10)
}

func markup(color string, fontSize int, text string) string {
	return fmt.Sprintf(
		"<span foreground='%s' font='%d'>%s</span>",
		color, fontSize, text)
}

func initValueLabels() {
	scoreValue, _ = gtk.LabelNew("")
	levelValue, _ = gtk.LabelNew("")
}

func initMovingButtons() (*gtk.Button, *gtk.Button, *gtk.Button, *gtk.Button) {
	btnRotate, _ := gtk.ButtonNewWithLabel("^")
	btnLeft, _ := gtk.ButtonNewWithLabel("<")
	btnRight, _ := gtk.ButtonNewWithLabel(">")
	btnDown, _ := gtk.ButtonNewWithLabel("v")

	btnRotate.SetActionName(ACTION_ROTATE)
	btnLeft.SetActionName(ACTION_LEFT)
	btnRight.SetActionName(ACTION_RIGHT)
	btnDown.SetActionName(ACTION_DOWN)

	return btnRotate, btnLeft, btnRight, btnDown
}

func addMovingButtonActions(win *gtk.ApplicationWindow, g *Game) {
	keyMap := map[uint]func(){
		KEY_LEFT:  func() { g.moveLeft() },
		KEY_UP:    func() { g.rotate() },
		KEY_RIGHT: func() { g.moveRight() },
		KEY_DOWN:  func() { g.dropDown() },
	}

	chanKey := make(chan uint, 1)

	win.Connect(SIGNAL_KEY_PRESS_EVENT, func(win *gtk.ApplicationWindow, ev *gdk.Event) {
		// Discard if chanKey not empty
		if len(chanKey) > 0 {
			return
		}
		keyEvent := &gdk.EventKey{ev}
		keyVal := keyEvent.KeyVal()
		if _, found := keyMap[keyVal]; found {
			chanKey <- keyVal
		}
	})

	go func() {
		for {
			k := <-chanKey
			if action, found := keyMap[k]; found {
				action()
			}
			// To avoid too frequently key events
			time.Sleep(time.Millisecond * 200)
		}
	}()

	addActionTo(win, simpleActionName4Win(ACTION_ROTATE), keyMap[KEY_UP])
	addActionTo(win, simpleActionName4Win(ACTION_LEFT), keyMap[KEY_LEFT])
	addActionTo(win, simpleActionName4Win(ACTION_RIGHT), keyMap[KEY_RIGHT])
	addActionTo(win, simpleActionName4Win(ACTION_DOWN), keyMap[KEY_DOWN])
}

func simpleActionName4Win(fullname string) string {
	return strings.TrimPrefix(fullname, "win.")
}
