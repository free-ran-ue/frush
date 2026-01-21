package manager

import (
	"context"

	loggergoUtil "github.com/Alonza0314/logger-go/v2/util"
	"github.com/free-ran-ue/free-ran-ue/v2/gnb"
	"github.com/free-ran-ue/free-ran-ue/v2/logger"
	"github.com/free-ran-ue/free-ran-ue/v2/model"
	"github.com/free-ran-ue/frush/constant"
)

type gnbContext struct {
	gnb    *gnb.Gnb
	name   string
	status constant.ContextStatus
}

func newGnbContext(gnbConfig model.GnbConfig) *gnbContext {
	logger := logger.NewGnbLogger(loggergoUtil.LogLevelString(gnbConfig.Logger.Level), "", true)
	return &gnbContext{
		gnb:    gnb.NewGnb(&gnbConfig, &logger),
		name:   gnbConfig.Gnb.GnbName,
		status: constant.CONTEXT_STATUS_GNB_STOPPED,
	}
}

func (c *gnbContext) GetName() string {
	return c.name
}

func (c *gnbContext) GetStatus() constant.ContextStatus {
	return c.status
}

func (c *gnbContext) SetStatus(status constant.ContextStatus) {
	c.status = status
}

func (c *gnbContext) Start(ctx context.Context) error {
	c.SetStatus(constant.CONTEXT_STATUS_GNB_STARTING)

	if err := c.gnb.Start(ctx); err != nil {
		c.SetStatus(constant.CONTEXT_STATUS_GNB_ERROR)
		return err
	}
	c.SetStatus(constant.CONTEXT_STATUS_GNB_RUNNING)

	return nil
}

func (c *gnbContext) Stop() {
	c.SetStatus(constant.CONTEXT_STATUS_GNB_STOPPING)

	c.gnb.Stop()

	c.SetStatus(constant.CONTEXT_STATUS_GNB_STOPPED)
}
