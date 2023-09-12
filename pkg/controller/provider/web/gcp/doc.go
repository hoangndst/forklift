package gcp

import (
	api "github.com/konveyor/forklift-controller/pkg/apis/forklift/v1beta1"
	"github.com/konveyor/forklift-controller/pkg/controller/provider/web/base"
)

// Routes
const (
	Root = base.ProvidersRoot + "/" + string(api.GCP)
)
