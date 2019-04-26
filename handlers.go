package main

import (
	"encoding/json"
	"fmt"
	"github.com/google/uuid"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"path"
	"path/filepath"
)

func getTablesHandler(w http.ResponseWriter, r *http.Request) {

	a := JsonTableOku(getProject(r.Header.Get("projectId")).Path)

	_ = json.NewEncoder(w).Encode(a)
}

func getProxyClassHandler(w http.ResponseWriter, r *http.Request) {
	a := JsonProxyClassOku(getProject(r.Header.Get("projectId")).Path)
	_ = json.NewEncoder(w).Encode(a)
}

func getProxyClassByNameHandler(w http.ResponseWriter, r *http.Request) {

	a := JsonProxyClassOku(getProject(r.Header.Get("projectId")).Path)
	name1 := r.Header.Get("className1")
	name2 := r.Header.Get("className2")

	var pc [2]TProxyClass
	for i := 0; i < len(a); i++ {
		if a[i].Name == name1 {
			pc[0] = a[i]
		} else if a[i].Name == name2 {
			pc[1] = a[i]
		}
	}

	_ = json.NewEncoder(w).Encode(pc)
}

func getEndPointHandler(w http.ResponseWriter, r *http.Request) {

	a := JsonEndPointOku(getProject(r.Header.Get("projectId")).Path)

	_ = json.NewEncoder(w).Encode(a)
}

func saveEndPointHandler(w http.ResponseWriter, r *http.Request) {

	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		panic(err)
	}
	var model []TEndPoint
	err = json.Unmarshal(body, &model)
	if err != nil {
		panic(err)
	}

	JsonEndPointKaydet(model, getProject(r.Header.Get("projectId")).Path)

	_ = json.NewEncoder(w).Encode(model)

}

func saveTablesHandler(w http.ResponseWriter, r *http.Request) {

	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		panic(err)
	}
	var project []TDataTable
	err = json.Unmarshal(body, &project)
	if err != nil {
		panic(err)
	}

	for i := 0; i < len(project); i++ {
		if project[i].Uid == "" {
			project[i].Uid = uuid.New().String()
		}
	}

	JsonTableKaydet(project, getProject(r.Header.Get("projectId")).Path)

	_ = json.NewEncoder(w).Encode(project)

}

func saveProxyClassHandler(w http.ResponseWriter, r *http.Request) {

	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		panic(err)
	}
	var model []TProxyClass
	err = json.Unmarshal(body, &model)
	if err != nil {
		panic(err)
	}

	JsonProxyClassKaydet(model, getProject(r.Header.Get("projectId")).Path)

	_ = json.NewEncoder(w).Encode(model)

}

func getProject(uid string) TProject {

	var r = TProject{}
	projects := JsonProjeOku()
	for i := 0; i < len(projects); i++ {
		if projects[i].Uid == uid {
			r = projects[i]
			break
		}
	}
	return r
}

func prgFormat(path string, w http.ResponseWriter) {

	cmd := "go fmt " + path + "/*.go"

	_, _ = fmt.Fprintf(w, "$: "+cmd+"\n")
	err, out, errout := Shellout(path, "bash", "-c", cmd)
	if err != nil {
		_, _ = fmt.Fprintf(w, "error: %v\n", err)
	}

	if len(out) > 0 {
		_, _ = fmt.Fprintf(w, out)
	}
	if len(errout) > 0 {
		_, _ = fmt.Fprintf(w, errout)

	}

}

func prgBuild(path string, w http.ResponseWriter) {

	cmd := "go build "

	_, _ = fmt.Fprintf(w, "$: "+cmd+"\n")
	//err, out, errout := Shellout(path,"go", "build", path+"/main.go")
	err, out, errout := Shellout(path, "go", "build")

	if err != nil {
		_, _ = fmt.Fprintf(w, "error: %v\n", err)
	}

	if len(out) > 0 {
		fmt.Fprintf(w, out)
	}
	if len(errout) > 0 {
		fmt.Fprintf(w, errout)
	}

}

