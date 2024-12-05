package splash

import (
	"errors"
	"fmt"
	"strconv"
	"strings"

	"github.com/onflow/cadence"
	"github.com/onflow/cadence/common"
)

func ToFloat64(value cadence.Value) float64 {
	val, _ := strconv.ParseFloat(value.(cadence.UFix64).String(), 64)
	return val
}

func UFix64ToString(v float64) string {
	vStr := strconv.FormatFloat(v, 'f', -1, 64)
	if strings.Contains(vStr, ".") {
		return vStr
	}

	return vStr + ".0"
}

func UFix64FromFloat64(v float64) cadence.Value {
	cv, err := cadence.NewUFix64(fmt.Sprintf("%.4f", v))
	if err != nil {
		panic(err)
	}
	return cv
}

func StringToPath(path string) (cadence.Path, error) {
	var val cadence.Path
	parts := strings.Split(path, "/")
	if len(parts) != 3 || (len(parts) > 0 && len(parts[0]) > 0) {
		return val, errors.New("bad Cadence path")
	}
	if parts[1] != "private" && parts[1] != "public" && parts[1] != "storage" {
		return val, errors.New("bad domain in Cadence path")
	}
	val.Domain = common.PathDomainFromIdentifier(parts[1])
	val.Identifier = parts[2]
	return val, nil
}

func ExtractStringValueFromEvent(txResult TransactionResult, eventName, key string) string {
	for _, e := range txResult.Events {
		if e.Name == eventName {
			v := e.Fields[key]
			if v == nil {
				panic(fmt.Sprintf("key %s not found in %s", key, eventName))
			}
			switch val := v.(type) {
			case string:
				return val
			default:
				panic(fmt.Sprintf("unexpected value type for %s in %s: %T", key, eventName, v))
			}
		}
	}

	return ""
}

func ExtractUInt64ValueFromEvent(txResult TransactionResult, eventName, key string) uint64 {
	for _, e := range txResult.Events {
		if e.Name == eventName {
			v := e.Fields[key]
			if v == nil {
				panic(fmt.Sprintf("key %s not found in %s", key, eventName))
			}
			switch val := v.(type) {
			case string:
				res, err := strconv.ParseUint(val, 10, 64)
				if err != nil {
					panic(err)
				}
				return res
			case uint64:
				return val
			default:
				panic(fmt.Sprintf("unexpected value type for %s in %s: %T", key, eventName, v))
			}
		}
	}

	panic(fmt.Sprintf("value not found for %s in %s", key, eventName))
}
