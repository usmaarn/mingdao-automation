package play

import (
	"fmt"
	"log"

	"github.com/mxschmitt/playwright-go"
)

func OpenFormModal(page playwright.Page) (playwright.Locator, error) {
	// First determine whether we're already authenticated.
	addRecord := page.Locator("span.addRecordEntry")

	waitOpt := playwright.LocatorWaitForOptions{Timeout: new(10000.0)}

	if err := addRecord.WaitFor(waitOpt); err != nil {
		log.Println("Not authenticated.")
		return nil, fmt.Errorf("Add Record button error: %w", err)
	}

	if err := addRecord.Click(); err != nil {
		return nil, fmt.Errorf("click add record: %w", err)
	}

	// Don't use `go btn.Click()` here.
	btn := page.Locator(".sc-bQmweE.sc-kIXKos.fqoxBO")

	if visible, _ := btn.IsVisible(); visible {
		if err := btn.Click(); err != nil {
			return nil, fmt.Errorf("click modal button: %w", err)
		}
	}

	modal := page.Locator(".workSheetNewRecord .ant-modal-content")

	if err := modal.WaitFor(waitOpt); err != nil {
		return nil, fmt.Errorf("modal not found: %w", err)
	}

	return modal, nil
}
