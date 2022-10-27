package fluentbitstorageexporter

import (
	"fmt"
	"io/ioutil"
	"strings"
	"testing"

	pkggodevclient "github.com/guseggert/pkggodev-client"
	"github.com/stretchr/testify/require"
	"golang.org/x/mod/modfile"
)

var (
	acceptedLicenses = map[string]struct{}{
		"MIT":              struct{}{},
		"Apache-2.0":       struct{}{},
		"BSD-3-Clause":     struct{}{},
		"BSD-2-Clause":     struct{}{},
		"ISC":              struct{}{},
		"Custom-Validated": struct{}{}, //Custom licence that has been validated manually that it complies with Adevinta type of licences accepted
	}

	knownUndectedLicenses = map[string]string{
		// bufpipe was later added the MIT license: https://github.com/acomagu/bufpipe/blob/cd7a5f79d3c413d14c0c60fd31dae7b397fc955a/LICENSE
		"github.com/acomagu/bufpipe@v1.0.3":                    "MIT",
		"github.com/xi2/xz@v0.0.0-20171230120015-48954b6210f8": "Custom-Validated",
	}
)

func TestLicenses(t *testing.T) {
	b, err := ioutil.ReadFile("go.mod")
	require.NoError(t, err)
	file, err := modfile.Parse("go.mod", b, nil)
	require.NoError(t, err)
	client := pkggodevclient.New()
	for _, req := range file.Require {
		pkg, err := client.DescribePackage(pkggodevclient.DescribePackageRequest{
			Package: req.Mod.Path,
		})
		if err != nil {
			fmt.Printf(req.Mod.Path)
		}
		require.NoError(t, err)
		licences := strings.Split(pkg.License, ",")
		for _, license := range licences {
			license = strings.TrimSpace(license)
			if license == "None detected" {
				if known, ok := knownUndectedLicenses[req.Mod.String()]; ok {
					license = known
				}
			}
			if _, ok := acceptedLicenses[license]; !ok {
				t.Errorf("dependency %s is using unexpected license %s. Check that this license complies with MIT in which maiao is released and update the checks accordingly or change dependency", req.Mod, license)
			}
		}
	}
}
