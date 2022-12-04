package services

import (
	"errors"
	"fmt"
	"gorm.io/gorm"
	"new-refferal/models"
	"strings"
)

func (s *ServiceFacade) LogInOrRegister(user *models.User) (*models.User, error) {
	usr, err := s.dao.GetUserByWalletAddress(user.WalletAddress)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			lastCode, err := s.dao.GetLastLink()
			if err != nil {
				return nil, fmt.Errorf("dao.GetLastLink: %s", err.Error())
			}

			usr, err = s.dao.CreateUserAndLink(user, getNewCode(lastCode.Code))
			if err != nil {
				return nil, fmt.Errorf("dao.CreateUser: %s", err.Error())
			}
			return usr, nil
		}
		return nil, fmt.Errorf("dao.GetUserByWalletAddress: %s", err.Error())
	}

	return usr, nil
}

// a = 97, A = 65, z = 122, Z = 90
func getNewCode(code string) string {
	newCode := make([]string, 0)

	if code == "" {
		return "aaa"
	}
	length := len(code)
	runes := []rune(code)
	lastRune := runes[length-1]

	if length > 1 {
		if lastRune == 'z' {
			runes[length-1] = 'A'
			return string(runes)
		}
		if lastRune == 'Z' {
			for i := length - 2; i >= 0; i-- {
				prevRune := rune(code[i])
				if prevRune == 'Z' {
					if i == 0 {
						for j := 0; j < length+1; j++ {
							newCode = append(newCode, "a")
						}
						return strings.Join(newCode, "")
					}
					continue
				} else if prevRune == 'z' {
					runes[i] = 'A'
					for j := i + 1; j <= length-1; j++ {
						runes[j] = 'a'
					}
					return string(runes)
				} else {
					runes[i] = runes[i] + 1
					for j := i + 1; j <= length-1; j++ {
						runes[j] = 'a'
					}
					return string(runes)
				}
			}
		} else {
			runes[length-1] = runes[length-1] + 1
			return string(runes)
		}
	} else {
		if lastRune == 'z' {
			return "A"
		}
		if lastRune == 'Z' {
			return "aa"
		}

		return string(lastRune + 1)
	}

	return strings.Join(newCode, "")
}
