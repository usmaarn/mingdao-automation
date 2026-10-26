
import { Mutex } from "async-mutex";
import { z } from "zod";

const RequestDataSchema = z.object({
    summary: z.string().trim().min(1, "Summary is required"),
    solution: z.string().trim().min(1, "Solution is required"),
    startDate: z.string()
        .min(1, "Start date is required")
        .refine(isValidDate, "Invalid start date. Expected YYYY-MM-DD HH:mm"),
    endDate: z.string()
        .min(1, "End date is required")
        .refine(isValidDate, "Invalid end date. Expected YYYY-MM-DD HH:mm"),
    type: z.enum([
        "issue locating",
        "documentantion",
        "documentation",
        "phone call",
        "meeting",
    ]),
});

const RequestPayloadSchema = z.object({
    merchantName: z.string().optional().default(""),
    data: z.array(RequestDataSchema).min(1, "At least one record is required"),
});

type RequestData = z.infer<typeof RequestDataSchema>;
type RequestPayload = z.infer<typeof RequestPayloadSchema>;

function isValidDate(value: string): boolean {
    if (!/^\d{4}-\d{2}-\d{2} \d{2}:\d{2}$/.test(value)) {
        return false;
    }

    const date = new Date(`${value.replace(" ", "T")}:00`);

    return (
        !Number.isNaN(date.getTime()) &&
        date.getFullYear() === Number(value.slice(0, 4)) &&
        date.getMonth() + 1 === Number(value.slice(5, 7)) &&
        date.getDate() === Number(value.slice(8, 10)) &&
        date.getHours() === Number(value.slice(11, 13)) &&
        date.getMinutes() === Number(value.slice(14, 16))
    );
}

interface AppAutomation {
    addRecord(record: {
        summary: string;
        merchantName: string;
        startDate: string;
        endDate: string;
        type: string;
        solution: string;
    }): Promise<void>;
}

export class App {
    private readonly mutex = new Mutex();

    constructor(private readonly automation: AppAutomation) { }

    async serve(): Promise<void> {
        const app = this;

        const server = Bun.serve({
            port: Number(Bun.env.PORT ?? 0),

            async fetch(request) {
                if (request.method !== "POST") {
                    return Response.json(
                        { error: "Method not allowed" },
                        { status: 405 },
                    );
                }

                return app.handleRequest(request);
            },
        });

        console.log(`Server running at ${server.url}`)
    }

    private async handleRequest(request: Request): Promise<Response> {
        let body: unknown;

        try {
            body = await request.json();
        } catch {
            return Response.json(
                { error: "Invalid JSON body" },
                { status: 400 },
            );
        }

        const result = RequestPayloadSchema.safeParse(body);

        if (!result.success) {
            return Response.json(
                { errors: result.error.flatten() },
                { status: 400 },
            );
        }

        const payload = result.data;

        // Schedule browser work without blocking the HTTP response.
        void this.processRecords(payload).catch((error: unknown) => {
            console.error("Unable to process records:", error);
        });

        return Response.json({ message: "queued" });
    }

    private async processRecords(payload: RequestPayload): Promise<void> {
        await this.mutex.runExclusive(async () => {
            let count = 1;

            for (const data of payload.data) {
                const type =
                    data.type.toLowerCase() === "documentation"
                        ? "documentantion"
                        : data.type.toLowerCase();

                await this.automation.addRecord({
                    summary: data.summary,
                    merchantName: payload.merchantName,
                    startDate: data.startDate,
                    endDate: data.endDate,
                    type,
                    solution: data.solution,
                });

                console.log(`Record ${count} completed.`);
                count++;
            }
        });
    }
}
