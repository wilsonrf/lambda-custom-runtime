package lambda_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/onsi/gomega"
	"github.com/paketo-buildpacks/packit/v2"
	"github.com/paketo-buildpacks/packit/v2/scribe"
	"github.com/sclevine/spec"
	"github.com/wilsonrf/lambda-custom-runtime/lambda"
	"github.com/wilsonrf/lambda-custom-runtime/lambda/utils"
)

func testBuild(t *testing.T, context spec.G, it spec.S) {

	var (
		Expect = gomega.NewWithT(t).Expect
		ctx    packit.BuildContext
		build  lambda.Build
	)

	it.Before(func() {
		cnbPath, err := os.MkdirTemp("", "cnb")
		Expect(err).NotTo(gomega.HaveOccurred())

		workingDir, err := os.MkdirTemp("", "workspace")
		Expect(err).NotTo(gomega.HaveOccurred())

		layerDir, err := os.MkdirTemp("", "layer")
		Expect(err).NotTo(gomega.HaveOccurred())

		ctx.CNBPath = cnbPath
		ctx.WorkingDir = workingDir
		ctx.Layers = packit.Layers{Path: layerDir}

		buildpackTomlContent := `
		[buildpack]
  		id = "com.wilsonfranca.lambda-custom-runtime"
  		name = "Lambda Custom Runtime Buildpack"
  		version = "0.0.1" 
		[metadata]
		[[metadata.dependencies]]
		id      = "aws-lambda-runtime-interface-emulator"
		version = "1.15"
		uri     = "https://github.com/aws/aws-lambda-runtime-interface-emulator/releases/download/v1.15/aws-lambda-rie"
		sha256  = "a2bca0ff67c5435a02bf28a85524a8ff2ec222f403c19d92fe304f3f7c7cce10"
		stacks  = ["*"]
		[[metadata.configurations]]
    	build = false
    	default = "false"
    	description = "whether to contribute with a emulator layer"
    	name = "BP_LAMBDA_CUSTOM_RUNTIME_INTERFACE_EMULATOR"
		`

		err = os.WriteFile(filepath.Join(cnbPath, "buildpack.toml"), []byte(buildpackTomlContent), 0644)

		if err != nil {
			t.Fatal(err)
		}

		err = utils.CopyFile("testdata/hello-lambda", filepath.Join(ctx.WorkingDir, "hello-lambda"))

		if err != nil {
			t.Fatal(err)
		}

		build = lambda.Build{Logger: scribe.NewLogger(os.Stdout)}
	})

	it.After(func() {
		Expect(os.RemoveAll(ctx.CNBPath)).To(gomega.Succeed())
		Expect(os.RemoveAll(ctx.WorkingDir)).To(gomega.Succeed())
		Expect(os.RemoveAll(ctx.Layers.Path)).To(gomega.Succeed())
	})

	it("contributes a layer containing the runtime interface", func() {
		result, err := build.Build(ctx)
		Expect(err).NotTo(gomega.HaveOccurred())
		Expect(result.Layers).To(gomega.HaveLen(1))
		Expect(result.Layers[0].Name).To(gomega.Equal("custom-runtime"))
	})

	it("contributes a layer containing the runtime interface emulator", func() {
		ctx.Plan.Entries = append(ctx.Plan.Entries, packit.BuildpackPlanEntry{Name: lambda.PlanEntryCustomRuntimeEmulator})
		result, err := build.Build(ctx)
		Expect(err).NotTo(gomega.HaveOccurred())
		Expect(result.Layers).To(gomega.HaveLen(2))
		Expect(result.Layers[0].Name).To(gomega.Equal("emulator"))
		Expect(result.Layers[1].Name).To(gomega.Equal("custom-runtime"))
	})

}
