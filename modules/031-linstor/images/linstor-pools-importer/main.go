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
	"bytes"
	"context"
	"flag"
	"fmt"
	"os"
	"os/exec"
	"runtime"
	"strings"
	"time"

	// Import all Kubernetes client auth plugins (e.g. Azure, GCP, OIDC, etc.)

	_ "k8s.io/client-go/plugin/pkg/client/auth"
	"k8s.io/client-go/tools/clientcmd"
	kclient "sigs.k8s.io/controller-runtime/pkg/client"

	v1 "k8s.io/api/core/v1"
	storagev1 "k8s.io/api/storage/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/klog"

	lclient "github.com/LINBIT/golinstor/client"
)

const (
	lvmConfig      = `devices {filter=["r|^/dev/drbd*|"]}`
	linstorPrefix  = "linstor"
	maxReplicasNum = 3
)

var (
	nodeName  = os.Getenv("NODE_NAME")
	podName   = os.Getenv("POD_NAME")
	namespace = os.Getenv("NAMESPACE")
)

var supportedProviderKinds = []lclient.ProviderKind{lclient.LVM, lclient.LVM_THIN}

func printVersion() {
	klog.Info(fmt.Sprintf("Go Version: %s", runtime.Version()))
	klog.Info(fmt.Sprintf("Go OS/Arch: %s/%s", runtime.GOOS, runtime.GOARCH))
}

func main() {
	delaySeconds := flag.Int("delay", 10, "Delay in seconds between scanning attempts")

	flag.Parse()
	printVersion()

	clientConfig := clientcmd.NewNonInteractiveDeferredLoadingClientConfig(
		clientcmd.NewDefaultClientConfigLoadingRules(),
		&clientcmd.ConfigOverrides{},
	)

	if nodeName == "" {
		klog.Fatalln("Required NODE_NAME env variable is not specified!")
	}
	klog.Infof("nodeName: %+v", nodeName)
	if podName == "" {
		klog.Fatalln("Required POD_NAME env variable is not specified!")
	}
	klog.Infof("podName: %+v", podName)
	if namespace == "" {
		klog.Fatalln("Required NAMESPACE env variable is not specified!")
	}
	klog.Infof("namespace: %+v", namespace)

	ctx := context.Background()

	// Get a config to talk to the apiserver
	config, err := clientConfig.ClientConfig()
	if err != nil {
		klog.Errorln("Failed to get kubernetes config:", err)
	}

	kc, err := kclient.New(config, kclient.Options{})
	if err != nil {
		klog.Fatalln("Failed to create Kubernetes client:", err)
	}

	var pod v1.Pod
	err = kc.Get(ctx, types.NamespacedName{Name: podName, Namespace: namespace}, &pod)
	if err != nil {
		klog.Fatalf("look up owner(s) of pod %s/%s: %v", namespace, podName, err)
	}
	owner := v1.ObjectReference{
		APIVersion: "v1",
		Kind:       "Pod",
		Name:       pod.GetName(),
		Namespace:  pod.GetNamespace(),
		UID:        pod.GetUID(),
	}
	klog.Infof("using %s/%s as owner for Kubernetes events", owner.Kind, owner.Name)

	lc, err := lclient.NewClient()
	if err != nil {
		klog.Fatalln("failed to create LINSTOR client:", err)
	}

	klog.Infof("Starting main loop (delay: %d seconds)", *delaySeconds)
	candidatesChannel := runCandidatesLoop(ctx, getCandidates, time.Duration(*delaySeconds)*time.Second)

	for storagePool := range candidatesChannel {
		klog.Infof("Got %s candidate: %+v\n", storagePool.ProviderKind, storagePool)
		_, err = lc.Nodes.Get(ctx, nodeName)
		if err != nil {
			klog.Fatalln("Failed to get LINSTOR node", err)
		}
		_, err = lc.Nodes.GetStoragePool(ctx, nodeName, storagePool.StoragePoolName)
		if err != nil && err == lclient.NotFoundError {
			err = lc.Nodes.CreateStoragePool(ctx, nodeName, storagePool)
			if err != nil {
				sendKubernetesEvent(ctx, kc, &owner, v1.EventTypeWarning, "Failed", "Failed to create LINSTOR storage pool: "+err.Error())
				klog.Fatalln("Failed to create LINSTOR storage pool", err)
			}
			sendKubernetesEvent(ctx, kc, &owner, v1.EventTypeNormal, "Created", "Created LINSTOR storage pool: "+nodeName+"/"+storagePool.StoragePoolName)
		} else if err != nil {
			klog.Fatalln("Failed to get LINSTOR storage pool", err)
		}

		// Get the maximum number of available replicas
		opts := lclient.ListOpts{
			StoragePool: []string{storagePool.StoragePoolName},
			Limit:       maxReplicasNum,
		}
		sps, err := lc.Nodes.GetStoragePoolView(ctx, &opts)
		if err != nil {
			klog.Fatalln("Failed to list LINSTOR storage pools", err)
		}
		replicasNum := len(sps)
		if replicasNum > maxReplicasNum {
			replicasNum = maxReplicasNum
		}

		// Create StorageClasses in Kubernetes
		for r := 1; r <= replicasNum; r++ {
			storageClass := newKubernetesStorageClass(&storagePool, r)
			err = kc.Get(ctx, types.NamespacedName{Name: storageClass.GetName()}, &storageClass)
			if err != nil && errors.IsNotFound(err) {
				err = kc.Create(ctx, &storageClass)
				if err != nil {
					sendKubernetesEvent(ctx, kc, &owner, v1.EventTypeWarning, "Failed", "Failed to create Kubernetes storage class: "+err.Error())
					klog.Fatalln("Failed to create Kubernetes storageClass", err)
				}
				sendKubernetesEvent(ctx, kc, &owner, v1.EventTypeNormal, "Created", "Created Kubernetes storage class: "+storageClass.Name)
			} else if err != nil {
				klog.Fatalln("Failed to get Kubernetes storage class:", err)
			}
		}
	}
}

