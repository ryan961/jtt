package jtt

import (
	"fmt"
	"testing"

	"github.com/shopspring/decimal"
)

// TestGetGeoPoint_Basic 测试基本功能
func TestGetGeoPoint_Basic(t *testing.T) {
	// 测试北纬东经（正常情况）
	lat, lng := GetGeoPointForWGS84(uint32(39909257), false, uint32(116397153), false)

	expectedLat := decimal.NewFromFloat(39.909257)
	expectedLng := decimal.NewFromFloat(116.397153)

	if !lat.Equal(expectedLat) {
		t.Errorf("纬度不匹配: 期望 %s, 实际 %s", expectedLat.String(), lat.String())
	}

	if !lng.Equal(expectedLng) {
		t.Errorf("经度不匹配: 期望 %s, 实际 %s", expectedLng.String(), lng.String())
	}
}

// TestGetGeoPoint_DirectionFlags 测试方向标志
func TestGetGeoPoint_DirectionFlags(t *testing.T) {
	tests := []struct {
		name        string
		lat         uint32
		south       bool
		lng         uint32
		west        bool
		expectedLat string
		expectedLng string
	}{
		{
			name:        "北纬东经",
			lat:         39909257,
			south:       false,
			lng:         116397153,
			west:        false,
			expectedLat: "39.909257",
			expectedLng: "116.397153",
		},
		{
			name:        "南纬东经",
			lat:         39909257,
			south:       true,
			lng:         116397153,
			west:        false,
			expectedLat: "-39.909257",
			expectedLng: "116.397153",
		},
		{
			name:        "北纬西经",
			lat:         39909257,
			south:       false,
			lng:         116397153,
			west:        true,
			expectedLat: "39.909257",
			expectedLng: "-116.397153",
		},
		{
			name:        "南纬西经",
			lat:         39909257,
			south:       true,
			lng:         116397153,
			west:        true,
			expectedLat: "-39.909257",
			expectedLng: "-116.397153",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			lat, lng := GetGeoPointForWGS84(tt.lat, tt.south, tt.lng, tt.west)

			if lat.String() != tt.expectedLat {
				t.Errorf("纬度不匹配: 期望 %s, 实际 %s", tt.expectedLat, lat.String())
			}

			if lng.String() != tt.expectedLng {
				t.Errorf("经度不匹配: 期望 %s, 实际 %s", tt.expectedLng, lng.String())
			}
		})
	}
}

// TestGetGeoPoint_BoundaryValues 测试边界值
func TestGetGeoPoint_BoundaryValues(t *testing.T) {
	tests := []struct {
		name        string
		lat         uint32
		south       bool
		lng         uint32
		west        bool
		expectedLat string
		expectedLng string
	}{
		{
			name:        "零值",
			lat:         0,
			south:       false,
			lng:         0,
			west:        false,
			expectedLat: "0",
			expectedLng: "0",
		},
		{
			name:        "最大纬度值",
			lat:         90000000, // 90度
			south:       false,
			lng:         180000000, // 180度
			west:        false,
			expectedLat: "90",
			expectedLng: "180",
		},
		{
			name:        "最大纬度值（南纬）",
			lat:         90000000,
			south:       true,
			lng:         180000000,
			west:        true,
			expectedLat: "-90",
			expectedLng: "-180",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			lat, lng := GetGeoPointForWGS84(tt.lat, tt.south, tt.lng, tt.west)

			if lat.String() != tt.expectedLat {
				t.Errorf("纬度不匹配: 期望 %s, 实际 %s", tt.expectedLat, lat.String())
			}

			if lng.String() != tt.expectedLng {
				t.Errorf("经度不匹配: 期望 %s, 实际 %s", tt.expectedLng, lng.String())
			}
		})
	}
}

