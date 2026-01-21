package manager

import "github.com/free-ran-ue/free-ran-ue/v2/model"

type ueContext struct {
}

func newUeContext(ueConfig model.UeConfig) *ueContext {
	return &ueContext{}
}
