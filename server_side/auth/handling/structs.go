package handling

type RegistationRequest struct {
	Nickname string `json:"nickname"`
	Password string `json:"password"`
}
