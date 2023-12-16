package business

import (
	"context"
	"fmt"
	"io/fs"
	"linux-search-bot/models"
	"linux-search-bot/utils"
	"os"
	"path"
	"strings"
	"syscall"
	"time"
)

type CrawlStorage interface {
	BatchCreate(context context.Context, files []models.File) error
	FindLastFileCreated(context context.Context) (*models.File, error)
}

type crawlBusiness struct {
	storage CrawlStorage
}

func NewCrawlBusiness(storage CrawlStorage) *crawlBusiness {
	return &crawlBusiness{storage: storage}
}

func (business *crawlBusiness) Crawl(context context.Context) {
	filesStack := initStack(context, business.storage)

	for !filesStack.IsEmpty() {
		parentFolderPath, ok := filesStack.Pop()

		if !ok {
			break
		}

		entries, err := os.ReadDir(parentFolderPath)

		if err != nil {
			fmt.Printf("cannot collect data in %s\n", parentFolderPath)
			continue
		}

		files := []models.File{}

		for _, entry := range entries {
			if entry.IsDir() {
				filesStack.Push(fmt.Sprintf("%s%s/", parentFolderPath, entry.Name()))
			} else {
				filePath := fmt.Sprintf("%s%s", parentFolderPath, entry.Name())

				fileInfo, err := entry.Info()

				if err != nil {
					fmt.Printf("cannot collect data in %s\n", filePath)
					continue
				}

				file := processFile(filePath, fileInfo)

				files = append(files, file)
			}
		}

		if len(files) > 0 {
			if err := business.storage.BatchCreate(context, files); err != nil {
				fmt.Println(err)
			}
		}
	}

}

func initStack(context context.Context, storage CrawlStorage) (filesStack utils.Stack[string]) {
	drives := []string{"/"}

	lastFileCreated, err := storage.FindLastFileCreated(context)

	if err != nil {
		filesStack.BatchPush(drives...)
		return
	}

	lastBackslashIdx := strings.LastIndex(lastFileCreated.Path, "/")

	lastFolderPathCreated := ""

	if lastBackslashIdx < 0 {
		lastFolderPathCreated = lastFileCreated.Path
	} else {
		lastFolderPathCreated = lastFileCreated.Path[:lastBackslashIdx+1]
	}

	for _, drive := range drives {
		filesStack.Push(drive)

		if strings.Contains(lastFolderPathCreated, drive) {
			break
		}
	}

	for !filesStack.IsEmpty() {
		parentFolderPath, ok := filesStack.Pop()

		if !ok || parentFolderPath == lastFolderPathCreated {
			break
		}

		entries, err := os.ReadDir(parentFolderPath)

		if err != nil {
			fmt.Printf("cannot collect data in %s\n", parentFolderPath)
			continue
		}

		for _, entry := range entries {
			if entry.IsDir() {
				subEntryPath := fmt.Sprintf("%s%s/", parentFolderPath, entry.Name())

				filesStack.Push(subEntryPath)

				if strings.Contains(lastFolderPathCreated, subEntryPath) {
					break
				}
			}
		}
	}

	entriesInLastFolder, err := os.ReadDir(lastFolderPathCreated)

	if err != nil {
		fmt.Printf("cannot collect data in %s\n", lastFolderPathCreated)
		return
	}

	lastFileNameCreated := lastFileCreated.Path[lastBackslashIdx+1:]

	allowCreate := false

	files := []models.File{}

	for _, entry := range entriesInLastFolder {
		if !entry.IsDir() {
			if allowCreate {
				filePath := fmt.Sprintf("%s%s", lastFolderPathCreated, entry.Name())

				fileInfo, err := entry.Info()

				if err != nil {
					fmt.Printf("cannot collect data in %s\n", filePath)
					continue
				}

				file := processFile(filePath, fileInfo)

				files = append(files, file)

				continue
			}

			if entry.Name() == lastFileNameCreated {
				allowCreate = true
			}
		}
	}

	if len(files) > 0 {
		if err := storage.BatchCreate(context, files); err != nil {
			fmt.Println(err)
		}
	}

	return
}

func processFile(filePath string, fileInfo fs.FileInfo) models.File {
	name := fileInfo.Name()

	stat := fileInfo.Sys().(*syscall.Stat_t)

	createdAt := time.Unix(int64(stat.Ctim.Sec), int64(stat.Ctim.Nsec))

	lastModifiedAt := fileInfo.ModTime()

	lastAccessedAt := time.Unix(int64(stat.Atim.Sec), int64(stat.Atim.Nsec))

	content := utils.GetDocumentContent(filePath, fileInfo.Size())

	return models.File{
		Name:               name,
		Path:               filePath,
		Extension:          path.Ext(name),
		Size:               fileInfo.Size(),
		Content:            content,
		ReadOnly:           fileInfo.Mode()&os.ModePerm == 0444,
		Hidden:             strings.HasPrefix(name, "."),
		FileCreatedAt:      createdAt,
		FileLastModifiedAt: lastModifiedAt,
		FileLastAccessedAt: lastAccessedAt,
	}
}
