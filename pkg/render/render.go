package render

import (
	"fmt"
	"time"
	"os"
	"image"
	"image/png"

	"github.com/faiface/pixel"
	"github.com/faiface/pixel/text"
	"github.com/faiface/pixel/pixelgl"
	"golang.org/x/image/colornames"
	"golang.org/x/image/font/basicfont"
)

type vec struct {
	x int
	y int
}

func iVec(v pixel.Vec) vec {
	return vec{int(v.X), int(v.Y)}
}

func fVec(v vec) pixel.Vec {
	return pixel.Vec{float64(v.x), float64(v.y)}
}

// GUI 
type GUI struct {
	atlas *text.Atlas
	brush *text.Text
	frameNr *text.Text
	sceneName *text.Text
	playbackFPS *text.Text
	brushBatch *pixel.Batch
}

// Canvas 
type Canvas struct {
	Win *pixelgl.Window
	
	title string
	width float64
	height float64

	// graphical user interface
	gui *GUI

	// set FPS
	FPS <-chan time.Time
	playbackFPS int

	// batch/sprite attributes
	batch *pixel.Batch
	prevBatch pixel.Batch
	batches []pixel.Batch
	brush *pixel.Sprite
	brushBuffer map[pixel.Vec]float64
	
	// painting/polling/framebuffer attributes
	frames [][]uint8
	frameNr int
	decay []uint8

	// canvas attributes
	erasing bool

	// brush attributes
	brushSize float64

}

// NewCanvas prepares a new Canvas
func NewCanvas(width float64, height float64, brushFile []byte, fontFile []byte) *Canvas {
	cfg := pixelgl.WindowConfig {
		Title:  "anim8",
		Bounds: pixel.R(0, 0, width, height),
		VSync:  false,
	}

	win, err := pixelgl.NewWindow(cfg)
	if err != nil {
		panic(err)
	}

	win.Canvas().SetSmooth(true)
	win.SetCursorVisible(false)

	// brush spritesheet
	spritesheet, err := loadPicture(brushFile)
	if err != nil {
		panic(err)
	}

	// gui
	face, err := loadTTF(fontFile, 52)
	if err != nil {
		panic(err)
	}

	screenNameAtlas := text.NewAtlas(face, text.ASCII)
	textAtlas := text.NewAtlas(basicfont.Face7x13, text.ASCII)

	gui := &GUI {
		textAtlas,
		text.New(pixel.V(width - 250, height - 30), textAtlas),
		text.New(pixel.V(width/2 - 50, 20), textAtlas),
		text.New(pixel.V(width/2 - 350, height/2), screenNameAtlas),
		text.New(pixel.V(30, height - 30), textAtlas),
		pixel.NewBatch(&pixel.TrianglesData{}, spritesheet),
	}


	batch := pixel.NewBatch(&pixel.TrianglesData{}, spritesheet)
	brush := pixel.NewSprite(spritesheet, spritesheet.Bounds())

	canvas := Canvas {
		win,
		cfg.Title,
		width,
		height,
		gui,
		time.Tick(time.Second / 120),
		15,
		batch,
		*batch,
		[]pixel.Batch{},
		brush,
		make(map[pixel.Vec]float64),
		[][]uint8{},
		1,
		nil,
		false,
		10,
	}

	canvas.gui.brush.Color = colornames.Red
	canvas.gui.frameNr.Color = colornames.Red
	canvas.gui.sceneName.Color = colornames.Red
	canvas.gui.playbackFPS.Color = colornames.Red

	return &canvas
}

// Paint draws or erases at the mouseposition
func (canvas *Canvas) Paint() {

	now := canvas.Win.MousePosition()
	prev := canvas.Win.MousePreviousPosition()

	// first draw as usual
	canvas.brush.Draw(canvas.batch, pixel.IM.Scaled(pixel.ZV, canvas.brushSize/20).Moved(now))

	// now if the mouse was moved, interpolate brush strokes inbetween

	// delta
	d := now.Sub(prev)

	// how many filling points to consider, depending on mouse delta
	points := d.Len()/(4*canvas.brushSize/20)
	if points > 16.0 {
		points = 16.0
	}
	
	// scaled delta
	delta := d.Scaled(1/points)

	// don't malloc a pixel.V every time
	paintPos := prev

	for i := float64(0); i < points; i = i+1 {
		paintPos = paintPos.Add(delta)
		canvas.brush.Draw(canvas.batch, pixel.IM.Scaled(pixel.ZV, canvas.brushSize/20).Moved(paintPos))
	}
}

