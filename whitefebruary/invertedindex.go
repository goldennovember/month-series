package whitefebruary

import (
	"bufio"
	"encoding/json"
	"fmt"
	"math"
	"os"
	"regexp"
	"sort"
	"strings"
	"sync"
)

type Abstracts struct {
	Title    string `json:"title"`
	URL      string `json:"url"`
	Abstract string `json:"abstract"`
}

type Data struct {
	Term         string
	Frequency    int
	DocumentList []*Abstracts
}

type InvertedIndex struct {
	Filename     string
	HashMap      map[string]*Data
	numAbstracts int
}

type SearchResult struct {
	pageUrl   string
	pageTitle string
	score     float64
}

var numPartitions = 8

func (invertedIndex *InvertedIndex) buildIndex() {

	// Read the file
	file, err := os.Open(invertedIndex.Filename)
	if err != nil {
		panic(err)
	}

	defer func(file *os.File) {
		err := file.Close()
		if err != nil {

		}
	}(file)

	// A scanner is used to read text from a Reader (such as files)
	scanner := bufio.NewScanner(file)

	data := []string{}

	for scanner.Scan() {
		data = append(data, scanner.Text())
		invertedIndex.numAbstracts++
	}

	// Create WaitGroup to wait for all goroutines to finish
	var wg sync.WaitGroup
	wg.Add(numPartitions)

	var lock sync.Mutex
	for i := 0; i < numPartitions; i++ {
		go func(partition int) {
			lowerBound := partition * (len(data) / numPartitions)
			upperBound := (partition + 1) * (len(data) / numPartitions)
			if partition == numPartitions-1 {
				upperBound = len(data)
			}

			for _, v := range data[lowerBound:upperBound] {
				abstract := Abstracts{}
				json.Unmarshal(
					[]byte(v),
					&abstract,
				)

				normalized := normalizeTerm(abstract.Abstract)
				words := strings.Fields(normalized)

				for _, word := range words {
					lock.Lock()
					word = normalizeTerm(word)
					if _, ok := invertedIndex.HashMap[word]; !ok {
						invertedIndex.HashMap[word] = &Data{
							Term:         word,
							Frequency:    1,
							DocumentList: []*Abstracts{&abstract},
						}
					} else {
						invertedIndex.HashMap[word].Frequency++
						if !containsDocument(invertedIndex.HashMap[word].DocumentList, &abstract) {
							invertedIndex.HashMap[word].DocumentList = append(invertedIndex.HashMap[word].DocumentList, &abstract)
						}
					}
					lock.Unlock()
				}
			}
			wg.Done()
		}(i)
	}

	wg.Wait()

	if err := scanner.Err(); err != nil {
		panic(err)
	}

}

func (invertedIndex *InvertedIndex) getPageUrlsForTerm(terms []string) []string {
	urlHashMap := make(map[string]bool)
	for _, term := range terms {
		word := strings.ToLower(term)
		if _, ok := invertedIndex.HashMap[word]; ok {
			for _, document := range invertedIndex.HashMap[word].DocumentList {
				urlHashMap[document.URL] = true
			}
		}
	}

	urlList := []string{}

	for url := range urlHashMap {
		urlList = append(urlList, url)
	}
	return urlList
}

func (invertedIndex *InvertedIndex) getSearchResult(terms []string) []SearchResult {
	urlHashMap := make(map[string]*SearchResult)
	for _, term := range terms {
		word := strings.ToLower(term)
		if _, ok := invertedIndex.HashMap[word]; ok {
			for _, document := range invertedIndex.HashMap[word].DocumentList {
				// Count words
				normalized := normalizeTerm(document.Abstract)
				words := strings.Fields(normalized)
				termCount := 0
				for _, s := range words {
					if s == word {
						termCount++
					}
				}

				if _, ok := urlHashMap[document.URL]; !ok {
					urlHashMap[document.URL] = &SearchResult{
						pageUrl:   document.URL,
						pageTitle: document.Title,
						score:     0,
					}
				}
				tf := float64(termCount) / float64(len(words))
				idf := math.Log(float64(invertedIndex.numAbstracts) / float64(len(invertedIndex.HashMap[word].DocumentList)))
				urlHashMap[document.URL].score += math.Floor(tf*idf*100000) / 100000
			}
		}
	}

	searchResults := []SearchResult{}

	for _, result := range urlHashMap {
		searchResults = append(searchResults, *result)
	}

	sort.Slice(searchResults, func(i, j int) bool {
		return searchResults[i].score >= searchResults[j].score
	})

	res := []SearchResult{}

	for i := 0; i < len(searchResults); i++ {
		if i == 5 {
			break
		}
		res = append(res, searchResults[i])
	}

	return res
}

func containsDocument(documentList []*Abstracts, document *Abstracts) bool {
	for _, doc := range documentList {
		if doc == document {
			return true
		}
	}
	return false
}

func normalizeTerm(term string) string {
	return regexp.MustCompile(`[^a-zA-Z0-9 ]+`).ReplaceAllString(strings.ToLower(term), " ")
}

func (s *SearchResult) String() string {
	return fmt.Sprintf("SearchResult{PageUrl='%s',pageTiTle='%s',Score='%v'", s.pageTitle, s.pageUrl, s.score)
}
