package docker

import "time"

type dockerInspect struct {
	ID              string          `json:"Id"`
	Created         time.Time       `json:"Created"`
	Path            string          `json:"Path"`
	Args            []string        `json:"Args"`
	State           State           `json:"State"`
	Image           string          `json:"Image"`
	ResolvConfPath  string          `json:"ResolvConfPath"`
	HostnamePath    string          `json:"HostnamePath"`
	HostsPath       string          `json:"HostsPath"`
	LogPath         string          `json:"LogPath"`
	Name            string          `json:"Name"`
	RestartCount    int             `json:"RestartCount"`
	Driver          string          `json:"Driver"`
	Platform        string          `json:"Platform"`
	MountLabel      string          `json:"MountLabel"`
	ProcessLabel    string          `json:"ProcessLabel"`
	AppArmorProfile string          `json:"AppArmorProfile"`
	ExecIDs         any             `json:"ExecIDs"`
	HostConfig      HostConfig      `json:"HostConfig"`
	GraphDriver     GraphDriver     `json:"GraphDriver"`
	Mounts          []Mounts        `json:"Mounts"`
	Config          Config          `json:"Config"`
	NetworkSettings NetworkSettings `json:"NetworkSettings"`
}
type State struct {
	Status     string    `json:"Status"`
	Running    bool      `json:"Running"`
	Paused     bool      `json:"Paused"`
	Restarting bool      `json:"Restarting"`
	OOMKilled  bool      `json:"OOMKilled"`
	Dead       bool      `json:"Dead"`
	Pid        int       `json:"Pid"`
	ExitCode   int       `json:"ExitCode"`
	Error      string    `json:"Error"`
	StartedAt  time.Time `json:"StartedAt"`
	FinishedAt time.Time `json:"FinishedAt"`
}

