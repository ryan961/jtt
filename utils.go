package jtt

import (
	"bytes"
	"crypto/rand"
	"crypto/rsa"
	"fmt"
	"hash"
	"math"
	"time"
	_ "time/tzdata" // load time zone data

	"github.com/shopspring/decimal"
	"golang.org/x/text/encoding/simplifiedchinese"
	"golang.org/x/text/transform"
)

const (
	IdentityBitChar       = "7e"   // 标识位
	EscapeBitChar         = "7d"   // 转义位
	IdentityBitEscapeChar = "7d02" // 7E 转义符
	EscapeBitEscapeChar   = "7d01" // 7D 转义符
)

const (
	boundaryMark = 0x7e
	escapeMark   = 0x7d
	escapeOne    = 0x01
	escapeTwo    = 0x02
)

// 坐标系转换相关常量
const (
	X_PI   = math.Pi * 3000.0 / 180.0 // 坐标转换中使用的π值
	OFFSET = 0.00669342162296594323   // WGS84椭球偏心率平方
	AXIS   = 6378245.0                // 长半轴
)

// Unescape 返回反转义后的数据包，不影响原始数据包（去头尾以及反转义）
func Unescape(src []byte) (res []byte) {
	dst := make([]byte, 0)
	i, n := 1, len(src)
	for i < n-1 {
		if i < n-2 && src[i] == 0x7d && src[i+1] == 0x02 {
			dst = append(dst, boundaryMark)
			i += 2
		} else if i < n-2 && src[i] == 0x7d && src[i+1] == 0x01 {
			dst = append(dst, escapeMark)
			i += 2
		} else {
			dst = append(dst, src[i])
			i++
		}
	}
	return dst
}

// Escape 返回转义以后的数据包，不影响原始数据包（加头尾以及转义）
func Escape(src []byte) (res []byte) {
	dst := make([]byte, 0)
	dst = append(dst, boundaryMark)
	for _, v := range src {
		switch v {
		case boundaryMark:
			dst = append(dst, escapeMark, escapeTwo)
		case escapeMark:
			dst = append(dst, escapeMark, escapeOne)
		default:
			dst = append(dst, v)
		}
	}
	dst = append(dst, boundaryMark)
	return dst
}

// Checksum 返回校验和
func Checksum(b []byte) (sum byte) {
	for _, i := range b {
		sum ^= i
	}
	return
}

// EncryptOAEP 使用 RSA-OAEP 加密.
func EncryptOAEP(hash hash.Hash, pub *rsa.PublicKey, msg []byte, label []byte) ([]byte, error) {
	buffer := bytes.NewBuffer(nil)
	chunks := bytesSplit(msg, pub.Size()-2*hash.Size()-2)
	for _, chunk := range chunks {
		ciphertext, err := rsa.EncryptOAEP(hash, rand.Reader, pub, chunk, label)
		if err != nil {
			return nil, err
		}
		buffer.Write(ciphertext)
	}
	return buffer.Bytes(), nil
}

// DecryptOAEP 使用 RSA-OAEP 解密.
func DecryptOAEP(hash hash.Hash, priv *rsa.PrivateKey, ciphertext []byte, label []byte) ([]byte, error) {
	buffer := bytes.NewBuffer(nil)
	chunks := bytesSplit(ciphertext, priv.Size())
	for _, chunk := range chunks {
		plaintext, err := rsa.DecryptOAEP(hash, rand.Reader, priv, chunk, label)
		if err != nil {
			return nil, err
		}
		buffer.Write(plaintext)
	}
	return buffer.Bytes(), nil
}

// bytes 切割.
func bytesSplit(data []byte, limit int) [][]byte {
	// 防御：limit 非法时直接返回整体，避免除零与不必要的分配
	if limit <= 0 || len(data) == 0 {
		if len(data) == 0 {
			return nil
		}
		return [][]byte{data}
	}

	// 预估分片数量，避免反复增长
	n := (len(data) + limit - 1) / limit
	chunks := make([][]byte, 0, n)
	for i := 0; i < len(data); i += limit {
		end := i + limit
		if end > len(data) {
			end = len(data)
		}
		chunks = append(chunks, data[i:end])
	}
	return chunks
}

