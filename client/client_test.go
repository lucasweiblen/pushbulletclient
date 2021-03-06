package client

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"reflect"
	"strings"
	"testing"
)

type FakeRoundTripper struct {
	message  string
	status   int
	header   map[string]string
	requests []*http.Request
}

func newTestClient(rt *FakeRoundTripper) *Client {
	client := &Client{
		token:      "foobar",
		HttpClient: &http.Client{Transport: rt},
	}
	return client
}

func (rt *FakeRoundTripper) RoundTrip(r *http.Request) (*http.Response, error) {
	body := strings.NewReader(rt.message)
	rt.requests = append(rt.requests, r)
	res := &http.Response{
		StatusCode: rt.status,
		Body:       ioutil.NopCloser(body),
		Header:     make(http.Header),
	}
	for k, v := range rt.header {
		res.Header.Set(k, v)
	}
	return res, nil
}

func (rt *FakeRoundTripper) Reset() {
	rt.requests = nil
}

func TestGetMe(t *testing.T) {
	body := `
	{
	  "iden": "ubdpjxxxOK0sKG",
	  "email": "ryan@pushbullet.com",
	  "email_normalized": "ryan@pushbullet.com",
	  "created": 1357941753.8287899,
	  "modified": 1399325992.1842301,
	  "name": "Ryan Oldenburg",
	  "image_url": "https://lh4.googleusercontent.com/-YGdcF2MteeI/AAAAAAAAAAI/AAAAAAAADPU/uo9z33FoEYs/photo.jpg",
	  "preferences": {
	    "onboarding": {
	      "app": false,
	      "friends": false,
	      "extension": false
	    },
	    "social": false
	  }
	}
    `
	var expected User
	err := json.Unmarshal([]byte(body), &expected)
	if err != nil {
		t.Errorf("Error unmarshaling JSON")
	}
	fakeRT := &FakeRoundTripper{message: body, status: http.StatusOK}
	client := newTestClient(fakeRT)
	got, _ := client.GetMe()
	if !reflect.DeepEqual(got, expected) {
		t.Errorf("Error, expected %#v, got %#v", expected, got)
	}
}

// TODO
func TestUpdateMe(t *testing.T) {
}

func TestSubscribeNoChannelTag(t *testing.T) {
	fakeRT := &FakeRoundTripper{message: "", status: http.StatusOK}
	client := newTestClient(fakeRT)
	_, err := client.Subscribe(Params{})
	if err != noChannelTagError {
		t.Errorf("Error, expected %#v, got %#v", noChannelTagError.Error(), err)
	}
}

func TestSubscribeChannel(t *testing.T) {
	body := `
	{
	  "iden": "udprOsjAoRtnM0jc",
	  "created": 1412047948.579029,
	  "modified": 1412047948.5790315,
	  "active": true,
	  "channel": {
	    "iden": "ujxPklLhvyKsjAvkMyTVh6",
	    "tag": "jblow",
	    "name": "Jonathan Blow",
	    "description": "New comments on the web by Jonathan Blow.",
	    "image_url": "https://pushbullet.imgix.net/ujxPklLhvyK-6fXf4O2JQ1dBKQedhypIKwPX0lyFfwXW/jonathan-blow.png"
	  }
	}
	`
	var expected Subscription
	err := json.Unmarshal([]byte(body), &expected)
	if err != nil {
		t.Errorf("Error unmarshaling JSON")
	}
	fakeRT := &FakeRoundTripper{message: body, status: http.StatusOK}
	client := newTestClient(fakeRT)
	got, _ := client.Subscribe(Params{"channel_tag": "jblow"})
	if !reflect.DeepEqual(got, expected) {
		t.Errorf("Error, expected %#v, got %#v", expected, got)
	}
}

func TestGetChannelNoChannel(t *testing.T) {
	fakeRT := &FakeRoundTripper{message: "", status: http.StatusOK}
	client := newTestClient(fakeRT)
	_, err := client.GetChannel(Params{})
	if err != noChannelTagError {
		t.Errorf("Error, expected %#v, got %#v", noChannelTagError.Error(), err)
	}
}