// TestGetGeoPoint_RealWorldCoordinates 测试真实世界坐标
func TestGetGeoPoint_RealWorldCoordinates(t *testing.T) {
	tests := []struct {
		name        string
		lat         uint32
		south       bool
		lng         uint32
		west        bool
		location    string
		expectedLat string
		expectedLng string
	}{
		{
			name:        "北京天安门",
			lat:         39909257, // 39.909257°N
			south:       false,
			lng:         116397153, // 116.397153°E
			west:        false,
			location:    "北京天安门",
			expectedLat: "39.909257",
			expectedLng: "116.397153",
		},
		{
			name:        "上海外滩",
			lat:         31234567, // 31.234567°N
			south:       false,
			lng:         121456789, // 121.456789°E
			west:        false,
			location:    "上海外滩",
			expectedLat: "31.234567",
			expectedLng: "121.456789",
		},
		{
			name:        "纽约时代广场",
			lat:         40758896, // 40.758896°N
			south:       false,
			lng:         73985130, // 73.985130°W
			west:        true,
			location:    "纽约时代广场",
			expectedLat: "40.758896",
			expectedLng: "-73.98513",
		},
		{
			name:        "悉尼歌剧院",
			lat:         33856784, // 33.856784°S
			south:       true,
			lng:         151215297, // 151.215297°E
			west:        false,
			location:    "悉尼歌剧院",
			expectedLat: "-33.856784",
			expectedLng: "151.215297",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			lat, lng := GetGeoPointForWGS84(tt.lat, tt.south, tt.lng, tt.west)

			if lat.String() != tt.expectedLat {
				t.Errorf("%s 纬度不匹配: 期望 %s, 实际 %s", tt.location, tt.expectedLat, lat.String())
			}

			if lng.String() != tt.expectedLng {
				t.Errorf("%s 经度不匹配: 期望 %s, 实际 %s", tt.location, tt.expectedLng, lng.String())
			}
		})
	}
}

// TestGetGeoPoint_TableDriven 表驱动测试
func TestGetGeoPoint_TableDriven(t *testing.T) {
	tests := []struct {
		lat         uint32
		south       bool
		lng         uint32
		west        bool
		expectedLat decimal.Decimal
		expectedLng decimal.Decimal
	}{
		{1000000, false, 2000000, false, decimal.NewFromFloat(1.0), decimal.NewFromFloat(2.0)},
		{500000, true, 1500000, true, decimal.NewFromFloat(-0.5), decimal.NewFromFloat(-1.5)},
		{123456, false, 654321, false, decimal.NewFromFloat(0.123456), decimal.NewFromFloat(0.654321)},
		{0, false, 0, false, decimal.NewFromFloat(0), decimal.NewFromFloat(0)},
		{90000000, true, 180000000, true, decimal.NewFromFloat(-90.0), decimal.NewFromFloat(-180.0)},
	}

	for i, tt := range tests {
		t.Run(fmt.Sprintf("测试用例_%d", i+1), func(t *testing.T) {
			lat, lng := GetGeoPointForWGS84(tt.lat, tt.south, tt.lng, tt.west)

			if !lat.Equal(tt.expectedLat) {
				t.Errorf("纬度不匹配: 期望 %s, 实际 %s", tt.expectedLat.String(), lat.String())
			}

			if !lng.Equal(tt.expectedLng) {
				t.Errorf("经度不匹配: 期望 %s, 实际 %s", tt.expectedLng.String(), lng.String())
			}
		})
	}
}

// TestGetGeoPointForGCJ02_Basic 测试GCJ02基本功能
func TestGetGeoPointForGCJ02_Basic(t *testing.T) {
	// 测试北京天安门坐标转换
	lat, lng := GetGeoPointForGCJ02(39909257, false, 116397153, false)

	// GCJ02坐标应该与WGS84有偏移
	wgs84Lat, wgs84Lng := GetGeoPointForWGS84(39909257, false, 116397153, false)

	// 验证坐标确实发生了转换（不应该相等）
	if lat.Equal(wgs84Lat) && lng.Equal(wgs84Lng) {
		t.Error("GCJ02坐标应该与WGS84坐标有偏移")
	}

	// 验证坐标在合理范围内（北京附近）
	latFloat, _ := lat.Float64()
	lngFloat, _ := lng.Float64()

	if latFloat < 39.8 || latFloat > 40.0 {
		t.Errorf("GCJ02纬度超出预期范围: %f", latFloat)
	}

	if lngFloat < 116.3 || lngFloat > 116.5 {
		t.Errorf("GCJ02经度超出预期范围: %f", lngFloat)
	}
}

// TestGetGeoPointForBD09_Basic 测试BD09基本功能
func TestGetGeoPointForBD09_Basic(t *testing.T) {
	// 测试北京天安门坐标转换
	lat, lng := GetGeoPointForBD09(39909257, false, 116397153, false)

	// BD09坐标应该与WGS84有偏移
	wgs84Lat, wgs84Lng := GetGeoPointForWGS84(39909257, false, 116397153, false)

	// 验证坐标确实发生了转换（不应该相等）
	if lat.Equal(wgs84Lat) && lng.Equal(wgs84Lng) {
		t.Error("BD09坐标应该与WGS84坐标有偏移")
	}

	// 验证坐标在合理范围内（北京附近）
	latFloat, _ := lat.Float64()
	lngFloat, _ := lng.Float64()

	if latFloat < 39.8 || latFloat > 40.0 {
		t.Errorf("BD09纬度超出预期范围: %f", latFloat)
	}

	if lngFloat < 116.3 || lngFloat > 116.5 {
		t.Errorf("BD09经度超出预期范围: %f", lngFloat)
	}
}

