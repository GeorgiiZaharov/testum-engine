package ldap

type LdapUserInfo struct {
	Login string
	Name  string
	Mail  string
	Group *string
}