// BytesToString bytes 转字符串.
func BytesToString(data []byte) string {
	n := bytes.IndexByte(data, 0)
	if n == -1 {
		return string(data)
	}
	return string(data[:n])
}

// StringToBCD 字符串转 BCD.
func StringToBCD(s string, size ...int) []byte {
	if (len(s) & 1) != 0 {
		s = "0" + s
	}

	data := []byte(s)
	bcdLen := len(s) / 2
	bcd := make([]byte, bcdLen)
	// 快路径：假设输入为数字字符，避免在循环内重复求 len 和除法
	for i, j := 0, 0; i < bcdLen; i, j = i+1, j+2 {
		high := data[j] - '0'
		low := data[j+1] - '0'
		bcd[i] = (high << 4) | low
	}

	if len(size) == 0 {
		return bcd
	}

	ret := make([]byte, size[0])
	if size[0] < len(bcd) {
		copy(ret, bcd)
	} else {
		copy(ret[len(ret)-len(bcd):], bcd)
	}
	return ret
}

// BcdToString BCD 转字符串， ignorePadding 是否忽略前置 0
//
//	ignorePadding 不传或传 false 表示忽略前置 0，例如终端手机号 015321115156 -> 15321115156
//	ignorePadding 传 true 表示不忽略前置 0，例如 GPS 中的 BCD 时间 081111090119 -> 081111090119，"[081111090119]定位时间": "2008-11-11 09:01:19",
func BcdToString(data []byte, ignorePadding ...bool) string {
	for {
		if len(data) == 0 {
			return ""
		}
		if data[0] != 0 {
			break
		}
		data = data[1:]
	}

	// 直接按目标长度分配
	buf := make([]byte, len(data)*2)
	for i := 0; i < len(data); i++ {
		b := data[i]
		buf[i*2] = ((b & 0xF0) >> 4) + '0'
		buf[i*2+1] = (b & 0x0F) + '0'
	}

	if len(ignorePadding) == 0 || !ignorePadding[0] {
		for idx := range buf {
			if buf[idx] != '0' {
				return string(buf[idx:])
			}
		}
	}
	return string(buf)
}

// ToBCDTime 转为 BCD 时间
func ToBCDTime(t time.Time) []byte {
	if t.IsZero() || t.Unix() == 0 { // 全 0 表示无时间条件
		return StringToBCD("000000000000", 6)
	}

	s := t.Format("20060102150405")[2:]
	return StringToBCD(s, 6)
}

// FromBCDTime bcd 时间转为 time.Time
func FromBCDTime(bcd []byte) (time.Time, error) {
	if len(bcd) != 6 {
		return time.Time{}, fmt.Errorf("invalid bcd time: %v", bcd)
	}
	s := BcdToString(bcd, true) // 不能忽略前置 0，例如："[081111090119]定位时间": "2008-11-11 09:01:19"
	if s == "" {
		return time.Time{}, nil
	}
	t, err := time.ParseInLocation("20060102150405", "20"+s, time.Local)
	if err != nil {
		return time.Time{}, fmt.Errorf("parse (%v) bcd time (%s) err: %w", bcd, s, err)
	}
	return t, nil
}

// FromBCDTimeInLocation bcd 时间转为 time.Time
func FromBCDTimeInLocation(bcd []byte, loc *time.Location) (time.Time, error) {
	if len(bcd) != 6 {
		return time.Time{}, fmt.Errorf("invalid bcd time: %v", bcd)
	}
	s := BcdToString(bcd, true) // 不能忽略前置 0，例如："[081111090119]定位时间": "2008-11-11 09:01:19"
	if s == "" {
		return time.Time{}, nil
	}
	t, err := time.ParseInLocation("20060102150405", "20"+s, loc)
	if err != nil {
		return time.Time{}, fmt.Errorf("parse (%v) bcd time (%s) err: %w", bcd, s, err)
	}
	return t, nil
}

