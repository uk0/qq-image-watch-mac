package qq_image_watch_mac

import (
	"fmt"
	"github.com/fsnotify/fsnotify"
	"html/template"
	"io/ioutil"
	"net/http"
	"os"
	"path/filepath"
	"strings"
)

type Image struct {
	Src string
}

type ImagePage struct {
	Title  string
	Images []Image
}

type Watch struct {
	watch *fsnotify.Watcher;
}

func (box *ImagePage) AddItem(item Image) []Image {
	box.Images = append(box.Images, item)
	return box.Images
}

var data = ImagePage{
	Title: "My Image list",
	Images: [] Image{
	},
}

func main() {
	tmpl := template.Must(template.ParseFiles("./template/index.html"))

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {

		fmt.Println("图片数量 : ", len(data.Images))
		tmpl.Execute(w, data)
	})
	watch, _ := fsnotify.NewWatcher()
	w := Watch{
		watch: watch,
	}
	var watchDirPath = "/Users/xxxx/Library/Containers/com.tencent.qq/Data/Library/Caches/Images"
	w.watchDir(watchDirPath);
	http.Handle("/html/", http.StripPrefix("/html/", http.FileServer(http.Dir(watchDirPath))))
	files, _ := ioutil.ReadDir(watchDirPath)
	for _, f := range files {
		addPath(string(f.Name()))
	}
	http.ListenAndServe(":8000", nil)
}


func addPath(name string){
	path := strings.Split(string(name), "/")
	data.AddItem(Image{Src:"/html/" + path[len(path)-1]})
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
						addPath(string(ev.Name))
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
