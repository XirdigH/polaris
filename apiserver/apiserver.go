/**
 * Tencent is pleased to support the open source community by making Polaris available.
 *
 * Copyright (C) 2019 THL A29 Limited, a Tencent company. All rights reserved.
 *
 * Licensed under the BSD 3-Clause License (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 * https://opensource.org/licenses/BSD-3-Clause
 *
 * Unless required by applicable law or agreed to in writing, software distributed
 * under the License is distributed on an "AS IS" BASIS, WITHOUT WARRANTIES OR
 * CONDITIONS OF ANY KIND, either express or implied. See the License for the
 * specific language governing permissions and limitations under the License.
 */

package apiserver

import (
	"context"
	"errors"
	"fmt"
	"github.com/polarismesh/polaris-server/common/log"
	"strings"
)

const (
	DiscoverAccess    string = "discover"
	RegisterAccess    string = "register"
	HealthcheckAccess string = "healthcheck"
)

/**
 * @brief API服务器配置
 */
type Config struct {
	Name   string
	Option map[string]interface{}
	API    map[string]APIConfig
}

/**
 * @brief API配置
 */
type APIConfig struct {
	Enable  bool
	Include []string
}

/**
 * @brief API服务器接口
 */
type Apiserver interface {
	GetProtocol() string
	GetPort() uint32
	Initialize(ctx context.Context, option map[string]interface{}, api map[string]APIConfig) error
	Run(errCh chan error)
	Stop()
	Restart(option map[string]interface{}, api map[string]APIConfig, errCh chan error) error
}

var (
	Slots = make(map[string]Apiserver)
)

/**
 * @brief 注册API服务器
 */
func Register(name string, server Apiserver) error {
	if _, exist := Slots[name]; exist {
		err := errors.New("apiserver name exist")
		return err
	}

	Slots[name] = server

	return nil
}

/**
 * @brief 获取客户端openMethod
 */
func GetClientOpenMethod(include []string, protocol string) (map[string]bool, error) {
	clientAccess := make(map[string][]string)
	clientAccess[DiscoverAccess] = []string{"Discover", "ReportClient"}
	clientAccess[RegisterAccess] = []string{"RegisterInstance", "DeregisterInstance"}
	clientAccess[HealthcheckAccess] = []string{"Heartbeat"}

	openMethod := make(map[string]bool)
	// 如果为空，开启全部接口
	if len(include) == 0 {
		for key := range clientAccess {
			include = append(include, key)
		}
	}

	for _, item := range include {
		if methods, ok := clientAccess[item]; ok {
			for _, method := range methods {
				method = "/v1.Polaris" + strings.ToUpper(protocol) + "/" + method
				openMethod[method] = true
			}
		} else {
			log.Errorf("method %s does not exist in %sserver client access", item, protocol)
			return nil, fmt.Errorf("method %s does not exist in %sserver client access", item, protocol)
		}
	}

	return openMethod, nil
}