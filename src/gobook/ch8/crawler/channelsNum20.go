package crawler

//keep 20 go routines without using tokens

func main()  {
	// go routine 收集到的 urls
	worklist := make(chan []string)

	// 从 workList 中过滤出的，没有见到过的 urls
	unseenUrl := make(chan string)

	//启动 20 个 routine
	for i:=0; i < 20; i ++ {
		go func() {
			for link := range unseenUrl {
				discoveriedUrls := craw(link)
				go func() {worklist <- discoveriedUrls}()
			}
		}()
	}

	//填充数据
	seen := make(map[string]bool)
	for list := range worklist {
		for _, link := range list {
			if !seen[link] {
				seen[link] = true
				unseenUrl <- link
			}
		}
	}
}

func craw(link string) []string {
	//empty
	return nil
}