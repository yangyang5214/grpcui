/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"context"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/yangyang5214/grpcui/pkg"
	"net/http"
)

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
		http.HandleFunc("/services", web.ListServices)
		http.HandleFunc("/service/methods", web.ListMethod)
		http.HandleFunc("/all/methods", web.AllMethods)
		http.HandleFunc("/method/fake_body", web.FakeBody)
		http.HandleFunc("/send", web.Send)
		log.Info("start web server http://127.0.0.1:8548")
		err = http.ListenAndServe(":8548", nil)
		if err != nil {
			log.Errorf("listen error: %v", err)
			return
		}
	},
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
