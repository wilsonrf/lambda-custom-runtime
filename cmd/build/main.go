package main

import (
	"os"

	"github.com/paketo-buildpacks/packit/v2"
	"github.com/paketo-buildpacks/packit/v2/scribe"
	"github.com/wilsonrf/lambda-custom-runtime/lambda"
)

func main() {
	logger := scribe.NewLogger(os.Stdout)
	build := &lambda.Build{Logger: logger}
	packit.Build(build.BuildFunc)
}
