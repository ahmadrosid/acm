# acm

Stop thinking about commit message. Let OpenAI GPT-3 handle your commit message, and auto-commit your git changes.

## Install

```bash
go install https://github.com/ahmadrosid/acm.git
```

## Usage

```bash
git add .
acm
```

Example result:
```bash
$ acm

The proposed commit message!
---------------------------------------------------------------------------------------------------------------------------------------------------------
Add post about setup macos for laravel development

This post is about how to setup macos for laravel development.
---------------------------------------------------------------------------------------------------------------------------------------------------------
Do you want to continue ? (y/n): y
```

## ENV
Add this env manually `OPENAI_API_KEY` :

```
export OPENAI_API_KEY="..."
```
