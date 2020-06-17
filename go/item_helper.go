package main

import (
	"fmt"
	"strconv"

	"github.com/gomodule/redigo/redis"
)

const (
	ITEM_IDS_KEY         string = "item_id_set"
	TRADING_ITEM_IDS_KEY string = "trading_item_id_set"
)

// item全体のidとtradingなstatusのitemのidセットを持つキャッシュ
func initializeItemIDs() error {
	rows, err := dbx.Queryx("SELECT * FROM items")
	if err != nil {
		return err
	}
	defer rows.Close()
	var item Item
	itemIDs := []interface{}{ITEM_IDS_KEY}
	tradingItemIDs := []interface{}{TRADING_ITEM_IDS_KEY}

	for rows.Next() {
		if err := rows.StructScan(&item); err != nil {
			return err
		}
		fmt.Println("DEBUG: ", item.Status)
		if item.Status == ItemStatusTrading {
			tradingItemIDs = append(tradingItemIDs, strconv.Itoa(int(item.ID)))
		}
		itemIDs = append(itemIDs, strconv.Itoa(int(item.ID)))
	}

	conn := redisPool.Get()
	defer conn.Close()
	_, err = conn.Do("SADD", itemIDs...)
	if err != nil {
		return err
	}
	_, err = conn.Do("SADD", tradingItemIDs...)
	return err
}

func addItemID(itemID int64) error {
	conn := redisPool.Get()
	_, err := conn.Do("SADD", ITEM_IDS_KEY, itemID)
	return err
}

func addTradingItemID(itemID int64) error {
	conn := redisPool.Get()
	_, err := conn.Do("SADD", TRADING_ITEM_IDS_KEY, itemID)
	return err
}
func removeTradingItemID(itemID int64) error {
	conn := redisPool.Get()
	// key が set型でなければerrorを返す
	// memberが存在していない場合は、何も実行しない
	_, err := conn.Do("SREM", TRADING_ITEM_IDS_KEY, itemID)
	return err
}

func isMemberTradingItemID(itemID interface{}) (ok bool, err error) {
	conn := redisPool.Get()
	ok, err = redis.Bool(conn.Do("SISMEMBER", TRADING_ITEM_IDS_KEY, itemID))
	return
}

func isMemberItemID(itemID interface{}) (ok bool, err error) {
	conn := redisPool.Get()
	ok, err = redis.Bool(conn.Do("SISMEMBER", ITEM_IDS_KEY, itemID))
	return
}
