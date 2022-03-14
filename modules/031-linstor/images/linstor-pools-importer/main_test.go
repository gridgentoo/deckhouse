/*
Copyright 2022 Flant JSC

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package main

import (
	"context"
	"reflect"
	"testing"
	"time"

	lclient "github.com/LINBIT/golinstor/client"
	v1 "k8s.io/api/core/v1"
	storagev1 "k8s.io/api/storage/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func TestParseThinPoolsEmpty(t *testing.T) {
	got, err := parseThinPools(``)
	if err != nil {
		t.Errorf("\nexpected no error\ngot: %s", err.Error())
	}
	if got != nil {
		t.Errorf("\nexpected nil\ngot: %+v", got)
	}
}

func TestParseThinPoolsWrong(t *testing.T) {
	_, err := parseThinPools(`a;b`)
	if err == nil {
		t.Errorf("\nexpected error\ngot: nil")
	}
}

func TestParseThinPoolsNoTags(t *testing.T) {
	got, err := parseThinPools(`  data;linstor_data;twi---tz--;
  pvc-ecc0e656-78ca-497f-8f7a-f9fe3b384748_00000;linstor_data;Vwi-aotz--;
  root;vg0;-wi-ao----;`)
	if err != nil {
		t.Errorf("\nexpected no error\ngot: %s", err.Error())
	}
	if got != nil {
		t.Errorf("\nexpected nil\ngot: %+v", got)
	}
}

func TestParseThinPoolsWithTags(t *testing.T) {
	got, err := parseThinPools(`  data;linstor_data;twi---tz--;linstor-ssd
  pvc-ecc0e656-78ca-497f-8f7a-f9fe3b384748_00000;linstor_data;Vwi-aotz--;linstor-ssd
  root;vg0;-wi-ao----;`)
	if err != nil {
		t.Errorf("\nexpected no error\ngot: %s", err.Error())
	}
	expected := []lclient.StoragePool{
		lclient.StoragePool{
			StoragePoolName: "ssd",
			ProviderKind:    lclient.LVM_THIN,
			NodeName:        nodeName,
			Props: map[string]string{
				"StorDriver/LvmVg":    "linstor_data",
				"StorDriver/ThinPool": "data",
			},
		},
	}
	if !reflect.DeepEqual(got, expected) {
		t.Errorf("\nexpected: %+v\ngot: %+v", expected, got)
	}
}

func TestParseVolumeGroupsEmpty(t *testing.T) {
	got, err := parseVolumeGroups(``)
	if err != nil {
		t.Errorf("\nexpected no error\ngot: %s", err.Error())
	}
	if got != nil {
		t.Errorf("\nexpected nil\ngot: %+v", got)
	}
}

func TestParseVolumeGroupsWrong(t *testing.T) {
	_, err := parseVolumeGroups(`avasd`)
	if err == nil {
		t.Errorf("\nexpected error\ngot: nil")
	}
}

func TestParseVolumeGroupsNoTags(t *testing.T) {
	got, err := parseVolumeGroups(`  linstor_data;
  vg0;`)
	if err != nil {
		t.Errorf("\nexpected no error\ngot: %s", err.Error())
	}
	if got != nil {
		t.Errorf("\nexpected nil\ngot: %+v", got)
	}
}

func TestParseVolumeGroupsWithTags(t *testing.T) {
	got, err := parseVolumeGroups(`  linstor_data;linstor-data
  vg0;`)
	if err != nil {
		t.Errorf("\nexpected no error\ngot: %s", err.Error())
	}
	expected := []lclient.StoragePool{
		lclient.StoragePool{
			StoragePoolName: "data",
			ProviderKind:    lclient.LVM,
			NodeName:        nodeName,
			Props: map[string]string{
				"StorDriver/LvmVg": "linstor_data",
			},
		},
	}
	if !reflect.DeepEqual(got, expected) {
		t.Errorf("\nexpected: %+v\ngot: %+v", expected, got)
	}
}

func TestNewKubernetesStorageClasses(t *testing.T) {
	tp := lclient.StoragePool{
		StoragePoolName: "ssd",
		ProviderKind:    lclient.LVM_THIN,
		NodeName:        nodeName,
		Props: map[string]string{
			"StorDriver/LvmVg":    "linstor_data",
			"StorDriver/ThinPool": "data",
		},
	}
	got := newKubernetesStorageClass(&tp, 2)

	volBindMode := storagev1.VolumeBindingImmediate
	allowVolumeExpansion := true
	reclaimPolicy := v1.PersistentVolumeReclaimDelete

	expected := storagev1.StorageClass{
		ObjectMeta: metav1.ObjectMeta{
			Name: "linstor-ssd-r2",
		},
		Provisioner:          "linstor.csi.linbit.com",
		VolumeBindingMode:    &volBindMode,
		AllowVolumeExpansion: &allowVolumeExpansion,
		ReclaimPolicy:        &reclaimPolicy,
		Parameters: map[string]string{
			"linstor.csi.linbit.com/storagePool":    "ssd",
			"linstor.csi.linbit.com/placementCount": "2",
		},
	}

	if !reflect.DeepEqual(got, expected) {
		t.Errorf("\nexpected: %+v\ngot: %+v", expected, got)
	}
}

func TestCandidatesLoop(t *testing.T) {

	vg1 := lclient.StoragePool{
		StoragePoolName: "data",
		ProviderKind:    lclient.LVM,
		NodeName:        nodeName,
		Props: map[string]string{
			"StorDriver/LvmVg": "linstor_data",
		},
	}

	vg2 := lclient.StoragePool{
		StoragePoolName: "hdd",
		ProviderKind:    lclient.LVM,
		NodeName:        nodeName,
		Props: map[string]string{
			"StorDriver/LvmVg": "linstor_hdd",
		},
	}

	tp := lclient.StoragePool{
		StoragePoolName: "ssd",
		ProviderKind:    lclient.LVM_THIN,
		NodeName:        nodeName,
		Props: map[string]string{
			"StorDriver/LvmVg":    "linstor_data",
			"StorDriver/ThinPool": "data",
		},
	}

	ctx, cancel := context.WithTimeout(context.TODO(), 10*time.Millisecond)
	defer cancel()
	var attempt int
	candidatesChannel := runCandidatesLoop(ctx, func() ([]lclient.StoragePool, error) {
		var s []lclient.StoragePool
		attempt++
		if attempt < 2 {
			s = []lclient.StoragePool{vg1, vg2}
		} else {
			s = []lclient.StoragePool{vg1, vg2, tp}
		}
		return s, nil
	}, time.Millisecond)

	var expectedStorageClassesNames = []string{
		"linstor-data-r1", "linstor-data-r2", "linstor-data-r3",
		"linstor-hdd-r1", "linstor-hdd-r2", "linstor-hdd-r3",
		"linstor-ssd-r1", "linstor-ssd-r2", "linstor-ssd-r3",
	}
	var gotStorageClassesNames []string

	for candidate := range candidatesChannel {
		// Create StorageClasses in Kubernetes
		for r := 1; r <= maxReplicasNum; r++ {
			storageClass := newKubernetesStorageClass(&candidate, r)
			gotStorageClassesNames = append(gotStorageClassesNames, storageClass.GetName())
		}
	}

	if !reflect.DeepEqual(gotStorageClassesNames, expectedStorageClassesNames) {
		t.Errorf("\nexpected Kubernetes storage classes: %+v\ngot: %+v", expectedStorageClassesNames, gotStorageClassesNames)
	}
}
