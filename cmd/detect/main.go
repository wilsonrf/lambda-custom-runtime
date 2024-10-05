package main

import (
	"os"

	"github.com/paketo-buildpacks/packit/v2"
	"github.com/paketo-buildpacks/packit/v2/scribe"
	"github.com/wilsonrf/lambda-custom-runtime/lambda"
)

func main() {
	logger := scribe.NewLogger(os.Stdout)
	detect := &lambda.Detect{Logger: logger}
	packit.Detect(detect.DetectFunc)
}
