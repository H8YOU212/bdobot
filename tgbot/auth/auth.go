package auth

import (
	db "bdobot/db"
)

type User struct{
	id int64		`bson:"id"`
	is_Auth bool	`bson:"is_Auth"`
	name string		`bson:"name"`
}


var filter, result any

func (u *User) CheckAuth(chatID int64) (bool, error) {
	
	db.Conn()
	defer db.Dconn()	
	filter = chatID
	res, err := db.SearchByID("users", "id", chatID)
	if err != nil{
		return false, err
	}

	if res == true {
		return true, nil
	} else{
		return false, nil
	}


}

func( u *User ) Auth(chatID int64, name string) User {
	db.Conn()
	defer db.Dconn()
	u.name = name 
	db.Insert("users", "name",)

	return User{}
}

func ( u *User ) UnAuth(chatID int64) {
	db.Conn()
	defer db.Dconn()

}

