//go:build ignore

// gen_versioninfo.go generates versioninfo.json for goversioninfo.
// Usage: go run ./cmd/git-worktree-manager/gen_versioninfo.go <version>
// Example: go run ./cmd/git-worktree-manager/gen_versioninfo.go 2.1.0
package main

import (
	"encoding/json"
	"fmt"
	"os"
	"strconv"
	"strings"
)

type FixedFileInfo struct {
	FileVersion    [4]uint16 `json:"FileVersion"`
	ProductVersion [4]uint16 `json:"ProductVersion"`
	FileFlagsMask  string    `json:"FileFlagsMask"`
	FileFlags      string    `json:"FileFlags"`
	FileOS         string    `json:"FileOS"`
	FileType       string    `json:"FileType"`
	FileSubType    string    `json:"FileSubType"`
}

type StringFileInfo struct {
	Comments         string `json:"Comments"`
	CompanyName      string `json:"CompanyName"`
	FileDescription  string `json:"FileDescription"`
	FileVersion      string `json:"FileVersion"`
	InternalName     string `json:"InternalName"`
	LegalCopyright   string `json:"LegalCopyright"`
	LegalTrademarks  string `json:"LegalTrademarks"`
	OriginalFilename string `json:"OriginalFilename"`
	PrivateBuild     string `json:"PrivateBuild"`
	ProductName      string `json:"ProductName"`
	ProductVersion   string `json:"ProductVersion"`
	SpecialBuild     string `json:"SpecialBuild"`
}

type Translation struct {
	LangID    string `json:"LangID"`
	CharsetID string `json:"CharsetID"`
}

type VarFileInfo struct {
	Translation Translation `json:"Translation"`
}

type VersionInfo struct {
	FixedFileInfo  FixedFileInfo  `json:"FixedFileInfo"`
	StringFileInfo StringFileInfo `json:"StringFileInfo"`
	VarFileInfo    VarFileInfo    `json:"VarFileInfo"`
	IconPath       string         `json:"IconPath"`
	ManifestPath   string         `json:"ManifestPath"`
}

func main() {
	if len(os.Args) < 2 {
		fmt.Fprintln(os.Stderr, "usage: gen_versioninfo.go <version>")
		os.Exit(1)
	}

	ver := strings.TrimPrefix(os.Args[1], "v")
	parts := strings.Split(ver, ".")
	if len(parts) < 3 {
		fmt.Fprintf(os.Stderr, "invalid version %q: expected Major.Minor.Patch\n", ver)
		os.Exit(1)
	}

	major, _ := strconv.ParseUint(parts[0], 10, 16)
	minor, _ := strconv.ParseUint(parts[1], 10, 16)
	patch, _ := strconv.ParseUint(parts[2], 10, 16)

	fileVer := fmt.Sprintf("%d.%d.%d.0", major, minor, patch)
	prodVer := fmt.Sprintf("%d.%d.%d", major, minor, patch)

	vi := VersionInfo{
		FixedFileInfo: FixedFileInfo{
			FileVersion:    [4]uint16{uint16(major), uint16(minor), uint16(patch), 0},
			ProductVersion: [4]uint16{uint16(major), uint16(minor), uint16(patch), 0},
			FileFlagsMask:  "3f",
			FileFlags:      "00",
			FileOS:         "040004",
			FileType:       "01",
			FileSubType:    "00",
		},
		StringFileInfo: StringFileInfo{
			CompanyName:      "Lucas Modrich",
			FileDescription:  "Git Worktree Manager",
			FileVersion:      fileVer,
			InternalName:     "gwtm",
			LegalCopyright:   "Copyright \u00a9 2026 Lucas Modrich",
			OriginalFilename: "gwtm.exe",
			ProductName:      "gwtm",
			ProductVersion:   prodVer,
		},
		VarFileInfo: VarFileInfo{
			Translation: Translation{
				LangID:    "0409",
				CharsetID: "04B0",
			},
		},
	}

	data, err := json.MarshalIndent(vi, "", "  ")
	if err != nil {
		fmt.Fprintln(os.Stderr, "json marshal error:", err)
		os.Exit(1)
	}

	outPath := "cmd/git-worktree-manager/versioninfo.json"
	if err := os.WriteFile(outPath, data, 0644); err != nil {
		fmt.Fprintln(os.Stderr, "write error:", err)
		os.Exit(1)
	}

	fmt.Println("wrote", outPath)
}
