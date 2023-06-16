/*-
 * #%L
 * OBKV Table Client Framework
 * %%
 * Copyright (C) 2021 OceanBase
 * %%
 * OBKV Table Client Framework is licensed under Mulan PSL v2.
 * You can use this software according to the terms and conditions of the Mulan PSL v2.
 * You may obtain a copy of Mulan PSL v2 at
 *          http//license.coscl.org.cn/MulanPSL2
 * THIS SOFTWARE IS PROVIDED ON AN "AS IS" BASIS, WITHOUT WARRANTIES OF ANY KIND,
 * EITHER EXPRESS OR IMPLIED, INCLUDING BUT NOT LIMITED TO NON-INFRINGEMENT,
 * MERCHANTABILITY OR FIT FOR A PARTICULAR PURPOSE.
 * See the Mulan PSL v2 for more details.
 * #L%
 */

package protocol

import (
	"bytes"

	"github.com/oceanbase/obkv-table-client-go/util"
)

type ObTableQueryRequest struct {
	ObUniVersionHeader
	ObPayloadBase
	credential           []byte
	tableName            string
	tableId              uint64
	partitionId          uint64
	entityType           ObTableEntityType
	consistencyLevel     ObTableConsistencyLevel
	tableQuery           *ObTableQuery
	returnRowKey         bool
	returnAffectedEntity bool
	returnAffectedRows   bool
}

func (r *ObTableQueryRequest) TableName() string {
	return r.tableName
}

func (r *ObTableQueryRequest) SetTableName(tableName string) {
	r.tableName = tableName
}

func (r *ObTableQueryRequest) TableId() uint64 {
	return r.tableId
}

func (r *ObTableQueryRequest) SetTableId(tableId uint64) {
	r.tableId = tableId
}

func (r *ObTableQueryRequest) PartitionId() uint64 {
	return r.partitionId
}

func (r *ObTableQueryRequest) SetPartitionId(partitionId uint64) {
	r.partitionId = partitionId
}

func (r *ObTableQueryRequest) EntityType() ObTableEntityType {
	return r.entityType
}

func (r *ObTableQueryRequest) SetEntityType(entityType ObTableEntityType) {
	r.entityType = entityType
}

func (r *ObTableQueryRequest) ConsistencyLevel() ObTableConsistencyLevel {
	return r.consistencyLevel
}

func (r *ObTableQueryRequest) SetConsistencyLevel(consistencyLevel ObTableConsistencyLevel) {
	r.consistencyLevel = consistencyLevel
}

func (r *ObTableQueryRequest) TableQuery() *ObTableQuery {
	return r.tableQuery
}

func (r *ObTableQueryRequest) SetTableQuery(tableQuery *ObTableQuery) {
	r.tableQuery = tableQuery
}

func (r *ObTableQueryRequest) ReturnRowKey() bool {
	return r.returnRowKey
}

func (r *ObTableQueryRequest) SetReturnRowKey(returnRowKey bool) {
	r.returnRowKey = returnRowKey
}

func (r *ObTableQueryRequest) ReturnAffectedEntity() bool {
	return r.returnAffectedEntity
}

func (r *ObTableQueryRequest) SetReturnAffectedEntity(returnAffectedEntity bool) {
	r.returnAffectedEntity = returnAffectedEntity
}

func (r *ObTableQueryRequest) ReturnAffectedRows() bool {
	return r.returnAffectedRows
}

func (r *ObTableQueryRequest) SetReturnAffectedRows(returnAffectedRows bool) {
	r.returnAffectedRows = returnAffectedRows
}

func (r *ObTableQueryRequest) PCode() ObTablePacketCode {
	return ObTableApiExecuteQuery
}

func (r *ObTableQueryRequest) PayloadLen() int {
	return r.PayloadContentLen() + r.ObUniVersionHeader.UniVersionHeaderLen() // Do not change the order
}

func (r *ObTableQueryRequest) PayloadContentLen() int {
	totalLen := 0
	if util.ObVersion() >= 4 {
		totalLen =
			util.EncodedLengthByBytesString(r.credential) +
				util.EncodedLengthByVString(r.tableName) +
				util.EncodedLengthByVi64(int64(r.tableId)) +
				8 + // partitionId
				2 + // entityType consistencyLevel
				r.tableQuery.PayloadLen()
	} else {
		totalLen =
			util.EncodedLengthByBytesString(r.credential) +
				util.EncodedLengthByVString(r.tableName) +
				util.EncodedLengthByVi64(int64(r.tableId)) +
				util.EncodedLengthByVi64(int64(r.partitionId)) + // partitionId
				2 + // entityType consistencyLevel
				r.tableQuery.PayloadLen()
	}
	r.ObUniVersionHeader.SetContentLength(totalLen)
	return r.ObUniVersionHeader.ContentLength()
}

func (r *ObTableQueryRequest) Credential() []byte {
	return r.credential
}

func (r *ObTableQueryRequest) SetCredential(credential []byte) {
	r.credential = credential
}

func (r *ObTableQueryRequest) Encode(buffer *bytes.Buffer) {
	r.ObUniVersionHeader.Encode(buffer)

	util.EncodeBytesString(buffer, r.credential)

	util.EncodeVString(buffer, r.tableName)

	util.EncodeVi64(buffer, int64(r.tableId))

	if util.ObVersion() >= 4 {
		util.PutUint64(buffer, r.partitionId)
	} else {
		util.EncodeVi64(buffer, int64(r.partitionId))
	}

	util.PutUint8(buffer, uint8(r.entityType))

	util.PutUint8(buffer, uint8(r.consistencyLevel))

	r.tableQuery.Encode(buffer)

}

func (r *ObTableQueryRequest) Decode(buffer *bytes.Buffer) {
	r.ObUniVersionHeader.Decode(buffer)

	r.credential = util.DecodeBytesString(buffer)

	r.tableName = util.DecodeVString(buffer)

	r.tableId = uint64(util.DecodeVi64(buffer))

	if util.ObVersion() >= 4 {
		r.partitionId = util.Uint64(buffer)
	} else {
		r.partitionId = uint64(util.DecodeVi64(buffer))
	}

	r.entityType = ObTableEntityType(util.Uint8(buffer))

	r.consistencyLevel = ObTableConsistencyLevel(util.Uint8(buffer))

	r.tableQuery.Decode(buffer)
}
