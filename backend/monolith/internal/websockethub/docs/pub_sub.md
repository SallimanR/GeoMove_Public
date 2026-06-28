User:
1. filters data
2. subscribes to relevant data events.
Server:
3. update happens and informantion caches (reduces DB and balancer load, can be scaled horizontally + distributed for local updates only, improves latency by elimination any i/o and by storing in-memory)
4. updated data compresses once and sends directly to each subscriber

Websocket vs REST:
1. no need for http server and load balancer
2. bidirectional
3. supports binary format
4. stateful => ability to server side cache subscriptions
