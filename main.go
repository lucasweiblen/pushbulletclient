package main

import (
	"fmt"

	"github.com/lucasweiblen/pushbulletclient/client"
)

func main() {
	///cli := client.NewClient("token")
	//resp, code, _ := cli.GetContacts()
	//fmt.Println(resp)
	//fmt.Println(code)
	//resp, code, _ = cli.CreateContact(client.Params{"name": "joao", "email": "foo@joao.com"})
	//fmt.Println(resp)
	//fmt.Println(code)
	//resp, code, _ = cli.UpdateContact(client.Params{"iden": "ujDigsFMxWesjAjjvzVOLY", "name": "blindpigs"})
	//fmt.Println(resp)
	//fmt.Println(code)
	//resp, code, _ = cli.DeleteContact(client.Params{"iden": "ujDigsFMxWesjz0F6E05sa"})
	//fmt.Println(resp)
	//fmt.Println(code)
	cli := client.NewClient("swbpcaTIjyV5eAYZnjfL2GZqFiiqrBHH")
	//contact, _ := cli.CreateContact(client.Params{"name": "joao", "email": "foo@joao.com"})
	//fmt.Println(contact)
	//updated, _ := cli.UpdateContact(client.Params{"iden": "ujvSxVpCjh6sjAiVsKnSTs", "name": "foo_updated"})
	//fmt.Println(updated)
	//device, _ := cli.CreateDevice(client.Params{"nickname": "teste", "type": "stream"})
	//fmt.Println(device)
	//updated, _ := cli.UpdateDevice(client.Params{"iden": "ujvSxVpCjh6sjz2gRnO9Aq", "nickname": "deviceupdated"})
	//bla, _ := cli.DeleteDevice(client.Params{"iden": "ujvSxVpCjh6sjAsoeMFET6"})
	//fmt.Println(updated)
	pushes, _ := cli.GetPushes()
	fmt.Println(pushes.Pushes)
	//b, _ := cli.CreatePush(client.Params{"type": "link", "title": "baz", "body": "foo", "url": "blaa"})
	//fmt.Println(b)
	//r, _ := cli.UpdatePush(client.Params{"iden": "ujvSxVpCjh6sjAcyrQ9Cmq", "title": "modified"})
	//fmt.Println(r)
	//o, _ := cli.CreatePush(client.Params{"type": "address", "name": "ok", "address": "bla"})
	//fmt.Println(o)
	c, _ := cli.CreatePush(client.Params{"type": "list", "title": "titulo", "items": "bla"})
}
