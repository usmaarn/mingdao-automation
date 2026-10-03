package rod

import (
	"context"
	"fmt"
	"log"
	"strings"

	"github.com/go-rod/rod"
	"github.com/go-rod/rod/lib/input"
	"github.com/go-rod/rod/lib/launcher"
	"github.com/usmaarn/mingdao/play"
)

type Server struct {
}

func Init(ctx context.Context, records []play.Record) error {
	var url = "https://mingdaoyun.palmpay-inc.com/app/e01c3b5c-5990-4881-8b24-d829db09234f/6724f8d3492cd03ba7a41173/6728e47d485926de9e240fd0"

	browser := rod.New().
		ControlURL(launcher.New().UserDataDir("./rod-data").MustLaunch()).
		MustConnect()
	defer browser.MustClose()

	page := browser.MustPage(url).Context(ctx)

	cookies := page.MustCookies()

	for _, cookie := range cookies {
		fmt.Println(cookie.Name, cookie.Domain, cookie.Expires)
	}

	log.Println("Initial URL:", page.MustInfo().URL)

	if page.MustHasR("span", "Log in using DingTalk") {
		page.MustElementR("span", "Log in using DingTalk").MustClick()

		page.MustWaitStable()

		tabSelector := `.dtm-segment-item[role='tab']`
		if page.MustHas(tabSelector) {
			tab := page.MustElements(tabSelector).Last()
			tab.MustClick()

			page.MustElement(".base-comp-check-box.module-qrcode-op-item").MustClick()

			page.MustElement(".module-qrcode-user-avatar-box").MustWaitVisible().
				MustClick()
		} else {
			log.Println("Skipped qrCode")
			page.MustElement(".module-confirm-button").MustClick()
		}

		page.MustNavigate(url)
	}

	page.MustElement("span.addRecordEntry").MustClick()

	editingBar := page.MustElement(".editingBar")

	style := editingBar.MustAttribute("style")

	if style != nil && strings.Contains(*style, "top: -140px") {
		// Hidden — skip
	} else {
		editingBar.MustElementR("div", "^Clear$").MustClick()
	}

	// continue with the rest of the automation

	for _, record := range records {
		textareas := page.MustElements("textarea.customFormTextarea.escclose")
		textareas.First().MustInput(record.Summary)

		page.MustElementR("span", "Please fill in the text content").MustClick()
		textareas.Last().MustInput(record.Solution)

		if strings.TrimSpace(record.MerchantName) != "" {
			page.MustElement(".relateRecordCardsMainBtn").MustClick()
			page.MustElement(".recordListKeyword").MustInput(record.MerchantName)

			m, err := page.Element(fmt.Sprintf(`span[title="%s" i]`, record.MerchantName))
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
	}

	log.Println("completed")
	return nil
}

func populateDates(page *rod.Page, record play.Record) {
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
