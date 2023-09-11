package gcp

import (
	"context"
	fb "github.com/konveyor/forklift-controller/pkg/lib/filebacked"
	libmodel "github.com/konveyor/forklift-controller/pkg/lib/inventory/model"
	"github.com/konveyor/forklift-controller/pkg/lib/logging"
)

// All adapters.
var adapterList []Adapter

// Updates the DB based on
// changes described by an Event.
type Updater func(tx *libmodel.Tx) error

// Adapter context.
type Context struct {
	// Context.
	ctx context.Context
	// DB client.
	db libmodel.DB
	// OpenStack client.
	client *Client
	// Log.
	log logging.LevelLogger
}

// The adapter request is canceled.
func (r *Context) canceled() (done bool) {
	select {
	case <-r.ctx.Done():
		done = true
	default:
	}

	return
}

// Model adapter.
// Provides integration between the REST resource
// model and the inventory model.
type Adapter interface {
	// List REST collections.
	List(ctx *Context) (itr fb.Iterator, err error)
	// Get object updates
	GetUpdates(ctx *Context) (updates []Updater, err error)
	// Clean unexisting objects within the database
	DeleteUnexisting(ctx *Context) (updates []Updater, err error)
}
