package types

func NewJwt(apps []string, userId uint, workspaceId uint, workspaceMemberId uint, isAdmin bool, expiresAt int64) Jwt {
	return &jwt{
		userId:            userId,
		workspaceId:       workspaceId,
		workspaceMemberId: workspaceMemberId,
		isAdmin:           isAdmin,
		expiresAt:         expiresAt,
	}
}

type Jwt interface {
	// UserId returns the user id
	UserId() uint
	// WorkspaceId returns the workspace id
	WorkspaceId() uint
	// WorkspaceMemberId returns the workspace member id
	WorkspaceMemberId() uint
	// IsAdmin returns if the user is an admin
	IsAdmin() bool
	// ExpiresAt returns the expires at
	ExpiresAt() int64
}

type jwt struct {
	userId            uint
	workspaceId       uint
	workspaceMemberId uint
	isAdmin           bool
	expiresAt         int64 // after this time, token will be expired, you cannot refresh the token
}

type Workspace struct {
	Id uint `json:"id"`
}

func (j *jwt) UserId() uint {
	return j.userId
}

func (j *jwt) WorkspaceId() uint {
	return j.workspaceId
}

func (j *jwt) WorkspaceMemberId() uint {
	return j.workspaceMemberId
}

func (j *jwt) IsAdmin() bool {
	return j.isAdmin
}

func (j *jwt) ExpiresAt() int64 {
	return j.expiresAt
}
