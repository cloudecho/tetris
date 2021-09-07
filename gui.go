package tetris

import (
	"fmt"
	"log"
	"os"
	"strings"

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
	ACTION_NEWGAME = "win.newgame"

	ACTION_ROTATE = "win.rotate"
	ACTION_LEFT   = "win.left"
	ACTION_RIGHT  = "win.right"
	ACTION_DOWN   = "win.down"

	LABEL_PAUSE   = "Pause"
	LABEL_RESUME  = "Resume"
	LABEL_NEWGAME = "New Game"

	UNIT_SIZE = 32
	SPAN_SIZE = UNIT_SIZE - 2
)

var (
	RGB_COLOR_GRAY = Rgb{231 / 255.0, 231 / 255.0, 231 / 255.0}
	RGB_COLOR_BLUE = Rgb{168 / 255.0, 202 / 255.0, 1}

	centralDa *gtk.DrawingArea
	rightDa   *gtk.DrawingArea

	scoreValue *gtk.Label
	levelValue *gtk.Label

	btnRotate *gtk.Button
	btnLeft   *gtk.Button
	btnRight  *gtk.Button
	btnDown   *gtk.Button

	game *Game
)

type Rgb [3]float64

func GUI() {
	const appID = "com.github.cloudecho.tetris"
	application, err := gtk.ApplicationNew(appID, glib.APPLICATION_FLAGS_NONE)
	if err != nil {
		log.Fatal("Could not create application:", err)
	}

	application.Connect(SIGNAL_ACTIVATE, func() {
		win := newWindow(application)

		aQuit := glib.SimpleActionNew("quit", nil)
		aQuit.Connect(SIGNAL_ACTIVATE, func() {
			application.Quit()
		})
		application.AddAction(aQuit)

		win.ShowAll()
	})

	// Initialize game
	game = NewGame()
	go showGame(game)

	os.Exit(application.Run(os.Args))
}

func showGame(g *Game) {
	// TODO showGame
	for {
		select {
		case pos := <-g.p:
			log.Println(pos) // debug
			showCurrentShape(pos)
		case <-g.n:
			showNextShape()
		}
	}
}

func showCurrentShape(pos Point) {
	shape := game.currShape

	// Hide the old shape
	opos := Point{left: pos.oleft, top: pos.otop}
	showShape(opos, shape, RGB_COLOR_GRAY, true, centralDa)

	// Show the current shape
	showShape(pos, shape, RGB_COLOR_BLUE, false, centralDa)
}

func showNextShape() {
	shape := game.nextShape

	// Hide the old shape
	pos := Point{left: 0, top: 0}
	showShape(pos, shape, RGB_COLOR_GRAY, true, rightDa)

	// Show the current shape
	showShape(pos, shape, RGB_COLOR_BLUE, false, rightDa)
}

