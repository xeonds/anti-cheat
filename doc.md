# 技术文档-AntiCheat

## 项目概述

本项目旨在利用OpenAI的davinci003和prompt构建一个识别潜在诈骗行为的辅助AI，帮助网络安全意识较低的普通人来识别和防范网络诈骗威胁。本项目的目标是接收一篇对话，并给出其中另一方可能是诈骗犯的置信度。

## 项目原理

本项目基于open ai的davinci003模型，这是一个具有强大NLP能力的大型语言模型（LLM），能够处理各种自然语言处理（NLP）任务，包括生成和分类文本、回答问题、翻译文本等¹。LLM是一种利用大量的无标注文本数据来训练的机器学习模型，它可以通过自监督学习或半监督学习的方式来学习语言的语法和语义知识，并且能够在不同的任务之间进行迁移学习²。

我们团队借助prompt，instruction和rules来对模型进行微调，从而让它适应于识别潜在的诈骗性对话。具体来说：

- prompt是一段文本，用来引导模型生成或分类文本。prompt可以包含一些示例或模板，来告诉模型我们想要什么样的输出格式或风格。
- instruction是一段文本，用来告诉模型如何完成某个任务。instruction可以包含一些具体的步骤或条件，来限制模型的行为或输出。
- rules是一些文本或符号，用来表示模型应当遵循的规则。rules可以包含一些逻辑或约束，来保证模型的输出符合我们的期望或要求。

在我们的项目中，我们在prompt中陈述了它的任务，即给出一个对话，并判断其中另一方是否是诈骗犯；在instruction中陈述了它的工作方法，即根据对话中的内容和语气，分析另一方是否有诈骗意图或手段；在rules中陈述了它应当遵循的规则，即给出一个0到1之间的数字作为置信度，并且根据置信度高低给出相应的建议或警告。

我们使用了一个概率预测模型（probability prediction model）来实现我们的项目。概率预测模型是一种机器学习模型，它可以根据输入数据的特征，给出一个介于0和1之间的连续概率值，表示某个事件发生的可能性。在我们的项目中，我们使用了davinci003模型作为我们的概率预测模型，并且通过prompt，instruction和rules来调整它的参数和输出。我们使用了逻辑回归（logistic regression）作为我们的概率预测方法，它是一种线性模型，可以通过拟合一个逻辑函数（logistic function），也叫做S形函数（sigmoid function），来生成一个介于0和1之间的概率值³。在我们的项目中，我们使用了逻辑回归来预测对话中另一方是否是诈骗犯的概率，并且根据概率高低给出相应的建议或警告。

## 项目优势和局限性

本项目有以下优势：

- LLM在NLP方面具有非常强大的潜力，这使得它给出的结果可信度有保障。
- 相比于人工的优势，它能做到连续工作，以及无人值守的工作，并且能同时应对大量的任务。这使得我们的项目作为一个反诈骗服务，具有相当高的可用性。
- 我们的AI能够识别基于对话的诈骗置信度预测，具体类型并非一成不变，我们可以不断添加可能的诈骗话术、套路等，来更新模型，也就是说，几乎是所有基于对话的诈骗类型都可以被识别。

本项目也有以下局限性：

- 对于过于内在逻辑过于复杂的语言，受限于LLM的自身限制，目前还难以识别。
- 训练数据来源的问题，需要保证数据质量和合法性。
- 识别结果置信度问题，AI难以保证绝对正确，所以应该明确告知用户AI的局限性。

## 参考资料

(1) Large language model - Wikipedia. <https://en.wikipedia.org/wiki/Large_language_model>
(2) 通向AGI之路：大型语言模型（LLM）技术精要 - 知乎. <https://zhuanlan.zhihu.com/p/597586623>.
(3) Logistic Regression: Calculating a Probability | Machine .... <https://developers.google.com/machine-learning/crash-course/logistic-regression/calculating-a-probability>.
(4) What Is the Difference Between Large Language Model <https://kennowak.com/llm-vs-npl/>.
(5) Large Language Models and Where to Use Them: Part 1 <https://txt.cohere.com/llm-use-cases/>.
(6) Best practices for prompt engineering with OpenAI API - OpenAI. <https://help.openai.com/en/articles/6654000-best-practices-for-prompt-engineering-with-openai-api>
(7) Logistic Regression: Calculating a Probability | Machine Learning | Google Developers. <https://developers.google.com/machine-learning/crash-course/logistic-regression/calculating-a-probability>
