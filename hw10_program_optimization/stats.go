package hw10programoptimization

import (
	"bufio"
	"errors"
	"io"
	"strings"

	"github.com/mailru/easyjson" //nolint:depguard
)

//easyjson:json
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

func GetDomainStat(r io.Reader, domain string) (DomainStat, error) {
	if r == nil {
		return DomainStat{}, nil
	}

	return countDomains(r, domain)
}

func countDomains(r io.Reader, domain string) (DomainStat, error) {
	result := make(DomainStat)
	reader := bufio.NewReader(r)
	var email string

	for {
		line, _, err := reader.ReadLine()
		if err != nil {
			if !errors.Is(err, io.EOF) {
				return result, err
			}

			break
		}

		var user User
		if err = easyjson.Unmarshal(line, &user); err != nil {
			return result, err
		}

		email = strings.ToLower(user.Email)
		if strings.HasSuffix(email, strings.ToLower(domain)) {
			result[strings.SplitN(email, "@", 2)[1]]++
		}
	}

	return result, nil
}
