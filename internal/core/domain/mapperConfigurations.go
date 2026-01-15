package domain

type BGPConfigTemplateDate struct {
	Hostname  string
	ASN       string
	RouterID  string
	Neighbors []NeighborData
}

type NeighborData struct {
	NeighborIP string
	NeighborAS string
}
