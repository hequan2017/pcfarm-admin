package request

import "github.com/flipped-aurora/gin-vue-admin/server/model/common/request"

type IPPoolSearch struct {
	request.PageInfo
}

type CreateIPPool struct {
	Name      string `json:"name" form:"name"`
	CIDR      string `json:"cidr" form:"cidr"`
	StartIP   string `json:"startIp" form:"startIp"`
	EndIP     string `json:"endIp" form:"endIp"`
	Gateway   string `json:"gateway" form:"gateway"`
	DNS       string `json:"dns" form:"dns"`
	BindIface string `json:"bindIface" form:"bindIface"`
	Enabled   bool   `json:"enabled" form:"enabled"`
}
