# Flowchart for the Project

```bash
(main)
server.NewServer 
       |-----(Routine)---- server.Start
wait for stop signal           |-----------(Routine)-----server.clearExpiryRotine
       |                  aof.LoadAOF          
server.Shutdown                | 
                        server.acceptLoop       
                            ^  |------------(Routine)-----server.handleConn (record conn to close on shutdown)
                            |__|                                |-----(Routine)-----server.ReadLoop (read -> parse -> execute)
                                                                                                      ^_________________|
```
