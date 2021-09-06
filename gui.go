package tetris

import (
	"fmt"
	"log"
	"os"

	"github.com/gotk3/gotk3/glib"
	"github.com/gotk3/gotk3/gtk"
)

const (
	ACTION_PAUSE   = "win.pause"
	ACTION_RESUME  = "win.resume"
	ACTION_NEWGAME = "win.newgame"

	LABEL_PAUSE   = "Pause"
	LABEL_RESUME  = "Resume"
	LABEL_NEWGAME = "New Game"

	DEFAULT_SPAN_COLOR = "#e7e7e7"

	UNIT_SIZE = 32
)
 
var (
	mainSpans [ROW][COL]*gtk.Label
	previewSpans  [PRE_ROW][PRE_ROW]*gtk.Label
)

func GUI() {
	const appID = "com.github.cloudecho.tetris"
	application, err := gtk.ApplicationNew(appID, glib.APPLICATION_FLAGS_NONE)
	if err != nil {
		log.Fatal("Could not create application:", err)
	}

	application.Connect("activate", func() {
		win := newWindow(application)

		aQuit := glib.SimpleActionNew("quit", nil)
		aQuit.Connect("activate", func() {
			application.Quit()
		})
		application.AddAction(aQuit)

		win.ShowAll()
	})

	os.Exit(application.Run(os.Args))
}

func newWindow(application *gtk.Application) *gtk.ApplicationWindow {
	win, err := gtk.ApplicationWindowNew(application)
	if err != nil {
		log.Fatal("Unable to create window:", err)
	}

	win.SetTitle("TETRIS")

	// Label text in the window
	lbl, err := gtk.LabelNew("Use the menu button to test the actions")
	if err != nil {
		log.Fatal("Could not create label:", err)
	}

	// Create a header bar
	header, err := gtk.HeaderBarNew()
	if err != nil {
		log.Fatal("Could not create header bar:", err)
	}
	header.SetShowCloseButton(false)
	header.SetTitle("TETRIS")
	header.SetSubtitle("github.com/cloudecho/tetris")

	// Create a new menu button
	mbtn, err := gtk.MenuButtonNew()
	if err != nil {
		log.Fatal("Could not create menu button:", err)
	}

	// Set up the menu model for the button
	menu := glib.MenuNew()
	if menu == nil {
		log.Fatal("Could not create menu (nil)")
	}

	// Actions with the prefix 'app' reference actions on the application
	// Actions with the prefix 'win' reference actions on the current window (specific to ApplicationWindow)
	// Other prefixes can be added to widgets via InsertActionGroup
	menu.Append(LABEL_NEWGAME, ACTION_NEWGAME)
	menu.Append("Quit", "app.quit")

	// Custom buttons
	btnPause := btnPause()
	btnNew := btnNewGame()

	// Create an action in the custom action group
	addActionTo(win, "newgame", func() {
		lbl.SetLabel("NEW GAME!")
	})

	addActionTo(win, "pause", func() {
		btnPause.SetLabel(LABEL_RESUME)
		btnPause.SetActionName(ACTION_RESUME)
		lbl.SetLabel("PAUSE!")
	})

	addActionTo(win, "resume", func() {
		btnPause.SetLabel(LABEL_PAUSE)
		btnPause.SetActionName(ACTION_PAUSE)
		lbl.SetLabel("RESUME!")
	})

	mbtn.SetMenuModel(&menu.MenuModel)

	// add the menu button to the header
	header.PackStart(mbtn)

	// Add custom buttons to the end
	buttonBox, err := gtk.ButtonBoxNew(gtk.ORIENTATION_HORIZONTAL)
	if err != nil {
		log.Fatal("Could not create button box:", err)
	}
	buttonBox.Add(btnNew)
	buttonBox.Add(btnPause)
	header.PackEnd(buttonBox)

	// Assemble the window
	box, _ := gtk.BoxNew(gtk.ORIENTATION_HORIZONTAL, 10)
	createCenterPanel(box, win)
	createRightPanel(box)
	win.Add(box)
	win.SetTitlebar(header)
	win.SetPosition(gtk.WIN_POS_MOUSE)
	win.SetDefaultSize(500, 600)
	return win
}

func btnPause() *gtk.Button {
	btn, err := gtk.ButtonNew()
	if err != nil {
		log.Fatal("Could not create 'Pause' button:", err)
	}

	btn.SetActionName(ACTION_PAUSE)
	btn.SetLabel(LABEL_PAUSE)
	return btn
}

func btnNewGame() *gtk.Button {
	btn, err := gtk.ButtonNew()
	if err != nil {
		log.Fatal("Could not create 'New Game' button:", err)
	}

	btn.SetActionName(ACTION_NEWGAME)
	btn.SetLabel("New")
	btn.SetTooltipText(LABEL_NEWGAME)
	return btn
}

// Create an action in the custom action group
func addActionTo(
	win *gtk.ApplicationWindow,
	actionName string,
	activateFunc func()) {
	a := glib.SimpleActionNew(actionName, nil)
	a.Connect("activate", activateFunc)
	win.AddAction(a)
}

func createCenterPanel(parent *gtk.Box, win *gtk.ApplicationWindow) {
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

func createRightPanel(parent *gtk.Box) {
	initPreviewSpans()
	grid0, _ := gtk.GridNew()
	for i := 0; i < PRE_ROW; i++ {
		for j := 0; j < PRE_ROW; j++ {
			grid0.Attach(previewSpans[i][j], j+1, i, 1, 1)
		}
	}

	scoreLabel, _ := gtk.LabelNew("")
	scoreValue, _ := gtk.LabelNew("")
	levelLabel, _ := gtk.LabelNew("")
	levelValue, _ := gtk.LabelNew("")
	separator1, _ := gtk.LabelNew("")
	separator2, _ := gtk.LabelNew("")
	separator3, _ := gtk.LabelNew("")
	btnRotate, _ := gtk.ButtonNewWithLabel("^")
	btnLeft, _ := gtk.ButtonNewWithLabel("<")
	btnRight, _ := gtk.ButtonNewWithLabel(">")
	btnDown, _ := gtk.ButtonNewWithLabel("v")

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

func span(color string) string {
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
			label.SetMarkup(span(DEFAULT_SPAN_COLOR))
			label.SetSizeRequest(UNIT_SIZE, UNIT_SIZE)
			mainSpans[i][j] = label
		}
	}
}

func initPreviewSpans() {
	for i := 0; i < PRE_ROW; i++ {
		for j := 0; j < PRE_ROW; j++ {
			label, _ := gtk.LabelNew("")
			label.SetMarkup(span(DEFAULT_SPAN_COLOR))
			label.SetSizeRequest(UNIT_SIZE, UNIT_SIZE)
			previewSpans[i][j] = label
		}
	}
}
