/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"context"
	"embed"
	"github.com/gorilla/mux"
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
		r := mux.NewRouter()

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
		//w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		//w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
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

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// webCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// webCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
