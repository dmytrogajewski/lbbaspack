package main

import (
	"testing"
	"time"

	"github.com/hajimehoshi/ebiten/v2"
)

// Helper to check if an image is all black (or nearly black)
func isAllBlack(img *ebiten.Image) bool {
	w, h := img.Size()
	pixels := make([]byte, 4*w*h)
	img.ReadPixels(pixels)
	for i := 0; i < len(pixels); i += 4 {
		r, g, b, a := pixels[i], pixels[i+1], pixels[i+2], pixels[i+3]
		if a > 0 && (r > 10 || g > 10 || b > 10) {
			return false
		}
	}
	return true
}

// TestGame wraps the real game and allows us to capture the screen after N frames
type TestGame struct {
	*Game
	frame    int
	maxFrame int
	img      *ebiten.Image
	done     chan struct{}
	captured bool
}

func (tg *TestGame) Update() error {
	tg.frame++
	if tg.frame == tg.maxFrame && !tg.captured {
		// Capture the image in the next Draw call
		tg.captured = true
	}
	if tg.frame > tg.maxFrame {
		close(tg.done)
	}
	return tg.Game.Update()
}

func (tg *TestGame) Draw(screen *ebiten.Image) {
	tg.Game.Draw(screen)
	if tg.captured && tg.img == nil {
		w, h := screen.Size()
		tg.img = ebiten.NewImage(w, h)
		tg.img.DrawImage(screen, nil)
	}
}

func (tg *TestGame) Layout(outsideWidth, outsideHeight int) (int, int) {
	return 800, 600
}

func TestGameRendersNonBlack(t *testing.T) {
	done := make(chan struct{})
	tg := &TestGame{
		Game:     NewGame(),
		frame:    0,
		maxFrame: 10,
		done:     done,
	}

	go func() {
		ebiten.RunGame(tg)
	}()

	select {
	case <-done:
		// Test finished
	case <-time.After(2 * time.Second):
		t.Fatal("Test timed out")
	}

	if tg.img == nil {
		t.Fatal("No image captured")
	}

	if isAllBlack(tg.img) {
		t.Error("Screen is all black after running game for several frames")
	}
}
