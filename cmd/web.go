/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"context"
	"embed"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/yangyang5214/grpcui/pkg"
	"io/fs"
	"net/http"
)

//go:embed build/*
var content embed.FS

// webCmd represents the web command
var webCmd = &cobra.Command{
	Use:   "web",
	Short: "web",
	Long:  `grpc-ui`,
	Run: func(cmd *cobra.Command, args []string) {
		log.Infof("listen rpc server %v", addr)
		ctx := context.Background()
		grpcCurl, err := pkg.NewGrpcCurl(ctx, addr)
		if err != nil {
			log.Errorf("init grpc curl error: %v", err)
			return
		}
		web := pkg.NewGrpcWeb(grpcCurl)

		resultContent, err := fs.Sub(content, "build")
		if err != nil {
			log.Errorf("fs .sub error: %v", err)
			return
		}
		r := http.DefaultServeMux

		r.Handle("/", http.FileServer(http.FS(resultContent)))
		r.HandleFunc("/services", web.ListServices)
		r.HandleFunc("/service/methods", web.ListMethod)
		r.HandleFunc("/all/methods", web.AllMethods)
		r.HandleFunc("/method/fake_body", web.FakeBody)
		r.HandleFunc("/send", web.Send)

		log.Info("start web server http://127.0.0.1:8548")

		server := &http.Server{
			Addr:    ":8548",
			Handler: allowCORS(r),
		}

		err = server.ListenAndServe()
		if err != nil {
			log.Errorf("listen error: %v", err)
			return
		}
	},
}

func allowCORS(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusOK)
			return
		}
		// 继续处理其他请求
		next.ServeHTTP(w, r)
	})
}

func init() {
	rootCmd.AddCommand(webCmd)
	webCmd.Flags().StringVarP(&addr, "addr", "a", "10.0.81.250:9000", "grpc server address")
}
