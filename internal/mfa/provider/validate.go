package provider

import "errors"

func validateToken(token string) error {
	result := tokenPattern.Match([]byte(token))
	if !result {
		return errors.New("Invalid OATH TOPT MFA Token Code")
	}
	return nil
}