// TestCoordinateSystemComparison 测试三种坐标系的对比
func TestCoordinateSystemComparison(t *testing.T) {
	// 使用北京天安门坐标进行测试
	testLat := uint32(39909257)
	testLng := uint32(116397153)

	wgs84Lat, wgs84Lng := GetGeoPointForWGS84(testLat, false, testLng, false)
	gcj02Lat, gcj02Lng := GetGeoPointForGCJ02(testLat, false, testLng, false)
	bd09Lat, bd09Lng := GetGeoPointForBD09(testLat, false, testLng, false)

	t.Logf("WGS84: 纬度=%s, 经度=%s", wgs84Lat.String(), wgs84Lng.String())
	t.Logf("GCJ02: 纬度=%s, 经度=%s", gcj02Lat.String(), gcj02Lng.String())
	t.Logf("BD09:  纬度=%s, 经度=%s", bd09Lat.String(), bd09Lng.String())

	// 验证三种坐标系都不相同
	if wgs84Lat.Equal(gcj02Lat) || wgs84Lng.Equal(gcj02Lng) {
		t.Error("WGS84和GCJ02坐标不应该相同")
	}

	if wgs84Lat.Equal(bd09Lat) || wgs84Lng.Equal(bd09Lng) {
		t.Error("WGS84和BD09坐标不应该相同")
	}

	if gcj02Lat.Equal(bd09Lat) || gcj02Lng.Equal(bd09Lng) {
		t.Error("GCJ02和BD09坐标不应该相同")
	}
}

// TestGetGeoPointForGCJ02_DirectionFlags 测试GCJ02方向标志
func TestGetGeoPointForGCJ02_DirectionFlags(t *testing.T) {
	tests := []struct {
		name  string
		lat   uint32
		south bool
		lng   uint32
		west  bool
	}{
		{"北纬东经", 39909257, false, 116397153, false},
		{"南纬东经", 39909257, true, 116397153, false},
		{"北纬西经", 39909257, false, 116397153, true},
		{"南纬西经", 39909257, true, 116397153, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			lat, lng := GetGeoPointForGCJ02(tt.lat, tt.south, tt.lng, tt.west)

			latFloat, _ := lat.Float64()
			lngFloat, _ := lng.Float64()

			// 验证方向标志的正确性
			if tt.south && latFloat > 0 {
				t.Errorf("南纬应该为负值，实际: %f", latFloat)
			}
			if !tt.south && latFloat < 0 {
				t.Errorf("北纬应该为正值，实际: %f", latFloat)
			}
			if tt.west && lngFloat > 0 {
				t.Errorf("西经应该为负值，实际: %f", lngFloat)
			}
			if !tt.west && lngFloat < 0 {
				t.Errorf("东经应该为正值，实际: %f", lngFloat)
			}
		})
	}
}

// TestGetGeoPointForBD09_DirectionFlags 测试BD09方向标志
func TestGetGeoPointForBD09_DirectionFlags(t *testing.T) {
	tests := []struct {
		name  string
		lat   uint32
		south bool
		lng   uint32
		west  bool
	}{
		{"北纬东经", 39909257, false, 116397153, false},
		{"南纬东经", 39909257, true, 116397153, false},
		{"北纬西经", 39909257, false, 116397153, true},
		{"南纬西经", 39909257, true, 116397153, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			lat, lng := GetGeoPointForBD09(tt.lat, tt.south, tt.lng, tt.west)

			latFloat, _ := lat.Float64()
			lngFloat, _ := lng.Float64()

			// 验证方向标志的正确性
			if tt.south && latFloat > 0 {
				t.Errorf("南纬应该为负值，实际: %f", latFloat)
			}
			if !tt.south && latFloat < 0 {
				t.Errorf("北纬应该为正值，实际: %f", latFloat)
			}
			if tt.west && lngFloat > 0 {
				t.Errorf("西经应该为负值，实际: %f", lngFloat)
			}
			if !tt.west && lngFloat < 0 {
				t.Errorf("东经应该为正值，实际: %f", lngFloat)
			}
		})
	}
}

