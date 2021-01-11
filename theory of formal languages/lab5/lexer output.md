<import>
<id, WebSocket>
<from>
<path, '../websocket'>
<;>
<import>
<id, delay>
<from>
<path, '../helpers/delay'>
<;>
<import>
<id, networkBridge>
<from>
<path, '../network-bridge'>
<;>
<import>
<{>
<id, createCloseEvent>
<,>
<id, createEvent>
<}>
<from>
<path, '../event/factory'>
<;>
<export>
<function>
<id, closeWebSocketConnection>
<(>
<id, context>
<,>
<id, code>
<,>
<id, reason>
<)>
<{>
<id, context>
<.>
<id, readyState>
<=>
<id, WebSocket>
<.>
<id, CLOSING>
<;>
<const>
<id, lol>
<=>
<(>
<(>
<num, 2>
<(>
<num, 3>
<)>
<=>>
<num, 6>
<)>
<num, 13>
<)>
<>
<num, 4>
<;>
<const>
<id, server>
<=>
<id, networkBridge>
<.>
<id, serverLookup>
<(>
<id, context>
<.>
<id, url>
<)>
<;>
<const>
<id, closeEvent>
<=>
<id, createCloseEvent>
<(>
<{>
<id, type>
<:>
<str, 'close'>
<,>
<id, target>
<:>
<id, context>
<,>
<id, code>
<,>
<id, reason>
<}>
<)>
<;>
<id, delay>
<(>
<(>
<)>
<=>>
<{>
<id, networkBridge>
<.>
<id, removeWebSocket>
<(>
<id, context>
<,>
<id, context>
<.>
<id, url>
<)>
<;>
<id, context>
<.>
<id, readyState>
<=>
<id, WebSocket>
<.>
<id, CLOSED>
<;>
<id, context>
<.>
<id, dispatchEvent>
<(>
<id, closeEvent>
<)>
<;>
<if>
<(>
<id, server>
<)>
<{>
<id, server>
<.>
<id, dispatchEvent>
<(>
<id, closeEvent>
<,>
<id, server>
<)>
<;>
<}>
<}>
<,>
<id, context>
<)>
<;>
<}>
<export>
<function>
<id, failWebSocketConnection>
<(>
<id, context>
<,>
<id, code>
<,>
<id, reason>
<)>
<{>
<id, context>
<.>
<id, readyState>
<=>
<id, WebSocket>
<.>
<id, CLOSING>
<;>
<const>
<id, server>
<=>
<id, networkBridge>
<.>
<id, serverLookup>
<(>
<id, context>
<.>
<id, url>
<)>
<;>
<const>
<id, closeEvent>
<=>
<id, createCloseEvent>
<(>
<{>
<id, type>
<:>
<str, 'close'>
<,>
<id, target>
<:>
<id, context>
<,>
<id, code>
<,>
<id, reason>
<,>
<id, wasClean>
<:>
<id, false>
<}>
<)>
<;>
<const>
<id, errorEvent>
<=>
<id, createEvent>
<(>
<{>
<id, type>
<:>
<str, 'error'>
<,>
<id, target>
<:>
<id, context>
<}>
<)>
<;>
<id, delay>
<(>
<(>
<)>
<=>>
<{>
<id, networkBridge>
<.>
<id, removeWebSocket>
<(>
<id, context>
<,>
<id, context>
<.>
<id, url>
<)>
<;>
<id, context>
<.>
<id, readyState>
<=>
<id, WebSocket>
<.>
<id, CLOSED>
<;>
<id, context>
<.>
<id, dispatchEvent>
<(>
<id, errorEvent>
<)>
<;>
<id, context>
<.>
<id, dispatchEvent>
<(>
<id, closeEvent>
<)>
<;>
<if>
<(>
<id, server>
<)>
<{>
<id, server>
<.>
<id, dispatchEvent>
<(>
<id, closeEvent>
<,>
<id, server>
<)>
<;>
<}>
<}>
<,>
<id, context>
<)>
<;>
<}>
