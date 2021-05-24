package entities

import "os"

type Course struct {
    Title string
    Items []Module
}

type Module struct {
    Title     string
    Files []*os.File
}
