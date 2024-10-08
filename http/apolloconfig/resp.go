package apolloconfig

type RespApp struct {
	Name                       string `json:"name"`
	AppId                      string `json:"appId"`
	OrgId                      string `json:"orgId"`
	OrgName                    string `json:"orgName"`
	OwnerName                  string `json:"ownerName"`
	OwnerEmail                 string `json:"ownerEmail"`
	Id                         int    `json:"id"`
	IsDeleted                  bool   `json:"isDeleted"`
	DeletedAt                  int    `json:"deletedAt"`
	DataChangeCreatedBy        string `json:"dataChangeCreatedBy"`
	DataChangeCreatedTime      string `json:"dataChangeCreatedTime"`
	DataChangeLastModifiedBy   string `json:"dataChangeLastModifiedBy"`
	DataChangeLastModifiedTime string `json:"dataChangeLastModifiedTime"`
}

type RespClusters struct {
	Code     int `json:"code"`
	Entities []struct {
		Code    int    `json:"code"`
		Message string `json:"message"`
		Body    struct {
			Env      string `json:"env"`
			Clusters []struct {
				Id                         int    `json:"id"`
				Name                       string `json:"name"`
				AppId                      string `json:"appId"`
				ParentClusterId            int    `json:"parentClusterId"`
				DataChangeCreatedBy        string `json:"dataChangeCreatedBy"`
				DataChangeLastModifiedBy   string `json:"dataChangeLastModifiedBy"`
				DataChangeCreatedTime      string `json:"dataChangeCreatedTime"`
				DataChangeLastModifiedTime string `json:"dataChangeLastModifiedTime"`
			} `json:"clusters"`
		} `json:"body"`
	} `json:"entities"`
}

// http://localhost:8070/apps/<appid>/envs/<env>/clusters/<cluster>/namespaces
type RespNamespaces struct {
	BaseInfo struct {
		Id                         int    `json:"id"`
		AppId                      string `json:"appId"`
		ClusterName                string `json:"clusterName"`
		NamespaceName              string `json:"namespaceName"`
		DataChangeCreatedBy        string `json:"dataChangeCreatedBy"`
		DataChangeLastModifiedBy   string `json:"dataChangeLastModifiedBy"`
		DataChangeCreatedTime      string `json:"dataChangeCreatedTime"`
		DataChangeLastModifiedTime string `json:"dataChangeLastModifiedTime"`
	} `json:"baseInfo"`
	ItemModifiedCnt int `json:"itemModifiedCnt"` // 修改的数量, 0: 未修改, >0: 修改的数量，未发布
	Items           []struct {
		Item struct {
			Id                                  int    `json:"id"`
			NamespaceId                         int    `json:"namespaceId"`
			Key                                 string `json:"key"`
			Type                                int    `json:"type"`
			Value                               string `json:"value"`
			Comment                             string `json:"comment"`
			LineNum                             int    `json:"lineNum"`
			DataChangeCreatedBy                 string `json:"dataChangeCreatedBy"`
			DataChangeLastModifiedBy            string `json:"dataChangeLastModifiedBy"`
			DataChangeLastModifiedByDisplayName string `json:"dataChangeLastModifiedByDisplayName"`
			DataChangeCreatedTime               string `json:"dataChangeCreatedTime"`
			DataChangeLastModifiedTime          string `json:"dataChangeLastModifiedTime"`
			DataChangeCreatedByDisplayName      string `json:"dataChangeCreatedByDisplayName,omitempty"`
		} `json:"item"`
		IsModified bool `json:"isModified"` // 是否修改, true: 修改完成未提交发布
		IsDeleted  bool `json:"isDeleted"`  // true: 删除未提交发布
	} `json:"items"`
	Format         string `json:"format"`
	IsPublic       bool   `json:"isPublic"`
	ParentAppId    string `json:"parentAppId"`
	Comment        string `json:"comment"`
	IsConfigHidden bool   `json:"isConfigHidden"`
}

type RespPublishedRelease struct {
	Id                         int    `json:"id"`
	ReleaseKey                 string `json:"releaseKey"`
	Name                       string `json:"name"`
	AppId                      string `json:"appId"`
	ClusterName                string `json:"clusterName"`
	NamespaceName              string `json:"namespaceName"`
	Configurations             string `json:"configurations"`
	Comment                    string `json:"comment"`
	IsAbandoned                bool   `json:"isAbandoned"`
	DataChangeCreatedBy        string `json:"dataChangeCreatedBy"`
	DataChangeLastModifiedBy   string `json:"dataChangeLastModifiedBy"`
	DataChangeCreatedTime      string `json:"dataChangeCreatedTime"`
	DataChangeLastModifiedTime string `json:"dataChangeLastModifiedTime"`
}

type RespHistoryRelease struct {
	Id                   int    `json:"id"`
	AppId                string `json:"appId"`
	ClusterName          string `json:"clusterName"`
	NamespaceName        string `json:"namespaceName"`
	BranchName           string `json:"branchName"`
	Operator             string `json:"operator"`
	OperatorDisplayName  string `json:"operatorDisplayName"`
	ReleaseId            int    `json:"releaseId"`
	ReleaseTitle         string `json:"releaseTitle"`
	ReleaseComment       string `json:"releaseComment"`
	ReleaseTime          string `json:"releaseTime"`
	ReleaseTimeFormatted string `json:"releaseTimeFormatted"`
	Configuration        []struct {
		FirstEntity  string `json:"firstEntity"`
		SecondEntity string `json:"secondEntity"`
	} `json:"configuration"`
	IsReleaseAbandoned bool `json:"isReleaseAbandoned"`
	PreviousReleaseId  int  `json:"previousReleaseId"`
	Operation          int  `json:"operation"`
	OperationContext   struct {
		IsEmergencyPublish bool `json:"isEmergencyPublish"`
	} `json:"operationContext"`
}

type RespItemChangeBody struct {
	Id                         int    `json:"id"`
	NamespaceId                int    `json:"namespaceId"`
	Key                        string `json:"key"`
	Value                      string `json:"value"`
	LineNum                    int    `json:"lineNum"`
	DataChangeCreatedBy        string `json:"dataChangeCreatedBy"`
	DataChangeLastModifiedBy   string `json:"dataChangeLastModifiedBy"`
	DataChangeCreatedTime      string `json:"dataChangeCreatedTime"`
	DataChangeLastModifiedTime string `json:"dataChangeLastModifiedTime"`
}
