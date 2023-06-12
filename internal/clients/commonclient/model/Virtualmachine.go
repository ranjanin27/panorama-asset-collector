package model

type VirtualMachine struct {
	AllowedOperations []string `json:"allowedOperations"`
	AppInfo           struct {
		Vmware struct {
			DatacenterInfo struct {
				ID    string `json:"id"`
				Moref string `json:"moref"`
				Name  string `json:"name"`
			} `json:"datacenterInfo"`
			DatastoresInfo []struct {
				DisplayName string `json:"displayName"`
				ID          string `json:"id"`
				Name        string `json:"name"`
				ResourceURI string `json:"resourceUri"`
				Type        string `json:"type"`
			} `json:"datastoresInfo"`
			Moref            string `json:"moref"`
			ResourcePoolInfo struct {
				DisplayName string `json:"displayName"`
				ID          string `json:"id"`
				Moref       string `json:"moref"`
				Name        string `json:"name"`
				ResourceURI string `json:"resourceUri"`
				Type        string `json:"type"`
			} `json:"resourcePoolInfo"`
			ToolsInfo struct {
				Status  string `json:"status"`
				Type    string `json:"type"`
				Version string `json:"version"`
			} `json:"toolsInfo"`
			Type string `json:"type"`
		} `json:"vmware"`
	} `json:"appInfo"`
	AppType         string `json:"appType"`
	CapacityInBytes int64  `json:"capacityInBytes"`
	ClusterInfo     struct {
		DisplayName string `json:"displayName"`
		ID          string `json:"id"`
		Name        string `json:"name"`
		ResourceURI string `json:"resourceUri"`
		Type        string `json:"type"`
	} `json:"clusterInfo"`
	ComputeInfo struct {
		MemorySizeInMib string `json:"memorySizeInMib"`
		NumCPUCores     int    `json:"numCpuCores"`
		NumCPUThreads   int    `json:"numCpuThreads"`
	} `json:"computeInfo"`
	CreatedAt   string `json:"createdAt"`
	CustomerID  string `json:"customerId"`
	DisplayName string `json:"displayName"`
	FolderInfo  struct {
		DisplayName string `json:"displayName"`
		ID          string `json:"id"`
		Name        string `json:"name"`
		ResourceURI string `json:"resourceUri"`
		Type        string `json:"type"`
	} `json:"folderInfo"`
	Generation int `json:"generation"`
	GuestInfo  struct {
		BuildVersion   string `json:"buildVersion"`
		Name           string `json:"name"`
		ReleaseVersion string `json:"releaseVersion"`
		Type           string `json:"type"`
	} `json:"guestInfo"`
	HciClusterUUID string `json:"hciClusterUuid"`
	HostInfo       struct {
		DisplayName string `json:"displayName"`
		ID          string `json:"id"`
		Name        string `json:"name"`
		ResourceURI string `json:"resourceUri"`
		Type        string `json:"type"`
	} `json:"hostInfo"`
	HypervisorManagerInfo struct {
		DisplayName string `json:"displayName"`
		ID          string `json:"id"`
		Name        string `json:"name"`
		ResourceURI string `json:"resourceUri"`
		Type        string `json:"type"`
	} `json:"hypervisorManagerInfo"`
	ID                string      `json:"id"`
	Name              string      `json:"name"`
	NetworkAdapters   interface{} `json:"networkAdapters"`
	NetworkAddress    string      `json:"networkAddress"`
	PowerState        string      `json:"powerState"`
	Protected         bool        `json:"protected"`
	ProtectionJobInfo struct {
		ID                   string `json:"id"`
		ProtectionPolicyInfo struct {
			ID          string `json:"id"`
			Name        string `json:"name"`
			ResourceURI string `json:"resourceUri"`
			Type        string `json:"type"`
		} `json:"protectionPolicyInfo"`
		ResourceURI string `json:"resourceUri"`
		Type        string `json:"type"`
	} `json:"protectionJobInfo"`
	ProtectionPolicyAppliedAtInfo struct {
		ID          string `json:"id"`
		Name        string `json:"name"`
		ResourceURI string `json:"resourceUri"`
		Type        string `json:"type"`
	} `json:"protectionPolicyAppliedAtInfo"`
	ProtectionStatus    string   `json:"protectionStatus"`
	RecoveryPointsExist bool     `json:"recoveryPointsExist"`
	ResourceURI         string   `json:"resourceUri"`
	Services            []string `json:"services"`
	State               string   `json:"state"`
	StateReason         string   `json:"stateReason"`
	Status              string   `json:"status"`
	Type                string   `json:"type"`
	UID                 string   `json:"uid"`
	UpdatedAt           string   `json:"updatedAt"`
	VclsVM              bool     `json:"vclsVm"`
	VirtualDisks        []struct {
		AppInfo struct {
			Vmware struct {
				DatastoreInfo struct {
					DisplayName string `json:"displayName"`
					ID          string `json:"id"`
					Name        string `json:"name"`
					ResourceURI string `json:"resourceUri"`
					Type        string `json:"type"`
				} `json:"datastoreInfo"`
				DiskUUIDEnabled bool   `json:"diskUuidEnabled"`
				Type            string `json:"type"`
			} `json:"vmware"`
		} `json:"appInfo"`
		CapacityInBytes int64  `json:"capacityInBytes"`
		FilePath        string `json:"filePath"`
		ID              string `json:"id"`
		Name            string `json:"name"`
		UID             string `json:"uid"`
	} `json:"virtualDisks"`
	VMClassification string `json:"vmClassification,omitempty"`
	VMConfigPath     string `json:"vmConfigPath"`
	VMPerfMetricInfo struct {
		AverageReadLatency   int `json:"averageReadLatency"`
		AverageWriteLatency  int `json:"averageWriteLatency"`
		CPUAllocatedInMhz    int `json:"cpuAllocatedInMhz"`
		CPUUsedInMhz         int `json:"cpuUsedInMhz"`
		MemoryAllocatedInMb  int `json:"memoryAllocatedInMb"`
		MemoryUsedInMb       int `json:"memoryUsedInMb"`
		StorageAllocatedInKb int `json:"storageAllocatedInKb"`
		StorageUsedInBytes   int `json:"storageUsedInBytes"`
		TotalReadIops        int `json:"totalReadIops"`
		TotalWriteIops       int `json:"totalWriteIops"`
	} `json:"vmPerfMetricInfo"`
	VMProtectionGroupsInfo interface{} `json:"vmProtectionGroupsInfo"`
	VolumesInfo            []struct {
		DisplayName       string `json:"displayName"`
		ID                string `json:"id"`
		Name              string `json:"name"`
		ResourceURI       string `json:"resourceUri"`
		ScsiIdentifier    string `json:"scsiIdentifier"`
		SizeInBytes       int    `json:"sizeInBytes"`
		StorageFolderInfo struct {
			DisplayName string `json:"displayName"`
			ID          string `json:"id"`
			Name        string `json:"name"`
			ResourceURI string `json:"resourceUri"`
			Type        string `json:"type"`
		} `json:"storageFolderInfo"`
		StoragePoolInfo struct {
			DisplayName string `json:"displayName"`
			ID          string `json:"id"`
			Name        string `json:"name"`
			ResourceURI string `json:"resourceUri"`
			Type        string `json:"type"`
		} `json:"storagePoolInfo"`
		StorageSystemInfo struct {
			DisplayName  string `json:"displayName"`
			ID           string `json:"id"`
			Managed      bool   `json:"managed"`
			Name         string `json:"name"`
			ResourceURI  string `json:"resourceUri"`
			SerialNumber string `json:"serialNumber"`
			Type         string `json:"type"`
			VendorName   string `json:"vendorName"`
		} `json:"storageSystemInfo"`
		Type          string `json:"type"`
		VolumeSetInfo struct {
			DisplayName string `json:"displayName"`
			ID          string `json:"id"`
			Name        string `json:"name"`
			ResourceURI string `json:"resourceUri"`
			Type        string `json:"type"`
		} `json:"volumeSetInfo"`
	} `json:"volumesInfo"`
}
