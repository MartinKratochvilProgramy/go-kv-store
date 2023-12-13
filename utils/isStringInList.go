package utils

func IsStringInList(s *string, list *[2]string) bool {
	for _, item := range *list {
		if item == *s {
			return true
		}
	}
	return false
}
