// Copyright 2024 Google LLC
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package winproducts

import "golang.org/x/sys/windows/registry"

const (
	// Registry information to find the flavor of Windows.
	regRoot = `SOFTWARE\Microsoft\Windows NT\CurrentVersion`
	regKey  = "InstallationType"
)

// WindowsFlavorFromRegistry returns the Windows flavor (e.g. server, client) from the registry.
// It will default to a "server" flavor if it cannot be determined.
func WindowsFlavorFromRegistry() string {
	k, err := registry.OpenKey(registry.LOCAL_MACHINE, regRoot, registry.QUERY_VALUE)
	if err != nil {
		return windowsFlavor("server")
	}
	defer k.Close()

	value, _, err := k.GetStringValue(regKey)
	if err != nil {
		return windowsFlavor("server")
	}

	return windowsFlavor(value)
}
