// (C) Copyright 2022 Hewlett Packard Enterprise Development LP

package model

type ScheduleStatus struct {
	CollectionID          string `json:"collection_id"`
	PlatformCustomerID    string `json:"platform_customer_id"`
	ApplicationCustomerID string `json:"application_customer_id"`
	CollectionType        string `json:"collection_type"`
	HaulerType            string `json:"hauler_type"`
	DeviceType            string `json:"device_type"`
	JSONVersion           string `json:"json_version"`
	UploadStatus          string `json:"upload_status"`
	CollectionStatus      string `json:"collection_status"`
	UploadFileSize        string `json:"upload_file_size"`
	S3Bucket              string `json:"s3_bucket"`
	FileName              string `json:"file_name"`
	Error                 string `json:"error"`
}

type KafkaCollectionRequest struct {
	PlatformCustomerID    string `json:"platform_customer_id"`
	CollectionType        string `json:"collection_type"`
	CollectionID          string `json:"collection_id"`
	CollectionTrigger     string `json:"collection_trigger"`
	Region                string `json:"region"`
	ApplicationCustomerID string `json:"application_customer_id"`
	ApplicationInstanceID string `json:"application_instance_id"`
}

type HarmonyStatusRequest struct {
	Source struct {
		S3 struct {
			Bucket string `json:"bucket"`
			Key    string `json:"key"`
		} `json:"s3"`
	} `json:"source"`
	Notification struct {
		FileDomain string `json:"fileDomain"`
		FileType   string `json:"fileType"`
		FileSize   int    `json:"fileSize"`
		EntityID   string `json:"entityId"`
		FileName   string `json:"fileName"`
	} `json:"notification"`
}
