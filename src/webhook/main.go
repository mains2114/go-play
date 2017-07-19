package main

import (
	"fmt"
	"net/http"
	"os/exec"
	"bytes"
	"time"
	"github.com/olebedev/config"
	"strconv"
	"errors"
)

const configFile = "./src/webhook/config.yaml"

type Project struct {
	key string
	dir string
	cmd []string
}

func main() {
	projects, err := initProject()
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(projects)
	//return
	//dir, err := projects.String("ims.dir")
	//fmt.Println(dir)
	//return

	cnt := 0
	http.HandleFunc("/", func (res http.ResponseWriter, req *http.Request) {
		cnt++
		fmt.Printf("#%d\t%s\t%s\n", cnt, req.RequestURI, req.RemoteAddr)
		fmt.Fprintf(res, "uri: %s", req.URL.Path)

		cmd := exec.Command("git", "pull")
		cmd.Dir = "D:/www/ims"

		var out bytes.Buffer
		cmd.Stdout = &out

		fmt.Println(time.Now())
		time.Sleep(5 * 1000 * 1000 * 1000)
		fmt.Println(time.Now())

		cmd.Run()
		fmt.Println(out.String())
	})
	fmt.Println("listening on 0.0.0.0:8080")
	http.ListenAndServe("0.0.0.0:8080", nil)
}

func initProject() (projects map[string]Project, err error) {
	cfg, err := config.ParseYamlFile(configFile)
	if err != nil {
		return
	}

	values, err := cfg.Map("projects")
	if err != nil {
		return
	}
	if len(values) <= 0 {
		return nil, errors.New(configFile + " has no valid project config")
	}


	projects = make(map[string]Project)

	for k, _ := range values {
		projectCfg, _ := cfg.Get("projects." + k)
		dir, _ := projectCfg.String("dir")
		cmd, _ := projectCfg.List("cmd")
		cmd2 := make([]string, 0)
		for i := 0; i < len(cmd); i++ {
			s, _ := projectCfg.String("cmd." + strconv.Itoa(i))
			cmd2 = append(cmd2, s)
		}

		projects[k] = Project{
			key: k,
			dir: dir,
			cmd: cmd2,
		}
		//fmt.Println(projects[k])
	}

	return projects, err
}
