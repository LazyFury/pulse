package utils

import "strings"

type stringUtils struct{}

var StringUtils = stringUtils{}

func (s stringUtils) IsEmpty(str string) bool {
	return len(str) == 0
}
func (m *stringUtils) HiddenEmail(str string) string {
	email := str
	if email == "" {
		return ""
	}

	arr := strings.Split(email, "@")
	if len(arr) != 2 {
		return email
	}

	prefix := strings.Split(email, "@")[0]
	suffix := strings.Split(email, "@")[1]
	if len(prefix) <= 2 {
		return prefix
	}
	return prefix[0:2] + "****" + prefix[len(prefix)-1:] + "@" + suffix
}

func (s *stringUtils) HiddenPhone(str string) string {
	phone := str
	if phone == "" {
		return ""
	}
	if len(phone) <= 4 {
		return phone
	}
	start := 3
	end := len(phone) - 4
	return phone[0:start] + "****" + phone[end:]
}
