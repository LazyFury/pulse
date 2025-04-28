package template

import (
	"embed"
	"encoding/json"
	"fmt"
	"html/template"
	"net/http"
	"os"
	"path"
	"path/filepath"
	"strings"

	"github.com/charmbracelet/log"
	"github.com/lazyfury/pulse/framework/constants"
	"github.com/lazyfury/pulse/framework/global"
)

/**
@ usege
```golang
html := template.Must(_template.ParseGlob(template.New("main"), "templates", "*.html")).Funcs(_template.TemplateFuncsMap)
engine.SetHTMLTemplate(html)
```
	**/

// TemplateRenderer is a custom html/template renderer for Echo framework
type (
	TemplateRenderer struct {
		Templates *template.Template
	}

	tplFile struct {
		Name string
		Path string
	}
)

func ParseGlob(tpl *template.Template, dir string, pattern string) (t *template.Template, err error) {
	return ParseGlobWithFiles(tpl, readOsFile, filesToTplFiles(dir, pattern, getOsDir())...)
}

func ParseGlobEmbedFS(tpl *template.Template, fs embed.FS, dir string, pattern string) (t *template.Template, err error) {
	return ParseGlobWithFiles(tpl, readEmbedFile(fs), filesToTplFiles(dir, pattern, getEmbedDir(fs))...)
}

// ParseGlob 自定义模版解析，扫描子目录
func ParseGlobWithFiles(tpl *template.Template, readFile func(string) ([]byte, error), files ...*tplFile) (t *template.Template, err error) {
	t = tpl
	for _, file := range files {
		log.Debug(fmt.Sprintf("挂载模板：%s", file.Path))
		b, err := readFile(file.Path)
		if err != nil {
			return t, err
		}
		s := string(b)
		name := file.Name
		var tmpl *template.Template
		if t == nil {
			t = template.New(name)
		}
		if name == t.Name() {
			tmpl = t
		} else {
			tmpl = t.New(name)
		}
		_, err = tmpl.Parse(s)
		if err != nil {
			return nil, err
		}
	}
	return t, nil
}

func readOsFile(path string) (b []byte, err error) {
	return os.ReadFile(path)
}

func readEmbedFile(fs embed.FS) func(string) ([]byte, error) {
	return func(path string) (b []byte, err error) {
		b, err = fs.ReadFile(path)
		return
	}
}

func getOsDir() func(string, string) (files []os.DirEntry, err error) {
	return func(dir string, suffix string) (files []os.DirEntry, err error) {
		files, err = os.ReadDir(dir)
		if err != nil {
			return
		}
		return
	}
}

func getEmbedDir(fs embed.FS) func(string, string) (files []os.DirEntry, err error) {
	return func(dir string, suffix string) (files []os.DirEntry, err error) {
		files, err = fs.ReadDir(dir)
		if err != nil {
			return
		}
		return
	}
}

// 目录下的所有文件
func filesToTplFiles(dir string, suffix string, getDirs func(string, string) (files []os.DirEntry, err error)) (arr []*tplFile) {

	files, err := getDirs(dir, suffix)
	if err != nil {
		return
	}
	for _, file := range files {
		if file.IsDir() {
			log.Debug("扫描子目录：" + file.Name())
			arr = append(arr, filesToTplFiles(path.Join(dir, file.Name()), suffix, getDirs)...)
		} else {
			ok, _ := filepath.Match(suffix, file.Name())
			if ok {
				pathName := path.Join(dir, file.Name())
				list := strings.Split(filepath.ToSlash(pathName), "/")
				if len(list) > 1 {
					list = list[1:]
				}
				// fmt.Println(pathName, list)
				name := strings.Join(list, "/")
				arr = append(arr, &tplFile{Name: name, Path: pathName})
			}
		}
	}

	return
}

func Funcs(manifestJSON []byte) template.FuncMap {
	return template.FuncMap{
		"hello": func() string {
			return "hello world by template funcs!"
		},
		"vite":        func(path string) template.HTML { return vite(path, manifestJSON) },
		"isActiveUrl": isActiveUrl,
	}
}

