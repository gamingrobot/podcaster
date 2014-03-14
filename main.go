package main

import (
	"bytes"
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"path"
	"strconv"
	"strings"
	"text/template"
	"time"
)

type Show struct {
	Title  string
	Date   string
	Url    string
	Length int64
}

type Podcast struct {
	Shows []Show
}

func getUrls(file string) map[string]string {
	var urls map[string]string
	urlfile, err := ioutil.ReadFile(file)
	if err != nil {
		panic("Could not find url file")
	}
	err = json.Unmarshal(urlfile, &urls)
	return urls
}

func formatXml(file string, object Podcast) string {
	xfile, err := ioutil.ReadFile(file)
	if err != nil {
		panic("Could not find template file")
	}
	t := template.New("template")
	t, _ = t.Parse(string(xfile))
	var xmltemp bytes.Buffer
	t.Execute(&xmltemp, object)
	return xmltemp.String()
}

func main() {
	baseurl := flag.String("url", "localhost/", "Base url for rss")
	configfile := flag.String("streams", "urls.json", "Stream Urls")
	rssfile := flag.String("rss", "rss.xml", "Rss feed output file")
	downloaddir := flag.String("downloads", "downloads/", "Downloads directory")
	script := flag.String("script", "./downloader.sh", "Script used to download it will be passed url, time, outputfile")
	timeout := flag.Duration("time", time.Duration(60)*time.Minute, "Time duration")
	flag.Parse()
	//json
	urls := getUrls(*configfile)
	//template
	x := Podcast{}
	for key, url := range urls {
		t := time.Now()
		//format outputfile
		name := strings.ToLower(key)
		name = strings.Replace(name, " ", "_", -1)
		year, tmonth, day := t.Date()
		month := strings.ToLower(tmonth.String())
		filename := name + "_" + month + "_" + strconv.Itoa(day) + "_" + strconv.Itoa(year) + ".mp4"
		outputpath := path.Clean(*downloaddir + "/" + filename)
		fmt.Println(outputpath)
		//execute function
		cmd := exec.Command(*script, url, strconv.FormatFloat(timeout.Seconds(), 'f', 6, 64), outputpath)
		err := cmd.Start()
		if err != nil {
			panic(err)
		}
		fmt.Println("Waiting for stream to finish")
		err = cmd.Wait()
		//get size of file
		file, err := os.Open(outputpath)
		if err != nil {
			panic(err)
		}
		stat, _ := file.Stat()
		x.Shows = append(x.Shows, Show{Title: key, Date: t.Format(time.RFC1123Z), Url: *baseurl + filename, Length: stat.Size()})

	}
	xmlout := formatXml("template.xml", x)
	f, err := os.Create(*rssfile)
	if err != nil {
		panic(err)
	}
	f.WriteString(xmlout)
	f.Sync()
}
