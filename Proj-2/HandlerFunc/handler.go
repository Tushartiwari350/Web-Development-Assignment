package HandlerFunc

import (
	"encoding/json"
	"net/http"

	"github.com/Tushar/myapi/database"
	"github.com/gorilla/mux"
)

func AddPost(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var newuser User
	_ = json.NewDecoder(r.Body).Decode(&newuser)
	database.AppendUserlist(newuser.Name, newuser.Msgs)
	json.NewEncoder(w).Encode(newuser)
}

func GetMyPosts(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	params := mux.Vars(r)
	// var x Resultstring
	flag1 := database.CheckValidUser(params["name"])
	if !flag1 {
		http.Error(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)
		return
	}

	// ans:= []string{}
	ans := database.ReturnMyPosts(params["name"])
	json.NewEncoder(w).Encode(ans)
}

func Follow(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var val FollowInfo
	_ = json.NewDecoder(r.Body).Decode(&val)

	flag1 := database.CheckValidUser(val.Username)
	flag2 := database.CheckValidUser(val.Following)
	if !flag1 || !flag2 {

		http.Error(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)
		return
	}
	ok := database.IterateFollowlist(val.Username, val.Following)
	if ok {
		var x Resultstring
		x.Result = "Already Following"
		json.NewEncoder(w).Encode(x)
		return
	}

	database.AppendFollowlist(val.Username, val.Following)
	json.NewEncoder(w).Encode(val)
}

func ReadPost(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	params := mux.Vars(r)
	ok := database.CheckValidUser(params["name"])
	if !ok {
		http.Error(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)
		return
	}
	allfriendsposts := database.ReturnMyFriendsPosts(params["name"])

	for idx, list := range allfriendsposts {
		var x PrintPost
		x.Name = allfriendsposts[idx][0]
		temp := list[1:]
		x.Post = temp
		json.NewEncoder(w).Encode(x)
	}

}
