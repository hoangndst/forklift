package gcp

import "github.com/konveyor/forklift-controller/pkg/controller/provider/web/base"

// Routes.
const (
	ProviderParam = base.ProviderParam
	ProvidersRoot = Root
	ProviderRoot  = ProvidersRoot + "/:" + ProviderParam
)
