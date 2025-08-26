package helper

func CensorID(id string) string {
	CensoredID := []rune(id)
	n := len(id)
	mid := n / 2
	if n%2 == 0 { // Even length
		CensoredID[mid-1] = '*'
		CensoredID[mid] = '*'
	} else { // Odd length
		CensoredID[mid] = '*'
		CensoredID[mid+1] = '*'
	}

	return string(CensoredID)
}
