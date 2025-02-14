package rules

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/ma111e/melody/internal/logging"

	"github.com/ma111e/melody/internal/config"

	"gopkg.in/yaml.v3"
)

var (
	// GlobalRules is the global object holding all the loaded rules
	GlobalRules = make(map[string][]Rules)
)

// GlobalRawRules describes a set of RawRules
type GlobalRawRules []RawRules

// LoadRulesDir walks the given directory to find rule files and load them into GlobalRules
func LoadRulesDir(rulesDir string) uint {
	var err error
	var globalRawRules GlobalRawRules
	var total uint

	skiplist := []string{
		".gitkeep",
	}

	globalRawRules, err = ParseRulesDir(rulesDir, skiplist)
	if err != nil {
		log.Println(fmt.Sprintf("Failed to parse rule directory [%s]", rulesDir))
		log.Println(err)
		os.Exit(1)
	}

	for _, rawRules := range globalRawRules {
		rules := Rules{}
		for ruleName, rawRule := range rawRules {
			rule, err := rawRule.Parse()
			if err != nil {
				logging.Warnings.Println(err)
				continue
			}
			rule.Name = ruleName

			rules = append(rules, rule)
		}

		for _, proto := range config.Cfg.MatchProtocols {
			GlobalRules[proto] = append(GlobalRules[proto], rules.Filter(func(rule Rule) bool { return rule.Layer == proto }))
		}
	}

	for _, protocolRules := range GlobalRules {
		for _, ruleset := range protocolRules {
			total += uint(len(ruleset))
		}
	}

	return total
}

// ParseRulesDir walks a directory and parses each of the rule file it encounters
func ParseRulesDir(rulesDir string, skiplist []string) ([]RawRules, error) {
	var rawRules []RawRules
	err := filepath.Walk(rulesDir,
		func(path string, info os.FileInfo, err error) error {
			if err != nil {
				return err
			}
			if info.IsDir() {
				return nil
			}

			for _, skipped := range skiplist {
				if info.Name() == skipped {
					return nil
				}
			}

			if strings.HasSuffix(path, ".yml") {
				log.Println("Parsing", path)
				parsed, err := ParseYAMLRulesFile(path)
				if err != nil {
					return fmt.Errorf("failed to read YAML rule file [%s] : %s", path, err)
				}

				rawRules = append(rawRules, parsed)
			} else {
				return fmt.Errorf("invalid rule file (wanted : .yml) : %s", path)
			}

			return nil
		})

	return rawRules, err
}

// ParseYAMLRulesFile is an helper that parses the given YAML file and return a set of raw rules as RawRules
func ParseYAMLRulesFile(filepath string) (RawRules, error) {
	rawRules := RawRules{}
	rulesData, err := ioutil.ReadFile(filepath)
	if err != nil {
		return RawRules{}, err
	}

	if err := yaml.Unmarshal(rulesData, &rawRules); err != nil {
		return RawRules{}, err
	}

	return rawRules, nil
}
