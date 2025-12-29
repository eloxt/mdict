package main

import (
	"crypto/sha1"
	"embed"
	"encoding/hex"
	"errors"
	"fmt"
	"go-mdict/lib"
	"go-mdict/lib/replacer"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	htmltomarkdown "github.com/JohannesKaufmann/html-to-markdown/v2"
	"github.com/gin-contrib/static"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// DictionaryInfo holds the runtime information for a single dictionary.
type DictionaryInfo struct {
	Id        string     `json:"id"`
	Name      string     `json:"name"`
	Path      string     `json:"-"`
	MdxPath   string     `json:"-"`
	MddPath   string     `json:"-"`
	DbPath    string     `json:"-"`
	Db        *gorm.DB   `json:"-"`
	MdxParser *lib.Mdict `json:"-"`
	MddParser *lib.Mdict `json:"-"`
}

var dictionaryManager map[string]*DictionaryInfo

var handler = []replacer.Replacer{
	&replacer.ReplacerCss{},
	&replacer.ReplacerJs{},
	&replacer.ReplacerImage{},
	&replacer.ReplacerSound{},
	&replacer.ReplacerEntry{},
}

//go:embed web/build/client
var buildFS embed.FS

func main() {
	dictionaryManager = make(map[string]*DictionaryInfo)

	loadAndIndexDictionaries()

	embedFs, err := static.EmbedFolder(buildFS, "web/build/client")
	if err != nil {
		log.Fatal("Failed to embed static files:", err)
	}

	router := gin.Default()
	router.Use(static.Serve("/", embedFs))
	apiRouter := router.Group("/api")
	{
		apiRouter.GET("/dictionaries", listDictionariesHandler)
		apiRouter.GET("/lookup/:word", lookupHandler)
		apiRouter.GET("/suggest/:prefix", suggestHandler)
		apiRouter.GET("/resource/:dict/:filename", resourceHandler)
		apiRouter.GET("/iframe/:dict/:filename", iframeHandler)
	}
	indexPageData, _ := buildFS.ReadFile("web/build/client/index.html")
	router.NoRoute(func(c *gin.Context) {
		if strings.HasPrefix(c.Request.RequestURI, "/api") {
			c.Data(http.StatusNotFound, "text/plain", []byte("Not Found"))
			return
		}
		c.Header("Cache-Control", "no-cache")
		c.Data(http.StatusOK, "text/html; charset=utf-8", indexPageData)
	})

	fmt.Println("Server started at :4200")
	fmt.Printf("Loaded %d dictionaries. Use /dictionaries to see the list.\n", len(dictionaryManager))
	if err := router.Run(":4200"); err != nil {
		log.Fatal("Failed to start server:", err)
	}
}

func loadAndIndexDictionaries() {
	const dictRoot = "dict"
	log.Println("Starting to scan for dictionaries in:", dictRoot)

	folders, err := os.ReadDir(dictRoot)
	if err != nil {
		log.Fatalf("Could not read dictionary directory '%s': %v", dictRoot, err)
	}

	for _, folder := range folders {
		if !folder.IsDir() {
			continue
		}

		dictName := folder.Name()
		dictPath := filepath.Join(dictRoot, dictName)
		log.Printf("Processing dictionary folder: %s", dictPath)

		var mdxPath, mddPath, dbPath string

		files, _ := os.ReadDir(dictPath)
		for _, file := range files {
			lowerName := strings.ToLower(file.Name())
			if strings.HasSuffix(lowerName, ".mdx") {
				mdxPath = filepath.Join(dictPath, file.Name())
			} else if strings.HasSuffix(lowerName, ".mdd") {
				mddPath = filepath.Join(dictPath, file.Name())
			} else if strings.HasSuffix(lowerName, ".db") {
				dbPath = filepath.Join(dictPath, file.Name())
			}
		}

		if mdxPath == "" {
			log.Printf("[WARN] No .mdx file found in '%s'. Skipping.", dictPath)
			continue
		}

		if dbPath == "" {
			dbPath = strings.TrimSuffix(mdxPath, filepath.Ext(mdxPath)) + ".db"
			log.Printf("Index file not found for '%s'. Will create at '%s'", mdxPath, dbPath)
			if err := buildIndexForDict(mdxPath, mddPath, dbPath); err != nil {
				log.Printf("[ERROR] Failed to build index for '%s': %v. Skipping.", mdxPath, err)
				continue
			}
		}

		log.Printf("Loading dictionary '%s'...", dictName)
		db, err := SetupDatabase(dbPath)
		if err != nil {
			log.Printf("[ERROR] Failed to connect to database '%s': %v. Skipping.", dbPath, err)
			continue
		}

		mdxParser, err := lib.New(mdxPath)
		if err != nil {
			log.Printf("[ERROR] Failed to create mdx parser for '%s': %v. Skipping.", mdxPath, err)
			continue
		}
		// Fully load the parser's metadata to enable lookups
		if err := mdxParser.LoadMetadata(); err != nil {
			log.Printf("[ERROR] Failed to load metadata for mdx '%s': %v. Skipping.", mdxPath, err)
			continue
		}

		var mddParser *lib.Mdict
		if mddPath != "" {
			mddParser, err = lib.New(mddPath)
			if err != nil {
				log.Printf("[ERROR] Failed to create mdd parser for '%s': %v. Skipping.", mddPath, err)
				continue
			}
			// Fully load the parser's metadata to enable lookups
			if err := mddParser.LoadMetadata(); err != nil {
				log.Printf("[ERROR] Failed to load metadata for mdd '%s': %v. Skipping.", mddPath, err)
				continue
			}
		}

		hashBytes := sha1.Sum([]byte(dictName))
		id := hex.EncodeToString(hashBytes[:])
		dictionaryManager[id] = &DictionaryInfo{
			Id:        id,
			Name:      dictName,
			Path:      dictPath,
			MdxPath:   mdxPath,
			MddPath:   mddPath,
			DbPath:    dbPath,
			Db:        db,
			MdxParser: mdxParser,
			MddParser: mddParser,
		}
		log.Printf("Successfully loaded dictionary: %s", dictName)
	}
}

func buildIndexForDict(mdxPath, mddPath, dbPath string) error {
	log.Printf("Building index for dictionary... This may take a while.")

	database, err := SetupDatabase(dbPath)
	if err != nil {
		return fmt.Errorf("error setting up database: %w", err)
	}

	// Index .mdx file
	log.Printf("Indexing keywords from '%s'...", mdxPath)
	mdxParser, err := lib.New(mdxPath)
	if err != nil {
		return fmt.Errorf("error creating mdx parser: %w", err)
	}
	if err := mdxParser.BuildIndex(); err != nil {
		return fmt.Errorf("error building mdx index: %w", err)
	}
	keywordEntries := mdxParser.KeyBlockData.KeyEntries
	keywordIndices := make([]*KeywordIndex, 0, len(keywordEntries))
	for _, entry := range keywordEntries {
		keywordIndices = append(keywordIndices, &KeywordIndex{
			Keyword:     entry.KeyWord,
			OffsetStart: entry.RecordStartOffset,
			OffsetEnd:   entry.RecordEndOffset,
		})
	}
	if err := BatchInsertKeywords(database, keywordIndices); err != nil {
		return fmt.Errorf("error batch inserting keywords: %w", err)
	}
	log.Printf("Successfully indexed %d keywords.", len(keywordIndices))

	// Index .mdd file if it exists
	if mddPath != "" {
		log.Printf("Indexing resources from '%s'...", mddPath)
		mddParser, err := lib.New(mddPath)
		if err != nil {
			return fmt.Errorf("error creating mdd parser: %w", err)
		}
		if err := mddParser.BuildIndex(); err != nil {
			return fmt.Errorf("error building mdd index: %w", err)
		}
		resourceEntries := mddParser.KeyBlockData.KeyEntries
		resourceIndices := make([]*ResourceIndex, 0, len(resourceEntries))
		for _, entry := range resourceEntries {
			resourceIndices = append(resourceIndices, &ResourceIndex{
				FileName:    entry.KeyWord,
				OffsetStart: entry.RecordStartOffset,
				OffsetEnd:   entry.RecordEndOffset,
			})
		}
		if err := BatchInsertResources(database, resourceIndices); err != nil {
			return fmt.Errorf("error batch inserting resources: %w", err)
		}
		log.Printf("Successfully indexed %d resources.", len(resourceIndices))
	}

	log.Printf("Successfully built index at '%s'", dbPath)
	return nil
}

func listDictionariesHandler(c *gin.Context) {
	dictInfos := make([]*DictionaryInfo, 0, len(dictionaryManager))
	for _, dictInfo := range dictionaryManager {
		dictInfos = append(dictInfos, dictInfo)
	}
	c.JSON(http.StatusOK, dictInfos)
}

func lookupHandler(c *gin.Context) {
	word := c.Param("word")
	dictId := c.Query("dict")

	if dictId == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Query parameter 'dict' is required."})
		return
	}

	dictInfo, ok := dictionaryManager[dictId]
	if !ok {
		c.JSON(http.StatusNotFound, gin.H{"error": fmt.Sprintf("Dictionary '%s' not found or not loaded.", dictId)})
		return
	}

	keywordIndex, err := FindKeyword(dictInfo.Db, word)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": fmt.Sprintf("Word '%s' not found in dictionary '%s'.", word, dictId)})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Database query failed: " + err.Error()})
		return
	}

	entry := &lib.MDictKeywordEntry{
		KeyWord:           keywordIndex.Keyword,
		RecordStartOffset: keywordIndex.OffsetStart,
		RecordEndOffset:   keywordIndex.OffsetEnd,
	}

	definition, err := dictInfo.MdxParser.LocateByKeywordEntry(entry)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to locate definition: " + err.Error()})
		return
	}
	html := string(definition)

	if newWord, ok0 := strings.CutPrefix(html, "@@@LINK="); ok0 {
		newWord = strings.TrimRight(newWord, "\r\n\000")
		keywordIndex, err := FindKeyword(dictInfo.Db, newWord)
		if err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				c.JSON(http.StatusNotFound, gin.H{"error": fmt.Sprintf("Word '%s' not found in dictionary '%s'.", word, dictId)})
				return
			}
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Database query failed: " + err.Error()})
			return
		}
		entry = &lib.MDictKeywordEntry{
			KeyWord:           keywordIndex.Keyword,
			RecordStartOffset: keywordIndex.OffsetStart,
			RecordEndOffset:   keywordIndex.OffsetEnd,
		}
		definition, err = dictInfo.MdxParser.LocateByKeywordEntry(entry)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to locate linked definition: " + err.Error()})
			return
		}
		html = string(definition)
	}

	for _, rep := range handler {
		html = rep.Replace(dictId, html)
	}
	html = fmt.Sprintf(replacer.WordDefinitionTemplate, html)

	markdown, err := htmltomarkdown.ConvertString(html)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to convert html to markdown: " + err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"word":       word,
		"dictionary": dictId,
		"html":       html,
		"markdown":   markdown,
	})
}

