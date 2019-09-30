package dkservice

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"regexp"
	"sort"
	"sync"
	"time"
)

type dockerService struct {
	username string
	password string
	host     string
	noAuth   bool
	verbose  bool
}

type Tag struct {
	Name        string
	CreatedTime time.Time
	Digest      string
}

func New(username, password, host string, verbose bool) dockerService {
	dks := dockerService{username: username, password: password, host: host, verbose: verbose}
	if dks.username == "" {
		dks.noAuth = true
	}
	return dks
}

func (dks dockerService) setAuth(req **http.Request) {
	if !dks.noAuth {
		(*req).SetBasicAuth(dks.username, dks.password)
	}
}

func (dks dockerService) getCreatedTime(repo, tag string) time.Time {
	url := fmt.Sprintf("%s/v2/%s/manifests/%s", dks.host, repo, tag)
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		panic(err)
	}
	dks.setAuth(&req)

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		panic(err)
	}

	type manifestResponse struct {
		History []struct {
			V1Compatibility string `json:"v1Compatibility"`
		} `json:"history"`
	}
	var mr manifestResponse
	err = json.Unmarshal(body, &mr)
	if err != nil {
		panic(err)
	}

	// Find latest time among all tags
	type V1Compatibility struct {
		Created string `json:"created"`
	}
	var latestTime string
	for _, v := range mr.History {
		var tmp V1Compatibility
		err = json.Unmarshal([]byte(v.V1Compatibility), &tmp)
		if err != nil {
			panic(err)
		}
		if tmp.Created > latestTime {
			latestTime = tmp.Created
		}
	}

	// Trim nanosecond, from "2019-09-13T09:35:06.944396815Z" to "2019-09-13T09:35:06"
	re := regexp.MustCompile(`\..*`)
	createdTimeStr := re.ReplaceAllString(latestTime, "")
	createdTime, _ := time.Parse("2006-01-02T15:04:05", createdTimeStr)

	return createdTime
}

func (dks dockerService) getCreatedTimeBatch(repo string, tags []Tag) []Tag {
	ch := make(chan Tag)
	var wg sync.WaitGroup
	wg.Add(len(tags))
	for _, t := range tags {
		go func(tag Tag) {
			defer wg.Done()
			createdTime := dks.getCreatedTime(repo, tag.Name)
			ch <- Tag{Name: tag.Name, CreatedTime: createdTime, Digest: tag.Digest}
		}(t)
	}

	go func() {
		wg.Wait()
		close(ch)
	}()

	var tagsWithCreatedTime []Tag
	for tag := range ch {
		tagsWithCreatedTime = append(tagsWithCreatedTime, tag)
	}

	return tagsWithCreatedTime
}

func (dks dockerService) getDigest(repo, tag string) string {
	url := fmt.Sprintf("%s/v2/%s/manifests/%s", dks.host, repo, tag)
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		panic(err)
	}
	dks.setAuth(&req)
	req.Header.Set("Accept", "application/vnd.docker.distribution.manifest.v2+json")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	return resp.Header.Get("Docker-Content-Digest")
}

func (dks dockerService) getDigestBatch(repo string, tags []Tag) []Tag {
	ch := make(chan Tag)
	var wg sync.WaitGroup
	wg.Add(len(tags))
	for _, t := range tags {
		go func(tag Tag) {
			defer wg.Done()
			digest := dks.getDigest(repo, tag.Name)
			ch <- Tag{Name: tag.Name, CreatedTime: tag.CreatedTime, Digest: digest}
		}(t)
	}

	go func() {
		wg.Wait()
		close(ch)
	}()

	var tagsWithDigest []Tag
	for tag := range ch {
		tagsWithDigest = append(tagsWithDigest, tag)
	}

	return tagsWithDigest
}

func (dks dockerService) deleteImageByDigest(repo, digest string) bool {
	url := fmt.Sprintf("%s/v2/%s/manifests/%s", dks.host, repo, digest)
	req, err := http.NewRequest("DELETE", url, nil)
	if err != nil {
		panic(err)
	}
	dks.setAuth(&req)

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	if dks.verbose {
		fmt.Printf("Digest: %s\n", digest)
	}

	if resp.StatusCode == 202 {
		if dks.verbose {
			fmt.Printf("Status: Deleted\n\n")
		}
		return true
	}

	if dks.verbose {
		fmt.Printf("Status: Error\n\n")
	}
	return false
}

func (dks dockerService) deleteImageByDigestBatch(repo string, tags []Tag) int {
	ch := make(chan bool)
	var wg sync.WaitGroup
	wg.Add(len(tags))
	for _, t := range tags {
		go func(tag Tag) {
			defer wg.Done()
			ch <- dks.deleteImageByDigest(repo, tag.Digest)
		}(t)
	}

	go func() {
		wg.Wait()
		close(ch)
	}()

	var nSuccess int
	for ok := range ch {
		if ok {
			nSuccess++
		}
	}

	return nSuccess
}

func (dks dockerService) getAllTag(repo string) []Tag {
	type allTagResponse struct {
		Tags []string `json:"tags"`
	}

	url := fmt.Sprintf("%s/v2/%s/tags/list", dks.host, repo)
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		panic(err)
	}
	dks.setAuth(&req)

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)

	var atr allTagResponse
	err = json.Unmarshal(body, &atr)
	if err != nil {
		panic(err)
	}

	var tags []Tag
	for _, tag := range atr.Tags {
		tags = append(tags, Tag{Name: tag})
	}

	return tags
}

func (dks dockerService) SweepOutdatedImages(repo string, thresholdAge, keepTag int) {
	fmt.Printf("Repo: %s\n", repo)

	// Get tag names and their created time
	tags := dks.getAllTag(repo)
	tagsWithCreatedTime := dks.getCreatedTimeBatch(repo, tags)

	fmt.Printf("%d tag(s)\n", len(tagsWithCreatedTime))

	// Check created time against minDateToKeep
	minDateToKeep := time.Now().UTC().AddDate(0, 0, -thresholdAge)
	var tagsToDelete []Tag
	for _, tag := range tagsWithCreatedTime {
		if tag.CreatedTime.Before(minDateToKeep) {
			tagsToDelete = append(tagsToDelete, tag)
		}
	}

	fmt.Printf("%d obsolete tag(s) – created before %v\n", len(tagsToDelete), minDateToKeep)

	sort.Slice(tagsToDelete, func(i, j int) bool {
		return tagsToDelete[i].CreatedTime.Before(tagsToDelete[j].CreatedTime)
	})

	maxNumberToDelete := len(tags) - keepTag
	endIndex := maxNumberToDelete
	if maxNumberToDelete > len(tagsToDelete) {
		endIndex = len(tagsToDelete)
	} else if maxNumberToDelete < 0 {
		endIndex = 0
	}

	tagsToDelete = tagsToDelete[0:endIndex]

	fmt.Printf("%d tag(s) to delete\n\n", len(tagsToDelete))

	tagsToDelete = dks.getDigestBatch(repo, tagsToDelete)

	nDeleted := dks.deleteImageByDigestBatch(repo, tagsToDelete)
	fmt.Printf("%d image(s) have been successfully deleted", nDeleted)
}
