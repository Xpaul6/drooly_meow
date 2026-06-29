package main

import (
	"log"
	"math"
	"strconv"

	. "github.com/gen2brain/raylib-go/raylib"
)

type GameState struct {
	clicks        int64
	clickRect     Rectangle
	isClicked     bool
	framesCounter int
}

const WINDOW_WIDTH = 800
const WINDOW_HEIGHT = 600
const CLICK_ANIM_DURATION = 10
const DEFORM_INTENSITY = 70

const CAT_TEXTURE_PATH = "./resources/drooly.png"
const CAT_SFX_PATH = "./resources/meow.mp3"

var gameState GameState = GameState{
	clicks: 0,
	clickRect: Rectangle{
		X:      100,
		Y:      200,
		Width:  200,
		Height: 200,
	},
	framesCounter: 0,
}

func Render(gamestate *GameState, catTexture Texture2D) {
	BeginDrawing()
	defer EndDrawing()

	ClearBackground(Black)

	if gamestate.framesCounter > CLICK_ANIM_DURATION {
		gamestate.framesCounter = 0
		gamestate.clickRect.Width = 200
		gamestate.clickRect.Height = 200
		gamestate.clickRect.X = 100
		gamestate.clickRect.Y = 200
	}
	if gameState.framesCounter <= CLICK_ANIM_DURATION && gamestate.framesCounter != 0 {
		progress := float64(gameState.framesCounter) / CLICK_ANIM_DURATION
		deform := float32(math.Sin(float64(progress)*math.Pi)) * DEFORM_INTENSITY

		gamestate.clickRect.Width = 200 + deform
		gamestate.clickRect.Height = 200 - deform
		gamestate.clickRect.X = 100 - (deform / 2)
		gamestate.clickRect.Y = 200 + (deform / 2)
		gameState.framesCounter += 1
	}

	// DrawRectangleRec(gameState.clickRect, Blue)
	sourceRect := Rectangle{X: 0, Y: 0, Width: float32(catTexture.Width), Height: float32(catTexture.Height)}
	DrawTexturePro(catTexture, sourceRect, gameState.clickRect, Vector2{X: 0, Y: 0}, 0, White)
	DrawText(strconv.Itoa(int(gameState.clicks)), 600, 280, 40, RayWhite)
	// DrawFPS(20, 20)
}

func RegisterMouseClick(gameState *GameState, catSfx *Sound) {
	mousePos := GetMousePosition()
	if CheckCollisionPointRec(mousePos, gameState.clickRect) {
		SetMouseCursor(MouseCursorPointingHand)
		if IsMouseButtonPressed(MouseLeftButton) {
			r := (float32(GetRandomValue(90, 110)) / 100)
			SetSoundPitch(*catSfx, r)
			PlaySound(*catSfx)
			gameState.clicks += 1
			gameState.framesCounter = 1
		}
	} else {
		SetMouseCursor(MouseCursorArrow)
	}
}

func MonitorScore(gameState *GameState) {

}

func main() {

	InitWindow(WINDOW_WIDTH, WINDOW_HEIGHT, "meow meow meow")
	defer CloseWindow()
	SetTargetFPS(60)

	var catTexture Texture2D = LoadTexture(CAT_TEXTURE_PATH)
	if !IsTextureValid(catTexture) {
		log.Fatalln("Texture invalid, aborting...")
	}

	InitAudioDevice()
	var catSfx Sound = LoadSound(CAT_SFX_PATH)
	if !IsSoundValid(catSfx) {
		log.Fatalln("Sound invalid, aborting...")
	}
	SetSoundVolume(catSfx, 0.5)

	// Game loop
	for !WindowShouldClose() {
		RegisterMouseClick(&gameState, &catSfx)
		Render(&gameState, catTexture)
	}
}
