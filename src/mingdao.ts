import { sleep } from "bun";
import { chromium, type BrowserContext, type Page } from "playwright";

interface Record {
    summary: string;
    solution: string;
    merchantName: string;
    type: string;
    screenshotPath?: string;
    startDate: string;
    endDate: string;
}

export class Mingdao {
    private constructor(
        private url: string,
        private browser: BrowserContext,
        private page: Page,
    ) {
    }

    static async init(url: string) {
        const browser = await chromium.launchPersistentContext("./browser-data", {
            headless: false,
            slowMo: 500,
        })

        const page = await browser.newPage();
        await page.goto(url);

        return new Mingdao(url, browser, page);
    }

    async openModal() {
        await this.page.locator("span.addRecordEntry").click();

        const editingBar = this.page.locator(".editingBar");
        const style = await editingBar.getAttribute("style");

        if (!style?.includes("top: -140px")) {
            await editingBar.getByText("Clear", { exact: true }).click();
        }
    }

    async addRecord(record: Record): Promise<void> {
        const page = this.page;

        const textareas = page.locator(
            "textarea.customFormTextarea.escclose",
        );

        await textareas.first().fill(record.summary);

        await page
            .getByText("Please fill in the text content", { exact: true })
            .click();

        await textareas.last().fill(record.solution);

        if (record.merchantName.trim() !== "") {
            await page.locator(".relateRecordCardsMainBtn").click();

            await page
                .locator(".recordListKeyword")
                .fill(record.merchantName);

            const merchant = page.locator(
                `span[title*="${record.merchantName.replaceAll('"', '\\"')}" i]`,
            ).first();

            await sleep(1000);

            if (await merchant.count() > 0) {
                await merchant.click();
            } else {
                console.log("No merchant name found for record", record);
                console.log(await merchant.all());
                // await page
                //     .locator(".Checkbox-box--default.Checkbox-box")
                //     .click();
            }
        } else {
            await page
                .locator(".Checkbox-box--default.Checkbox-box")
                .click();

            if (record.screenshotPath) {
                await page.locator(".triggerTraget .icon-attachment").click();

                await page
                    .locator('input[type="file"]')
                    .setInputFiles(record.screenshotPath);

                await page
                    .getByRole("button", { name: /Confirm/i })
                    .click();
            }
        }

        await page
            .locator(
                `.RadioGroupCon label.Radio[title="${record.type.replaceAll('"', '\\"')}" i]`,
            )
            .click();

        await this.populateDates(page, record);

        const saveButton = page.locator(".saveAndContinueBtn");

        await saveButton.getByText("Save & New").click();

        // const message = page.locator(".ant-message-notice");

        // if (await message.isVisible().catch(() => false)) {
        //     await textareas.first().fill(record.summary);
        //     await saveButton.getByText("Save & New").click();
        // }
    }

    private async populateDates(
        page: Page,
        record: Record,
    ): Promise<void> {
        const pickers = page.locator(".dateTimeIcon");
        const dates = [record.startDate, record.endDate];

        for (let i = 0; i < dates.length; i++) {
            await pickers.nth(i).click();

            const input = page.locator(".ant-picker-input input").first();

            await input.waitFor({ state: "visible" });
            await input.fill(dates[i] as string);
            await input.press("Enter");
        }
    }
}