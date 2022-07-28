package client

import (
	"encoding/json"
	"schnegge/internal/base"
	"schnegge/internal/config"
)

func ReadOvertime(cfg config.Config) base.Overtime {
	if !checkConfig(cfg) {
		return base.Overtime{}
	}
	resp, err := doGetRequest(cfg, "/rest/overtimes/status?includeToday=true")
	if err != nil {
		base.Log.Panic(err)
	}
	defer resp.Body.Close()

	base.Log.Println("Response status:", resp.Status)

	var overtime base.Overtime
	err = json.NewDecoder(resp.Body).Decode(&overtime)
	if err != nil {
		base.Log.Panic(err)
	}
	base.Log.Printf("%+v\n", overtime)
	return overtime
}
