/*
Copyright 2017, 2019 the Velero contributors.

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
	"github.com/ljtbbt/fksdr-plugin/internal/plugin"
	"github.com/sirupsen/logrus"
	"github.com/vmware-tanzu/velero/pkg/plugin/framework"
)

func main() {
	framework.NewServer().
		RegisterObjectStore("fksdr.io/object-store-plugin", newObjectStorePlugin).
		RegisterVolumeSnapshotter("fksdr.io/volume-snapshotter-plugin", newNoOpVolumeSnapshotterPlugin).
		RegisterRestoreItemAction("fksdr.io/restore-plugin", newRestorePlugin).
		RegisterBackupItemAction("fksdr.io/backup-plugin", newBackupPlugin).
		Serve()
}

func newBackupPlugin(logger logrus.FieldLogger) (interface{}, error) {
	return plugin.NewBackupPlugin(logger), nil
}

func newObjectStorePlugin(logger logrus.FieldLogger) (interface{}, error) {
	return plugin.NewFileObjectStore(logger), nil
}

func newRestorePlugin(logger logrus.FieldLogger) (interface{}, error) {
	return plugin.NewRestorePlugin(logger), nil
}

func newNoOpVolumeSnapshotterPlugin(logger logrus.FieldLogger) (interface{}, error) {
	return plugin.NewNoOpVolumeSnapshotter(logger), nil
}
