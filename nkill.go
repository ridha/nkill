/*

 Kills all processes listening on the given TCP ports.

*/

package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
)

const (
	PROC_TCP6    = "/proc/net/tcp6"
	PROC_TCP     = "/proc/net/tcp"
	LISTEN_STATE = "0A"
)

type Process struct {
	Name  string
	Pid   string
	State string
	Port  int64
}

func readFile(tcpfile string) []string {
	//  Read the table of tcp connections & remove header
	content, err := ioutil.ReadFile(tcpfile)
	if err != nil {
		log.Fatalln(err, content)
	}
	return strings.Split(string(content), "\n")[1:]
}

func hexToDec(h string) int64 {
	dec, err := strconv.ParseInt(h, 16, 32)
	if err != nil {
		log.Fatalln(err)
	}
	return dec
}

func deleteEmpty(s []string) []string {
	var r []string
	for _, str := range s {
		if str != "" {
			r = append(r, str)
		}
	}
	return r
}

func netstat(portToKill int64) []Process {
	tcpStats := statTCP(portToKill, PROC_TCP)
	tcp6Stats := statTCP(portToKill, PROC_TCP6)
	return append(tcpStats, tcp6Stats...)
}

func statTCP(portToKill int64, tcpfile string) []Process {
	// To get pid of all network process running on system, you must run this script
	// as superuser
	content := readFile(tcpfile)
	var processes []Process

	for _, line := range content {
		if line == "" {
			continue
		}
		parts := deleteEmpty(strings.Split(strings.TrimSpace(line), " "))
		localAddress := parts[1]
		state := parts[3]
		if state != LISTEN_STATE {
			continue
		}
		inode := parts[9]
		localPort := hexToDec(strings.Split(localAddress, ":")[1])
		if localPort != portToKill {
			continue
		}

		pid := getPIDFromInode(inode)
		exe := getProcessExe(pid)
		p := Process{Name: exe, Pid: pid, State: state, Port: localPort}
		processes = append(processes, p)
	}

	return processes
}

func getPIDFromInode(inode string) string {
	// To retrieve the pid, check every running process and look for one using
	// the given inode
	pid := "-"

	d, err := filepath.Glob("/proc/[0-9]*/fd/[0-9]*")
	if err != nil {
		log.Fatalln(err)
	}

	re := regexp.MustCompile(inode)
	for _, item := range d {
		path, _ := os.Readlink(item)
		out := re.FindString(path)
		if len(out) != 0 {
			pid = strings.Split(item, "/")[2]
		}
	}
	return pid
}

func getProcessExe(pid string) string {
	exe := fmt.Sprintf("/proc/%s/exe", pid)
	path, _ := os.Readlink(exe)
	return path
}

func killPort(portToKill int64) {
	killed := false
	for _, conn := range netstat(portToKill) {
		iport, _ := strconv.Atoi(conn.Pid)
		p, _ := os.FindProcess(iport)

		if err := p.Kill(); err != nil {
			log.Println(err)
		} else {
			log.Printf("Killed %s (pid: %s) listening on port %d", conn.Name, conn.Pid, iport)
			killed = true
		}

	}
	if !killed {
		log.Printf("No process found listening on port %d\n", portToKill)
	}
}

func init() {
	log.SetFlags(0)
}

func main() {
	if len(os.Args) <= 1 {
		log.Fatalln("Kills all processes listening on the given TCP ports.\nusage: nkill port")
	}

	// if os.Getpid() != 0 {
	// 	log.Println("WARNING: You are not running this script as superuser, expect some things to not work")
	// }

	for _, port := range os.Args[1:] {
		p, err := strconv.ParseInt(port, 10, 64)
		if err != nil {
			log.Printf("%s is not a valid port number\n", port)
			continue
		}
		killPort(p)
	}

}
