package input

import (
	"io/ioutil"
	"os"
	"fmt"
	"os/exec"
)

func check(e error) {
    if e != nil {
        panic(e)
    }
}

// UserInput get raw input from a user with a text editor and a temporary file
func UserInput(in string) string {
	userInput := []byte(in)
	tmpDir := os.TempDir()
	tmpFile, tmpFileErr := ioutil.TempFile(tmpDir, "krontab_input_")
	if tmpFileErr != nil {
		fmt.Printf("Error %s while creating tempFile", tmpFileErr)
	}
	err := ioutil.WriteFile(tmpFile.Name(), userInput, 0644)
	check(err)

	return UserEdit(tmpFile.Name())
}

// UserEdit allow the user to edit a file
func UserEdit(path string) string {
	// TODO Check for EDITOR/VISUAL env var
	editor := "vim"
	editorPath, err := exec.LookPath(editor)
	if err != nil {
		fmt.Printf("Error %s while looking up for %s!!", editorPath, editor)
	}

	cmd := exec.Command(editorPath, path)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	err = cmd.Start()
	if err != nil {
		fmt.Printf("Start failed: %s", err)
	}
	err = cmd.Wait()

	dat, err := ioutil.ReadFile(path)
	check(err)
	return string(dat)
}

func init() {

}
