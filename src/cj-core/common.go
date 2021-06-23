package cj_core

import "net/http"

// Client http request client.
var Client http.Client

// LuckyList Final lucky boys.
var LuckyList []reply

type NewestComment struct {
	Data commentData `json:"data"`
}

type commentData struct {
	Cursor  cursor  `json:"cursor"`
	Replies []reply `json:"replies"`
}

type cursor struct {
	AllCount int `json:"all_count"`
	Prev     int `json:"prev"`
	Next     int `json:"next"`
}

type reply struct {
	Floor   int     `json:"floor"`
	Member  member  `json:"member"`
	Content content `json:"content"`
}
type member struct {
	Mid    string `json:"mid"`
	Uname  string `json:"uname"`
	Avatar string `json:"avatar"`
}

type content struct {
	Message string `json:"message"`
}
