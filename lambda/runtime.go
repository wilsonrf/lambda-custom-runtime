package lambda

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/paketo-buildpacks/packit/v2"
	"github.com/paketo-buildpacks/packit/v2/postal"
	"github.com/paketo-buildpacks/packit/v2/scribe"
	"github.com/wilsonrf/lambda-custom-runtime/lambda/utils"
)

const (
	EmulatorLayerName  = "emulator"
	EmulatorProccess   = "aws-lambda-rie"
	EmulatorDependency = "aws-lambda-runtime-interface-emulator"
)

type Runtime interface {
	Create() (packit.Layer, packit.Process, error)
}

type CustomRuntime struct {
	Logger  *scribe.Logger
	Context *packit.BuildContext
	Service *postal.Service
}

func NewCustomRuntime(logger *scribe.Logger, context *packit.BuildContext, service *postal.Service) CustomRuntime {
	return CustomRuntime{
		Logger:  logger,
		Context: context,
		Service: service,
	}
}

func (e *CustomRuntime) Create() (packit.Layer, packit.Process, error) {
	e.Logger.Process("Creating runtime layer")

	layer, err := e.Context.Layers.Get("custom-runtime")

	if err != nil {
		return packit.Layer{}, packit.Process{}, fmt.Errorf("failed to get custom runtime layer: %w", err)
	}

	layer.Launch = true
	layer.Cache = true
	layer.Build = false

	// wd is the working directory where the binary file is located
	wd, err := os.ReadDir(e.Context.WorkingDir)

	if err != nil {
		return packit.Layer{}, packit.Process{}, fmt.Errorf("failed to read working directory: %w", err)
	}

	if len(wd) == 0 {
		return packit.Layer{}, packit.Process{}, fmt.Errorf("no files found in the working directory")
	}

	cmd := fmt.Sprintf("%s none", filepath.Join(e.Context.WorkingDir, wd[0].Name()))

	e.Logger.Process("Custom Runtime Command: %s", cmd)

	p := packit.Process{
		Type:    "custom-runtime",
		Command: cmd,
		Args:    []string{},
		Default: true,
		Direct:  true,
	}

	return layer, p, nil
}

type Emulator struct {
	Logger  *scribe.Logger
	Context *packit.BuildContext
	Service *postal.Service
}

func NewEmulator(logger *scribe.Logger, context *packit.BuildContext, service *postal.Service) Emulator {

	return Emulator{
		Logger:  logger,
		Context: context,
		Service: service,
	}
}

func (e *Emulator) Create() (packit.Layer, packit.Process, error) {

	layer, err := createLayer(e.Context, e.Logger, e.Service)

	if err != nil {
		return packit.Layer{}, packit.Process{}, err
	}

	lf, err := os.ReadDir(layer.Path)

	if err != nil {
		return packit.Layer{}, packit.Process{}, fmt.Errorf("failed to read layer directory: %w", err)
	}

	if len(lf) == 0 {
		return packit.Layer{}, packit.Process{}, fmt.Errorf("no files found in the emulator layer")
	}

	// wd is the working directory where the binary file is located
	wd, err := os.ReadDir(e.Context.WorkingDir)

	if err != nil {
		return packit.Layer{}, packit.Process{}, fmt.Errorf("failed to read working directory: %w", err)
	}

	if len(wd) == 0 {
		return packit.Layer{}, packit.Process{}, fmt.Errorf("no files found in the working directory")
	}

	cmd := filepath.Join(layer.Path, lf[0].Name())
	args := fmt.Sprintf("%s none", filepath.Join(e.Context.WorkingDir, wd[0].Name()))

	emulatorProccess := createProcess(cmd, args, true, e.Logger)

	return layer, emulatorProccess, nil
}

func createLayer(context *packit.BuildContext, logger *scribe.Logger, service *postal.Service) (packit.Layer, error) {
	logger.Process("Creating emulator layer")
	layer, err := context.Layers.Get(EmulatorLayerName)

	layer.Reset()

	layer.Build = false
	layer.Launch = true
	layer.Cache = true

	if err != nil {
		return packit.Layer{}, fmt.Errorf("failed to get emulator layer: %w", err)
	}

	ie, err := service.Resolve(filepath.Join(context.CNBPath, "buildpack.toml"), EmulatorDependency, "", "*")

	if err != nil {
		return packit.Layer{}, fmt.Errorf("failed to resolve emulator dependency: %w", err)
	}

	err = utils.DownloadFile(ie.URI, filepath.Join(layer.Path, ie.ID))

	if err != nil {
		logger.Process("Failed to deliver emulator dependency")
		return packit.Layer{}, fmt.Errorf("failed to download emulator file: %w", err)
	}

	return layer, nil
}

func createProcess(command string, arguments string, isDefault bool, logger *scribe.Logger) packit.Process {

	logger.Process("Emulator Command: %s %s", command, arguments)

	args := strings.Fields(arguments)

	emultorProccess := packit.Process{
		Type:    "emulator",
		Args:    args,
		Command: command,
		Default: isDefault,
		Direct:  true,
	}

	return emultorProccess
}