func showShape(pos Point, shape *Shape, rgb Rgb, full bool, da *gtk.DrawingArea) {
	if pos.top < 0 {
		return
	}

	// Show the current shape
	da.Connect(SIGNAL_DRAW, func(da *gtk.DrawingArea, cr *cairo.Context) {
		cr.SetSourceRGB(rgb[0], rgb[1], rgb[2])
		for i := 0; i < SHAPE_SIZE; i++ { // left
			for j := 0; j < SHAPE_SIZE; j++ { // top
				if (full || shape.data[j][i] > 0) &&
					!checkOutOfBound(pos.left+i, pos.top+j) {
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

func newWindow(application *gtk.Application) *gtk.ApplicationWindow {
	win, err := gtk.ApplicationWindowNew(application)
	if err != nil {
		log.Fatal("Unable to create window:", err)
	}

	win.SetTitle("TETRIS")
	initTitleBar(win)

	// Centrol & Right panels
	box, _ := gtk.BoxNew(gtk.ORIENTATION_HORIZONTAL, 10)
	initCentralPanel(box)
	initRightPanel(box)
	addMovingButtonActions(win)

	// Assemble the window
	win.Add(box)
	win.SetPosition(gtk.WIN_POS_MOUSE)
	win.SetDefaultSize(500, 600)
	return win
}

func initTitleBar(win *gtk.ApplicationWindow) {
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
	menu.Append(LABEL_NEWGAME, ACTION_NEWGAME)
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
	btnNew := btnNewGame()

	// Add title buttons to the end
	buttonBox, _ := gtk.ButtonBoxNew(gtk.ORIENTATION_HORIZONTAL)
	buttonBox.Add(btnNew)
	buttonBox.Add(btnPause)
	header.PackEnd(buttonBox)

	addTitleButtonActions(win, btnPause)
	win.SetTitlebar(header)
}

func btnPause() *gtk.Button {
	btn, _ := gtk.ButtonNew()
	btn.SetActionName(ACTION_PAUSE)
	btn.SetLabel(LABEL_PAUSE)
	return btn
}

func btnNewGame() *gtk.Button {
	btn, _ := gtk.ButtonNew()
	btn.SetActionName(ACTION_NEWGAME)
	btn.SetLabel("New")
	btn.SetTooltipText(LABEL_NEWGAME)
	return btn
}

func addTitleButtonActions(win *gtk.ApplicationWindow, btnPause *gtk.Button) {
	addActionTo(win, simpleActionName4Win(ACTION_NEWGAME), func() {
		log.Println("Start to game.")
		game.start()
	})

	addActionTo(win, simpleActionName4Win(ACTION_PAUSE), func() {
		btnPause.SetLabel(LABEL_RESUME)
		btnPause.SetActionName(ACTION_RESUME)
		log.Println("TODO PAUSE!")
	})

	addActionTo(win, simpleActionName4Win(ACTION_RESUME), func() {
		btnPause.SetLabel(LABEL_PAUSE)
		btnPause.SetActionName(ACTION_PAUSE)
		log.Println("TODO RESUME!")
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
	da.Connect(SIGNAL_DRAW, func(da *gtk.DrawingArea, cr *cairo.Context) {
		cr.SetSourceRGB(RGB_COLOR_GRAY[0], RGB_COLOR_GRAY[1], RGB_COLOR_GRAY[2])
		for i := 0; i < ROW; i++ {
			for j := 0; j < COL; j++ {
				cr.Rectangle(float64(j*UNIT_SIZE), float64(i*UNIT_SIZE), SPAN_SIZE, SPAN_SIZE)
			}
		}
		cr.Fill()
	})
	da.SetSizeRequest(COL*UNIT_SIZE, (ROW+1)*UNIT_SIZE)
	centralDa = da

	parent.PackStart(da, true, true, 10)
}

func initRightPanel(parent *gtk.Box) {
	initValueLabels()
	initMovingButtons()

	da, _ := gtk.DrawingAreaNew()
	da.Connect(SIGNAL_DRAW, func(da *gtk.DrawingArea, cr *cairo.Context) {
		cr.SetSourceRGB(RGB_COLOR_GRAY[0], RGB_COLOR_GRAY[1], RGB_COLOR_GRAY[2])
		for i := 0; i < SHAPE_SIZE; i++ {
			for j := 0; j < SHAPE_SIZE; j++ {
				cr.Rectangle(float64(j*UNIT_SIZE), float64(i*UNIT_SIZE), SPAN_SIZE, SPAN_SIZE)
			}
		}
		cr.Fill()
	})
	da.SetSizeRequest(SHAPE_SIZE*UNIT_SIZE, SHAPE_SIZE*UNIT_SIZE)
	da.SetMarginTop(UNIT_SIZE)
	rightDa = da

	scoreLabel, _ := gtk.LabelNew("")
	levelLabel, _ := gtk.LabelNew("")
	separator1, _ := gtk.LabelNew("")
	separator2, _ := gtk.LabelNew("")
	separator3, _ := gtk.LabelNew("")

	scoreLabel.SetMarkup(markup("#000", UNIT_SIZE, "SCORE"))
	scoreValue.SetMarkup(markup("#000", UNIT_SIZE, "0"))
	levelLabel.SetMarkup(markup("#000", UNIT_SIZE, "LEVEL"))
	levelValue.SetMarkup(markup("#000", UNIT_SIZE, "0"))
	separator1.SetMarkup(markup("#000", 4*UNIT_SIZE, " "))
	separator2.SetMarkup(markup("#000", UNIT_SIZE, " "))
	separator3.SetMarkup(markup("#000", UNIT_SIZE, " "))

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

func initMovingButtons() {
	btnRotate, _ = gtk.ButtonNewWithLabel("^")
	btnLeft, _ = gtk.ButtonNewWithLabel("<")
	btnRight, _ = gtk.ButtonNewWithLabel(">")
	btnDown, _ = gtk.ButtonNewWithLabel("v")

	btnRotate.SetActionName(ACTION_ROTATE)
	btnLeft.SetActionName(ACTION_LEFT)
	btnRight.SetActionName(ACTION_RIGHT)
	btnDown.SetActionName(ACTION_DOWN)
}

func addMovingButtonActions(win *gtk.ApplicationWindow) {
	keyMap := map[uint]func(){
		KEY_LEFT:  func() { game.moveLeft() },
		KEY_UP:    func() { game.rotate() },
		KEY_RIGHT: func() { game.moveRight() },
		KEY_DOWN:  func() { game.dropDown() },
	}

	win.Connect(SIGNAL_KEY_PRESS_EVENT, func(win *gtk.ApplicationWindow, ev *gdk.Event) {
		keyEvent := &gdk.EventKey{ev}
		if action, found := keyMap[keyEvent.KeyVal()]; found {
			action()
		}
	})

	addActionTo(win, simpleActionName4Win(ACTION_ROTATE), keyMap[KEY_UP])
	addActionTo(win, simpleActionName4Win(ACTION_LEFT), keyMap[KEY_LEFT])
	addActionTo(win, simpleActionName4Win(ACTION_RIGHT), keyMap[KEY_RIGHT])
	addActionTo(win, simpleActionName4Win(ACTION_DOWN), keyMap[KEY_DOWN])
}

func simpleActionName4Win(fullname string) string {
	return strings.TrimPrefix(fullname, "win.")
}
