# Anti-Cheat

An generative AI for offensive anti-cheat in WeChat.

## Deploy

First, install npm.

```bash
pnpm i --registry=https://registry.npm.taobao.org/
```

Then open `config.js` and configure your apiKey in config.js:

```js
let apiKey = 'your-key-here'
```

Then initialize & run the bot with:

```bash
pnpm run anticheat
```

After login, you should see it works.

## Special Thanks

- <https://github.com/wechaty/wechaty/>
- <https://github.com/transitive-bullshit/chatgpt-api>
- <https://github.com/sunshanpeng/wechaty-chatgpt>

## License

MIT License
