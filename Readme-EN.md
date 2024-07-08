
# tcpdump_go

## Introduction

* Written in Golang, this tool directly captures packets from the network interface and saves them to a pcap file.
* The main goal is to capture traffic without installing tools like Wireshark, so no filtering parameters are added, and all traffic is saved.
* To use this tool, you need to have npcap or libpcap installed on your system.
* Supports Windows 7+, Windows 2008+, Linux, macOS.
* For other systems and architectures, compile it yourself.

## Dependencies

* For running on Windows, it is recommended to use the 32-bit version (64-bit also works).
* On Windows, npcap must be installed, otherwise it won't run.
* On Linux, the libpcap library must be installed, otherwise it won't run.

```shell
# apt with administrator privileges
apt install -y libpcap-dev 

# yum with administrator privileges
yum install -y libpcap-devel

# For offline installation of rpm and deb packages, please download the corresponding package files from the official website.

# Windows installation of npcap with administrator privileges
Directly double-click npcap-1.79.exe in the dist directory and install with default settings. Do not change settings.

# Windows command-line installation of npcap with administrator privileges
npcap-1.79.exe /S loopback_support=yes winpcap_mode=yes npf_start_on_boot=yes
```

## Usage

* Manually specify the network interface and save file.
* The network interface can be specified by name, IP address, or description.

```shell
# On Linux and macOS, use the network interface name or IP address
./tcpdump_go -i eth0 -w 111.pcap

./tcpdump_go -i 102.11.11.100 -w linux.pcap

# On Windows, use the network interface description or IP address
./tcpdump_go -i 10.211.11.100 -w windows.pcap

# First run ipconfig /all
./tcpdump_go -i "Parallels VirtIO Ethernet Adapter" -w windows.pcap
```

* It has an interactive terminal to select the network interface for operations.
* Just run the program directly.

```shell
# run
tcpdump_go-windows-386.exe

# tips 
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

