package models

type Auth struct {
	ID       int    `gorm:"primary_key" json:"id"`
	Username string `json:"username"`
	Password string `json:"password"`
}

//检查username和password是否在对应数据库的表中
func CheckAuth(username, password string) bool {
	var auth Auth
	db.Select("id").Where(Auth{Username: username, Password: password}).First(&auth)

	return auth.ID > 0
}
