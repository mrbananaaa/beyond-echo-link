package auth

type signupRequest struct {
	Email          string `json:"email" validate:"required,email"`
	Username       string `json:"username" validate:"required,min=6,max=21"`
	Password       string `json:"password" validate:"required,max=25"`
	Bio            string `json:"bio"`
	ProfilePicture string `json:"profile_picture"`
}
