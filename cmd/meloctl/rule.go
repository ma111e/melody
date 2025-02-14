package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/ma111e/melody/internal/meloctl/prompt"

	"github.com/google/uuid"
	"github.com/k0kubun/pp"
	"github.com/ma111e/melody/internal/rules"
	"github.com/spf13/cobra"
	"gopkg.in/yaml.v3"
)

var (
	ruleCmd = &cobra.Command{
		Use:   "rule",
		Short: "Handle Melody rule files",
		Long:  `This subcommand is used to handle Melody rule files`,
	}
	checkRuleCmd = &cobra.Command{
		Use:   "check",
		Args:  cobra.ExactArgs(1),
		Short: "Check rules",
		Long:  `This subcommand is used to check rules`,
		Run:   checkRule,
	}
	initRuleCmd = &cobra.Command{
		Use:   "init",
		Args:  cobra.ExactArgs(1),
		Short: "Create a basic rule template",
		Long:  `This subcommand is used to create a basic rule template`,
		Run:   initRule,
	}
	addRuleCmd = &cobra.Command{
		Use:   "add",
		Args:  cobra.ExactArgs(1),
		Short: "Add a rule in the specified rule file",
		Long:  `This subcommand is used to add a rule template in an existing rule file`,
		Run:   addRule,
	}
	//pullRuleCmd = &cobra.Command{
	//	Use:   "pull",
	//	Args:  cobra.ExactArgs(1),
	//	Short: "Update local melody rules in the $melody.rules.home dir",
	//	Long:  `This subcommand is used to update the rulesets for each folder in the $melody.rules.home directory with a git pull`,
	//	Run:   pullRules,
	//}

	ruleTemplateStringVars = map[string]string{
		"layer":       "",
		"name":        "",
		"author":      "",
		"description": "",
		"status":      "",
	}

	// Check
	isVerbose bool

	// Init
	layer string
	name  string
	//version     string
	author      string
	description string
	status      string
	references  []string
	tags        map[string]string

	force       bool
	interactive bool
)

// RuleTemplateMeta represents the 'meta' section of a rule template
type RuleTemplateMeta struct {
	Version     string `yaml:"version" pretty:"Version"`
	ID          string `yaml:"id"`
	Author      string `yaml:"author" pretty:"Author"`
	Status      string `yaml:"status" pretty:"Status" choices:"stable,experimental,testing"`
	Created     string `yaml:"created" pretty:"Created" validate:"date"`
	Modified    string `yaml:"modified" pretty:"Modified" validate:"date"`
	Description string `yaml:"description" pretty:"Description"`
}

// RuleTemplate represents a rule template
type RuleTemplate struct {
	Layer string           `yaml:"layer" pretty:"Layer" choices:"http,icmp,tcp,udp,ip"`
	Meta  RuleTemplateMeta `yaml:"meta"`
	//References  []string                       `yaml:"references" pretty:"References"`
	Match      map[string]map[string][]string `yaml:"match"`
	References []string                       `yaml:"references"`
	//Tags        map[string]string              `yaml:"tags" pretty:"Tags" hint:"ie. 'cve: cve-2020-9054'"`
	Tags map[string]string `yaml:"tags"`
}

func init() {
	RootCmd.AddCommand(ruleCmd)

	ruleCmd.AddCommand(checkRuleCmd)
	checkRuleCmd.Flags().BoolVarP(&isVerbose, "verbose", "v", false, "Print additional info while checking")

	ruleCmd.AddCommand(initRuleCmd)
	initRuleCmd.Flags().StringVarP(&layer, "layer", "l", "http", `Layer field for new rule`)
	initRuleCmd.Flags().StringVarP(&name, "name", "n", "Changeme", `Name field for new rule`)
	//initRuleCmd.Flags().StringVarP(&version, "version", "V", "1.0", `Version field for new rule`)
	initRuleCmd.Flags().StringVarP(&author, "author", "a", "Changeme", `Author field for new rule`)
	initRuleCmd.Flags().StringVarP(&description, "description", "d", "", `Description field for new rule`)

	initRuleCmd.Flags().StringVarP(&status, "status", "s", "experimental", `Status field for new rule`)
	initRuleCmd.Flags().StringArrayVarP(&references, "references", "r", []string{}, `References fields new rule`)
	initRuleCmd.Flags().StringToStringVarP(&tags, "tag", "t", map[string]string{}, `Tags fields for new rule`)

	initRuleCmd.Flags().BoolVarP(&force, "force", "f", false, `Do not ask permission to overwrite if a rule already defined`)
	initRuleCmd.Flags().BoolVarP(&interactive, "interactive", "i", false, `Ask for each parameter for the new rule`)

	ruleCmd.AddCommand(addRuleCmd)
	addRuleCmd.Flags().StringVarP(&layer, "layer", "l", "http", `Layer field for new rule`)
	addRuleCmd.Flags().StringVarP(&name, "name", "n", "Changeme", `Name field for new rule`)
	//addRuleCmd.Flags().StringVarP(&version, "version", "V", "1.0", `Version field for new rule`)
	addRuleCmd.Flags().StringVarP(&author, "author", "a", "Changeme", `Author field for new rule`)
	addRuleCmd.Flags().StringVarP(&description, "description", "d", "", `Description field for new rule`)

	addRuleCmd.Flags().StringVarP(&status, "status", "s", "experimental", `Status field for new rule`)
	addRuleCmd.Flags().StringArrayVarP(&references, "references", "r", []string{}, `References fields new rule`)
	addRuleCmd.Flags().StringToStringVarP(&tags, "tag", "t", map[string]string{}, `Tags fields for new rule`)

	addRuleCmd.Flags().BoolVarP(&interactive, "interactive", "i", false, `Ask for each parameter for the new rule`)

	//ruleCmd.AddCommand(pullRuleCmd)
}

