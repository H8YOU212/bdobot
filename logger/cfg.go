package logfdb

type logapi struct{
	is_success		bool		`bson:"is_success"`
	req_body		string		`bson:"req_body"`
	res_body		string		`bson:"res_body"`

}

type logtg struct{
	conn_is_success		bool			`bson:"is_success"`
	update_content	interface{}		`bson:"update_content"`
}

type logdb struct{
	conn_is_success		bool		`bson:"is_success"`
}	