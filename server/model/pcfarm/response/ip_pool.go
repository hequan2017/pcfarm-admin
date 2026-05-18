package response

import "github.com/flipped-aurora/gin-vue-admin/server/model/pcfarm"

type IPPoolSummary struct {
	Pool           pcfarm.IPPool `json:"pool"`
	AllocatedCount int64         `json:"allocatedCount"`
	AvailableCount int64         `json:"availableCount"`
}
