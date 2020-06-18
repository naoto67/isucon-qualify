package main

import (
	"strconv"
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
	itemIDs := []interface{}{"SADD", ITEM_IDS_KEY}
	tradingItemIDs := []interface{}{"SADD", TRADING_ITEM_IDS_KEY}

	for rows.Next() {
		if err := rows.StructScan(&item); err != nil {
			return err
		}
		if item.Status == ItemStatusTrading {
			tradingItemIDs = append(tradingItemIDs, strconv.Itoa(int(item.ID)))
		}
		itemIDs = append(itemIDs, strconv.Itoa(int(item.ID)))
	}

	err = redisCluster.Do(redisCtx, itemIDs...).Err()
	if err != nil {
		return err
	}
	err = redisCluster.Do(redisCtx, tradingItemIDs...).Err()
	return err
}

func addItemID(itemID int64) error {
	err := redisCluster.Do(redisCtx, "SADD", ITEM_IDS_KEY, itemID).Err()
	return err
}

func addTradingItemID(itemID int64) error {
	err := redisCluster.Do(redisCtx, "SADD", TRADING_ITEM_IDS_KEY, itemID).Err()
	return err
}
func removeTradingItemID(itemID int64) error {
	// key が set型でなければerrorを返す
	// memberが存在していない場合は、何も実行しない
	err := redisCluster.Do(redisCtx, "SREM", TRADING_ITEM_IDS_KEY, itemID).Err()
	return err
}

func isMemberTradingItemID(itemID interface{}) (ok bool, err error) {
	ok, err = redisCluster.Do(redisCtx, "SISMEMBER", TRADING_ITEM_IDS_KEY, itemID).Bool()
	return
}

func isMemberItemID(itemID interface{}) (ok bool, err error) {
	ok, err = redisCluster.Do(redisCtx, "SISMEMBER", ITEM_IDS_KEY, itemID).Bool()
	return
}
