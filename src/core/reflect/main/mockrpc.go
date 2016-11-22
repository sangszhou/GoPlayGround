package main
//
//import (
//	"reflect"
//	"bytes"
//	"encoding/gob"
//	"log"
//)
//
//type reqMsg struct {
//	endname  interface{} // name of sending ClientEnd
//	svcMeth  string      // e.g. "Raft.AppendEntries"
//	argsType reflect.Type
//	args     []byte
//	replyCh  chan replyMsg
//}
//
//type replyMsg struct {
//	ok    bool
//	reply []byte
//}
//
//type ClientEnd struct {
//	endname interface{} // this end-point's name
//	ch      chan reqMsg // copy of Network.endCh
//}
//
//// send an RPC, wait for the reply.
//// the return value indicates success; false means the server couldn't be contacted.
//func (e *ClientEnd) Call(svcMeth string, args interface{}, reply interface{}) bool {
//	req := reqMsg{}
//	req.endname = e.endname
//	req.svcMeth = svcMeth
//	req.argsType = reflect.TypeOf(args)
//	req.replyCh = make(chan replyMsg)
//
//	qb := new(bytes.Buffer)
//	qe := gob.NewEncoder(qb)
//	qe.Encode(args)
//	req.args = qb.Bytes()
//
//	// 发送消息
//	e.ch <- req
//
//	// 等待结果返回
//	rep := <-req.replyCh
//
//	if rep.ok {
//		rb := bytes.NewBuffer(rep.reply)
//		rd := gob.NewDecoder(rb)
//		if err := rd.Decode(reply); err != nil {
//			log.Fatalf("ClientEnd.Call(): decode reply: %v\n", err)
//		}
//		return true
//	} else {
//		return false
//	}
//}
//
//
