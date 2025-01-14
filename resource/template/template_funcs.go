package template

import (
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"net"
	"os"
	"sort"
	"strconv"
	"strings"
	"time"

	util "github.com/haad/confd/util"
	"github.com/kelseyhightower/memkv"
)

func newFuncMap() map[string]interface{} {

	m := map[string]interface{}{
		"json":           UnmarshalJsonObject,
		"jsonArray":      UnmarshalJsonArray,
		"map":            CreateMap,
		"getenv":         Getenv,
		"datetime":       time.Now,
		"toUpper":        strings.ToUpper,
		"toLower":        strings.ToLower,
		"lookupIP":       LookupIP,
		"lookupIPV4":     LookupIPV4,
		"lookupIPV6":     LookupIPV6,
		"lookupSRV":      LookupSRV,
		"fileExists":     util.IsFileExist,
		"base64Encode":   Base64Encode,
		"base64Decode":   Base64Decode,
		"parseBool":      strconv.ParseBool,
		"reverse":        Reverse,
		"sortByLength":   SortByLength,
		"sortKVByLength": SortKVByLength,
		"seq":            Seq,
		"hostname": 	  GetHostname,
		"lookupIfaceIPV4": LookupIfaceIPV4,
		"lookupIfaceIPV6": LookupIfaceIPV6,
		"printf":         fmt.Sprintf,
		"unixTS":         func() string { return strconv.FormatInt(time.Now().Unix(), 10) },
		"dateRFC3339":    func() string { return time.Now().Format(time.RFC3339) },
	}

	return m
}

func addFuncs(out, in map[string]interface{}) {
	for name, fn := range in {
		out[name] = fn
	}
}

// Seq creates a sequence of integers. It's named and used as GNU's seq.
// Seq takes the first and the last element as arguments. So Seq(3, 5) will generate [3,4,5]
func Seq(first, last int) []int {
	var arr []int
	for i := first; i <= last; i++ {
		arr = append(arr, i)
	}
	return arr
}

type byLengthKV []memkv.KVPair

func (s byLengthKV) Len() int {
	return len(s)
}

func (s byLengthKV) Swap(i, j int) {
	s[i], s[j] = s[j], s[i]
}

func (s byLengthKV) Less(i, j int) bool {
	return len(s[i].Key) < len(s[j].Key)
}

func SortKVByLength(values []memkv.KVPair) []memkv.KVPair {
	sort.Sort(byLengthKV(values))
	return values
}

type byLength []string

func (s byLength) Len() int {
	return len(s)
}
func (s byLength) Swap(i, j int) {
	s[i], s[j] = s[j], s[i]
}
func (s byLength) Less(i, j int) bool {
	return len(s[i]) < len(s[j])
}

func SortByLength(values []string) []string {
	sort.Sort(byLength(values))
	return values
}

//Reverse returns the array in reversed order
//works with []string and []KVPair
func Reverse(values interface{}) interface{} {
	switch values.(type) {
	case []string:
		v := values.([]string)
		for left, right := 0, len(v)-1; left < right; left, right = left+1, right-1 {
			v[left], v[right] = v[right], v[left]
		}
	case []memkv.KVPair:
		v := values.([]memkv.KVPair)
		for left, right := 0, len(v)-1; left < right; left, right = left+1, right-1 {
			v[left], v[right] = v[right], v[left]
		}
	}
	return values
}

// Getenv retrieves the value of the environment variable named by the key.
// It returns the value, which will the default value if the variable is not present.
// If no default value was given - returns "".
func Getenv(key string, v ...string) string {
	defaultValue := ""
	if len(v) > 0 {
		defaultValue = v[0]
	}

	value := os.Getenv(key)
	if value == "" {
		return defaultValue
	}
	return value
}

func GetHostname() (string, error) {
	value, error := os.Hostname()
	return value, error
}

// CreateMap creates a key-value map of string -> interface{}
// The i'th is the key and the i+1 is the value
func CreateMap(values ...interface{}) (map[string]interface{}, error) {
	if len(values)%2 != 0 {
		return nil, errors.New("invalid map call")
	}
	dict := make(map[string]interface{}, len(values)/2)
	for i := 0; i < len(values); i += 2 {
		key, ok := values[i].(string)
		if !ok {
			return nil, errors.New("map keys must be strings")
		}
		dict[key] = values[i+1]
	}
	return dict, nil
}

func UnmarshalJsonObject(data string) (map[string]interface{}, error) {
	var ret map[string]interface{}
	err := json.Unmarshal([]byte(data), &ret)
	return ret, err
}

func UnmarshalJsonArray(data string) ([]interface{}, error) {
	var ret []interface{}
	err := json.Unmarshal([]byte(data), &ret)
	return ret, err
}

func LookupIP(data string) []string {
	ips, err := net.LookupIP(data)
	if err != nil {
		return nil
	}
	// "Cast" IPs into strings and sort the array
	ipStrings := make([]string, len(ips))

	for i, ip := range ips {
		ipStrings[i] = ip.String()
	}
	sort.Strings(ipStrings)
	return ipStrings
}

func LookupIPV6(data string) []string {
	var addresses []string
	for _, ip := range LookupIP(data) {
		if strings.Contains(ip, ":") {
			addresses = append(addresses, ip)
		}
	}
	return addresses
}

func LookupIPV4(data string) []string {
	var addresses []string
	for _, ip := range LookupIP(data) {
		if strings.Contains(ip, ".") {
			addresses = append(addresses, ip)
		}
	}
	return addresses
}

func LookupIfaceIPV4(data string) (addr string) {
	var (
		ief      *net.Interface
		addrs    []net.Addr
		ipv4Addr net.IP
	)
	ief, err := net.InterfaceByName(data)
	if err != nil {
		return
	}
	addrs, err = ief.Addrs()
	if err != nil { // get addresses
		return
	}
	for _, addr := range addrs { // get ipv4 address
		ipv4Addr = addr.(*net.IPNet).IP.To4()
		if ipv4Addr != nil {
			break
		}
	}
	return ipv4Addr.String()
}

func LookupIfaceIPV6(data string) (addr string) {
	var (
		ief      *net.Interface
		addrs    []net.Addr
		ipv6Addr net.IP
	)
	ief, err := net.InterfaceByName(data)
	if err != nil {
		return
	}
	addrs, err = ief.Addrs()
	if err != nil { // get addresses
		return
	}
	for _, addr := range addrs { // get ipv6 address
		ipv6Addr = addr.(*net.IPNet).IP.To16()
		if ipv6Addr != nil {
			break
		}
	}
	return ipv6Addr.String()
}

type sortSRV []*net.SRV

func (s sortSRV) Len() int {
	return len(s)
}

func (s sortSRV) Swap(i, j int) {
	s[i], s[j] = s[j], s[i]
}

func (s sortSRV) Less(i, j int) bool {
	str1 := fmt.Sprintf("%s%d%d%d", s[i].Target, s[i].Port, s[i].Priority, s[i].Weight)
	str2 := fmt.Sprintf("%s%d%d%d", s[j].Target, s[j].Port, s[j].Priority, s[j].Weight)
	return str1 < str2
}

func LookupSRV(service, proto, name string) []*net.SRV {
	_, addrs, err := net.LookupSRV(service, proto, name)
	if err != nil {
		return []*net.SRV{}
	}
	sort.Sort(sortSRV(addrs))
	return addrs
}

func Base64Encode(data string) string {
	return base64.StdEncoding.EncodeToString([]byte(data))
}

func Base64Decode(data string) (string, error) {
	s, err := base64.StdEncoding.DecodeString(data)
	return string(s), err
}
