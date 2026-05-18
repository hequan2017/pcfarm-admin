package request

type AgentRegisterRequest struct {
	SerialNumber    string `json:"serialNumber" form:"serialNumber"`
	PxeMac          string `json:"pxeMac" form:"pxeMac"`
	IP              string `json:"ip" form:"ip"`
	AgentVersion    string `json:"agentVersion" form:"agentVersion"`
	HardwareSummary string `json:"hardwareSummary" form:"hardwareSummary"`
	Token           string `json:"token" form:"token"`
}

type AgentHeartbeatRequest struct {
	SerialNumber string `json:"serialNumber" form:"serialNumber"`
	PxeMac       string `json:"pxeMac" form:"pxeMac"`
	Token        string `json:"token" form:"token"`
}
