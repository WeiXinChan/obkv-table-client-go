package route

import (
	"strconv"
	"strings"
	"time"

	"github.com/pkg/errors"

	"github.com/oceanbase/obkv-table-client-go/protocol"
	"github.com/oceanbase/obkv-table-client-go/table"
	"github.com/oceanbase/obkv-table-client-go/util"
)

type obKeyPartDesc struct {
	obPartDescCommon
	partSpace     int
	partNum       int
	partNameIdMap map[string]int64
}

func (d *obKeyPartDesc) SetPartNum(partNum int) {
	d.partNum = partNum
}

func newObKeyPartDesc() *obKeyPartDesc {
	return &obKeyPartDesc{}
}

func (d *obKeyPartDesc) partFuncType() obPartFuncType {
	return d.PartFuncType
}

func (d *obKeyPartDesc) orderedPartColumnNames() []string {
	return d.OrderedPartColumnNames
}

func (d *obKeyPartDesc) setOrderedPartColumnNames(partExpr string) {
	// eg:"c1, c2", need to remove ' '
	str := strings.ReplaceAll(partExpr, " ", "")
	d.OrderedPartColumnNames = strings.Split(str, ",")
}

func (d *obKeyPartDesc) orderedPartRefColumnRowKeyRelations() []*obColumnIndexesPair {
	return d.OrderedPartRefColumnRowKeyRelations
}

func (d *obKeyPartDesc) rowKeyElement() *table.ObRowKeyElement {
	return d.RowKeyElement
}

func (d *obKeyPartDesc) setRowKeyElement(rowKeyElement *table.ObRowKeyElement) {
	d.setCommRowKeyElement(rowKeyElement)
}

func (d *obKeyPartDesc) setPartColumns(partColumns []*obColumn) {
	d.PartColumns = partColumns
}

func (d *obKeyPartDesc) GetPartId(rowKey []interface{}) (int64, error) {
	if len(rowKey) == 0 {
		return ObInvalidPartId, errors.New("rowKey size is 0")
	}
	evalValues, err := evalPartKeyValues(d, rowKey)
	if err != nil {
		return ObInvalidPartId, errors.WithMessagef(err, "eval part key value, part desc:%s", d.String())
	}
	if len(evalValues) < len(d.OrderedPartRefColumnRowKeyRelations) {
		return ObInvalidPartId, errors.Errorf("invalid eval values length, "+
			"evalValues length:%d, OrderedPartRefColumnRowKeyRelations length: %d", len(evalValues), len(d.OrderedPartRefColumnRowKeyRelations))
	}
	var hashValue int64
	for i := 0; i < len(d.OrderedPartRefColumnRowKeyRelations); i++ {
		hashValue, err = d.toHashCode(
			evalValues[i],
			d.OrderedPartRefColumnRowKeyRelations[i].column,
			hashValue,
			d.PartFuncType,
		)
		if err != nil {
			return ObInvalidPartId, errors.WithMessagef(err, "convert to hash code, part desc:%s", d.String())
		}
	}
	if hashValue < 0 {
		hashValue = -hashValue
	}
	return (int64(d.partSpace) << ObPartIdBitNum) | (hashValue % int64(d.partNum)), nil
}

func intToInt64(value interface{}) (int64, error) {
	switch v := value.(type) {
	case bool:
		if v {
			return int64(1), nil
		} else {
			return int64(0), nil
		}
	case int8:
		return int64(v), nil
	case uint8:
		return int64(v), nil
	case int16:
		return int64(v), nil
	case uint16:
		return int64(v), nil
	case int:
		return int64(v), nil
	case uint:
		return int64(v), nil
	case int32:
		return int64(v), nil
	case uint32:
		return int64(v), nil
	case int64:
		return v, nil
	default:
		return -1, errors.Errorf("invalid type to convert to int64， value：%T", value)
	}
}

