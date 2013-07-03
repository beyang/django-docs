package django_docs

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"regexp"
	"strings"
)

type SymbolID string

type DjangoDoc struct {
	Symbol     SymbolID
	Body       string
	SourceFile string
	Start      int
	End        int
}

type moduleInfo struct {
	Name  string
	Start int
}

var moduleAnchorRegex = regexp.MustCompile("\\.\\. module:: (?P<moduleName>.*)")

func ExtractDocs(docDir string) (docs map[SymbolID]DjangoDoc, errs []error) {

	docs = make(map[SymbolID]DjangoDoc)

	refsDir := filepath.Join(docDir, "ref")
	filepath.Walk(refsDir, func(path string, info os.FileInfo, err error) error {
		if !info.IsDir() && filepath.Ext(path) == ".txt" {
			if docBytes, err := ioutil.ReadFile(path); err != nil {
				errs = append(errs, err)
			} else {
				docBody := string(docBytes)
				submatches := moduleAnchorRegex.FindAllStringSubmatchIndex(docBody, -1)

				mInfos := make([]moduleInfo, len(submatches))
				for s, submatch := range submatches {
					anchorIdx := submatch[0]
					mInfos[s].Name = docBody[submatch[2]:submatch[3]]
					mInfos[s].Start = findModuleStart(docBody, anchorIdx)
				}

				for m, mInfo := range mInfos {
					var end int
					if m+1 < len(mInfos) {
						end = mInfos[m+1].Start
					} else {
						end = len(docBody)
					}

					if srcFile, err := filepath.Rel(docDir, path); err != nil {
						errs = append(errs, err)
					} else {
						symbolId := SymbolID(mInfo.Name)
						if prev, ok := docs[symbolId]; ok {
							errs = append(errs,
								fmt.Errorf("Duplicate symbol: %s.  Previous location: %s:%d:%d.  New location: %s:%d:%d",
									symbolId, prev.SourceFile, prev.Start, prev.End, srcFile, mInfo.Start, end),
							)
						}
						docs[symbolId] = DjangoDoc{
							Symbol:     symbolId,
							Body:       docBody[mInfo.Start:end],
							SourceFile: srcFile,
							Start:      mInfo.Start,
							End:        end,
						}
					}
				}
			}
		}
		return nil
	})

	return
}

func findModuleStart(docBody string, moduleAnchorIdx int) int {
	nLinesUp := 4
	idx := moduleAnchorIdx
	for i := 0; i < nLinesUp; i += 1 {
		prefix := docBody[:idx]
		idx = strings.LastIndex(prefix, "\n")
		if idx < 0 {
			idx = 0
			break
		}
	}

	return idx
}
