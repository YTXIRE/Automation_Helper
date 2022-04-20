package auth

type DTO struct {
	ID       string `json:"id,omitempty" bson:"_id,omitempty"`
	Username string `json:"username" bson:"username"`
	Password string `json:"password" bson:"password"`
}

type Tokens struct {
	ID           string `json:"id,omitempty" bson:"_id,omitempty"`
	UserId       string `json:"user_id" bson:"user_id"`
	AccessToken  string `json:"access_token" bson:"access_token"`
	RefreshToken string `json:"refresh_token" bson:"refresh_token"`
}

type Respond struct {
	ID           string `json:"id,omitempty" bson:"_id,omitempty"`
	Username     string `json:"username" bson:"username"`
	AccessToken  string `json:"access_token" bson:"access_token"`
	RefreshToken string `json:"refresh_token" bson:"refresh_token"`
	Email        string `json:"email" bson:"email"`
	CreatedAt    int64  `json:"created_at" bson:"created_at"`
	UpdatedAt    int64  `json:"updated_at" bson:"updated_at"`
	LastLogin    int64  `json:"last_login" bson:"last_login"`
	Password     string `json:"-" bson:"password"`
}

type Token struct {
	AccessToken  string `json:"access_token" bson:"access_token"`
	RefreshToken string `json:"refresh_token" bson:"refresh_token"`
}
