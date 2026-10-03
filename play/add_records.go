package play

import (
	"fmt"
	"strings"

	"github.com/mxschmitt/playwright-go"
)

type Record struct {
	Summary        string
	MerchantName   string
	StartDate      string
	EndDate        string
	Type           string
	Solution       string
	ScreenshotPath *string
}

func AddRecords(records []Record, page playwright.Page) error {
	for _, record := range records {
		textarea := page.Locator(".customFormTextarea")

		if err := textarea.Nth(0).WaitFor(); err != nil {
			return err
		}

		if err := textarea.Nth(0).Fill(record.Summary); err != nil {
			return err
		}

		if record.MerchantName != "" {
			page.Locator(".relateRecordCardsMainBtn").Click()

			input := page.Locator("input.recordListKeyword")
			input.Fill(record.MerchantName)
			input.Press("Enter")

			if el := page.Locator(fmt.Sprintf("span[title='%s' i]", strings.TrimSpace(record.MerchantName))); el != nil {
				el.Click()
			} else {
				page.Locator(`button[aria-label="Close"]`).Click()
				// UploadFile(page, record.Screenshot)
			}
		} else {
			// UploadFile(page, record.Screenshot)
		}

		page.Locator(".sc-kOHvGY.fpKcCw .RadioGroupCon").Locator("")

		page.Locator(fmt.Sprintf(`.sc-dWcDbm label[title='%s' i]`, record.Type)).Click()

		if err := textarea.Nth(1).Fill(record.Solution); err != nil {
			return err
		}

		dfs, err := page.Locator(".classtabfocus.dateTimeIcon").All()
		if err != nil {
			return err
		}

		dates := []string{record.StartDate, record.EndDate}
		for i := range len(dates) {
			dfs[i].Click()
			inputEl := page.Locator(".sc-cZLAbK.fwfisz").Locator("input")
			if err := inputEl.Fill(dates[i]); err != nil {
				return err
			}
			inputEl.Press("Enter")
		}

		page.Locator("button.saveAndContinueBtn").Click()
	}

	return nil
}

func UploadFile(page playwright.Page, screenshot []byte) {
	page.Locator(".Checkbox-box--default.Checkbox-box").Click()
	if screenshot != nil {
		page.Locator(".triggerTraget").Click()
		page.Locator(".html5_1k3r0ke0e1t9ht9h153quv66ltb").SetInputFiles(playwright.InputFile{
			Name:   "Screenshot.png",
			Buffer: screenshot,
		})
		page.GetByRole("button", playwright.PageGetByRoleOptions{
			Name: "Confirm",
		}).Click()
	}
}
