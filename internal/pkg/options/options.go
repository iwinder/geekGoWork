package configs

type Option struct {
	ServerOption *ServerOption `yaml:"server" mapstructure:"server"`
	MysqlOption  *MysqlOption  `yaml:"mysql" mapstructure:"mysql"`
}

type ServerOption struct {
	GRpcServerOption *GRpcServerOption `yaml:"grpc" mapstructure:"grpc"`
	HttpServerOption *HttpServerOption `yaml:"http" mapstructure:"http"`
}
