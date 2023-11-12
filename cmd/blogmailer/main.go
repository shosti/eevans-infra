package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"

	"gitea.eevans.me/shosti/eevans-infra/pkg/blogmailer"
)

func main() {
	listmonkUrl := mustGetenv("LISTMONK_URL")
	listmonkUser := mustGetenv("LISTMONK_USER")
	listmonkPassord := mustGetenv("LISTMONK_PASSWORD")
	feedURL := mustGetenv("FEED_URL")
	stateBucket := mustGetenv("BUCKET_NAME")
	s3Host := mustGetenv("BUCKET_HOST")
	s3Region := mustGetenv("BUCKET_REGION")
	s3Endpoint := fmt.Sprintf("http://%s", s3Host)
	categoryLists, err := parseCategoryLists(mustGetenv("CATEGORY_LISTS"))
	if err != nil {
		log.Fatalf("error parsing category lists: %v", err)
	}

	mailer, err := blogmailer.New(&blogmailer.Config{
		ListmonkURL:      listmonkUrl,
		ListmonkUser:     listmonkUser,
		ListmonkPassword: listmonkPassord,
		FeedURL:          feedURL,
		StateBucket:      stateBucket,
		S3Endpoint:       s3Endpoint,
		S3Region:         s3Region,
		CategoryLists:    categoryLists,
	})

	if err != nil {
		log.Fatalf("error initializing mailer: %v", err)
	}

	ctx := context.Background()
	if err := mailer.Run(ctx); err != nil {
		log.Fatalf("error running mailer: %v", err)
	}
}

func mustGetenv(key string) string {
	val := os.Getenv(key)
	if val == "" {
		log.Fatalf("environment variable %s must be set", key)
	}

	return val
}

func parseCategoryLists(val string) (map[string][]int, error) {
	ret := map[string][]int{}
	entries := strings.Split(val, ",")
	for _, entry := range entries {
		parts := strings.Split(entry, "=")
		if len(parts) != 2 {
			return nil, fmt.Errorf("category entry %s could not be split into k/v pair", entry)
		}
		k := parts[0]
		v, err := strconv.Atoi(parts[1])
		if err != nil {
			return nil, fmt.Errorf("category entry %s does not have integer value", entry)
		}
		ret[k] = append(ret[k], v)
	}

	return ret, nil
}
