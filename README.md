# demo-golang-websocket

Websocket implementation
Server accept connection at 15000/tcp, receives text string and echos it back to WS client.

Features
========

Server
------
- Server leverages lightweight raw TCP connection instead of net/http Server (by means of github.com/gobwas/ws lib);
- Server implements Zero copy upgrade;
- Server omits goroutine blocking read by means of Linux epoll facility. It uses mail.ru library easygo/netpoll,
  that allow to wait for data readiness in websocket and run handler only when it need (omiting run blocking)
- Server control goroutine resourses. It manages Go Routine pool, allowing to create only specicifed number of concurrent
  Go routines
  

Used libraries
==============

- github.com/gobwas/ws
- github.com/gobwas/ws/wsutil
- github.com/mailru/easygo/netpoll (https://godoc.org/github.com/mailru/easygo/netpoll)


Used ideas:
===========
https://habr.com/ru/company/mailru/blog/331784/
