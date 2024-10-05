package lambda

import (
	"path/filepath"

	"github.com/paketo-buildpacks/packit/v2"
	"github.com/paketo-buildpacks/packit/v2/cargo"
	"github.com/paketo-buildpacks/packit/v2/postal"
	"github.com/paketo-buildpacks/packit/v2/scribe"
	"github.com/wilsonrf/lambda-custom-runtime/lambda/utils"
)

type Build struct {
	Logger  scribe.Logger
	Context *packit.BuildContext
}

func (b *Build) Build(context packit.BuildContext) (packit.BuildResult, error) {

	b.Logger.Title("%s: %s", context.BuildpackInfo.Name, context.BuildpackInfo.Version)
	b.Logger.Process(context.BuildpackInfo.Homepage)

	m, err := utils.ReadBuildpackMetadata(filepath.Join(context.CNBPath, "buildpack.toml"))

	if err != nil {
		b.Logger.Process("Error reading buildpack metadata: %s", err)
		return packit.BuildResult{}, err
	}

	logMetadata(&b.Logger, m)

	result := packit.BuildResult{}

	pe := context.Plan.Entries

	transport := cargo.NewTransport()
	service := postal.NewService(transport)

	for _, entry := range pe {
		switch entry.Name {
		case PlanEntryCustomRuntimeEmulator:
			emulator := NewEmulator(&b.Logger, &context, &service)
			layer, process, err := emulator.Create()
			if err != nil {
				return result, err
			}
			result.Layers = append(result.Layers, layer)
			result.Launch.Processes = append(result.Launch.Processes, process)
		}
	}

	cr := NewCustomRuntime(&b.Logger, &context, &service)
	layer, process, err := cr.Create()

	if err != nil {
		b.Logger.Process("Error creating custom runtime: %s", err)
		return result, err
	}

	result.Layers = append(result.Layers, layer)
	result.Launch.Processes = append(result.Launch.Processes, process)

	return result, nil
}

func (b *Build) BuildFunc(context packit.BuildContext) (packit.BuildResult, error) {
	return b.Build(context)
}

func logMetadata(logger *scribe.Logger, m utils.BuildpackMetadata) {
	logger.Process("Buildpack Configurations:")
	for _, c := range m.Metadata.Configurations {
		logger.Subprocess("%-15s %-15s %-15s", c.Name, c.Default, c.Description)
	}
}
