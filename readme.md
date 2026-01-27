This project is a realisation of a telegram client along with a bittorrent client. Using event-driven architecture.


The basic flow I wanted to produce is as follows:
1. User writes a message into tg client
2. Tg client reads the message and sends back found torrents
3. User selects some of the torrent files
4. The tg client writes a message into Kafka message queue cause we don't care about order of torrents
5. Torrent client consumer sees the message and parses it
6. Torrents begin downloading using a simple torrent client
7. As soon as a certain download finishes the torrent sever sends message back to kafka
8. The tg client consumer reads it and writes a message to the user
