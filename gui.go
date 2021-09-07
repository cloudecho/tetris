package tetris

import (
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/gotk3/gotk3/glib"
	"github.com/gotk3/gotk3/gtk"
)

const (
	SIGNAL_ACTIVATE = "activate"

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

	DEFAULT_SPAN_COLOR = "#e7e7e7"
	SPAN_COLOR_BLUE    = "#a8caff"

	UNIT_SIZE = 32
)

var (
	mainSpans    [ROW][COL]*gtk.Label
	previewSpans [SHAPE_SIZE][SHAPE_SIZE]*gtk.Label

	scoreValue *gtk.Label
	levelValue *gtk.Label

	btnRotate *gtk.Button
	btnLeft   *gtk.Button
	btnRight  *gtk.Button
	btnDown   *gtk.Button

	game *Game
)

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
			log.Println(pos)
			showCurrentShape(pos)
		}
	}
}

func showCurrentShape(pos Point) {
	shape := game.currShape

	// Hide the old shape
	opos := Point{left: pos.oLeft, top: pos.oTop}
	showShape(opos, shape, DEFAULT_SPAN_COLOR, true)

	// Show the current shape
	showShape(pos, shape, SPAN_COLOR_BLUE, false)
}

func showShape(pos Point, shape *Shape, color string, full bool) {
	if pos.top < 0 || pos.left < 0 {
		return
	}

	// Show the current shape
	for i := 0; i < SHAPE_SIZE; i++ { // left
		for j := 0; j < SHAPE_SIZE; j++ { // top
			if (full || shape.data[j][i] > 0) &&
				!checkOutOfBound(pos.left+i, pos.top+j) {
				span := mainSpans[pos.top+j][pos.left+i]
				span.SetMarkup(markupSpan(color))
			}
		}
	}
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
	initCentralPanel(box, win)
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

func initCentralPanel(parent *gtk.Box, win *gtk.ApplicationWindow) {
	initMainSpans()
	grid, _ := gtk.GridNew()
	grid.SetMarginBottom(10)
	for i := 0; i < ROW; i++ {
		for j := 0; j < COL; j++ {
			grid.Attach(mainSpans[i][j], j, i, 1, 1)
		}
	}

	parent.PackStart(grid, true, true, 10)
}

func initRightPanel(parent *gtk.Box) {
	initPreviewSpans()
	initValueLabels()
	initMovingButtons()

	grid0, _ := gtk.GridNew()
	for i := 0; i < SHAPE_SIZE; i++ {
		for j := 0; j < SHAPE_SIZE; j++ {
			grid0.Attach(previewSpans[i][j], j+1, i, 1, 1)
		}
	}

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
	grid.Attach(grid0, 0, 0, 3, 1)
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

func markupSpan(color string) string {
	return fmt.Sprintf(
		"<span background='%s' foreground='%s' font='%d'>✿</span>",
		color, color, UNIT_SIZE)
}

func markup(color string, fontSize int, text string) string {
	return fmt.Sprintf(
		"<span foreground='%s' font='%d'>%s</span>",
		color, fontSize, text)
}

func initMainSpans() {
	for i := 0; i < ROW; i++ {
		for j := 0; j < COL; j++ {
			label, _ := gtk.LabelNew("")
			label.SetMarkup(markupSpan(DEFAULT_SPAN_COLOR))
			label.SetSizeRequest(UNIT_SIZE, UNIT_SIZE)
			mainSpans[i][j] = label
		}
	}
}

func initPreviewSpans() {
	for i := 0; i < SHAPE_SIZE; i++ {
		for j := 0; j < SHAPE_SIZE; j++ {
			label, _ := gtk.LabelNew("")
			label.SetMarkup(markupSpan(DEFAULT_SPAN_COLOR))
			label.SetSizeRequest(UNIT_SIZE, UNIT_SIZE)
			previewSpans[i][j] = label
		}
	}
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
	addActionTo(win, simpleActionName4Win(ACTION_ROTATE), func() {
		game.rotate()
	})

	addActionTo(win, simpleActionName4Win(ACTION_LEFT), func() {
		log.Println("TODO ", ACTION_LEFT)
	})

	addActionTo(win, simpleActionName4Win(ACTION_RIGHT), func() {
		log.Println("TODO ", ACTION_RIGHT)
	})

	addActionTo(win, simpleActionName4Win(ACTION_DOWN), func() {
		log.Println("TODO ", ACTION_DOWN)
	})
}

func simpleActionName4Win(fullname string) string {
	return strings.TrimPrefix(fullname, "win.")
}
