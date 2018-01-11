package main

import (
	"io"
	"log"
	"os"
	"text/template"

	"github.com/AlecAivazis/survey"
	kingpin "gopkg.in/alecthomas/kingpin.v2"
)

var version string

type wizard struct {
	Output io.Writer
}

var (
	app = kingpin.New("tmplwizard", "Generate an interactive wizard from text/template")

	templateFile = app.Arg("template", "source template file").Required().ExistingFile()

	outputFile = app.Flag("output", "name of file to send output to").Short('o').OpenFile(os.O_RDWR|os.O_CREATE, 0755)
)

func main() {

	app.Version(version).VersionFlag.Short('V')
	app.HelpFlag.Short('h')
	app.UsageTemplate(kingpin.SeparateOptionalFlagsUsageTemplate)
	kingpin.MustParse(app.Parse(os.Args[1:]))

	tmpl, err := template.ParseFiles(*templateFile)

	w := wizard{}
	w.Output = os.Stdout //output to stdout by default
	//file, _ := os.Create("/tmp/test.md") //write to a file
	if *outputFile != nil {
		w.Output = *outputFile
	}

	err = tmpl.Execute(w.Output, w)
	if err != nil {
		log.Println(err)
	}

}

//PromptBool asks the user a yes/no question and returns the answer as a bool
func (w wizard) PromptBool(question string) bool {

	yesno := false
	prompt := &survey.Confirm{
		Message: question,
	}

	survey.AskOne(prompt, &yesno, nil)
	return yesno

}

//PromptString asks the user a question and returns the string the user enters
func (w wizard) PromptString(question string) string {

	answer := ""
	prompt := &survey.Input{
		Message: question,
	}
	survey.AskOne(prompt, &answer, nil)
	//return strings.TrimSpace(answer)
	return answer

}
