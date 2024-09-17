package lambda_test

import (
	"os"
	"testing"

	"github.com/onsi/gomega"
	"github.com/paketo-buildpacks/packit"
	"github.com/paketo-buildpacks/packit/scribe"
	"github.com/sclevine/spec"
	"github.com/wilsonrf/lambda-custom-runtime/lambda"
)

func testDetect(t *testing.T, context spec.G, it spec.S) {
	var (
		Expect = gomega.NewWithT(t).Expect
		ctx    packit.DetectContext
		detect lambda.Detect
	)

	context("BP_NATIVE_IMAGE is not set", func() {
		it.Before(func() {
			Expect(os.Unsetenv("BP_NATIVE_IMAGE")).To(gomega.Succeed())
		})
		it("fails without BP_NATIVE_IMAGE", func() {
			Expect(detect.Detect(ctx)).To(gomega.Equal(packit.DetectResult{}))
		})
	})

	context("BP_NATIVE_IMAGE is set", func() {
		it.Before(func() {
			Expect(os.Setenv("BP_NATIVE_IMAGE", "true")).To(gomega.Succeed())
			detect = lambda.Detect{Logger: scribe.NewLogger(os.Stdout)}
		})
		it("passes with BP_NATIVE_IMAGE is set to true", func() {
			Expect(detect.Detect(ctx)).To(gomega.Equal(packit.DetectResult{
				Plan: packit.BuildPlan{
					Provides: []packit.BuildPlanProvision{
						{Name: lambda.PlanEntryCustomRuntime},
					},
					Requires: []packit.BuildPlanRequirement{
						{Name: lambda.PlanEntryNativeImage},
						{Name: lambda.PlanEntryCustomRuntime},
					},
				},
			}))
		})

		it("fails with BP_NATIVE_IMAGE is set to false", func() {
			Expect(os.Setenv("BP_NATIVE_IMAGE", "false")).To(gomega.Succeed())
			Expect(detect.Detect(ctx)).To(gomega.Equal(packit.DetectResult{}))
		})
	})
}
