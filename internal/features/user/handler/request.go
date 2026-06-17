package handler

type createUserRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type updateUserRequest struct {
	Email     *string `json:"email,omitempty"`
	Password  *string `json:"password,omitempty"`
	HeightCm  *int    `json:"height_cm,omitempty"`
	BirthDate *string `json:"birth_date,omitempty"`
	Sex       *string `json:"sex,omitempty"`
}

type addProfilePhotoRequest struct {
	Url string `json:"url"`
}
