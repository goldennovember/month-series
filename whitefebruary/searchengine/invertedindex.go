package searchengine

import (
	"bufio"
	"encoding/json"
	"os"
	"strings"
)

var stopwords = []string{"a", "able", "about",
	"across", "after", "all", "almost", "also", "am", "among", "an",
	"and", "any", "are", "as", "at", "be", "because", "been", "but",
	"by", "can", "cannot", "could", "dear", "did", "do", "does",
	"either", "else", "ever", "every", "for", "from", "get", "got",
	"had", "has", "have", "he", "her", "hers", "him", "his", "how",
	"however", "i", "if", "in", "into", "is", "it", "its", "just",
	"least", "let", "like", "likely", "may", "me", "might", "most",
	"must", "my", "neither", "no", "nor", "not", "of", "off", "often",
	"on", "only", "or", "other", "our", "own", "rather", "said", "say",
	"says", "she", "should", "since", "so", "some", "than", "that",
	"the", "their", "them", "then", "there", "these", "they", "this",
	"tis", "to", "too", "twas", "us", "wants", "was", "we", "were",
	"what", "when", "where", "which", "while", "who", "whom", "why",
	"will", "with", "would", "yet", "you", "your"}

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

var InvertedIndex = map[string]*Data{}

func buildIndex(filename string) {

	// Convert stopwords to a map for faster lookup
	stopwordsMap := make(map[string]bool)
	for _, word := range stopwords {
		stopwordsMap[word] = true
	}

	abstractList := []Abstracts{}

	// Read the file
	file, err := os.Open(filename)
	if err != nil {
		panic(err)
	}

	defer file.Close()

	// A scanner is used to read text from a Reader (such as files)
	scanner := bufio.NewScanner(file)

	// optionally, resize scanner's capacity for lines over 64K, see next example
	for scanner.Scan() {
		value, _ := os.ReadFile(scanner.Text())
		var abstract Abstracts
		json.Unmarshal(value, &abstract)
		abstractList = append(abstractList, abstract)
	}

	if err := scanner.Err(); err != nil {
		panic(err)
	}

	for _, abstract := range abstractList {
		words := strings.Fields(abstract.Abstract)
		for _, word := range words {
			word = strings.ToLower(word)
			if !stopwordsMap[word] {
				if InvertedIndex[word] == nil {
					InvertedIndex[word] = &Data{Term: word, Frequency: 1, DocumentList: []*Abstracts{&abstract}}
				} else {
					InvertedIndex[word].Frequency++
					if !containsDocument(InvertedIndex[word].DocumentList, abstract) {
						InvertedIndex[word].DocumentList = append(InvertedIndex[word].DocumentList, &abstract)
					}
				}
			}
		}
	}

}

func getPageUrlsForTerm(terms []string) []string {
	urlList := []string{}
	for _, term := range terms {
		word := strings.ToLower(term)
		if InvertedIndex[word] != nil {
			for _, document := range InvertedIndex[word].DocumentList {
				urlList = append(urlList, document.URL)
			}
		}
	}
	return urlList
}

func containsDocument(documentList []*Abstracts, document Abstracts) bool {
	for _, doc := range documentList {
		if doc.URL == document.URL {
			return true
		}
	}
	return false
}
