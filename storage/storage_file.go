package storage

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/fsnotify/fsnotify"
	"io/fs"
	"io/ioutil"
	"log"
	"path/filepath"
	"runtime"
	"strings"
	"sync"
	"xhttp/pkg/tools"
)

const (
	projectDir     = "demo"
	routesFileName = "routes.json"
)

var _, fileName, _, _ = runtime.Caller(0)

type FileStorage struct {
	once     sync.Once
	Projects map[string]*Project
	lock     sync.Mutex
	Event    chan Event
}

func NewFileStorage() *FileStorage {
	return &FileStorage{
		Projects: make(map[string]*Project),
		Event:    make(chan Event, 10),
	}
}

func (f *FileStorage) Init() (err error) {
	f.once.Do(func() {
		projects := make(map[string]*Project)
		projects, err = f.GetAll()
		if err != nil {
			return
		}
		f.Projects = projects
	})
	return
}

func (f *FileStorage) WatchEvent() chan Event {
	return f.Event
}

func (f *FileStorage) Watch() (err error) {
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		return
	}
	defer func() {
		err = watcher.Close()
		if err != nil {
			log.Println("err in Watch.Close(),err:" + err.Error())
		}
	}()
	watcherFullPath := filepath.Dir(fileName) + "/" + projectDir
	done := make(chan bool)
	go func() {
		defer close(done)
		for {
			select {
			case e, ok := <-watcher.Events:
				if !ok {
					return
				}
				changeFileName := e.Name
				if strings.HasSuffix(changeFileName, routesFileName+"~") {
					changeFileName = strings.TrimRight(changeFileName, "~")
				}

				log.Println(fmt.Sprintf("监听到文件 %s 变化| ", changeFileName))

				var project string
				pathArr := strings.Split(changeFileName, "/")
				if strings.HasSuffix(changeFileName, routesFileName) {
					project = pathArr[len(pathArr)-2]
				} else {
					project = pathArr[len(pathArr)-1]
				}
				log.Println("project is:", project)

				if err = watcher.Add(watcherFullPath + "/" + project); err != nil {
					log.Println("watcher.Add inner err:" + err.Error())
					f.lock.Lock()
					delete(f.Projects, project)
					if err = watcher.Remove(watcherFullPath + "/" + project); err != nil {
						log.Println("remove watch project:" + project)
					}
					f.lock.Unlock()
					p, _ := f.Get(project)
					evt := Event{
						Op:      OpDel,
						Project: p,
					}
					f.Event <- evt
					continue
				}
				switch e.Op {
				case fsnotify.Create:
					log.Println("创建事件", e.Op)
					f.updateProject(project, changeFileName)
				case fsnotify.Write:
					log.Println("写入事件", e.Op)
					f.updateProject(project, changeFileName)
				case fsnotify.Remove:
					log.Println("删除事件", e.Op)
					f.updateProject(project, changeFileName)
				case fsnotify.Rename:
					log.Println("重命名事件", e.Op)
					f.updateProject(project, changeFileName)
				case fsnotify.Chmod:
					log.Println("属性修改事件", e.Op)
					f.updateProject(project, changeFileName)
				default:
					log.Println("some thing else")
					f.updateProject(project, changeFileName)
				}
				p, _ := f.Get(project)
				evt := Event{
					Op:      OpMod,
					Project: p,
				}
				log.Println("new event send")
				f.Event <- evt
			case err, ok := <-watcher.Errors:
				if !ok {
					return
				}
				log.Println("error:", err)
			}
		}
	}()
	err = watcher.Add(watcherFullPath)
	if err != nil {
		log.Println("watcher.Add err:" + err.Error())
		return
	}
	var fsInfo []fs.FileInfo
	fsInfo, err = ioutil.ReadDir(watcherFullPath)
	for _, o := range fsInfo {
		if o.IsDir() {
			log.Println("add watch project: ", o.Name())
			if err = watcher.Add(watcherFullPath + "/" + o.Name()); err != nil {
				log.Println("watcher.Add project err:", o.Name())
				return
			}
		}
	}
	<-done
	return nil
}

func (f *FileStorage) Get(project string) (p *Project, err error) {
	if o, ok := f.Projects[project]; ok {
		return o, nil
	}
	return nil, errors.New("no project:" + project)
}

func (f *FileStorage) updateProject(project, fullPath string) {
	routesUpdateEvent := false
	if strings.HasSuffix(fullPath, routesFileName) {
		routesUpdateEvent = true
	}
	item := &Project{
		Name: project,
		APIs: nil,
	}
	oldProject, _ := f.Get(project)
	if oldProject != nil {
		item.APIs = oldProject.APIs
	}
	if routesUpdateEvent {
		_, apis, err := f.readRoutes(fullPath)
		if err != nil {
			return
		}
		log.Println("update apis:", apis)
		item.APIs = apis
	}
	f.lock.Lock()
	f.Projects[project] = item
	f.lock.Unlock()
}

func (f *FileStorage) readRoutes(path string) (project string, apis []*API, err error) {
	if !tools.CheckFileIsExist(path) {
		log.Println("no file:" + path)
		return project, apis, errors.New("no file:" + path)
	}
	pathArr := strings.Split(path, "/")
	project = pathArr[len(pathArr)-2]
	var bs []byte
	bs, err = ioutil.ReadFile(path)
	if err != nil {
		log.Println(" readRoutes path:" + path + ",err:" + err.Error())
		return
	}
	if len(bs) == 0 {
		return project, apis, errors.New("no config content:" + path)
	}
	if err = json.Unmarshal(bs, &apis); err != nil {
		log.Println("readRoutes json.Unmarshal routes.json err:" + err.Error())
		return
	}
	return
}

func (f *FileStorage) GetAll() (projects map[string]*Project, err error) {
	parentDir := filepath.Dir(fileName)
	var fsInfo []fs.FileInfo
	fsInfo, err = ioutil.ReadDir(parentDir + "/" + projectDir)
	if err != nil {
		return
	}
	projectsMap := make(map[string]*Project)
	for _, o := range fsInfo {
		if o.IsDir() {
			name := o.Name()
			var apis []*API
			_, apis, err = f.readRoutes(parentDir + "/" + projectDir + "/" + name + "/" + routesFileName)
			if err != nil {
				continue
			}
			item := &Project{
				Name: name,
				APIs: apis,
			}
			if _, ok := projectsMap[name]; !ok {
				projectsMap[name] = item
			}
		}
	}
	f.Projects = projectsMap
	return f.Projects, nil
}
