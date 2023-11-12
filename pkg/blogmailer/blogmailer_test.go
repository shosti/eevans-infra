package blogmailer

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"regexp"
	"slices"
	"strconv"
	"strings"
	"testing"

	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/johannesboyne/gofakes3"
	"github.com/johannesboyne/gofakes3/backend/s3mem"
)

const (
	testBucket      = "blog-mailer"
	technicalListID = 2
	miscListID      = 3
)

func TestInitialRun(t *testing.T) {
	ctx := context.Background()
	feedServer := httptest.NewServer(http.HandlerFunc(fakeFeed))
	defer feedServer.Close()

	lm := fakeListmonk{t: t, scheduled: map[int]bool{}}
	lmServer := httptest.NewServer(&lm)
	defer lmServer.Close()

	fakeS3 := gofakes3.New(s3mem.New())
	s3Server := httptest.NewServer(fakeS3.Server())
	defer s3Server.Close()
	s3Client := s3.New(s3.Options{
		BaseEndpoint: &s3Server.URL,
		Region:       "us-east-1",
		UsePathStyle: true,
	})
	prepS3Bucket(ctx, s3Client)

	m, err := New(&Config{
		FeedURL:          feedServer.URL,
		ListmonkURL:      lmServer.URL,
		ListmonkUser:     "listmonk",
		ListmonkPassword: "abc123",
		CategoryLists: map[string][]int{
			"technical":    []int{technicalListID},
			"book reviews": []int{miscListID},
		},
		StateBucket: testBucket,
	})
	if err != nil {
		t.Fatal(err)
	}
	m.s3Client = s3Client
	if err := m.Run(ctx); err != nil {
		t.Fatal(err)
	}
	if err := m.refreshProcessed(ctx); err != nil {
		t.Fatal(err)
	}

	wantProcessed := []string{
		"https://eevans.co/blog/theoretical-minimum-review/",
		"https://eevans.co/blog/bright-ages-review/",
		"https://eevans.co/blog/garage/",
		"https://eevans.co/blog/wraft/",
		"https://eevans.co/blog/kubernetes-multi-node/",
		"https://eevans.co/blog/deconstructing-kubernetes-networking/",
		"https://eevans.co/blog/minimum-viable-kubernetes/",
	}
	for _, want := range wantProcessed {
		if slices.Index(m.processed, want) == -1 {
			t.Errorf("%s should have been added to the processed list but wasn't", want)
		}
	}
}

func TestSubsequentRun(t *testing.T) {
	ctx := context.Background()
	feedServer := httptest.NewServer(http.HandlerFunc(fakeFeed))
	defer feedServer.Close()

	lm := fakeListmonk{t: t, scheduled: map[int]bool{}}
	lmServer := httptest.NewServer(&lm)
	defer lmServer.Close()

	fakeS3 := gofakes3.New(s3mem.New())
	s3Server := httptest.NewServer(fakeS3.Server())
	defer s3Server.Close()
	s3Client := s3.New(s3.Options{
		BaseEndpoint: &s3Server.URL,
		Region:       "us-east-1",
		UsePathStyle: true,
	})
	prepS3StateFile(ctx, s3Client)

	m, err := New(&Config{
		FeedURL:          feedServer.URL,
		ListmonkURL:      lmServer.URL,
		ListmonkUser:     "listmonk",
		ListmonkPassword: "abc123",
		CategoryLists: map[string][]int{
			"technical":    {technicalListID},
			"book reviews": {miscListID},
		},
		StateBucket: testBucket,
	})
	if err != nil {
		t.Fatal(err)
	}
	m.s3Client = s3Client
	if err := m.Run(ctx); err != nil {
		t.Fatal(err)
	}

	type blogData struct {
		title string
		lists []int
	}
	wantProcessed := []struct {
		url        string
		title      string
		lists      []int
		bodyphrase string
	}{
		{
			url:        "https://eevans.co/blog/bright-ages-review/",
			title:      "Book Review: The Bright Ages",
			lists:      []int{miscListID},
			bodyphrase: "big step backwards in terms of math",
		},
		{
			url:        "https://eevans.co/blog/wraft/",
			title:      "Implementing Raft for Browsers with Rust and WebRTC",
			lists:      []int{technicalListID},
			bodyphrase: "stubbornly persistent",
		},
	}

	if err := m.refreshProcessed(ctx); err != nil {
		t.Fatal(err)
	}
	if len(lm.reqs) != len(wantProcessed) {
		t.Fatalf("unexpected number of listmonk requests: want %v, got %v", len(wantProcessed), len(lm.reqs))
	}
	for i, blog := range wantProcessed {
		if slices.Index(m.processed, blog.url) == -1 {
			t.Errorf("%s should have been added to the processed list but wasn't", blog.url)
		}
		lmReq := lm.reqs[i]
		if lmReq.Subject != blog.title {
			t.Errorf("wrong campaign title: want '%v', got '%v'", blog.title, lmReq.Subject)
		}
		if !strings.Contains(lmReq.Body, blog.bodyphrase) {
			t.Errorf("expected %v to contain the phrase '%s' but it didn't", blog.title, blog.bodyphrase)
		}
		if !strings.Contains(lmReq.Altbody, blog.bodyphrase) {
			t.Errorf("expected altbody for %v to contain the phrase '%s' but it didn't", blog.title, blog.bodyphrase)
		}
		if len(lmReq.Lists) != len(blog.lists) {
			t.Errorf("wrong number of campaign lists: want %v, got %v", blog.lists, lmReq.Lists)
		} else {
			for i := range lmReq.Lists {
				if lmReq.Lists[i].ID != blog.lists[i] {
					t.Errorf("wrong list requested: want %v, got %v", lmReq.Lists[i].ID, blog.lists[i])
				}
			}
		}
		if !lm.scheduled[lmReq.ID] {
			t.Errorf("campaign %v should have been scheduled but wasn't", lmReq.ID)
		}
	}
}