// Sends event to Kubernetes
func sendKubernetesEvent(ctx context.Context, kc kclient.Client, owner *v1.ObjectReference, eventType, reason, message string) {
	event := newKubernetesEvent(owner, eventType, reason, message)
	err := kc.Create(ctx, &event)
	if err != nil {
		klog.Fatalln("Failed to create event", err)
	}
}

// Makes loop over storage pool candidates, retruns channel of changed ones
func runCandidatesLoop(ctx context.Context, f func() ([]lclient.StoragePool, error), delay time.Duration) <-chan lclient.StoragePool {
	ch := make(chan lclient.StoragePool)
	go func() {
		var oldCandidates []lclient.StoragePool
		for {
			select {
			case <-ctx.Done():
				close(ch)
				return
			default:
				candidates, err := f()
				if err != nil {
					klog.Fatalln("failed to load candidates:", err)
				}
			LOOP:
				for _, candidate := range candidates {
					for _, oldCandidate := range oldCandidates {
						if candidate.StoragePoolName == oldCandidate.StoragePoolName {
							continue LOOP
						}
					}
					ch <- candidate
				}
				oldCandidates = candidates
			}
			time.Sleep(delay)
		}
	}()
	return ch
}

type VolumeGroups struct{}

func getLVMThinStoragePools() ([]lclient.StoragePool, error) {
	cmd := exec.Command("vgs", "-oname,tags", "--separator=;", "--noheadings", "--config="+lvmConfig)

	var out bytes.Buffer
	cmd.Stdout = &out
	err := cmd.Run()
	if err != nil {
		return nil, err
	}
	return parseVolumeGroups(out.String())
}

type ThinPools struct{}

func getLVMStoragePools() ([]lclient.StoragePool, error) {
	cmd := exec.Command("lvs", "-oname,vg_name,lv_attr,tags", "--separator=;", "--noheadings", "--config="+lvmConfig)
	var out bytes.Buffer
	cmd.Stdout = &out
	err := cmd.Run()
	if err != nil {
		return nil, err
	}
	return parseThinPools(out.String())
}

// Collects all storage pool candidates from the node
func getCandidates() ([]lclient.StoragePool, error) {
	var candidates []lclient.StoragePool

	// Getting LVM storage pools
	cs, err := getLVMStoragePools()
	if err != nil {
		return nil, fmt.Errorf("failed to get LVM storage pools: %s", err)
	}
	candidates = append(candidates, cs...)

	// Getting LVM thin storage pools
	cs, err = getLVMThinStoragePools()
	if err != nil {
		return nil, fmt.Errorf("failed to get LVMthin storage pools: %s", err)
	}
	candidates = append(candidates, cs...)

	return candidates, nil
}

