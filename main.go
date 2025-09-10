package main

import (
	"ingress-test-suite/k8s"
	"ingress-test-suite/logger"
	"ingress-test-suite/pkg/messages"
	"ingress-test-suite/runner"
	"ingress-test-suite/test_load"
	"os"

	"github.com/spf13/pflag"
)

var version = "dev"

var log = logger.GetLogger()

var (
	testCasesPath string
)

func init() {
	testCasesPathEnv := os.Getenv("TESTS_PATH")

	pflag.StringVarP(&testCasesPath,
		"tests-path", "t", testCasesPathEnv, messages.TestCasesPathMessageString)

	versionFlag := pflag.BoolP("version", "v", false, messages.VersionMessageString)
	helpFlag := pflag.BoolP("help", "h", false, messages.HelpMessageString)

	pflag.Parse()

	if *helpFlag {
		log.Infof(messages.Usage)
		pflag.PrintDefaults()
		log.Exit(0)
	}

	if *versionFlag {
		log.Infof(messages.Version, version)
	}

	if testCasesPath == "" {
		log.Fatalf(messages.TestCasesDirPathVariableError)
	}
}

func main() {

	testCases := test_load.LoadDir(testCasesPath)
	k8sClient := k8s.SetupK8SClient()

	manager := runner.NewK8sIngressManager(k8sClient)
	tester := &runner.HTTPTester{}
	r := runner.NewRunner(manager, tester)

	results := r.Run(testCases)

	var exitCode = 0

	for caseName, caseResults := range results {
		log.Infof(messages.ResultTestCaseInfo, caseName)
		for _, r := range caseResults {
			if r.Success {
				log.Infof(messages.Status, r.Host, r.Path, messages.OK, r.StatusCode)
			} else {
				log.Errorf(messages.Status, r.Host, r.Path, messages.Fail, r.StatusCode)
				exitCode = 12
			}
		}
	}

	log.ExitFunc(exitCode)
}