// GetGeoPointForWGS84 获取经纬度（地球坐标系，国际通用）
//
//	lat 纬度，以度为单位的维度值乘以10的6次方，精确到百万分之一度
//	lon 经度，以度为单位的维度值乘以10的6次方，精确到百万分之一度
func GetGeoPointForWGS84(lat uint32, south bool, lng uint32, west bool) (decimal.Decimal, decimal.Decimal) {
	div := decimal.NewFromFloat(1000000)
	fLat := decimal.NewFromInt(int64(lat)).Div(div)
	fLon := decimal.NewFromInt(int64(lng)).Div(div)
	if south {
		fLat = decimal.Zero.Sub(fLat)
	}
	if west {
		fLon = decimal.Zero.Sub(fLon)
	}
	return fLat.Truncate(6), fLon.Truncate(6)
}

// GetGeoPointForGCJ02 获取经纬度（GCJ02 坐标系，高德、腾讯）
//
//	lat 纬度，以度为单位的维度值乘以 10 的 6 次方，精确到百万分之一度
//	lon 经度，以度为单位的维度值乘以 10 的 6 次方，精确到百万分之一度
func GetGeoPointForGCJ02(lat uint32, south bool, lng uint32, west bool) (decimal.Decimal, decimal.Decimal) {
	// 先获取 WGS84 坐标
	wgs84Lat, wgs84Lng := GetGeoPointForWGS84(lat, south, lng, west)

	// 转换为 float64 进行坐标系转换
	wgs84LatFloat, _ := wgs84Lat.Float64()
	wgs84LngFloat, _ := wgs84Lng.Float64()

	// 转换为 GCJ02 坐标系
	gcj02LngFloat, gcj02LatFloat := WGS84toGCJ02(wgs84LngFloat, wgs84LatFloat)

	// 转换回 decimal.Decimal 并截断到6位小数
	gcj02Lat := decimal.NewFromFloat(gcj02LatFloat).Truncate(6)
	gcj02Lng := decimal.NewFromFloat(gcj02LngFloat).Truncate(6)

	return gcj02Lat, gcj02Lng
}

// GetGeoPointForBD09 获取经纬度（BD09 坐标系，百度）
//
//	lat 纬度，以度为单位的维度值乘以 10 的 6 次方，精确到百万分之一度
//	lon 经度，以度为单位的维度值乘以 10 的 6 次方，精确到百万分之一度
func GetGeoPointForBD09(lat uint32, south bool, lng uint32, west bool) (decimal.Decimal, decimal.Decimal) {
	// 先获取 WGS84 坐标
	wgs84Lat, wgs84Lng := GetGeoPointForWGS84(lat, south, lng, west)

	// 转换为 float64 进行坐标系转换
	wgs84LatFloat, _ := wgs84Lat.Float64()
	wgs84LngFloat, _ := wgs84Lng.Float64()

	// 转换为 BD09 坐标系
	bd09LngFloat, bd09LatFloat := WGS84toBD09(wgs84LngFloat, wgs84LatFloat)

	// 转换回 decimal.Decimal 并截断到 6 位小数
	bd09Lat := decimal.NewFromFloat(bd09LatFloat).Truncate(6)
	bd09Lng := decimal.NewFromFloat(bd09LngFloat).Truncate(6)

	return bd09Lat, bd09Lng
}

// 坐标系转换辅助函数

// isOutOfChina 判断坐标是否在中国境外
func isOutOfChina(lon, lat float64) bool {
	return !(lon > 72.004 && lon < 135.05 && lat > 3.86 && lat < 53.55)
}

