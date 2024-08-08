package hw10programoptimization

import (
	"bufio"
	"fmt"
	"io"
	"strings"
	"sync"

	"github.com/goccy/go-json"
)

const (
	EmptyString = ""
	Atsign      = "@"
)

type User struct {
	ID       int
	Name     string
	Username string
	Email    string
	Phone    string
	Password string
	Address  string
}

type DomainStat map[string]int

var userPool = sync.Pool{
	New: func() interface{} { return new(User) },
}

func GetDomainStat(r io.Reader, domain string) (DomainStat, error) {
	domain = "." + domain
	result, err := processData(r, domain)
	if err != nil {
		return nil, fmt.Errorf("get users error: %w", err)
	}

	return result, nil
}

func processData(r io.Reader, domain string) (result DomainStat, err error) {
	result = make(DomainStat)
	scanner := bufio.NewScanner(r)
	for scanner.Scan() {
		content := scanner.Bytes()
		user := userPool.Get().(*User)
		if err = json.Unmarshal(content, user); err != nil {
			continue
		}
		if err = countDomains(*user, domain, &result); err != nil {
			continue
		}
		user.resetUser()
		userPool.Put(user)
	}
	if err = scanner.Err(); err != nil {
		return
	}
	return result, nil
}

func countDomains(user User, domain string, result *DomainStat) (err error) {
	str := strings.SplitN(user.Email, Atsign, 2)
	if len(str) < 2 {
		return fmt.Errorf("invalid email: %s", user.Email)
	}
	key := strings.ToLower(str[1])
	if strings.Contains(key, domain) {
		(*result)[key]++
	}
	return nil
}

func (user *User) resetUser() {
	user.ID = 0
	user.Name = EmptyString
	user.Username = EmptyString
	user.Email = EmptyString
	user.Phone = EmptyString
	user.Password = EmptyString
	user.Address = EmptyString
}
