/*
 Copyright 2020 The Kubernetes Authors.

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

package config

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

/*
	TODO:
	When the INI based cloud-config is deprecated. This file should be deleted.
*/

func TestReadINIConfig(t *testing.T) {
	contents := `
[LoadBalancer]
ip-pool-name = pool1
size = MEDIUM
lb-service-id = 4711
tier1-gateway-path = 1234
tcp-app-profile-name = default-tcp-lb-app-profile
udp-app-profile-name = default-udp-lb-app-profile
snat-disabled = false
tags = {\"tag1\": \"value1\", \"tag2\": \"value 2\"}

[LoadBalancerClass "public"]
ip-pool-name = poolPublic

[LoadBalancerClass "private"]
ip-pool-name = poolPrivate
tcp-app-profile-name = tcp2
udp-app-profile-name = udp2
`
	config, err := ReadRawConfigINI([]byte(contents))
	if err != nil {
		t.Error(err)
		return
	}

	assertEquals := func(name, left, right string) {
		if left != right {
			t.Errorf("%s %s != %s", name, left, right)
		}
	}
	assertEquals("LoadBalancer.ipPoolName", config.LoadBalancer.IPPoolName, "pool1")
	assertEquals("LoadBalancer.lbServiceId", config.LoadBalancer.LBServiceID, "4711")
	assertEquals("LoadBalancer.tier1GatewayPath", config.LoadBalancer.Tier1GatewayPath, "1234")
	assertEquals("LoadBalancer.tcpAppProfileName", config.LoadBalancer.TCPAppProfileName, "default-tcp-lb-app-profile")
	assertEquals("LoadBalancer.udpAppProfileName", config.LoadBalancer.UDPAppProfileName, "default-udp-lb-app-profile")
	assertEquals("LoadBalancer.size", config.LoadBalancer.Size, "MEDIUM")
	assert.Equal(t, false, config.LoadBalancer.SnatDisabled)
	if len(config.LoadBalancerClass) != 2 {
		t.Errorf("expected two LoadBalancerClass subsections, but got %d", len(config.LoadBalancerClass))
	}
	assertEquals("LoadBalancerClass.public.ipPoolName", config.LoadBalancerClass["public"].IPPoolName, "poolPublic")
	assertEquals("LoadBalancerClass.private.tcpAppProfileName", config.LoadBalancerClass["private"].TCPAppProfileName, "tcp2")
	assertEquals("LoadBalancerClass.private.udpAppProfileName", config.LoadBalancerClass["private"].UDPAppProfileName, "udp2")
	if len(config.LoadBalancer.AdditionalTags) != 2 || config.LoadBalancer.AdditionalTags["tag1"] != "value1" || config.LoadBalancer.AdditionalTags["tag2"] != "value 2" {
		t.Errorf("unexpected additionalTags %v", config.LoadBalancer.AdditionalTags)
	}
}

func TestReadINIConfigOnVMC(t *testing.T) {
	contents := `
[LoadBalancer]
ip-pool-id = 123-456
size = MEDIUM
tier1-gateway-path = 1234
tcp-app-profile-path = infra/xxx/tcp1234
udp-app-profile-path = infra/xxx/udp1234
snat-disabled = false
`
	config, err := ReadRawConfigINI([]byte(contents))
	if err != nil {
		t.Error(err)
		return
	}
	assertEquals := func(name, left, right string) {
		if left != right {
			t.Errorf("%s %s != %s", name, left, right)
		}
	}
	assertEquals("LoadBalancer.ip-pool-id", config.LoadBalancer.IPPoolID, "123-456")
	assertEquals("LoadBalancer.size", config.LoadBalancer.Size, "MEDIUM")
	assertEquals("LoadBalancer.tier1-gateway-path", config.LoadBalancer.Tier1GatewayPath, "1234")
	assertEquals("LoadBalancer.tcp-app-profile-path", config.LoadBalancer.TCPAppProfilePath, "infra/xxx/tcp1234")
	assertEquals("LoadBalancer.udp-app-profile-path", config.LoadBalancer.UDPAppProfilePath, "infra/xxx/udp1234")
	assert.Equal(t, false, config.LoadBalancer.SnatDisabled)
}
