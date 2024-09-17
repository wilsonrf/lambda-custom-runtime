package lambda_test

import (
	"testing"

	"github.com/sclevine/spec"
	"github.com/sclevine/spec/report"
)

func TestUnit(t *testing.T) {
	suite := spec.New("lambda runtime", spec.Report(report.Terminal{}), spec.Parallel())
	suite("Detect", testDetect)
	suite.Run(t)
}
