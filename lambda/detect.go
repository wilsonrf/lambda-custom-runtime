package lambda

import (
	"os"

	"github.com/paketo-buildpacks/packit"
	"github.com/paketo-buildpacks/packit/scribe"
)

const (
	ConfigNativeImage      = "BP_NATIVE_IMAGE"
	PlanEntryNativeImage   = "native-image-application"
	PlanEntryCustomRuntime = "lambda-custom-runtime"
)

type Detect struct {
	Logger scribe.Logger
}

func (d *Detect) Detect(context packit.DetectContext) (packit.DetectResult, error) {

	if env, ok := os.LookupEnv(ConfigNativeImage); ok {
		if env == "true" {
			d.Logger.Process("PASSED: BP_NATIVE_IMAGE is true")
			return packit.DetectResult{
				Plan: packit.BuildPlan{
					Provides: []packit.BuildPlanProvision{
						{Name: PlanEntryCustomRuntime},
					},
					Requires: []packit.BuildPlanRequirement{
						{Name: PlanEntryNativeImage},
						{Name: PlanEntryCustomRuntime},
					},
				},
			}, nil
		} else {
			return packit.DetectResult{}, nil
		}
	} else {
		return packit.DetectResult{}, nil
	}
}

func (d *Detect) DetectFunc(context packit.DetectContext) (packit.DetectResult, error) {
	return d.Detect(context)
}
