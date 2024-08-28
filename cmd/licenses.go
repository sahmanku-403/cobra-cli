

package cmd

import (
	"fmt"
	"strings"
	"time"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// Licenses contains all possible licenses a user can choose from.
var Licenses = make(map[string]License)

type License struct {
	Name            string   // The type of license in use
	PossibleMatches []string // Similar names to guess
	Text            string   // License text data
	Header          string   // License header for source files
}

func init() {
	// Allows a user to not use a license.
	Licenses["none"] = License{"None", []string{"none", "false"}, "", ""}

	initApache2()
	initMit()
	initBsdClause3()
	initBsdClause2()
	initGpl2()
	initGpl3()
	initLgpl()
	initAgpl()
}

// getLicense returns license specified by user in flag or in config.
// If user didn't specify the license, it returns none
//
// TODO: Inspect project for existing license
func getLicense() License {
	// If explicitly flagged, use that.
	if userLicense != "" {
		return findLicense(userLicense)
	}

	// If user wants to have custom license, use that.
	if viper.IsSet("license.header") || viper.IsSet("license.text") {
		return License{Header: viper.GetString("license.header"),
			Text: viper.GetString("license.text")}
	}

	// If user wants to have built-in license, use that.
	if viper.IsSet("license") {
		return findLicense(viper.GetString("license"))
	}

	// If user didn't set any license, use none by default
	return Licenses["none"]
}

func copyrightLine() string {
	author := viper.GetString("author")

	year := viper.GetString("year") // For tests.
	if year == "" {
		year = time.Now().Format("2006")
	}

	return "Copyright © " + year + " " + author
}

// findLicense looks for License object of built-in licenses.
// If it didn't find license, then the app will be terminated and
// error will be printed.
func findLicense(name string) License {
	found := matchLicense(name)
	if found == "" {
		cobra.CheckErr(fmt.Errorf("unknown license: " + name))
	}
	return Licenses[found]
}

// matchLicense compares the given a license name
// to PossibleMatches of all built-in licenses.
// It returns blank string, if name is blank string or it didn't find
// then appropriate match to name.
func matchLicense(name string) string {
	if name == "" {
		return ""
	}

	for key, lic := range Licenses {
		for _, match := range lic.PossibleMatches {
			if strings.EqualFold(name, match) {
				return key
			}
		}
	}

	return ""
}
