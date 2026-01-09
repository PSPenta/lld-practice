const { EventEmitter } = require('events');

let pubsub = new EventEmitter();
pubsub.setMaxListeners(5);

const produce = (topic, data) => {
  pubsub.emit(topic, data);
};

const consume = (topic) => {
  return {
    on(event, callback) {
      if (event == 'data') {
        pubsub.on(topic, callback);
      }
    }
  }
};

let data = 0;
const interval = setInterval(() => {
	produce("myTopic", data);
       data += 1;
	if (data > 10) {
		clearInterval(interval);
	}
}, 1000);

const consumer = consume("myTopic");

consumer.on("data", (data) => {
	console.log("got data", data);
});
