# p2p_tool

#### 介绍

采用Go语言编写的支持各种P2P协议的下载工具

##### BitTorrent DHT 协议

- BT 协议像 TCP/IP 协议一样是一个协议簇
- DHT 协议是在 UDP 通信协议的基础上使用 Kademila（俗称 Kad 算法）算法实现
- Tracker 服务器保存和 torrent 文件相关的 peer 的信息
- 一个 peer 节点是一个实现了 BT 协议并且开启了 TCP 监听端口的 BT 客户端或者服务器
- 一个 node 节点是一个实现了 DHT 协议并且开启了 UDP 监听端口的 BT 客户端或者服务器
- DHT 由很多 node 节点以及这些 node 节点保存的 peer 地址信息组成
- 一个 BT 客户端包括了一个 DHT node 节点，通过这些节点来和 DHT 网络中的其它节点通信来获取 peer 节点的信息，然后再通过 BT 协议从 peer 节点下载文件
- DHT 协议通过从附近的 node 节点获取 peer 信息，而不是从 tracker 服务器获取 peer 信息，这就是所谓的 trackerless.

在 Kademlia 算法中，每个节点互相通信， 会接收到下面四种请求，也会发出以下四种请求：

- PING：测试一个节点是否在线。相当于打个电话，看还能打通不；不能打通就把这个节点记录删除掉
- find_node：查询某个节点
- get_peers： 根据文件hash，问下哪个节点知道文件放哪了
- announce_peer: 告知应该知道文件该存储哪些节点

##### P2P 协议

- [BT下载的未来 (magnet协议)](http://www.ruanyifeng.com/blog/2009/11/future_of_bittorrent.html)
- [BitTorrent 协议簇概述](https://www.addesp.com/archives/5236)
- [网络协议 15 - P2P 协议](https://zhuanlan.zhihu.com/p/87327257)
- [BitTorrent DHT 协议简述](https://www.barretlee.com/blog/2017/06/16/bittorrent-dht-protocal/)
- [入门DHT协议理论篇](https://l1905.github.io/p2p/dht/2021/04/23/dht01/)
- [BitTorrent 分布式散列表（DHT）协议详解](https://www.addesp.com/archives/5428)
- [DHT 协议 - 译](https://www.lyyyuna.com/2016/03/26/dht01/)
  -[DHT协议(翻译)](http://blog.leanote.com/post/simon88/DHT%E5%8D%8F%E8%AE%AE-%E7%BF%BB%E8%AF%91)
- [DHT协议代码分析篇](https://l1905.github.io/p2p/dht/2021/04/24/dht02/)
- [分布式散列表协议 —— Kademlia 详解](https://www.addesp.com/archives/5338)
- [P2P 网络核心技术：Kademlia 协议](https://zhuanlan.zhihu.com/p/40286711)
- [聊聊分布式散列表（DHT）的原理 — — 以 Kademlia（Kad） 和 Chord 为例](https://program-think.mediumcom/%E8%81%8A%E8%81%8A%E5%88%86%E5%B8%83%E5%BC%8F%E6%95%A3%E5%88%97%E8%A1%A8-dht-%E7%9A%84%E5%8E%9F%E7%90%86-%E4%BB%A5-kademlia-kad-%E5%92%8C-chord-%E4%B8%BA%E4%BE%8B-8e648d853288)
- [Kademlia、DHT、KRPC、BitTorrent 协议、DHT Sniffer](https://www.daimajiaoliu.com/daima/34017a519900408)
- [BitTorrent协议与MagNet协议原理](https://www.cnblogs.com/wpjamer/articles/10788222.html)
- [BitTorrent Tracker 协议详解](https://www.addesp.com/archives/5313)
- [Golang从零开发BitTorrent客户端](https://mojotv.cn/go/golang-torrent)
- [BitTorrent 伙伴（Peer）协议详解](https://www.addesp.com/archives/5271)
- [一步一步教你写BT种子嗅探器](https://yushuangqi.com/blog/2016/yi-bu-yi-bu-jiao-ni-xie-btchong-zi-xiu-tan-qi-zhi-yi----yuan-lipian.html)
- [使用WireShark进行磁力链接协议分析](https://www.aneasystone.com/archives/2015/05/analyze-magnet-protocol-using-wireshark.html)

##### P2P libraries Or OpenSource Tools

- [torrent (Go library)](https://github.com/anacrolix/torrent)
- [BitTorrent client (Go library)](https://github.com/cenkalti/rain)
- [BitTorrent DHT (Go tool)](https://github.com/boramalper/magnetico)
- [DHT 磁力爬虫 (Java tool)](https://github.com/BrightStarry/zx-bt)
- [DHT BT种子爬虫 (Java tool)](https://github.com/kaiscript/dht)
- [DHT 磁力爬虫 (Go tool)](https://github.com/shiyanhui/dht/blob/master/README_CN.md)
- [DHT 磁力爬虫 (Python tool)](https://github.com/chenjiandongx/magnet-dht)
- [Go torrent-client](https://github.com/veggiedefender/torrent-client/)

#### Questions

1. Download files info of magnet URI
   > torrent.GotInfo() and then torrent.Info()
2. Download progress query

- The progress of the file being downloaded
- The download progress of the downloaded file has been interrupted

3. How to resume the last download after the interruption
4. The download real-time speed:
   > Torrent.Stats.DataBytesRead