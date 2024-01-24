package main

import (
	"encoding/json"
	"net/http"
	"strconv"
)

type User struct {
	ID       int
	Username string
}

type Server struct {
	db    map[int]*User // this is our "database"
	cache map[int]*User // this is our "cache"
	hit   int
}

func NewServer() *Server {

	db := make(map[int]*User)

	for i := 1; i <= 100; i++ {
		db[i] = &User{ID: i, Username: "user" + strconv.Itoa(i)}
	}

	return &Server{
		db:    db,
		cache: make(map[int]*User),
	}
}

func (s *Server) checkCache(id int) (*User, bool) {

	user, ok := s.cache[id]
	if !ok {
		return nil, false
	}

	return user, true
}

func (s *Server) handleGetUser(w http.ResponseWriter, r *http.Request) {

	// we get the id from the query string and then convert it to an int
	idStr := r.URL.Query().Get("id")
	// we're skipping error handling for simplicity but in
	// production code you should always handle errors
	id, _ := strconv.Atoi(idStr)

	// first we try to get the user from the cache
	user, ok := s.checkCache(id)
	if ok {
		json.NewEncoder(w).Encode(user)
		return
	}

	// then we try to get the user from the database
	user, ok = s.db[id]
	if !ok {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	// increment the hit counter
	s.hit++

	// we update the cache
	s.cache[id] = user

	json.NewEncoder(w).Encode(user)
}

func main() {}
