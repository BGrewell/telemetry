package timesync

import (
	api "github.com/BGrewell/tgams/api/go"
)

func Offset(t1, t2, t3, t4 int64) int64 {
	return ((t2 - t1) + (t3 - t4)) / 2
}

func Delay(t1, t2, t3, t4 int64) int64 {
	return (t4 - t1) - (t3 - t2)
}

func CalcOffset(response *api.TimeSyncResponse) int64 {
    return Offset(response.T1, response.T2, response.T3, response.T4)
}

func CalcDelay(response *api.TimeSyncResponse) int64 {
    return Delay(response.T1, response.T2, response.T3, response.T4)
}
