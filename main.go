package main

import (
	"fmt"
	"github.com/fsnotify/fsnotify"
	"github.com/gorilla/mux"
	"io/ioutil"
	"net/http"
	"os"
	"path/filepath"
	"sort"
	"strconv"
	"strings"
)


type Watch struct {
	watch *fsnotify.Watcher;
}
var imgslice = []string{}

func main() {
	r := mux.NewRouter()
	r.Handle("/", http.FileServer(http.Dir("static")))
	watch, _ := fsnotify.NewWatcher()
	w := Watch{
		watch: watch,
	}
	var watchDirPath = "/Users/zhangjianxin/Library/Containers/com.tencent.qq/Data/Library/Caches/Images"
	w.watchDir(watchDirPath);
	http.Handle("/img/", http.StripPrefix("/img/", http.FileServer(http.Dir(watchDirPath))))
	http.HandleFunc("/watcher",watcher)
	r.HandleFunc("/random/{args}", random)
	files, _ := ioutil.ReadDir(watchDirPath)
	for _, f := range files {
		imgslice = append(imgslice, string(f.Name()))
	}
	http.Handle("/", r)
	http.ListenAndServe(":8000", nil)
}

func watcher(w http.ResponseWriter, r *http.Request) {
	for {
		select {
		default:
			w.Write([]byte("64 bytes or fewer"));
		}
	}
}

func random(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	params := vars["args"];
	page, err := strconv.Atoi(strings.Split(params, ",")[0])
	size, err := strconv.Atoi(strings.Split(params, ",")[1])
	if err != nil {
		fmt.Println("is Error ")
	}
	// s = 1 * 10 - 10 = 0
	// en = 1 * 10
	endIndex := (page * size);
	nowIndex := endIndex - size;
	str := imgslice[nowIndex:endIndex];
	var x = []byte{}
	for i := 0; i < len(str); i++ {
		b := []byte(str[i])
		for j := 0; j < len(b); j++ {
			x = append(x,b[j])
		}
	}
	justString := strings.Join(str, ",")
	w.Write([]byte("[" + justString + "]"))
}
//监控目录
func (w *Watch) watchDir(dir string) {
	filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
		if info.IsDir() {
			path, err := filepath.Abs(path);
			if err != nil {
				return err;
			}
			err = w.watch.Add(path);
			if err != nil {
				return err;
			}
			fmt.Println("watch : ", path);
		}
		return nil;
	});
	go func() {
		for {
			select {
			case ev := <-w.watch.Events:
				{
					if ev.Op&fsnotify.Create == fsnotify.Create {
						fmt.Println("create File : ", ev.Name);
						imgslice = append(imgslice, string(ev.Name))
						// 先搞事情在排序- -
						sort.Sort(sort.Reverse(sort.StringSlice(imgslice)))
						fi, err := os.Stat(ev.Name);
						if err == nil && fi.IsDir() {
							w.watch.Add(ev.Name);
							fmt.Println("add Watch : ", ev.Name);
						}
					}
				}
			case err := <-w.watch.Errors:
				{
					fmt.Println("error : ", err);
					return;
				}
			}
		}
	}();
}