func TestGetChannel(t *testing.T) {
	body := `
	{
	  "iden": "ujxPklLhvyKsjAvkMyTVh6",
	  "tag": "jblow",
	  "name": "Jonathan Blow",
	  "description": "New comments on the web by Jonathan Blow.",
	  "image_url": "https://pushbullet.imgix.net/ujxPklLhvyK-6fXf4O2JQ1dBKQedhypIKwPX0lyFfwXW/jonathan-blow.png"
	}
	`
	var expected Channel
	err := json.Unmarshal([]byte(body), &expected)
	if err != nil {
		t.Errorf("Error unmarshaling JSON")
	}
	fakeRT := &FakeRoundTripper{message: body, status: http.StatusOK}
	client := newTestClient(fakeRT)
	got, _ := client.GetChannel(Params{"tag": "jblow"})
	if !reflect.DeepEqual(got, expected) {
		t.Errorf("Error, expected %#v, got %#v", expected, got)
	}
}

func TestUnsubscribeNoIden(t *testing.T) {
	fakeRT := &FakeRoundTripper{message: "", status: http.StatusOK}
	client := newTestClient(fakeRT)
	err := client.Unsubscribe(Params{})
	if err != noIdenError {
		t.Errorf("Error, expected %#v, got %#v", noIdenError, err)
	}
}

func TestUnsubscribe(t *testing.T) {
	fakeRT := &FakeRoundTripper{message: "", status: http.StatusOK}
	client := newTestClient(fakeRT)
	err := client.Unsubscribe(Params{"iden": "0xyz"})
	if err != nil {
		t.Errorf("Expected no error, got %#v", err)
	}
}

func TestGetContacts(t *testing.T) {
	body := `
	{
	    "contacts": [
	        {
		  "iden": "ubdcjAfszs0Smi",
		  "name": "Ryan Oldenburg",
		  "created": 1399011660.4298899,
		  "modified": 1399011660.42976,
		  "email": "ryanjoldenburg@gmail.com",
		  "email_normalized": "ryanjoldenburg@gmail.com",
		  "active": true
		}
		]
	      }
	      `
	var expected Contacts
	err := json.Unmarshal([]byte(body), &expected)
	if err != nil {
		t.Errorf("Error unmarshaling JSON")
	}
	fakeRT := &FakeRoundTripper{message: body, status: http.StatusOK}
	client := newTestClient(fakeRT)
	got, _ := client.GetContacts()
	if !reflect.DeepEqual(got, expected.Contacts) {
		t.Errorf("Error, expected %#v, got %#v", expected, got)
	}
}

func TestCreateContactError(t *testing.T) {
	client := Client{}
	_, err := client.CreateContact(Params{})
	if err.Error() != "no name has been given" {
		t.Errorf("Error, expected no name has been given, got %#v", err)
	}
	_, err = client.CreateContact(Params{"name": "foo"})
	if err.Error() != "no email has been given" {
		t.Errorf("Error, expected no email has been given, got %#v", err)
	}
}

func TestCreateContact(t *testing.T) {
	body := `
	{
	  "iden": "ubdcjAfszs0Smi",
	  "name": "Ryan Oldenburg",
	  "created": 1399011660.4298899,
	  "modified": 1399011660.42976,
	  "email": "ryanjoldenburg@gmail.com",
	  "email_normalized": "ryanjoldenburg@gmail.com",
	  "active": true
	}
  `
	var expected Contact
	err := json.Unmarshal([]byte(body), &expected)
	if err != nil {
		t.Errorf("Error unmarshaling JSON")
	}
	fakeRT := &FakeRoundTripper{message: body, status: http.StatusOK}
	client := newTestClient(fakeRT)
	got, _ := client.CreateContact(Params{"name": "foo", "email": "bar"})
	if !reflect.DeepEqual(got, expected) {
		t.Errorf("Error, expected %#v, got %#v", expected, got)
	}
}

func TestUpdateContactError(t *testing.T) {
	client := Client{}
	_, err := client.UpdateContact(Params{})
	if err != noIdenError {
		t.Errorf("Error, expected noIdenError, got %#v", err)
	}
}

