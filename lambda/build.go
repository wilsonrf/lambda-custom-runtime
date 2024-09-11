package lambda

import (
	"github.com/paketo-buildpacks/packit"
	"github.com/paketo-buildpacks/packit/scribe"
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
