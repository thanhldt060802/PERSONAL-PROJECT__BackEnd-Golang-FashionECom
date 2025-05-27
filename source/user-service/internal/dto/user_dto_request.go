package dto

type GetUserByIdRequest struct {
	Id int64 `path:"id" required:"true" doc:"Id of user."`
}

type CreateUserRequest struct {
	Body struct {
		FullName string `json:"full_name" required:"true" minLength:"1" doc:"Full name of user."`
		Email    string `json:"email" required:"true" minLength:"1" format:"email" doc:"Email of user."`
		Username string `json:"username" required:"true" minLength:"1" doc:"Username of user."`
		Password string `json:"password" required:"true" minLength:"1" doc:"Password of user."`
		Address  string `json:"address" required:"true" minLength:"1" doc:"Address of user."`
		RoleName string `json:"role_name" required:"true" enum:"ADMIN,STAFF,CUSTOMER" doc:"Role name of user."`
	}
}

type UpdateUserByIdRequest struct {
	Id   int64 `path:"id" required:"true" doc:"Id of user."`
	Body struct {
		FullName *string `json:"fullname,omitempty" minLength:"1" doc:"Full name of user."`
		Email    *string `json:"email,omitempty" minLength:"1" format:"email" doc:"Email of user."`
		Password *string `json:"password,omitempty" minLength:"1" doc:"Password of user."`
		Address  *string `json:"address,omitempty" minLength:"1" doc:"Address of user."`
		RoleName *string `json:"role_name,omitempty" enum:"ADMIN,CUSTOMER" doc:"Role name of user."`
	}
}

type DeleteUserByIdRequest struct {
	Id int64 `path:"id" required:"true" doc:"Id of user."`
}

type LoginAccountRequest struct {
	Body struct {
		Username string `json:"username" required:"true" minLength:"1" doc:"Username of user account."`
		Password string `json:"password" required:"true" minLength:"1" doc:"Password of user account."`
	}
}

type RegisterAccountRequest struct {
	Body struct {
		FullName string `json:"full_name" required:"true" minLength:"1" doc:"Full name of user account."`
		Email    string `json:"email" required:"true" minLength:"1" format:"email" doc:"Email of user acount."`
		Username string `json:"username" required:"true" minLength:"1" doc:"Username of user account."`
		Password string `json:"password" required:"true" minLength:"1" doc:"Password of user account."`
		Address  string `json:"address" required:"true" minLength:"1" doc:"Address of user account."`
	}
}

type UpdateAccountRequest struct {
	Body struct {
		FullName *string `json:"fullname,omitempty" minLength:"1" doc:"Full name of user account."`
		Email    *string `json:"email,omitempty" minLength:"1" format:"email" doc:"Email of user account."`
		Password *string `json:"password,omitempty" minLength:"1" doc:"Password of user acccount."`
		Address  *string `json:"address,omitempty" minLength:"1" doc:"Address of user account."`
	}
}

type DeleteLoggedInAccountRequest struct {
	Id int64 `path:"id" required:"true" doc:"Id of user."`
}

// // Integrate with Elasticsearch

// type GetUsersWithElasticsearchRequest struct {
// 	Offset int    `query:"offset" default:"0" minimum:"0" example:"0" doc:"Skip item by offset."`
// 	Limit  int    `query:"limit" default:"5" minimum:"1" maximum:"10" example:"10" doc:"Limit item from offset."`
// 	SortBy string `query:"sort_by" default:"created_at:asc" example:"full_name:desc,created_at" doc:"Sort by one or more fields separated by commas. For example: sort_by=full_name:desc,created_at will sort by full_name in descending order, then by created_at in ascending order."`
// 	// Filter
// 	FullName     string `query:"full_name" example:"Thành Lê" doc:"Filter by full name."`
// 	Email        string `query:"email" example:"thanhle" doc:"Filter by email."`
// 	Username     string `query:"username" example:"thanhle" doc:"Filter by username."`
// 	Address      string `query:"address" example:"Quận 7, Hồ Chí Minh" doc:"Filter by address."`
// 	RoleName     string `query:"role_name" enum:"ADMIN,STAFF,CUSTOMER" example:"CUSTOMER" doc:"Filter by role name."`
// 	CreatedAtGTE string `query:"created_at_gte" example:"2024-01-15T00:00:00" doc:"Filter by created_at greater than or equal, with format is YYYY-MM-ddTHH:mm:ss."`
// 	CreatedAtLTE string `query:"created_at_lte" example:"2024-02-05T23:59:59" doc:"Filter by created_at less than or equal, with format is YYYY-MM-ddTHH:mm:ss."`
// }
// type GetUsersRequest struct {
// 	Offset int    `query:"offset" default:"0" minimum:"0" example:"0" doc:"Skip item by offset."`
// 	Limit  int    `query:"limit" default:"5" minimum:"1" maximum:"10" example:"10" doc:"Limit item from offset."`
// 	SortBy string `query:"sort_by" default:"created_at:asc" example:"full_name:desc,created_at" doc:"Sort by one or more fields separated by commas. For example: sort_by=full_name:desc,created_at will sort by full_name in descending order, then by created_at in ascending order."`
// }
