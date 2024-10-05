package utils_test

import (
	"testing"

	"github.com/sclevine/spec"
	"github.com/sclevine/spec/report"
)

func TestUnit(t *testing.T) {
	suite := spec.New("util", spec.Report(report.Terminal{}), spec.Parallel())
	suite("BuildpackMetadata", testBuildpackMetadata)
	suite.Run(t)
}
