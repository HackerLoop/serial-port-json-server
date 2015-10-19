Rotonde fork
====

This fork brings compatibility with [rotonde](https://github.com/HackerLoop/rotonde), and cleans the API in the process.

All features from the excellent [SPJS](https://github.com/johnlauer/serial-port-json-server) are now available through rotonde in a clean json API.

Please read the README of rotonde before going further. [here](https://github.com/HackerLoop/rotonde).

Actions
====

This modules exposes the following actions:

SERIAL_LIST tells the server to send a list of available serial ports
```
{  
   "type":"def",
   "payload":{  
      "identifier":"SERIAL_LIST",
      "type":"action",
      "fields":null
   }
}
```

SERIAL_OPEN opens the given device at the given baudrat and buffer
algorithm.
TODO: explain buffer algorithms
```
{  
   "type":"def",
   "payload":{  
      "identifier":"SERIAL_OPEN",
      "type":"action",
      "fields":[  
         {  
            "name":"port",
            "type":"string",
            "units":""
         },
         {  
            "name":"baudrate",
            "type":"number",
            "units":""
         },
         {  
            "name":"buffer",
            "type":"string",
            "units":""
         }
      ]
   }
}
```

SERIAL_SENDJSON: TODO link to documentation in SPJS
```
{  
   "type":"def",
   "payload":{  
      "identifier":"SERIAL_SENDJSON",
      "type":"action",
      "fields":[  
         {  
            "name":"json",
            "type":"string",
            "units":""
         }
      ]
   }
}
```

SERIAL_SEND: send a string to a given port
```
{  
   "type":"def",
   "payload":{  
      "identifier":"SERIAL_SEND",
      "type":"action",
      "fields":[  
         {  
            "name":"port",
            "type":"string",
            "units":""
         },
         {  
            "name":"data",
            "type":"string",
            "units":""
         }
      ]
   }
}
```

SERIAL_SENDJSON: TODO link to documentation in SPJS
```
{  
   "type":"def",
   "payload":{  
      "identifier":"SERIAL_SENDNOBUF",
      "type":"action",
      "fields":[  
         {  
            "name":"port",
            "type":"string",
            "units":""
         },
         {  
            "name":"data",
            "type":"string",
            "units":""
         }
      ]
   }
}
```

SERIAL_CLOSE closes the connection to a port
```
{  
   "type":"def",
   "payload":{  
      "identifier":"SERIAL_CLOSE",
      "type":"action",
      "fields":[  
         {  
            "name":"port",
            "type":"string",
            "units":""
         }
      ]
   }
}
```

SERIAL_BUFFERALGORITHMS lists the available buffer algorithms
TODO link to documentation in SPJS
```
{  
   "type":"def",
   "payload":{  
      "identifier":"SERIAL_BUFFERALGORITHMS",
      "type":"action",
      "fields":[  
         {  
            "name":"port",
            "type":"string",
            "units":""
         }
      ]
   }
}
```

SERIAL_BAUDRATES lists the available baudrates
```
{  
   "type":"def",
   "payload":{  
      "identifier":"SERIAL_BAUDRATES",
      "type":"action",
      "fields":null
   }
}
```

SERIAL_RESTART triggers a restart (?)
```
{  
   "type":"def",
   "payload":{  
      "identifier":"SERIAL_RESTART",
      "type":"action",
      "fields":null
   }
}
```

SERIAL_EXIT exits the server
```
{  
   "type":"def",
   "payload":{  
      "identifier":"SERIAL_EXIT",
      "type":"action",
      "fields":null
   }
}
```

TODO link to documentation in SPJS
```
{  
   "type":"def",
   "payload":{  
      "identifier":"SERIAL_FRO",
      "type":"action",
      "fields":[  
         {  
            "name":"port",
            "type":"string",
            "units":""
         },
         {  
            "name":"mult",
            "type":"number",
            "units":""
         }
      ]
   }
}
```

SERIAL_MEMSTATS returns the memstat object
```
{  
   "type":"def",
   "payload":{  
      "identifier":"SERIAL_MEMSTATS",
      "type":"action",
      "fields":null
   }
}
```

SERIAL_VERSION ask the server to send a SERIAL_VERSION event
```
{  
   "type":"def",
   "payload":{  
      "identifier":"SERIAL_VERSION",
      "type":"action",
      "fields":null
   }
}
```

SERIAL_HOSTNAME ask the server to send SERIAL_HOSTNAME event
```
{  
   "type":"def",
   "payload":{  
      "identifier":"SERIAL_HOSTNAME",
      "type":"action",
      "fields":null
   }
}
```

SERIAL_PROGRAM programs an arduino
TODO link to documentation in SPJS
```
{  
   "type":"def",
   "payload":{  
      "identifier":"SERIAL_PROGRAM",
      "type":"action",
      "fields":[  
         {  
            "name":"port",
            "type":"string",
            "units":""
         },
         {  
            "name":"arch",
            "type":"string",
            "units":""
         },
         {  
            "name":"path",
            "type":"string",
            "units":""
         }
      ]
   }
}
```

SERIAL_PROGRAM programs an arduino from an URL
TODO link to documentation in SPJS
```
{  
   "type":"def",
   "payload":{  
      "identifier":"SERIAL_PROGRAMFROMURL",
      "type":"action",
      "fields":[  
         {  
            "name":"port",
            "type":"string",
            "units":""
         },
         {  
            "name":"arch",
            "type":"string",
            "units":""
         },
         {  
            "name":"url",
            "type":"string",
            "units":""
         }
      ]
   }
}
```

Events
===

This module sends the following events:

```
{  
   "type":"def",
   "payload":{  
      "identifier":"SERIAL_COMMANDS",
      "type":"event",
      "fields":[  
         {  
            "name":"Commands",
            "type":"",
            "units":""
         }
      ]
   }
}
```

```
{  
   "type":"def",
   "payload":{  
      "identifier":"SERIAL_ERROR",
      "type":"event",
      "fields":[  
         {  
            "name":"Error",
            "type":"",
            "units":""
         }
      ]
   }
}
```

```
{  
   "type":"def",
   "payload":{  
      "identifier":"SERIAL_BUFFLOWDEBUG",
      "type":"event",
      "fields":[  
         {  
            "name":"BufFlowDebug",
            "type":"",
            "units":""
         }
      ]
   }
}
```

```
{  
   "type":"def",
   "payload":{  
      "identifier":"SERIAL_MEMSTATS",
      "type":"event",
      "fields":[  
         {  
            "name":"Alloc",
            "type":"number",
            "units":""
         },
         {  
            "name":"TotalAlloc",
            "type":"number",
            "units":""
         }
      ]
   }
}
```

```
{  
   "type":"def",
   "payload":{  
      "identifier":"SERIAL_HOSTNAME",
      "type":"event",
      "fields":[  
         {  
            "name":"Hostname",
            "type":"",
            "units":""
         }
      ]
   }
}
```

```
{  
   "type":"def",
   "payload":{  
      "identifier":"SERIAL_VERSION",
      "type":"event",
      "fields":[  
         {  
            "name":"Version",
            "type":"",
            "units":""
         }
      ]
   }
}
```

```
{  
   "type":"def",
   "payload":{  
      "identifier":"SERIAL_GC",
      "type":"event",
      "fields":[  
         {  
            "name":"gc",
            "type":"",
            "units":""
         }
      ]
   }
}
```

```
{  
   "type":"def",
   "payload":{  
      "identifier":"SERIAL_EXIT",
      "type":"event",
      "fields":[  
         {  
            "name":"Exiting",
            "type":"",
            "units":""
         }
      ]
   }
}
```

```
{  
   "type":"def",
   "payload":{  
      "identifier":"SERIAL_RESTART",
      "type":"event",
      "fields":[  
         {  
            "name":"Restarting",
            "type":"",
            "units":""
         }
      ]
   }
}
```

```
{  
   "type":"def",
   "payload":{  
      "identifier":"SERIAL_RESTARTED",
      "type":"event",
      "fields":[  
         {  
            "name":"Restarted",
            "type":"",
            "units":""
         }
      ]
   }
}
```

```
{  
   "type":"def",
   "payload":{  
      "identifier":"SERIAL_BROADCAST",
      "type":"event",
      "fields":[  
         {  
            "name":"Cmd",
            "type":"",
            "units":""
         },
         {  
            "name":"Msg",
            "type":"",
            "units":""
         }
      ]
   }
}
```

```
{  
   "type":"def",
   "payload":{  
      "identifier":"SERIAL_PROGRAMMERSTATUS",
      "type":"event",
      "fields":[  
         {  
            "name":"ProgrammerStatus",
            "type":"",
            "units":""
         },
         {  
            "name":"Url",
            "type":"",
            "units":""
         }
      ]
   }
}
```

```
{  
   "type":"def",
   "payload":{  
      "identifier":"SERIAL_CMDCOMPLETE",
      "type":"event",
      "fields":[  
         {  
            "name":"Cmd",
            "type":"",
            "units":""
         },
         {  
            "name":"Id",
            "type":"",
            "units":""
         },
         {  
            "name":"P",
            "type":"",
            "units":""
         },
         {  
            "name":"BufSize",
            "type":"",
            "units":""
         },
         {  
            "name":"D",
            "type":"",
            "units":""
         }
      ]
   }
}
```

```
{  
   "type":"def",
   "payload":{  
      "identifier":"SERIAL_PORTMESSAGE",
      "type":"event",
      "fields":[  
         {  
            "name":"P",
            "type":"",
            "units":""
         },
         {  
            "name":"D",
            "type":"",
            "units":""
         }
      ]
   }
}
```

```
{  
   "type":"def",
   "payload":{  
      "identifier":"SERIAL_OUTPUT",
      "type":"event",
      "fields":[  
         {  
            "name":"Cmd",
            "type":"",
            "units":""
         },
         {  
            "name":"Desc",
            "type":"",
            "units":""
         },
         {  
            "name":"Port",
            "type":"",
            "units":""
         }
      ]
   }
}
```

```
{  
   "type":"def",
   "payload":{  
      "identifier":"SERIAL_PORTLIST",
      "type":"event",
      "fields":[  
         {  
            "name":"SerialPorts",
            "type":"",
            "units":""
         }
      ]
   }
}
```
