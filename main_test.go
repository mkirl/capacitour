package main

import (
	"testing"

	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	api "github.com/mkirl/capacitour/api"
	"github.com/stretchr/testify/assert"
)

func TestUpdate(t *testing.T) {
	// Initialize the model
	m := initialModel()

	// Test window size message
	windowSizeMsg := tea.WindowSizeMsg{Width: 80, Height: 24}
	updatedModel, _ := m.Update(windowSizeMsg)
	assert.Equal(t, 80, updatedModel.(model).list.Width())
	assert.Equal(t, 24, updatedModel.(model).list.Height())

	// Test key press "down"
	keyDownMsg := tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'j'}}
	updatedModel, _ = m.Update(keyDownMsg)
	assert.Equal(t, 1, updatedModel.(model).list.Index())

	// Test key press "up"
	keyUpMsg := tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'k'}}
	updatedModel, _ = m.Update(keyUpMsg)
	assert.Equal(t, 0, updatedModel.(model).list.Index())

	// Test key press "enter"
	keyEnterMsg := tea.KeyMsg{Type: tea.KeyEnter}
	updatedModel, _ = m.Update(keyEnterMsg)
	assert.Equal(t, "WTF", updatedModel.(model).choice)

	// Test key press "q"
	keyQuitMsg := tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'q'}}
	updatedModel, cmd := m.Update(keyQuitMsg)
	assert.True(t, updatedModel.(model).quitting)
	assert.Equal(t, tea.Quit(), cmd) // Corrected comparison
}

func TestSpacesResponseMsg(t *testing.T) {
	// Initialize the model
	m := initialModel()

	// Create a spacesResponseMsg
	spaces := []api.Space{
		{Title: "Space 1"},
		{Title: "Space 2"},
	}
	items := []list.Item{
		item{title: "Space 1", desc: "Description for Space 1"},
		item{title: "Space 2", desc: "Description for Space 2"},
	}
	msg := spacesResponseMsg{spaces: spaces, items: items}

	// Update the model with the message
	updatedModel, _ := m.Update(msg)

	// Check that the model was updated correctly
	assert.Equal(t, "spaces", updatedModel.(model).currentView)
	assert.Equal(t, spaces, updatedModel.(model).spaces)
	assert.Equal(t, items, updatedModel.(model).list.Items())
}
