package lambda_test

import (
	"fmt"
	"os"
	"path/filepath"
	"testing"

	"github.com/onsi/gomega"
	"github.com/paketo-buildpacks/packit/v2"
	"github.com/paketo-buildpacks/packit/v2/cargo"
	"github.com/paketo-buildpacks/packit/v2/postal"
	"github.com/paketo-buildpacks/packit/v2/scribe"
	"github.com/sclevine/spec"
	"github.com/wilsonrf/lambda-custom-runtime/lambda"
	"github.com/wilsonrf/lambda-custom-runtime/lambda/utils"
)

func testRuntime(t *testing.T, context spec.G, it spec.S) {
	var (
		Expect  = gomega.NewWithT(t).Expect
		ctx     packit.BuildContext
		logger  scribe.Logger
		service postal.Service
	)

	it.Before(func() {

		logger = scribe.NewLogger(os.Stdout)

		t := cargo.NewTransport()
		service = postal.NewService(t)

		cnbPath, err := os.MkdirTemp("", "cnb")
		Expect(err).NotTo(gomega.HaveOccurred())

		workingDir, err := os.MkdirTemp("", "workspace")
		Expect(err).NotTo(gomega.HaveOccurred())

		layerDir, err := os.MkdirTemp("", "layer")
		Expect(err).NotTo(gomega.HaveOccurred())

		err = utils.CopyFile("testdata/hello-lambda", filepath.Join(workingDir, "hello-lambda"))

		Expect(err).NotTo(gomega.HaveOccurred())

		ctx = packit.BuildContext{
			WorkingDir: workingDir,
			Layers:     packit.Layers{Path: layerDir},
			CNBPath:    cnbPath,
		}

	})

	it.After(func() {
		os.RemoveAll(ctx.WorkingDir)
		os.RemoveAll(ctx.Layers.Path)
	})

	context("Custom Runtime", func() {
		it("creates a new customruntime", func() {

			r := lambda.NewCustomRuntime(&logger, &ctx, &service)
			l, p, err := r.Create()
			Expect(err).NotTo(gomega.HaveOccurred())
			Expect(l).NotTo(gomega.BeNil())
			Expect(p).NotTo(gomega.BeNil())
			Expect(l.Name).To(gomega.Equal("custom-runtime"))
			Expect(l.Launch).To(gomega.BeTrue())
			Expect(l.Cache).To(gomega.BeTrue())
			Expect(l.Build).To(gomega.BeFalse())
			Expect(p.Command).To(gomega.Equal(fmt.Sprintf("%s none", filepath.Join(ctx.WorkingDir, "hello-lambda"))))
		})
	})

	context("Emulator Runtime", func() {

		it.Before(func() {

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
    	name = "BP_LAMBDA_CUSTOM_RUNTIME_INTERFACE_EMULATOR"`

			err := os.WriteFile(filepath.Join(ctx.CNBPath, "buildpack.toml"), []byte(buildpackTomlContent), 0644)

			Expect(err).NotTo(gomega.HaveOccurred())
		})

		it("creates a new emulator runtime", func() {
			r := lambda.NewEmulator(&logger, &ctx, &service)
			l, p, err := r.Create()
			Expect(err).NotTo(gomega.HaveOccurred())
			Expect(l).NotTo(gomega.BeNil())
			Expect(p).NotTo(gomega.BeNil())
			Expect(l.Name).To(gomega.Equal("emulator"))
			Expect(l.Launch).To(gomega.BeTrue())
			Expect(l.Cache).To(gomega.BeTrue())
			Expect(l.Build).To(gomega.BeFalse())
			Expect(p.Command).To(gomega.Equal(fmt.Sprintf("%s/%s", filepath.Join(ctx.Layers.Path, "emulator"), "aws-lambda-runtime-interface-emulator")))
			Expect(p.Args).To(gomega.Equal([]string{filepath.Join(ctx.WorkingDir, "hello-lambda"), "none"}))
		})
	})

}
