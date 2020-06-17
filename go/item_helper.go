package main

import "github.com/gomodule/redigo/redis"

const (
	ITEM_IDS_KEY         = "item_id_set"
	TRADING_ITEM_IDS_KEY = "trading_item_id_set"
)

// item全体のidとtradingなstatusのitemのidセットを持つキャッシュ
func initializeItemIDs() error {
	rows, err := dbx.Queryx("SELECT * FROM items")
	if err != nil {
		return err
	}
	defer rows.Close()
	var item Item
	var itemIDs []int64
	var tradingItemIDs []int64

	for rows.Next() {
		if err := rows.Scan(&item); err != nil {
			return err
		}

		if item.Status == ItemStatusTrading {
			tradingItemIDs = append(tradingItemIDs, item.ID)
		}
		itemIDs = append(itemIDs, item.ID)
	}

	conn := redisPool.Get()
	_, err = conn.Do("SADD", ITEM_IDS_KEY, itemIDs)
	if err != nil {
		return err
	}
	_, err = conn.Do("SADD", TRADING_ITEM_IDS_KEY, tradingItemIDs)
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

func isMemberTradingItemID(itemID int64) (ok bool, err error) {
	conn := redisPool.Get()
	ok, err = redis.Bool(conn.Do("SISMEMBER", TRADING_ITEM_IDS_KEY, itemID))
	return
}
