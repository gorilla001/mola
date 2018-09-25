package main

import (
	"log"
	"net/http"
	"time"

	"golang.org/x/net/context"
)

func main() {
	etcdClient, err := newETCD3Client(etcdConfig{
		CertFile:   "ca/peer.crt",
		KeyFile:    "ca/peer.key",
		CAFile:     "ca/ca.crt",
		ServerList: []string{"192.168.210.108:2379"},
	})

	if err != nil {
		log.Fatal("Error: cannot connec to etcd:", err)
	}
	defer etcdClient.Close()

	go func() {
		for {
			_, err := etcdClient.Put(context.TODO(), "logkeywords", "something")
			if err != nil {
				log.Println(err)
			}
			time.Sleep(time.Second)
		}
	}()

	go func() {
		watchChan := etcdClient.Watch(context.Background(), "logkeywords")
		for watchResp := range watchChan {
			for _, event := range watchResp.Events {
				log.Printf("Event received! %s executed on %q with value %q\n", event.Type, event.Kv.Key, event.Kv.Value)
			}
		}
	}()

	http.HandleFunc("/", serveHTTP)

	log.Println("Listen on :24444")
	http.ListenAndServe(":24444", nil)
}

func serveHTTP(w http.ResponseWriter, r *http.Request) {
}