func TestUpdateContact(t *testing.T) {
	body := `
	{
	  "iden": "ubdcjAfszs0Smi",
	  "name": "Ryan Oldenburg",
	  "created": 1399011660.4298899,
	  "modified": 1399011660.42976,
	  "email": "ryanjoldenburg@gmail.com",
	  "email_normalized": "ryanjoldenburg@gmail.com",
	  "active": true
	}
  `
	var expected Contact
	err := json.Unmarshal([]byte(body), &expected)
	if err != nil {
		t.Errorf("Error unmarshaling JSON")
	}
	fakeRT := &FakeRoundTripper{message: body, status: http.StatusOK}
	client := newTestClient(fakeRT)
	got, _ := client.UpdateContact(Params{"iden": "0xyz", "email": "bar"})
	if !reflect.DeepEqual(got, expected) {
		t.Errorf("Error, expected %#v, got %#v", expected, got)
	}
}

func TestDeleteContact(t *testing.T) {
	fakeRT := &FakeRoundTripper{message: "", status: http.StatusOK}
	client := newTestClient(fakeRT)
	err := client.DeleteContact(Params{"iden": "0xyz"})
	if err != nil {
		t.Errorf("Error, expected nil, got %#v", err)
	}
}

func TestGetDevices(t *testing.T) {
	body :=
		`
		{
		  "devices": [
		  {
		    "iden": "u1qSJddxeKwOGuGW",
		    "push_token": "u1qSJddxeKwOGuGWu1qdxeKwOGuGWu1qSJddxeK",
		    "app_version": 74,
		    "fingerprint": "<json_string>",
		    "active": true,
		    "nickname": "Galaxy S4",
		    "manufacturer": "samsung",
		    "type": "android",
		    "created": 1394748080.0139201,
		    "modified": 1399008037.8487799,
		    "model": "SCH-I545",
		    "pushable": true
		  }
		  ]
		}
		`
	var expected Devices
	err := json.Unmarshal([]byte(body), &expected)
	if err != nil {
		t.Errorf("Error unmarshaling JSON")
	}
	fakeRT := &FakeRoundTripper{message: body, status: http.StatusOK}
	client := newTestClient(fakeRT)
	got, _ := client.GetDevices()
	if !reflect.DeepEqual(got, expected.Devices) {
		t.Errorf("Error, expected %#v, got %#v", expected, got)
	}
}

func TestCreateDeviceError(t *testing.T) {
	client := Client{}
	_, err := client.CreateDevice(Params{})
	if err.Error() != "no nickname has been given" {
		t.Errorf("Error, expected no nickname has been given, got %#v", err)
	}
	_, err = client.CreateDevice(Params{"nickname": "foo"})
	if err.Error() != "no type has been given" {
		t.Errorf("Error, expected no type has been given, got %#v", err)
	}
}

func TestCreateDevice(t *testing.T) {
	body := `
	{
	  "iden": "udm0Tdjz5A7bL4NM",
	  "nickname": "stream_device",
	  "created": 1401840789.2369599,
	  "modified": 1401840789.2369699,
	  "active": true,
	  "type": "stream",
	  "pushable": true
	}
  `
	var expected Device
	err := json.Unmarshal([]byte(body), &expected)
	if err != nil {
		t.Errorf("Error unmarshaling JSON")
	}
	fakeRT := &FakeRoundTripper{message: body, status: http.StatusOK}
	client := newTestClient(fakeRT)
	got, _ := client.CreateDevice(Params{"nickname": "foo", "type": "stream"})
	if !reflect.DeepEqual(got, expected) {
		t.Errorf("Error, expected %#v, got %#v", expected, got)
	}
}

func TestUpdateDeviceError(t *testing.T) {
	client := Client{}
	_, err := client.UpdateDevice(Params{})
	if err != noIdenError {
		t.Errorf("Expected %#v, got %#v", noIdenError, err)
	}
}

