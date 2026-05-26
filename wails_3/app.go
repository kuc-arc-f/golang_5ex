package main

import (
	"context"
	"fmt"

	"github.com/wailsapp/wails/v2/pkg/runtime"
)

// App struct
type App struct {
	ctx context.Context
}

// NewApp creates a new App application struct
func NewApp() *App {
	return &App{}
}

// startup is called when the app starts. The context is saved
// so we can call the runtime methods
func (a *App) startup(ctx context.Context) {
	a.ctx = ctx
}

// Greet returns a greeting for the given name
func (a *App) Greet(name string) string {
	return fmt.Sprintf("Hello %s, It's show time!", name)
}

// JavaScript から呼ばれる
func (a *App) SendMessage(msg string) string {

	println("JSから受信:", msg)

	// JavaScriptへ返信イベント送信
	runtime.EventsEmit(a.ctx, "go-message", "GOから返信: "+msg)

	return "JSから受信:" + msg
	//return "OK"
}