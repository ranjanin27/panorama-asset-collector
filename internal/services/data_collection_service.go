// (C) Copyright 2023 Hewlett Packard Enterprise Development LP

package services

import (
	"context"
	"fmt"
	"os"
	"strings"
	"time"

	jsoniter "github.com/json-iterator/go"
	"github.com/rs/xid"
	"github.hpe.com/nimble-dcs/panorama-common-hauler/internal/clients/commonclient"
	"github.hpe.com/nimble-dcs/panorama-common-hauler/internal/handlers"
	"github.hpe.com/nimble-dcs/panorama-common-hauler/internal/utils/configs"
	"github.hpe.com/nimble-dcs/panorama-common-hauler/internal/utils/kafka/producer"
	"github.hpe.com/nimble-dcs/panorama-common-hauler/internal/utils/logging"
)

var (
	log  = logging.GetLogger()
	json = jsoniter.ConfigCompatibleWithStandardLibrary
)

const (
	Success               = "Success"
	Partial               = "Partial"
	Failed                = "Failed"
	DeviceType1           = "deviceType1"
	DeviceType2           = "deviceType2"
	UnknownDeviceType     = "Unknown"
	HaulerType            = "Fleet"
	GRPCError             = "GRPCError"
	AuthError             = "AuthError"
	Systems               = "Systems"
	SystemCapacity        = "SystemCapacity"
	StoragePools          = "StoragePools"
	Volumes               = "Volumes"
	Snapshots             = "Snapshots"
	VolumeCollections     = "VolumeCollections"
	VolumePerformance     = "VolumePerformance"
	Vluns                 = "Vluns"
	Applicationsets       = "Applicationsets"
	Inventory             = "Inventory"
	VolumePerf            = "Volume-Perf"
	UnknownCollectionType = "UnknownCT"
	Clones                = "Clones"
	DeviceType4           = "deviceType4"
	Common                = "Common"
	jsonExt               = ".json"
	PARTITIONKEY          = "partitionKey"
	equalExt              = "="
)

type dataCollectionService struct {
	commonClient commonclient.ArcusInterface
	ctx          context.Context
}

func NewDataCollectionService(ctx context.Context,
	aClient commonclient.ArcusInterface) DataCollectionServiceInterface {
	return &dataCollectionService{
		commonClient: aClient,
		ctx:          ctx,
	}
}

type DataCollectionServiceInterface interface {
	CollectDeviceInformation(ctx context.Context, consumerDetails handlers.ConsumerDetails,
		schedulerProducer producer.Producer, harmonyProducer producer.Producer)
}

var UploadToS3 = func(ctx context.Context, jsonContent *[]byte, bucketName, awsS3Region, awsAccessKeyID,
	awsSecretAccessKey string, uploadType UploadType) (string, int, error) {
	log.Debugf("bucketName: %v awsS3Region %v awsAccessKeyID %v awsSecretAccessKey %v uploadType %v\n",
		bucketName, awsS3Region, awsAccessKeyID, awsSecretAccessKey, uploadType)
	keyname, filesize, err := UploadToAwsS3(ctx, jsonContent, bucketName, awsS3Region, awsAccessKeyID,
		awsSecretAccessKey, uploadType)
	return keyname, filesize, err
}

func ConstructS3Object() string {
	now := time.Now().UTC()
	nows := fmt.Sprintf("%v", now)
	word1 := strings.Split(string(nows), " ")
	word2 := strings.Split(word1[1], ":")
	partitionKey := fmt.Sprintf("%v", word1[0]+"-"+word2[0]+"-"+word2[1])
	collectionID := xid.New().String()
	return PARTITIONKEY + equalExt + partitionKey + "/" + collectionID + jsonExt
}

