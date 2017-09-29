package main

import (
	"reflect"
	"bytes"
	"encoding/gob"
	"log"
	"fmt"
)


type reqMsg struct {
	endname  interface{} // name of sending ClientEnd
	svcMeth  string      // e.g. "Raft.AppendEntries"
	argsType reflect.Type
	args     []byte
	replyCh  chan replyMsg
}

type replyMsg struct {
	ok    bool
	reply []byte
}

func arg2ReqMsg(args interface{}, methodName string) reqMsg {
	req := reqMsg{}
	req.svcMeth = methodName
	req.argsType = reflect.TypeOf(args)

	qb := new(bytes.Buffer)
	qe := gob.NewEncoder(qb)
	qe.Encode(args)
	req.args = qb.Bytes()

	return req
}

type Service struct {
	name    string
	rcvr    reflect.Value
	typ     reflect.Type
	methods map[string]reflect.Method
}

//参数是一个 类和它的方法
func MakeService(rcvr interface{}) *Service {

	//log.Printf("type of argument is %s", reflect.TypeOf(rcvr))

	svc := &Service{}
	svc.typ = reflect.TypeOf(rcvr)
	svc.rcvr = reflect.ValueOf(rcvr)
	svc.name = reflect.Indirect(svc.rcvr).Type().Name()
	//反射的方法
	svc.methods = map[string]reflect.Method{}

	log.Printf("num of method is %d", svc.typ.NumMethod())

	for m := 0; m < svc.typ.NumMethod(); m ++ {
		method := svc.typ.Method(m)
		//方法为什么还有类型，是调用方的类型么
		mtype := method.Type
		mname := method.Name

		fmt.Printf("[%v] pkgPath [%v] ni [%v] 1k [%v] 2k [%v] no [%v]\n",
			mname, method.PkgPath, mtype.NumIn(), mtype.In(1).Kind(), mtype.In(2).Kind(), mtype.NumOut())

		if method.PkgPath != "" || // capitalized?
			mtype.NumIn() != 3 ||
		//mtype.In(1).Kind() != reflect.Ptr ||
			mtype.In(2).Kind() != reflect.Ptr ||
			mtype.NumOut() != 0 { // 不能有返回值
			// the method is not suitable for a handler
			fmt.Printf("bad method: %v\n", mname)
		} else {
			// the method looks like a handler
			svc.methods[mname] = method
		}
	}

	return svc
}

func (svc *Service) dispatch(methodName string, req reqMsg) replyMsg {
	if method, ok := svc.methods[methodName]; ok {
		// prepare space into which to read the argument.
		// the Value's type will be a pointer to req.argsType.
		args := reflect.New(req.argsType)

		// decode the argument.
		ab := bytes.NewBuffer(req.args)
		ad := gob.NewDecoder(ab)
		ad.Decode(args.Interface())

		// allocate space for the reply.
		// 第二个参数肯定是返回值，就是这么定义的
		replyType := method.Type.In(2)
		replyType = replyType.Elem()
		replyv := reflect.New(replyType)

		// call the method
		// 反射调用方法的用法
		function := method.Func
		function.Call([]reflect.Value {svc.rcvr, args.Elem(), replyv} )

		// encode the reply.
		rb := new(bytes.Buffer)
		re := gob.NewEncoder(rb)
		re.EncodeValue(replyv)

		return replyMsg {true, rb.Bytes()}
	} else {
		choices := []string{}
		for k, _ := range svc.methods {
			choices = append(choices, k)
		}
		log.Fatalf("labrpc.Service.dispatch(): unknown method %v in %v; expecting one of %v\n",
			methodName, req.svcMeth, choices)
		return replyMsg{false, nil}
	}
}