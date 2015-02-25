package main

import (
	"bufio"
	"log"
	"net"
	"os/exec"
	"strings"
	"time"
)

func getIPNeighbors(host string) []string {
	cidrMask := getCIDRMask(host)
	LANIPS := []string{}
	ip, ipnet, err := net.ParseCIDR(cidrMask)
	if err != nil {
		log.Fatal(err)
	}
	for ip := ip.Mask(ipnet.Mask); ipnet.Contains(ip); inc(ip) {
		if ip.String() != host {
			LANIPS = append(LANIPS, ip.String())
		}
	}
	//return LANIPS
	return []string{"192.168.1.148"}
}

func inc(ip net.IP) {
	for j := len(ip) - 1; j >= 0; j-- {
		ip[j]++
		if ip[j] > 0 {
			break
		}
	}
}

func getLANIP() string {
	cmd := exec.Command("ifconfig", "en0")
	stdout, _ := cmd.Output()
	out := string(stdout)
	ip := strings.Split(out, "inet ")[1]
	ip = strings.Split(ip, " netmask")[0]
	return net.ParseIP(ip).String()
}

func getCIDRMask(ip string) string {
	octets := strings.Split(ip, ".")
	octets[3] = "0/24"
	return strings.Join(octets, ".")
}

func initiateUsernameExchange(conn net.Conn) (string, error) {
	sendMessage(conn, _USERNAME)
	reader := bufio.NewReader(conn)
	user, err := reader.ReadString('\n')
	if err != nil {
		return "", err
	}
	return strings.TrimSpace(user), nil
}

func (chatApp ChatApp) InitiateConversation(remoteAddress string) error {
	timeout, _ := time.ParseDuration("200ms")
	conn, err := net.DialTimeout("tcp", remoteAddress, timeout)
	if err != nil {
		println("NOP")
		return err
	}
	defer conn.Close()

	reader := bufio.NewReader(conn)
	remoteUser, err := initiateUsernameExchange(conn)
	if err != nil {
		return err
	}
	chatApp.Connections[remoteUser] = conn
	chatApp.Gui.Flush()
	for {
		message, err := reader.ReadString('\n')
		if err != nil || message == "" {
			break
		}
		chatApp.QueueUserMessageForDisplay(_USERNAME, message)
	}
	delete(chatApp.Connections, remoteUser)
	chatApp.Gui.Flush()
	return nil
}

func (chatApp ChatApp) BootstrapPeers() {
	for _, ip := range getIPNeighbors(getLANIP()) {
		remoteAddress := net.JoinHostPort(ip, _PORT)
		go chatApp.InitiateConversation(remoteAddress)
	}
}
