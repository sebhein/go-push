package main

var pool = make(map[string]*PrivateChannel)

func getOrCreateChannel(channel_id string) *PrivateChannel {
  channel, ok := pool[channel_id]
  if !ok {
    channel = newPrivateChannel()
    pool[channel_id] = channel
    go channel.run()
  }
  return channel
}

func getChannel(channel_id string) (*PrivateChannel, bool) {
  channel, ok := pool[channel_id]
  return channel, ok
}
