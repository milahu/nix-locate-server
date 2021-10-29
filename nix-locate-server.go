// TODO nix-locate is rather slow -> make it faster with SQLite and full text search index
// TODO remove -d and --db from queryArgs

package main

import (
	"github.com/gin-gonic/gin"
	"github.com/google/shlex"

	"net/http"
	"strings"
	"fmt"
	"log"
	"os/exec"
	"regexp"
	"net/url"
	"bytes"
)

func main() {

	router := gin.Default()

	regexNonAscii := regexp.MustCompile("[[:^ascii:]]")

	router.GET("/", func(context *gin.Context) {

		context.Redirect(302, "/nix-locate.txt?--help") // moved temporarily
	})

	router.GET("/nix-locate.txt", func(context *gin.Context) {

		//queryRaw := context.Query("q")

		// use full query string
		queryRawEncoded := context.Request.URL.RawQuery // https://pkg.go.dev/net/url#URL
		//fmt.Printf("queryRawEncoded: %s\n", queryRawEncoded)

		queryRaw, err := url.PathUnescape(queryRawEncoded)
		if err != nil {
			log.Fatal(err)
		}
		//fmt.Printf("queryRaw: %s\n", queryRaw)

		query := strings.TrimSpace(regexNonAscii.ReplaceAllLiteralString(queryRaw, ""))

		if query == "" {
			//fmt.Printf("error: empty query\n")
			context.String(404, "error: empty query")
			return
		}
	
		if len(query) < 3 {
			//fmt.Printf("error: query is too short: %s\n", query) // WTF? %!(EXTRA string=)
			context.String(404, "error: query is too short: %s\n", query)
			return
		}

		fmt.Printf("query = '%s'\n", query)

		responseHeader := ""

		if query == "--help" {
			responseHeader = "$ echo have fun breaking my system :D >/dev/null\n\n"
		}

		//shlex.Split("one \"two three\" four") -> []string{"one", "two three", "four"}
		queryArgs, err := shlex.Split(query)
		if err != nil {
			fmt.Printf("error: failed to parse query: %s\n", query)
			context.String(404, "error: failed to parse query")
			//log.Fatal(err)
			return
		}
	
		// capture both stdout and stderr of subprocess
		// https://stackoverflow.com/a/39968254/10440128
		cmd := exec.Command("nix-locate", queryArgs...)
		var outb, errb bytes.Buffer
		cmd.Stdout = &outb
		cmd.Stderr = &errb
		//err := cmd.Run() // FIXME reassign err: no new variables on left side of :=
		err = cmd.Run()
		cmdOut := outb.String()
		cmdErr := errb.String()
		if err != nil {
			log.Printf("error: nix-locate returned nonzero\n", query)
			context.String(404, "$ nix-locate %s\n\n%s", query, cmdErr)
			return
		}
	
		//fmt.Printf("nix-locate result:\n%s\n", cmdOut)
		context.String(http.StatusOK, "%s$ nix-locate %s\n\n%s", responseHeader, query, cmdOut)

	})

	//router.Run() // serves on :8080 unless env.PORT is set
	router.Run("localhost:8080")

}
