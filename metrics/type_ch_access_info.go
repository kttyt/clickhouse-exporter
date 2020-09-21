// Copyright 2019 Altinity Ltd and/or its affiliates. All rights reserved.
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

package metrics

type CHAccessInfo struct {
	Username string
	Password string
	Hostname string
	Port     int
}

func NewCHAccessInfo(username, password, hostname string, port int) *CHAccessInfo {
	return &CHAccessInfo{
		Username: username,
		Password: password,
		Hostname: hostname,
		Port:     port,
	}
}
