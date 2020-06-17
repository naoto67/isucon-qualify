package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/gomodule/redigo/redis"
	"github.com/jmoiron/sqlx"
)

var (
	defaultLastBump, _ = time.Parse("2006-01-02", "2000-01-01")
)

func getUser(r *http.Request) (user User, errCode int, errMsg string) {
	session := getSession(r)
	userID, ok := session.Values["user_id"]
	if !ok {
		return user, http.StatusNotFound, "no session"
	}

	conn := redisPool.Get()
	key := fmt.Sprintf("%s%v", USER_KEY, userID)
	data, err := redis.Bytes(conn.Do("GET", key))
	if err != nil {
		return user, http.StatusNotFound, "user not found"
	}

	err = json.Unmarshal(data, &user)
	if err != nil {
		return user, http.StatusInternalServerError, "MarshalError"
	}

	return user, http.StatusOK, ""
}

func getUserByID(userID int64) (user User, err error) {
	conn := redisPool.Get()
	key := fmt.Sprintf("%s%v", USER_KEY, userID)
	data, err := redis.Bytes(conn.Do("GET", key))
	if err != nil {
		return user, err
	}

	err = json.Unmarshal(data, &user)
	if err != nil {
		return user, err
	}

	return user, nil
}

func getUserSimpleByID(q sqlx.Queryer, userID int64) (userSimple UserSimple, err error) {
	user := User{}

	conn := redisPool.Get()
	key := fmt.Sprintf("%s%v", USER_KEY, userID)
	data, err := redis.Bytes(conn.Do("GET", key))
	if err != nil {
		return UserSimple{}, err
	}

	err = json.Unmarshal(data, &user)
	if err != nil {
		return UserSimple{}, err
	}

	userSimple.ID = user.ID
	userSimple.AccountName = user.AccountName
	userSimple.NumSellItems = user.NumSellItems
	return userSimple, err
}

func CreateUser(accountName, address string, hashedPassword []byte) (User, error) {
	var user User
	now := time.Now()
	result, err := dbx.Exec("INSERT INTO `users` (`account_name`, `hashed_password`, `address`, `last_bump`, `created_at`) VALUES (?, ?, ?, ?, ?)",
		accountName,
		hashedPassword,
		address,
		defaultLastBump,
		now,
	)
	if err != nil {
		log.Print(err)
		return user, err
	}

	userID, err := result.LastInsertId()
	if err != nil {
		log.Print(err)
		return user, err
	}

	user = User{
		ID:             userID,
		AccountName:    accountName,
		Address:        address,
		NumSellItems:   0,
		HashedPassword: hashedPassword,
		LastBump:       defaultLastBump,
		CreatedAt:      now,
	}
	key := fmt.Sprintf("%s%v", USER_KEY, userID)
	conn := redisPool.Get()
	_, err = conn.Do("SET", key, user.toJSON())
	if err != nil {
		return user, err
	}

	return user, nil
}

func (u User) toJSON() []byte {
	m := make(map[string]interface{})
	m["id"] = u.ID
	m["account_name"] = u.AccountName
	m["hashed_password"] = u.HashedPassword
	m["address"] = u.Address
	m["num_sell_items"] = u.NumSellItems
	m["last_bump"] = u.LastBump
	m["created_at"] = u.CreatedAt

	b, _ := json.Marshal(m)
	return b
}

func initializeUsersCache() error {
	var users []User
	err := dbx.Select(&users, "SELECT * FROM `users`")
	if err != nil {
		return err
	}
	conn := redisPool.Get()

	for _, user := range users {
		key := fmt.Sprintf("%s%v", USER_KEY, user.ID)
		_, err = conn.Do("SET", key, user.toJSON())
		if err != nil {
			return err
		}
	}

	return nil
}

func updateNumSellItems(q *sqlx.Tx, userID int64, num int) error {
	now := time.Now()
	_, err := q.Exec("UPDATE `users` SET `num_sell_items`=?, `last_bump`=? WHERE `id`=?",
		num,
		now,
		userID,
	)
	if err != nil {
		return err
	}

	conn := redisPool.Get()
	key := fmt.Sprintf("%s%v", USER_KEY, userID)
	data, err := redis.Bytes(conn.Do("GET", key))
	if err != nil {
		return err
	}

	var user User
	err = json.Unmarshal(data, &user)
	if err != nil {
		return err
	}

	user.NumSellItems = num
	user.LastBump = now

	_, err = conn.Do("SET", key, user.toJSON())
	return err
}
