package main

import (
    "fmt"
    "io"
    "os"
    // "bufio"
    "log"
    // "sort"
    "strings"

    "github.com/charmbracelet/bubbles/list"
    tea "github.com/charmbracelet/bubbletea"
    "github.com/charmbracelet/lipgloss"
	todo "github.com/1set/todotxt"

    
	// "github.com/knadh/koanf/v2"
	// "github.com/knadh/koanf/parsers/toml"
	// "github.com/knadh/koanf/providers/file"
)

const todoPath = "/home/z/todo/todo.txt"
const donePath = "/home/z/todo/done.txt"
const listHeight = 20

// Global koanf instance. Use "." as the key path delimiter. This can be "/" or any character.
// var k = koanf.New(".")
var (
    titleStyle        = lipgloss.NewStyle().MarginLeft(2)
    itemStyle         = lipgloss.NewStyle().PaddingLeft(1)
    selectedItemStyle = lipgloss.NewStyle().PaddingLeft(1).Foreground(lipgloss.Color("170"))
    paginationStyle   = list.DefaultStyles().PaginationStyle.PaddingLeft(4)
    helpStyle         = list.DefaultStyles().HelpStyle.PaddingLeft(4).PaddingBottom(1)
    quitTextStyle     = lipgloss.NewStyle().Margin(1, 0, 2, 4)
)

type item string

func (i item) FilterValue() string { return "" }

type itemDelegate struct{}

func (d itemDelegate) Height() int                                                       { return 1 }
func (d itemDelegate) Spacing() int                                                      { return 0 }
func (d itemDelegate) Update(_ tea.Msg, _ *list.Model) tea.Cmd { return nil }
func (d itemDelegate) Render(w io.Writer, m list.Model, index int, listItem list.Item) {
    i, ok := listItem.(item)
    if !ok {
        return
    }

    str := fmt.Sprintf("%s", i)

    fn := itemStyle.Render
    if index == m.Index() {
        fn = func(s ...string) string {
            return selectedItemStyle.Render("" + strings.Join(s, " "))
        }
    }

    fmt.Fprint(w, fn(str))
}

type model struct {
    list         list.Model
    choice   string
    quitting bool
}

func (m model) Init() tea.Cmd {
    return nil
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
    switch msg := msg.(type) {

    case tea.WindowSizeMsg:
        m.list.SetWidth(msg.Width)
        return m, nil

    case tea.KeyMsg:
        switch keypress := msg.String(); keypress {
        case "q", "ctrl+c":
            m.quitting = true
            return m, tea.Quit
        // mark as done
        case "d":
            return m, tea.Quit
        // edit
        case "e":
            return m, tea.Quit
        case "enter":
            i, ok := m.list.SelectedItem().(item)
            if ok {
                m.choice = string(i)
            }
            return m, tea.Quit
        }
    }

    var cmd tea.Cmd
    m.list, cmd = m.list.Update(msg)
    return m, cmd
}

func (m model) View() string {
    if m.choice != "" {
        return quitTextStyle.Render(fmt.Sprintf("%s? Sounds good to me.", m.choice))
    }
    if m.quitting {
        return quitTextStyle.Render("Not hungry? That’s cool.")
    }
    return "\n" + m.list.View()
}

func main() {
	// Load config.
	// if err := k.Load(file.Provider("godo.toml"), toml.Parser()); err != nil {
	// 	log.Fatalf("error loading config: %v", err)
	// }
    
	// fmt.Println("parent's name is = ", k.String("parent1.name"))
	// fmt.Println("parent's ID is = ", k.Int("parent1.id"))

    items := []list.Item{}

    const defaultWidth = 20

    // parse todo.txt
    if tasklist, err := todo.LoadFromPath(todoPath); err != nil {
        log.Fatal(err)
    } else {
        tasks := tasklist.Filter(todo.FilterNotCompleted)
        _ = tasks.Sort(todo.SortPriorityAsc, todo.SortProjectAsc)
        for _, t := range tasks {

            // prepend priority
            priority := ""
            if t.HasPriority() {
                priority = t.Priority
            } else {
                priority = " "
            }

            // prepend id
            taskString := fmt.Sprintf("%3d %s %s", t.ID, priority, t.Todo)
            
            // append due date
            if t.HasDueDate() {
                taskString = taskString + " due:" + t.DueDate.Format("2006-01-02")
            }
            items = append(items, item(taskString))
        }
        // if err = tasks.WriteToPath("today-todo.txt"); err != nil {
        //     log.Fatal(err)
        // }
    }

    l := list.New(items, itemDelegate{}, defaultWidth, listHeight)
    l.Title = "Todo list"
    l.SetShowStatusBar(false)
    l.SetFilteringEnabled(false)
    l.Styles.Title = titleStyle
    l.Styles.PaginationStyle = paginationStyle
    l.Styles.HelpStyle = helpStyle

    m := model{list: l}

    if _, err := tea.NewProgram(m).Run(); err != nil {
        fmt.Println("Error running program:", err)
        os.Exit(1)
    }
}
