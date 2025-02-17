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
	"bytes"
	"time"

	"github.com/oceanbase/obkv-table-client-go/util"
)

type ObTableBatchOperationRequest struct {
	ObUniVersionHeader
	ObPayloadBase
	credential              []byte
	tableName               string
	tableId                 uint64
	obTableEntityType       ObTableEntityType
	obTableBatchOperation   *ObTableBatchOperation
	obTableConsistencyLevel ObTableConsistencyLevel
	returnRowKey            bool
	returnAffectedEntity    bool
	returnAffectedRows      bool
	partitionId             uint64
	atomicOperation         bool
}

func NewObTableBatchOperationRequest() *ObTableBatchOperationRequest {
	return &ObTableBatchOperationRequest{
		ObUniVersionHeader: ObUniVersionHeader{
			version:       1,
			contentLength: 0,
		},
		ObPayloadBase: ObPayloadBase{
			uniqueId:  0,
			sequence:  0,
			tenantId:  1,
			sessionId: 0,
			flag:      7,
			timeout:   10 * 1000 * time.Millisecond,
		},
		credential:              nil,
		tableName:               "",
		tableId:                 0,
		obTableEntityType:       0,
		obTableBatchOperation:   NewObTableBatchOperation(),
		obTableConsistencyLevel: 0,
		returnRowKey:            false,
		returnAffectedEntity:    false,
		returnAffectedRows:      false,
		partitionId:             0,
		atomicOperation:         true,
	}
}

func NewObTableBatchOperationRequestWithParams(
	tableName string,
	tableId uint64,
	partitionId uint64,
	obTableBatchOperation *ObTableBatchOperation,
	timeout time.Duration,
	flag uint16,
	entityType ObTableEntityType) *ObTableBatchOperationRequest {

	return &ObTableBatchOperationRequest{
		ObUniVersionHeader: ObUniVersionHeader{
			version:       1,
			contentLength: 0,
		},
		ObPayloadBase: ObPayloadBase{
			uniqueId:  0,
			sequence:  0,
			tenantId:  1,
			sessionId: 0,
			flag:      flag,
			timeout:   timeout,
		},
		credential:              nil, // when execute set
		tableName:               tableName,
		tableId:                 tableId,
		obTableEntityType:       entityType,
		obTableBatchOperation:   obTableBatchOperation,
		obTableConsistencyLevel: ObTableConsistencyLevelStrong,
		returnRowKey:            false,
		returnAffectedEntity:    false,
		returnAffectedRows:      false,
		partitionId:             partitionId,
		atomicOperation:         true,
	}
}

func (r *ObTableBatchOperationRequest) TableName() string {
	return r.tableName
}

func (r *ObTableBatchOperationRequest) SetTableName(tableName string) {
	r.tableName = tableName
}

func (r *ObTableBatchOperationRequest) TableId() uint64 {
	return r.tableId
}

func (r *ObTableBatchOperationRequest) SetTableId(tableId uint64) {
	r.tableId = tableId
}

func (r *ObTableBatchOperationRequest) ObTableEntityType() ObTableEntityType {
	return r.obTableEntityType
}

func (r *ObTableBatchOperationRequest) SetObTableEntityType(obTableEntityType ObTableEntityType) {
	r.obTableEntityType = obTableEntityType
}

func (r *ObTableBatchOperationRequest) ObTableBatchOperation() *ObTableBatchOperation {
	return r.obTableBatchOperation
}

func (r *ObTableBatchOperationRequest) SetObTableBatchOperation(obTableBatchOperation *ObTableBatchOperation) {
	r.obTableBatchOperation = obTableBatchOperation
}

func (r *ObTableBatchOperationRequest) ObTableConsistencyLevel() ObTableConsistencyLevel {
	return r.obTableConsistencyLevel
}

func (r *ObTableBatchOperationRequest) SetObTableConsistencyLevel(obTableConsistencyLevel ObTableConsistencyLevel) {
	r.obTableConsistencyLevel = obTableConsistencyLevel
}

func (r *ObTableBatchOperationRequest) ReturnRowKey() bool {
	return r.returnRowKey
}

func (r *ObTableBatchOperationRequest) SetReturnRowKey(returnRowKey bool) {
	r.returnRowKey = returnRowKey
}

func (r *ObTableBatchOperationRequest) ReturnAffectedEntity() bool {
	return r.returnAffectedEntity
}

func (r *ObTableBatchOperationRequest) SetReturnAffectedEntity(returnAffectedEntity bool) {
	r.returnAffectedEntity = returnAffectedEntity
}

