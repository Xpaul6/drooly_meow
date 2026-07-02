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

// Constants & globals
const WINDOW_WIDTH = 800
const WINDOW_HEIGHT = 600
const CLICK_ANIM_DURATION = 10
const DEFORM_INTENSITY = 70

var TEXTURES_MAP = map[string]string {
	"catTexture": "./resources/drooly.png",
}

var SOUNDS_MAP = map[string]string {
	"catSfx":  "./resources/meow.mp3",
	"levelUp": "./resources/levelup.mp3",
}

// ------------------------------------
// ---------- Game functions ----------
// ------------------------------------

func Render(gameState *GameState, textures *map[string]Texture2D) {
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
	sourceRect := Rectangle{X: 0, Y: 0, Width: float32((*textures)["catTexture"].Width), Height: float32((*textures)["catTexture"].Height)}
	DrawTexturePro((*textures)["catTexture"], sourceRect, gameState.clickRect, Vector2{X: 0, Y: 0}, 0, White)
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

func InitTextures() map[string]Texture2D {
	var textures map[string]Texture2D = make(map[string]Texture2D)

	for k, v := range TEXTURES_MAP {
		textures[k] = LoadTexture(v)
	}

	for k, v := range textures {
		if !IsTextureValid(v) {
			log.Fatalf("Error importing %s texture, aborting...", k)
		}
	}

	return textures
}

func InitSounds() map[string]Sound {
	var sounds map[string]Sound = make(map[string]Sound)

	for k, v := range SOUNDS_MAP {
		sounds[k] = LoadSound(v)
	}

	// Volume controls
	SetSoundVolume(sounds["catSfx"], 0.5)
	SetSoundVolume(sounds["levelUp"], 0.3)

	// Samples validation
	for k, v := range sounds {
		if !IsSoundValid(v) {
			log.Fatalf("Error importing %s sfx, aborting...", k)
		}
	}

	return sounds
}

// -------------------------------
// ---------- Game loop ----------
// -------------------------------

func main() {
	// Game state init
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

	// Window init
	InitWindow(WINDOW_WIDTH, WINDOW_HEIGHT, "meow meow meow")
	defer CloseWindow()
	SetTargetFPS(60)

	// Textures init
	var textures map[string]Texture2D = InitTextures()

	// Sounds init
	InitAudioDevice()
	var sounds map[string]Sound = InitSounds()

	// Game loop
	for !WindowShouldClose() {
		RegisterMouseClick(&gameState, &sounds)
		MonitorScore(&gameState, &sounds)
		Render(&gameState, &textures)
	}
}
