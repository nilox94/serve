/*
Serve is a very simple static file server in Go.
Based on the Gist https://gist.github.com/paulmach/7271283/2a1116ca15e34ee23ac5a3a87e2a626451424993
by Paul Mach (https://github.com/paulmach)

Usage:
  -d string
        The directory of static file to host (default ".")
  -p string
        Port to serve on (default "8100")
  -v    Print the version

Navigating to http://localhost:8100 will display the index.html
or directory listing file.
*/
package main

import (
	"flag"
	"fmt"
	"log"
	"net"
	"net/http"
	"os"
	"strings"
)

// Increment and remove "+" in release commits.
// Add "+" after release commits.
const version = "v0.1.0+"

func main() {
	// Flags
	port := flag.String("p", "8100", "Port to serve on")
	directory := flag.String("d", ".", "The directory of static file to host")
	printVersion := flag.Bool("v", false, "Print the version")
	flag.Parse()

	// If the "v" flag was used, only print the version and exit
	if *printVersion {
		fmt.Printf("serve version: %v\n", version)
		os.Exit(0)
	}

	http.Handle("/", http.FileServer(http.Dir(*directory)))

	fmt.Printf("\nServing \"%s\" on all network interfaces (0.0.0.0) on HTTP port: %s\n", *directory, *port)

	// Print local IP addresses
	printAddrs(*port)

	log.Fatal(http.ListenAndServe(":"+*port, nil))
}

// printAddrs prints the local network interfaces and their IP addresses
func printAddrs(port string) {
	fmt.Println("\nLocal network interfaces and their IP addresses so you can pass one to your colleagues:\n")
	ifaces, err := net.Interfaces()
	if err != nil {
		log.Fatal(err)
	}
	// We want the interface + IP address list to look like a table:
	//    Interface    |  IPv4 Address   |  IPv6 Address
	// ----------------|-----------------|-------------------
	// eth0            | 192.168.178.123 | ...
	//
	// "docker_gwbridge" and other Docker interfaces like "br-2b303075a67e" have 15 characters.
	// 123.123.123.123 are also 15 characters.
	// the IPv6 Address can be open ended, because its the rightmost value
	fmt.Println("   Interface    |  IPv4 Address   | IPv6 Address   ")
	fmt.Println("----------------|-----------------|----------------")
	fav := ""
	for _, iface := range ifaces {
		fmt.Printf("%-15v |", iface.Name)
		addrs, err := iface.Addrs()
		if err != nil {
			log.Fatal(err)
		}
		// The interesting interfaces like eth0 and wlan0 typically have 2 addresses: one IPv4 and one IPv6 address.
		// But some interfaces just have one of them, or if an interface is deactivated it doesn't have any.
		// We take care of these mentioned possibilities. We ignore the third and more addresses.
		ipv4 := ""
		ipv6 := ""
		for i := 0; i < 2 && i < len(addrs); i++ {
			// In the case of two addresses they could potentially be of the same type.
			// We want to show the first address. overwriteIfEmpty() doesn't overwrite existing values.
			addrWithoutMask := strings.Split(addrs[i].String(), "/")[0]
			if strings.Contains(addrWithoutMask, ":") {
				overwriteIfEmpty(&ipv4, "")
				overwriteIfEmpty(&ipv6, addrWithoutMask)
			} else {
				overwriteIfEmpty(&ipv4, addrWithoutMask)
				overwriteIfEmpty(&ipv6, "")
				if (iface.Name == "eth0" || iface.Name == "wlan0") && fav == "" {
					fav = addrWithoutMask
				}
			}
		}
		fmt.Printf(" %-15v | %v\n", ipv4, ipv6)
	}

	// Show probable favorite
	if fav != "" {
		fmt.Printf("\nYou probably want to share:\nhttp://%v:%v\n", fav, port)
	}
}

// overwriteIfEmpty only overwrites the string s with the string overwrite if s is empty
func overwriteIfEmpty(s *string, overwrite string) {
	if *s == "" {
		*s = overwrite
	}
}
