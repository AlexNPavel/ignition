// Copyright 2017 CoreOS, Inc.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

// The vagrant-virtualbox provider fetches the configuration from raw data on a partition
// with the GUID 99570a8a-f826-4eb0-ba4e-9dd72d55ea13

package vagrant_virtualbox

import (
	"bytes"
	"io/ioutil"
	"os"
	"time"

	"github.com/coreos/ignition/config"
	"github.com/coreos/ignition/config/types"
	"github.com/coreos/ignition/config/validate/report"
	"github.com/coreos/ignition/internal/providers/util"
	"github.com/coreos/ignition/internal/resource"
)

const (
	configPath = "/dev/disk/by-partuuid/99570a8a-f826-4eb0-ba4e-9dd72d55ea13"
)

func FetchConfig(f resource.Fetcher) (types.Config, report.Report, error) {
	f.Logger.Debug("Attempting to read config drive")
	var err error
	var rawConfig []byte
	rawConfig = nil
	i := 0
	for i = 0; i < 30; i++ {
		rawConfig, err = ioutil.ReadFile(configPath)
		if os.IsNotExist(err) {
			f.Logger.Info("Path to ignition config does not exist, waiting 1 second")
			time.Sleep(1)
		} else if err != nil {
			f.Logger.Err("Error reading ignition config: %v", err)
			return types.Config{}, report.Report{}, err
		} else {
			break
		}
	}
	if rawConfig != nil {
		trimmedConfig := bytes.TrimRight(rawConfig, "\x00")
		return util.ParseConfig(f.Logger, trimmedConfig)
	} else {
		f.Logger.Info("Path to ignition config does not exist after 30 seconds, assuming no config")
		return types.Config{}, report.Report{}, config.ErrEmpty
	}
}