// BenchmarkGetGeoPoint 基准测试
func BenchmarkGetGeoPoint(b *testing.B) {
	for i := 0; i < b.N; i++ {
		GetGeoPointForWGS84(39909257, false, 116397153, false)
	}
}

// BenchmarkGetGeoPointForGCJ02 GCJ02基准测试
func BenchmarkGetGeoPointForGCJ02(b *testing.B) {
	for i := 0; i < b.N; i++ {
		GetGeoPointForGCJ02(39909257, false, 116397153, false)
	}
}

// BenchmarkGetGeoPointForBD09 BD09基准测试
func BenchmarkGetGeoPointForBD09(b *testing.B) {
	for i := 0; i < b.N; i++ {
		GetGeoPointForBD09(39909257, false, 116397153, false)
	}
}

// --------------------
// utils.go helper tests
// --------------------

func Test_bytesSplit(t *testing.T) {
	// empty data
	if got := bytesSplit(nil, 4); got != nil {
		t.Fatalf("expected nil, got %v", got)
	}

	// limit <= 0
	data := []byte{1, 2, 3}
	got := bytesSplit(data, 0)
	if len(got) != 1 || string(got[0]) != string(data) {
		t.Fatalf("limit<=0: expected 1 chunk equal to data, got %v", got)
	}

	// limit greater than len(data)
	got = bytesSplit([]byte{1, 2, 3}, 10)
	if len(got) != 1 || len(got[0]) != 3 {
		t.Fatalf("limit>len: expected 1 chunk of len 3, got %v", got)
	}

	// exact multiple
	got = bytesSplit([]byte{1, 2, 3, 4}, 2)
	if len(got) != 2 || got[0][0] != 1 || got[0][1] != 2 || got[1][0] != 3 || got[1][1] != 4 {
		t.Fatalf("exact multiple: unexpected chunks %v", got)
	}

	// non-multiple
	got = bytesSplit([]byte{1, 2, 3, 4, 5}, 2)
	if len(got) != 3 || len(got[2]) != 1 || got[2][0] != 5 {
		t.Fatalf("non-multiple: unexpected chunks %v", got)
	}
}

func Test_BytesToString(t *testing.T) {
	// no zero terminator
	if s := BytesToString([]byte("hello")); s != "hello" {
		t.Fatalf("expected hello, got %q", s)
	}
	// with zero terminator
	if s := BytesToString([]byte{'a', 'b', 0, 'c'}); s != "ab" {
		t.Fatalf("expected ab, got %q", s)
	}
	// starting with zero
	if s := BytesToString([]byte{0, 'x', 'y'}); s != "" {
		t.Fatalf("expected empty, got %q", s)
	}
}

func Test_StringToBCD(t *testing.T) {
	// even length
	b := StringToBCD("1234")
	if len(b) != 2 || b[0] != 0x12 || b[1] != 0x34 {
		t.Fatalf("even: expected [0x12 0x34], got %v", b)
	}

	// odd length -> left-pad 0
	b = StringToBCD("123")
	if len(b) != 2 || b[0] != 0x01 || b[1] != 0x23 {
		t.Fatalf("odd: expected [0x01 0x23], got %v", b)
	}

	// size shorter than bcd -> truncated copy from left
	b = StringToBCD("1234", 1)
	if len(b) != 1 || b[0] != 0x12 {
		t.Fatalf("size short: expected [0x12], got %v", b)
	}

	// size longer than bcd -> right aligned
	b = StringToBCD("1234", 4)
	if len(b) != 4 || b[0] != 0x00 || b[1] != 0x00 || b[2] != 0x12 || b[3] != 0x34 {
		t.Fatalf("size long: expected [0 0 0x12 0x34], got %v", b)
	}
}

func Test_BcdToString(t *testing.T) {
	// leading zero byte(s) should be stripped before conversion per current logic
	if s := BcdToString([]byte{0x00, 0x12, 0x34}); s != "1234" {
		t.Fatalf("leading zero byte strip: expected 1234, got %q", s)
	}

	// without ignorePadding: trim leading '0' chars from result
	if s := BcdToString([]byte{0x01, 0x23}); s != "123" {
		t.Fatalf("default trim zeros: expected 123, got %q", s)
	}

	// with ignorePadding=true: do not trim leading '0' chars
	if s := BcdToString([]byte{0x01, 0x23}, true); s != "0123" {
		t.Fatalf("keep zeros: expected 0123, got %q", s)
	}
}