type LogConfig struct {
	Type   string `json:"Type"`
	Config Config `json:"Config"`
}
type PortBindings struct {
}
type RestartPolicy struct {
	Name              string `json:"Name"`
	MaximumRetryCount int    `json:"MaximumRetryCount"`
}
type HostConfig struct {
	Binds                []string      `json:"Binds"`
	ContainerIDFile      string        `json:"ContainerIDFile"`
	LogConfig            LogConfig     `json:"LogConfig"`
	NetworkMode          string        `json:"NetworkMode"`
	PortBindings         PortBindings  `json:"PortBindings"`
	RestartPolicy        RestartPolicy `json:"RestartPolicy"`
	AutoRemove           bool          `json:"AutoRemove"`
	VolumeDriver         string        `json:"VolumeDriver"`
	VolumesFrom          any           `json:"VolumesFrom"`
	ConsoleSize          []int         `json:"ConsoleSize"`
	CapAdd               any           `json:"CapAdd"`
	CapDrop              any           `json:"CapDrop"`
	CgroupnsMode         string        `json:"CgroupnsMode"`
	DNS                  []any         `json:"Dns"`
	DNSOptions           []any         `json:"DnsOptions"`
	DNSSearch            []any         `json:"DnsSearch"`
	ExtraHosts           any           `json:"ExtraHosts"`
	GroupAdd             any           `json:"GroupAdd"`
	IpcMode              string        `json:"IpcMode"`
	Cgroup               string        `json:"Cgroup"`
	Links                any           `json:"Links"`
	OomScoreAdj          int           `json:"OomScoreAdj"`
	PidMode              string        `json:"PidMode"`
	Privileged           bool          `json:"Privileged"`
	PublishAllPorts      bool          `json:"PublishAllPorts"`
	ReadonlyRootfs       bool          `json:"ReadonlyRootfs"`
	SecurityOpt          any           `json:"SecurityOpt"`
	UTSMode              string        `json:"UTSMode"`
	UsernsMode           string        `json:"UsernsMode"`
	ShmSize              int           `json:"ShmSize"`
	Runtime              string        `json:"Runtime"`
	Isolation            string        `json:"Isolation"`
	CPUShares            int           `json:"CpuShares"`
	Memory               int           `json:"Memory"`
	NanoCpus             int           `json:"NanoCpus"`
	CgroupParent         string        `json:"CgroupParent"`
	BlkioWeight          int           `json:"BlkioWeight"`
	BlkioWeightDevice    []any         `json:"BlkioWeightDevice"`
	BlkioDeviceReadBps   []any         `json:"BlkioDeviceReadBps"`
	BlkioDeviceWriteBps  []any         `json:"BlkioDeviceWriteBps"`
	BlkioDeviceReadIOps  []any         `json:"BlkioDeviceReadIOps"`
	BlkioDeviceWriteIOps []any         `json:"BlkioDeviceWriteIOps"`
	CPUPeriod            int           `json:"CpuPeriod"`
	CPUQuota             int           `json:"CpuQuota"`
	CPURealtimePeriod    int           `json:"CpuRealtimePeriod"`
	CPURealtimeRuntime   int           `json:"CpuRealtimeRuntime"`
	CpusetCpus           string        `json:"CpusetCpus"`
	CpusetMems           string        `json:"CpusetMems"`
	Devices              []any         `json:"Devices"`
	DeviceCgroupRules    any           `json:"DeviceCgroupRules"`
	DeviceRequests       any           `json:"DeviceRequests"`
	MemoryReservation    int           `json:"MemoryReservation"`
	MemorySwap           int           `json:"MemorySwap"`
	MemorySwappiness     any           `json:"MemorySwappiness"`
	OomKillDisable       any           `json:"OomKillDisable"`
	PidsLimit            any           `json:"PidsLimit"`
	Ulimits              []any         `json:"Ulimits"`
	CPUCount             int           `json:"CpuCount"`
	CPUPercent           int           `json:"CpuPercent"`
	IOMaximumIOps        int           `json:"IOMaximumIOps"`
	IOMaximumBandwidth   int           `json:"IOMaximumBandwidth"`
	MaskedPaths          []string      `json:"MaskedPaths"`
	ReadonlyPaths        []string      `json:"ReadonlyPaths"`
}
type GraphDriver struct {
	Data any    `json:"Data"`
	Name string `json:"Name"`
}
type Mounts struct {
	Type        string `json:"Type"`
	Source      string `json:"Source"`
	Destination string `json:"Destination"`
	Mode        string `json:"Mode"`
	Rw          bool   `json:"RW"`
	Propagation string `json:"Propagation"`
}
type Eight0TCP struct {
}
type ExposedPorts struct {
	Eight0TCP Eight0TCP `json:"80/tcp"`
}
type Labels struct {
	Maintainer string `json:"maintainer"`
}
type Config struct {
	Hostname     string       `json:"Hostname"`
	Domainname   string       `json:"Domainname"`
	User         string       `json:"User"`
	AttachStdin  bool         `json:"AttachStdin"`
	AttachStdout bool         `json:"AttachStdout"`
	AttachStderr bool         `json:"AttachStderr"`
	ExposedPorts ExposedPorts `json:"ExposedPorts"`
	Tty          bool         `json:"Tty"`
	OpenStdin    bool         `json:"OpenStdin"`
	StdinOnce    bool         `json:"StdinOnce"`
	Env          []string     `json:"Env"`
	Cmd          []string     `json:"Cmd"`
	Image        string       `json:"Image"`
	Volumes      any          `json:"Volumes"`
	WorkingDir   string       `json:"WorkingDir"`
	Entrypoint   []string     `json:"Entrypoint"`
	OnBuild      any          `json:"OnBuild"`
	Labels       Labels       `json:"Labels"`
	StopSignal   string       `json:"StopSignal"`
}
type Ports struct {
	Eight0TCP any `json:"80/tcp"`
}
type Bridge struct {
	IPAMConfig          any    `json:"IPAMConfig"`
	Links               any    `json:"Links"`
	Aliases             any    `json:"Aliases"`
	MacAddress          string `json:"MacAddress"`
	DriverOpts          any    `json:"DriverOpts"`
	NetworkID           string `json:"NetworkID"`
	EndpointID          string `json:"EndpointID"`
	Gateway             string `json:"Gateway"`
	IPAddress           string `json:"IPAddress"`
	IPPrefixLen         int    `json:"IPPrefixLen"`
	IPv6Gateway         string `json:"IPv6Gateway"`
	GlobalIPv6Address   string `json:"GlobalIPv6Address"`
	GlobalIPv6PrefixLen int    `json:"GlobalIPv6PrefixLen"`
	DNSNames            any    `json:"DNSNames"`
}
type Networks struct {
	Bridge Bridge `json:"bridge"`
}
type NetworkSettings struct {
	Bridge                 string   `json:"Bridge"`
	SandboxID              string   `json:"SandboxID"`
	SandboxKey             string   `json:"SandboxKey"`
	Ports                  Ports    `json:"Ports"`
	HairpinMode            bool     `json:"HairpinMode"`
	LinkLocalIPv6Address   string   `json:"LinkLocalIPv6Address"`
	LinkLocalIPv6PrefixLen int      `json:"LinkLocalIPv6PrefixLen"`
	SecondaryIPAddresses   any      `json:"SecondaryIPAddresses"`
	SecondaryIPv6Addresses any      `json:"SecondaryIPv6Addresses"`
	EndpointID             string   `json:"EndpointID"`
	Gateway                string   `json:"Gateway"`
	GlobalIPv6Address      string   `json:"GlobalIPv6Address"`
	GlobalIPv6PrefixLen    int      `json:"GlobalIPv6PrefixLen"`
	IPAddress              string   `json:"IPAddress"`
	IPPrefixLen            int      `json:"IPPrefixLen"`
	IPv6Gateway            string   `json:"IPv6Gateway"`
	MacAddress             string   `json:"MacAddress"`
	Networks               Networks `json:"Networks"`
}