//nolint:funlen,gocyclo,ineffassign // Code can be lengthy without a need for decomposition
func (dc *dataCollectionService) CollectDeviceInformation(ctx context.Context, consumerDetails handlers.ConsumerDetails,
	schedulerProducer producer.Producer, harmonyProducer producer.Producer) {
	//var collectionStartTime = time.Now().UTC().String()
	var mainErrorMap = make(map[string]map[string]string)
	var mainVMMap = make([]interface{}, 0)
	var mainDSMap = make([]interface{}, 0)
	var mainPPMap = make([]interface{}, 0)
	var mainVPGMap = make([]interface{}, 0)
	var mainPVMMap = make([]interface{}, 0)
	var mainCSPMIMap = make([]interface{}, 0)
	var mainPSMap = make([]interface{}, 0)
	var mainPSGMap = make([]interface{}, 0)
	var mainSOMap = make([]interface{}, 0)
	var mainCSPVMap = make([]interface{}, 0)
	var mainCSPAMap = make([]interface{}, 0)
	var mainZertoVPGMap = make([]interface{}, 0)
	var mainVMBackupMap = make(map[string][]interface{})
	var mainVMSnapshotMap = make(map[string][]interface{})
	var mainDSBackupMap = make(map[string][]interface{})
	var mainDSSnapshotMap = make(map[string][]interface{})

	dc.commonClient.SetCustomerIDForRest(consumerDetails.ApplicationCustomerID)

	authHeader, authErr := dc.commonClient.GetAuthHeaderForRest()
	if authErr != nil {
		log.WithContext(ctx).Errorf("Auth request failed : %v", authErr)
		return
	}

	// authHeader := "eyJhbGciOiJSUzI1NiIsImtpZCI6IlZ5WXdidVRLZnhwanhHbVVXUmtVbnZ1NU5xdyIsInBpLmF0bSI6IjFmN28ifQ.eyJzY29wZSI6Im9wZW5pZCBwcm9maWxlIGVtYWlsIiwiY2xpZW50X2lkIjoiZGFmZDdmOWYtN2NhYy00MjBhLWI5MTAtZmQyYjc4NDZmZmU0IiwiaXNzIjoiaHR0cHM6Ly9kZXYtc3NvLmNjcy5hcnViYXRoZW5hLmNvbSIsImF1ZCI6ImF1ZCIsImxhc3ROYW1lIjoiaHBlIiwic3ViIjoiaHBlLmF1dGgudGVzdEBnbWFpbC5jb20iLCJ1c2VyX2N0eCI6Ijg5MjJhZmE2NzIzMDExZWJiZTAxY2EzMmQzMmI2Yjc3IiwiYXV0aF9zb3VyY2UiOiJwMTRjIiwiZ2l2ZW5OYW1lIjoiYXV0aHoiLCJpYXQiOjE2MTU5MTAzOTksImV4cCI6MTYxNTkxNzU5OX0.jLniCfT7DbPsZpzVBuYKrUvQ02VFEYhtULAd4NmT1ohPtiy3ybhY1oEjG6GsxMeOvD-6wMNokZqae3Zrt4BJrlENm0G00TF-jcbsKGkRHfqRxdpjS5yifOCySIwykcierd_32O0saTkNKj1FP56NzVKoRa8REdfgHawaFjsMhQ9nwDvftTwiANQqWF9tu1icIFjAuXJV5SVeOKf05ypnYLPtaMn5feTmxbteJh6fhsDx2y9SHDFgx6N8TkIDTu6yTKIFvNo85MdvDnzCFRNj6zzbCGIHPyjiL0hBuXyXQlI9j5FMjC2m7JICM2PSyR1BGD7Y7IULAlf_kaIMST4UNQ"

	pKey := ConstructS3Object()

	virtualmachines, vErr := dc.commonClient.GetVMs(ctx, authHeader)
	if vErr != nil {
		log.WithContext(ctx).Errorf("GetVMs request failed : %v", vErr)
		handlers.SetNested(mainErrorMap, "VirtualMachines", "VirtualMachines", vErr.Error())
	}

	if virtualmachines != nil {
		log.WithContext(ctx).Infof("VM count - %v", len(virtualmachines))
		for v := range virtualmachines {
			mainVMMap = append(mainVMMap, virtualmachines[v])
		}
		u, err := json.Marshal(mainVMMap)
		if err != nil {
			log.WithContext(ctx).Errorf("err is %v", err)
		}
		err = os.WriteFile("./vmsTest", u, 0644)
		log.WithContext(ctx).Errorf("file write err is %v", err)
		key, filesize, err := UploadToS3(ctx, &u, configs.GetAWSS3BucketName(), configs.GetAWSRegion(),
			configs.GetAWSAccessKey(), configs.GetAWSSecretAccessKey(), UploadType(configs.GetSourceType()+"/VM/"+pKey))
		log.WithContext(ctx).Infof("Err = %v, key= %v, filesize = %v", err, key, filesize)
	}

	for v := range virtualmachines {
		vmbackups, vmbkpErr := dc.commonClient.GetVMBackups(ctx, virtualmachines[v].ID, authHeader)
		if vmbkpErr != nil {
			log.WithContext(ctx).Errorf("GetVMBackups request failed - %v : %v",
				virtualmachines[v].ID, vmbkpErr.Error())
			handlers.SetNested(mainErrorMap, "VMBackups", virtualmachines[v].ID, vmbkpErr.Error())
		}
		for vb := range vmbackups {
			mainVMBackupMap[virtualmachines[v].ID] = append(mainVMBackupMap[virtualmachines[v].ID], vmbackups[vb])
		}
	}
	bkup, err := json.Marshal(mainVMBackupMap)
	if err != nil {
		log.WithContext(ctx).Errorf("err is %v", err)
	}
	err = os.WriteFile("./vmBackups", bkup, 0644)
	log.WithContext(ctx).Errorf("file write err is %v", err)

	UploadToS3(ctx, &bkup, configs.GetAWSS3BucketName(), configs.GetAWSRegion(),
		configs.GetAWSAccessKey(), configs.GetAWSSecretAccessKey(), UploadType(configs.GetSourceType()+"/VMBK/"+pKey))

	for v := range virtualmachines {
		vmSnapshots, vmsnapErr := dc.commonClient.GetVMSnapshots(ctx, virtualmachines[v].ID, authHeader)
		if vmsnapErr != nil {
			log.WithContext(ctx).Errorf("GetVMSnapshots request failed - %v : %v",
				virtualmachines[v].ID, vmsnapErr.Error())
			handlers.SetNested(mainErrorMap, "VMSnapshots", virtualmachines[v].ID, vmsnapErr.Error())
		}
		for vs := range vmSnapshots {
			mainVMSnapshotMap[virtualmachines[v].ID] = append(mainVMSnapshotMap[virtualmachines[v].ID], vmSnapshots[vs])
		}
	}
	snap, err := json.Marshal(mainVMSnapshotMap)
	if err != nil {
		log.WithContext(ctx).Errorf("err is %v", err)
	}
	err = os.WriteFile("./vmSnapshots", snap, 0644)
	log.WithContext(ctx).Errorf("file write err is %v", err)
	key, filesize, err := UploadToS3(ctx, &snap, configs.GetAWSS3BucketName(), configs.GetAWSRegion(),
		configs.GetAWSAccessKey(), configs.GetAWSSecretAccessKey(), UploadType(configs.GetSourceType()+"/VMSNP/"+pKey))
	log.WithContext(ctx).Infof("Err = %v, ke:y= %v, filesize = %v", err, key, filesize)
	virtualmachines = nil

	datastores, dErr := dc.commonClient.GetDatastores(ctx, authHeader)
	if dErr != nil {
		log.WithContext(ctx).Errorf("GetDatastores request failed : %v", dErr)
		handlers.SetNested(mainErrorMap, "Datastores", "Datastores", dErr.Error())
	}

	if datastores != nil {
		log.WithContext(ctx).Infof("Datastores count - %v", len(datastores))
		for d := range datastores {
			mainDSMap = append(mainDSMap, datastores[d])
		}
		u, err := json.Marshal(mainDSMap)
		if err != nil {
			log.WithContext(ctx).Errorf("err is %v", err)
		}
		err = os.WriteFile("./datastoreTest", u, 0644)
		log.WithContext(ctx).Errorf("file write err is %v", err)
		key, filesize, err := UploadToS3(ctx, &u, configs.GetAWSS3BucketName(), configs.GetAWSRegion(),
			configs.GetAWSAccessKey(), configs.GetAWSSecretAccessKey(), UploadType(configs.GetSourceType()+"/DS/"+pKey))
		log.WithContext(ctx).Infof("Err = %v, key= %v, filesize = %v", err, key, filesize)
	}

	for d := range datastores {
		dsbackups, dsbkpErr := dc.commonClient.GetDSBackups(ctx, datastores[d].ID, authHeader)
		if dsbkpErr != nil {
			log.WithContext(ctx).Errorf("GetDSBackups request failed - %v : %v",
				datastores[d].ID, dsbkpErr.Error())
			handlers.SetNested(mainErrorMap, "DatastoreBackups", datastores[d].ID, dsbkpErr.Error())
		}
		for db := range dsbackups {
			mainDSBackupMap[datastores[d].ID] = append(mainDSBackupMap[datastores[d].ID], dsbackups[db])
		}
	}
	u, err := json.Marshal(mainDSBackupMap)
	if err != nil {
		log.WithContext(ctx).Errorf("err is %v", err)
	}
	err = os.WriteFile("./datastoreBKTest", u, 0644)
	log.WithContext(ctx).Errorf("file write err is %v", err)
	key, filesize, err = UploadToS3(ctx, &u, configs.GetAWSS3BucketName(), configs.GetAWSRegion(),
		configs.GetAWSAccessKey(), configs.GetAWSSecretAccessKey(), UploadType(configs.GetSourceType()+"/DSBK/"+pKey))
	log.WithContext(ctx).Infof("Err = %v, key= %v, filesize = %v", err, key, filesize)

	for d := range datastores {
		dsSnapshots, dssnapErr := dc.commonClient.GetDSSnapshots(ctx, datastores[d].ID, authHeader)
		if dssnapErr != nil {
			log.WithContext(ctx).Errorf("GetDSSnapshots request failed - %v : %v",
				datastores[d].ID, dssnapErr.Error())
			handlers.SetNested(mainErrorMap, "DatastoreSnapshots", datastores[d].ID, dssnapErr.Error())
		}
		for ds := range dsSnapshots {
			mainDSSnapshotMap[datastores[d].ID] = append(mainDSSnapshotMap[datastores[d].ID], dsSnapshots[ds])
		}
	}
	dssnap, err := json.Marshal(mainDSSnapshotMap)
	if err != nil {
		log.WithContext(ctx).Errorf("err is %v", err)
	}
	err = os.WriteFile("./datastoreSnapshotTest", dssnap, 0644)
	log.WithContext(ctx).Errorf("file write err is %v", err)
	key, filesize, err = UploadToS3(ctx, &dssnap, configs.GetAWSS3BucketName(), configs.GetAWSRegion(),
		configs.GetAWSAccessKey(), configs.GetAWSSecretAccessKey(), UploadType(configs.GetSourceType()+"/DSSNP/"+pKey))
	log.WithContext(ctx).Infof("Err = %v, key= %v, filesize = %v", err, key, filesize)

	datastores = nil
	protectionpolicies, pErr := dc.commonClient.GetProtectionPolicies(ctx, authHeader)
	if pErr != nil {
		log.WithContext(ctx).Errorf("GetProtectionPolicies request failed : %v", pErr)
		handlers.SetNested(mainErrorMap, "ProtectionPolicies", "ProtectionPolicies", pErr.Error())
	}

	if protectionpolicies != nil {
		log.WithContext(ctx).Infof("ProtectionPolicies count - %v", len(protectionpolicies))
		for p := range protectionpolicies {
			mainPPMap = append(mainPPMap, protectionpolicies[p])
		}
		pp, err := json.Marshal(mainPPMap)
		if err != nil {
			log.WithContext(ctx).Errorf("err is %v", err)
		}
		err = os.WriteFile("./ppTest", pp, 0644)
		log.WithContext(ctx).Errorf("file write err is %v", err)
		key, filesize, err := UploadToS3(ctx, &pp, configs.GetAWSS3BucketName(), configs.GetAWSRegion(),
			configs.GetAWSAccessKey(), configs.GetAWSSecretAccessKey(), UploadType(configs.GetSourceType()+"/PP/"+pKey))
		log.WithContext(ctx).Infof("Err = %v, key= %v, filesize = %v", err, key, filesize)

		protectionpolicies = nil
	}

	vpg, vpgErr := dc.commonClient.GetVMProtectionGroups(ctx, authHeader)
	if vpgErr != nil {
		log.WithContext(ctx).Errorf("GetVMProtectionGroups request failed : %v", vpgErr)
		handlers.SetNested(mainErrorMap, "VMProtectionGroup", "VMProtectionGroup", vpgErr.Error())
	}

	if vpg != nil {
		log.WithContext(ctx).Infof("GetVMProtectionGroups count - %v", len(vpg))
		for v := range vpg {
			mainVPGMap = append(mainVPGMap, vpg[v])
		}
		vpg, err := json.Marshal(mainVPGMap)
		if err != nil {
			log.WithContext(ctx).Errorf("err is %v", err)
		}
		err = os.WriteFile("./vpgTest", vpg, 0644)
		log.WithContext(ctx).Errorf("file write err is %v", err)
		key, filesize, err := UploadToS3(ctx, &vpg, configs.GetAWSS3BucketName(), configs.GetAWSRegion(),
			configs.GetAWSAccessKey(), configs.GetAWSSecretAccessKey(), UploadType(configs.GetSourceType()+"/VMPG/"+pKey))
		log.WithContext(ctx).Infof("Err = %v, key= %v, filesize = %v", err, key, filesize)

		vpg = nil
	}

	pvms, pvmErr := dc.commonClient.GetProtectedVMs(ctx, authHeader)
	if pvmErr != nil {
		log.WithContext(ctx).Errorf("GetProtectedVMs request failed : %v", pvmErr)
		handlers.SetNested(mainErrorMap, "ProtectedVMs", "ProtectedVMs", pvmErr.Error())
	}

	if pvms != nil {
		log.WithContext(ctx).Infof("GetProtectedVMs count - %v", len(pvms))
		for p := range pvms {
			mainPVMMap = append(mainPVMMap, pvms[p])
		}
		pvm, err := json.Marshal(mainPVMMap)
		if err != nil {
			log.WithContext(ctx).Errorf("err is %v", err)
		}
		err = os.WriteFile("./pvmTest", pvm, 0644)
		log.WithContext(ctx).Errorf("file write err is %v", err)
		key, filesize, err := UploadToS3(ctx, &pvm, configs.GetAWSS3BucketName(), configs.GetAWSRegion(),
			configs.GetAWSAccessKey(), configs.GetAWSSecretAccessKey(), UploadType(configs.GetSourceType()+"/PVM/"+pKey))
		log.WithContext(ctx).Infof("Err = %v, key= %v, filesize = %v", err, key, filesize)

		pvms = nil
	}

	cspMI, cspErr := dc.commonClient.GetCSPMachineInstances(ctx, authHeader)
	if cspErr != nil {
		log.WithContext(ctx).Errorf("GetCSPMachineInstances request failed : %v", cspErr)
		handlers.SetNested(mainErrorMap, "CSPMachineInstances", "CSPMachineInstances", cspErr.Error())
	}

	if cspMI != nil {
		log.WithContext(ctx).Infof("GetCSPMachineInstances count - %v", len(cspMI))
		for csp := range cspMI {
			mainCSPMIMap = append(mainCSPMIMap, cspMI[csp])
		}
		csp, err := json.Marshal(mainCSPMIMap)
		if err != nil {
			log.WithContext(ctx).Errorf("err is %v", err)
		}
		err = os.WriteFile("./pvmTest", csp, 0644)
		log.WithContext(ctx).Errorf("file write err is %v", err)
		key, filesize, err := UploadToS3(ctx, &csp, configs.GetAWSS3BucketName(), configs.GetAWSRegion(),
			configs.GetAWSAccessKey(), configs.GetAWSSecretAccessKey(), UploadType(configs.GetSourceType()+"/EC2/"+pKey))
		log.WithContext(ctx).Infof("Err = %v, key= %v, filesize = %v", err, key, filesize)
		cspMI = nil
	}

	zvpgs, zvpgErr := dc.commonClient.GetZertoVPGs(ctx, authHeader)
	if zvpgErr != nil {
		log.WithContext(ctx).Errorf("GetZertoVPGs request failed : %v", zvpgErr)
		handlers.SetNested(mainErrorMap, "ZertoVPGs", "ZertoVPGs", zvpgErr.Error())
	}

	if zvpgs != nil {
		log.WithContext(ctx).Infof("GetZertoVPGs count - %v", len(zvpgs))
		for zvpg := range zvpgs {
			mainZertoVPGMap = append(mainZertoVPGMap, zvpgs[zvpg])
		}
		zmap, err := json.Marshal(mainZertoVPGMap)
		if err != nil {
			log.WithContext(ctx).Errorf("err is %v", err)
		}
		err = os.WriteFile("./pvmTest", zmap, 0644)
		log.WithContext(ctx).Errorf("file write err is %v", err)
		key, filesize, err := UploadToS3(ctx, &zmap, configs.GetAWSS3BucketName(), configs.GetAWSRegion(),
			configs.GetAWSAccessKey(), configs.GetAWSSecretAccessKey(), UploadType(configs.GetSourceType()+"/ZERTO/"+pKey))
		log.WithContext(ctx).Infof("Err = %v, key= %v, filesize = %v", err, key, filesize)
		zvpgs = nil
	}

	ps, psErr := dc.commonClient.GetProtectionStores(ctx, authHeader)
	if psErr != nil {
		log.WithContext(ctx).Errorf("GetProtectionStores request failed : %v", psErr)
		handlers.SetNested(mainErrorMap, "ProtectionStores", "ProtectionStores", psErr.Error())
	}

	if ps != nil {
		log.WithContext(ctx).Infof("GetProtectionStores count - %v", len(ps))
		for pst := range ps {
			mainPSMap = append(mainPSMap, ps[pst])
		}
		u, err := json.Marshal(mainPSMap)
		if err != nil {
			log.WithContext(ctx).Errorf("err is %v", err)
		}
		err = os.WriteFile("./psTest", u, 0644)
		log.WithContext(ctx).Errorf("file write err is %v", err)
		key, filesize, err := UploadToS3(ctx, &u, configs.GetAWSS3BucketName(), configs.GetAWSRegion(),
			configs.GetAWSAccessKey(), configs.GetAWSSecretAccessKey(), UploadType(configs.GetSourceType()+"/PS/"+pKey))
		log.WithContext(ctx).Infof("Err = %v, key= %v, filesize = %v", err, key, filesize)
		ps = nil
	}

	psg, psgErr := dc.commonClient.GetProtectionStoreGateways(ctx, authHeader)
	if psgErr != nil {
		log.WithContext(ctx).Errorf("GetProtectionStoreGateways request failed : %v", psgErr)
		handlers.SetNested(mainErrorMap, "ProtectionStoreGateways", "ProtectionStoreGateways", psgErr.Error())
	}

	if psg != nil {
		log.WithContext(ctx).Infof("GetProtectionStoreGateways count - %v", len(psg))
		for psgy := range psg {
			mainPSGMap = append(mainPSGMap, psg[psgy])
		}
		m, err := json.Marshal(mainPSGMap)
		if err != nil {
			log.WithContext(ctx).Errorf("err is %v", err)
		}
		err = os.WriteFile("./psmTest", m, 0644)
		log.WithContext(ctx).Errorf("file write err is %v", err)
		key, filesize, err := UploadToS3(ctx, &m, configs.GetAWSS3BucketName(), configs.GetAWSRegion(),
			configs.GetAWSAccessKey(), configs.GetAWSSecretAccessKey(), UploadType(configs.GetSourceType()+"/PSG/"+pKey))
		log.WithContext(ctx).Infof("Err = %v, key= %v, filesize = %v", err, key, filesize)
		psg = nil
	}

	sos, sosErr := dc.commonClient.GetStoreonces(ctx, authHeader)
	if sosErr != nil {
		log.WithContext(ctx).Errorf("GetStoreonces request failed : %v", sosErr)
		handlers.SetNested(mainErrorMap, "Storeonces", "Storeonces", sosErr.Error())
	}

	if sos != nil {
		log.WithContext(ctx).Infof("GetStoreonces count - %v", len(sos))
		for so := range sos {
			mainSOMap = append(mainSOMap, sos[so])
		}
		sMap, err := json.Marshal(mainSOMap)
		if err != nil {
			log.WithContext(ctx).Errorf("err is %v", err)
		}
		err = os.WriteFile("./sosTest", sMap, 0644)
		log.WithContext(ctx).Errorf("file write err is %v", err)
		key, filesize, err := UploadToS3(ctx, &sMap, configs.GetAWSS3BucketName(), configs.GetAWSRegion(),
			configs.GetAWSAccessKey(), configs.GetAWSSecretAccessKey(), UploadType(configs.GetSourceType()+"/STOREONCE/"+pKey))
		log.WithContext(ctx).Infof("Err = %v, key= %v, filesize = %v", err, key, filesize)
		sos = nil
	}

	cspv, cspvErr := dc.commonClient.GetCSPVolumes(ctx, authHeader)
	if cspvErr != nil {
		log.WithContext(ctx).Errorf("GetCSPVolumes request failed : %v", cspvErr)
		handlers.SetNested(mainErrorMap, "CSPVolumes", "CSPVolumes", cspvErr.Error())
	}

	if cspv != nil {
		log.WithContext(ctx).Infof("GetCSPVolumes count - %v", len(cspv))
		for so := range cspv {
			mainCSPVMap = append(mainCSPVMap, cspv[so])
		}
		sMap, err := json.Marshal(mainCSPVMap)
		if err != nil {
			log.WithContext(ctx).Errorf("err is %v", err)
		}
		err = os.WriteFile("./cspTest", sMap, 0644)
		log.WithContext(ctx).Errorf("file write err is %v", err)
		key, filesize, err := UploadToS3(ctx, &sMap, configs.GetAWSS3BucketName(), configs.GetAWSRegion(),
			configs.GetAWSAccessKey(), configs.GetAWSSecretAccessKey(), UploadType(configs.GetSourceType()+"/EBS/"+pKey))
		log.WithContext(ctx).Infof("Err = %v, key= %v, filesize = %v", err, key, filesize)
		cspv = nil
	}

	cspa, cspaErr := dc.commonClient.GetCSPAccounts(ctx, authHeader)
	if cspaErr != nil {
		log.WithContext(ctx).Errorf("GetCSPAccounts request failed : %v", cspaErr)
		handlers.SetNested(mainErrorMap, "CSPAccounts", "CSPAccounts", cspaErr.Error())
	}

	if cspa != nil {
		log.WithContext(ctx).Infof("GetCSPAccounts count - %v", len(cspa))
		for so := range cspa {
			mainCSPAMap = append(mainCSPAMap, cspa[so])
		}
		cspMap, err := json.Marshal(mainCSPAMap)
		if err != nil {
			log.WithContext(ctx).Errorf("err is %v", err)
		}
		err = os.WriteFile("./cspTest", cspMap, 0644)
		log.WithContext(ctx).Errorf("file write err is %v", err)
		key, filesize, err := UploadToS3(ctx, &cspMap, configs.GetAWSS3BucketName(), configs.GetAWSRegion(),
			configs.GetAWSAccessKey(), configs.GetAWSSecretAccessKey(), UploadType(configs.GetSourceType()+"/ACC/"+pKey))
		log.WithContext(ctx).Infof("Err = %v, key= %v, filesize = %v", err, key, filesize)
		cspa = nil
	}

}
