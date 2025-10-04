package hw10programoptimization

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io"
	"strings"
)

type User struct {
	Email string `json:"email"`
}

type DomainStat map[string]int

func GetDomainStat(r io.Reader, domain string) (DomainStat, error) {
	u, err := getUsers(r)
	if err != nil {
		return nil, fmt.Errorf("get users error: %w", err)
	}
	return countDomains(u, domain)
}

type users []User

func getUsers(r io.Reader) (result users, err error) {
	scanner := bufio.NewScanner(r)
	result = make(users, 0, 1000)
	i := 0
	for scanner.Scan() {
		var user User
		line := scanner.Bytes()
		if err = json.Unmarshal(line, &user); err != nil {
			return
		}
		result = append(result, user)
		i++
	}

	return
}

func countDomains(u users, domain string) (DomainStat, error) {
	result := make(DomainStat)

	for _, user := range u {
		matched := strings.HasSuffix(user.Email, "."+domain)

		if matched {
			result[strings.ToLower(strings.SplitN(user.Email, "@", 2)[1])]++
		}
	}
	return result, nil
}
