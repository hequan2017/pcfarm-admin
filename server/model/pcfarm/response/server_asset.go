package response

import "github.com/flipped-aurora/gin-vue-admin/server/model/pcfarm"

type ServerAssetResponse struct {
	Server pcfarm.ServerAsset `json:"server"`
}

type ServerAssetListItem struct {
	pcfarm.ServerAsset
	HasBmcPassword bool `json:"hasBmcPassword"`
}
