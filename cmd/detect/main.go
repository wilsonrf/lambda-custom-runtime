package main

import (
	"os"

	"github.com/paketo-buildpacks/packit/v2"
	"github.com/paketo-buildpacks/packit/v2/scribe"
	"github.com/wilsonrf/lambda-custom-runtime/lambda"
)

func main() {
	detect := &lambda.Detect{}
	packit.Detect(detect.DetectFunc)
}
