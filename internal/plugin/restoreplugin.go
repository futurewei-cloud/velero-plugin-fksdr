/*
Copyright 2018, 2019 the Velero contributors.

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

package plugin

import (
	"github.com/sirupsen/logrus"
	"github.com/vmware-tanzu/velero/pkg/plugin/velero"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
)

// RestorePlugin is a restore item action plugin for Velero
type RestorePlugin struct {
	log logrus.FieldLogger
}

// NewRestorePlugin instantiates a RestorePlugin.
func NewRestorePlugin(log logrus.FieldLogger) *RestorePlugin {
	return &RestorePlugin{log: log}
}

// AppliesTo returns information about which resources this action should be invoked for.
// A RestoreItemAction's Execute function will only be invoked on items that match the returned
// selector. A zero-valued ResourceSelector matches all resources.g
func (p *RestorePlugin) AppliesTo() (velero.ResourceSelector, error) {
	//res := make ( []string, 1 )
	//res[0] = "PersistentVolume"
	return velero.ResourceSelector{}, nil
}

// Execute allows the RestorePlugin to perform arbitrary logic with the item being restored,
// in this case, setting a custom annotation on the item being restored.
func (p *RestorePlugin) Execute(input *velero.RestoreItemActionExecuteInput) (*velero.RestoreItemActionExecuteOutput, error) {
	p.log.Info("FKSDR RestorePlugin!")

	obj, ok := input.Item.(*unstructured.Unstructured)

	if !ok {
		return velero.NewRestoreItemActionExecuteOutput(nil), nil
	}

	log := p.log.WithFields(logrus.Fields{
		"TonyContent": obj.Object,
		"TonyKind":    obj.GetKind(),
	})

	log.Info("FKSDR Restore Item++++++")

	if obj.GetKind() == "PersistentVolume" {

		annotations := obj.GetAnnotations()

		if annotations == nil {
			annotations = make(map[string]string)
		}

		annotations["velero.io/my-restore-plugin"] = "1"

		obj.SetAnnotations(annotations)

		if newClusterID, ok := annotations["tony.io/dr-protected-pv"]; ok {

			if newClusterID == "Huawei-DR-Protected" {
				annotations["velero.io/my-restore-plugin"] = "Huawei-DR-Restore"
			} else {
				//Replace source clusterID with target cluster ID
				unstructured.SetNestedField(obj.Object, newClusterID, "spec", "csi", "volumeAttributes", "clusterID")
				annotations["velero.io/my-restore-plugin"] = "Ceph-DR-Restore"
			}

			log := p.log.WithFields(logrus.Fields{
				"Kind":    obj.GetKind(),
				"Content": obj.UnstructuredContent(),
			})
			log.Info("FKSDR Restore Item -- Made Changes!!!")
		}
	}

	log = p.log.WithFields(logrus.Fields{
		"TonyContent": obj.Object,
		"TonyKind":    obj.GetKind(),
	})

	log.Info("FKSDR After Restore Item------")

	return velero.NewRestoreItemActionExecuteOutput(obj), nil
}
