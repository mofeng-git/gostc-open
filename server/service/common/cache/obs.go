package cache

import (
	"fmt"
	"server/global"
	"sort"
	"sync"
	"time"
)

const (
	gost_obs_tunnel_key = "gost_plugins:obs:tunnel:%s:%s"
	gost_obs_client_key = "gost_plugins:obs:client:%s:%s"
	gost_obs_node_key   = "gost_plugins:obs:node:%s:%s"
	gost_obs_user_key   = "gost_plugins:obs:user:%s:%s"
)

var (
	tunnelObsLock       = &sync.RWMutex{}
	MONTH_DATEONLY_LIST []string
)

func generateMonthDateList() (dateList []string) {
	i := 0
	for {
		dateList = append(dateList, time.Now().Add(time.Duration(-i)*time.Hour*24).Format(time.DateOnly))
		i++
		if i >= 30 {
			break
		}
	}
	sort.Strings(dateList)
	return dateList
}

func init() {
	MONTH_DATEONLY_LIST = generateMonthDateList()
	go func() {
		timer := time.NewTimer(time.Hour)
		for {
			now := time.Now()
			nextDate := time.Date(now.Year(), now.Month(), now.Day()+1, 0, 0, 0, 0, time.Local)
			timer.Reset(nextDate.Sub(now))
			<-timer.C
			MONTH_DATEONLY_LIST = generateMonthDateList()
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

func GetClientObs(dateOnly, code string) TunnelObs {
	tunnelObsLock.RLock()
	defer tunnelObsLock.RUnlock()
	var data TunnelObs
	_ = global.Cache.GetStruct(fmt.Sprintf(gost_obs_client_key, dateOnly, code), &data)
	return data
}

func GetNodeObs(dateOnly, code string) TunnelObs {
	tunnelObsLock.RLock()
	defer tunnelObsLock.RUnlock()
	var data TunnelObs
	_ = global.Cache.GetStruct(fmt.Sprintf(gost_obs_node_key, dateOnly, code), &data)
	return data
}

func GetUserObs(dateOnly, code string) TunnelObs {
	tunnelObsLock.RLock()
	defer tunnelObsLock.RUnlock()
	var data TunnelObs
	_ = global.Cache.GetStruct(fmt.Sprintf(gost_obs_user_key, dateOnly, code), &data)
	return data
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
	var data TunnelObs
	for _, dateOnly := range dateOnlyList {
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
