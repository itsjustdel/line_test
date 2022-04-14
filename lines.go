package main

import (
	"image/color"
	"math"
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
)

var currentWidth = 2
var currentSides = 4
var rotation = 0
var maxSides = 12
var lines = []*canvas.Line{}

func rotatePoint(point fyne.Position, rads float64) fyne.Position {

	//some maths to rotate a point round a zero origin
	//math library wants float64 number types
	px := float64(point.X)
	py := float64(point.Y)

	//rotate
	x := math.Cos(rads)*(px) - math.Sin(rads)*(py)
	y := math.Sin(rads)*(px) + math.Cos(rads)*(py)

	//and place in centre of screen
	x += 250
	y += 250

	//NewPos wants float32s
	newPos := fyne.NewPos(float32(x), float32(y))

	return newPos
}

func buttonContainer() *fyne.Container {

	//buttons
	sidesPlusButton := widget.NewButton("Sides +", func() {
		if currentSides < maxSides {
			currentSides++
		}
	})

	sidesMinusButton := widget.NewButton("Sides -", func() {
		if currentSides > 2 {
			currentSides--
		}
	})

	widthPlusButton := widget.NewButton("Width + ", func() {
		if currentWidth < 30 {
			currentWidth += 2
		}
	})

	widthMinusButton := widget.NewButton("Width - ", func() {
		if currentWidth > 2 {
			currentWidth -= 2
		}
	})

	rotatePlusButton := widget.NewButton("Rotate + ", func() {
		rotation++
	})

	rotateMinusButton := widget.NewButton("Rotate - ", func() {
		rotation--
	})

	//container to present in a vertical array
	buttonContainer := container.NewVBox(sidesMinusButton, sidesPlusButton, widthMinusButton, widthPlusButton, rotateMinusButton, rotatePlusButton)

	return buttonContainer
}

func shapeContainer() *fyne.Container {

	//create lines which we will then manipulate
	for i := 0; i < maxSides; i++ {

		line := canvas.NewLine(color.White)
		line.StrokeWidth = 5
		lines = append(lines, line)
	}

	//use polymorphism to create a slice of Canvas Object to set in the in the container
	canvasObjects := []fyne.CanvasObject{}
	for i := 0; i < len(lines); i++ {
		canvasObjects = append(canvasObjects, lines[i])
	}

	//set all lineObjects
	shapeContainer := container.NewWithoutLayout(canvasObjects...)

	return shapeContainer
}

func rotateShape(angle float64, armToRotate fyne.Position) {

	for i, line := range lines {

		if i >= currentSides {
			line.Hide()
			continue
		}

		line.Show()

		//adjust width
		line.StrokeWidth = float32(currentWidth)

		//step round the circle by how many sides we have chosen
		tempAngle := angle + (math.Pi*2)/float64(currentSides)*float64(i)
		line.Position1 = rotatePoint(armToRotate, tempAngle)

		//add another step to find the next point
		tempAngle += (math.Pi * 2) / float64(currentSides)
		line.Position2 = rotatePoint(armToRotate, tempAngle)
		line.Refresh()

	}
}

func main() {
	myApp := app.New()
	myWindow := myApp.NewWindow("Better Make Shapes")

	buttonContainer := buttonContainer()

	shapeContainer := shapeContainer()

	angle := float64(math.Pi / 4)
	armToRotate := fyne.NewPos(0, 100)
	go func() {
		for range time.Tick(time.Millisecond * 10) {

			rotateShape(angle, armToRotate)

			//add to counter/angle for next frame
			angle += float64(rotation) / 100
		}
	}()

	hBox := container.NewHBox(buttonContainer, shapeContainer)
	myWindow.SetContent(hBox)

	myWindow.Resize(fyne.NewSize(700, 500))
	myWindow.ShowAndRun()
}
