package network

import (
	"fmt"
	"log"
	"net"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/hashicorp/memberlist"
)

var cluster *Cluster

type Cluster struct {
	list *memberlist.Memberlist
	stop chan bool
}
type Members struct {
	Name string
	IP   net.IPAddr
}

func NewCluster(local bool) *Cluster {
	mconfig := memberlist.DefaultLANConfig()
	if local {
		mconfig = memberlist.DefaultLocalConfig()
	}
	fmt.Printf("Creating cluster on tcp://%s:%d", mconfig.BindAddr, mconfig.BindPort)
	list, err := memberlist.Create(memberlist.DefaultWANConfig())
	if err != nil {
		panic("Failed to create memberlist: " + err.Error())
	}
	return &Cluster{
		list: list,
		stop: make(chan bool),
	}
}

func (s *Cluster) Join(host string) error {
	count, err := s.list.Join([]string{host})
	if err != nil {
		return err
	}
	fmt.Printf("Found %d other nodes", count)
	return nil
}

func (s *Cluster) Members() []Members {
	// Ask for members of the cluster
	r := make([]Members, 0)
	for _, member := range s.list.Members() {
		ip, err := net.ResolveIPAddr("tcp", member.Address())
		if err != nil {
			log.Printf("[error] %v", err)
		}
		r = append(r, Members{
			Name: member.Name,
			IP:   *ip,
		})
	}
	return r
}

func (s *Cluster) PeerCount() int {
	return s.list.NumMembers()
}
func Join(address string) error {
	return cluster.Join(address)
}

func ClusterSetup(local bool) {
	cluster = NewCluster(local)
	go cluster.poll()
}

func Serve() {
	sigchan := make(chan os.Signal, 1)
	signal.Notify(sigchan)
	exitsig := make(chan int)

	go func() {
		for {
			sig := <-sigchan
			switch sig {
			case syscall.SIGKILL:
			case syscall.SIGTERM:
			case syscall.SIGINT:
				log.Printf("[info] terminating cluster, please stand by")
				clusterStop()
				os.Exit(0)
				return
			default:
				continue
			}
		}
	}()

	exitcode := <-exitsig
	log.Printf("[info] exiting with %d", exitcode)
	os.Exit(exitcode)
}

func (s *Cluster) poll() {
	defer func() {
		s.list.Leave(10 * time.Second)
		s.list.Shutdown()
	}()
	for {
		select {
		case <-s.stop:
			return
		default:
			log.Printf("[info] node update: %d nodes in list", s.list.NumMembers())
			time.Sleep(10 * time.Second)
		}
	}
}

func clusterStop() {
	cluster.stop <- true
}
