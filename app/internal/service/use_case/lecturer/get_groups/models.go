package getgroups

type GetGroupsRequest struct {
	UserID int
	TestID int
	Year   int
}

type GroupInfo struct {
	GroupName    string
	MembersCount int
}

type GetGroupsResponse struct {
	Groups []GroupInfo
}
