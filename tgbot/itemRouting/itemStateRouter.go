package itemrouting

import "log"

var (

)

func ItemStateRouter(chatID int64, state string) {
	switch state{
	case "swords":
		
	case "armours":

	case "jewelerys":	

	default:
		log.Print("Item not found")
	}
}