package test_load

import (
	"encoding/json"
	"fmt"
	"ingress-test-suite/logger"
	"ingress-test-suite/pkg/messages"
	"os"
	"path/filepath"

	networkingv1 "k8s.io/api/networking/v1"
)

var log = logger.GetLogger()

type IngressTestsFile struct {
	IngressClassName string             `json:"ingressClassName"`
	Tests            []IngressTestEntry `json:"tests"`
}

type IngressTestEntry struct {
	Host           string `json:"host"`
	Path           string `json:"path"`
	Service        string `json:"service"`
	PathType       string `json:"pathType"`
	ExpectedStatus int    `json:"expectedStatus"`
	Namespace      string `json:"namespace"`
	Port           int    `json:"port"`
	ExtPort        int    `json:"extPort"`
	Create         bool   `json:"create"`
}

func LoadDir(dir string) []IngressTestsFile {
	var result []IngressTestsFile

	entries, err := os.ReadDir(dir)
	if err != nil {
		log.Fatalf(messages.FailedReadDir, dir, err)
	}

	for _, entry := range entries {
		if entry.IsDir() {
			continue
		}

		path := filepath.Join(dir, entry.Name())
		data, err := os.ReadFile(path)
		if err != nil {
			log.Fatalf(messages.FailedReadFile, path, err)
		}

		var tests IngressTestsFile
		if err := json.Unmarshal(data, &tests); err != nil {
			log.Fatalf(messages.FailedUnmarshalFile, path, err)
		}
		for _, t := range tests.Tests {
			if _, err := ValidatePathType(t.PathType); err != nil {
				log.Fatalf(messages.InvalidPathType, path)
			}
		}

		result = append(result, tests)
	}

	return result
}

func ValidatePathType(s string) (networkingv1.PathType, error) {
	switch s {
	case "Exact":
		pt := networkingv1.PathTypeExact
		return pt, nil
	case "Prefix":
		pt := networkingv1.PathTypePrefix
		return pt, nil
	case "ImplementationSpecific":
		pt := networkingv1.PathTypeImplementationSpecific
		return pt, nil
	default:
		return "", fmt.Errorf(messages.InvalidPathType, s)
	}
}
