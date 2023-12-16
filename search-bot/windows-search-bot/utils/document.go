package utils

import (
	"archive/zip"
	"fmt"
	"io"
	"os"
	"regexp"
	"strings"

	"github.com/h2non/filetype"
)

var types = map[string]string{
	"docx": "application/vnd.openxmlformats-officedocument.wordprocessingml.document",
	"doc":  "application/msword",
	"xlsx": "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet",
	"xls":  "application/vnd.ms-excel",
	"pptx": "application/vnd.openxmlformats-officedocument.presentationml.presentation",
	"ppt":  "application/vnd.ms-powerpoint",
}

func GetDocumentContent(filePath string, fileSize int64) string {
	content := ""

	file, err := os.Open(filePath)

	if err != nil {
		return ""
	}

	defer file.Close()

	buf := make([]byte, fileSize)

	if _, err := file.Read(buf); err != nil {
		return ""
	}

	documentType, err := filetype.Match(buf)

	if err != nil {
		return ""
	}

	reader, err := zip.NewReader(file, fileSize)

	if err != nil {
		return ""
	}

	// printAllDocumentFile(reader)

	regexPattern := ``
	regexNormalizeOptions := map[string]string{}

	switch documentType.MIME.Value {
	case types["doc"], types["docx"]:
		regexPattern = `word/document.xml`
	case types["xlsx"], types["xls"]:
		regexPattern = `xl/(worksheets/sheet|sharedStrings).*\.xml`
		regexNormalizeOptions[`t="s"><v>.*?</v></c>`] = ">"
	case types["pptx"], types["ppt"]:
		regexPattern = `ppt/slides/slide.*\.xml`
	default:
		return ""
	}

	for _, f := range reader.File {
		r := regexp.MustCompile(regexPattern)

		if r.MatchString(f.Name) {
			documentFile, err := f.Open()

			if err != nil {
				return ""
			}

			defer documentFile.Close()

			contentBytes, err := io.ReadAll(documentFile)

			if err != nil {
				return ""
			}

			content = fmt.Sprintf("%s %s", content, string(contentBytes))
		}
	}
	return normalizeContent(content, regexNormalizeOptions)
}

// func printAllDocumentFile(reader *zip.Reader) {
// 	for _, f := range reader.File {
// 		fmt.Println(f.Name)
// 	}
// }

func normalizeQuotes(in rune) rune {
	switch in {
	case '“', '”':
		return '"'
	case '‘', '’':
		return '\''
	}
	return in
}

func normalizeContent(text string, regexOptions map[string]string) string {
	for regexPattern := range regexOptions {
		pattern := regexp.MustCompile(regexPattern)

		text = pattern.ReplaceAllString(text, regexOptions[regexPattern])
	}

	brakets := regexp.MustCompile(`<.*?>`)
	quotes := regexp.MustCompile(`&quot;`)
	spaces := regexp.MustCompile(`\s+`)

	text = brakets.ReplaceAllString(text, " ")
	text = quotes.ReplaceAllString(text, "\"")
	text = spaces.ReplaceAllString(text, " ")

	text = strings.Map(normalizeQuotes, text)

	return strings.TrimSpace(text)
}
