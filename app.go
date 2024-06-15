package main

import (
	"context"
	"evil-soundcloud/pkg/soundcloud"
	"os"
)

type App struct {
	baseDir string
	ctx     context.Context
}

func NewApp() *App {
	return &App{}
}

func (a *App) startup(ctx context.Context) {
	a.ctx = ctx
	home, err := os.UserHomeDir()
	if err != nil {
		home = ""
	}
	a.baseDir = home + "/Downloads"
}

func (a *App) FetchPlaylist(url string) {
	soundcloud.GetTracks(url, a.baseDir)
}
