package main

import (
	"embed"
	"hamster-client/ctx"
	"log"

	"github.com/wailsapp/wails/v2"
	"github.com/wailsapp/wails/v2/pkg/logger"
	"github.com/wailsapp/wails/v2/pkg/options"
	"github.com/wailsapp/wails/v2/pkg/options/windows"
)

//go:embed frontend/dist/*
var assets embed.FS

func main() {

	app := ctx.NewApp()
	// Create application with options
	err := wails.Run(&options.App{
		Title:  "hamster-client",
		Width:  1024,
		Height: 768,
		// MinWidth:          720,
		// MinHeight:         570,
		// MaxWidth:          1280,
		// MaxHeight:         740,
		DisableResize:     false,
		Fullscreen:        false,
		Frameless:         false,
		StartHidden:       false,
		HideWindowOnClose: false,
		RGBA:              &options.RGBA{R: 255, G: 255, B: 255, A: 255},
		Assets:            assets,
		LogLevel:          logger.DEBUG,
		OnStartup:         app.Startup,
		OnDomReady:        app.DomReady,
		OnShutdown:        app.Shutdown,
		Bind: []interface{}{
			&app.AccountApp,
			&app.P2pApp,
			&app.ResourceApp,
			&app.SettingApp,
			&app.WalletApp,
			&app.DeployApp,
			&app.ApplicationApp,
			&app.GraphApp,
			&app.KeyStorageApp,
			&app.QueueApp,
			&app.ChainManagerApp,
		},
		// Windows platform specific options
		Windows: &windows.Options{
			WebviewIsTransparent: false,
			WindowIsTranslucent:  false,
			DisableWindowIcon:    false,
		},
	})

	if err != nil {
		log.Fatal(err)
	}
}
