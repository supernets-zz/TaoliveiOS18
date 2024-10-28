package TaoliveOp

import "strings"

const (
	// 直播广告
	ADType_Live = iota
	// 每6秒向上滑动广告
	ADType_Scroll6s
	// 每10秒向上滑动广告
	ADType_Scroll10s
	ADType_Unknown
)

func ADType(result []interface{}) int {
	for _, v := range result {
		txt := v.([]interface{})[1].([]interface{})[0].(string)
		if txt == "更多直播" {
			return ADType_Live
		} else if strings.Contains(txt, "滑动浏览") {
			return ADType_Scroll6s
		}
	}

	return ADType_Unknown
}
