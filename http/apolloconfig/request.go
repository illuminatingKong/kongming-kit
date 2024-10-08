package apolloconfig

type ReqReleaseBody struct {
	ReleaseTitle       string `json:"releaseTitle"`
	ReleaseComment     string `json:"releaseComment"`
	IsEmergencyPublish bool   `json:"isEmergencyPublish"`
}

type ReqAddPropertyItem struct {
	TableViewOperType  string `json:"tableViewOperType" validate:"required"`
	Key                string `json:"key" validate:"required" `
	Value              string `json:"value"`
	AddItemBtnDisabled bool   `json:"addItemBtnDisabled"`
}
