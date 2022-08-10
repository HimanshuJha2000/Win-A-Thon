package utils

type UpdateProfileRequest struct {
	Email        string `json:"email"`
	LinkedIn     string `json:"linked_in"`
	GitHub       string `json:"github"`
	WebLink      string `json:"web_link"`
	Organisation string `json:"organisation"`
}
