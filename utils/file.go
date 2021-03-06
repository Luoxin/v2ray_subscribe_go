package utils

import (
	"io"
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/eddieivan01/nic"
	"github.com/schollz/progressbar/v3"
	log "github.com/sirupsen/logrus"
)

func IsDir(path string) bool {
	s, err := os.Stat(path)
	if err != nil {
		return false
	}
	return s.IsDir()
}

func FileExists(path string) bool {
	_, err := os.Stat(path)
	if err != nil {
		return os.IsExist(err)
	}
	return true
}

func IsFile(path string) bool {
	return !IsDir(path)
}

func CopyFile(from, to string) error {
	input, err := ioutil.ReadFile(from)
	if err != nil {
		return err
	}

	err = ioutil.WriteFile(to, input, 0777)
	if err != nil {
		return err
	}

	return nil
}

func GetConfigDir() string {
	homeDir, _ := os.UserHomeDir()
	return filepath.Join(homeDir, ".config", "Eutamias")
}

func FileWrite(filename string, content string) error {
	out, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer out.Close()

	_, err = out.WriteString(content)
	if err != nil {
		return err
	}
	return nil
}

func FileRead(filename string) (content string, err error) {
	buf, err := ioutil.ReadFile(filename)
	if err != nil {
		log.Errorf("read file %s err %s", filename, err)
		return "", err
	}
	return string(buf), err
}

func DownloadWithProgressbar(fileUrl, fileName string) error {
	resp, err := nic.Get(fileUrl, nic.H{
		AllowRedirect: true,
		Timeout:       300,
		Chunked:       true,
		SkipVerifyTLS: true,
	})
	if err != nil {
		log.Errorf("err:%v", err)
		return err
	}
	defer resp.Body.Close()

	f, err := os.OpenFile(fileName, os.O_CREATE|os.O_WRONLY, 0777)
	if err != nil {
		log.Errorf("err:%v", err)
		return err
	}
	defer f.Close()

	bar := progressbar.DefaultBytes(
		resp.ContentLength,
		"downloading",
	)

	_, err = io.Copy(io.MultiWriter(f, bar), resp.Body)
	if err != nil {
		log.Errorf("err:%v", err)
		return err
	}

	return nil
}

func GetExecPath() string {
	execPath, err := os.Executable()
	if err != nil {
		log.Errorf("err:%v", err)
		return ""
	}
	return filepath.Dir(execPath)
}

func GetPwd() string {
	pwd, err := os.Getwd()
	if err != nil {
		log.Errorf("err:%v", err)
		return ""
	}
	return pwd
}
