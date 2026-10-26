import { App } from "./app";
import { Mingdao } from "./mingdao";

const url = "https://mingdaoyun.palmpay-inc.com/app/e01c3b5c-5990-4881-8b24-d829db09234f/6724f8d3492cd03ba7a41173/6728e47d485926de9e240fd0";

const mingdao = await Mingdao.init(url);

console.log("Opening browser modal");
await mingdao.openModal()

const app = new App(mingdao);

await app.serve();