func parseThinPools(out string) ([]lclient.StoragePool, error) {
	var sps []lclient.StoragePool
	for _, line := range strings.Split(out, "\n") {
		if line == "" {
			continue
		}
		a := strings.Split(line, ";")
		if len(a) != 4 {
			return nil, fmt.Errorf("wrong line: \"%s\"", line)
		}
		lvName, vgName, lvAttr, tags := strings.TrimSpace(a[0]), a[1], a[2], strings.Split(a[3], ",")
		if lvName == "" {
			return nil, fmt.Errorf("LV name can't be empty (line: \"%s\")", line)
		}
		if vgName == "" {
			return nil, fmt.Errorf("vgName can't be empty (line: \"%s\")", line)
		}
		if lvAttr == "" {
			return nil, fmt.Errorf("lvAttr can't be empty (line: \"%s\")", line)
		}
		if lvAttr[0:1] != "t" {
			continue
		}
		name := parseNameFromLVMTags(&tags)
		if name == "" {
			continue
		}
		sps = append(sps, lclient.StoragePool{
			StoragePoolName: name,
			NodeName:        nodeName,
			ProviderKind:    lclient.LVM_THIN,
			Props: map[string]string{
				"StorDriver/LvmVg":    vgName,
				"StorDriver/ThinPool": lvName,
			},
		})
	}
	return sps, nil
}

func parseVolumeGroups(out string) ([]lclient.StoragePool, error) {
	var sps []lclient.StoragePool
	for _, line := range strings.Split(out, "\n") {
		if line == "" {
			continue
		}
		a := strings.Split(line, ";")
		if len(a) != 2 {
			return nil, fmt.Errorf("wrong line: %s", line)
		}
		vgName, tags := strings.TrimSpace(a[0]), strings.Split(a[1], ",")
		if vgName == "" {
			return nil, fmt.Errorf("VG name can't be empty (line: \"%s\")", line)
		}
		name := parseNameFromLVMTags(&tags)
		if name == "" {
			continue
		}
		sps = append(sps, lclient.StoragePool{
			StoragePoolName: name,
			NodeName:        nodeName,
			ProviderKind:    lclient.LVM,
			Props: map[string]string{
				"StorDriver/LvmVg": vgName,
			},
		})
	}
	return sps, nil
}

func parseNameFromLVMTags(tags *[]string) string {
	for _, tag := range *tags {
		t := strings.Split(tag, "-")
		if t[0] == linstorPrefix && t[1] != "" {
			return t[1]
		}
	}
	return ""
}

func newKubernetesEvent(owner *v1.ObjectReference, eventType, reason, message string) v1.Event {
	eventTime := metav1.Now()
	event := v1.Event{
		ObjectMeta: metav1.ObjectMeta{
			Namespace:    owner.Namespace,
			GenerateName: owner.Name,
			Labels: map[string]string{
				"app": "linstor-pools-importer",
			},
		},
		Reason:         reason,
		Message:        message,
		InvolvedObject: *owner,
		Source: v1.EventSource{
			Component: "linstor-pools-importer",
			Host:      nodeName,
		},
		Count:          1,
		FirstTimestamp: eventTime,
		LastTimestamp:  eventTime,
		Type:           v1.EventTypeNormal,
	}
	return event
}

func newKubernetesStorageClass(sp *lclient.StoragePool, r int) storagev1.StorageClass {
	volBindMode := storagev1.VolumeBindingImmediate
	allowVolumeExpansion := true
	reclaimPolicy := v1.PersistentVolumeReclaimDelete
	return storagev1.StorageClass{
		ObjectMeta: metav1.ObjectMeta{
			Name: fmt.Sprintf("%s-%s-r%d", linstorPrefix, sp.StoragePoolName, r),
		},
		Provisioner:          "linstor.csi.linbit.com",
		VolumeBindingMode:    &volBindMode,
		AllowVolumeExpansion: &allowVolumeExpansion,
		ReclaimPolicy:        &reclaimPolicy,
		Parameters: map[string]string{
			"linstor.csi.linbit.com/storagePool":    sp.StoragePoolName,
			"linstor.csi.linbit.com/placementCount": fmt.Sprintf("%d", r),
		},
	}
}
