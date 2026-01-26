package manager

import (
	"context"
	"fmt"
	"sync"
	"time"

	loggergoUtil "github.com/Alonza0314/logger-go/v2/util"
	"github.com/free-ran-ue/free-ran-ue/v2/logger"
	"github.com/free-ran-ue/free-ran-ue/v2/model"
	"github.com/free-ran-ue/free-ran-ue/v2/ue"
	"github.com/free-ran-ue/frush/constant"
	"github.com/free-ran-ue/util"
)

type ueContext struct {
	ue     *ue.Ue
	imsi   string
	tunnel string
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
		tunnel: ueConfig.Ue.UeTunnelDevice,
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

func (c *ueContext) Ping(dn string) error {
	ranDataPlaneConn, ueIp := c.ue.GetRanDataPlaneConn(), c.ue.GetUeIp()

	receiveIcmpReplyChan := make(chan bool, 1)
	go func() {
		if err := ranDataPlaneConn.SetReadDeadline(time.Now().Add(5 * time.Second)); err != nil {
			receiveIcmpReplyChan <- false
			return
		}
		icmpEchoReply := make([]byte, 1024)
		if _, err := ranDataPlaneConn.Read(icmpEchoReply); err != nil {
			receiveIcmpReplyChan <- false
			return
		}
		if ok, err := util.IsIcmpEchoReply(icmpEchoReply, ueIp, dn); err != nil {
			receiveIcmpReplyChan <- false
			return
		} else if !ok {
			receiveIcmpReplyChan <- false
			return
		}
		receiveIcmpReplyChan <- true
	}()

	time.Sleep(300 * time.Millisecond)

	icmpEchoPacket, err := util.BuildIcmpEchoPacket(ueIp, dn)
	if err != nil {
		return fmt.Errorf("build icmp echo packet failed: %w", err)
	}

	if _, err = ranDataPlaneConn.Write(icmpEchoPacket); err != nil {
		return fmt.Errorf("send icmp echo packet failed: %w", err)
	}

	if !<-receiveIcmpReplyChan {
		return fmt.Errorf("ping %s failed", dn)
	}

	return nil
}