// delta 计算 WGS84 到 GCJ02 的偏移量
func delta(lon, lat float64) (float64, float64) {
	dlat, dlon := coordTransform(lon-105.0, lat-35.0)
	radlat := lat / 180.0 * math.Pi
	magic := math.Sin(radlat)
	magic = 1 - OFFSET*magic*magic
	sqrtmagic := math.Sqrt(magic)
	dlat = (dlat * 180.0) / ((AXIS * (1 - OFFSET)) / (magic * sqrtmagic) * math.Pi)
	dlon = (dlon * 180.0) / (AXIS / sqrtmagic * math.Cos(radlat) * math.Pi)
	mgLat := lat + dlat
	mgLon := lon + dlon
	return mgLon, mgLat
}

// coordTransform 坐标转换的核心算法
func coordTransform(lon, lat float64) (x, y float64) {
	var lonlat = lon * lat
	var absX = math.Sqrt(math.Abs(lon))
	var lonPi, latPi = lon * math.Pi, lat * math.Pi
	var d = 20.0*math.Sin(6.0*lonPi) + 20.0*math.Sin(2.0*lonPi)
	x, y = d, d
	x += 20.0*math.Sin(latPi) + 40.0*math.Sin(latPi/3.0)
	y += 20.0*math.Sin(lonPi) + 40.0*math.Sin(lonPi/3.0)
	x += 160.0*math.Sin(latPi/12.0) + 320*math.Sin(latPi/30.0)
	y += 150.0*math.Sin(lonPi/12.0) + 300.0*math.Sin(lonPi/30.0)
	x *= 2.0 / 3.0
	y *= 2.0 / 3.0
	x += -100.0 + 2.0*lon + 3.0*lat + 0.2*lat*lat + 0.1*lonlat + 0.2*absX
	y += 300.0 + lon + 2.0*lat + 0.1*lon*lon + 0.1*lonlat + 0.1*absX
	return
}

// WGS84toGCJ02 WGS84 坐标系->火星坐标系
func WGS84toGCJ02(lon, lat float64) (float64, float64) {
	if isOutOfChina(lon, lat) {
		return lon, lat
	}
	mgLon, mgLat := delta(lon, lat)
	return mgLon, mgLat
}

// GCJ02toBD09 火星坐标系->百度坐标系
func GCJ02toBD09(lon, lat float64) (float64, float64) {
	z := math.Sqrt(lon*lon+lat*lat) + 0.00002*math.Sin(lat*X_PI)
	theta := math.Atan2(lat, lon) + 0.000003*math.Cos(lon*X_PI)
	bdLon := z*math.Cos(theta) + 0.0065
	bdLat := z*math.Sin(theta) + 0.006
	return bdLon, bdLat
}

// WGS84toBD09 WGS84坐标系->百度坐标系
func WGS84toBD09(lon, lat float64) (float64, float64) {
	lon, lat = WGS84toGCJ02(lon, lat)
	return GCJ02toBD09(lon, lat)
}

// GB18030Length 计算字符串 s 的 GB18030 编码字节长度。
// 它返回编码后的字节长度和一个可能的错误。
func GB18030Length(s string) (int, error) {
	// 如果输入字符串为空，GB18030 编码长度也为 0
	if len(s) == 0 {
		return 0, nil
	}

	// 创建 GB18030 编码器
	encoder := simplifiedchinese.GB18030.NewEncoder()

	// 使用编码器将字符串转换为 GB18030 编码的字节切片
	// transform.Bytes 会执行转换并返回结果字节和任何错误
	encodedBytes, _, err := transform.Bytes(encoder, []byte(s))
	if err != nil {
		// 如果编码过程中发生错误，返回 0 和一个包装后的错误信息
		return 0, fmt.Errorf("failed to encode string to GB18030: %w", err)
	}

	// 返回编码后字节切片的长度和 nil 错误
	return len(encodedBytes), nil
}
