package builder

import (
	"encoding/json"
	"fmt"
	"log"
	"os"

	"github.com/go-telegram/bot/models"
)

const (
	KeyboardStudentMainMenuInlineButtons string = "student_main_menu_inline_buttons"
	KeyboardTeacherMainMenuInlineButtons string = "teacher_main_menu_inline_buttons"

	KeyboardFinishOrInsertWordsButtons string = "teacher_finish_or_insert_words"
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

func InlineKeyboardPaginator(testSetsID []int64, offSet int32, isLastPage bool) *models.InlineKeyboardMarkup {
	markup := &models.InlineKeyboardMarkup{
		InlineKeyboard: make([][]models.InlineKeyboardButton, 3),
	}
	markup.InlineKeyboard[0] = make([]models.InlineKeyboardButton, 0, len(testSetsID))
	for id, testSetID := range testSetsID {
		fmt.Printf("paginator test id: %d\n", id)
		button := models.InlineKeyboardButton{
			Text:         fmt.Sprintf("%d", int(offSet)+id+1),
			CallbackData: fmt.Sprintf("teacher_test_set_%d", testSetID),
		}
		markup.InlineKeyboard[0] = append(markup.InlineKeyboard[0], button)
	}
	markup.InlineKeyboard[1] = make([]models.InlineKeyboardButton, 2)

	if offSet >= 5 {
		markup.InlineKeyboard[1][0] = models.InlineKeyboardButton{
			Text:         "Prev",
			CallbackData: fmt.Sprintf("test_sets_page_%d", offSet-5),
		}
	}
	if !isLastPage {
		markup.InlineKeyboard[1][1] = models.InlineKeyboardButton{
			Text:         "Next",
			CallbackData: fmt.Sprintf("test_sets_page_%d", offSet+5),
		}
	}

	markup.InlineKeyboard[2] = make([]models.InlineKeyboardButton, 1)

	markup.InlineKeyboard[2][0] = models.InlineKeyboardButton{
		Text:         "Dashboard",
		CallbackData: "dashboard",
	}
	log.Println("Paginator before return")
	return markup

}
