package dnssd

import (
	"log"
	"os"
	"os/signal"
	"schnoddelbotz/k12-booter/utility"
	"time"

	"github.com/oleksandr/bonjour"
)

// all based on example code at https://github.com/oleksandr/bonjour/blob/master/README.md

func Browse() {
	resolver, err := bonjour.NewResolver(nil)
	if err != nil {
		log.Println("Failed to initialize resolver:", err.Error())
		os.Exit(1)
	}

	results := make(chan *bonjour.ServiceEntry)

	go func(results chan *bonjour.ServiceEntry, exitCh chan<- bool) {
		for e := range results {
			log.Printf("%s", e.Instance)
			exitCh <- true
			time.Sleep(1e9)
			os.Exit(0)
		}
	}(results, resolver.Exit)

	err = resolver.Browse("_k12booter._tcp", "local.", results)
	if err != nil {
		log.Println("Failed to browse:", err.Error())
	}

	select {}
}

func BrowseSingle() *bonjour.ServiceEntry {
	resolver, err := bonjour.NewResolver(nil)
	utility.Fatal(err)
	results := make(chan *bonjour.ServiceEntry)
	err = resolver.Browse("_k12booter._tcp", "local.", results)
	utility.Fatal(err)
	for e := range results {
		if e.AddrIPv4 == nil && e.Instance != "" && e.Domain != "" {
			// first response of k12-booter after startup+registration had IP=nil :-/
			// happened on windows and linux if waiting for a teacher to pop up
			// joining later would work. hack, did not look ... just retry m(
			time.Sleep(1 * time.Second)
			return BrowseSingle()
		}
		resolver.Exit <- true
		return e
	}
	return nil
}

func Lookup(instance string) {
	resolver, err := bonjour.NewResolver(nil)
	if err != nil {
		log.Println("Failed to initialize resolver:", err.Error())
		os.Exit(1)
	}

	results := make(chan *bonjour.ServiceEntry)

	go func(results chan *bonjour.ServiceEntry, exitCh chan<- bool) {
		for e := range results {
			log.Printf("%s", e.Instance)
			exitCh <- true
			time.Sleep(1e9)
			os.Exit(0)
		}
	}(results, resolver.Exit)

	err = resolver.Lookup(instance, "_k12booter._tcp", "", results)
	if err != nil {
		log.Println("Failed to browse:", err.Error())
	}

	select {}
}

func RegisterTeacherService(teacherName string) {
	// could/should be run from within k12-booter cui / CLI
	// console should provide feedback about pupils connected etc.
	// imagine k12-booter CLI commands:
	// classroom open https://mit.edu/
	// ... and all pupils' devices would do just that.

	// Run registration (blocking call)
	// &{ServiceRecord:{Instance:jan Service:_k12booter._tc|
	// |p Domain:local. serviceName: serviceInstanceName: serviceTypeName:} HostName:j|
	// |mbair.local. Port:9999 Text:[txtv=1 app=k12booter] TTL:3200 AddrIPv4:192.168.7|
	// |8.144 AddrIPv6:fd2f:bf00:ae78:...}
	s, err := bonjour.Register(teacherName, "_k12booter._tcp", "", 8888, []string{"txtv=1", "app=k12booter"}, nil)
	if err != nil {
		log.Fatalln(err.Error())
	}
	// security m( ...
	// - initial trust / handshake ... MUST exchange ssh key ONCE for later trust
	// - sign control commands using that trusted key?

	//////////////////////////// ahem ... todo -> cui F10
	// Ctrl+C handling
	handler := make(chan os.Signal, 1)
	signal.Notify(handler, os.Interrupt)
	for sig := range handler {
		if sig == os.Interrupt {
			s.Shutdown()
			time.Sleep(1e9)
			break
		}
	}
}
