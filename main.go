package main

import (
	"fmt"
	"net/url"
	"os"
	"time"

	"odysee-sitemap-generator/chainquery"
	"odysee-sitemap-generator/configs"

	"github.com/ikeikeikeike/go-sitemap-generator/stm"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

var (
	playerEndpoint  string
	embedEndpoint   string
	website         string
	sitemapEndpoint string
	sitemapType     string
)
var cqApi *chainquery.CQApi

func main() {
	cmd := &cobra.Command{
		Use:   "sitemap-generator",
		Short: "builds and uploads a sitemap for odysee.com",
		Run:   crawler,
		Args:  cobra.RangeArgs(0, 0),
	}

	cmd.Flags().StringVar(&playerEndpoint, "player-endpoint", "https://player.odycdn.com/", "endpoint of the player")
	cmd.Flags().StringVar(&embedEndpoint, "embed-endpoint", "https://odysee.com/$/embed/", "endpoint for embeds")
	cmd.Flags().StringVar(&website, "website", "https://odysee.com", "endpoint for embeds")
	cmd.Flags().StringVar(&sitemapEndpoint, "sitemap-endpoint", "https://sitemaps.odysee.com", "endpoint for embeds")
	cmd.Flags().StringVar(&sitemapType, "type", "global", "type of sitemap to generate (global, monthly, weekly, daily, hourly")

	if err := cmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
func crawler(cmd *cobra.Command, args []string) {
	err := configs.Init("./config.json")
	if err != nil {
		panic(err)
	}
	cqApi, err = chainquery.Init()
	if err != nil {
		panic(err)
	}

	prevId := uint64(0)
	newPrevId := uint64(0)
	sm := stm.NewSitemap()
	sm.SetSitemapsPath(sitemapType)
	sm.SetDefaultHost(website)
	sm.SetSitemapsHost(sitemapEndpoint)
	sm.SetCompress(true)
	sm.Create()

	topId, err := cqApi.GetTopID()
	if err != nil {
		panic(err)
	}

	since := time.Time{}

	switch sitemapType {
	case "global":
		newPrevId = 0
	case "monthly":
		since = time.Now().AddDate(0, -1, 0)
	case "weekly":
		since = time.Now().AddDate(0, 0, -7)
	case "daily":
		since = time.Now().AddDate(0, 0, -1)
	case "hourly":
		since = time.Now().Add(-1 * time.Hour)
	default:
		panic("invalid sitemap type")
	}

	newPrevId, err = cqApi.GetIdFor(since)
	if err != nil {
		panic(err)
	}
	//newPrevId = 24750000

	for prevId < topId {
		if prevId == newPrevId {
			newPrevId++
		}
		prevId = newPrevId
		var claims []*chainquery.Claim
		claims, newPrevId, err = cqApi.GetVideoStreamsBatch(prevId, 50000, since)
		if err != nil {
			panic(err)
		}
		processClaims(sm, claims)
	}
	sm.Finalize().PingSearchEngines()
}

func processClaims(sm *stm.Sitemap, claims []*chainquery.Claim) {
	for _, c := range claims {
		_, err := url.ParseRequestURI(c.ThumbnailURL)
		if err != nil {
			logrus.Errorf("invalid thumbnail URL found: %s", c.ThumbnailURL)
			continue
		}
		description := parseDescription(c.Description.String)
		releaseTime := time.Unix(c.TransactionTime, 0)
		if !c.ReleaseTime.IsNull() {
			releaseTime = time.Unix(c.ReleaseTime.Int64, 0)
		}
		if len(c.SdHash) < 10 {
			continue
		}
		videoURL := stm.URL{
			"thumbnail_loc":    c.ThumbnailURL,
			"title":            c.Title,
			"description":      description,
			"publication_date": releaseTime.Format("2006-01-02T15:04:05-07:00"),
			"player_loc":       stm.Attrs{fmt.Sprintf("%s%s/%s", embedEndpoint, c.Name, c.ClaimID), map[string]string{"allow_embed": "Yes", "autoplay": "autoplay=1"}},
			"content_loc":      fmt.Sprintf("%sapi/v3/streams/free/%s/%s/%s", playerEndpoint, c.Name, c.ClaimID, c.SdHash[:6]),
		}
		if !c.Duration.IsNull() && c.Duration.Int >= 1 && c.Duration.Int < 28800 {
			videoURL["duration"] = c.Duration.Int
		}
		urlToAdd := stm.URL{
			"loc":   fmt.Sprintf("/%s/%s", url.QueryEscape(c.Name), c.ClaimID),
			"video": videoURL,
		}
		sm.Add(urlToAdd)
	}
}

func parseDescription(description string) string {
	if len(description) > 1000 {
		return description[:1000] + "..."
	}
	return description
}
