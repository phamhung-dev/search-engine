package business

import (
	"context"
	"fmt"
	"linux-search-bot/models"
	"log"
	"os"
	"strings"
	"time"

	"github.com/radovskyb/watcher"
)

type MonitorStorage interface {
	Create(context context.Context, file *models.File) error
	Update(context context.Context, file *models.File) error
	Find(context context.Context, conditions map[string]interface{}) (*models.File, error)
	Delete(context context.Context, file *models.File) error
	ListFilesInFolder(context context.Context, folderName string) ([]models.File, error)
}

type monitorBusiness struct {
	storage MonitorStorage
}

func NewMonitorBusiness(storage MonitorStorage) *monitorBusiness {
	return &monitorBusiness{storage: storage}
}

func (business *monitorBusiness) Monitor(context context.Context) {
	w := watcher.New()
	go func() {
		for {
			select {
			case event := <-w.Event:
				if !event.IsDir() && canMonitor(event) {
					triggerEventChange(context, event, business.storage)
				}
			case err := <-w.Error:
				fmt.Println(err)
			case <-w.Closed:
				return
			}
		}
	}()

	w.IgnoreHiddenFiles(true)

	if err := w.AddRecursive("/"); err != nil {
		fmt.Println(err)
	}

	go func() {
		w.Wait()
		w.TriggerEvent(watcher.Create, nil)
		w.TriggerEvent(watcher.Remove, nil)
		w.TriggerEvent(watcher.Rename, nil)
		w.TriggerEvent(watcher.Chmod, nil)
		w.TriggerEvent(watcher.Write, nil)
		w.TriggerEvent(watcher.Move, nil)
	}()

	if err := w.Start(time.Millisecond * 100); err != nil {
		log.Fatalln(err)
	}
}

func triggerEventChange(context context.Context, event watcher.Event, storage MonitorStorage) {
	if !event.IsDir() {
		switch event.Op.String() {
		case "CREATE":
			eventCreate(context, event, storage)
		case "REMOVE":
			eventRemove(context, event, storage)
		case "RENAME", "CHMOD", "WRITE", "MOVE":
			eventUpdate(context, event, storage)
		default:
			return
		}
	}
}

func eventCreate(context context.Context, event watcher.Event, storage MonitorStorage) {
	lastBackslashIdx := 0
	folderName := event.Path

	for lastBackslashIdx > -1 {
		lastBackslashIdx = strings.LastIndex(folderName, "/")

		if lastBackslashIdx < 0 {
			break
		}

		folderName = fmt.Sprintf("%s%%", folderName[:lastBackslashIdx])

		if files, err := storage.ListFilesInFolder(context, folderName); err == nil && len(files) > 0 {
			file := processFile(event.Path, event.FileInfo)

			if err := storage.Create(context, &file); err != nil {
				fmt.Println(err)
			}

			break
		}
	}
}

func eventRemove(context context.Context, event watcher.Event, storage MonitorStorage) {
	conditions := map[string]interface{}{"path": event.OldPath}

	file, err := storage.Find(context, conditions)

	if err != nil {
		fmt.Println(err)
		return
	}

	if err := storage.Delete(context, file); err != nil {
		fmt.Println(err)
	}
}

func eventUpdate(context context.Context, event watcher.Event, storage MonitorStorage) {
	conditions := map[string]interface{}{"path": event.OldPath}

	file, err := storage.Find(context, conditions)

	if err != nil {
		fmt.Println(err)
		return
	}

	fileInfo, err := os.Lstat(event.Path)

	if err != nil {
		fmt.Println(err)
		return
	}

	fileUpdate := processFile(event.Path, fileInfo)

	file.Name = fileUpdate.Name
	file.Path = fileUpdate.Path
	file.Extension = fileUpdate.Extension
	file.Size = fileUpdate.Size
	file.Content = fileUpdate.Content
	file.ReadOnly = fileUpdate.ReadOnly
	file.Hidden = fileUpdate.Hidden
	file.FileCreatedAt = fileUpdate.FileCreatedAt
	file.FileLastModifiedAt = fileUpdate.FileLastModifiedAt
	file.FileLastAccessedAt = fileUpdate.FileLastAccessedAt

	if err := storage.Update(context, file); err != nil {
		fmt.Println(err)
	}
}

func canMonitor(event watcher.Event) bool {
	if event.Sys() == nil {
		fmt.Printf("cannot monitor data in %s\n", event.Path)
		return false
	}

	return true
}
