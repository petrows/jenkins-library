package cmd

import (
	"bytes"
	"io"

	"github.com/SAP/jenkins-library/pkg/command"
	"github.com/SAP/jenkins-library/pkg/log"
	"github.com/SAP/jenkins-library/pkg/telemetry"
)

const mavenExecutable = "mvn"

func mavenExecute(config mavenExecuteOptions, telemetryData *telemetry.CustomData) string {
	var stdOut io.Writer
	var stdOutBuf *bytes.Buffer

	c := command.Command{}
	stdOut = log.Entry().Writer()

	if config.ReturnStdout {
		stdOutBuf = new(bytes.Buffer)
		stdOut = io.MultiWriter(log.Entry().Writer(), stdOutBuf)
	}

	c.Stdout(stdOut)
	c.Stderr(log.Entry().Writer())

	err := runMavenExecute(&config, &c)
	if err != nil {
		log.Entry().WithError(err).Fatal("step execution failed")
	}

	if stdOutBuf == nil {
		return ""
	}
	return string(stdOutBuf.Bytes())
}

func runMavenExecute(config *mavenExecuteOptions, command execRunner) error {

	parameters := []string{}

	if config.GlobalSettingsFile != "" {
		globalSettingsFileParameter := "--global-settings " + config.GlobalSettingsFile
		parameters = append(parameters, globalSettingsFileParameter)
	}

	if config.ProjectSettingsFile != "" {
		projectSettingsFileParameter := "--settings " + config.ProjectSettingsFile
		parameters = append(parameters, projectSettingsFileParameter)
	}

	if config.M2Path != "" {
		m2PathParameter := "-Dmaven.repo.local=" + config.M2Path
		parameters = append(parameters, m2PathParameter)
	}

	if config.PomPath != "" {
		pomPathParameter := "--file " + config.PomPath
		parameters = append(parameters, pomPathParameter)
	}

	if config.Flags != nil {
		parameters = append(parameters, config.Flags...)
	}

	parameters = append(parameters, "--batch-mode")

	if config.LogSuccessfulMavenTransfers {
		parameters = append(parameters, "-Dorg.slf4j.simpleLogger.log.org.apache.maven.cli.transfer.Slf4jMavenTransferListener=warn")
	}

	parameters = append(parameters, config.Goals...)

	err := command.RunExecutable(mavenExecutable, parameters...)
	if err != nil {
		log.Entry().
			WithError(err).
			WithField("command", append([]string{mavenExecutable}, parameters...)).
			Fatal("failed to execute run command")
	}

	return nil
}
