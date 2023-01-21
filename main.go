// spserve - Serve files to current network with ease.
// Copyright (C) 2023 Safal Piya

// This program is free software: you can redistribute it and/or modify
// it under the terms of the GNU General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.

// This program is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
// GNU General Public License for more details.

// You should have received a copy of the GNU General Public License
// along with this program.  If not, see <https://www.gnu.org/licenses/>.

package main

import (
	"embed"
	"errors"
	"flag"
	"fmt"
	"html/template"
	"log"
	"net"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
)

// Server variables
var (
	serverIP   string
	serverPort int
	rootPath   string
)

// Templates

//go:embed templates/*.html
var tmpl embed.FS

var (
	dirListTmpl *template.Template
	errTmpl     *template.Template
)

type Response struct {
	StatusCode    int
	StatusMessage string
	Message       string
}

type DirEntry struct {
	Name     string
	PrevDirs []FileEntry
	Dirs     []FileEntry
	RegFiles []FileEntry
	IsRoot   bool
}

type FileEntry struct {
	Name     string
	Location string
	IsDir    bool
}

// Server
// ------

func getMyInterfaceAddr() (net.IP, error) {
	ifaces, err := net.Interfaces()
	if err != nil {
		return nil, err
	}
	addresses := []net.IP{}
	for _, iface := range ifaces {

		if iface.Flags&net.FlagUp == 0 {
			continue // interface down
		}
		if iface.Flags&net.FlagLoopback != 0 {
			continue // loopback interface
		}
		addrs, err := iface.Addrs()
		if err != nil {
			continue
		}

		for _, addr := range addrs {
			var ip net.IP
			switch v := addr.(type) {
			case *net.IPNet:
				ip = v.IP
			case *net.IPAddr:
				ip = v.IP
			}
			if ip == nil || ip.IsLoopback() {
				continue
			}
			ip = ip.To4()
			if ip == nil {
				continue // not an ipv4 address
			}
			addresses = append(addresses, ip)
		}
	}
	if len(addresses) == 0 {
		return nil, fmt.Errorf("no address Found, net.InterfaceAddrs: %v", addresses)
	}
	return addresses[0], nil
}

func serveFile(w http.ResponseWriter, r *http.Request) {
	reqFilePath := filepath.Join(rootPath, r.URL.Path)
	log.Println(reqFilePath)

	info, err := os.Open(reqFilePath)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			errTmpl.Execute(w, Response{
				StatusCode:    http.StatusNotFound,
				StatusMessage: http.StatusText(http.StatusNotFound),
				Message:       fmt.Sprintf("Given file '%s' doesn't exist.", reqFilePath),
			})
			return
		}
		errTmpl.Execute(w, Response{
			StatusCode:    http.StatusInternalServerError,
			StatusMessage: http.StatusText(http.StatusInternalServerError),
			Message:       err.Error(),
		})
		return
	}
	stat, err := info.Stat()
	if err != nil {
		errTmpl.Execute(w, Response{
			StatusCode:    http.StatusInternalServerError,
			StatusMessage: http.StatusText(http.StatusInternalServerError),
			Message:       err.Error(),
		})
		return
	}

	// If the file requested is a regular file
	if !stat.IsDir() {
		http.ServeFile(w, r, reqFilePath)
		return
	}

	// If the file requested is a directory
	dirEntry := DirEntry{
		Name:   stat.Name(),
		IsRoot: r.URL.Path == "/",
	}
	files, err := info.Readdir(0)
	if err != nil {
		errTmpl.Execute(w, Response{
			StatusCode:    http.StatusInternalServerError,
			StatusMessage: http.StatusText(http.StatusInternalServerError),
			Message:       err.Error(),
		})
		return
	}
	for _, file := range files {
		fileEntry := FileEntry{
			Name:     file.Name(),
			IsDir:    file.IsDir(),
			Location: filepath.Join(r.URL.Path, file.Name()),
		}
		if file.IsDir() {
			dirEntry.Dirs = append(dirEntry.Dirs, fileEntry)
		} else {
			dirEntry.RegFiles = append(dirEntry.RegFiles, fileEntry)
		}
	}
	dirListTmpl.Execute(w, dirEntry)
}

func setHTTPHandlers() {
	http.HandleFunc("/", serveFile)
}

func startServer(port int) {
	log.Printf("Serving \"%s\" in %s:%s", rootPath, serverIP, strconv.Itoa(serverPort))
	err := http.ListenAndServe(":"+strconv.Itoa(port), nil)
	if err != nil {
		log.Fatalf("[ERROR] Couldn't start server: %s", err)
	}
}

// Core
// ----

func usage() {
	fmt.Printf("Usage: %s [options] root_dir\n\noptions:\n", os.Args[0])
	flag.PrintDefaults()
}

func parseArgs() {
	flag.IntVar(&serverPort, "port", 8080, "Port for the server to listen")
	flag.Parse()
}

func getRootPathCleaned(root string) (string, error) {
	path, err := filepath.Abs(root)
	if err != nil {
		return "", err
	}

	info, err := os.Stat(path)
	if err != nil {
		return "", err
	}
	if !info.IsDir() {
		return "", fmt.Errorf("Given path is not an directory")
	}

	return path, nil
}

func initVariables() {
	parseArgs()

	nonFlagArgs := flag.Args()
	if len(nonFlagArgs) != 1 {
		usage()
		os.Exit(1)
	}

	var err error
	rootPath, err = getRootPathCleaned(nonFlagArgs[0])
	if err != nil {
		log.Fatalf("[ERROR] Couldn't set root path: %s", err)
	}

	interfaceAddr, err := getMyInterfaceAddr()
	if err != nil {
		log.Fatalf("[ERROR] Couldn't get network interface address: %s", err)
	}
	serverIP = interfaceAddr.String()
}

func executeTemplates() {
	var err error
	dirListTmpl, err = template.ParseFS(tmpl, "templates/directory_listing.html")
	if err != nil {
		log.Printf("[ERROR] Couldn't parse the directory listing template: %s", err)
		os.Exit(1)
	}

	errTmpl, err = template.ParseFS(tmpl, "templates/error.html")
	if err != nil {
		log.Printf("[ERROR] Couldn't parse the error template: %s", err)
		os.Exit(1)
	}
}

func main() {
	initVariables()
	executeTemplates()
	setHTTPHandlers()
	startServer(serverPort)
}
