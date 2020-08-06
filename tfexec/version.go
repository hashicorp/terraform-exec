package tfexec

import (
	"bytes"
	"context"
	"fmt"
	"regexp"
	"strings"

	"github.com/hashicorp/go-version"
)

// Version returns structured output from the terraform version command including both the Terraform CLI version
// and any initialized provider versions. This will read cached values when present unless the skipCache parameter
// is set to true.
func (tf *Terraform) Version(ctx context.Context, skipCache bool) (tfVersion *version.Version, providerVersions map[string]*version.Version, err error) {
	tf.versionLock.Lock()
	defer tf.versionLock.Unlock()

	if tf.execVersion == nil || skipCache {
		tf.execVersion, tf.provVersions, err = tf.version(ctx)
		if err != nil {
			return nil, nil, err
		}
	}

	return tf.execVersion, tf.provVersions, nil
}

// version does not use the locking on the Terraform instance and should probably not be used directly, prefer Version.
func (tf *Terraform) version(ctx context.Context) (*version.Version, map[string]*version.Version, error) {
	// TODO: 0.13.0-beta2? and above supports a `-json` on the version command, should add support
	// for that here and fallback to string parsing

	versionCmd := tf.buildTerraformCmd(ctx, "version")
	var errBuf strings.Builder
	var outBuf bytes.Buffer
	versionCmd.Stderr = &errBuf
	versionCmd.Stdout = &outBuf

	err := versionCmd.Run()
	if err != nil {
		return nil, nil, parseError(err, errBuf.String())
	}

	tfVersion, providerVersions, err := parseVersionOutput(outBuf.String())
	if err != nil {
		return nil, nil, fmt.Errorf("unable to parse version: %w", err)
	}

	return tfVersion, providerVersions, nil
}

var (
	simpleVersionRe = `v?(?P<version>[0-9]+(?:\.[0-9]+)*(?:-[A-Za-z0-9\.]+)?)`

	versionOutputRe         = regexp.MustCompile(`^Terraform ` + simpleVersionRe)
	providerVersionOutputRe = regexp.MustCompile(`(\n\+ provider[\. ](?P<name>\S+) ` + simpleVersionRe + `)`)
)

func parseVersionOutput(stdout string) (*version.Version, map[string]*version.Version, error) {
	stdout = strings.TrimSpace(stdout)

	submatches := versionOutputRe.FindStringSubmatch(stdout)
	if len(submatches) != 2 {
		return nil, nil, fmt.Errorf("unexpected number of version matches %d for %s", len(submatches), stdout)
	}
	v, err := version.NewVersion(submatches[1])
	if err != nil {
		return nil, nil, fmt.Errorf("unable to parse version %q: %w", submatches[1], err)
	}

	allSubmatches := providerVersionOutputRe.FindAllStringSubmatch(stdout, -1)
	provV := map[string]*version.Version{}

	for _, submatches := range allSubmatches {
		if len(submatches) != 4 {
			return nil, nil, fmt.Errorf("unexpected number of providerion version matches %d for %s", len(submatches), stdout)
		}

		v, err := version.NewVersion(submatches[3])
		if err != nil {
			return nil, nil, fmt.Errorf("unable to parse provider version %q: %w", submatches[3], err)
		}

		provV[submatches[2]] = v
	}

	return v, provV, err
}
