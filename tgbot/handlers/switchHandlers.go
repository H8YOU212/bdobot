package handlers

func SwitchHandler(action string, curIndex *int) *int {
	switch action {
	case "prev":
		*curIndex = *curIndex - 1
	case "next":
		*curIndex = *curIndex + 1
	}
	return curIndex
}
