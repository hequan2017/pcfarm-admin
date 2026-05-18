package request

import (
	"github.com/flipped-aurora/gin-vue-admin/server/model/common/request"
)

type ServerAssetSearch struct {
	Keyword string `json:"keyword" form:"keyword"`
	Status  string `json:"status" form:"status"`
	request.PageInfo
}

type CreateServerAsset struct {
	AssetCode     string `json:"assetCode" form:"assetCode"`
	SerialNumber  string `json:"serialNumber" form:"serialNumber"`
	PxeMac        string `json:"pxeMac" form:"pxeMac"`
	BmcAddress    string `json:"bmcAddress" form:"bmcAddress"`
	BmcUsername   string `json:"bmcUsername" form:"bmcUsername"`
	BmcPassword   string `json:"bmcPassword" form:"bmcPassword"`
	PowerProtocol string `json:"powerProtocol" form:"powerProtocol"`
}

type UpdateBootPolicy struct {
	ID         uint   `json:"id" form:"id"`
	BootPolicy string `json:"bootPolicy" form:"bootPolicy"`
}

type PowerActionRequest struct {
	ID     uint   `json:"id" form:"id"`
	Action string `json:"action" form:"action"`
}
