package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"regexp"
	"time"

	"github.com/hculpan/theearthisflatnomic/entity"
	"github.com/hculpan/theearthisflatnomic/utils"
	"github.com/hculpan/theearthisflatnomic/web"
)

//go:generate stringer -type=NextTurnAction ./entity

// Config contains the runtime configuration
// options for the application
type Config struct {
	Refresh bool
	Watch   bool
}

var filesWatched map[string]os.FileInfo = map[string]os.FileInfo{}

func main() {
	var config *Config
	if c, err := processCommandLine(); err != nil {
		fmt.Println(err)
		return
	} else {
		config = c
	}

	if config.Refresh {
		setupTestData()
	}

	if config.Watch {
		setupWatcher()
	}

	if os.Getenv("TEIFN_SECRET_KEY") == "" {
		fmt.Println("TEIFN_SECRET_KEY is not setup")
		return
	}

	initialize()

	fmt.Println("Server is started and listening on port 8080")
	web.SetupRoutes()
	log.Fatal(http.ListenAndServe(":8080", nil))
}

func setupTestData() {
	fmt.Println("Setting up test db")
	if err := utils.RemoveContents("./dbdata"); err != nil {
		panic(err)
	}
	initialize()
	fmt.Println("Loading user")
	entity.AddNewUser("Harry Culpan", "Kabluey", "harry@culpan.org", "happy")
	entity.AddNewUser("Harry Culpan 2", "Sarsenon", "harry@tesseract.org", "happy")
	entity.AddNewUser("Chris Duignan", "ElBriano", "duignan_chris@yahoo.com", "happy")
	if _, err := entity.SetFirstTurn(); err != nil {
		panic(err)
	}
	if _, err := entity.SetNextTurn(); err != nil {
		panic(err)
	}

	if err := loadRuleset(); err != nil {
		panic(err)
	}
}

func loadRuleset() error {
	fmt.Println("Loading initial ruleset")
	f, err := os.OpenFile("./resources/data/initial-rules.json", os.O_RDONLY, 0755)
	if err != nil {
		return err
	}
	defer f.Close()
	dec := json.NewDecoder(f)
	var list []interface{}
	if err := dec.Decode(&list); err != nil {
		return err
	}
	for _, v := range list {
		e := v.(map[string]interface{})
		r := entity.Rule{
			RuleNumber: int(e["number"].(float64)),
			RuleText:   buildString(e["text"].([]interface{})),
			Initial:    true,
			Active:     true,
		}
		r.Insert()
	}

	return nil
}

func buildString(list []interface{}) string {
	result := ""

	for _, v := range list {
		result += v.(string)
	}

	return result
}

func initialize() error {
	folderinfo, err := os.Stat("./dbdata")
	if err != nil && os.IsNotExist(err) {
		if err = os.Mkdir("./dbdata", 0755); err != nil {
			return err
		}
	} else if err != nil {
		return err
	} else if !folderinfo.IsDir() {
		return fmt.Errorf("dbdata exists, but is not a directory")
	}

	if err := entity.InitDB(); err != nil {
		fmt.Println(err)
		return err
	}

	return nil
}

func processCommandLine() (*Config, error) {
	result := &Config{}

	for _, v := range os.Args[1:] {
		switch v {
		case "--refresh":
			result.Refresh = true
		case "--watch":
			result.Watch = true
		default:
			return nil, fmt.Errorf("Unrecognized parameter: '%s'", v)
		}
	}

	return result, nil
}

func setupWatcher() {
	fmt.Println("Watching files")
	filesWatched = loadFiles()

	go func() {
		for {
			if changed, currFiles := filesHaveChanged(); changed {
				fmt.Print("Template files have changed, reloading...")
				web.LoadTemplates()
				filesWatched = currFiles
				fmt.Println("done")
			}
			time.Sleep(time.Second * 3)
		}
	}()
}

func loadFiles() map[string]os.FileInfo {
	result := map[string]os.FileInfo{}

	match := regexp.MustCompile(".gohtml$")
	dir, _ := os.Open("./templates")
	defer dir.Close()
	files, _ := dir.ReadDir(0)
	for _, f := range files {
		if match.MatchString(f.Name()) {
			i, _ := f.Info()
			result[f.Name()] = i
		}
	}

	return result
}

func filesInfosAreEqual(f1, f2 os.FileInfo) bool {
	return f1.IsDir() == f2.IsDir() &&
		f1.ModTime() == f2.ModTime() &&
		f1.Size() == f2.Size()
}

func filesHaveChanged() (bool, map[string]os.FileInfo) {
	currFileWatch := loadFiles()
	if len(currFileWatch) != len(filesWatched) {
		return true, currFileWatch
	}

	for _, v1 := range currFileWatch {
		v2, ok := filesWatched[v1.Name()]
		if !ok {
			return true, currFileWatch
		} else if !filesInfosAreEqual(v1, v2) {
			return true, currFileWatch
		}
	}

	return false, currFileWatch
}
