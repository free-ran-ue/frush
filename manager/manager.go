package manager

type manager struct {
	gnbContext *gnbContext
	ueContext *ueContext
}

func NewManager() *manager {
	return &manager{
		gnbContext: NewGnbContext(),
		ueContext: NewUeContext(),
	}
}