# github-a2a

## client

```go
message, err := newClient.SendMessage(types.MessageSendParam{
	Message: &types.Message{
		ContextID: "context-1",
		Role:      types.User,
		Parts: []types.Part{
			&types.TextPart{Text: "Show recent commits for repository 'facebook/react", Kind: "text"},
        },
    },
})
```

## server

```shell
export DEEPSEEK_API_KEY = ""
export GITHUB_TOKEN = ""
```

start the server
```go
go run sever.go
```

## output

```shell
1. **[97cdd5d3](https://github.com/facebook/react/commit/97cdd5d3c33eda77be4f96a43f72d6916d3badbb)**  
   `[eslint] Do not allow useEffectEvent fns to be called in arbitrary closures (#33544)`  
   By: Jordan Brown  
   Date: 2025-07-10  

2. **[eb7f8b42](https://github.com/facebook/react/commit/eb7f8b42c92ed804bbf7f700d2bdda276d591007)**  
   `[Flight] Add Separate Outgoing Debug Channel (#33754)`  
   By: Sebastian Markbåge  
   Date: 2025-07-10  

3. **[eed25607](https://github.com/facebook/react/commit/eed25607629f5e67f13f53e91edec12b3388559f)**  
   `[Flight] Treat empty message as a close signal (#33756)`  
   By: Sebastian Markbåge  
   Date: 2025-07-10  

4. **[463b8081](https://github.com/facebook/react/commit/463b808176ad7c9429a4981bb45a1da225fd4b85)**  
   `[Fizz] Reset the segment id assignment when postponing the root (#33755)`  
   By: Josh Story  
   Date: 2025-07-10  

5. **[96c61b7f](https://github.com/facebook/react/commit/96c61b7f1f145b9fe5103051b636959cdeb20cc8)**  
   `[compiler] Add CompilerError.UnsupportedJS variant (#33750)`  
   By: Joseph Savona  
   Date: 2025-07-10  

6. **[0bfa404b](https://github.com/facebook/react/commit/0bfa404bacbad78af5b39c080ba67535f2e53044)**  
   `[compiler] More precise errors for invalid import/export/namespace statements (#33748)`  
   By: Joseph Savona  
   Date: 2025-07-10  

7. **[81e1ee74](https://github.com/facebook/react/commit/81e1ee7476a68fdf13c63d3002e5ef1b699b6842)**  
   `[compiler] Support inline enums (flow/ts), type declarations (#33747)`  
   By: Joseph Savona  
   Date: 2025-07-10  

8. **[4a3ff8ee](https://github.com/facebook/react/commit/4a3ff8eed65f96cda7617150f92de3544d5ddf6a)**  
   `[compiler] Errors for eval(), with statements, class declarations (#33746)`  
   By: Joseph Savona  
   Date: 2025-07-10  

9. **[ec4374c3](https://github.com/facebook/react/commit/ec4374c3872b320af60f322289c30cd3d7066bdf)**  
   `[compiler] Show logged errors in playground (#33740)`  
   By: Joseph Savona  
   Date: 2025-07-09  

10. **[60b5271a](https://github.com/facebook/react/commit/60b5271a9ad0e9eec2489b999ce774d39d09285b)**  
    `[Flight] Call finishHaltedTask on sync aborted tasks in stream abort listeners (#33743)`  
    By: Sebastian Markbåge  
    Date: 2025-07-09  

Let me know if you'd like more details about any of these commits!
```