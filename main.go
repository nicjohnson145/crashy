package main

import (
	"fmt"
	"net"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/spf13/viper"
)

func main() {
	viper.AutomaticEnv()
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_", "-", "_"))

	dieAfter := viper.GetDuration("die-after")

	mux := http.NewServeMux()
	mux.Handle("/crash", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Println("recieved crash request, dying")
		os.Exit(1)
	}))

	lis, err := net.Listen("tcp4", ":8080")
	if err != nil {
		fmt.Printf("error starting listener: %v\n", err)
		os.Exit(1)
	}

	srv := http.Server{
		Addr:    ":8080",
		Handler: mux,
	}

	if dieAfter.Nanoseconds() != 0 {
		fmt.Printf("delayed crash configured, will die in %v\n", dieAfter)
		go func() {
			<-time.After(dieAfter)
			fmt.Println("delay elapsed, dying")
			os.Exit(1)
		}()
	}

	fmt.Println("starting server")
	if err := srv.Serve(lis); err != nil {
		fmt.Printf("error serving: %v\n", err)
		os.Exit(1)
	}
}
