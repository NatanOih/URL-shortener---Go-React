package helpers

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"strings"

	"github.com/go-redis/redis/v8"
)

type Data struct {
	URL    string `json:"url"`
	Clicks int    `json:"clicks"`
}

func EnforceHTTP(url string) string {
	if url[:4] != "http" {
		return "http://" + url
	}

	return url
}

func RemoveDomainError(url string) bool {
	if url == os.Getenv("DOMAIN") {
		return false
	}

	newURL := strings.Replace(url, "http://", "", 1)
	newURL = strings.Replace(newURL, "https://", "", 1)
	newURL = strings.Replace(newURL, "www", "", 1)
	newURL = strings.Split(newURL, "/")[0]

	return newURL != os.Getenv("DOMAIN")

}

func FetchDataFromRedis(ctx context.Context, r *redis.Client) (map[string]Data, error) {

	result := make(map[string]Data)

	// Iterate over all keys in Redis
	iter := r.Scan(ctx, 0, "*", 0).Iterator()
	for iter.Next(ctx) {
		key := iter.Val()

		// Get the JSON data for the current key
		jsonData, err := r.Get(ctx, key).Result()
		if err != nil {
			return nil, fmt.Errorf("error getting value for key %s: %v", key, err)
		}

		// Unmarshal the JSON data into a Data struct
		var data Data
		if err := json.Unmarshal([]byte(jsonData), &data); err != nil {
			continue // Skip if not JSON or there's an unmarshaling error
		}

		// Add the key and data to the result map
		result[key] = data
	}

	if err := iter.Err(); err != nil {
		return nil, fmt.Errorf("error during iteration: %v", err)
	}

	return result, nil
}
