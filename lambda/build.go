package lambda

import (
	"github.com/paketo-buildpacks/packit/v2"
	"github.com/paketo-buildpacks/packit/v2/cargo"
	"github.com/paketo-buildpacks/packit/v2/postal"
	"github.com/paketo-buildpacks/packit/v2/scribe"
)

type Build struct {
	Logger  scribe.Logger
	Context *packit.BuildContext
}

func (b *Build) Build(context packit.BuildContext) (packit.BuildResult, error) {
	return packit.BuildResult{}, nil
}

func (b *Build) BuildFunc(context packit.BuildContext) (packit.BuildResult, error) {
	return b.Build(context)
}
