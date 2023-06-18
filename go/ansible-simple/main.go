package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"net"
	"os"
	"path"
	"path/filepath"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/pkg/sftp"
	"golang.org/x/crypto/ssh"
	"gopkg.in/yaml.v3"
)

var ipAddrs chan string = make(chan string)
var wg sync.WaitGroup

type HostInfo struct {
	Roomname string   `yaml:"roomname"`
	IpList   []string `yaml:"ip"`
}

type Group struct {
	List []*HostInfo `yaml:"groups"`
}

func OpencfgFile() {
	databytes, err := ioutil.ReadFile("roomname.yaml")
	if err != nil {
		log.Fatal(err)
	}
	g := new(Group)
	if err = yaml.Unmarshal(databytes, &g); err != nil {
		log.Fatal(err)
	}
	for _, hostinfo := range g.List {
		if os.Args[1] == "all" {
			for _, i := range hostinfo.IpList {
				ipAddrs <- hostinfo.Roomname + " " + i
			}
		} else if hostinfo.Roomname == os.Args[1] {
			for _, i := range hostinfo.IpList {
				ipAddrs <- hostinfo.Roomname + " " + i
			}
		}
	}
	close(ipAddrs)
}

func Opsee() {
	for s := range ipAddrs {
		wg.Add(1)
		str1 := strings.Split(s, " ")
		room := &str1[0]
		ip := &str1[1]
		user := &str1[2]
		passwd := &str1[3]
		comm := os.Args[2]
		if os.Args[2] == "scp" {
			go SftpRun(*room, *user, *passwd, *ip)
		} else {
			go Runshell(*room, *user, *passwd, *ip, comm)
		}
	}
	wg.Wait()
	if len(ipAddrs) == 0 {
		fmt.Println("End of channel data reading")
	}
}

func sshSession(user, passwd, host string, port int) (sshSession *ssh.Session, err error) {
	sshClinet, err := connector(user, passwd, host, port)
	if err != nil {
		fmt.Printf("连接失败:%s\n", host)
		return
	}
	if sshSession, err = sshClinet.NewSession(); err != nil {
		fmt.Printf("创建客户端fail:%s\n", host)
		return
	}
	return
}

func connector(user, passwd, host string, port int) (sshClinet *ssh.Client, err error) {
	auth := make([]ssh.AuthMethod, 0)
	auth = append(auth, ssh.Password(passwd))

	clientConfig := &ssh.ClientConfig{
		User:    user,
		Auth:    auth,
		Timeout: 1 * time.Second,
		HostKeyCallback: func(hostname string, remote net.Addr, key ssh.PublicKey) error {
			return nil
		},
	}
	addr := host + ":" + strconv.Itoa(port)
	sshClinet, err = ssh.Dial("tcp", addr, clientConfig)
	if err != nil {
		return
	}
	return
}

func sftpconnect(user, password, host string, port int) (*sftp.Client, error) {
	var (
		auth         []ssh.AuthMethod
		addr         string
		clientConfig *ssh.ClientConfig
		sshClinet    *ssh.Client
		sftpClient   *sftp.Client
		err          error
	)

	auth = make([]ssh.AuthMethod, 0)
	auth = append(auth, ssh.Password(password))
	clientConfig = &ssh.ClientConfig{
		User:    user,
		Auth:    auth,
		Timeout: 30 * time.Second,
		HostKeyCallback: func(hostname string, remote net.Addr, key ssh.PublicKey) error {
			return nil
		},
	}
	addr = host + ":" + strconv.Itoa(port)
	if sshClinet, err = ssh.Dial("tcp", addr, clientConfig); err != nil {
		return nil, err
	}
	if sftpClient, err = sftp.NewClient(sshClinet); err != nil {
		return nil, err
	}
	return sftpClient, nil
}

func SftpRun(room, user, passwd, ip string) {
	defer wg.Done()
	var (
		sftpClient *sftp.Client
	)
	sftpClient, _ = sftpconnect(user, passwd, ip, 22)
	defer sftpClient.Close()
	localFilePath := os.Args[3]
	remoteDir := os.Args[4]
	srcFile, err := os.Open(localFilePath)
	if err != nil {
		log.Fatal(err)
	}
	defer srcFile.Close()
	_, remoteFileName := filepath.Split(localFilePath)
	dstFile, err := sftpClient.Create(path.Join(remoteDir, remoteFileName))
	if err != nil {
		log.Fatal(err)
	}
	defer dstFile.Close()
	reader, _ := ioutil.ReadAll(srcFile)
	if err := ioutil.WriteFile(dstFile.Name(), reader, 0766); err != nil {
		log.Fatal(err)
	}
	fmt.Println(room + "." + ip + "\n" + "copy file finished!")
}

func Runshell(room, user, passwd, ip, comm string) {
	defer wg.Done()
	session, err := sshSession(user, passwd, ip, 22)
	if err != nil {
		log.Println(room+"."+ip+"\n"+"result:", err)
		return
	}
	buf, _ := session.CombinedOutput(comm)
	defer session.Close()
	fmt.Println(room + "." + ip + "\n" + "result:" + string(buf))
}

func main() {
	// ipAddrs <- "1"
	// close(ipAddrs)
	i := 1
	fmt.Print(&i)
	a := &i
	fmt.Print(*a)
	fmt.Println(i, 1)
}
