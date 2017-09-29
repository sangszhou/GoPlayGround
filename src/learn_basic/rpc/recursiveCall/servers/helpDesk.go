package servers

import (
	"strings"
	"net/rpc"
	"log"
)

type Information struct {
	caller []string
}

func (infor *Information) Call(request *RequestArgs, reply *ReplyArgs) error {
	log.Print("called by client")
	if strings.Contains(request.Content, "direction") {
		/**
			这个地方必须是 *reply = xxx 而不能是 reply = &xxx
			坑
		 */
		*reply = ReplyArgs {
			Content: "walk straight this street and turn left",
		}
	} else if strings.Contains(request.Content,"911") {
		callPolice("us", request, reply)
	} else if strings.Contains(request.Content, "110") {
		callPolice("china", request, reply)
	} else {
		log.Print("cannot help you")
	}

	return nil
}

func callPolice(country string, request *RequestArgs, reply *ReplyArgs)  {
	client, err := rpc.DialHTTP("tcp", "localhost:9110")
	if err != nil {
		log.Printf("error connect 911")
	}

	// 可以直接传递到孩子节点
	if country == "us" {
		err = client.Call("PoliceMan.Call911", request, reply)
	} else {

		err = client.Call("PoliceMan.Call110", request, reply)
	}

	if err != nil {
		log.Print("failed to call police man")
	}
}