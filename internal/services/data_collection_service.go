// (C) Copyright 2022 Hewlett Packard Enterprise Development LP

package services

import (
	"context"
	"time"

	jsoniter "github.com/json-iterator/go"
	"github.hpe.com/nimble-dcs/panorama-common-hauler/internal/clients/commonclient"
	"github.hpe.com/nimble-dcs/panorama-common-hauler/internal/handlers"
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

//nolint:funlen,gocyclo,ineffassign // Code can be lengthy without a need for decomposition
func (dc *dataCollectionService) CollectDeviceInformation(ctx context.Context, consumerDetails handlers.ConsumerDetails,
	schedulerProducer producer.Producer, harmonyProducer producer.Producer) {
	var collectionStartTime = time.Now().UTC().String()
	var mainErrorMap = make(map[string]map[string]string)
	var mainVMMap = make([]interface{}, 0)
	var mainDSMap = make([]interface{}, 0)
	var mainPPMap = make([]interface{}, 0)
	var mainVPGMap = make([]interface{}, 0)
	var mainPVMMap = make([]interface{}, 0)
	var mainCSPMIMap = make([]interface{}, 0)
	var mainZertoVPGMap = make([]interface{}, 0)
	var mainVMBackupMap = make(map[string][]interface{})
	var mainVMSnapshotMap = make(map[string][]interface{})
	var mainDSBackupMap = make(map[string][]interface{})
	var mainDSSnapshotMap = make(map[string][]interface{})

	dc.commonClient.SetCustomerIDForRest(consumerDetails.ApplicationCustomerID)

	// authHeader, authErr := dc.commonClient.GetAuthHeaderForRest()
	// if authErr != nil {
	// 	log.WithContext(ctx).Errorf("Auth request failed : %v", authErr)
	// 	return
	// }

	authHeader := "eyJhbGciOiJSUzI1NiIsImtpZCI6IlZ5WXdidVRLZnhwanhHbVVXUmtVbnZ1NU5xdyIsInBpLmF0bSI6IjFmN28ifQ.eyJzY29wZSI6Im9wZW5pZCBwcm9maWxlIGVtYWlsIiwiY2xpZW50X2lkIjoiZGFmZDdmOWYtN2NhYy00MjBhLWI5MTAtZmQyYjc4NDZmZmU0IiwiaXNzIjoiaHR0cHM6Ly9kZXYtc3NvLmNjcy5hcnViYXRoZW5hLmNvbSIsImF1ZCI6ImF1ZCIsImxhc3ROYW1lIjoiaHBlIiwic3ViIjoiaHBlLmF1dGgudGVzdEBnbWFpbC5jb20iLCJ1c2VyX2N0eCI6Ijg5MjJhZmE2NzIzMDExZWJiZTAxY2EzMmQzMmI2Yjc3IiwiYXV0aF9zb3VyY2UiOiJwMTRjIiwiZ2l2ZW5OYW1lIjoiYXV0aHoiLCJpYXQiOjE2MTU5MTAzOTksImV4cCI6MTYxNTkxNzU5OX0.jLniCfT7DbPsZpzVBuYKrUvQ02VFEYhtULAd4NmT1ohPtiy3ybhY1oEjG6GsxMeOvD-6wMNokZqae3Zrt4BJrlENm0G00TF-jcbsKGkRHfqRxdpjS5yifOCySIwykcierd_32O0saTkNKj1FP56NzVKoRa8REdfgHawaFjsMhQ9nwDvftTwiANQqWF9tu1icIFjAuXJV5SVeOKf05ypnYLPtaMn5feTmxbteJh6fhsDx2y9SHDFgx6N8TkIDTu6yTKIFvNo85MdvDnzCFRNj6zzbCGIHPyjiL0hBuXyXQlI9j5FMjC2m7JICM2PSyR1BGD7Y7IULAlf_kaIMST4UNQ"

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
		cspMI = nil
	}

	zvpgs, zvpgErr := dc.commonClient.GetZertoVPGs(ctx, authHeader)
	if zvpgErr != nil {
		log.WithContext(ctx).Errorf("GetZertoVPGs request failed : %v", zvpgErr)
		handlers.SetNested(mainErrorMap, "ZertoVPGs", "ZertoVPGs", cspErr.Error())
	}

	if zvpgs != nil {
		log.WithContext(ctx).Infof("GetZertoVPGs count - %v", len(zvpgs))
		for zvpg := range zvpgs {
			mainZertoVPGMap = append(mainZertoVPGMap, zvpgs[zvpg])
		}
		zvpgs = nil
	}

	var data = handlers.ConstructCommonJSON(consumerDetails, collectionStartTime, Common,
		mainVMMap, mainDSMap, mainPPMap, mainVPGMap, mainPVMMap, mainCSPMIMap, mainZertoVPGMap,
		mainVMBackupMap, mainDSBackupMap, mainErrorMap)
	file, _ := json.MarshalIndent(data, "", " ")
	// _ = ioutil.WriteFile("test.json", file, 0644)
	handlers.UploadToServer(ctx, file, consumerDetails, Common, mainErrorMap, schedulerProducer,
		harmonyProducer)
	file = nil
}
