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

package route

import (
	"strconv"
	"strings"

	"github.com/pkg/errors"

	"github.com/oceanbase/obkv-table-client-go/protocol"
)

func newObSimpleColumn(
	columnName string,
	index int,
	objType protocol.ObjType,
	collType protocol.CollationType) *obColumn {
	return &obColumn{
		columnName:     columnName,
		index:          index,
		objType:        objType,
		collationType:  collType,
		refColumnNames: []string{columnName},
		isGenColumn:    false,
		columnExpress:  nil,
	}
}

type obColumn struct {
	columnName    string
	index         int
	objType       protocol.ObjType
	collationType protocol.CollationType
	// refColumnNames: Represents which columns are referenced by the current column
	// 1. generate column: key_prefix VARCHAR(1024) GENERATED ALWAYS AS (SUBSTRING(`col1`, 1, 4)
	// 					   refColumnNames = ["col1"]
	// 2. normal column: col_normal int
	//					   refColumnNames = ["col_normal"]
	refColumnNames []string
	isGenColumn    bool
	columnExpress  obGeneratedColumnSimpleFunc // only support 'SUBSTRING' expr now
}

func (c *obColumn) EvalValue(refs ...interface{}) (interface{}, error) {
	if !c.isGenColumn {
		if len(refs) == 0 || len(refs) > 1 {
			return nil, errors.Errorf("simple column is refer to itself "+
				"so that the length of the refs must be 1, len:%d", len(refs))
		}
		return c.objType.CheckTypeForValue(refs[0], c.collationType)
	} else {
		if len(refs) != len(c.refColumnNames) {
			return nil, errors.Errorf("input refs count is not equal to refColumnNames count "+
				"when column is generate column, refs len:%d, refColumnNames len:%d", len(refs), len(c.refColumnNames))
		}
		return nil, errors.New("not support generate column now")
	}
}

func (c *obColumn) String() string {
	// isGenColumn to string
	var isGenColumnStr string
	if c.isGenColumn {
		isGenColumnStr = "true"
	} else {
		isGenColumnStr = "false"
	}

	// objType to string
	var objTypeStr string
	if c.objType != nil {
		objTypeStr = c.objType.String()
	} else {
		objTypeStr = "nil"
	}

	// columnExpress to string
	var columnExpressStr string
	if c.isGenColumn {
		columnExpressStr = c.columnExpress.String()
	} else {
		columnExpressStr = "nil"
	}
	return "obColumn{" +
		"columnName:" + c.columnName + ", " +
		"index:" + strconv.Itoa(c.index) + ", " +
		"objType:" + objTypeStr + ", " +
		"collationType:" + c.collationType.String() + ", " +
		"refColumnNames:" + "[" + strings.Join(c.refColumnNames, ",") + "]" + ", " +
		"isGenColumn:" + isGenColumnStr + ", " +
		"columnExpress:" + columnExpressStr +
		"}"
}
