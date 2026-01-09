class PubSub {
  constructor() {
    this.events = new Map(); //{ event: [cb1, cb2....] }
  }

  subscribe(event, cb) {
    if(!this.events.get(event)) {
      this.events.set(event, []);
    }

    this.events.get(event).push(cb);
    return {
      unsubscribe : () => this.unsubscribe(event, cb) // return unsub method
    }
  }

  unsubscribe(event, cb) {
    if(!this.events.get(event)) return;

    this.events.set(event, this.events.get(event).filter(eventCb => eventCb !== cb));
  }

  publish(event, data) {
    if(!this.events.get(event)) return;

    this.events.get(event).forEach(cb => {
      cb(data);
    });
  }
}


const pubsub = new PubSub()

const consumer1 = pubsub.subscribe('movie', (data)=> {
  console.log('from 1st movie', data)
})

pubsub.subscribe('movie', (data)=> {
  console.log('from 2nd movie', data)
})

pubsub.subscribe('music', (data)=> {
  console.log('from 1st music', data)
})

pubsub.publish('movie', 'Fight Club')


consumer1.unsubscribe();

pubsub.publish('movie', 'Fight Club 2')



