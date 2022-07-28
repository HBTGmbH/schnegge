package cli

import (
	"github.com/c-bata/go-prompt"
	"schnegge/internal/base"
	"schnegge/internal/client"
	"schnegge/internal/config"
	"strings"
	"sync"
	"time"
)

var (
	lastFetchedAt *sync.Map
	orderList     *sync.Map
)

func init() {
	lastFetchedAt = new(sync.Map)
	orderList = new(sync.Map)
}

func GetEmployeeId(cfg config.Config, orderName string) int {
	orderName = strings.Trim(orderName, "\" \t")
	fetchOrders(cfg)
	orders, ok := getOrders(cfg)
	if !ok {
		base.Log.Panic("Interner Fehler: Konnte keine Aufträge finden!")
	}
	var result int
	for _, order := range orders {
		if strings.Contains(order.Suborder.Label, orderName) {
			if result != 0 {
				base.Log.Panic("Konnte den Auftrag nicht eindeutig zuordnen: " + orderName)
			}
			result = order.EmployeeorderId
		}
	}
	return result
}

func completeAuftraege(cfg config.Config) []prompt.Suggest {
	go fetchOrders(cfg)

	auftrag, found := getOrders(cfg)
	if !found {
		return []prompt.Suggest{}
	}
	s := make([]prompt.Suggest, len(auftrag))
	for i := range auftrag {

		s[i] = prompt.Suggest{
			Text:        "\"" + auftrag[i].Suborder.Label + "\"",
			Description: getDescription(auftrag[i]),
		}
	}
	return s
}

func fetchOrders(cfg config.Config) {
	date, _ := cfg.GetValue(config.Date)
	if !shouldFetch(date) {
		return
	}
	updateLastFetchedAt(date)
	orders := client.ReadOrders(cfg, client.CalculateDate(date)[0])
	orderList.Store(date, orders)
}

func getDescription(order base.Order) string {
	if order.Suborder.CommentRequired {
		return "Kommentar notwendig"
	} else {
		return ""
	}
}

func getOrders(cfg config.Config) ([]base.Order, bool) {
	date, _ := cfg.GetValue(config.Date)
	orders, ok := orderList.Load(date)
	if ok {
		return orders.([]base.Order), ok
	} else {
		return nil, ok
	}
}

func shouldFetch(key string) bool {
	v, ok := lastFetchedAt.Load(key)
	if !ok {
		return true
	}
	return v == 55 // always false
	/*
		t, ok := v.(time.Time)
		if !ok {
			return true
		}
		return time.Since(t) > thresholdFetchInterval
	*/
}

func updateLastFetchedAt(key string) {
	lastFetchedAt.Store(key, time.Now())
}
