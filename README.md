# bash-butler
 <pre style="color: #75baff; font-weight: bold;">
 __                        __                __               __    ___                   
/\ \                      /\ \              /\ \             /\ \__/\_ \                  
\ \ \____     __      ____\ \ \___          \ \ \____  __  __\ \ ,_\//\ \      __   _ __  
 \ \ '__`\  /'__`\   /',__\\ \  _ `\  _______\ \ '__`\/\ \/\ \\ \ \/ \ \ \   /'__`\/\`'__\
  \ \ \L\ \/\ \L\.\_/\__, `\\ \ \ \ \/\______\\ \ \L\ \ \ \_\ \\ \ \_ \_\ \_/\  __/\ \ \/ 
   \ \_,__/\ \__/.\_\/\____/ \ \_\ \_\/______/ \ \_,__/\ \____/ \ \__\/\____\ \____\\ \_\ 
    \/___/  \/__/\/_/\/___/   \/_/\/_/          \/___/  \/___/   \/__/\/____/\/____/ \/_/ 
</pre>

UI for chatting with LLM models in the terminal.
Purpose of this build is to familiarize myself with the bubbletea TUI framework.
---
This application is in development and is not intended for use as a final product.

### [Demo](https://s3.us-east-1.amazonaws.com/images.proaistudios.com/bash-butler-demo.mov)

To run:
```bash
export LOG_FILE_PATH=./logs/
export OPENROUTER_API_KEY=<your OPENROUTER key>
go run ./cmd/bash-butler/main.go
```
---

Right not the app will load ollama, and you would need to change the main.go file to 
initialize the open router provider - which is commented out.

For Debug mode and more verbose logging:
```bash
export DEBUG=1
```

The chat supports markdown using charmbracelets glamour library.

You can toggle between different available models for each provider using the 
Ctrl+F key.
