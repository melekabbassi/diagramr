package simple

import "time"

type Address struct {
	City string
}

type User struct {
	ID        string
	Name      string
	Address   *Address
	Tags      []string
	Meta      map[string]int
	CreatedAt time.Time
	secret    string
}

func (u *User) FullName(prefix string) string {
	return prefix + u.Name
}

func (u User) id() string {
	return u.ID
}
