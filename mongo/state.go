package mongo

const (
	StateStartup    = 0
	StatePrimary    = 1
	StateSecondary  = 2
	StateRecovering = 3
	StateStartup2   = 5
	StateUnknown    = 6
	StateArbiter    = 7
	StateDown       = 8
	StateRollback   = 9
	StateRemoved    = 10
)
