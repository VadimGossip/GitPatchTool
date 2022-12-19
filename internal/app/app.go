package app

import (
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
}

func NewApp(name, configDir string) *App {
	return &App{
		name:      name,
		configDir: configDir,
	}
}

func (app *App) Run() {
	var err error
	app.Factory, err = newFactory()
	if err != nil {
		logrus.Fatalf("Fail to init gpt service %s", err)
	}
	files, err := app.Factory.walkerService.Walk("e:\\WorkSpace\\TCS_Oracle\\", domain.OracleFileType)
	if err != nil {
		logrus.Fatalf("Fail to collect files %s", err)
	}
	logrus.Infof("File received %d", len(files))
}