func makeRule() (string, string, error) {
	now := time.Now()
	currentDate := now.Format(prompt.DateFormat)
	ruleID := uuid.NewString()

	ruleTemplateStringVars["layer"] = layer
	ruleTemplateStringVars["name"] = name
	ruleTemplateStringVars["author"] = author
	ruleTemplateStringVars["description"] = description
	ruleTemplateStringVars["status"] = status
	ruleTemplateStringVars["created"] = currentDate
	ruleTemplateStringVars["modified"] = currentDate
	ruleTemplateStringVars["version"] = "1.0"

	if interactive {
		err := prompt.AskAll(RuleTemplate{}, &ruleTemplateStringVars)
		if err != nil {
			fmt.Println(err)
			return "", ruleID, err
		}

		err = prompt.AskAll(RuleTemplateMeta{}, &ruleTemplateStringVars)
		if err != nil {
			fmt.Println(err)
			return "", ruleID, err
		}
	}

	vars := RuleTemplate{
		Layer: ruleTemplateStringVars["layer"],
		Meta: RuleTemplateMeta{
			Version:     ruleTemplateStringVars["version"],
			ID:          ruleID,
			Author:      ruleTemplateStringVars["author"],
			Status:      ruleTemplateStringVars["status"],
			Created:     ruleTemplateStringVars["created"],
			Modified:    ruleTemplateStringVars["modified"],
			Description: ruleTemplateStringVars["description"],
		},
		References: references,
		Tags:       tags,
	}

	vars.Match = map[string]map[string][]string{}
	vars.Match["http.uri"] = map[string][]string{
		"startswith|any": {
			"",
		},
		"contains|nocase": {
			"",
		},
		"is|regex": {
			"",
		},
		"endswith": {
			"",
		},
	}

	rule := map[string]RuleTemplate{
		name: vars,
	}

	formattedRule, err := yaml.Marshal(&rule)
	if err != nil {
		fmt.Println(err)
		return "", ruleID, err
	}

	return string(formattedRule), ruleID, nil
}

func makeRuleFilepath(rulePath string, ruleID string, append bool) (string, error) {
	var rulePathWithExt string

	if !(strings.HasSuffix(rulePath, ".yml")) {
		rulePathWithExt = rulePath + ".yml"
	} else {
		rulePathWithExt = rulePath
		rulePath = strings.TrimSuffix(rulePath, ".yml")
	}

	if _, err := os.Stat(rulePathWithExt); !os.IsNotExist(err) && !force && !append {
		overwrite, err := prompt.AskConfirmation(fmt.Sprintf(`"%s" is already defined. Overwrite ?`, rulePathWithExt), false)
		if err != nil {
			return "", err
		}

		if !overwrite {
			rulePathWithExt = fmt.Sprintf("%s_%s.yml", rulePath, ruleID[:4])
			fmt.Printf("Using \"%s\" instead\n", rulePathWithExt)
		}
	}

	rulePath, err := filepath.Abs(rulePathWithExt)
	if err != nil {
		return "", err
	}

	return rulePath, err
}

func checkRuleName(filepath string, name string) (bool, error) {
	targetRule := make(map[string]interface{})

	data, err := ioutil.ReadFile(filepath)
	if err != nil {
		return true, err
	}

	err = yaml.Unmarshal(data, &targetRule)
	if err != nil {
		return true, err
	}

	for key := range targetRule {
		if key == name {
			return true, nil
		}
	}

	return false, nil
}

