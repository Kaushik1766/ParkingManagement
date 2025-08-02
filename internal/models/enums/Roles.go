package enums

type Role int

const (
	Customer Role = iota
	Admin
)

func (r Role) String() string {
	switch r {
	case Admin:
		return "Admin"
	default:
		return "Customer"
	}
}
