package utils_test

import (
	"os"
	"testing"

	"github.com/onsi/gomega"
	"github.com/sclevine/spec"
	util "github.com/wilsonrf/lambda-custom-runtime/lambda/utils"
)

func testBuildpackMetadata(t *testing.T, context spec.G, it spec.S) {
	var (
		Expect = gomega.NewWithT(t).Expect
		path   string
	)

	it.Before(func() {
		file, err := os.CreateTemp("", "buildpack.toml")
		Expect(err).NotTo(gomega.HaveOccurred())
		_, err = file.WriteString(`
[buildpack]
id = "buldpack-id"
name = "buildpack-name"
version = "buildpack-version"

[metadata]
[[metadata.configurations]]
name = "name"
default = "default"
description = "description"
`)
		Expect(err).NotTo(gomega.HaveOccurred())

		Expect(file.Close()).To(gomega.Succeed())

		path = file.Name()

	})

	it.After(func() {
		Expect(os.RemoveAll(path)).To(gomega.Succeed())
	})

	it("ReadBuildpackMetadata", func() {
		cm, err := util.ReadBuildpackMetadata(path)
		Expect(err).NotTo(gomega.HaveOccurred())
		Expect(cm.Buildpack.ID).To(gomega.Equal("buldpack-id"))
		Expect(cm.Buildpack.Name).To(gomega.Equal("buildpack-name"))
		Expect(cm.Buildpack.Version).To(gomega.Equal("buildpack-version"))
		Expect(cm.Metadata.Configurations).To(gomega.HaveLen(1))
		Expect(cm.Metadata.Configurations[0].Name).To(gomega.Equal("name"))
		Expect(cm.Metadata.Configurations[0].Default).To(gomega.Equal("default"))
		Expect(cm.Metadata.Configurations[0].Description).To(gomega.Equal("description"))
	})
}
