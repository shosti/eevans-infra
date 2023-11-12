package blogmailer

import (
	"bufio"
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"os"
	"slices"
	"strings"
	"text/template"
	"time"

	awsConfig "github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/aws/aws-sdk-go-v2/service/s3/types"
	"github.com/mmcdole/gofeed"
	"jaytaylor.com/html2text"
)

const (
	sendDelay = 1 * time.Hour
)

var (
	logger = slog.New(slog.NewJSONHandler(os.Stdout, nil))

	processedKeyName = "processed.txt"
	htmlContentType  = "html"
)

type BlogMailer struct {
	lmURL         string
	lmUser        string
	lmPassword    string
	s3Client      *s3.Client
	stateBucket   string
	feedURL       string
	categoryLists map[string][]int

	processed []string
}

type Config struct {
	ListmonkURL      string
	ListmonkUser     string
	ListmonkPassword string
	FeedURL          string
	StateBucket      string
	S3Endpoint       string
	S3Region         string
	CategoryLists    map[string][]int
}

func New(config *Config) (*BlogMailer, error) {
	s3Cfg, err := awsConfig.LoadDefaultConfig(context.Background())
	if err != nil {
		return nil, fmt.Errorf("could not load S3 config: %w", err)
	}
	s3Client := s3.New(s3.Options{
		Credentials:  s3Cfg.Credentials,
		BaseEndpoint: &config.S3Endpoint,
		Region:       config.S3Region,
		UsePathStyle: true,
	})

	return &BlogMailer{
		lmURL:         config.ListmonkURL,
		lmUser:        config.ListmonkUser,
		lmPassword:    config.ListmonkPassword,
		feedURL:       config.FeedURL,
		s3Client:      s3Client,
		stateBucket:   config.StateBucket,
		categoryLists: config.CategoryLists,
	}, nil
}

func (m *BlogMailer) Run(ctx context.Context) error {
	if err := m.run(ctx); err != nil {
		logger.Error(fmt.Sprintf("unable to process blog feed: %v", err))
		return err
	}
	return nil
}

func (m *BlogMailer) run(ctx context.Context) error {
	fp := gofeed.NewParser()
	feed, err := fp.ParseURLWithContext(m.feedURL, ctx)
	if err != nil {
		return fmt.Errorf("could not parse feed at %s: %w", m.feedURL, err)
	}
	if err := m.refreshProcessed(ctx); err != nil {
		return err
	}
	for _, item := range feed.Items {
		logger.Info("handling feed item", "guid", item.GUID)
		if m.isProcessed(item) {
			logger.Info("feed item already processed", "guid", item.GUID)
			continue
		}
		logger.Info("creating listmonk campaign", "guid", item.GUID)
		lists := m.listsForItem(item)
		if len(lists) == 0 {
			logger.Info("feed item has no associated mailing lists", "guid", item.GUID)
			continue
		}
		campaign, err := m.createCampaign(ctx, item, lists)
		if err != nil {
			return fmt.Errorf("error creating campaign: %w", err)
		}
		logger.Info("successfully created listmonk campaign", "guid", item.GUID, "campaign_id", campaign.ID)

		if err := m.scheduleCampaign(ctx, campaign); err != nil {
			return fmt.Errorf("error scheduling campaign: %w", err)
		}
		logger.Info("successfully scheduled listmonk campaign", "guid", item.GUID, "campaign_id", campaign.ID, "sheduled_for", campaign.SendAt)

		m.processed = append(m.processed, item.GUID)
		if err := m.updateProcessed(ctx); err != nil {
			return fmt.Errorf("error updating processed campaign file (for %s): %w", item.GUID, err)
		}
	}

	return nil
}

func (m *BlogMailer) listmonkAuthRequestEditor(ctx context.Context, req *http.Request) error {
	req.SetBasicAuth(m.lmUser, m.lmPassword)
	return nil
}

func (m *BlogMailer) refreshProcessed(ctx context.Context) error {
	resp, err := m.s3Client.GetObject(ctx, &s3.GetObjectInput{
		Bucket: &m.stateBucket,
		Key:    &processedKeyName,
	})
	if err != nil {
		var notFound *types.NoSuchKey
		if errors.As(err, &notFound) {
			logger.Info("state file does not exist, creating initial version")

			_, err := m.s3Client.PutObject(ctx, &s3.PutObjectInput{
				Bucket: &m.stateBucket,
				Key:    &processedKeyName,
				Body:   strings.NewReader(""),
			})
			if err != nil {
				return fmt.Errorf("could not create initial state file: %w", err)
			}

			logger.Info("initial state created successfully")
			m.processed = []string{}
			return nil
		} else {
			return fmt.Errorf("unknown error reading state file: %w", err)
		}
	}
	defer resp.Body.Close()

	processed := []string{}
	scanner := bufio.NewScanner(resp.Body)
	for scanner.Scan() {
		processed = append(processed, scanner.Text())
	}
	if err := scanner.Err(); err != nil {
		return fmt.Errorf("error parsing state file: %w", err)
	}
	m.processed = processed

	return nil
}