func (d *obKeyPartDesc) toHashCode(
	value interface{},
	refColumn *obColumn,
	hashCode int64,
	partFuncType obPartFuncType) (int64, error) {
	typeValue := refColumn.objType.GetValue()
	if typeValue >= protocol.ObTinyIntTypeValue && typeValue <= protocol.ObUInt64TypeValue {
		i64, err := intToInt64(value)
		if err != nil {
			return -1, errors.WithMessagef(err, "convert int to int64, value:%T", typeValue)
		}
		arr := d.longToByteArray(i64)
		return murmurHash64A(arr, len(arr), hashCode), nil
	} else if typeValue == protocol.ObDateTimeTypeValue || typeValue == protocol.ObTimestampTypeValue {
		t, ok := value.(time.Time)
		if !ok {
			return -1, errors.Errorf("invalid timestamp type, value:%T", value)
		}
		return d.timeStampHash(t, hashCode), nil
	} else if typeValue == protocol.ObDateTypeValue {
		date, ok := value.(time.Time)
		if !ok {
			return -1, errors.Errorf("invalid date type, value:%T", value)
		}
		return d.dateHash(date, hashCode), nil
	} else if typeValue == protocol.ObVarcharTypeValue || typeValue == protocol.ObCharTypeValue {
		return d.varcharHash(value, refColumn.collationType, hashCode, partFuncType)
	} else {
		return -1, errors.Errorf("unsupported type for key hash, objType:%s", refColumn.objType.String())
	}
}

func (d *obKeyPartDesc) longToByteArray(l int64) []byte {
	return []byte{(byte)(l & 0xFF), (byte)((l >> 8) & 0xFF), (byte)((l >> 16) & 0xFF),
		(byte)((l >> 24) & 0xFF), (byte)((l >> 32) & 0xFF), (byte)((l >> 40) & 0xFF),
		(byte)((l >> 48) & 0xFF), (byte)((l >> 56) & 0xFF)}
}

func (d *obKeyPartDesc) longHash(value int64, hashCode int64) int64 {
	arr := d.longToByteArray(value)
	return murmurHash64A(arr, len(arr), hashCode)
}

func (d *obKeyPartDesc) timeStampHash(ts time.Time, hashCode int64) int64 {
	return d.longHash(ts.UnixMilli(), hashCode)
}

func (d *obKeyPartDesc) dateHash(ts time.Time, hashCode int64) int64 {
	return d.longHash(ts.UnixMilli(), hashCode)
}

func (d *obKeyPartDesc) varcharHash(
	value interface{},
	collType protocol.ObCollationType,
	hashCode int64,
	partFuncType obPartFuncType) (int64, error) {
	var seed uint64 = 0xc6a4a7935bd1e995
	var bytes []byte
	if v, ok := value.(string); ok {
		// Right Now, only UTF8 String is supported, aligned with the Serialization.
		// string and []byte is utf8 default in go language
		bytes = []byte(v)
	} else if v, ok := value.([]byte); ok {
		bytes = v
	} else if v, ok := value.(protocol.ObBytesString); ok {
		bytes = v.BytesVal()
	} else {
		return -1, errors.Errorf("invalid varchar value for calc hash value, value:%T", value)
	}
	switch collType.Value() {
	case protocol.CsTypeUtf8mb4GeneralCi:
		if partFuncType == partFuncTypeKeyV3 ||
			partFuncType == partFuncTypeKeyImplV2 ||
			util.ObVersion() >= 4 {
			hashCode = hashSortUtf8Mb4(bytes, hashCode, seed, true)
		} else {
			hashCode = hashSortUtf8Mb4(bytes, hashCode, seed, false)
		}
	case protocol.CsTypeUtf8mb4Bin:
	case protocol.CsTypeBinary:
		if partFuncType == partFuncTypeKeyV3 ||
			partFuncType == partFuncTypeKeyImplV2 ||
			util.ObVersion() >= 4 {
			hashCode = murmurHash64A(bytes, len(bytes), hashCode)
		} else {
			hashCode = hashSortMbBin(bytes, hashCode, seed)
		}
	case protocol.CsTypeInvalid:
	case protocol.CsTypeCollationFree:
	case protocol.CsTypeMax:
		return -1, errors.Errorf("not supported collation type, collType:%d", collType.Value())
	}
	return hashCode, nil
}

func (d *obKeyPartDesc) String() string {
	// partNameIdMap to string
	var partNameIdMapStr string
	partNameIdMapStr = partNameIdMapStr + "{"
	var i = 0
	for k, v := range d.partNameIdMap {
		if i > 0 {
			partNameIdMapStr += ", "
		}
		i++
		partNameIdMapStr += "m[" + k + "]=" + strconv.Itoa(int(v))
	}
	partNameIdMapStr += "}"
	return "obKeyPartDesc{" +
		"comm:" + d.CommString() + ", " +
		"partSpace:" + strconv.Itoa(d.partSpace) + ", " +
		"partNum:" + strconv.Itoa(d.partNum) + ", " +
		"partNameIdMap:" + partNameIdMapStr +
		"}"
}
