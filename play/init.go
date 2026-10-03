package play

import (
	"fmt"
	"log"

	"github.com/mxschmitt/playwright-go"
)

func OpenPage(pw *playwright.Playwright) (playwright.BrowserContext, playwright.Page, error) {
	browser, err := pw.Chromium.LaunchPersistentContext(
		"./browser-data",
		playwright.BrowserTypeLaunchPersistentContextOptions{
			Headless: playwright.Bool(false),
		},
	)
	if err != nil {
		return nil, nil, fmt.Errorf("Error launching browser: %w", err)
	}

	page, err := browser.NewPage()
	if err != nil {
		return nil, nil, fmt.Errorf("Error opening page: %w", err)
	}

	return browser, page, nil
}

func Init(records []Record) error {
	var url = "https://mingdaoyun.palmpay-inc.com/app/e01c3b5c-5990-4881-8b24-d829db09234f/6724f8d3492cd03ba7a41173/6728e47d485926de9e240fd0"

	pw, err := playwright.Run()
	if err != nil {
		return fmt.Errorf("Error running playwright: %w", err)
	}
	defer pw.Stop()

	browser, page, err := OpenPage(pw)
	if err != nil {
		return err
	}
	defer browser.Close()

	if _, err := page.Goto(url); err != nil {
		return fmt.Errorf("unable to open url: %w", err)
	}

	if err := LoginToMingdao(page, url); err != nil {
		return err
	}

	_, err = OpenFormModal(page)
	if err != nil {
		return err
	}

	if err := AddRecords(records, page); err != nil {
		log.Println("Unable to add record", err)
	}

	return nil
}
