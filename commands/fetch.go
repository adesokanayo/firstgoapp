package commands

import (
	"fmt"
	"log"
	"net/url"
	"os"
	"os/signal"
	"strings"

	"time"

	rss "github.com/jteeuwen/go-pkg-rss"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

//Config struct holds a slice of feeds which is a slice of strings and integer port number
type Config struct {
	Feeds []string
	Port  int
}

var fetchCmd = &cobra.Command{
	Use:   "fetch",
	Short: "Fetch feeds",
	Long:  `Dagobah will fetch all feeds listed in the config file.`,
	Run:   fetchRun,
}

func init() {
	fetchCmd.Flags().Int("rsstimeout", 5, "Timeout (in min) for RSS retrival")
	viper.BindPFlag("rsstimeout", fetchCmd.Flags().Lookup("rsstimeout"))
}

func fetchRun(cmd *cobra.Command, args []string) {

	Fetcher()
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt)
	<-sigChan
}

//Fetcher read from the configuration file
func Fetcher() {
	var config Config
	err := viper.Unmarshal(&config)
	if err != nil {
		fmt.Println(err)
	}

	for _, feed := range config.Feeds {
		go PollFeed(feed)
	}

}

//PollFeed determines the interval
func PollFeed(uri string) {
	timeout := viper.GetInt("RSSTimeout")
	if timeout < 1 {
		timeout = 1
	}
	feed := rss.New(timeout, true, chanHandler, itemHandler)

	for {
		if err := feed.Fetch(uri, nil); err != nil {
			fmt.Fprintf(os.Stderr, "[e] %s: %s", uri, err)
			return
		}

		fmt.Printf("Sleeping for %d seconds on %s\n", feed.SecondsTillUpdate(), uri)
		time.Sleep(time.Duration(feed.SecondsTillUpdate() * 1e9))
	}

}

func chanHandler(feed *rss.Feed, newchannels []*rss.Channel) {
	fmt.Printf("%d new channel(s) in %s\n", len(newchannels), feed.Url)
	for _, ch := range newchannels {
		chnl := chnlify(ch)
		if err := Channels().Insert(chnl); err != nil {
			if !strings.Contains(err.Error(), "E11000") {
				fmt.Printf("Database error. Err:%v", err)
			}
		}
	}

}

func itemHandler(feed *rss.Feed, ch *rss.Channel, newitems []*rss.Item) {
	fmt.Printf("%d new item(s) in %s\n", len(newitems), feed.Url)
	for _, item := range newitems {
		itm := itmify(item, ch)
		if err := Items().Insert(itm); err != nil {
			if !strings.Contains(err.Error(), "E11000") {
				fmt.Printf("Database error. Err:%v", err)
			}
		}
	}

}

//Itm is all item in DB
type Itm struct {
	Date         time.Time
	Key          string
	ChannelKey   string
	Title        string
	FullContent  string
	Links        []*rss.Link
	Description  string
	Author       rss.Author
	Categories   []*rss.Category
	Comments     string
	Enclosures   []*rss.Enclosure
	GUID         *string `bson:",omitempty"`
	Source       *rss.Source
	PubDate      string
	ID           string `bson:",omitempty"`
	Generator    *rss.Generator
	Contributors []string
	Content      *rss.Content
	Extensions   map[string]map[string][]rss.Extension
}

func itmify(o *rss.Item, ch *rss.Channel) Itm {
	var x Itm
	x.Title = o.Title
	x.Links = o.Links
	x.ChannelKey = ch.Key()
	x.Description = o.Description
	x.Author = o.Author
	x.Categories = o.Categories
	x.Comments = o.Comments
	x.Enclosures = o.Enclosures
	x.GUID = o.Guid
	x.PubDate = o.PubDate
	x.ID = o.Id
	x.Key = o.Key()
	x.Generator = o.Generator
	x.Contributors = o.Contributors
	x.Content = o.Content
	x.Extensions = o.Extensions
	x.Date, _ = o.ParsedPubDate()

	if o.Content != nil && o.Content.Text != "" {
		x.FullContent = o.Content.Text
	} else {
		x.FullContent = o.Description
	}

	return x
}

//Chnl is used to hold channels
type Chnl struct {
	Key            string
	Title          string
	Links          []rss.Link
	Description    string
	Language       string
	Copyright      string
	ManagingEditor string
	WebMaster      string
	PubDate        string
	LastBuildDate  string
	Docs           string
	Categories     []*rss.Category
	Generator      rss.Generator
	TTL            int
	Rating         string
	SkipHours      []int
	SkipDays       []int
	Image          rss.Image
	ItemKeys       []string
	Cloud          rss.Cloud
	TextInput      rss.Input
	Extensions     map[string]map[string][]rss.Extension
	ID             string
	Rights         string
	Author         rss.Author
	SubTitle       rss.SubTitle
}

func chnlify(o *rss.Channel) Chnl {
	var x Chnl
	x.Key = o.Key()
	x.Title = o.Title
	x.Links = o.Links
	x.Description = o.Description
	x.Language = o.Language
	x.Copyright = o.Copyright
	x.ManagingEditor = o.ManagingEditor
	x.WebMaster = o.WebMaster
	x.PubDate = o.PubDate
	x.LastBuildDate = o.LastBuildDate
	x.Docs = o.Docs
	x.Categories = o.Categories
	x.Generator = o.Generator
	x.TTL = o.TTL
	x.Rating = o.Rating
	x.SkipHours = o.SkipHours
	x.SkipDays = o.SkipDays
	x.Image = o.Image
	x.Cloud = o.Cloud
	x.TextInput = o.TextInput
	x.Extensions = o.Extensions
	x.ID = o.Id
	x.Rights = o.Rights
	x.Author = o.Author
	x.SubTitle = o.SubTitle

	var keys []string
	for _, y := range o.Items {
		keys = append(keys, y.Key())
	}
	x.ItemKeys = keys

	return x

}

//FirstLink is used to
func (i Itm) FirstLink() (link rss.Link) {
	if len(i.Links) == 0 || i.Links[0] == nil {
		return
	}
	return *i.Links[0]
}

//WorthShowing is to display
func (i Itm) WorthShowing() bool {
	if len(i.FullContent) > 100 {
		return true
	}
	return false
}

// HomePage displays HomePage
func (c Chnl) HomePage() string {
	if len(c.Links) == 0 {
		return ""
	}

	url, err := url.Parse(c.Links[0].Href)
	if err != nil {
		log.Println(err)
	}
	return url.Scheme + "://" + url.Host
}
