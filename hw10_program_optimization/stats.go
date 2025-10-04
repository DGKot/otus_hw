package hw10programoptimization

import (
	"bufio"
	"encoding/json"
	"io"
	"strings"
)

type User struct {
	Email string `json:"email"`
}

type DomainStat map[string]int

func GetDomainStat(r io.Reader, domain string) (DomainStat, error) {
	result := make(DomainStat)

	scanner := bufio.NewScanner(r)
	i := 0
	var user User
	suffix := "." + domain
	for scanner.Scan() {
		line := scanner.Bytes()
		if err := json.Unmarshal(line, &user); err != nil {
			return result, err
		}
		matched := strings.HasSuffix(user.Email, suffix)
		if matched {
			result[strings.ToLower(strings.SplitN(user.Email, "@", 2)[1])]++
		}
		i++
	}
	return result, nil
}
