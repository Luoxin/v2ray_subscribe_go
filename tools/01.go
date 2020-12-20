package main

import (
	"fyne.io/fyne/app"
	"fyne.io/fyne/widget"
)

func main() {
	a := app.New()

	w := a.NewWindow("Hello")
	
	hello := widget.NewLabel("Hello Fyne!")
	w.SetContent(widget.NewVBox(
		hello,
	))
	
	

	w.ShowAndRun()
}