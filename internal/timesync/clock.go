package timesync

//TODO: Implement internal clock that can use timesync probes to bring clocks into relative sync so that we can do
//      one-way delay measurements.

type Clock interface {
	Now() int64
	EnableCompensation()
	DisableCompensation()
}
