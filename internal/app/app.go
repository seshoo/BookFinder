package app

import (
	"fmt"
	"github.com/seshoo/bookFinder/internal/service"
)

type LogOptions struct {
	ToFile    bool `long:"to-file" env:"TO_FILE" description:"log in file"`
	ToConsole bool `long:"to-console" env:"TO_CONSOLE" description:"log in console"`
}

type DataProvider struct {
	UrlTemplate string `long:"dp-url-tmp" env:"DP_URL_TMP" default:"https://rutracker.net/forum/viewtopic.php?t=%s" description:"Template of url"`
}

type Options struct {
	Dbg       bool       `long:"dbg" env:"DBG" description:"debug mode"`
	LogConfig LogOptions `group:"logger" namespace:"logger" env-namespace:"LOGGER"`
	Dp        DataProvider
}

func Run(opts Options) error {
	fmt.Printf("opts: %+v", opts)

	services := service.NewServices(service.Deps{
		DpUrlTmp: opts.Dp.UrlTemplate,
	})

	_ = services

	return nil
}
