package handlers

import (
	"unicode"
	"regexp"
)

func userRegisterDataValidation(data []string) (bool, string){
	if len(data[0]) >= 3{
		for _, j := range data[0]{
			if !unicode.IsLetter(j){
				return false, "First name should not contain digits"
			}
		}
	}else{
		return false, "First name must be at least 3 characters length"
	}
	if len(data[1]) >= 3{
		for _, j := range data[1]{
			if !unicode.IsLetter(j){
				return false, "Last name should not contain digits"
			}
		}
	}else{
		return false, "Last name must be at least 3 characters length"
	}
	if len(data[2]) < 3{
		return false, "User name must be at least 3 characters length"
	}
	if m,_ := regexp.MatchString(`^([\w\.\_]{2,10})@(\w{1,}).([a-z]{2,4})$`, data[3]); !m{
		return false, "E-mail is incorrect"
	}
	if len(data[4]) > 5{
		u, d, c := 0, 0, 0
		for _, j := range data[4]{
			if unicode.IsDigit(j){
				d++
			}
			if unicode.IsUpper(j){
				u++
			}
			if unicode.IsLetter(j){
				c++
			}
		}
		if d == 0 || c == 0 || u == 0{
			return false, "Password should contain at least one letter in upper case and one digit"
		}
	}else{
		return false, "Password must be at least 5 characters length"
	}
	return true, ""
}
