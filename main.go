package main
import (
	"bufio"
	"embed"
	"fmt"
	"os"
	"os/user"
	"path/filepath"
	"time"
	"strconv"
	"strings"
	"io/ioutil"
	"net/http"

	"github.com/sqweek/dialog"
	"github.com/gen2brain/beeep"
	"gopkg.in/yaml.v3"
)
//go:embed assets/*
var embedDirStatic embed.FS
func checkError(err error) {
	if err != nil {
		fmt.Println(err)
		panic(err)
	}
}
func showDialog() (bool) {
	return dialog.Message("%s", "Are you working?").Title("workreg").YesNo()
}
func showFistRun() {
	dialog.Message("%s", "check .workreg/config.yaml in your homefolder, and make a config").Title("workreg").Info()
}
func showToast(location string) {
	fileContent, err := embedDirStatic.ReadFile("assets/lego_headphone.png")
	checkError(err)
	tempFile, err := ioutil.TempFile("", "lego_headphone.png")
	checkError(err)
	defer os.Remove(tempFile.Name())
	_, err = tempFile.Write(fileContent)
	checkError(err)
	err = beeep.Alert("Workreg", "Workday registered at " + location, tempFile.Name())
	checkError(err)
}
func getWorkregDir() (string) {
	currentUser, err := user.Current()
	checkError(err)
	return filepath.Join(currentUser.HomeDir, ".workreg")
}
func getIp() (string, error) {
	response, err := http.Get("https://api.ipify.org/")
	if err != nil {
		return "", err
	}
	defer response.Body.Close()
	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return "", err
	}
	ip := string(body)
	return ip, nil
}
func readLastLine(filePath string) (string, error) {
	file, err := os.Open(filePath)
	checkError(err)
	defer file.Close()
	var lastLine string
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		lastLine = scanner.Text()
	}
	err = scanner.Err()
	checkError(err)
	return lastLine, nil
}
func isDayRegistered()(bool) {
	current := time.Now()
	year := strconv.Itoa(current.Year())
	month := current.Month().String()
	day := strconv.Itoa(current.Day())

	path := filepath.Join(getWorkregDir(), year)
	filename := filepath.Join(path, month + ".txt")
	line, _ := readLastLine(filename);
	return strings.HasPrefix(line, day + " ")
}
func writeFile(location string)(bool) {
	current := time.Now()
	year := strconv.Itoa(current.Year())
	month := current.Month().String()
	day := strconv.Itoa(current.Day())
	
	path := filepath.Join(getWorkregDir(), year)
	err := os.MkdirAll(path, 0755)
	checkError(err)
	filename := filepath.Join(path, month + ".txt")
	file, err := os.OpenFile(filename, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	checkError(err)
	defer file.Close()
	_, err = file.WriteString(day + " " + location + "\r\n")
	checkError(err)
	return true
}
func writeAndNotify(location string) {
	isNewDay := writeFile(location)
	if isNewDay {
		showToast(location)
	}
}
func getLocation(filename string)(string, bool) {
	yamlContent, err := ioutil.ReadFile(filename)
	checkError(err) 
	type Entry struct {
		IP  string `yaml:"ip"`
		Ask bool   `yaml:"ask"`
	}
	type Config struct {
		Entries map[string]Entry `yaml:",inline"`
	}
	var config Config
	err = yaml.Unmarshal(yamlContent, &config)
	checkError(err) 
	ip, err := getIp()
	checkError(err)
	for location, entry := range config.Entries {
		// fmt.Printf("%s IP: %s, Ask: %v\n", location, entry.IP, entry.Ask)
		if(entry.IP == ip) {
			return location, entry.Ask;
		}
	}
	return "other", true;
}
func loadConfig() {
	filename := filepath.Join(getWorkregDir(), "config.yaml")
	// fmt.Println(filename);
	if _, err := os.Stat(filename); os.IsNotExist(err) {
		fileContent, _ := embedDirStatic.ReadFile("assets/config.yaml")
		ioutil.WriteFile(filename, fileContent, 0644)
		showFistRun();
		return;
	} 
	location, ask := getLocation(filename)
	if ask && !showDialog() {
		return;
	}
	writeAndNotify(location)
	// file, err := os.OpenFile(filename, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
}
func main() {
	if(!isDayRegistered()) {
		loadConfig();
	}
}