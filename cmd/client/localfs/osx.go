// +build darwin

package localfs

import (
	"encoding/xml"
	"fmt"
	"os"
)

type PlistArray struct {
	String []string `xml:"string"`
}

var stdExclusion = make(map[string][]string)

const stdExclusionPath = "/System/Library/CoreServices/backupd.bundle/Contents/Resources/StdExclusions.plist"

func init() {
	xmlFile, err := os.Open(stdExclusionPath)
	if err != nil {
		fmt.Println("Error opening file:", err)
		return
	}
	defer xmlFile.Close()

	dec := xml.NewDecoder(xmlFile)
	var workingKey string

	for {
		token, _ := dec.Token()
		if token == nil {
			break
		}
		switch start := token.(type) {
		case xml.StartElement:
			switch start.Name.Local {
			case "key":
				var k string
				err := dec.DecodeElement(&k, &start)
				if err != nil {
					fmt.Println(err.Error())
				}
				workingKey = k
			case "array":
				var ai PlistArray
				err := dec.DecodeElement(&ai, &start)
				if err != nil {
					fmt.Println(err.Error())
				}
				stdExclusion[workingKey] = ai.String
				workingKey = ""
			}
		}
	}

	fmt.Println(stdExclusion)
}

func ignoreByStd(dirname string) bool {
	return contains(stdExclusion["PathsExcluded"], dirname) ||
		contains(stdExclusion["ContentsExcluded"], dirname) ||
		contains(stdExclusion["FileContentsExcluded"], dirname)
}