func buildHandler(w http.ResponseWriter, r *http.Request) {

	// Kill it:

	if c.Process != nil {
		if err := c.Process.Kill(); err != nil {
			log.Fatal("failed to kill process: ", err)
		}
	}
	projectId := r.Header.Get("projectId")
	project := getProject(projectId)
	PrgDir = project.Path

	os.MkdirAll(project.Path+"/gocreator/db", os.ModePerm)

	tables := JsonTableOku(PrgDir)
	proxyclass := JsonProxyClassOku(PrgDir)
	endpoint := JsonEndPointOku(PrgDir)

	TamplateFile := "main.gohtml"
	HedefeKaydet(path.Base(project.Path), (project.Path + "/main.go"), ("./templates/" + TamplateFile), TamplateFile)

	TamplateFile = "InitDB_oto.gohtml"
	HedefeKaydet(tables, (project.Path + "/gocreator/InitDB.go"), ("./templates/" + TamplateFile), TamplateFile)

	TamplateFile = "entity_oto.gohtml"
	HedefeKaydet(tables, (project.Path + "/gocreator/" + "entity_oto.go"), ("./templates/" + TamplateFile), TamplateFile)

	TamplateFile = "crud_oto.gohtml"
	HedefeKaydet(tables, (project.Path + "/gocreator/" + "crud_oto.go"), ("./templates/" + TamplateFile), TamplateFile)

	TamplateFile = "proxyclass_oto.gohtml"
	HedefeKaydet(proxyclass, (project.Path + "/gocreator/" + "proxyclass_oto.go"), ("./templates/" + TamplateFile), TamplateFile)

	TamplateFile = "handler_oto.gohtml"
	HedefeKaydet(endpoint, (project.Path + "/gocreator/" + "handler_oto.go"), ("./templates/" + TamplateFile), TamplateFile)

	for i := 0; i < len(endpoint); i++ {
		TamplateFile = "handlerMap.gohtml"
		HedefFileName := project.Path + "/gocreator/" + "handlerMap_" + endpoint[i].Name + ".go"
		HedefeKaydetEgerDosyaYoksa(endpoint[i], HedefFileName, ("./templates/" + TamplateFile), TamplateFile)

	}

	prgFormat(project.Path+"/gocreator", w)
	prgFormat(project.Path, w)

	prgBuild(project.Path, w)

	//json.NewEncoder(w).Encode(project)

}

func getProjecthandler(w http.ResponseWriter, r *http.Request) {
	goprojects := JsonProjeOku()
	json.NewEncoder(w).Encode(goprojects)
}

func saveProjectHandler(w http.ResponseWriter, r *http.Request) {

	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		panic(err)
	}
	var project TProject
	err = json.Unmarshal(body, &project)
	if err != nil {
		panic(err)
	}

	project.Uid = uuid.New().String()
	//project.Ad= "Deneme"

	goprojects := JsonProjeOku()
	goprojects = append(goprojects, project)
	JsonProjeKaydet(goprojects)

	json.NewEncoder(w).Encode(goprojects)

}

func getDir(w http.ResponseWriter, r *http.Request) {

	projectId := r.Header.Get("projectId")
	project := getProject(projectId)

	a, _ := NewTree(project.Path)

	json.NewEncoder(w).Encode(a.Children)

}

type TFile struct {
	Path string
}

func getFile(w http.ResponseWriter, r *http.Request) {

	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		panic(err)
	}
	var request TFile
	err = json.Unmarshal(body, &request)
	if err != nil {
		panic(err)
	}

	buf, _ := ioutil.ReadFile(request.Path)

	json.NewEncoder(w).Encode(string(buf))

}

type TFileSave struct {
	Path  string
	Value string
}

func saveFile(w http.ResponseWriter, r *http.Request) {

	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		panic(err)
	}
	var request TFileSave
	err = json.Unmarshal(body, &request)
	if err != nil {
		panic(err)
	}

	ioutil.WriteFile(request.Path, []byte(request.Value), 0644)

	prgFormat2(request.Path)

	buf, _ := ioutil.ReadFile(request.Path)
	json.NewEncoder(w).Encode(string(buf))

}

func prgFormat2(path string) {

	cmd := "go fmt " + path

	s := filepath.Dir(path)

	err, out, errout := Shellout(s, "bash", "-c", cmd)
	if err != nil {
		fmt.Print("err:")
		fmt.Println(err)
	}

	if len(out) > 0 {
		fmt.Println("out:" + out)
	}
	if len(errout) > 0 {
		fmt.Println("errout:" + errout)

	}

}
