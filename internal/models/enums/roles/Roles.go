package roles

type Role int

const (
	Customer Role = iota
	Admin
)

//
// func (v *Role) Scan(value any) error {
//
// }

// func (r Role) String() string {
// 	switch r {
// 	case Admin:
// 		return "Admin"
// 	default:
// 		return "Customer"
// 	}
// }
