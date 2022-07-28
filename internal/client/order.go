package client

import (
	"encoding/json"
	"schnegge/internal/base"
	"schnegge/internal/config"
	"time"
)

func ReadOrders(cfg config.Config, date time.Time) []base.Order {
	if !checkConfig(cfg) {
		return nil
	}
	resp, err := doGetRequest(cfg, "/rest/employee-orders/list?refDate="+FormatDateForSalat(date))

	if err != nil {
		base.Log.Panic(err)
	}
	defer resp.Body.Close()

	base.Log.Println("Response status:", resp.Status)

	var order []base.Order
	err = json.NewDecoder(resp.Body).Decode(&order)
	if err != nil {
		base.Log.Panic(err)
	}
	base.Log.Printf("%+v\n", order)
	return order
}
