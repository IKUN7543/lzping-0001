package types

type RegisterReq struct {
	Username string `json:"username"`
	Password string `json:"password"`
	Nickname string `json:"nickname"`
	Mobile   string `json:"mobile"`
	Email    string `json:"email"`
}

type RegisterResp struct {
	Id       int64  `json:"id"`
	Username string `json:"username"`
	Nickname string `json:"nickname"`
}

type LoginReq struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type LoginResp struct {
	AccessToken  string `json:"accessToken"`
	RefreshToken string `json:"refreshToken"`
	ExpiresAt    int64  `json:"expiresAt"`
	Id           int64  `json:"id"`
	Username     string `json:"username"`
	Nickname     string `json:"nickname"`
}

type UserInfoReq struct {
	Id int64 `form:"id"`
}

type UserInfoResp struct {
	Id        int64  `json:"id"`
	Username  string `json:"username"`
	Nickname  string `json:"nickname"`
	Mobile    string `json:"mobile"`
	Email     string `json:"email"`
	Gender    int32  `json:"gender"`
	Avatar    string `json:"avatar"`
	Status    int32  `json:"status"`
	CreatedAt string `json:"createdAt"`
	UpdatedAt string `json:"updatedAt"`
}

type UpdateUserReq struct {
	Nickname string `json:"nickname"`
	Avatar   string `json:"avatar"`
	Gender   int32  `json:"gender"`
}

type UpdateUserResp struct {
	Id       int64  `json:"id"`
	Nickname string `json:"nickname"`
	Avatar   string `json:"avatar"`
	Gender   int32  `json:"gender"`
}
