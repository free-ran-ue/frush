package manager

import (
	"context"
	"fmt"
	"sync"

	loggergoUtil "github.com/Alonza0314/logger-go/v2/util"
	"github.com/free-ran-ue/free-ran-ue/v2/logger"
	"github.com/free-ran-ue/free-ran-ue/v2/model"
	"github.com/free-ran-ue/free-ran-ue/v2/ue"
	"github.com/free-ran-ue/frush/constant"
)

type ueContext struct {
	ue     *ue.Ue
	imsi   string
	status constant.ContextStatus
	wg     *sync.WaitGroup
	ctx    context.Context
	cancel context.CancelFunc
}

func newUeContext(ueConfig model.UeConfig) *ueContext {
	logger := logger.NewUeLogger(loggergoUtil.LogLevelString(ueConfig.Logger.Level), "", true)
	return &ueContext{
		ue:     ue.NewUe(&ueConfig, &logger),
		imsi:   fmt.Sprintf("imsi-%s%s%s", ueConfig.Ue.PlmnId.Mcc, ueConfig.Ue.PlmnId.Mnc, ueConfig.Ue.Msin),
		status: constant.CONTEXT_STATUS_UE_STOPPED,
		wg:     &sync.WaitGroup{},
		ctx:    nil,
		cancel: nil,
	}
}

func (c *ueContext) GetImsi() string {
	return c.imsi
}

func (c *ueContext) GetStatus() constant.ContextStatus {
	return c.status
}

func (c *ueContext) SetStatus(status constant.ContextStatus) {
	c.status = status
}

func (c *ueContext) Start(ctx context.Context) error {
	c.SetStatus(constant.CONTEXT_STATUS_UE_REGISTERING)

	c.ctx, c.cancel = context.WithCancel(ctx)

	if err := c.ue.Start(c.ctx, c.wg); err != nil {
		c.SetStatus(constant.CONTEXT_STATUS_UE_ERROR)
		return err
	}
	c.SetStatus(constant.CONTEXT_STATUS_UE_REGISTERED)

	return nil
}

func (c *ueContext) Stop() {
	c.SetStatus(constant.CONTEXT_STATUS_UE_DE_REGISTERING)

	c.cancel()
	c.wg.Wait()

	c.ue.Stop()

	c.SetStatus(constant.CONTEXT_STATUS_UE_DE_REGISTERED)
}
