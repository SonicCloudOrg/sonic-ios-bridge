package conn

import (
	"github.com/Masterminds/semver"
	log "github.com/sirupsen/logrus"
)

func RunWebDriverAgent(device iDevice, webDriverAgentBundleID, webDriverAgentBundleIDTestRunnerBundleID string) error {
	version, err := device.GetSemverProductVersion()
	if err != nil {
		return err
	}
	if version.LessThan(semver.MustParse("14.0")) {
		log.Infof("udId: %s iOSVersion: %s, WebDriverAgent will run without secure...", device.Properties.SerialNumber, version)
	}
	return nil
}
