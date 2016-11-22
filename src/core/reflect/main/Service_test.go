package main

import (
	"fmt"
	"testing"
)

type Person struct {
	Age int `json:"age"`
	Name string `json:"name"`
}

func (person *Person) GetName(arg interface{}, reply *string)  {
	*reply = person.Name
}

func (person *Person) GetAge(arg interface{}, reply *int)  {
	*reply = person.Age
}

func (person *Person) Speak(words string, reply *string)  {
	//fmt.Printf("%+v\n", *person)
	fmt.Printf(person.Name + "\n")
	*reply = person.Name + " said: " + words

}

func TestPersonStruct(t *testing.T) {
	boy := Person { 25, "xinszhou"}
	var reply string
	boy.Speak("hello world", &reply)
}

func TestMakeService(t *testing.T) {
	boy := Person{25, "xinszhou"}

	service := MakeService(&boy)

	getNameRequest := arg2ReqMsg("", "GetName")
	getNameReply := service.dispatch(getNameRequest.svcMeth, getNameRequest)
	fmt.Printf("name is [%v]\n", string(getNameReply.reply))

	speakRequest := arg2ReqMsg("hello world", "Speak")
	speakReply := service.dispatch(speakRequest.svcMeth, speakRequest)
	fmt.Printf("speak reply: [%v]\n", string(speakReply.reply))
}