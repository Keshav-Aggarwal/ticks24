package config

type Config struct {
	AppName        string        `json:"AppName"`
	AppToken       string        `json:"AppToken"`
	MachineIp      string        `json:"MachineIp"`
	ConfigFilePath string        `json:"ConfigFilePath"`
	FrontEnd       FrontEnd      `json:"FrontEnd"`
	WebServer      WebServer     `json:"WebServer"`
	AuthService    AuthService   `json:"AuthService"`
	LoginService   AuthService   `json:"LoginService"`
	AppDatabases   []AppDatabase `json:"AppDatabases"`
	LogConfig      LogConfig     `json:"LogConfig"`
}

type AuthService struct {
	Ip     string `json:"Ip"`
	Port   int32  `json:"Port"`
	IsTcp  bool   `json:"IsTcp"`
	IsHttp bool   `json:"IsHttp"`
}

type WebServer struct {
	Ip         string `json:"Ip"`
	Port       int32  `json:"Port"`
	StopUrl    string `json:"StopUrl"`
	RestartUrl string `json:"RestartUrl"`
	AuthKey    string `json:"AuthKey"`
	Mode       string `json:"Mode"`
}

type AppDatabase struct {
	Ip           string `json:"Ip"`
	Port         int32  `json:"Port"`
	DatabaseName string `json:"DatabaseName"`
	MaxBatchSize int32  `json:"MaxBatchSize"`
}

type FrontEnd struct {
	ViewsPath              string `json:"ViewsPath"`
	TemplatesPath          string `json:"TemplatesPath"`
	TemplateDelimiterStart string `json:"TemplateDelimiterStart"`
	TemplateDelimiterEnd   string `json:"TemplateDelimiterEnd"`
}

type LogConfig struct {
	Level string `json:"Level"`
	Path  string `json:"Path"`
	Days  int32  `json:"Days"`
}
