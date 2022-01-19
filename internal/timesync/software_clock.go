package timesync

import (
	"github.com/BGrewell/go-conversions"
	log "github.com/BGrewell/tgams/internal/logging"
	"time"
)

type SoftwareClock struct {
	TimeSyncFunc func(t1 int64) (t2, t3 int64, err error)
	Interval     time.Duration
	Bursts       int
	offset       int64
	delay        int64
	enabled      bool
	compEnabled  bool
}

func (c *SoftwareClock) EnableCompensation() {
	c.compEnabled = true
}

func (c *SoftwareClock) DisableCompensation() {
	c.compEnabled = false
}

func (c *SoftwareClock) Enable() {
	c.enabled = true
	go c.update()
}

func (c *SoftwareClock) Disable() {
	c.enabled = false
}

func (c *SoftwareClock) Enabled() bool {
	return c.enabled
}

func (c *SoftwareClock) Now() int64 {
	if !c.compEnabled {
		return time.Now().UnixNano()
	}

	return time.Now().UnixNano() + c.offset
}

func (c *SoftwareClock) update() {
	for c.enabled {
		// TODO: enable burst measurements
		t1 := c.Now()
		t2, t3, err := c.TimeSyncFunc(t1)
		if err != nil {
			log.Error("Error syncing time: %v", err)
			time.Sleep(c.Interval)
			continue
		}
		t4 := c.Now()
		offset := Offset(t1, t2, t3, t4)
		delay := Delay(t1, t2, t3, t4)
		c.offset += offset
		c.delay += delay
		log.TraceWithFields(map[string]interface{}{
			"offset": conversions.ConvertNanosecondsToStringTime(offset),
			"delay":  conversions.ConvertNanosecondsToStringTime(delay)}, "clock updated")
		time.Sleep(c.Interval)

	}
}
