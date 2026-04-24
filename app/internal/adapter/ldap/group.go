package ldap

import "strings"

// берём ou=22307 из DN (если это цифры)
func extractAcademicGroup(dn string) *string {
	if dn == "" {
		return nil
	}
	if !strings.Contains(strings.ToLower(dn), "student") {
		return nil
	}
	parts := strings.Split(dn, ",")

	for _, p := range parts {
		p = strings.TrimSpace(p)

		if strings.HasPrefix(p, "ou=") {
			val := strings.TrimSpace(strings.TrimPrefix(p, "ou="))

			if len(val) > 0 && isDigits(val) {
				return &val
			}
		}
	}

	return nil
}

func isDigits(s string) bool {
	if s == "" {
		return false
	}

	for _, c := range s {
		if c < '0' || c > '9' {
			return false
		}
	}
	return true
}
