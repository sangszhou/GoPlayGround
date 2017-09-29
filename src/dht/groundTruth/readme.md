## Implementation Notes

[1] cluster_buggly 的实现是有 bug 的

至少不能分为 upper part 和 down part 来找 predecessor，并且函数的实现也不太优雅，更好的做法是先排个序，按照 clock wise 或者
reverse clockwise, 然后再找元素

此文件设置 ground truth