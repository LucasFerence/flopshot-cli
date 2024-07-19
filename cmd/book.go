package cmd

import (
	"fmt"
	"os"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/huh"
	"github.com/spf13/cobra"
)

// bookCmd represents the book command
var bookCmd = &cobra.Command{
	Use:   "book",
	Short: "Book a tee time round",
	Run: func(cmd *cobra.Command, args []string) {

		model := newModel()

		p := tea.NewProgram(model)
		if _, err := p.Run(); err != nil {
        	fmt.Printf("Error: %v", err)
         	os.Exit(1)
    	}

     	fmt.Println(model.form.GetString("user_id"))
	},
}

func init() {
	rootCmd.AddCommand(bookCmd)
}

// -------------------------------------------------

type model struct {
	form *huh.Form
}

func newModel() model {
	m := model{}

	m.form = huh.NewForm(
		huh.NewGroup(
			huh.NewInput().
				Key("user_id").
				Title("User ID"),
			huh.NewInput().
				Key("course_id").
				Title("Course ID"),
			huh.NewConfirm(),
		),
	)

	return m
}

func (m model) Init() tea.Cmd {
	return m.form.Init()
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {

	var cmds []tea.Cmd

	form, cmd := m.form.Update(msg)
    if f, ok := form.(*huh.Form); ok {
        m.form = f
        cmds = append(cmds, cmd)
    }

    if m.form.State == huh.StateCompleted {
		cmds = append(cmds, tea.Quit)
	}

    return m, tea.Batch(cmds...)
}

func (m model) View() string {

    if m.form.State == huh.StateCompleted {
        user := m.form.GetString("user_id")
        course := m.form.GetString("course_id")

        return fmt.Sprintf("You selected: User %s, Course %s", user, course)
    }

    return m.form.View()
}
