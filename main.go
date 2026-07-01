package main

import (
	"log"
	"math"
	"strconv"

	. "github.com/gen2brain/raylib-go/raylib"
)

// Custom types
type GameState struct {
	clicks        int64
	clickRect     Rectangle
	isNotified    bool
	framesCounter int
}

// Constants
const WINDOW_WIDTH = 800
const WINDOW_HEIGHT = 600
const CLICK_ANIM_DURATION = 10
const DEFORM_INTENSITY = 70

const CAT_TEXTURE_PATH = "./resources/drooly.png"
const CAT_SFX_PATH = "./resources/meow.mp3"
const LEVELUP_SFX_PATH = "./resources/levelup.mp3"

// Global game state
var gameState GameState = GameState{
	clicks: 0,
	clickRect: Rectangle{
		X:      100,
		Y:      200,
		Width:  200,
		Height: 200,
	},
	framesCounter: 0,
	isNotified:    false,
}

// ------------------------------------
// ---------- Game functions ----------
// ------------------------------------

func Render(gameState *GameState, catTexture Texture2D) {
	BeginDrawing()
	defer EndDrawing()

	ClearBackground(Black)

	// Click animation
	if gameState.framesCounter > CLICK_ANIM_DURATION {
		gameState.framesCounter = 0
		gameState.clickRect.Width = 200
		gameState.clickRect.Height = 200
		gameState.clickRect.X = 100
		gameState.clickRect.Y = 200
	}
	if gameState.framesCounter <= CLICK_ANIM_DURATION && gameState.framesCounter != 0 {
		progress := float64(gameState.framesCounter) / CLICK_ANIM_DURATION
		deform := float32(math.Sin(float64(progress)*math.Pi)) * DEFORM_INTENSITY

		gameState.clickRect.Width = 200 + deform
		gameState.clickRect.Height = 200 - deform
		gameState.clickRect.X = 100 - (deform / 2)
		gameState.clickRect.Y = 200 + (deform / 2)
		gameState.framesCounter += 1
	}

	// DrawRectangleRec(gameState.clickRect, Blue)
	sourceRect := Rectangle{X: 0, Y: 0, Width: float32(catTexture.Width), Height: float32(catTexture.Height)}
	DrawTexturePro(catTexture, sourceRect, gameState.clickRect, Vector2{X: 0, Y: 0}, 0, White)
	DrawText(strconv.Itoa(int(gameState.clicks)), 600, 280, 40, RayWhite)
	// DrawFPS(20, 20)
}

func RegisterMouseClick(gameState *GameState, sounds *map[string]Sound) {
	mousePos := GetMousePosition()
	if CheckCollisionPointRec(mousePos, gameState.clickRect) {
		// Hovering cursor
		SetMouseCursor(MouseCursorPointingHand)
		// Click register
		if IsMouseButtonPressed(MouseLeftButton) {
			gameState.clicks += 1
			gameState.framesCounter = 1
			gameState.isNotified = false

			SetSoundPitch((*sounds)["catSfx"], (float32(GetRandomValue(90, 110)) / 100))
			PlaySound((*sounds)["catSfx"])
		}
	} else {
		SetMouseCursor(MouseCursorArrow)
	}
}

func MonitorScore(gameState *GameState, sounds *map[string]Sound) {
	if gameState.clicks == 0 {
		return
	}
	if !gameState.isNotified && gameState.clicks%100 == 0 {
		PlaySound((*sounds)["levelUp"])
		gameState.isNotified = true
	}
}

func InitSounds() map[string]Sound {
	var sounds map[string]Sound = make(map[string]Sound)

	sounds["catSfx"] = LoadSound(CAT_SFX_PATH)
	if !IsSoundValid(sounds["catSfx"]) {
		log.Fatalln("Sound invalid, aborting...")
	}
	SetSoundVolume(sounds["catSfx"], 0.5)

	sounds["levelUp"] = LoadSound(LEVELUP_SFX_PATH)
	if !IsSoundValid(sounds["levelUp"]) {
		log.Fatalln("Sound invalid, aborting...")
	}
	SetSoundVolume(sounds["levelUp"], 0.3)

	return sounds
}

// -------------------------------
// ---------- Game loop ----------
// -------------------------------

func main() {
	// Window init
	InitWindow(WINDOW_WIDTH, WINDOW_HEIGHT, "meow meow meow")
	defer CloseWindow()
	SetTargetFPS(60)

	// Textures init
	var catTexture Texture2D = LoadTexture(CAT_TEXTURE_PATH)
	if !IsTextureValid(catTexture) {
		log.Fatalln("Texture invalid, aborting...")
	}
	// Sounds init
	InitAudioDevice()
	var sounds map[string]Sound = InitSounds()

	// Game loop
	for !WindowShouldClose() {
		RegisterMouseClick(&gameState, &sounds)
		MonitorScore(&gameState, &sounds)
		Render(&gameState, catTexture)
	}
}
