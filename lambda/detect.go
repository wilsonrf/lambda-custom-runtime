package lambda

import "github.com/paketo-buildpacks/packit"

type Detect struct {
}

func (d *Detect) Detect(context packit.DetectContext) (packit.DetectResult, error) {
	return packit.DetectResult{}, nil
}

func (d *Detect) DetectFunc(context packit.DetectContext) (packit.DetectResult, error) {
	return d.Detect(context)
}
