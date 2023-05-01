import { WechatyBuilder } from "wechaty";
import qrcodeTerminal from "qrcode-terminal";
import { ChatGPTAPI } from "chatgpt";
import { apiKey, requestBuilder, prompt } from "./config.js";

const api = new ChatGPTAPI({ apiKey: apiKey || process.env.OPENAI_API_KEY }); // import key from apiKey.js or environment variable
const conversationPool = new Map();
const wechaty = WechatyBuilder.build({
  name: "anti-cheat",
  puppet: "wechaty-puppet-wechat",
  puppetOptions: {
    uos: true,
  },
});
let opts = {};
let initialized = false;

wechaty
  .on("scan", async (qrcode, status) => {
    qrcodeTerminal.generate(qrcode); // show login qrcode on console
    const qrcodeImageUrl = [
      "https://api.qrserver.com/v1/create-qr-code/?data=",
      encodeURIComponent(qrcode),
    ].join("");
    console.log(qrcodeImageUrl);
  })
  .on("login", async (user) => {
    console.log(`User ${user} logged in`);
    for (let i = 0; i < prompt.length; i++) {
      await chatgptReply(prompt[i], "null");
      console.log(`Initialize process(${i}/${prompt.length}).`);
    }
    console.log("initialize success.");
    initialized = true;
  })
  .on("logout", (user) => console.log(`User ${user} has logged out`))
  .on("friendship", async (friendship) => {
    try {
      console.log(
        `received friend event from ${friendship
          .contact()
          .name()}, messageType: ${friendship.type()}`
      );
    } catch (e) {
      console.error(e);
    }
  })
  .on("message", async (message) => {
    if (message.self() || !initialized) {
      return;
    }
    const contact = message.talker();
    const content = message.text();
    if (
      !(message.type() == wechaty.Message.Type.Text && message.room() == null)
    ) {
      return;
    }
    console.log(`Receiving msg from ${contact}: ${content}`);
    let text = await chatgptReply(content, "user");
    const replyRegex = /(?<=reply:)(.*)(?=predict)/i;
    const predictRegex = /(?<=predict:)(.*)/i;
    let reply = text.match(replyRegex)[0];
    let predict = text.match(predictRegex)[0];
    console.log(predict);
    send(contact, reply);
  });

async function chatgptReply(request, prefix) {
  let response = "Something unexpected happened, please try again later...";
  try {
    opts.timeoutMs = 2 * 60 * 1000;
    request = requestBuilder(request, prefix);
    let res = await api.sendMessage(request, opts);
    response = res.text;
    opts = {
      conversationId: res.conversationId,
      parentMessageId: res.id,
    };
  } catch (e) {
    if (e.message === "ChatGPTAPI error 429") {
      response =
        "Previous request is still processing, please try again later...";
    }
    console.error(e);
  }
  return response;
}

async function send(contact, message) {
  try {
    await contact.say(message);
  } catch (e) {
    console.error(e);
  }
}

function sleep(ms) {
  const start = new Date().getTime();
  while (new Date().getTime() < start + ms);
}

// main process
wechaty
  .start()
  .then(() => console.log("Start to log in wechat..."))
  .catch((e) => console.error(e));
