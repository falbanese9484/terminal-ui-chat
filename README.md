# bash-butler
 __                        __                __               __    ___                   
/\ \                      /\ \              /\ \             /\ \__/\_ \                  
\ \ \____     __      ____\ \ \___          \ \ \____  __  __\ \ ,_\//\ \      __   _ __  
 \ \ '__`\  /'__`\   /',__\\ \  _ `\  _______\ \ '__`\/\ \/\ \\ \ \/ \ \ \   /'__`\/\`'__\
  \ \ \L\ \/\ \L\.\_/\__, `\\ \ \ \ \/\______\\ \ \L\ \ \ \_\ \\ \ \_ \_\ \_/\  __/\ \ \/ 
   \ \_,__/\ \__/.\_\/\____/ \ \_\ \_\/______/ \ \_,__/\ \____/ \ \__\/\____\ \____\\ \_\ 
    \/___/  \/__/\/_/\/___/   \/_/\/_/          \/___/  \/___/   \/__/\/____/\/____/ \/_/ 


UI for chatting with LLM models in the terminal.
Purpose of this build is to familiarize myself with the bubbletea TUI framework.

To run:
```bash
export LOG_FILE_PATH=./logs/
export OPENROUTER_API_KEY=<your OPENROUTER key>
go run ./sandbox/ui-chat/main.go
```

For Debug mode and more verbose logging:
```bash
export DEBUG=1
```

The chat supports markdown using charmbracelets glamour library.

### Current Ideas for Steps Forward
1. Configure OpenRouter [x]
2. Selectable Models [ ]
3. Tooling [ ]
4. Persistant Sessions / Storage / Recall [ ]
5. Enhanced Debugger Frontend Log Viewer [ ]
6. UI Frontend Refactor [ ]
