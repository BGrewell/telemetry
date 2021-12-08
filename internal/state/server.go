package state

type ServerState struct {
	ListenAddr string `json:"listen_addr" yaml:"listen_addr" xml:"listen_addr" bson:"listen_addr"`
	ListenPort int    `json:"listen_port" yaml:"listen_port" xml:"listen_port" bson:"listen_port"`
}
