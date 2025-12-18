// Package main is a module which serves the fps overlay module
package main

import (
	"context"

	"github.com/viam-labs/overlay-fps/overlayfps"
	"go.viam.com/rdk/components/camera"
	"go.viam.com/rdk/logging"
	"go.viam.com/rdk/module"
	"go.viam.com/utils"
)

func main() {
	utils.ContextualMain(mainWithArgs, module.NewLoggerFromArgs("overlay-fps"))
}

func mainWithArgs(ctx context.Context, args []string, logger logging.Logger) (err error) {
	oMod, err := module.NewModuleFromArgs(ctx, logger)
	if err != nil {
		return err
	}

	err = oMod.AddModelFromRegistry(ctx, camera.API, overlayfps.Model)
	if err != nil {
		return err
	}

	err = oMod.Start(ctx)
	defer oMod.Close(ctx)
	if err != nil {
		return err
	}
	<-ctx.Done()
	return nil
}
