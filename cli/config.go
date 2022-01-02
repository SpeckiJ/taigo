package main

import (
	"log"
	"os"

	"github.com/gookit/config/v2"
)

var configStruct = struct {
	TaigoHost         string            `mapstructure:"TAIGO_HOST"`
	TaigoToken        string            `mapstructure:"TAIGO_TOKEN"`
	TaigoUsername     string            `mapstructure:"TAIGO_USERNAME"`
	FavouriteProjects FavouriteProjects `mapstructure:"FavouriteProjects"`
}{}

func dumpConfigToFile(path string) error {
	// dump default config b/c it doesn't exist
	f, err := os.Create(path)
	if err != nil {
		return err
	}
	defer f.Close()
	config.Set("AppName", AppName)
	if _, err = config.WriteTo(f); err != nil {
		return err
	}
	log.Printf("config written to %s\n", path)
	return nil
}