// vite 使用 vite 组织 js 和 css 文件
func vite(_path string, manifestJSON []byte) template.HTML {
	ext := path.Ext(_path)

	vite_url := global.GetConfig().GetString(constants.VITE_URL)
	vite_build_dir := global.GetConfig().GetString(constants.VITE_BUILD_DIR)
	debug := global.GetConfig().GetBool(constants.VITE_DEBUG)
	vite_assets_root := global.GetConfig().GetString(constants.VITE_ASSETS_ROOT)
	// log.Info(fmt.Sprintf("vite url: %s, vite build dir: %s", vite_url, vite_build_dir))

	if debug {
		if ext == ".js" {
			return template.HTML(`<script>/** Vite **/</script> <script type="module" src="` + (vite_url + _path) + `"></script> <script>/** Vite end **/</script>`)
		}
		if ext == ".css" || ext == ".scss" {
			return template.HTML(`<style>/** Vite **/</style> <link rel="stylesheet" href="` + (vite_url + _path) + `"> <style>/** Vite end **/</style>`)
		}
		return template.HTML(`unknown ext: ` + ext)

	}
	manifestData := make(map[string]map[string]any)

	if len(manifestJSON) > 0 {
		if err := json.Unmarshal(manifestJSON, &manifestData); err != nil {
			log.Error("embed fs parse manifest error: " + err.Error() + "\n" + string(manifestJSON))
			return template.HTML("parse manifest error: " + err.Error())
		}
	} else {
		manifestPath := path.Join(vite_build_dir, "manifest.json")
		manifestFile, err := os.Open(manifestPath)
		if err != nil {
			log.Error("open manifest error: " + err.Error() + "\n" + manifestPath)
			return template.HTML("open manifest error: " + err.Error())
		}
		defer manifestFile.Close()
		if err := json.NewDecoder(manifestFile).Decode(&manifestData); err != nil {
			log.Error("parse manifest error: " + err.Error())
			return template.HTML("parse manifest error: " + err.Error())
		}
	}

	if _, ok := manifestData[_path]; !ok {
		log.Error("path not found in manifest: " + _path)
		return template.HTML("path not found in manifest: " + _path)
	}

	if _, ok := manifestData[_path]["file"]; !ok {
		log.Error("file not found in manifest: " + _path)
		return template.HTML("file not found in manifest: " + _path)
	}

	if ext == ".js" {

		jsFile := manifestData[_path]["file"].(string)
		if !strings.HasSuffix(jsFile, ".js") {
			log.Error("file not found in manifest: " + _path)
			return template.HTML("file not found in manifest: " + _path)
		}

		cssFiles := manifestData[_path]["css"]
		if cssFiles != nil {
			for _, cssFile := range cssFiles.([]any) {
				if !strings.HasSuffix(cssFile.(string), ".css") {
					log.Error("file not found in manifest: " + _path)
					return template.HTML("file not found in manifest: " + _path)
				}
			}
		}

		_html := `<script>/** Vite **/</script> <script type="module" src="` + (vite_assets_root + jsFile) + `"></script>`

		if cssFiles != nil {
			for _, cssFile := range cssFiles.([]any) {
				_html += `<link rel="stylesheet" href="` + (vite_assets_root + cssFile.(string)) + `">`
			}
		}

		_html += ` <script>/** Vite end **/</script>`

		return template.HTML(_html)

	}

	if ext == ".css" || ext == ".scss" {
		cssFile := manifestData[_path]["file"].(string)
		if !strings.HasSuffix(cssFile, ".css") {
			log.Error("file not found in manifest: " + _path)
			return template.HTML("file not found in manifest: " + _path)
		}

		return template.HTML(`<style>/** Vite **/</style> <link rel="stylesheet" href="` + (vite_assets_root + cssFile) + `"> <style>/** Vite end **/</style>`)
	}

	return template.HTML("unknown ext: " + ext)
}

func isActiveUrl(path string, request *http.Request) bool {
	return request.URL.Path == path
}