func fakeFeed(w http.ResponseWriter, r *http.Request) {
	content, err := os.ReadFile("fixtures/index.xml")
	if err != nil {
		panic(err)
	}
	w.Header().Add("content-type", "application/rss+xml")
	if _, err := w.Write(content); err != nil {
		panic(err)
	}
}

type fakeListmonk struct {
	t         *testing.T
	reqs      []*Campaign
	scheduled map[int]bool
}

var campaignStatusRe = regexp.MustCompile("/api/campaigns/([0-9]+)/status")

func (lm *fakeListmonk) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	u, p, ok := r.BasicAuth()
	if !ok {
		lm.t.Error("listmonk auth not provided")
	} else if u != "listmonk" || p != "abc123" {
		lm.t.Error("incorrect auth provided for listmonk")
	}
	if r.Method == http.MethodPost && r.URL.Path == "/api/campaigns" {
		req := &Campaign{}
		dec := json.NewDecoder(r.Body)
		if err := dec.Decode(&req); err != nil {
			lm.t.Fatal(err)
		}
		req.ID = len(lm.reqs) + 1
		lm.reqs = append(lm.reqs, req)
		respBody := fmt.Sprintf(`{"data":{"id":%v,"status":"draft","lists":[{"name":"technical","id":2}]}}`, req.ID)

		w.Write([]byte(respBody))
	} else if r.Method == http.MethodPut && campaignStatusRe.MatchString(r.URL.Path) {
		idStr := campaignStatusRe.FindStringSubmatch(r.URL.Path)[1]
		id, err := strconv.Atoi(idStr)
		if err != nil {
			lm.t.Fatal(err)
		}
		var req CampaignStatusUpdate
		dec := json.NewDecoder(r.Body)
		if err := dec.Decode(&req); err != nil {
			lm.t.Fatal(err)
		}
		if req.Status == CampaignStatusScheduled {
			lm.scheduled[id] = true
		} else {
			lm.t.Errorf("unexpected campaign status update: %v", req.Status)
		}
		resp := CreateCampaignResponse{Data: &Campaign{}}
		respBody, err := json.Marshal(&resp)
		if err != nil {
			lm.t.Fatal(err)
		}

		w.Write(respBody)
	} else {
		lm.t.Errorf("unknown API path: %v", r.URL)
		w.WriteHeader(http.StatusNotFound)
	}
}

func prepS3Bucket(ctx context.Context, c *s3.Client) {
	_, err := c.CreateBucket(ctx, &s3.CreateBucketInput{
		Bucket: aws.String(testBucket),
	})
	if err != nil {
		panic(err)
	}
}

func prepS3StateFile(ctx context.Context, c *s3.Client) {
	prepS3Bucket(ctx, c)
	_, err := c.PutObject(ctx, &s3.PutObjectInput{
		Bucket: aws.String(testBucket),
		Key:    aws.String(processedKeyName),
		Body:   strings.NewReader(processed),
	})
	if err != nil {
		panic(err)
	}

}

const processed = `https://eevans.co/blog/theoretical-minimum-review/
https://eevans.co/blog/garage/
https://eevans.co/blog/kubernetes-multi-node/
https://eevans.co/blog/deconstructing-kubernetes-networking/
https://eevans.co/blog/minimum-viable-kubernetes/`
