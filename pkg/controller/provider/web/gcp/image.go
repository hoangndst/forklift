package gcp

import (
	"errors"
	"github.com/gin-gonic/gin"
	api "github.com/konveyor/forklift-controller/pkg/apis/forklift/v1beta1"
	model "github.com/konveyor/forklift-controller/pkg/controller/provider/model/gcp"
	"github.com/konveyor/forklift-controller/pkg/controller/provider/web/base"
	libmodel "github.com/konveyor/forklift-controller/pkg/lib/inventory/model"
	"net/http"
)

// Routes
const (
	ImageParam      = "image"
	ImageCollection = "images"
	ImagesRoot      = ProviderRoot + "/" + ImageCollection
	ImageRoot       = ImagesRoot + "/:" + ImageParam
)

// Image handler.
type ImageHandler struct {
	Handler
}

type Image struct {
	Resource
	Architecture      string            `json:"architecture"`
	ArchiveSizeBytes  int64             `json:"archiveSizeBytes"`
	CreationTimestamp string            `json:"creationTimestamp"`
	Description       string            `json:"description"`
	DiskSizeGb        int64             `json:"diskSizeGb"`
	Family            string            `json:"family"`
	Kind              string            `json:"kind"`
	LabelFingerprint  string            `json:"labelFingerprint"`
	Labels            map[string]string `json:"labels"`
	LicenseCodes      []int64           `json:"licenseCodes"`
	Licenses          []string          `json:"licenses"`
	SatisfiesPzs      bool              `json:"satisfiesPzs"`
	SelfLink          string            `json:"selfLink"`
	SourceDisk        string            `json:"sourceDisk"`
	SourceDiskId      string            `json:"sourceDiskId"`
	SourceImage       string            `json:"sourceImage"`
	SourceImageId     string            `json:"sourceImageId"`
	SourceSnapshot    string            `json:"sourceSnapshot"`
	SourceSnapshotId  string            `json:"sourceSnapshotId"`
	SourceType        string            `json:"sourceType"`
	Status            string            `json:"status"`
	StorageLocations  []string          `json:"storageLocations"`
}

// Add routes to the `gin` router.
func (h *ImageHandler) AddRoutes(e *gin.Engine) {
	e.GET(ImagesRoot, h.List)
	e.GET(ImagesRoot+"/", h.List)
	e.GET(ImageRoot, h.Get)
}

// Build the resource using the model.
func (r *Image) With(m *model.Image) {
	r.Resource.ID = m.ID
	r.Resource.Revision = m.Revision
	r.Resource.Name = *m.Name
	r.Status = *m.Status
	r.Architecture = *m.Architecture
	r.ArchiveSizeBytes = *m.ArchiveSizeBytes
	r.CreationTimestamp = *m.CreationTimestamp
	r.Description = *m.Description
	r.DiskSizeGb = *m.DiskSizeGb
	r.Family = *m.Family
	r.Kind = *m.Kind
	r.LabelFingerprint = *m.LabelFingerprint
	r.Labels = m.Labels
	r.LicenseCodes = m.LicenseCodes
	r.Licenses = m.Licenses
	r.SatisfiesPzs = *m.SatisfiesPzs
	r.SelfLink = *m.SelfLink
	r.SourceDisk = *m.SourceDisk
	r.SourceDiskId = *m.SourceDiskId
	r.SourceImage = *m.SourceImage
	r.SourceImageId = *m.SourceImageId
	r.SourceSnapshot = *m.SourceSnapshot
	r.SourceSnapshotId = *m.SourceSnapshotId
	r.SourceType = *m.SourceType
	r.StorageLocations = m.StorageLocations
}

// List resources in a REST collection.
// A GET onn the collection that includes the `X-Watch`
// header will negotiate an upgrade of the connection
// to a websocket and push watch events.
func (h ImageHandler) List(ctx *gin.Context) {
	status, err := h.Prepare(ctx)
	if status != http.StatusOK {
		ctx.Status(status)
		base.SetForkliftError(ctx, err)
		return
	}
	if h.WatchRequest {
		h.watch(ctx)
		return
	}
	db := h.Collector.DB()
	list := []model.Image{}
	err = db.List(&list, h.ListOptions(ctx))
	if err != nil {
		log.Trace(
			err,
			"url",
			ctx.Request.URL)
		ctx.Status(http.StatusInternalServerError)
		return
	}
	content := []interface{}{}
	for _, m := range list {
		r := &Image{}
		r.With(&m)
		r.Link(h.Provider)
		content = append(content, r.Content(h.Detail))
	}

	ctx.JSON(http.StatusOK, content)
}

// Get a specific REST resource.
func (h ImageHandler) Get(ctx *gin.Context) {
	status, err := h.Prepare(ctx)
	if status != http.StatusOK {
		ctx.Status(status)
		base.SetForkliftError(ctx, err)
		return
	}
	h.Detail = model.MaxDetail
	m := &model.Image{
		Base: model.Base{
			ID: ctx.Param(ImageParam),
		},
	}
	db := h.Collector.DB()
	err = db.Get(m)
	if errors.Is(err, model.NotFound) {
		ctx.Status(http.StatusNotFound)
		return
	}
	if err != nil {
		log.Trace(
			err,
			"url",
			ctx.Request.URL)
		ctx.Status(http.StatusInternalServerError)
		return
	}
	r := &Image{}
	r.With(m)
	r.Link(h.Provider)
	content := r.Content(h.Detail)

	ctx.JSON(http.StatusOK, content)
}

// Watch.
func (h *ImageHandler) watch(ctx *gin.Context) {
	db := h.Collector.DB()
	err := h.Watch(
		ctx,
		db,
		&model.Image{},
		func(in libmodel.Model) (r interface{}) {
			m := in.(*model.Image)
			image := &Image{}
			image.With(m)
			image.Link(h.Provider)
			r = image
			return
		})
	if err != nil {
		log.Trace(
			err,
			"url",
			ctx.Request.URL)
		ctx.Status(http.StatusInternalServerError)
	}
}

// Build self link (URI).
func (r *Image) Link(p *api.Provider) {
	r.SelfLink = base.Link(
		ImageRoot,
		base.Params{
			base.ProviderParam: string(p.UID),
			ImageParam:         r.ID,
		})
}

// As content.
func (r *Image) Content(detail int) interface{} {
	if detail == 0 {
		return r.Resource
	}

	return r
}
