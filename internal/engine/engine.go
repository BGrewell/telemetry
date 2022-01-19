package engine

import (
	"fmt"
	"github.com/BGrewell/tgams/internal/grpc"
	"github.com/BGrewell/tgams/internal/timesync"
	"time"
)

func NewCoreEngine(host string, port int, timeout int) (engine *CoreEngine, err error) {
	ctrl := &grpc.ControlClient{}
	err = ctrl.Connect(fmt.Sprintf("%s:%d", host, port), timeout)
	if err != nil {
		return nil, err
	}

	clk := &timesync.SoftwareClock{
		TimeSyncFunc: ctrl.SendTimeSyncRequest,
		Interval:     5 * time.Second,
		Bursts:       5,
	}
	clk.EnableCompensation()
	clk.Enable()

	ce := &CoreEngine{
		clock:   clk,
		control: ctrl,
		running: true,
	}
	return ce, nil
}

type CoreEngine struct {
	clock   timesync.Clock
	control *grpc.ControlClient
	running bool
}

func (ce *CoreEngine) EnableCompensatedTime() error {
	ce.clock.EnableCompensation()
	return nil
}

func (ce *CoreEngine) DisableCompensatedTime() error {
	ce.clock.DisableCompensation()
	return nil
}

func (ce *CoreEngine) IsRunning() bool {
	return ce.running
}
