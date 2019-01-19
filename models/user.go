package models

import (
	"../config"
	"../utils"

	"encoding/json"
	"strconv"
	"strings"
	"time"

	"github.com/dgraph-io/badger"
)

type User struct {
	Model
	Name     string `json:"name"`
	Login    string `gorm:"unique_index; not null" json:"login"`
	Password string `gorm:"not null" json:"-"`
	IsAdmin  bool   `gorm:"default:false; not null" json:"is_admin"`

	CreatorGroups []Group `gorm:"foreignkey:CreatorID" json:"-"`
	Groups        []Group `gorm:"many2many:group_user" json:"groups"`
	Tasks         []Task  `json:"tasks"`
}

type UserTokens struct {
	Token string `json:"token"`
	Ip    string `json:"ip"`
}

const (
	keyUser      = "user."
	keyUserToken = "user.token."
)

var currentUser User

// Return token
func (user User) AddToken(ip string) (string, error) {
	rnd := utils.RandStringBytesMaskImpr(60)
	token := user.GetId() + ":" + rnd

	badgerDb := config.GetBadgerDb()
	err := badgerDb.Update(func(txn *badger.Txn) error {
		return txn.Set([]byte(keyUserToken+token), []byte(ip))
	})

	if err != nil {
		return "", err
	}

	return token, nil
}

func (user User) RemoveToken(token string) error {
	badgerDb := config.GetBadgerDb()
	return badgerDb.Update(func(txn *badger.Txn) error {
		return txn.Delete([]byte(keyUserToken + token))
	})
}

func (user User) RemoveTokens() error {
	badgerDb := config.GetBadgerDb()
	return badgerDb.Update(func(txn *badger.Txn) error {
		opts := badger.DefaultIteratorOptions
		opts.PrefetchValues = false
		it := txn.NewIterator(opts)
		defer it.Close()

		prefix := []byte(keyUserToken + user.GetId() + ":")
		println("--->", prefix)
		for it.Seek(prefix); it.ValidForPrefix(prefix); it.Next() {
			key := it.Item().KeyCopy(nil)
			err := txn.Delete(key)

			if err != nil {
				return err
			}
		}

		return nil
	})
}

func (user User) GetTokens() ([]UserTokens, error) {
	var tokens [] UserTokens

	badgerDb := config.GetBadgerDb()
	err := badgerDb.View(func(txn *badger.Txn) error {
		opts := badger.DefaultIteratorOptions
		opts.PrefetchValues = false
		it := txn.NewIterator(opts)
		defer it.Close()

		prefix := []byte(keyUserToken + user.GetId() + ":")
		for it.Seek(prefix); it.ValidForPrefix(prefix); it.Next() {
			item := it.Item()
			k := item.Key()
			v, err := item.ValueCopy(nil)
			if err == nil {
				tokens = append(tokens, UserTokens{Token: string(k[:]), Ip: string(v[:])})
			}
		}

		return nil
	})

	if err != nil {
		return nil, err
	}

	return tokens, nil
}

func (user User) GetId() string {
	return strconv.FormatUint(uint64(user.ID), 10)
}

func GetUsers() []User {
	db := config.GetDb()

	var users []User
	db.Find(&users)

	return users
}

// True - is validated
// Set/Clear currentUser from parsed token
func ValidateUserToken(token string) bool {
	badgerDb := config.GetBadgerDb()
	currentUser = User{}

	if token == "" {
		return false
	}

	// TODO Protect Brute-force

	err := badgerDb.View(func(txn *badger.Txn) error {
		_, keyIsExists := txn.Get([]byte(keyUserToken + token))
		return keyIsExists
	})

	if err == nil {
		tokenSplit := strings.Split(token, ":")
		id, err := strconv.Atoi(tokenSplit[0])
		if err == nil {
			currentUser = GetUserById(id)
		}
	}

	return err == nil
}

// ip - Remote Address
func GetUserById(id int) User {
	var user User

	badgerDb := config.GetBadgerDb()
	db := config.GetDb()

	// Get user from badgerDB
	err := badgerDb.View(func(txn *badger.Txn) error {
		item, err := txn.Get([]byte(keyUser + strconv.Itoa(id)))
		if err != nil {
			return err
		}

		val, err := item.ValueCopy(nil)
		if err != nil {
			return err
		}

		return json.Unmarshal(val, &user)
	})

	// Get user from DB and save to badgerDB
	if err != nil {
		db.First(&user, id)

		// Add/Update with TTL
		if user.ID > 0 {
			byteUser, err := json.Marshal(user)
			if err == nil {
				_ = badgerDb.Update(func(txn *badger.Txn) error {
					return txn.SetWithTTL([]byte(keyUser+strconv.Itoa(id)), byteUser, time.Hour*24*7)
				})
			}
		}
	}

	return user
}

func GetUserByLogin(login string) User {
	var user User

	db := config.GetDb()
	db.Where("login = ?", login).First(&user)

	return user
}

func GetCurrentUser() User {
	return currentUser
}
