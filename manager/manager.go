package manager

import (
	"context"

	"github.com/free-ran-ue/free-ran-ue/v2/model"
	"github.com/free-ran-ue/frush/constant"
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

func (m *manager) GnbStart(ctx context.Context) error {
	return m.gnbContext.start(ctx)
}

func (m *manager) GnbStop() {
	m.gnbContext.stop()
}

func (m *manager) GnbStatus() constant.ContextStatus {
	return m.gnbContext.getStatus()
}
