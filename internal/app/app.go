package app

import (
	"github.com/sirupsen/logrus"
	"os"
)

func init() {
	logrus.SetFormatter(&logrus.JSONFormatter{})
	logrus.SetOutput(os.Stdout)
	logrus.SetLevel(logrus.InfoLevel)
}

type App struct {
	*Factory
	name      string
	configDir string
}

func NewApp(name, configDir string) *App {
	return &App{
		name:      name,
		configDir: configDir,
	}
}

func (app *App) Run() {
	var err error

	dbAdapter := NewDBAdapter()
	if err := dbAdapter.Connect(); err != nil {
		logrus.Fatalf("Fail to connect db %s", err)
	}

	app.Factory, err = newFactory(dbAdapter)
	if err != nil {
		logrus.Fatalf("Fail to init gpt service %s", err)
	}

	if err := app.Factory.oraToolService.CreatePatch(); err != nil {
		logrus.Fatalf("Patch creation failed with err %s", err)
	}
}
