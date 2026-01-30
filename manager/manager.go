package manager

import (
	"context"

	"github.com/free-ran-ue/free-ran-ue/v2/model"
)

var (
	Manager    *manager
	RootCtx    context.Context
	RootCancel context.CancelFunc
)

type manager struct {
	gnbContext *gnbContext
	ueContext  *ueContext
}

func NewManager(gnbConfig model.GnbConfig, ueConfig model.UeConfig) *manager {
	return &manager{
		gnbContext: newGnbContext(gnbConfig),
		ueContext:  newUeContext(ueConfig),
	}
}

func (m *manager) GnbContext() *gnbContext {
	return m.gnbContext
}

func (m *manager) UeContext() *ueContext {
	return m.ueContext
}
