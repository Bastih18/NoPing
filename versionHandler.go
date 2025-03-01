package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/fatih/color"
)

type VersionCache struct {
	LatestVersion string    `json:"latest_version"`
	Timestamp     time.Time `json:"timestamp"`
}

func getCacheFilePath() string {
	return filepath.Join(os.TempDir(), "noping-version-cache.json")
}

func getLatestVersion(current string) (string, string, string) {
	if cached, err := loadCache(); err == nil {
		return compareVersionsAndReturn(cached.LatestVersion, current)
	}

	client := http.Client{Timeout: 5 * time.Second}
	resp, err := client.Get("https://api.github.com/repos/Bastih18/NoPing/releases/latest")
	if err != nil || resp.StatusCode != http.StatusOK {
		return "Could not fetch", color.YellowString("are maybe"), "the latest"
	}
	defer resp.Body.Close()

	var release struct {
		TagName string `json:"tag_name"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&release); err != nil || release.TagName == "" {
		return "Could not decode", color.YellowString("are maybe"), "the latest"
	}
	saveCache(release.TagName)

	return compareVersionsAndReturn(release.TagName, current)
}

func compareVersionsAndReturn(latest, current string) (string, string, string) {
	if current == latest {
		return latest, color.GreenString("are"), "the latest"
	}
	if current == "dev" {
		return latest, color.GreenString("are"), color.CyanString("the dev")
	}

	const pattern = `^v\d+\.\d+\.\d+$`
	re := regexp.MustCompile(pattern)
	if re.MatchString(current) && re.MatchString(latest) {
		if cmp := compareVersions(current, latest); cmp < 0 {
			return latest, color.RedString("are"), color.RedString("an outdated")
		} else if cmp > 0 {
			return latest, color.GreenString("are"), color.CyanString("an ahead")
		}
	}
	return latest, color.YellowString("are maybe"), "the latest"
}

func loadCache() (*VersionCache, error) {
	cacheFile := getCacheFilePath()
	data, err := os.ReadFile(cacheFile)
	if err != nil {
		return nil, err
	}

	var cache VersionCache
	if err := json.Unmarshal(data, &cache); err != nil {
		return nil, err
	}

	if time.Since(cache.Timestamp) > 10*time.Minute {
		return nil, fmt.Errorf("cache expired")
	}
	return &cache, nil
}

func saveCache(version string) {
	cacheFile := getCacheFilePath()
	cache := VersionCache{
		LatestVersion: version,
		Timestamp:     time.Now(),
	}
	data, _ := json.Marshal(cache)
	_ = os.WriteFile(cacheFile, data, 0644)
}

// 📌 Version comparison function
func compareVersions(v1, v2 string) int {
	p1, p2 := strings.Split(v1[1:], "."), strings.Split(v2[1:], ".")
	for i := 0; i < 3; i++ {
		n1, _ := strconv.Atoi(p1[i])
		n2, _ := strconv.Atoi(p2[i])
		if n1 < n2 {
			return -1
		} else if n1 > n2 {
			return 1
		}
	}
	return 0
}
