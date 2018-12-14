package main

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"path"
	"sort"
	"sync"
)

var ALLOW_EXT = []string{".jpg", ".png"}
var SKIP_DIR = []string{"tmp"}

type Library struct {
	BaseDir      string `json:"base_dir,"`
	sync.RWMutex `json:"-"`
	Books        map[string]Book `json:"books,"`
}

type Book struct {
	Chapters map[string]Pages `json:"chapters,"`
}

type Pages []string

// Contains tells whether a contains x.
func Contains(a []string, x string) bool {
	for _, n := range a {
		if x == n {
			return true
		}
	}
	return false
}

func (b *Book) Clone() *Book {
	newBook := Book{}
	newBook.Chapters = make(map[string]Pages)
	for k, v := range b.Chapters {
		newBook.Chapters[k] = v.Clone()
	}
	return &newBook
}

func (p *Pages) Clone() []string {
	newPages := make([]string, len(*p))
	copy(newPages, *p)
	return newPages
}

func LoadPages(baseDir string, book string, chapter string) Pages {
	fullPath := baseDir + "/" + book + "/" + chapter
	pageEntries, err := ioutil.ReadDir(fullPath)
	if err != nil {
		log.Fatal(err)
		return nil
	}
	pages := make([]string, 0)
	for _, file := range pageEntries {
		if !file.IsDir() {
			fullFilename := fullPath + "/" + file.Name()
			ext := path.Ext(fullFilename)

			if Contains(ALLOW_EXT, ext) {
				pages = append(pages, book+"/"+chapter+"/"+file.Name())
				//log.Printf("Add page[%d]:%s\n", idx, fullFilename)
			}
		}
	}
	return pages
}

func NewBook(baseDir string, dir string) *Book {
	fullPath := baseDir + "/" + dir
	chapterEntries, err := ioutil.ReadDir(fullPath)
	if err != nil {
		log.Fatal(err)
		return nil
	}
	chapters := make(map[string]Pages, len(chapterEntries))
	for _, file := range chapterEntries {
		if file.IsDir() {
			if !Contains(SKIP_DIR, file.Name()) {
				chapters[file.Name()] = LoadPages(baseDir, dir, file.Name())
			}
		}
	}
	return &Book{Chapters: chapters}
}

func NewLibrary(baseDir string) *Library {
	l := make(map[string]Book)
	c := Library{
		BaseDir: baseDir,
		Books:   l,
	}
	return &c
}

// return a device state for an account
func (l *Library) LoadAll() bool {
	ok := false
	l.RLock()
	defer l.RUnlock()
	files, err := ioutil.ReadDir(l.BaseDir)
	if err != nil {
		log.Fatal(err)
		return false
	}

	for _, file := range files {
		if file.IsDir() {
			l.Books[file.Name()] = *NewBook(l.BaseDir, file.Name())
		}
	}
	ok = true
	return ok
}

func (l *Library) Delete(book string) {
	l.Lock()
	delete(l.Books, book)
	l.Unlock()
}

func (l *Library) Store(bookDir string) bool {
	ok := false
	l.Lock()

	defer l.Unlock()
	return ok
}

func (l *Library) Load(book string) (Book, bool) {
	ok := false
	l.Lock()
	defer l.Unlock()
	result, ok := l.Books[book]
	if ok {
		log.Println("loadresult:", result)
		b := result.Clone()
		return *b, ok
	}
	return Book{}, false
}

func (l *Library) GetBooks() []string {
	l.Lock()
	defer l.Unlock()
	returnLib := make([]string, len(l.Books))

	i := 0
	for k, _ := range l.Books {
		returnLib[i] = k
		i = i + 1
	}
	sort.Strings(returnLib)
	return returnLib
}

func (l *Library) GetChapters(book string) []string {
	l.Lock()
	defer l.Unlock()
	result, ok := l.Books[book]
	if ok {
		returnChapters := make([]string, len(result.Chapters))
		i := 0
		for k, _ := range result.Chapters {
			returnChapters[i] = k
			i = i + 1
		}
		sort.Strings(returnChapters)
		return returnChapters
	}
	return nil
}

func (l *Library) GetPages(book string, chapter string) []string {
	l.Lock()
	defer l.Unlock()
	foundBook, ok := l.Books[book]
	if ok {
		foundChapter, ok := foundBook.Chapters[chapter]
		if ok {
			returnPages := make([]string, len(foundChapter))
			for k, v := range foundChapter {
				returnPages[k] = v
			}
			sort.Strings(returnPages)
			return returnPages
		}
		return nil
	}
	return nil
}

func (l *Library) Dump() string {
	b, err := json.Marshal(*l)
	if err != nil {
		log.Println("Marshalling error:", err)
	}
	return string(b)
}
