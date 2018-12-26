package main

import (
	"fmt"
	"os"

	"github.com/gxlog/gxlog"
	"github.com/gxlog/gxlog/formatter/json"
	"github.com/gxlog/gxlog/formatter/text"
	"github.com/gxlog/gxlog/writer"
	"github.com/gxlog/gxlog/writer/file"
)

func main() {
	// create a new Logger with default config
	log := gxlog.New(gxlog.NewConfig())

	// create a new Config and customize it
	// config := gxlog.NewConfig().
	// 	WithDisabled(gxlog.DynamicContexts | gxlog.Limit).
	// 	WithTrackLevel(gxlog.Off)

	// another equivalent way
	// config := gxlog.NewConfig()
	// config.Flags &^= gxlog.DynamicContexts | gxlog.Limit
	// config.TrackLevel = gxlog.Off

	// create a new Logger with custom config
	// gxlog.New(config)

	// create a file writer, logs output to /tmp/gxlog
	fileWriter, err := file.Open(file.NewConfig("/tmp/gxlog", "base"))
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	defer fileWriter.Close()

	log.Link(gxlog.Slot0, text.New(text.NewConfig().WithEnableColor(true)),
		writer.Wrap(os.Stderr))
	log.Link(gxlog.Slot1, json.New(json.NewConfig()), fileWriter)

	log.Info("test")
}
