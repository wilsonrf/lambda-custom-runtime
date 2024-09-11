package main

import (
	"github.com/paketo-buildpacks/packit"
	"github.com/wilsonrf/lambda-custom-runtime/lambda"
)

func main() {
	detect := &lambda.Detect{}
	packit.Detect(detect.DetectFunc)
}
