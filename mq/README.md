

### benchmark

client_pub_test.go
6.1号之前的版本，avg大于13ms

移除了for-sleep逻辑，修改成chan wake方式处理listcache，
移除了anetWriteMessageWithCallback方法对connectionType为server模式的支持
由于anet移除了server版本对WriteMessageWithCallback的支持（此模式增加了大量的注册回调map数据，并且每个IO都需要查询map，增加了大量开销），
综合考虑，移除x_pub_impl内部的WriteMessageWithCallback的调用，改为使用IOnProcessPack回调中查询源publisher，并转发结果。

以上优化调整，将benchmark速度提升了20倍（save 95% of avg time ms）

client_pub_test.go
totalCount:10000, timeoutCount:0, rightCount:10000, time cost ms:6866, avg:0.686600

### 设计目标

考虑到anet的callback模式每秒接近19万的测试值，mq的普通扇出模式设定每秒5万的优化目标。


## client_pub_test.go

thead:1,totalCount:50000,failCount:0,passCount:50000, timeCost:7044, avg:0.140880
thead:2,totalCount:100000,failCount:0,passCount:100000, timeCost:7397, avg:0.073970
thead:4,totalCount:200000,failCount:0,passCount:200000, timeCost:8701, avg:0.043505
thead:6,totalCount:300000,failCount:0,passCount:300000, timeCost:11339, avg:0.037797
thead:8,totalCount:400000,failCount:0,passCount:400000, timeCost:13379, avg:0.033447
thead:16,totalCount:800000,failCount:0,passCount:800000, timeCost:23360, avg:0.029200