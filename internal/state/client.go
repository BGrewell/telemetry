package state

type ClientState struct {
	HostAddr string `json:"host_addr" yaml:"host_addr" xml:"host_addr" bson:"host_addr"`
	HostPort int    `json:"host_port" yaml:"host_port" xml:"host_port" bson:"host_port"`
}
