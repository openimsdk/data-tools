package utils

import (
	"github.com/OpenIMSDK/tools/errs"
	"gopkg.in/yaml.v3"
	"os"
)

func ParseConfig[T any](path string) (*T, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, errs.Wrap(err, "ReadFile config failed")
	}
	var conf T
	if err := yaml.Unmarshal(data, &conf); err != nil {
		return nil, errs.Wrap(err, "config parse failed")
	}
	for _, c := range []any{conf, &conf} {
		if checker, ok := c.(interface{ Check() error }); ok {
			if err := checker.Check(); err != nil {
				return nil, err
			}
		}
	}
	return &conf, nil
}
