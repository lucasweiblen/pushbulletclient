package client

import "net/http"

type Device struct {
	Iden         string `json:"iden"`
	PushToken    string `json:"push_token"`
	AppVersion   int    `json:"app_version"`
	FingerPrint  string `json:"fingerprint"`
	Active       bool   `json:"active"`
	Nickname     string `json:"nickname"`
	Manufacturer string `json:"manufacturer"`
	Type         string `json:"type"`
	Model        string `json:"model"`
	Pushable     bool   `json:"pushable"`
}

type Devices struct {
	Devices []Device `json:"devices"`
}

type Note struct {
	title string
	body  string
	kind  string
}

type Link struct {
	title string
	body  string
	url   string
	kind  string
}

type Address struct {
	name    string
	address string
	kind    string
}

type List struct {
	title string
	items []string
	kind  string
}

type File struct {
	fileName string
	fileType string
	fileUrl  string
	body     string
	kind     string
}

type Client struct {
	token      string
	HttpClient *http.Client
}

type Params map[string]string

type Endpoint map[string]string

func NewNote(title, body string) *Note {
	return &Note{title: title, body: body, kind: "note"}
}

func NewLink(title, body, url string) *Link {
	return &Link{title: title, body: body, url: url, kind: "link"}
}

func NewAddress(name, address string) *Address {
	return &Address{name: name, address: address, kind: "address"}
}

func NewList(title string, items ...string) *List {
	return &List{title: title, items: items, kind: "list"}
}

func NewFile(fname, ftype, furl, body string) *File {
	return &File{fileName: fname, fileType: ftype, fileUrl: furl, body: body, kind: "file"}
}

func NewClient(token string) *Client {
	httpClient := &http.Client{}
	return &Client{token: token, HttpClient: httpClient}
}