package hw10programoptimization

import (
	"bufio"
	"errors"
	"io"
	"regexp"
	"strings"

	"github.com/mailru/easyjson"
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
	regex := regexp.MustCompile("(?i)\\." + domain)
	reader := bufio.NewReader(r)

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

		if regex.MatchString(user.Email) {
			result[strings.ToLower(strings.SplitN(user.Email, "@", 2)[1])]++
		}
	}

	return result, nil
}
