package main

const (
	greetMsg string = "Jack's here, mate. May I help you with those green stacks?"
	helpMsg  string = "*Need help with Jack? Here are some info:*\n" +
		"/start - Wake Jack's mind. You'll receive a greeting from him.\n" +
		"/help - Show this messages, obviously. What are you expecting?\n" +
		"/earn - Log your income. Probably the best command here.\n" +
		"/spend - Log your expense. Jack hates subtraction, just saying.\n\n" +
		"*Got it?* Now time to fully utilize Jack's simple mind."
	defaultMsg string = "Jack doesn't know that command, sadly.\nMaybe you can try to type /help."
	startLog   string = "Waiting to get updates (from her)..."
	aboutMsg   string = "Who's Jack?\n\n*imalrightjack*, usually called Jack, is a bot that helps you track " +
		"your financial flow, like your income and outcome. His mind is very simple and intuitive though, like " +
		"basic income and expense logs, daily report, et cetera.\nMind you, Jack loves anything that smells money, " +
		"so I think you can trust him."
	trscErrMsgArg string = "Please specify the amount, and optionally the notes."
	trscErrMsgNum string = "Invalid amount! Use positive numbers only (e.g. 67000)\nNo commas, dots, and letters allowed."
)
