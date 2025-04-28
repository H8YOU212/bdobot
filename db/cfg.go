package db

type User struct{
	id int64		`bson:"id"`
	is_Auth bool	`bson:"is_Auth"`
	name string		`bson:"name"`
}