package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"github.com/mattjmcnaughton/go-youtube-feed/internal/feed"
	"github.com/mattjmcnaughton/go-youtube-feed/internal/server"
	"github.com/mattjmcnaughton/go-youtube-feed/internal/youtube"
)

var rootCmd = &cobra.Command{Use: "go-youtube-feed"}

func init() {
	viperConfig := viper.New()

	viperConfig.SetConfigName(".env")
	viperConfig.SetConfigType("dotenv")
	viperConfig.AddConfigPath(".")

	err := viperConfig.ReadInConfig()
	if err != nil {
		log.Fatalf("Error reading config file, %s", err)
	}

	var checkFeed bool
	var numEntries int
	var generateFeedCmd = &cobra.Command{
		Use:   "generate-feed [handle]",
		Short: "Generate an Atom feed for the given handle.",
		Args:  cobra.ExactArgs(1),
		Run: func(_ *cobra.Command, args []string) {
			handle := args[0]

			err := generateFeed(handle, checkFeed, numEntries, viperConfig)
			if err != nil {
				fmt.Printf("Error executing generateFeed: %s", err)
				os.Exit(1)
			}
		},
	}
	generateFeedCmd.Flags().BoolVarP(&checkFeed, "check-feed", "c", true, "Validate the feed")
	generateFeedCmd.Flags().IntVarP(&numEntries, "num-entries", "n", 0, "Print first N titles/urls")

	var portNumber int
	var serverCmd = &cobra.Command{
		Use:   "server",
		Short: "Run a server for interacting (instead of using CLI)",
		Args:  cobra.ExactArgs(0),
		Run: func(_ *cobra.Command, args []string) {
			runServer(portNumber)
		},
	}

	serverCmd.Flags().IntVarP(&portNumber, "port-number", "p", 8080, "Port number")

	rootCmd.AddCommand(generateFeedCmd)
	rootCmd.AddCommand(serverCmd)
}

func main() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func generateFeed(handle string, checkFeed bool, numEntries int, viperConfig *viper.Viper) error {
	apiKey := viperConfig.GetString("YOUTUBE_API_KEY")
	if apiKey == "" {
		return fmt.Errorf("`YOUTUBE_API_KEY` must be defined in config.")
	}
	youtubeClient := youtube.NewYoutubeClient(apiKey)

	ctx := context.Background()
	feedURL, err := youtubeClient.GenerateAtomFeedURL(ctx, handle)

	if err != nil {
		return fmt.Errorf("Error generating Atom Feed URL for %s: %w", handle, err)
	}

	if checkFeed {
		if err = feed.ValidateAtomFeed(ctx, feedURL); err != nil {
			return fmt.Errorf("Error validating Atom Feed URL (%s): %w", feedURL, err)
		}
	}
	fmt.Println(feedURL)

	if numEntries > 0 {
		entryList, err := feed.ListEntry(ctx, feedURL, numEntries)
		if err != nil {
			return fmt.Errorf("Error listing entries for Atom Feed URL (%s): %w", feedURL, err)
		}

		for _, entry := range entryList {
			fmt.Printf("%s:%s}\n", entry.Title, entry.Link.Href)
		}
	}

	return nil
}

func runServer(portNumber int) {
	router := server.GetRouter()

	router.Run(fmt.Sprintf(":%d", portNumber))
}
