package itemspecrouting

import "bdobot/db"



func GetUser(chatID int64) *db.User {
	user, err := db.FindUserByID(chatID)
	if err != nil{
		return nil
	}
	return user
}