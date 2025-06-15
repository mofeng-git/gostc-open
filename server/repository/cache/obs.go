package cache

import (
	"fmt"
	"server/global"
	"sync"
	"time"
)

const (
	gost_obs_tunnel_key = "gost_plugins:obs:tunnel:%s:%s"
	gost_obs_client_key = "gost_plugins:obs:client:%s:%s"
	gost_obs_node_key   = "gost_plugins:obs:node:%s:%s"
	gost_obs_user_key   = "gost_plugins:obs:user:%s:%s"

	gost_obs_node_limit_key = "gost_plugins:obs:node:limit:%s" // 流量预警
)

var (
	tunnelObsLock       = &sync.RWMutex{}
	MONTH_DATEONLY_LIST []string
)

func init() {
	MONTH_DATEONLY_LIST = getDateRange(time.Now().AddDate(0, 0, -29), time.Now(), time.DateOnly)
	go func() {
		timer := time.NewTimer(time.Hour)
		for {
			now := time.Now()
			nextDate := time.Date(now.Year(), now.Month(), now.Day()+1, 0, 0, 0, 0, time.Local)
			timer.Reset(nextDate.Sub(now))
			<-timer.C
			MONTH_DATEONLY_LIST = getDateRange(time.Now().AddDate(0, 0, -29), time.Now(), time.DateOnly)
		}
	}()
}

type TunnelObs struct {
	InputBytes  int64
	OutputBytes int64
}

func IncreaseObs(dateOnly, code, clientCode, nodeCode, userCode string, obs TunnelObs) {
	tunnelObsLock.Lock()
	defer tunnelObsLock.Unlock()
	{
		var data TunnelObs
		key := fmt.Sprintf(gost_obs_tunnel_key, dateOnly, code)
		_ = global.Cache.GetStruct(key, &data)
		data.InputBytes += obs.InputBytes
		data.OutputBytes += obs.OutputBytes
		global.Cache.SetStruct(key, data, time.Hour*24*40)
	}
	{
		var data TunnelObs
		key := fmt.Sprintf(gost_obs_client_key, dateOnly, clientCode)
		_ = global.Cache.GetStruct(key, &data)
		data.InputBytes += obs.InputBytes
		data.OutputBytes += obs.OutputBytes
		global.Cache.SetStruct(key, data, time.Hour*24*40)
	}
	{
		var data TunnelObs
		key := fmt.Sprintf(gost_obs_node_key, dateOnly, nodeCode)
		_ = global.Cache.GetStruct(key, &data)
		data.InputBytes += obs.InputBytes
		data.OutputBytes += obs.OutputBytes
		global.Cache.SetStruct(key, data, time.Hour*24*40)
	}
	{
		var data TunnelObs
		key := fmt.Sprintf(gost_obs_user_key, dateOnly, userCode)
		_ = global.Cache.GetStruct(key, &data)
		data.InputBytes += obs.InputBytes
		data.OutputBytes += obs.OutputBytes
		global.Cache.SetStruct(key, data, time.Hour*24*40)
	}
}

func GetTunnelObs(dateOnly, code string) TunnelObs {
	tunnelObsLock.RLock()
	defer tunnelObsLock.RUnlock()
	var data TunnelObs
	_ = global.Cache.GetStruct(fmt.Sprintf(gost_obs_tunnel_key, dateOnly, code), &data)
	return data
}

func OverflowTunnelObs(dateOnly string, code string, data TunnelObs) {
	tunnelObsLock.Lock()
	defer tunnelObsLock.Unlock()
	key := fmt.Sprintf(gost_obs_tunnel_key, dateOnly, code)
	global.Cache.SetStruct(key, data, time.Hour*24*40)
}

func GetClientObs(dateOnly, code string) TunnelObs {
	tunnelObsLock.RLock()
	defer tunnelObsLock.RUnlock()
	var data TunnelObs
	_ = global.Cache.GetStruct(fmt.Sprintf(gost_obs_client_key, dateOnly, code), &data)
	return data
}

func OverflowClientObs(dateOnly string, code string, data TunnelObs) {
	tunnelObsLock.Lock()
	defer tunnelObsLock.Unlock()
	key := fmt.Sprintf(gost_obs_client_key, dateOnly, code)
	global.Cache.SetStruct(key, data, time.Hour*24*40)
}

func GetNodeObs(dateOnly, code string) TunnelObs {
	tunnelObsLock.RLock()
	defer tunnelObsLock.RUnlock()
	var data TunnelObs
	_ = global.Cache.GetStruct(fmt.Sprintf(gost_obs_node_key, dateOnly, code), &data)
	return data
}

func OverflowNodeObs(dateOnly string, code string, data TunnelObs) {
	tunnelObsLock.Lock()
	defer tunnelObsLock.Unlock()
	key := fmt.Sprintf(gost_obs_node_key, dateOnly, code)
	global.Cache.SetStruct(key, data, time.Hour*24*40)
}

