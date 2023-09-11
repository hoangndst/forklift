package gcp

import (
	"github.com/konveyor/forklift-controller/pkg/controller/plan/adapter/base"
	plancontext "github.com/konveyor/forklift-controller/pkg/controller/plan/context"
)

type Adapter struct{}

func (r *Adapter) Builder(ctx *plancontext.Context) (builder base.Builder, err error) {
	builder = &Builder{Context: ctx}
	return
}
