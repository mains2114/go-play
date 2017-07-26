package main

import (
	"bytes"
	"errors"
	"flag"
	"fmt"
	"github.com/olebedev/config"
	"net/http"
	"os/exec"
	"strconv"
	"strings"
	"time"
	//"os"
)

var configFile string

type Project struct {
	key string
	dir string
	cmd []string
}

func main() {
	flag.StringVar(&configFile, "c", "./src/webhook/config.yaml", "config file path")
	flag.Parse()
	fmt.Println("config file: " + configFile)

	addr, projects, err := getConfig()
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(addr, projects)
	//return
	//dir, err := projects.String("ims.dir")
	//fmt.Println(dir)
	//return

	cnt := 0
	http.HandleFunc("/", func(res http.ResponseWriter, req *http.Request) {
		cnt++
		fmt.Printf("%v\t#%d\t%s\t%s\n", time.Now(), cnt, req.RequestURI, req.RemoteAddr)
		defer fmt.Fprintf(res, "uri: %s", req.URL.Path)

		req.ParseForm()
		//fmt.Println(req.Form)

		project := strings.Join(req.Form["project"], "")
		_ = strings.Join(req.Form["action"], "")

		projectCfg, ok := projects[project]
		if !ok {
			res.WriteHeader(400)
			return
		}
		//fmt.Println(projectCfg)

		// 异步执行构建命令
		go func() {
			fmt.Println("project: ", project)
			fmt.Println("dir: ", projectCfg.dir)
			for _, v := range projectCfg.cmd {
				fmt.Println("execute: " + v)

				args := strings.Split(v, " ")
				cmd := exec.Command(args[0], args[1:]...)
				cmd.Dir = projectCfg.dir

				var out bytes.Buffer
				cmd.Stdout = &out
				cmd.Stderr = &out

				err = cmd.Run()
				if err != nil {
					fmt.Println(err)
					break
				} else {
					fmt.Println(out.String())
				}
			}
		}()
	})

	//addr := "127.0.0.1:8081"
	fmt.Println("listening on " + addr)
	err = http.ListenAndServe(addr, nil)
	if err != nil {
		fmt.Println(err)
	}
}

func getConfig() (addr string, projects map[string]Project, err error) {
	cfg, err := config.ParseYamlFile(configFile)
	if err != nil {
		return
	}

	addr, _ = cfg.String("bindAddr")
	values, err := cfg.Map("projects")
	if err != nil {
		return
	}
	if len(values) <= 0 {
		return "", nil, errors.New(configFile + " has no valid project config")
	}

	projects = make(map[string]Project)

	for k := range values {
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

	return addr, projects, err
}