func GetUserObs(dateOnly, code string) TunnelObs {
	tunnelObsLock.RLock()
	defer tunnelObsLock.RUnlock()
	var data TunnelObs
	_ = global.Cache.GetStruct(fmt.Sprintf(gost_obs_user_key, dateOnly, code), &data)
	return data
}

func OverflowUserObs(dateOnly string, code string, data TunnelObs) {
	tunnelObsLock.Lock()
	defer tunnelObsLock.Unlock()
	key := fmt.Sprintf(gost_obs_user_key, dateOnly, code)
	global.Cache.SetStruct(key, data, time.Hour*24*40)
}

func GetTunnelObsDateRange(dateOnlyList []string, code string) (result TunnelObs) {
	tunnelObsLock.RLock()
	defer tunnelObsLock.RUnlock()
	var data TunnelObs
	for _, dateOnly := range dateOnlyList {
		_ = global.Cache.GetStruct(fmt.Sprintf(gost_obs_tunnel_key, dateOnly, code), &data)
		result.InputBytes += data.InputBytes
		result.OutputBytes += data.OutputBytes
	}
	return result
}

func GetClientObsDateRange(dateOnlyList []string, code string) (result TunnelObs) {
	tunnelObsLock.RLock()
	defer tunnelObsLock.RUnlock()
	var data TunnelObs
	for _, dateOnly := range dateOnlyList {
		_ = global.Cache.GetStruct(fmt.Sprintf(gost_obs_client_key, dateOnly, code), &data)
		result.InputBytes += data.InputBytes
		result.OutputBytes += data.OutputBytes
	}
	return result
}

func GetNodeObsDateRange(dateOnlyList []string, code string) (result TunnelObs) {
	tunnelObsLock.RLock()
	defer tunnelObsLock.RUnlock()
	for _, dateOnly := range dateOnlyList {
		var data TunnelObs
		_ = global.Cache.GetStruct(fmt.Sprintf(gost_obs_node_key, dateOnly, code), &data)
		result.InputBytes += data.InputBytes
		result.OutputBytes += data.OutputBytes
	}
	return result
}

func GetUserObsDateRange(dateOnlyList []string, code string) (result TunnelObs) {
	tunnelObsLock.RLock()
	defer tunnelObsLock.RUnlock()
	var data TunnelObs
	for _, dateOnly := range dateOnlyList {
		_ = global.Cache.GetStruct(fmt.Sprintf(gost_obs_user_key, dateOnly, code), &data)
		result.InputBytes += data.InputBytes
		result.OutputBytes += data.OutputBytes
	}
	return result
}

//func DelTunnelObs(dateOnly, tunnel string) {
//	global.Cache.Del(gost_obs_tunnel_key + tunnel + dateOnly)
//}

func GetNodeObsLimit(code string) (result TunnelObs) {
	_ = global.Cache.GetStruct(fmt.Sprintf(gost_obs_node_limit_key, code), &result)
	return result
}

func RefreshNodeObsLimit(code string, eventIndex int) {
	if eventIndex == 0 {
		global.Cache.SetStruct(fmt.Sprintf(gost_obs_node_limit_key, code), &TunnelObs{
			InputBytes:  -1,
			OutputBytes: -1,
		}, time.Hour*24*3)
		return
	}
	now := time.Now()

	if now.Day() == eventIndex {
		global.Cache.SetStruct(fmt.Sprintf(gost_obs_node_limit_key, code), TunnelObs{}, time.Hour*24*3)
		return
	}

	eventDate := time.Date(now.Year(), now.Month(), eventIndex, 0, 0, 0, 0, time.Local)
	var start, end = eventDate, eventDate
	if eventDate.After(now) {
		start = time.Date(eventDate.Year(), eventDate.Month()-1, eventDate.Day(), 0, 0, 0, 0, time.Local)
	} else {
		end = time.Date(eventDate.Year(), eventDate.Month()+1, eventDate.Day(), 0, 0, 0, 0, time.Local)
	}
	// 不计算今日的流量
	end = end.AddDate(0, -1, 0)
	obs := GetNodeObsDateRange(getDateRange(start, end, time.DateOnly), code)
	global.Cache.SetStruct(fmt.Sprintf(gost_obs_node_limit_key, code), obs, time.Hour*24*3)
}

func getDateRange(start, end time.Time, dateFormat string) []string {
	// 为避免修改原始时间对象，进行拷贝
	start = start.Local()
	end = end.Local()
	if start.After(end) {
		return nil
	}
	var dates []string
	// 循环生成日期
	for d := start; !d.After(end); d = d.AddDate(0, 0, 1) {
		dates = append(dates, d.Format(dateFormat))
	}
	return dates
}
