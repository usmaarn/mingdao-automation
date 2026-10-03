package play

import (
	"fmt"
	"log"

	"github.com/mxschmitt/playwright-go"
)

func LoginToMingdao(page playwright.Page, redirect string) error {
	log.Println("Starting DingTalk login...")

	loginBtn := page.GetByText("Log in using DingTalk")

	btnOpt := playwright.LocatorWaitForOptions{Timeout: new(5000.0)}

	if err := loginBtn.WaitFor(btnOpt); err != nil {
		return fmt.Errorf("Login Button Error: %w", err)
	}

	if visible, err := loginBtn.IsVisible(); err == nil && visible {
		if err := loginBtn.Click(); err != nil {
			return fmt.Errorf("click DingTalk login: %w", err)
		}
	}

	confirmBtn := page.Locator(".module-confirm-button")

	if err := confirmBtn.WaitFor(btnOpt); err != nil {
		return fmt.Errorf("Confirm Button Error: %w", err)
	}

	if visible, err := confirmBtn.IsVisible(); err == nil && visible {
		if err := confirmBtn.Click(); err != nil {
			return fmt.Errorf("click login confirmation: %w", err)
		}
	}

	if err := page.WaitForURL(
		"**mingdaoyun.palmpay-inc.com/dashboard/**",
		playwright.PageWaitForURLOptions{
			Timeout: new(5_000.0),
		},
	); err != nil {
		return fmt.Errorf("Mingdao did not finish authentication: %w", err)
	}

	log.Println("After authentication:", page.URL())

	return nil
}
