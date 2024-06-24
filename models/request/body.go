package request

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/kataras/iris/v12"
	"go-instaloader/utils/rlog"
	"reflect"
	"strconv"
	"strings"
)

func ReadJsonBody(ctx iris.Context, outptr interface{}) error {
	var derr error
	defer func() {
		if derr != nil {
			rlog.Error(ctx.RemoteAddr(), "ReadBody error", derr.Error())
		}
	}()

	data, err := ctx.GetBody()
	if err != nil {
		derr = err
		return errors.New("error reading body content")
	}

	if len(data) <= 0 {
		return errors.New("parameter cannot be empty")
	}

	str := string(data)
	if err := JsonUnmarshal(str, outptr); err != nil {
		derr = err
		return errors.New("parse data format error: " + err.Error())
	}

	return nil
}

// parse body to map
func GetBodyToMap(ctx iris.Context) map[string]interface{} {
	if ctx.Method() == "GET" {
		params := ctx.URLParams()
		ret := make(map[string]interface{})
		for k, v := range params {
			ret[k] = v
		}
		return ret
	}

	byteData, err := ctx.GetBody()
	if err != nil {
		rlog.Error(err)
	}
	strData := string(byteData[:])

	var data map[string]interface{}
	err = json.Unmarshal([]byte(strData), &data)
	if err != nil {
		return nil
	}
	return data
}

// String data
func GetValueString(data map[string]interface{}, key string) string {
	value, ok := data[key].(string)
	if !ok {
		return ""
	}
	return value
}

func GetValueStringDefault(data map[string]interface{}, key string, def string) string {
	value, ok := data[key].(string)
	if !ok || value == "" {
		return def
	}
	return value
}

// Bool data
func GetValueBoolDefault(data map[string]interface{}, key string, def bool) bool {
	value, ok := data[key].(bool)
	if !ok {
		return def
	}
	return value
}

// Int data
func GetValueInt(data map[string]interface{}, key string) (int, error) {
	switch data[key].(type) {
	case float64:
		value, ok := data[key].(float64)
		if !ok {
			return 0, fmt.Errorf("unable to find param value '%s'", key)
		}
		return int(value), nil
	case string:
		value, ok := data[key].(string)
		if !ok {
			return 0, fmt.Errorf("unable to find param value '%s'", key)
		}
		return strconv.Atoi(value)
	default:
		return 0, fmt.Errorf("unable to find param data type '%s'", key)
	}
}

func GetValueIntDefault(data map[string]interface{}, key string, def int) int {
	val, err := GetValueInt(data, key)
	if err != nil {
		return def
	}
	return val
}

// Int64 data
func GetValueInt64(data map[string]interface{}, key string) (int64, error) {
	switch data[key].(type) {
	case float64:
		value, ok := data[key].(float64)
		if !ok {
			return 0, fmt.Errorf("unable to find param value '%s'", key)
		}
		return int64(value), nil
	case string:
		value, ok := data[key].(string)
		if !ok {
			return 0, fmt.Errorf("unable to find param value '%s'", key)
		}
		return strconv.ParseInt(value, 10, 64)
	default:
		return 0, fmt.Errorf("unable to find param data type '%s'", key)
	}
}

func GetValueInt64Default(data map[string]interface{}, key string, def int64) int64 {
	val, err := GetValueInt64(data, key)
	if err != nil {
		return def
	}
	return val
}

func GetValueInt64Array(data map[string]interface{}, key string) []int64 {
	switch v := data[key].(type) {
	case []interface{}:
		arr, err := interfaceToInt64Array(v)
		if err != nil {
			return nil
		}
		return arr
	case string:
		arr, err := stringToInt64Array(v)
		if err != nil {
			return nil
		}
		return arr
	default:
		return nil
	}
}

// Float64 data
func GetValueFloat64(data map[string]interface{}, key string) (float64, error) {
	value, ok := data[key].(float64)
	if !ok {
		return 0, fmt.Errorf("unable to find param value '%s'", key)
	}
	return value, nil
}

func GetValueFloat64Default(data map[string]interface{}, key string, def float64) float64 {
	value, ok := data[key].(float64)
	if !ok {
		return def
	}
	return value
}

// Pagination data
func GetValuePageInfo(data map[string]interface{}) (page, limit int) {
	const maxLimit = 100
	// 判断page 类型
	p, ok := data["page"].(string)
	if !ok {
		page = 1
	}
	l, ok := data["limit"].(string)
	if !ok {
		limit = 10
	}
	page, err := strconv.Atoi(p)
	if err != nil {
		page = 1
	}
	limit, err = strconv.Atoi(l)
	if err != nil {
		limit = 10
	}
	if limit > maxLimit {
		limit = maxLimit
	}
	return
}

// Utils
func interfaceToInt64Array(intrfcArray []interface{}) ([]int64, error) {
	int64Arr := make([]int64, len(intrfcArray))
	var err error

	for i, v := range intrfcArray {
		switch val := v.(type) {
		case int:
			int64Arr[i] = int64(val)
		case int64:
			int64Arr[i] = val
		case float64:
			int64Arr[i] = int64(val)
		case string:
			int64Arr[i], err = strconv.ParseInt(val, 10, 64)
		default:
			int64Arr = nil
		}
	}

	return int64Arr, err
}

func stringToInt64Array(str string) ([]int64, error) {
	strArr := strings.Split(str, ",")
	int64Arr := make([]int64, len(strArr))

	var err error
	for i, v := range strArr {
		int64Arr[i], err = strconv.ParseInt(v, 10, 64)
		if err != nil {
			break
		}
	}

	return int64Arr, err
}

func JsonUnmarshal(str string, data interface{}) error {
	t := reflect.TypeOf(data)
	if t.Kind() != reflect.Ptr {
		return errors.New("parameter structure must be a pointer")
	}

	dec := json.NewDecoder(strings.NewReader(str))
	dec.DisallowUnknownFields()
	if err := dec.Decode(data); err != nil {
		msg := fmt.Sprintf("type:[%s] err:[%+v] json:%+v ", strings.ToLower(t.String()), err, str)
		return errors.New(msg)
	}

	return nil
}
