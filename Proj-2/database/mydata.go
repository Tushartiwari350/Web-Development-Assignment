package database

var UserList = map[string][]string{}
var FollowList = map[string][]string{}

func AppendUserlist(Name, Msgs string) {
	UserList[Name] = append(UserList[Name], Msgs)
}
func AppendFollowlist(Name string, Following string) {
	FollowList[Name] = append(FollowList[Name], Following)
}

func CheckValidUser(Name string) bool {
	_, ok := UserList[Name]
	return ok
}

func IterateFollowlist(Username, Following string) bool {
	for _, item := range FollowList[Username] {
		if Following == item {
			return true
		}
	}
	return false
}

func ReturnMyPosts(Name string) []string {
	return UserList[Name]
}

func ReturnMyFriendsPosts(Name string) [][]string {
	posts := [][]string{}
	for _, list := range FollowList[Name] {
		tempslice := []string{}
		tempslice = append(tempslice, list)
		for _, msgs := range UserList[list] {
			tempslice = append(tempslice, msgs)
		}
		posts = append(posts, tempslice)
	}
	return posts
}