func TestUpdateDevice(t *testing.T) {
	body := `
	{
	  "iden": "udm0Tdjz5A7bL4NM",
	  "nickname": "stream_device",
	  "created": 1401840789.2369599,
	  "modified": 1401840789.2369699,
	  "active": true,
	  "type": "stream",
	  "pushable": true
	}
  `
	var expected Device
	err := json.Unmarshal([]byte(body), &expected)
	if err != nil {
		t.Errorf("Error unmarshaling JSON")
	}
	fakeRT := &FakeRoundTripper{message: body, status: http.StatusOK}
	client := newTestClient(fakeRT)
	got, _ := client.UpdateDevice(Params{"iden": "0xyz", "nickname": "bar"})
	if !reflect.DeepEqual(got, expected) {
		t.Errorf("Error, expected %#v, got %#v", expected, got)
	}
}

func TestDeleteDeviceError(t *testing.T) {
	client := Client{}
	err := client.DeleteDevice(Params{})
	if err != noIdenError {
		t.Errorf("Expected %#v, got %#v", noIdenError, err)
	}
}

func TestDeleteDevice(t *testing.T) {
	fakeRT := &FakeRoundTripper{message: "", status: http.StatusOK}
	client := newTestClient(fakeRT)
	err := client.DeleteDevice(Params{"iden": "0xyz"})
	if err != nil {
		t.Errorf("Error, expected nil, got %#v", err)
	}
}

func TestGetPushes(t *testing.T) {
	body := `
	{
	  "pushes": [
	  {
	    "iden": "ubdprOsjAhOzf0XYq",
	    "type": "link",
	    "title": "Pushbullet",
	    "body": "Documenting our API",
	    "url": "http://docs.pushbullet.com",
	    "created": 1411595135.9685705,
	    "modified": 1411595135.9686127,
	    "active": true,
	    "dismissed": false,
	    "sender_iden": "ubd",
	    "sender_email": "ryan@pushbullet.com",
	    "sender_email_normalized": "ryan@pushbullet.com",
	    "receiver_iden": "ubd",
	    "receiver_email": "ryan@pushbullet.com",
	    "receiver_email_normalized": "ryan@pushbullet.com"
	  }
	  ]
	}
  `
	var expected Pushes
	err := json.Unmarshal([]byte(body), &expected)
	if err != nil {
		t.Errorf("Error unmarshaling JSON", err)
	}
	fakeRT := &FakeRoundTripper{message: body, status: http.StatusOK}
	client := newTestClient(fakeRT)
	got, _ := client.GetPushes()
	if !reflect.DeepEqual(got, expected.Pushes) {
		t.Errorf("Error, expected %#v, got %#v", expected, got)
	}
}

func TestCreatePushError(t *testing.T) {
	client := Client{}
	_, err := client.CreatePush(Params{})
	if err != pushNoTypeError {
		t.Errorf("Expected %#v, got %#v", pushNoTypeError, err)
	}
	_, err = client.CreatePush(Params{"type": "link"})
	if err != pushNoLinkError {
		t.Errorf("Expected %#v, got %#v", pushNoLinkError, err)
	}
	_, err = client.CreatePush(Params{"type": "address"})
	if err != pushNoAddressError {
		t.Errorf("Expected %#v, got %#v", pushNoAddressError, err)
	}
	_, err = client.CreatePush(Params{"type": "list"})
	if err != pushNoItemsError {
		t.Errorf("Expected %#v, got %#v", pushNoItemsError, err)
	}
	_, err = client.CreatePush(Params{"type": "file"})
	if err != pushNoFileNameError {
		t.Errorf("Expected %#v, got %#v", pushNoFileNameError, err)
	}
	_, err = client.CreatePush(Params{"type": "file", "file_name": "foo.txt"})
	if err != pushNoFileTypeError {
		t.Errorf("Expected %#v, got %#v", pushNoFileTypeError, err)
	}
}

func TestCreatePushNote(t *testing.T) {
	body := `
	{
	  "iden": "ubdpj29aOK0sKG",
	  "type": "note",
	  "title": "Note Title",
	  "body": "Note Body",
	  "created": 1399253701.9744401,
	  "modified": 1399253701.9746201,
	  "active": true,
	  "dismissed": false,
	  "sender_iden": "ubd",
	  "sender_email": "ryan@pushbullet.com",
	  "sender_email_normalized": "ryan@pushbullet.com",
	  "receiver_iden": "ubd",
	  "receiver_email": "ryan@pushbullet.com",
	  "receiver_email_normalized": "ryan@pushbullet.com"
	}
  `
	var expected Push
	err := json.Unmarshal([]byte(body), &expected)
	if err != nil {
		t.Errorf("Error unmarshaling JSON", err)
	}
	fakeRT := &FakeRoundTripper{message: body, status: http.StatusOK}
	client := newTestClient(fakeRT)
	got, _ := client.CreatePush(Params{"type": "note"})
	if !reflect.DeepEqual(got, expected) {
		t.Errorf("Expected %#v, got %#v", expected, got)
	}
	if got.Type != "note" {
		t.Errorf("Got %#v, expected note", got.Type)
	}
}

