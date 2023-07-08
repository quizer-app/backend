package models

type User struct {
	Id        string `json:"id" bson:"_id,omitempty"`
	CreatedAt int64  `json:"createdAt" bson:"createdAt"`
	Username  string `json:"username" bson:"username"`
	Email     string `json:"email" bson:"email"`
	Password  string `json:"password" bson:"password"`
	Verified  bool   `json:"verified" bson:"verified"`
	Role      string `json:"role" bson:"role"`
}
