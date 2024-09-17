package lambda

import (
	"os"

	"github.com/paketo-buildpacks/packit"
	"github.com/paketo-buildpacks/packit/scribe"
)

const (
	ConfigNativeImage              = "BP_NATIVE_IMAGE"
	PlanEntryNativeImage           = "native-image-application"
	PlanEntryCustomRuntime         = "lambda-custom-runtime"
	PlanEntryCustomRuntimeEmulator = "lambda-custom-runtime-emulator"
)

type Detect struct {
	Logger scribe.Logger
}

func (d *Detect) Detect(context packit.DetectContext) (packit.DetectResult, error) {

	result := packit.DetectResult{}

	if env, ok := os.LookupEnv(ConfigNativeImage); ok {
		if env == "true" {
			d.Logger.Process("PASSED: BP_NATIVE_IMAGE is true")
			result = packit.DetectResult{
				Plan: packit.BuildPlan{
					Provides: []packit.BuildPlanProvision{
						{Name: PlanEntryCustomRuntime},
					},
					Requires: []packit.BuildPlanRequirement{
						{Name: PlanEntryNativeImage},
						{Name: PlanEntryCustomRuntime},
					},
				},
			}

			if emu, ok := os.LookupEnv("BP_LAMBDA_CUSTOM_RUNTIME_INTERFACE_EMULATOR"); ok {
				if emu == "true" {
					d.Logger.Process("PASSED: BP_LAMBDA_CUSTOM_RUNTIME_INTERFACE_EMULATOR is true")
					result.Plan.Provides = append(result.Plan.Provides, packit.BuildPlanProvision{Name: PlanEntryCustomRuntimeEmulator})
				}
			}
		} else {
			result = packit.DetectResult{}
		}
	} else {
		result = packit.DetectResult{}
	}

	return result, nil
}

func (d *Detect) DetectFunc(context packit.DetectContext) (packit.DetectResult, error) {
	return d.Detect(context)
}