func suggestHandler(c *gin.Context) {
	prefix := c.Param("prefix")
	dictName := c.Query("dict")

	if dictName == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Query parameter 'dict' is required."})
		return
	}

	dictInfo, ok := dictionaryManager[dictName]
	if !ok {
		c.JSON(http.StatusNotFound, gin.H{"error": fmt.Sprintf("Dictionary '%s' not found or not loaded.", dictName)})
		return
	}

	if len(prefix) < 2 {
		c.JSON(http.StatusOK, []string{})
		return
	}

	keywords, err := FuzzyFindKeywords(dictInfo.Db, prefix, 20)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Database query failed: " + err.Error()})
		return
	}

	c.JSON(http.StatusOK, keywords)
}

func resourceHandler(c *gin.Context) {
	filename := c.Param("filename")
	dictId := c.Param("dict")

	if dictId == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Query parameter 'dict' is required."})
		return
	}

	dictInfo, ok := dictionaryManager[dictId]
	if !ok {
		c.JSON(http.StatusNotFound, gin.H{"error": fmt.Sprintf("Dictionary '%s' not found or not loaded.", dictId)})
		return
	}

	if dictInfo.MddParser == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": fmt.Sprintf("No .mdd file loaded for dictionary '%s'.", dictId)})
		return
	}

	resourceIndex, err := FindResource(dictInfo.Db, filename)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": fmt.Sprintf("Resource '%s' not found in dictionary '%s'.", filename, dictId)})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Database query failed for resource: " + err.Error()})
		return
	}

	entry := &lib.MDictKeywordEntry{
		KeyWord:           resourceIndex.FileName,
		RecordStartOffset: resourceIndex.OffsetStart,
		RecordEndOffset:   resourceIndex.OffsetEnd,
	}

	data, err := dictInfo.MddParser.LocateByKeywordEntry(entry)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to locate resource data: " + err.Error()})
		return
	}

	// Determine Content-Type from file extension
	ext := strings.ToLower(filepath.Ext(filename))
	var contentType string
	switch ext {
	case ".png":
		contentType = "image/png"
	case ".jpg", ".jpeg":
		contentType = "image/jpeg"
	case ".gif":
		contentType = "image/gif"
	case ".mp3":
		contentType = "audio/mpeg"
	case ".wav":
		contentType = "audio/wav"
	case ".spx":
		contentType = "audio/ogg"
	default:
		contentType = "application/octet-stream"
	}

	c.Data(http.StatusOK, contentType, data)
}

func iframeHandler(c *gin.Context) {
	dictId := c.Param("dict")
	filename := c.Param("filename")
	if dictId == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Query parameter 'dict' is required."})
	}

	if filename == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Query parameter 'filename' is required."})
	}

	dictInfo, ok := dictionaryManager[dictId]
	if !ok {
		c.JSON(http.StatusNotFound, gin.H{"error": fmt.Sprintf("Dictionary '%s' not found or not loaded.", dictId)})
		return
	}

	c.File(filepath.Join(dictInfo.Path, filename))
}
