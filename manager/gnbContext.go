package manager

import (
	"context"
	"time"

	loggergoUtil "github.com/Alonza0314/logger-go/v2/util"
	"github.com/free-ran-ue/free-ran-ue/v2/gnb"
	"github.com/free-ran-ue/free-ran-ue/v2/logger"
	"github.com/free-ran-ue/free-ran-ue/v2/model"
	"github.com/free-ran-ue/frush/constant"
)

type gnbContext struct {
	gnb       *gnb.Gnb
	gnbConfig model.GnbConfig
	name      string
	status    constant.ContextStatus
	ctx       context.Context
	cancel    context.CancelFunc
}

func newGnbContext(gnbConfig model.GnbConfig) *gnbContext {
	logger := logger.NewGnbLogger(loggergoUtil.LogLevelString(gnbConfig.Logger.Level), "", true)
	return &gnbContext{
		gnb:       gnb.NewGnb(&gnbConfig, &logger),
		gnbConfig: gnbConfig,
		name:      gnbConfig.Gnb.GnbName,
		status:    constant.CONTEXT_STATUS_GNB_STOPPED,
		ctx:       nil,
		cancel:    nil,
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

func (c *gnbContext) GetContext() context.Context {
	return c.ctx
}

func (c *gnbContext) Start(ctx context.Context) error {
	c.SetStatus(constant.CONTEXT_STATUS_GNB_STARTING)

	// Create a new context for this gNB instance
	c.ctx, c.cancel = context.WithCancel(ctx)

	// Create a fresh gNB instance
	logger := logger.NewGnbLogger(loggergoUtil.LogLevelString(c.gnbConfig.Logger.Level), "", true)
	c.gnb = gnb.NewGnb(&c.gnbConfig, &logger)

	if err := c.gnb.Start(c.ctx); err != nil {
		c.SetStatus(constant.CONTEXT_STATUS_GNB_ERROR)
		return err
	}
	c.SetStatus(constant.CONTEXT_STATUS_GNB_RUNNING)

	return nil
}

func (c *gnbContext) Stop() error {
	c.SetStatus(constant.CONTEXT_STATUS_GNB_STOPPING)

	// Cancel context first to signal goroutines to stop
	if c.cancel != nil {
		c.cancel()
	}

	c.gnb.Stop()

	// Wait for goroutines to fully exit after connections are closed
	time.Sleep(200 * time.Millisecond)

	// Clear references
	c.gnb = nil
	c.cancel = nil
	c.ctx = nil

	c.SetStatus(constant.CONTEXT_STATUS_GNB_STOPPED)

	return nil
}
