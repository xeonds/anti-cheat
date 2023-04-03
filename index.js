import { WechatyBuilder } from 'wechaty';
import qrcodeTerminal from 'qrcode-terminal';
import { ChatGPTAPI } from 'chatgpt';
import { apiKey, requestBuilder } from './config.js';

const api = new ChatGPTAPI({ apiKey: apiKey || process.env.OPENAI_API_KEY });	// import key from apiKey.js or environment variable
const conversationPool = new Map();
const wechaty = WechatyBuilder.build({
  name: 'anti-cheat',
  puppet: 'wechaty-puppet-wechat',
  puppetOptions: {
    uos: true,
  },
});

wechaty
  .on('scan', async (qrcode, status) => {
    qrcodeTerminal.generate(qrcode); 	// show login qrcode on console
    const qrcodeImageUrl = ['https://api.qrserver.com/v1/create-qr-code/?data=', encodeURIComponent(qrcode)].join('');
    console.log(qrcodeImageUrl);
  })
  .on('login', user => console.log(`User ${user} logged in`))
  .on('logout', user => console.log(`User ${user} has logged out`))
  .on('friendship', async friendship => {
    try {
      console.log(`received friend event from ${friendship.contact().name()}, messageType: ${friendship.type()}`);
    } catch (e) {
      console.error(e);
    }
  })
  .on('message', async message => {
    const contact = message.talker();
    const receiver = message.listener();
    let content = message.text();
    const room = message.room();
    const isText = message.type() === wechaty.Message.Type.Text;
    if (!isText) {
      return;
    }
    if (room) { } // room message
	else {
      console.log(`contact: ${contact} content: ${content}`);
      reply(null, contact, content);
    }
  });

async function reply(room, contact, content) {
  content = content.trim();
  if (content === 'ding') {		// bot function testing, should be disabled at production environment
    const target = room || contact;
    await send(target, 'dong');
  }
  if (content.startsWith('-c ')) {
    const request = content.replace('-c ', '');
    await chatgptReply(room, contact, request);
  }
}

async function chatgptReply(room, contact, request) {
  console.log(`contact: ${contact} request: ${request}`);
  let response = 'Something unexpected happened, please try again later...';
  try {
    let opts = {};
    // conversation
    let conversation = conversationPool.get(contact.id);
    if (conversation) {
      opts = conversation;
    }
    opts.timeoutMs = 2 * 60 * 1000;
	request = requestBuilder(request);
    let res = await api.sendMessage(request, opts);
    response = res.text;
    console.log(`contact: ${contact} response: ${response}`);
    conversation = {
      conversationId: res.conversationId,
      parentMessageId: res.id,
    };
    conversationPool.set(contact.id, conversation);
  } catch (e) {
    if (e.message === 'ChatGPTAPI error 429') {
      response = 'Previous request is still processing, please try again later...';
    }
    console.error(e);
  }
  //response = `${request} \n ------------------------ \n` + response;
  const target = room || contact;
  await send(target, response);
}

async function send(contact, message) {
  try {
    await contact.say(message);
  } catch (e) {
    console.error(e);
  }
}

// main process
wechaty
  .start()
  .then(() => console.log('Start to log in wechat...'))
  .catch(e => console.error(e));