func (m *BlogMailer) updateProcessed(ctx context.Context) error {
	body := strings.Join(m.processed, "\n")
	_, err := m.s3Client.PutObject(ctx, &s3.PutObjectInput{
		Bucket: &m.stateBucket,
		Key:    &processedKeyName,
		Body:   strings.NewReader(body),
	})
	if err != nil {
		return err
	}
	return nil
}

func (m *BlogMailer) isProcessed(item *gofeed.Item) bool {
	for _, guid := range m.processed {
		if item.GUID == guid {
			return true
		}
	}
	return false
}

func (m *BlogMailer) listsForItem(item *gofeed.Item) []int {
	var lists []int
	for _, cat := range item.Categories {
		if l, ok := m.categoryLists[cat]; ok {
			for _, id := range l {
				lists = append(lists, id)
			}
		}
	}
	slices.Sort(lists)
	slices.Compact(lists)

	return lists
}

func (m *BlogMailer) createCampaign(ctx context.Context, item *gofeed.Item, lists []int) (*Campaign, error) {
	sendAt := time.Now().Add(sendDelay).UTC().Format(time.RFC3339)
	body, err := emailBody(item)
	if err != nil {
		return nil, err
	}
	altbody, err := plaintextBody(item)
	if err != nil {
		return nil, err
	}
	var sendLists []List
	for _, listID := range lists {
		l := List{ID: listID}
		sendLists = append(sendLists, l)
	}
	campaign := &Campaign{
		Name:        fmt.Sprintf("Blog Campaign for %s", item.GUID),
		Subject:     item.Title,
		Tags:        item.Categories,
		ContentType: htmlContentType,
		Lists:       sendLists,
		Body:        body,
		Altbody:     altbody,
		SendAt:      sendAt,
	}
	reqBody, err := json.Marshal(&campaign)
	if err != nil {
		return nil, fmt.Errorf("error campaign marshalling body json: %w", err)
	}

	respBody, err := m.doListmonkCall(ctx, http.MethodPost, m.lmURL+"/api/campaigns", reqBody)
	if err != nil {
		return nil, fmt.Errorf("error creating campaign: %w", err)
	}
	defer respBody.Close()
	var resp CreateCampaignResponse
	if err := json.NewDecoder(respBody).Decode(&resp); err != nil {
		return nil, fmt.Errorf("error decoding campaign response body: %w", err)
	}

	return resp.Data, nil
}

func (m *BlogMailer) scheduleCampaign(ctx context.Context, campaign *Campaign) error {
	update := CampaignStatusUpdate{Status: CampaignStatusScheduled}
	reqBody, err := json.Marshal(&update)
	if err != nil {
		return fmt.Errorf("error creating campaign update request body: %w", err)
	}
	url := fmt.Sprintf("%v/api/campaigns/%v/status", m.lmURL, campaign.ID)
	if _, err := m.doListmonkCall(ctx, http.MethodPut, url, reqBody); err != nil {
		return fmt.Errorf("error updating campaign status: %w", err)
	}

	return nil
}

func (m *BlogMailer) doListmonkCall(ctx context.Context, method string, url string, body []byte) (io.ReadCloser, error) {
	req, err := http.NewRequestWithContext(ctx, method, url, bytes.NewBuffer(body))
	if err != nil {
		return nil, fmt.Errorf("error creating campaign request body: %w", err)
	}
	req.Header.Add("content-type", "application/json")
	req.Header.Add("accept", "application/json")
	req.SetBasicAuth(m.lmUser, m.lmPassword)
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("error calling listmonk api: %w", err)
	}
	if resp.StatusCode != http.StatusOK {
		message, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("listmonk api error: status %v, message %s", resp.StatusCode, message)
	}

	return resp.Body, nil
}

const bodyTemplate = `
<h1>{{.Title}}</h1>
<hr />
<a href="{{.Link}}">Read this on my blog</a>
<hr />
{{.Description}}
`

var emailBodyTmpl = template.Must(template.New("email-body").Parse(bodyTemplate))

func emailBody(item *gofeed.Item) (string, error) {
	var b strings.Builder
	if err := emailBodyTmpl.Execute(&b, item); err != nil {
		return "", fmt.Errorf("error rendering email body: %w", err)
	}

	return b.String(), nil
}

func plaintextBody(item *gofeed.Item) (string, error) {
	text, err := html2text.FromString(item.Description, html2text.Options{PrettyTables: true})
	if err != nil {
		return "", fmt.Errorf("error formatting email alt body: %w", err)
	}
	return text, nil
}
