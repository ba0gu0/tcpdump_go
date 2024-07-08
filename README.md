# tcpdump_go

## 介绍

* Golang编写，但文件，直接抓取网卡的数据包，保存到pcap文件中。
* 主要是为了不安装wireshark之类的工具直接抓取流量，所以不考虑添加各种过滤参数，直接保存所有流量。
* 使用此工具，需要在系统中安装npcap或者libpcap。
* 支持win7+、win2008+、linux、macos。
* 其他系统、架构自行编译

## 依赖

* 在windows上运行时，推荐直接使用32位版本（64也行）。
* 在windows上必须安装npcap，不然无法运行。
* 在linux必须安装libpcap库，不然无法运行。

```shell

# apt 管理员权限
apt install -y libpcap-dev 

# yum 管理员权限
yum install -y libpcap-devel

# rpm和deb包离线安装方式，请自行在项目官网下载对应的包文件。

# windows 安装 npcap 管理员权限
直接双击dist目录下的 npcap-1.79.exe 默认安装即可 不要乱点

# windows 无图形化 安装 npcap 管理员权限
npcap-1.79.exe /S loopback_support=yes winpcap_mode=yes npf_start_on_boot=yes

```

## 使用方法

* 手动指定网卡和保存文件
* 网卡支持网卡名字、网卡IP、网卡说明。

```shell

# linux macos 直接使用网卡名字，或者IP地址
./tcpdump_go -i eth0 -w 111.pcap

./tcpdump_go -i 102.11.11.100 -w linux.pcap

# windows因为无法获取网卡名字，只支持 网卡说明 和 网卡IP地址
./tcpdump_go -i 10.211.11.100 -w windows.pcap

# 先执行 ipconfig /all
./tcpdump_go -i "Parallels VirtIO Ethernet Adapter" -w windows.pcap


```

* 有交互式的终端，可以选择网卡进行操作
* 直接运行程序即可

```shell
# 直接运行即可
tcpdump_go-windows-386.exe

# 示例
./tcpdump_go-darwin-arm64
0: en0 -
	IP Address: de81::13a1:a82e:7d24:1a44
	IP Address: 172.20.10.2
	
1: bridge102 -
	IP Address: 10.37.129.2
	IP Address: 13a1:2c26:13a1:1::1
	
2: lo0 -
	IP Address: 127.0.0.1
	IP Address: ::1
	IP Address: fe80::1
	
Choose interface number: 1
Starting packet capture, press Ctrl+C to stop...
Packets captured: 4994, Source addresses: 7, Destination addresses: 7
interrupt

Packet capture stopped
```
