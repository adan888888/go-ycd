package utils

import (
	"encoding/json"
	"strconv"
	"strings"
)

// ParseDecimalString 解析金额字符串，兼容带 '+' 前缀的历史数据。
func ParseDecimalString(s string) (float64, error) {
	s = strings.TrimSpace(s)
	s = strings.TrimPrefix(s, "+")
	if s == "" {
		return 0, nil
	}
	return strconv.ParseFloat(s, 64)
}

// FlexDecimal 兼容 JSON 数字或字符串（旧客户端可能传 "200.0"）。
type FlexDecimal float64

func (f *FlexDecimal) UnmarshalJSON(b []byte) error {
	if len(b) == 0 || string(b) == "null" {
		*f = 0
		return nil
	}
	var n float64
	if err := json.Unmarshal(b, &n); err == nil {
		*f = FlexDecimal(n)
		return nil
	}
	var s string
	if err := json.Unmarshal(b, &s); err != nil {
		return err
	}
	v, err := ParseDecimalString(s)
	if err != nil {
		return err
	}
	*f = FlexDecimal(v)
	return nil
}

func (f FlexDecimal) Float64() float64 {
	return float64(f)
}

// NullableFlexDecimal 可空金额，兼容 JSON null / 数字 / 字符串。
type NullableFlexDecimal struct {
	Value *float64
}

func (n *NullableFlexDecimal) UnmarshalJSON(b []byte) error {
	if len(b) == 0 || string(b) == "null" {
		n.Value = nil
		return nil
	}
	var f FlexDecimal
	if err := json.Unmarshal(b, &f); err != nil {
		return err
	}
	v := f.Float64()
	n.Value = &v
	return nil
}

// FormatDecimal 格式化金额为 JSON/展示字符串，去掉多余尾零。
func FormatDecimal(v float64) string {
	return strconv.FormatFloat(v, 'f', -1, 64)
}

// DecimalPtrValue 读取可空 DECIMAL 指针，nil 视为 0。
func DecimalPtrValue(p *float64) float64 {
	if p == nil {
		return 0
	}
	return *p
}

// HasDecimalPtr 判断可空 DECIMAL 是否有值。
func HasDecimalPtr(p *float64) bool {
	return p != nil
}
