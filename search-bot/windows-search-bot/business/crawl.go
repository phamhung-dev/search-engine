package business

import (
	"context"
	"fmt"
	"math"
	"os"
	"path"
	"strings"
	"time"
	"windows-search-bot/models"
	"windows-search-bot/utils"

	"golang.org/x/sys/windows"
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

		files := func(parentFolderPath string, filesStack *utils.Stack[string]) []models.File {
			parentFolderPathPtr, err := windows.UTF16PtrFromString(fmt.Sprintf("%s*", parentFolderPath))

			if err != nil {
				return nil
			}

			var fileInfo windows.Win32finddata

			handle, err := windows.FindFirstFile(
				parentFolderPathPtr,
				&fileInfo,
			)

			if err != nil || handle == windows.InvalidHandle {
				return nil
			}

			defer windows.FindClose(handle)

			files := []models.File{}

			for {
				err := windows.FindNextFile(
					handle,
					&fileInfo,
				)

				if err != nil {
					break
				}

				fileName := windows.UTF16ToString(fileInfo.FileName[:])

				if fileName != "." && fileName != ".." {
					if fileInfo.FileAttributes&windows.FILE_ATTRIBUTE_DIRECTORY != 0 {
						filesStack.Push(fmt.Sprintf("%s%s\\", parentFolderPath, fileName))
					} else {
						file := processFileCrawl(parentFolderPath, fileInfo)

						files = append(files, file)
					}
				}
			}

			return files
		}(parentFolderPath, &filesStack)

		if len(files) > 0 {
			if err := business.storage.BatchCreate(context, files); err != nil {
				fmt.Println(err)
			}
		}
	}
}

func initStack(context context.Context, storage CrawlStorage) (filesStack utils.Stack[string]) {
	drives := getExistDrives()

	lastFileCreated, err := storage.FindLastFileCreated(context)

	if err != nil {
		filesStack.BatchPush(drives...)
		return
	}

	lastBackslashIdx := strings.LastIndex(lastFileCreated.Path, "\\")

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

		func(parentFolderPath string, filesStack *utils.Stack[string]) {
			parentFolderPathPtr, err := windows.UTF16PtrFromString(fmt.Sprintf("%s*", parentFolderPath))

			if err != nil {
				return
			}

			var fileInfo windows.Win32finddata

			handle, err := windows.FindFirstFile(
				parentFolderPathPtr,
				&fileInfo,
			)

			if err != nil || handle == windows.InvalidHandle {
				return
			}

			defer windows.FindClose(handle)

			for {
				err := windows.FindNextFile(
					handle,
					&fileInfo,
				)

				if err != nil {
					break
				}

				fileName := windows.UTF16ToString(fileInfo.FileName[:])

				if fileName != "." && fileName != ".." && fileInfo.FileAttributes&windows.FILE_ATTRIBUTE_DIRECTORY != 0 {
					folderPath := fmt.Sprintf("%s%s\\", parentFolderPath, fileName)

					filesStack.Push(folderPath)

					if strings.Contains(lastFolderPathCreated, folderPath) {
						break
					}
				}
			}

		}(parentFolderPath, &filesStack)
	}

	files := func(lastFolderPathCreated string) []models.File {
		parentFolderPathPtr, err := windows.UTF16PtrFromString(fmt.Sprintf("%s*", lastFolderPathCreated))

		if err != nil {
			return nil
		}

		var fileInfo windows.Win32finddata

		handle, err := windows.FindFirstFile(
			parentFolderPathPtr,
			&fileInfo,
		)

		if err != nil || handle == windows.InvalidHandle {
			return nil
		}

		defer windows.FindClose(handle)

		lastFileNameCreated := lastFileCreated.Path[lastBackslashIdx+1:]

		allowCreate := false

		files := []models.File{}

		for {
			err := windows.FindNextFile(
				handle,
				&fileInfo,
			)

			if err != nil {
				break
			}

			fileName := windows.UTF16ToString(fileInfo.FileName[:])

			if fileName != "." && fileName != ".." && fileInfo.FileAttributes&windows.FILE_ATTRIBUTE_DIRECTORY == 0 {
				if allowCreate {
					file := processFileCrawl(lastFolderPathCreated, fileInfo)

					files = append(files, file)

					continue
				}

				if fileName == lastFileNameCreated {
					allowCreate = true
				}
			}
		}

		return files
	}(lastFolderPathCreated)

	if len(files) > 0 {
		if err := storage.BatchCreate(context, files); err != nil {
			fmt.Println(err)
		}
	}

	return
}

func getExistDrives() []string {
	fakeDrives := []string{"A", "B", "C", "D", "E", "F", "G", "H", "I", "J", "K", "L", "M", "N", "O", "P", "Q", "R", "S", "T", "U", "V", "W", "X", "Y", "Z"}
	realDrives := []string{}

	for _, drive := range fakeDrives {
		pathDrive := fmt.Sprintf("%s:\\", drive)

		f, err := os.Open(pathDrive)

		if err == nil {
			realDrives = append(realDrives, pathDrive)
			f.Close()
		}
	}

	return realDrives
}

func processFileCrawl(parentFolderPath string, fileInfo windows.Win32finddata) models.File {
	fileName := windows.UTF16ToString(fileInfo.FileName[:])

	filePath := fmt.Sprintf("%s%s", parentFolderPath, fileName)

	extension := path.Ext(fileName)

	size := int64(fileInfo.FileSizeHigh)*int64(math.MaxUint32+1) + int64(fileInfo.FileSizeLow)

	content := utils.GetDocumentContent(filePath, size)

	readOnly := fileInfo.FileAttributes&windows.FILE_ATTRIBUTE_READONLY != 0

	hidden := fileInfo.FileAttributes&windows.FILE_ATTRIBUTE_HIDDEN != 0

	createdAt := time.Unix(0, fileInfo.CreationTime.Nanoseconds())
	lastModifiedAt := time.Unix(0, fileInfo.LastWriteTime.Nanoseconds())
	lastAccessedAt := time.Unix(0, fileInfo.LastAccessTime.Nanoseconds())

	return models.File{
		Name:               fileName,
		Path:               filePath,
		Extension:          extension,
		Size:               size,
		Content:            content,
		ReadOnly:           readOnly,
		Hidden:             hidden,
		FileCreatedAt:      createdAt,
		FileLastModifiedAt: lastModifiedAt,
		FileLastAccessedAt: lastAccessedAt,
	}
}