// BrustType returns the string corresponding to the current brush type
func (canvas *Canvas) BrushType() string {
	if canvas.erasing {
		return "Erasor"
	}

	return "Default Brush"
}

// Clear canvas by using the decaying previous frame
func (canvas *Canvas) Clear() {
	if canvas.decay == nil {
		canvas.Win.Clear(colornames.Black)
	} else {
		canv := canvas.Win.Canvas()
		canv.SetPixels(canvas.decay)
	}
}


// Dump saves the animation as a set of PNGs using `sceneName` as the naming prefix
func (canvas *Canvas) Dump(sceneName string) {
	if _, err := os.Stat(sceneName); os.IsNotExist(err) {
		os.Mkdir(sceneName, 0700)
	}
	
	for i := 0; i < len(canvas.frames); i++ {
		img := image.NewRGBA(image.Rectangle{image.Point{0, 0}, image.Point{int(canvas.width), int(canvas.height)}})
		img.Pix = canvas.frames[i]

		file, err := os.Create(fmt.Sprintf(sceneName + "/%s%06d.png", sceneName, i))
		if err != nil {
			panic(err)
		}

		if err := png.Encode(file, img); err != nil {
			file.Close()
			panic(err)
		}
	}
}

// Poll user input
func (canvas *Canvas) Poll() {
	// paint at mouseclick
	if canvas.Win.Pressed(pixelgl.MouseButtonLeft) {
		canvas.Paint()	
	}
	
	// go into erasing mode at keypress E
	if canvas.Win.JustPressed(pixelgl.KeyE) {
		if canvas.erasing {
			canvas.erasing = false
			canvas.batch.SetColorMask(colornames.White)
		} else {
			canvas.erasing = true
			canvas.batch.SetColorMask(colornames.Black)
		}
	}

	// save canvas to animation buffer at keypress SPACE
	if canvas.Win.JustPressed(pixelgl.KeySpace) {
		// cache batch incase user wants to reuse the previous sketch
		canvas.prevBatch = *canvas.batch

		// clear screen except for canvas
		canvas.Win.Clear(colornames.Black)
		canvas.batch.Draw(canvas.Win)
		canvas.Win.Update()

		// now get canvas pixels
		canv := canvas.Win.Canvas()
		pixels := canv.Pixels()
		canvas.frames = append(canvas.frames, pixels)
		canvas.frameNr++
		
		// clear the target
		canvas.batch.Clear()

		// as an aid for drawing, indicate the previous frame
		decay := canv.Pixels()
		for i := 0; i < len(decay); i++ {
			decay[i] = uint8(float64(decay[i]) * 0.3)
		}

		canvas.decay = decay
	}

	// play animation at keypress P
	if canvas.Win.JustPressed(pixelgl.KeyP) {
		
		canv := canvas.Win.Canvas()
		
		// show animation at 15 FPS
		fps15 := time.Tick(time.Second/time.Duration(canvas.playbackFPS))
		for i := 0; i < len(canvas.frames); i++ {
			canv.SetPixels(canvas.frames[i])
			canvas.Win.Update()
			<-fps15
		}
	}

	if canvas.Win.JustPressed(pixelgl.KeyL) {

		canv := canvas.Win.Canvas()
		skipped := false

		for {
			// show animation at 15 FPS
			fps15 := time.Tick(time.Second/time.Duration(canvas.playbackFPS))
			for i := 0; i < len(canvas.frames); i++ {
				canv.SetPixels(canvas.frames[i])
				canvas.Win.Update()
				if canvas.Win.JustPressed(pixelgl.KeyL) {
					skipped = true
					break
				}	
				if canvas.Win.Pressed(pixelgl.KeyUp) {
					canvas.playbackFPS = canvas.playbackFPS + 1
				}
				if canvas.Win.Pressed(pixelgl.KeyDown) {
					canvas.playbackFPS = canvas.playbackFPS - 1
					if canvas.playbackFPS < 5 {
						canvas.playbackFPS = 5
					}
				}
				<-fps15
			}
			if skipped || canvas.Win.JustPressed(pixelgl.KeyL) {
				break
			}
		}		
	}

	// load previous batch at keypress C
	if canvas.Win.JustPressed(pixelgl.KeyC) {
		canvas.decay = nil
		canvas.batch = &canvas.prevBatch
		// TODO
		fmt.Println(*canvas.batch)
	}

	// move canvas to certain directon at keypresses up, down, left, right
	if canvas.Win.JustPressed(pixelgl.KeyUp) {
		// TODO
		canvas.batch.SetMatrix(pixel.IM.Moved(pixel.V(0, 1)))
	}

	if canvas.Win.JustPressed(pixelgl.KeyDown) {
		// TODO
		canvas.batch.SetMatrix(pixel.IM.Moved(pixel.V(0, -1)))	
	}

	if canvas.Win.JustPressed(pixelgl.KeyLeft) {
		// TODO
		canvas.batch.SetMatrix(pixel.IM.Moved(pixel.V(-1, 0)))	
	}

	if canvas.Win.JustPressed(pixelgl.KeyRight) {
		// TODO
		canvas.batch.SetMatrix(pixel.IM.Moved(pixel.V(1, 0)))
	}

	// reset animation at keypress R
	if canvas.Win.JustPressed(pixelgl.KeyS) {
		canvas.frames = [][]uint8{}
		canvas.frameNr = 1
		canvas.decay = nil
	}

	// delete current frame at keypress D
	if canvas.Win.JustPressed(pixelgl.KeyD) {
		canvas.batch.Clear()
	}

	// dump animation at keypress ENTER
	if canvas.Win.JustPressed(pixelgl.KeyEnter) {
		// remember previous frame state
		canv := canvas.Win.Canvas()
		pixels := canv.Pixels()
		sceneName := ""
		for {
			sceneName = sceneName + canvas.Win.Typed()
			canvas.gui.sceneName.WriteString(canvas.Win.Typed())
			canv.SetPixels(pixels)
			canvas.gui.sceneName.Draw(canvas.Win, pixel.IM)
			canvas.Win.Update()
			if canvas.Win.JustPressed(pixelgl.KeyEnter) {
				break
			}
		}
		canvas.Dump(sceneName)
	}

	if canvas.Win.JustPressed(pixelgl.KeyEscape) {
		canvas.Win.Destroy()
	}

	// adjust brush size at mousescroll
	scroll := canvas.Win.MouseScroll()
	canvas.brushSize = canvas.brushSize - scroll.X + scroll.Y 
	if canvas.brushSize < 1 {
		canvas.brushSize = 1
	}

	// update GUI
	fmt.Fprintf(canvas.gui.brush, "Brush\nSize\t%.0f\nType\t%s", canvas.brushSize, canvas.BrushType())
	fmt.Fprintf(canvas.gui.frameNr, "Frame Nr. %d", canvas.frameNr)
	fmt.Fprintf(canvas.gui.playbackFPS, "Playback-FPS\t%d", canvas.playbackFPS)
}

// Draw renders the canvas onto the window
func (canvas *Canvas) Draw() {
	canvas.Clear()

	canvas.batch.Draw(canvas.Win)

	// draw GUI
	canvas.brush.Draw(canvas.gui.brushBatch, pixel.IM.Scaled(pixel.ZV, canvas.brushSize/20).Moved(canvas.Win.MousePosition()))
	canvas.gui.brushBatch.Draw(canvas.Win)
	canvas.gui.brushBatch.Clear()

	canvas.gui.brush.Draw(canvas.Win, pixel.IM.Scaled(canvas.gui.brush.Orig, 1.4))
	canvas.gui.frameNr.Draw(canvas.Win, pixel.IM.Scaled(canvas.gui.frameNr.Orig, 1.4))
	canvas.gui.brush.Clear()
	canvas.gui.frameNr.Clear()
	canvas.gui.playbackFPS.Draw(canvas.Win, pixel.IM.Scaled(canvas.gui.playbackFPS.Orig, 1.4))
	canvas.gui.playbackFPS.Clear()

	// update window
	canvas.Win.Update()
}
