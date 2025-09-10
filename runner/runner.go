package runner

import (
	"fmt"
	"ingress-test-suite/logger"
	"ingress-test-suite/pkg/messages"
	"ingress-test-suite/test_load"
	"time"
)

var log = logger.GetLogger()

type TestResult struct {
	Host         string
	Path         string
	Success      bool
	StatusCode   int
	ErrorMessage string
}
type Tester interface {
	Test(entry test_load.IngressTestEntry) TestResult
}

type IngressManager interface {
	Create(entry test_load.IngressTestEntry, c test_load.IngressTestsFile) error
	Delete(entry test_load.IngressTestEntry) error
	CheckExist(entry test_load.IngressTestEntry) (bool, error)
}

type Runner struct {
	manager IngressManager
	tester  Tester
}

func NewRunner(manager IngressManager, tester Tester) *Runner {
	return &Runner{manager: manager, tester: tester}
}

func (r *Runner) Run(cases []test_load.IngressTestsFile) map[string][]TestResult {
	results := make(map[string][]TestResult)

	for _, c := range cases {
		log.Printf(messages.RunningIngressClass, c.IngressClassName)
		var classResults []TestResult

		for _, t := range c.Tests {
			res := r.runSingleTest(c, t)
			classResults = append(classResults, res...)
		}

		results[c.IngressClassName] = classResults
	}
	return results
}

func (r *Runner) runSingleTest(c test_load.IngressTestsFile, t test_load.IngressTestEntry) []TestResult {
	var results []TestResult

	if t.Create {
		if err := r.ensureIngressCreated(c, t); err != nil {
			results = append(results, TestResult{
				Host:         t.Host,
				Path:         t.Path,
				Success:      false,
				ErrorMessage: err.Error(),
			})
			return results
		}
		time.Sleep(2 * time.Second)
	}

	testRes := r.tester.Test(t)
	results = append(results, testRes)

	if t.Create {
		if err := r.cleanupIngress(t); err != nil {
			log.Printf(messages.FailedIngressRuleDelete, t.Host, err)
		}
	}

	return results
}

func (r *Runner) ensureIngressCreated(c test_load.IngressTestsFile, t test_load.IngressTestEntry) error {
	exists, err := r.manager.CheckExist(t)
	if err != nil {
		log.Errorf(messages.IngressRuleExist, err)
		return fmt.Errorf(messages.FailedIngressRuleCreate, err)
	}
	if exists {
		return nil
	}

	if err := r.manager.Create(t, c); err != nil {
		return fmt.Errorf(messages.FailedIngressRuleCreate, err)
	}
	return nil
}

func (r *Runner) cleanupIngress(t test_load.IngressTestEntry) error {
	exists, err := r.manager.CheckExist(t)
	if err != nil {
		return r.manager.Delete(t)
	}
	if !exists {
		return fmt.Errorf(messages.IngressRuleNotExist, t.Host)
	}
	return r.manager.Delete(t)
}
