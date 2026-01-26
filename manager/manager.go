package manager

import (
	"sync"

	"github.com/free-ran-ue/free-ran-ue/v2/model"
)

type manager struct {
	gnbContext *gnbContext
	ueContext  *ueContext
}

func NewManager(gnbConfig model.GnbConfig, ueConfig model.UeConfig, managerWg *sync.WaitGroup) *manager {
	return &manager{
		gnbContext: newGnbContext(gnbConfig, managerWg),
		ueContext:  newUeContext(ueConfig, managerWg),
	}
}

func (m *manager) GnbContext() *gnbContext {
	return m.gnbContext
}

func (m *manager) UeContext() *ueContext {
	return m.ueContext
}
