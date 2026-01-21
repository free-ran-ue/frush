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
	status constant.ContextStatus
}

func newGnbContext(gnbConfig model.GnbConfig) *gnbContext {
	logger := logger.NewGnbLogger(loggergoUtil.LogLevelString(gnbConfig.Logger.Level), "", true)
	return &gnbContext{
		gnb:    gnb.NewGnb(&gnbConfig, &logger),
		status: constant.Context_Stopped,
	}
}

func (c *gnbContext) getStatus() constant.ContextStatus {
	return c.status
}

func (c *gnbContext) setStatus(status constant.ContextStatus) {
	c.status = status
}

func (c *gnbContext) start(ctx context.Context) error {
	c.setStatus(constant.Context_Starting)

	if err := c.gnb.Start(ctx); err != nil {
		c.setStatus(constant.Context_Error)
		return err
	}
	c.setStatus(constant.Context_Running)

	return nil
}

func (c *gnbContext) stop() {
	c.setStatus(constant.Context_Stopping)

	c.gnb.Stop()

	c.setStatus(constant.Context_Stopped)
}
