package atlas

type Project struct {
	Id string
}

type AutoScale struct {
	MinInstanceSize string `json:"minInstanceSize,omitempty"`
	MaxInstanceSize string `json:"maxInstanceSize,omitempty"`
}

type ElectableSpecs struct {
	InstanceSize  string     `json:"instanceSize,omitempty"`
	DiskIOPS      int        `json:"diskIOPS,omitempty"`
	EBSVolumeType string     `json:"ebsVolumeType,omitempty"`
	NodeCount     int        `json:"nodeCount,omitempty"`
	AutoScale     *AutoScale `json:"autoScale,omitempty"`
}

func (e *ElectableSpecs) GetInstanceSize() string {
	return e.InstanceSize
}

func (e *ElectableSpecs) SetInstanceSize(size string) {
	e.InstanceSize = size
}

func (e *ElectableSpecs) GetDiskIOPS() int {
	return e.DiskIOPS
}

func (e *ElectableSpecs) SetDiskIOPS(diskIOPS int) {
	e.DiskIOPS = diskIOPS
}

func (e *ElectableSpecs) GetEbsVolumeType() string {
	return e.EBSVolumeType
}

func (e *ElectableSpecs) SetEbsVolumeType(ebsVolumeType string) {
	e.EBSVolumeType = ebsVolumeType
}

func (e *ElectableSpecs) GetNodeCount() int {
	return e.NodeCount
}

func (e *ElectableSpecs) SetNodeCount(nodeCount int) {
	e.NodeCount = nodeCount
}

type ReadOnlySpecs struct {
	InstanceSize  string `json:"instanceSize,omitempty"`
	DiskIOPS      int    `json:"diskIOPS,omitempty"`
	EBSVolumeType string `json:"ebsVolumeType,omitempty"`
	NodeCount     int    `json:"nodeCount,omitempty"`
}

func (r *ReadOnlySpecs) GetInstanceSize() string {
	return r.InstanceSize
}

func (r *ReadOnlySpecs) SetInstanceSize(size string) {
	r.InstanceSize = size
}

func (r *ReadOnlySpecs) GetDiskIOPS() int {
	return r.DiskIOPS
}

func (r *ReadOnlySpecs) SetDiskIOPS(diskIOPS int) {
	r.DiskIOPS = diskIOPS
}

func (r *ReadOnlySpecs) GetEbsVolumeType() string {
	return r.EBSVolumeType
}

func (r *ReadOnlySpecs) SetEbsVolumeType(ebsVolumeType string) {
	r.EBSVolumeType = ebsVolumeType
}

func (r *ReadOnlySpecs) GetNodeCount() int {
	return r.NodeCount
}

func (r *ReadOnlySpecs) SetNodeCount(nodeCount int) {
	r.NodeCount = nodeCount
}

type AnalyticsSpecs struct {
	InstanceSize  string     `json:"instanceSize,omitempty"`
	DiskIOPS      int        `json:"diskIOPS,omitempty"`
	EBSVolumeType string     `json:"ebsVolumeType,omitempty"`
	NodeCount     int        `json:"nodeCount,omitempty"`
	AutoScale     *AutoScale `json:"autoScale,omitempty"`
}

func (a *AnalyticsSpecs) GetInstanceSize() string {
	return a.InstanceSize
}

func (a *AnalyticsSpecs) SetInstanceSize(size string) {
	a.InstanceSize = size
}

func (a *AnalyticsSpecs) GetDiskIOPS() int {
	return a.DiskIOPS
}

func (a *AnalyticsSpecs) SetDiskIOPS(diskIOPS int) {
	a.DiskIOPS = diskIOPS
}

func (a *AnalyticsSpecs) GetEbsVolumeType() string {
	return a.EBSVolumeType
}

func (a *AnalyticsSpecs) SetEbsVolumeType(ebsVolumeType string) {
	a.EBSVolumeType = ebsVolumeType
}

func (a *AnalyticsSpecs) GetNodeCount() int {
	return a.NodeCount
}

func (a *AnalyticsSpecs) SetNodeCount(nodeCount int) {
	a.NodeCount = nodeCount
}

type InstanceSizer interface {
	GetInstanceSize() string
	SetInstanceSize(size string)
	GetDiskIOPS() int
	SetDiskIOPS(diskIOPS int)
	GetEbsVolumeType() string
	SetEbsVolumeType(ebsVolumeType string)
	GetNodeCount() int
	SetNodeCount(nodeCount int)
}

type MyEvent struct {
	Project        string          `json:"project"`
	Cluster        string          `json:"cluster"`
	ReadOnlySpecs  *ReadOnlySpecs  `json:"readOnlySpecs,omitempty"`
	ElectableSpecs *ElectableSpecs `json:"electableSpecs,omitempty"`
	AnalyticsSpecs *AnalyticsSpecs `json:"analyticsSpecs,omitempty"`
}
