// Copyright (c) Qualcomm Technologies, Inc. and/or its subsidiaries.
// SPDX-License-Identifier: BSD-3-Clause

package main

import (
	"encoding/json"
	"log"
	"os"

	"github.com/arduino/go-paths-helper"
)

type InfoJSON struct {
	Latest   Release   `json:"latest"`
	Releases []Release `json:"releases"`
}

type Release struct {
	Version string `json:"version"`
	Url     string `json:"url"`
	Sha256  string `json:"sha256"`
	Kernel  string `json:"kernel,omitempty"`
	OS      string `json:"os,omitempty"`
	Board   string `json:"board,omitempty"`
}

func main() {
	version := os.Getenv("BUILD_ID")
	if version == "" {
		log.Fatal("BUILD_ID environment variable must be set")
	}
	url := os.Getenv("URL")
	if url == "" {
		log.Fatal("URL environment variable must be set")
	}
	sha256 := os.Getenv("CHECKSUM")
	if sha256 == "" {
		log.Fatal("CHECKSUM environment variable must be set")
	}
	osVersion := os.Getenv("OS_VERSION")
	if osVersion == "" {
		log.Fatal("OS_VERSION environment variable must be set")
	}
	kernel := os.Getenv("KERNEL_VERSION")
	if kernel == "" {
		log.Fatal("KERNEL_VERSION environment variable must be set")
	}
	board := os.Getenv("BOARD")
	if board == "" {
		log.Fatal("BOARD environment variable must be set")
	}

	infoJSONPath := paths.New("info.json")
	infoJSON, err := infoJSONPath.ReadFile()
	if err != nil {
		log.Fatal(err)
	}
	var info InfoJSON
	if err := json.Unmarshal(infoJSON, &info); err != nil {
		log.Fatal(err)
	}
	newRelease := Release{
		Version: version,
		Url:     url,
		Sha256:  sha256,
		OS:      osVersion,
		Kernel:  kernel,
		Board:   board,
	}
	info.Releases = append(info.Releases, newRelease)
	info.Latest = newRelease
	updatedInfoJSON, err := json.MarshalIndent(info, "", "  ")
	if err != nil {
		log.Fatal(err)
	}
	err = infoJSONPath.WriteFile(updatedInfoJSON)
	if err != nil {
		log.Fatal(err)
	}
}
