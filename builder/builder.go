package builder

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/go-telegram/bot/models"
)

const (
	StudentMainMenuInlineButtons string = "student_main_menu_inline_buttons"
	StudentTestNameInlineButton  string = "student_test_Name"
	TeacherMainMenuInlineButtons string = "teacher_main_menu_inline_buttons"
)

// InlineKeyboardButton represents a single button in the inline keyboard
type InlineKeyboardButton struct {
	Text         string `json:"text"`
	CallbackData string `json:"callback_data,omitempty"`
}

// InlineKeyboardConfig holds the configuration for the keyboard layout
type InlineKeyboardConfig struct {
	Buttons [][]InlineKeyboardButton `json:"buttons"`
}

// InlineKeyboardBuilder builds Telegram inline keyboards from JSON configuration
type InlineKeyboardBuilder struct {
	config InlineKeyboardConfig
}

// NewInlineKeyboardBuilder creates a new builder instance from a JSON file
func NewInlineKeyboardBuilder(inlineKeyboard string) (*InlineKeyboardBuilder, error) {
	configData, err := os.ReadFile(fmt.Sprintf("builder/%s.json", inlineKeyboard))
	if err != nil {
		return nil, err
	}

	var config InlineKeyboardConfig
	if err := json.Unmarshal(configData, &config); err != nil {
		return nil, err
	}

	return &InlineKeyboardBuilder{config: config}, nil
}

// Build creates a Telegram inline keyboard markup from the configuration
func (b *InlineKeyboardBuilder) Build() *models.InlineKeyboardMarkup {
	markup := &models.InlineKeyboardMarkup{
		InlineKeyboard: make([][]models.InlineKeyboardButton, len(b.config.Buttons)),
	}

	for rowIdx, row := range b.config.Buttons {
		markup.InlineKeyboard[rowIdx] = make([]models.InlineKeyboardButton, len(row))
		for colIdx, btn := range row {
			markup.InlineKeyboard[rowIdx][colIdx] = models.InlineKeyboardButton{
				Text:         btn.Text,
				CallbackData: btn.CallbackData,
			}
		}
	}

	return markup
}
