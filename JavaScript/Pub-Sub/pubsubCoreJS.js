let pubsub = {};
const produce = (topic, data) => {
  let subscribers = pubsub[topic] || [];
  subscribers.forEach((sub) => sub(data));
}

const consume = (topic) => {
  if (!pubsub[topic]) pubsub[topic] = [];

  return {
    on(event, callback) {
      if (event === "data") {
        pubsub[topic].push(callback);
      }
    }
  };
}

let data = 0;
const interval = setInterval(() => {
  produce('myTopic', data);
  data += 1;
  if (data > 10) {
    clearInterval(interval);
  }
}, 1000);

const consumer = consume("myTopic");

consumer.on("data", (data) => {
  console.log('got data', data);
});
