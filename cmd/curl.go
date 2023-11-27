/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"context"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/yangyang5214/grpcui/pkg"
)

// curlCmd represents the curl command
var curlCmd = &cobra.Command{
	Use:   "curl",
	Short: "grpc curl",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		ctx := context.Background()
		grpcCurl, err := pkg.NewGrpcCurl(ctx, "10.0.81.250:9000")
		if err != nil {
			log.Errorf("init grpc curl error: %v", err)
			return
		}
		services, err := grpcCurl.ListServices()
		if err != nil {
			log.Errorf("get list of methods: %v", err)
			return
		}
		for _, service := range services {
			log.Info(service)
		}
	},
}

func init() {
	rootCmd.AddCommand(curlCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// curlCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// curlCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
