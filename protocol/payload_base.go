/*-
 * #%L
 * OBKV Table Client Framework
 * %%
 * Copyright (C) 2021 OceanBase
 * %%
 * OBKV Table Client Framework is licensed under Mulan PSL v2.
 * You can use this software according to the terms and conditions of the Mulan PSL v2.
 * You may obtain a copy of Mulan PSL v2 at:
 *          http://license.coscl.org.cn/MulanPSL2
 * THIS SOFTWARE IS PROVIDED ON AN "AS IS" BASIS, WITHOUT WARRANTIES OF ANY KIND,
 * EITHER EXPRESS OR IMPLIED, INCLUDING BUT NOT LIMITED TO NON-INFRINGEMENT,
 * MERCHANTABILITY OR FIT FOR A PARTICULAR PURPOSE.
 * See the Mulan PSL v2 for more details.
 * #L%
 */

package protocol

import (
	"time"
)

// ObPayloadBase payload base
type ObPayloadBase struct {
	uniqueId uint64 // rpc header traceId0
	sequence uint64 // rpc header traceId1

	flag      uint16
	tenantId  uint64
	sessionId uint64

	timeout time.Duration
}

func (p *ObPayloadBase) UniqueId() uint64 {
	return p.uniqueId
}

func (p *ObPayloadBase) SetUniqueId(uniqueId uint64) {
	p.uniqueId = uniqueId
}

func (p *ObPayloadBase) Sequence() uint64 {
	return p.sequence
}

func (p *ObPayloadBase) SetSequence(sequence uint64) {
	p.sequence = sequence
}

func (p *ObPayloadBase) Flag() uint16 {
	return p.flag
}

func (p *ObPayloadBase) SetFlag(flag uint16) {
	p.flag = flag
}

func (p *ObPayloadBase) TenantId() uint64 {
	return p.tenantId
}

func (p *ObPayloadBase) SetTenantId(tenantId uint64) {
	p.tenantId = tenantId
}

func (p *ObPayloadBase) SessionId() uint64 {
	return p.sessionId
}

func (p *ObPayloadBase) SetSessionId(sessionId uint64) {
	p.sessionId = sessionId
}

func (p *ObPayloadBase) Timeout() time.Duration {
	return p.timeout
}

func (p *ObPayloadBase) SetTimeout(timeout time.Duration) {
	p.timeout = timeout
}
