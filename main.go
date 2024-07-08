package main

import (
	"bufio"
	"flag"
	"fmt"
	"github.com/google/gopacket"
	"github.com/google/gopacket/pcap"
	"github.com/google/gopacket/pcapgo"
	"log"
	"os"
	"os/signal"
	"strconv"
	"strings"
	"syscall"
	"time"
)

func listInterfaces() []pcap.Interface {
	devices, err := pcap.FindAllDevs()
	if err != nil {
		log.Fatal(err)
	}
	return devices
}

func containsIPAddress(device pcap.Interface, ip string) bool {
	for _, address := range device.Addresses {
		if address.IP != nil && address.IP.String() == ip {
			return true
		}
	}
	return false
}

func printDeviceInfo(devices []pcap.Interface) {
	for i, device := range devices {
		fmt.Printf("%d: %s - %s\n", i, device.Name, device.Description)
		if len(device.Addresses) > 0 {
			for _, address := range device.Addresses {
				if address.IP != nil {
					fmt.Printf("\tIP Address: %s\n", address.IP)
				}
			}
		} else {
			fmt.Println("\tNo IP Address found.")
		}
	}
}

func chooseInterface(devices []pcap.Interface) pcap.Interface {
	printDeviceInfo(devices)

	reader := bufio.NewReader(os.Stdin)
	fmt.Print("Choose interface number: ")
	input, _ := reader.ReadString('\n')
	input = strings.TrimSpace(input) // Remove newline and spaces
	index, err := strconv.Atoi(input)
	if err != nil || index < 0 || index >= len(devices) {
		log.Fatalf("Invalid interface number: %s ", input)
	}

	return devices[index]
}

func main() {
	// Define command-line parameters
	iface := flag.String("i", "", "Interface to capture packets from (can be interface name, description, or IP address)")
	filename := flag.String("w", "", "Output pcap file")

	// Custom usage function to include tips
	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "Usage of %s:\n", os.Args[0])
		flag.PrintDefaults()
		fmt.Println("\nTips:")
		fmt.Fprintf(os.Stderr, "%s -i eth0 -w out.pcap \n", os.Args[0])
		fmt.Println("- Interface Name: The specific name of the network interface, such as 'eth0' or 'wlan0'.")
		fmt.Println("- Description: The user-friendly description of the network interface, such as 'Intel(R) Ethernet Connection'.")
		fmt.Println("- IP Address: The IP address associated with the network interface, such as '192.168.1.10'.")
		fmt.Println("- Output pcap file: The filename to save captured packets, such as 'output.pcap'.")
	}

	flag.Parse()

	// Check if a network interface was specified
	var selectedDevice pcap.Interface
	if *iface == "" {
		// List all interfaces and choose one
		devices := listInterfaces()
		if len(devices) == 0 {
			log.Fatal("No network interfaces found")
		}
		selectedDevice = chooseInterface(devices)
	} else {
		// Verify the specified network interface exists
		devices := listInterfaces()
		for _, device := range devices {
			if device.Name == *iface || device.Description == *iface || containsIPAddress(device, *iface) {
				selectedDevice = device
				break
			}
		}
		if selectedDevice.Name == "" {
			log.Fatalf("Specified network interface %s does not exist", *iface)
		}
	}

	// Check if a filename was specified
	if *filename == "" {
		*filename = fmt.Sprintf("output_%s.pcap", time.Now().Format("20060102150405"))
	}

	// Open the specified network interface
	handle, err := pcap.OpenLive(selectedDevice.Name, 1600, true, pcap.BlockForever)
	if err != nil {
		log.Fatal(err)
	}
	defer handle.Close()

	// Open output file
	f, err := os.Create(*filename)
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()

	// Create pcap writer
	w := pcapgo.NewWriter(f)
	w.WriteFileHeader(1600, handle.LinkType())

	// Create signal channel
	sigs := make(chan os.Signal, 1)
	done := make(chan bool, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)

	// Initialize statistics
	packetCount := 0
	sourceAddresses := make(map[string]bool)
	destAddresses := make(map[string]bool)

	// Start a goroutine to listen for signals
	go func() {
		sig := <-sigs
		fmt.Println()
		fmt.Println(sig)
		handle.Close()
		done <- true
	}()

	// Start a goroutine to dynamically output capture information
	go func() {
		for {
			select {
			case <-done:
				return
			default:
				fmt.Printf("\rPackets captured: %d, Source addresses: %d, Destination addresses: %d", packetCount, len(sourceAddresses), len(destAddresses))
				time.Sleep(1 * time.Second)
			}
		}
	}()

	// Capture and write packets
	packetSource := gopacket.NewPacketSource(handle, handle.LinkType())
	fmt.Println("Starting packet capture, press Ctrl+C to stop...")
	for packet := range packetSource.Packets() {
		w.WritePacket(packet.Metadata().CaptureInfo, packet.Data())
		packetCount++
		if netLayer := packet.NetworkLayer(); netLayer != nil {
			sourceAddresses[netLayer.NetworkFlow().Src().String()] = true
			destAddresses[netLayer.NetworkFlow().Dst().String()] = true
		}
	}

	// Wait for signal
	<-done
	fmt.Println("\nPacket capture stopped")
}
