package app

import (
	"context"
	"fmt"
	"log"
	"strings"
	"sync"

	"github.com/go-rod/rod"
	"github.com/go-rod/rod/lib/input"
	"github.com/go-rod/rod/lib/launcher"
)

type App struct {
	Page    *rod.Page
	Browser *rod.Browser
	mutex   *sync.Mutex
}

func New(ctx context.Context) *App {
	var url = "https://mingdaoyun.palmpay-inc.com/app/e01c3b5c-5990-4881-8b24-d829db09234f/6724f8d3492cd03ba7a41173/6728e47d485926de9e240fd0"

	browser := rod.New().
		ControlURL(launcher.New().UserDataDir("./rod-data").MustLaunch()).
		MustConnect()

	page := browser.MustPage(url).Context(ctx)

	return &App{
		Browser: browser,
		Page:    page,
		mutex:   &sync.Mutex{},
	}
}

func (app *App) OpenFormModal() {
	app.Page.MustElement("span.addRecordEntry").MustClick()

	editingBar := app.Page.MustElement(".editingBar")

	style := editingBar.MustAttribute("style")

	if style != nil && strings.Contains(*style, "top: -140px") {
		// Hidden — skip
	} else {
		editingBar.MustElementR("div", "^Clear$").MustClick()
	}
}

type Record struct {
	Summary        string
	MerchantName   string
	StartDate      string
	EndDate        string
	Type           string
	Solution       string
	ScreenshotPath *string
}

func (app *App) AddRecord(record Record) {
	page := app.Page

	textareas := page.MustElements("textarea.customFormTextarea.escclose")
	textareas.First().MustInput(record.Summary)

	page.MustElementR("span", "Please fill in the text content").MustClick()
	textareas.Last().MustInput(record.Solution)

	if strings.TrimSpace(record.MerchantName) != "" {
		page.MustElement(".relateRecordCardsMainBtn").MustClick()
		page.MustElement(".recordListKeyword").MustInput(record.MerchantName)

		m, err := page.Element(fmt.Sprintf(`span[title*="%s" i]`, record.MerchantName))
		if m != nil && err == nil {
			m.MustClick()
		} else {
			log.Println("No merchant name found for record", record)
			page.MustElement(".Checkbox-box--default.Checkbox-box").MustClick()
		}
	} else {
		page.MustElement(".Checkbox-box--default.Checkbox-box").MustClick()
		if record.ScreenshotPath != nil {
			page.MustElement(".triggerTraget .icon-attachment").MustClick()
			page.MustElement(`input[type='file']`).SetFiles([]string{*record.ScreenshotPath})
			page.MustElementR(`button[type='button'] > span`, "/Confirm/i").MustClick()
		}
	}

	page.MustElement(fmt.Sprintf(`.RadioGroupCon label.Radio[title="%s" i]`, record.Type)).MustClick()

	populateDates(page, record)

	page.MustElementR(".saveAndContinueBtn", "Save & New").MustClick()

	log.Println("completed")
}

func populateDates(page *rod.Page, record Record) {
	pickers := page.MustElements(".dateTimeIcon")
	dates := []string{record.StartDate, record.EndDate}

	for i := range len(pickers) {
		pickers[i].MustClick()
		inputEl := page.MustElement(".ant-picker-input input")
		inputEl.MustWaitVisible()

		// Replace it.
		inputEl.SelectAllText()
		inputEl.MustInput(dates[i])

		inputEl.MustType(input.Enter)
	}
}
