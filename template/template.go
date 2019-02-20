package template

import (
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"

	"github.com/jacobtomlinson/krontab/config"
	"github.com/jacobtomlinson/krontab/input"
)

var defaultCronYaml = `
apiVersion: batch/v1beta1
kind: CronJob
metadata:
  name: test
spec:
  schedule: 0 */6 * * *
  jobTemplate:
    metadata:
      name: test
      krontabTemplate: default
    spec:
      template:
        metadata:
        spec:
          restartPolicy: Never
          containers:
            - name: test
              image: busybox
              command: ["echo", "hello"]
              resources:
                limits:
                  cpu: "1"
                  memory: 2G
                requests:
                  cpu: "0.25"
                  memory: 0.5G
`

func contains(s []string, e string) bool {
	for _, a := range s {
		if a == e {
			return true
		}
	}
	return false
}

// ListTemplates gives a list of the current cron templates
func ListTemplates() []string {
	files, err := ioutil.ReadDir(filepath.Join(config.ConfigDir.Path, "templates"))
	if err != nil {
		panic(err)
	}

	var templates []string

	for _, f := range files {
		templates = append(templates, strings.Replace(f.Name(), ".yaml", "", 1))
	}
	return templates
}

// IsTemplate checks whether a template exists
func IsTemplate(template string) bool {
	templates := ListTemplates()
	return contains(templates, template)
}

// EditTemplate opens a template file for editing
func EditTemplate(template string) error {
	// TODO Check if template in use and fail if so
	if template == "default" {
		fmt.Println("You cannot edit the default template. Creating a new one instead.")
		return errors.New("you cannot edit the default template")
	}
	if IsTemplate(template) {
		path := filepath.Join(config.ConfigDir.Path, "templates", template+".yaml")
		input.UserEdit(path)
		// TODO Validate is valid CronJob template
	} else {
		fmt.Println(fmt.Sprintf("Template %s doesn't exist.", template))
		return errors.New("template doesn't exist")
	}
	return nil
}

// CreateTemplate opens a new template file for editing
func CreateTemplate(template string) error {
	if !IsTemplate(template) {
		path := filepath.Join(config.ConfigDir.Path, "templates", template+".yaml")
		input.UserEdit(path)
		// TODO Validate is valid CronJob template
	} else {
		fmt.Println(fmt.Sprintf("Template %s already exists.", template))
		return errors.New("template already exists")
	}
	return nil
}

// DeleteTemplate opens a new template file for editing
func DeleteTemplate(template string) error {
	// TODO Check if template in use and fail if so
	if template == "default" {
		fmt.Println("You cannot delete the default template.")
		return errors.New("you cannot delete the default template")
	}
	if IsTemplate(template) {
		path := filepath.Join(config.ConfigDir.Path, "templates", template+".yaml")
		os.Remove(path)
	} else {
		fmt.Println(fmt.Sprintf("Template %s doesn't exist.", template))
		return errors.New("template doesn't exist")
	}
	return nil
}

// GetTemplate opens a template and reads as a string
func GetTemplate(template string) (string, error) {
	if IsTemplate(template) {
		path := filepath.Join(config.ConfigDir.Path, "templates", template+".yaml")
		dat, err := ioutil.ReadFile(path)
		if err != nil {
			panic(err)
		}
		return string(dat), nil
	}
	fmt.Println(fmt.Sprintf("Template %s doesn't exist.", template))
	return "", errors.New("template doesn't exist")
}

func init() {
	config.ConfigDir.WriteFile(filepath.Join("templates", "default.yaml"), []byte(defaultCronYaml))
}
