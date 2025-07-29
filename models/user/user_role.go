package user

type Role string

const (
	RoleAuthor      Role = "AUTHOR"
	RoleContributor Role = "CONTRIBUTOR"
	RoleAdmin       Role = "ADMIN"
	RoleReader      Role = "READER"
)
