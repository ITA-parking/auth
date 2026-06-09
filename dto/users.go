package dto

type UserListItemDto struct {
	ID       string `json:"id"`
	Username string `json:"username"`
	Email    string `json:"email"`
}

type UserListResponse struct {
	Body UserListResponseBody
}

type UserListResponseBody struct {
	Users []UserListItemDto `json:"users"`
}

type AddUserToChatRequest struct {
	UserID string `path:"userId"`
}

type AddUserToChatResponse struct {
}
