package models

type RefreshToken struct {
	Id           string `json:"id" bson:"_id,omitempty"`
	CreatedAt    int64  `json:"createdAt" bson:"createdAt"`
	ExpiresAt    int64  `json:"expiresAt" bson:"expiresAt"`
	UserId       string `json:"userId" bson:"userId"`
	Token        string `json:"token" bson:"token"`
	UserPassword string `json:"userPassword" bson:"userPassword"`
}
