package servers

import "log"

type PoliceMan struct {
	Address string
	name string
}

func (man *PoliceMan)Call911(request *RequestArgs, reply *ReplyArgs) error {
	log.Printf("call 911, people said that %v", request.Content)
	*reply = ReplyArgs{
		Content: "us policy will do that for you",
	}

	return nil
}

func (man *PoliceMan)Call110(request *RequestArgs, reply *ReplyArgs) error  {
	log.Printf("call 110, people said that %v", request.Content)

	*reply = ReplyArgs{
		Content: "人民警察为您服务",
	}

	return nil
}
