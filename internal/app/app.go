package app

import (
	"github.com/VadimGossip/gitPatchTool/internal/config"
	"github.com/VadimGossip/gitPatchTool/internal/domain"
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
	cfg       *domain.Config
}

func NewApp(name, configDir string) *App {
	return &App{
		name:      name,
		configDir: configDir,
	}
}

func (app *App) Run() {
	cfg, err := config.Init(app.configDir)
	if err != nil {
		logrus.Fatalf("Config initialization error %s", err)
	}
	app.cfg = cfg

	dbAdapter := NewDBAdapter(cfg)
	if err := dbAdapter.Connect(); err != nil {
		logrus.Fatalf("Fail to connect data source %s", err)
	}

	app.Factory, err = newFactory(cfg, dbAdapter)
	if err != nil {
		logrus.Fatalf("Fail to init gpt service %s", err)
	}

	if err := app.Factory.oraToolService.CreatePatch(); err != nil {
		logrus.Fatalf("Patch creation failed with err %s", err)
	}
}
