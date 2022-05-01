package assumable

import (
	"fmt"

	"github.com/aripalo/vegas-credentials/internal/assumable/awsini"
	"github.com/aripalo/vegas-credentials/internal/checksum"
)

// New returns a struct representing all the information required to assume an
// IAM role with MFA. Effectively it parses the given dataSource
// (either file name with string type or raw data in []byte) and finds the
// correct configuration by looking up the given profileName.
func New[D awsini.DataSource](dataSource D, profileName string) (Opts, error) {
	var opts Opts
	var role awsini.Role
	var user awsini.User

	err := awsini.LoadProfile(dataSource, profileName, &role)
	if err != nil {
		return opts, err
	}

	if role.SourceProfile == "" {
		return opts, fmt.Errorf(`Profile "%s" does not contain "vegas_source_profile"`, profileName)
	}

	err = awsini.LoadProfile(dataSource, role.SourceProfile, &user)
	if err != nil {
		return opts, err
	}

	opts = Opts{
		ProfileName:     profileName,
		MfaSerial:       user.MfaSerial,
		YubikeySerial:   user.YubikeySerial,
		YubikeyLabel:    resolveYubikeyLabel(user.YubikeyLabel, user.MfaSerial),
		Region:          resolveRegion(role.Region, user.Region),
		SourceProfile:   role.SourceProfile,
		RoleArn:         role.RoleArn,
		DurationSeconds: resolveDurationSeconds(role.DurationSeconds),
		RoleSessionName: role.RoleSessionName,
		ExternalID:      role.ExternalID,
	}

	err = opts.validate()
	if err != nil {
		return opts, err
	}

	checksum, err := checksum.Generate(opts)
	if err != nil {
		return opts, err
	}

	opts.Checksum = checksum

	return opts, nil
}