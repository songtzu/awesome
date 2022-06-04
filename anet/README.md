
all test avg is ms.

### tcp_client_block_test.go benchmark
test on laptop i4700,16G
pSize = 1
totalCount:10000,failCount:0,passCount:10000, timeCost:585, avg:0.058500

pSize = 2
totalCount:20000,failCount:0,passCount:20000, timeCost:1243, avg:0.062150

pSize = 5
totalCount:50000,failCount:0,passCount:50000, timeCost:4090, avg:0.081800

pSize = 50
totalCount:50000,failCount:0,passCount:50000, timeCost:15924, avg:0.318480

test on desktop 8700k 24G
pSize= 2
totalCount:20000,failCount:0,passCount:20000, timeCost:621, avg:0.031050


### tcp_client_cb_test.go benchmark
test on laptop i4700,16G
totalCount:10000,failCount:0,passCount:9999, timeCost:64, avg:0.006400

test on desktop i8700k 24G
totalCount:10000,failCount:0,passCount:9999, timeCost:35, avg:0.003500

thead:1,totalCount:1000000,failCount:0,passCount:1000000, timeCost:7309, avg:0.007309

thead:2,totalCount:2000000,failCount:0,passCount:2000000, timeCost:7174, avg:0.003587

thead:2,totalCount:2000000,failCount:0,passCount:2000000, timeCost:6856, avg:0.003428

thead:2,totalCount:10000000,failCount:0,passCount:10000000, timeCost:36419, avg:0.003642

thead:3,totalCount:15000000,failCount:0,passCount:15000000, timeCost:45576, avg:0.003038

cpu:80%
thead:3,totalCount:15000000,failCount:0,passCount:15000000, timeCost:44798, avg:0.002987

thead:3,totalCount:3000000,failCount:0,passCount:3000000, timeCost:8226, avg:0.002742

thead:3,totalCount:3000000,failCount:0,passCount:3000000, timeCost:8238, avg:0.002746

thead:4,totalCount:4000000,failCount:0,passCount:4000000, timeCost:10150, avg:0.002537

thead:4,totalCount:4000000,failCount:0,passCount:4000000, timeCost:11276, avg:0.002819

thead:4,totalCount:20000000,failCount:0,passCount:20000000, timeCost:45362, avg:0.002268

thead:4,totalCount:20000000,failCount:0,passCount:20000000, timeCost:43909, avg:0.002195, cpu:90%~95%



tcp_client_block_test.go

test on desktop i8700k 24G

thead:1,totalCount:10000,failCount:0,passCount:10000, timeCost:353, avg:0.035300

thead:2,totalCount:20000,failCount:0,passCount:20000, timeCost:296, avg:0.014800

thead:2,totalCount:200000,failCount:0,passCount:200000, timeCost:3241, avg:0.016205

thead:3,totalCount:300000,failCount:0,passCount:300000, timeCost:3590, avg:0.011967

thead:4,totalCount:400000,failCount:0,passCount:400000, timeCost:4250, avg:0.010625

thead:4,totalCount:400000,failCount:0,passCount:400000, timeCost:4412, avg:0.011030

thead:5,totalCount:500000,failCount:0,passCount:500000, timeCost:4819, avg:0.009638

thead:5,totalCount:500000,failCount:0,passCount:500000, timeCost:4762, avg:0.009524

thead:8,totalCount:800000,failCount:0,passCount:800000, timeCost:5982, avg:0.007477

thead:10,totalCount:1000000,failCount:0,passCount:1000000, timeCost:6596, avg:0.006596, cpu:100%

thead:10,totalCount:1000000,failCount:0,passCount:1000000, timeCost:6578, avg:0.006578, cpu:100%

thead:12,totalCount:1200000,failCount:0,passCount:1200000, timeCost:7402, avg:0.006168, cpu:100%

thead:12,totalCount:1200000,failCount:0,passCount:1200000, timeCost:7327, avg:0.006106

thead:14,totalCount:1400000,failCount:0,passCount:1400000, timeCost:7815, avg:0.005582

thead:20,totalCount:2000000,failCount:0,passCount:2000000, timeCost:9772, avg:0.004886, cpu:100%