func initRule(_ *cobra.Command, args []string) {
	rulePath := args[0]
	formattedRule, ruleID, err := makeRule()
	if err != nil {
		fmt.Println(err)
		return
	}

	absRulePath, err := makeRuleFilepath(rulePath, ruleID, false)
	if err != nil {
		fmt.Println(err)
		return
	}

	fmt.Println("\nWriting :\n", formattedRule)

	err = ioutil.WriteFile(absRulePath, []byte(formattedRule), 0644)
	if err != nil {
		fmt.Println(err)
		return
	}

	fmt.Printf("✅ [%s]: Rule file created\n", absRulePath)
}

func addRule(_ *cobra.Command, args []string) {
	rulePath := args[0]
	var ruleNameExists bool

	formattedRule, ruleID, err := makeRule()
	if err != nil {
		fmt.Println(err)
		return
	}

	absRulePath, err := makeRuleFilepath(rulePath, ruleID, true)
	if err != nil {
		fmt.Println(err)
		return
	}

	if ruleNameExists, err = checkRuleName(absRulePath, name); err != nil {
		fmt.Println(err)
		return
	}

	if ruleNameExists {
		fmt.Printf("❌ [%s]: Rule name \"%s\" is already defined\n", absRulePath, name)
		return
	}

	fmt.Println("\nWriting :\n", formattedRule)

	formattedRule = "\n" + formattedRule

	f, err := os.OpenFile(absRulePath, os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0644)
	if err != nil {
		fmt.Println(err)
		return
	}

	defer f.Close()

	if _, err = f.WriteString(formattedRule); err != nil {
		fmt.Println(err)
		return
	}

	fmt.Printf("✅ [%s]: Rule added\n", absRulePath)
}

//func pullRules(_ *cobra.Command, args []string) {
// Walk all dir in melodyConf.RulesDir
// If there is a .git folder, trigger a git pull, with --force if the -f arg is given to "meloctl rule pull -f"
// We instantiate a new repository targeting the given path (the .git folder)
//r, err := git.PlainOpen(meloctlConf.MelodyHomeDir)
//if err != nil {
//	fmt.Println(err)
//	return
//}
//
//// Get the working directory for the repository
//w, err := r.Worktree()
//if err != nil {
//	fmt.Println(err)
//	return
//}
//
//// Pull the latest changes from the origin remote and merge into the current branch
//err = w.Pull(&git.PullOptions{RemoteName: "origin"})
//if err != nil {
//	fmt.Println(err)
//	return
//}
//// Print the latest commit that was just pulled
//ref, err := r.Head()
//if err != nil {
//	fmt.Println(err)
//	return
//}
//commit, err := r.CommitObject(ref.Hash())
//if err != nil {
//	fmt.Println(err)
//	return
//}
//fmt.Println(commit)
//}

func checkRule(_ *cobra.Command, args []string) {
	var err error
	rulePath := args[0]

	target, err := os.Stat(rulePath)
	if err != nil {
		fmt.Println(err)
		return
	}

	switch mode := target.Mode(); {
	case mode.IsDir():
		checkRuleDir(rulePath)

	case mode.IsRegular():
		checkRuleFile(rulePath)
	}

}

func checkRuleFile(rulePath string) {
	var atLeastOneError bool

	parsedRawRules, err := rules.ParseYAMLRulesFile(rulePath)
	if err != nil {
		fmt.Printf("❌ [%s]: %s\n", rulePath, err)
		return
	}

	for _, rawRule := range parsedRawRules {
		parsed, err := rawRule.Parse()
		if err != nil {
			atLeastOneError = true
			fmt.Printf("❌ [%s]: %s\n", rulePath, err)
			continue
		}

		if isVerbose {
			_, _ = pp.Println(parsed)
		}
	}

	if !atLeastOneError {
		fmt.Printf("✅ [%s]: OK\n", rulePath)
	}
}

func checkRuleDir(rulesDir string) {
	err := filepath.Walk(rulesDir,
		func(path string, info os.FileInfo, err error) error {
			if err != nil {
				return err
			}
			if info.IsDir() {
				return nil
			}

			if strings.HasSuffix(path, ".yml") {
				checkRuleFile(path)
				return nil
			}

			fmt.Printf("== Skipping %s (does not ends with '.yml')\n", path)
			return nil
		})

	if err != nil {
		fmt.Printf("❌ [%s]: %s\n", rulesDir, err)
	}
}
