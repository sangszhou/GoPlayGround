@todo

[1] 实现一个 Bep  5

bep protocol
http://www.bittorrent.org/beps/bep_0005.html

https://github.com/nictuku/dht/blob/f929f301e23c84bb4ddaac58b92588d88237a38c/dht.go

希望国庆节期间可以完成这个项目, 当然是使用 go 语言

[2] 需要注意的问题

在编码时，想到几个问题

a. 如果出现递归调用了该怎么办，尤其是在初始化的时候
b. FindSuccess 和 FindPredecessor 在论文中收敛到了一起，在实际中，有没有必要再放开呢?

