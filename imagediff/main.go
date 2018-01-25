// +build main

package main

import (
	"github.com/strickyak/web-cam-web/imagediff"

	"os"
)

func main() {
	imagediff.DiffFilenames(os.Args[1], os.Args[2], os.Args[3])
}