func TestCreatePushList(t *testing.T) {
	body := `
	{
	  "iden": "ubdpjAkaGXvUl2",
	  "type": "list",
	  "title": "foo",
	  "items": ["foo"],
	  "created": 1411595195.1267679,
	  "modified": 1411595195.1268303,
	  "active": true,
	  "dismissed": false,
	  "sender_iden": "ubd",
	  "sender_email": "ryan@pushbullet.com",
	  "sender_email_normalized": "ryan@pushbullet.com",
	  "receiver_iden": "ubd",
	  "receiver_email": "ryan@pushbullet.com",
	  "receiver_email_normalized": "ryan@pushbullet.com"
	}
  `
	var expected Push
	err := json.Unmarshal([]byte(body), &expected)
	if err != nil {
		t.Errorf("Error unmarshaling JSON", err)
	}
	fakeRT := &FakeRoundTripper{message: body, status: http.StatusOK}
	client := newTestClient(fakeRT)
	got, _ := client.CreatePush(Params{"type": "list",
		"title": "foo",
		"items": []string{"foo"}})
	if !reflect.DeepEqual(got, expected) {
		t.Errorf("Expected %#v, got %#v", expected, got)
	}
	if got.Type != "list" {
		t.Errorf("Got %#v, expected note", got.Type)
	}
	if got.Title != "foo" {
		t.Errorf("Got %#v, expected foo", got.Title)
	}
	//if len(got.Items) != 1 {
	//t.Errorf("Error")
	//}
}

func TestCreatePushLink(t *testing.T) {

}

func TestCreatePushFile(t *testing.T) {

}

func TestUpdatePushError(t *testing.T) {
	client := Client{}
	_, err := client.UpdatePush(Params{})
	if err != noIdenError {
		t.Errorf("Expected %#v, got %#v", noIdenError, err)
	}
}

func TestUpdatePush(t *testing.T) {
	body := `
	{
	  "iden": "ubdpj29aOK0sKG",
	  "type": "note",
	  "title": "Note Title",
	  "body": "Note Body",
	  "created": 1399253701.9744401,
	  "modified": 1399253701.9746201,
	  "active": true,
	  "dismissed": false,
	  "sender_iden": "ubd",
	  "sender_email": "ryan@pushbullet.com",
	  "sender_email_normalized": "ryan@pushbullet.com",
	  "receiver_iden": "ubd",
	  "receiver_email": "ryan@pushbullet.com",
	  "receiver_email_normalized": "ryan@pushbullet.com"
	}
  `
	var expected Push
	err := json.Unmarshal([]byte(body), &expected)
	if err != nil {
		t.Errorf("Error unmarshaling JSON")
	}
	fakeRT := &FakeRoundTripper{message: body, status: http.StatusOK}
	client := newTestClient(fakeRT)
	got, _ := client.UpdatePush(Params{"iden": "0xyz", "title": "foobaz"})
	if !reflect.DeepEqual(got, expected) {
		t.Errorf("Error, expected %#v, got %#v", expected, got)
	}
}

func TestDeletePushError(t *testing.T) {
	client := Client{}
	err := client.DeletePush(Params{})
	if err != noIdenError {
		t.Errorf("Expected %#v, got %#v", noIdenError, err)
	}
}

func TestDeletePush(t *testing.T) {
	fakeRT := &FakeRoundTripper{message: "", status: http.StatusOK}
	client := newTestClient(fakeRT)
	err := client.DeletePush(Params{"iden": "0xyz"})
	if err != nil {
		t.Errorf("Error, expected nil, got %#v", err)
	}
}

func TestUploadRequestError(t *testing.T) {

}
