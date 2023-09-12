package gcp

import (
	"github.com/konveyor/forklift-controller/pkg/controller/provider/model/base"
	libref "github.com/konveyor/forklift-controller/pkg/lib/ref"
)

// Kind
var (
	VMKind      = libref.ToKind(VM{})
	ImageKind   = libref.ToKind(Image{})
	NetworkKind = libref.ToKind(Network{})
)

// Types.
type Tree = base.Tree
type TreeNode = base.TreeNode
type BranchNavigator = base.BranchNavigator
type ParentNavigator = base.ParentNavigator
