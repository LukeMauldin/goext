package ext

import (
	"crypto/rand"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"math"
	"reflect"
	"strconv"
	"strings"
	"time"
)

//Returns a pointer to the string
func SPtr(s string) *string { return &s }

//Returns a pointer to the time
func TimePtr(t time.Time) *time.Time { return &t }

//Flattens a pointer to a string to a string
func StrElem(s *string) string {
	if s == nil {
		return ""
	}
	return *s
}

//Returns a pointer to the int64
func Int64Ptr(i int64) *int64 { return &i }

//Returns a pointer to the int32
func Int32Ptr(i int32) *int32 { return &i }

//Returns a pointer to the int
func IntPtr(i int) *int { return &i }

//Constructs a substring from the string
func Substring(s string, start int, length int) string {
	if length > len(s) {
		length = len(s)
	}
	return s[start:length]
}

//Truncate the string to the specified length
func Truncate(s string, length int) string {
	if length > len(s) {
		length = len(s)
	}
	return s[0:length]
}

//Checks the type v for nil
func IsNil(v interface{}) bool {
	if v == nil {
		return true
	}
	val := reflect.ValueOf(v)
	switch val.Kind() {
	case reflect.Chan, reflect.Func, reflect.Interface, reflect.Map,
		reflect.Ptr, reflect.Slice, reflect.UnsafePointer:
		return val.IsNil()
	}
	return false
}

//Flattens a pointer to an int to an int
func IntElem(i *int) int {
	if i == nil {
		return 0
	}
	return *i
}

func Int64Elem(i *int64) int64 {
	if i == nil {
		return 0
	}
	return *i
}

//Flattens a pointer to a float to a float
func FloatElem(f *float64) float64 {
	if f == nil {
		return 0.0
	}
	return *f
}

//Flattens a pointer to a time.Time to a zero time.Time
func TimeElem(t *time.Time) time.Time {
	if t == nil {
		return time.Time{}
	} else {
		return *t
	}
}

//Truncate the time component of a datetime
func TruncateTime(t time.Time) time.Time {
	return time.Date(t.Year(), t.Month(), t.Day(), 0, 0, 0, 0, time.Local)
}

func TimeConvertUTCToLocal(in int64) time.Time {
	tUTC := time.Unix(in, 0).UTC()
	return time.Date(tUTC.Year(), tUTC.Month(), tUTC.Day(), tUTC.Hour(), tUTC.Minute(),
		tUTC.Second(), tUTC.Nanosecond(), time.Local)
}

//Generate a GUID - concept from http://www.ashishbanerjee.com/home/go/go-generate-uuid
type Guid string

func GenerateGuid() (Guid, error) {
	uuid := make([]byte, 16)
	n, err := rand.Read(uuid)
	if n != len(uuid) || err != nil {
		return "", err
	}
	// TODO: verify the two lines implement RFC 4122 correctly
	uuid[8] = 0x80 // variant bits see page 5
	uuid[4] = 0x40 // version 4 Pseudo Random, see page 7

	return Guid(hex.EncodeToString(uuid)), nil
}

//Compares two floats to the specified precision
func FloatEqual(x, y float64, precision int) bool {
	//Convert numbers to their string representation
	strX, strY := fmt.Sprintf("%."+strconv.Itoa(precision)+"f", x), fmt.Sprintf("%."+strconv.Itoa(precision)+"f", y)
	return strings.EqualFold(strX, strY)
}

//Round a float to the specified precision
func FloatRound(val float64, precision int) float64 {
	var rounder float64
	intermed := val * math.Pow(10, float64(precision))
	if _, mod := math.Modf(intermed); mod >= 0.5 {
		rounder = math.Ceil(intermed)
	} else {
		rounder = math.Floor(intermed)
	}
	return rounder / math.Pow(10, float64(precision))
}

func MustParseFloat(s string) float64 {
	if s == "" {
		return 0.0
	}
	f, err := strconv.ParseFloat(s, 64)
	if err != nil {
		panic(fmt.Errorf("Unable to string: '%s' as float -- Error: %v", s, err))
	}
	return f
}

func MustParseInt(s string) int {
	if s == "" {
		return 0
	}
	i, err := strconv.Atoi(s)
	if err != nil {
		panic(fmt.Errorf("Unable to string: '%s' as int -- Error: %v", s, err))
	}
	return i
}

func StartFunc(chErr chan error, f func() error) {
	go func() {
		defer func() {
			if rec := recover(); rec != nil {
				if err, ok := rec.(error); ok {
					chErr <- err
				} else {
					chErr <- fmt.Errorf("%v", rec)
				}
			}
		}()
		err := f()
		chErr <- err
	}()
}

func MustMarshalJSON(v interface{}) []byte {
	ret, err := json.Marshal(v)
	if err != nil {
		panic(fmt.Errorf("Error json.Marshal: %v -- V: %+v", err, v))
	}
	return ret
}

func MustUnmarshalJSON(data []byte, v interface{}) {
	err := json.Unmarshal(data, v)
	if err != nil {
		panic(fmt.Errorf("Error json.Unmarshal: %v -- data: %s", err, string(data)))
	}
}

func InterfaceToInt(v interface{}) int {
	switch t := v.(type) {
	case float32:
		return int(t)
	case float64:
		return int(t)
	case int:
		return t
	case int8:
		return int(t)
	case int16:
		return int(t)
	case int32:
		return int(t)
	case int64:
		return int(t)
	case string:
		return MustParseInt(t)
	default:
		panic(fmt.Errorf("Type not covertable to int: %T", v))
	}
}
