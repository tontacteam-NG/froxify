package binary

import (
	"os/exec"

	"github.com/kr/pretty"
)

const PATH = "./binary/"

func X8(url string) {
	cmd := exec.Command(PATH+"x8", url)
	err := cmd.Run()
	if err != nil {
		pretty.Println(err)
	}
}