func (r *ObTableBatchOperationRequest) ReturnAffectedRows() bool {
	return r.returnAffectedRows
}

func (r *ObTableBatchOperationRequest) SetReturnAffectedRows(returnAffectedRows bool) {
	r.returnAffectedRows = returnAffectedRows
}

func (r *ObTableBatchOperationRequest) PartitionId() uint64 {
	return r.partitionId
}

func (r *ObTableBatchOperationRequest) SetPartitionId(partitionId uint64) {
	r.partitionId = partitionId
}

func (r *ObTableBatchOperationRequest) AtomicOperation() bool {
	return r.atomicOperation
}

func (r *ObTableBatchOperationRequest) SetAtomicOperation(atomicOperation bool) {
	r.atomicOperation = atomicOperation
}

func (r *ObTableBatchOperationRequest) PCode() ObTablePacketCode {
	return ObTableApiBatchExecute
}

func (r *ObTableBatchOperationRequest) PayloadLen() int {
	return r.PayloadContentLen() + r.ObUniVersionHeader.UniVersionHeaderLen() // Do not change the order
}

func (r *ObTableBatchOperationRequest) PayloadContentLen() int {
	totalLen := 0
	if util.ObVersion() >= 4 {
		totalLen =
			util.EncodedLengthByBytesString(r.credential) +
				util.EncodedLengthByVString(r.tableName) +
				util.EncodedLengthByVi64(int64(r.tableId)) +
				6 + // obTableEntityType obTableConsistencyLevel returnRowKey returnAffectedEntity returnAffectedRows atomicOperation
				8 + // partitionId
				r.obTableBatchOperation.PayloadLen()
	} else {
		totalLen =
			util.EncodedLengthByBytesString(r.credential) +
				util.EncodedLengthByVString(r.tableName) +
				util.EncodedLengthByVi64(int64(r.tableId)) +
				6 + // obTableEntityType obTableConsistencyLevel returnRowKey returnAffectedEntity returnAffectedRows atomicOperation
				util.EncodedLengthByVi64(int64(r.partitionId)) + // partitionId
				r.obTableBatchOperation.PayloadLen()
	}

	r.ObUniVersionHeader.SetContentLength(totalLen)
	return r.ObUniVersionHeader.ContentLength()
}

func (r *ObTableBatchOperationRequest) Credential() []byte {
	return r.credential
}

func (r *ObTableBatchOperationRequest) SetCredential(credential []byte) {
	r.credential = credential
}

func (r *ObTableBatchOperationRequest) Encode(buffer *bytes.Buffer) {
	r.ObUniVersionHeader.Encode(buffer)

	util.EncodeBytesString(buffer, r.credential)

	util.EncodeVString(buffer, r.tableName)

	util.EncodeVi64(buffer, int64(r.tableId))

	util.PutUint8(buffer, uint8(r.obTableEntityType))

	r.obTableBatchOperation.Encode(buffer)

	util.PutUint8(buffer, uint8(r.obTableConsistencyLevel))

	util.PutUint8(buffer, util.BoolToByte(r.returnRowKey))

	util.PutUint8(buffer, util.BoolToByte(r.returnAffectedEntity))

	util.PutUint8(buffer, util.BoolToByte(r.returnAffectedRows))

	if util.ObVersion() >= 4 {
		util.PutUint64(buffer, r.partitionId)
	} else {
		util.EncodeVi64(buffer, int64(r.partitionId))
	}

	util.PutUint8(buffer, util.BoolToByte(r.atomicOperation))
}

func (r *ObTableBatchOperationRequest) Decode(buffer *bytes.Buffer) {
	r.ObUniVersionHeader.Decode(buffer)

	r.credential = util.DecodeBytesString(buffer)

	r.tableName = util.DecodeVString(buffer)

	r.tableId = uint64(util.DecodeVi64(buffer))

	r.obTableEntityType = ObTableEntityType(util.Uint8(buffer))

	r.obTableBatchOperation.Decode(buffer)

	r.obTableConsistencyLevel = ObTableConsistencyLevel(util.Uint8(buffer))

	r.returnRowKey = util.ByteToBool(util.Uint8(buffer))

	r.returnAffectedEntity = util.ByteToBool(util.Uint8(buffer))

	r.returnAffectedRows = util.ByteToBool(util.Uint8(buffer))

	if util.ObVersion() >= 4 {
		r.partitionId = util.Uint64(buffer)
	} else {
		r.partitionId = uint64(util.DecodeVi64(buffer))
	}

	r.atomicOperation = util.ByteToBool(util.Uint8(buffer))
}

func (r *ObTableBatchOperationRequest) String() string {
	return ""
}
