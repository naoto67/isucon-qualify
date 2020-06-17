package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

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

	err := dbx.Get(&user, "SELECT * FROM `users` WHERE `id` = ?", userID)
	if err == sql.ErrNoRows {
		return user, http.StatusNotFound, "user not found"
	}
	if err != nil {
		log.Print(err)
		return user, http.StatusInternalServerError, "db error"
	}

	return user, http.StatusOK, ""
}

func getUserSimpleByID(q sqlx.Queryer, userID int64) (userSimple UserSimple, err error) {
	user := User{}
	err = sqlx.Get(q, &user, "SELECT * FROM `users` WHERE `id` = ?", userID)
	if err != nil {
		return userSimple, err
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
