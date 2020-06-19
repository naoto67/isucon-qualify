package main

import (
	"encoding/json"
	"fmt"

	"google.golang.org/appengine/memcache"
)

const (
	ITEM_KEY         = "item_id:"
	TRADING_ITEM_KEY = "trading_item_id:"
)

func initializeItems() error {
	rows, err := dbx.Queryx("SELECT * FROM items")
	if err != nil {
		return err
	}
	defer rows.Close()
	var item Item
	data, _ := json.Marshal(1)

	for rows.Next() {
		if err := rows.StructScan(&item); err != nil {
			return err
		}
		if item.Status == ItemStatusTrading {
			key := fmt.Sprintf("%s%d", TRADING_ITEM_KEY, item.ID)
			err = cacheClient.Set(key, data)
			if err != nil {
				return err
			}
		}
		key := fmt.Sprintf("%s%d", ITEM_KEY, item.ID)
		err = cacheClient.Set(key, data)
		if err != nil {
			return err
		}
	}
	return nil
}

func addItemID(itemID int64) error {
	data, _ := json.Marshal(1)
	key := fmt.Sprintf("%s%d", ITEM_KEY, itemID)
	return cacheClient.Set(key, data)
}

func addTradingItemID(itemID int64) error {
	data, _ := json.Marshal(1)
	key := fmt.Sprintf("%s%d", TRADING_ITEM_KEY, itemID)
	return cacheClient.Set(key, data)
}

func removeTradingItemID(itemID int64) error {
	key := fmt.Sprintf("%s%d", TRADING_ITEM_KEY, itemID)
	return cacheClient.Delete(key)
}

func isMemberOfItemIDs(itemID int64) (bool, error) {
	key := fmt.Sprintf("%s%d", ITEM_KEY, itemID)
	_, err := cacheClient.Get(key)
	if err != nil && err == memcache.ErrCacheMiss {
		return false, nil
	}
	if err != nil {
		return false, err
	}
	return true, nil
}

func isMemberOfTradingItemIDs(itemID int64) (bool, error) {
	key := fmt.Sprintf("%s%d", TRADING_ITEM_KEY, itemID)
	_, err := cacheClient.Get(key)
	if err != nil && err == memcache.ErrCacheMiss {
		return false, nil
	}
	if err != nil {
		return false, err
	}
	return true, nil
}
