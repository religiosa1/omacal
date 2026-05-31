package main

import (
	"fmt"
	"os"
	"strings"
	"time"

	tea "charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"
)

// duh
const daysInWeek = 7

var _ tea.Model = (*Calendar)(nil)

type Calendar struct {
	// current date time to show as a highlight
	now time.Time
	// selected date
	selectedDate time.Time
	// month to display
	displayMonth time.Time
}

func NewCalendar() *Calendar {
	now := time.Now()
	today := midnightDate(now)
	return &Calendar{
		now:          today,
		selectedDate: today,
		displayMonth: firstDayOfMonth(now),
	}
}

type dayChangedMsg struct{}

func tickTomorrow() tea.Cmd {
	tomorrow := midnightDate(time.Now()).
		AddDate(0, 0, 1).
		Add(time.Second) // small slop so we never fire a hair *before* midnight
	return tea.Tick(time.Until(tomorrow), func(time.Time) tea.Msg {
		return dayChangedMsg{}
	})
}

// Init implements [tea.Model].
func (c *Calendar) Init() tea.Cmd {
	return nil
}

func (c *Calendar) normalizeSelectedMonth() {
	if sameMonth(c.selectedDate, c.displayMonth) {
		return
	}
	switch {
	case c.selectedDate.Before(firstDateInTable(c.displayMonth)):
		c.displayMonth = c.displayMonth.AddDate(0, -1, 0)
	case lastDateInTable(c.displayMonth).Before(c.selectedDate.Add(time.Second)):
		c.displayMonth = c.displayMonth.AddDate(0, 1, 0)
	}
}

// Update implements [tea.Model].
func (c *Calendar) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyPressMsg:
		switch msg.String() {
		case "ctrl+c", "q", "esc":
			return c, tea.Quit
		case "up", "k":
			c.selectedDate = c.selectedDate.AddDate(0, 0, -daysInWeek)
			c.normalizeSelectedMonth()
		case "down", "j":
			c.selectedDate = c.selectedDate.AddDate(0, 0, daysInWeek)
			c.normalizeSelectedMonth()
		case "left", "h":
			c.selectedDate = c.selectedDate.AddDate(0, 0, -1)
			c.normalizeSelectedMonth()
		case "right", "l":
			c.selectedDate = c.selectedDate.AddDate(0, 0, 1)
			c.normalizeSelectedMonth()
		case "shift+left", "H":
			c.displayMonth = addMonths(c.displayMonth, -1)
			c.selectedDate = c.selectedDate.AddDate(0, -1, 0)
		case "shift+right", "L":
			c.displayMonth = addMonths(c.displayMonth, 1)
			c.selectedDate = c.displayMonth.AddDate(0, 1, 0)
		}
	case dayChangedMsg:
		c.now = midnightDate(time.Now())
		return c, tickTomorrow()
	}
	return c, nil
}

const headerString = "Mo Tu We Th Fr Sa Su"

// View implements [tea.Model].
func (c *Calendar) View() tea.View {
	var b strings.Builder
	monthTitle := lipgloss.NewStyle().
		Width(len(headerString)).
		Align(lipgloss.Center).
		Render(c.displayMonth.Format("January 2006"))
	b.WriteString(monthTitle)
	b.WriteByte('\n')
	b.WriteString(lipgloss.NewStyle().Bold(true).Render(headerString))
	b.WriteByte('\n')

	first := firstDateInTable(c.displayMonth)
	for y := range calcRowsInATable(c.displayMonth) {
		for x := range daysInWeek {
			date := first.AddDate(0, 0, (y*daysInWeek)+x)
			day := date.Day()
			st := lipgloss.NewStyle()
			switch {
			case sameDay(date, c.now):
				st = st.Reverse(true)
			case date.Month() != c.displayMonth.Month():
				st = st.Faint(true)
			}
			if sameDay(c.selectedDate, date) {
				st = st.Underline(true)
			}
			b.WriteString(st.Render(fmt.Sprintf("%2d", day)))
			b.WriteByte(' ')
		}
		fmt.Fprintln(&b)
	}
	return tea.NewView(b.String())
}

func main() {
	p := tea.NewProgram(NewCalendar())
	_, err := p.Run()
	if err != nil {
		fmt.Printf("error occurred during the start: %s", err)
		os.Exit(2)
	}
}

func sameDay(a, b time.Time) bool {
	ay, am, ad := a.Date()
	by, bm, bd := b.Date()
	return ay == by && am == bm && ad == bd
}

func sameMonth(a, b time.Time) bool {
	ay, am, _ := a.Date()
	by, bm, _ := b.Date()
	return ay == by && am == bm
}

// addMonths adds month, being aware of max days in a month
func addMonths(t time.Time, n int) time.Time {
	y, m, d := t.Date()
	first := time.Date(y, m, 1, 0, 0, 0, 0, t.Location()).AddDate(0, n, 0)
	if last := lastDayOfMonth(first).Day(); d > last {
		d = last
	}
	return time.Date(first.Year(), first.Month(), d, 0, 0, 0, 0, t.Location())
}

func firstDateInTable(t time.Time) time.Time {
	fdom := firstDayOfMonth(t)
	offset := (int(fdom.Weekday()) + 6) % daysInWeek // Mon=0 ... Sun=6
	return fdom.AddDate(0, 0, -offset)
}

func lastDateInTable(t time.Time) time.Time {
	fdom := firstDayOfMonth(t)
	nRows := calcRowsInATable(t)
	offset := (int(fdom.Weekday()) + 6) % daysInWeek // Mon=0 ... Sun=6
	return fdom.AddDate(0, 0, daysInWeek*nRows-offset)
}

// midnightDate returns a date, with time component set to midnight
func midnightDate(t time.Time) time.Time {
	return time.Date(t.Year(), t.Month(), t.Day(), 0, 0, 0, 0, t.Location())
}

// calcRowsInATable calculates the amount of rows, required to display a month
// Most months require 5 rows, but sometimes you need 6 rows, if a month starts
// on a later weekday, or 4 in case of february
func calcRowsInATable(t time.Time) int {
	fdom := firstDayOfMonth(t)
	ldom := lastDayOfMonth(t)
	offset := (int(fdom.Weekday()) + 6) % daysInWeek
	daysInMonth := ldom.Day()
	return (offset + daysInMonth + 6) / daysInWeek
}

func firstDayOfMonth(t time.Time) time.Time {
	return time.Date(
		t.Year(),
		t.Month(),
		1, // day
		0, // hour
		0, // min
		0, // sec
		0, // nano
		t.Location(),
	)
}

func lastDayOfMonth(t time.Time) time.Time {
	fdom := firstDayOfMonth(t)
	return fdom.AddDate(0, 1, -1)
}
