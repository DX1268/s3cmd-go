package controllers

type FileObject struct {
	Source   int64 // used by sync
	Name     string
	Size     int64
	Checksum string
}

const DATE_FMT = "2006-01-02 15:04"
