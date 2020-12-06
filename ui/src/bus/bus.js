import { EventEmitter } from 'events';

const bus = new EventEmitter();
export const busState = {
  pwd: '/',
  files: [],
  showNotifications: false,
  notifications: [],
  chosen: {},
  fileSize: null,
  filesPositions: [],
  boxChosen: {},
  cutFiles: [],
  copyFiles: [],
};
export const busValue = {};

window.busState = busState;
window.busValue = busValue;

const triggers = {};

export function addTrigger(stateName, fn) {
  if (triggers[stateName]) {
    triggers[stateName].push(fn);
  } else {
    triggers[stateName] = [fn];
  }
  console.log('triggers', triggers);
}

export function setState(states) {
  console.log('---setState---', states);
  Object.keys(states).forEach((k) => {
    const v = states[k];
    if (typeof v === 'function') {
      console.log('---setState before---', busState[k]);
      v(busState[k]);
      console.log('---setState after---', busState[k]);
      bus.emit(`state-${k}`, busState[k]);
    } else if (!Object.prototype.hasOwnProperty.call(busState, k) || v !== busState[k]) {
      busState[k] = v;
      bus.emit(`state-${k}`, busState[k]);
    }
    const functions = triggers[k];
    functions && functions.forEach((fn) => fn(busState[k]));
  });
}

export function inState(...bindStates) {
  return function (component) {
    console.log('inState', bindStates);
    const listeners = {};
    const newComponent = class extends component {
      constructor(props, ctx) {
        super(props, ctx);
        if (!this.state) {
          this.state = {};
        }
        bindStates.forEach((k) => {
          if (Object.prototype.hasOwnProperty.call(busState, k)) {
            this.state[k] = busState[k];
          }
          const key = `state-${k}`;
          listeners[key] = (v) => {
            console.log('---inState this.setState---', k, v);
            this.setState({ [k]: v });
          };
          bus.addListener(key, listeners[key]);
        });
        console.log('---inState constructor---', this.state);
      }
    };
    const { componentWillUnmount } = component.prototype;
    component.prototype.componentWillUnmount = function () {
      Object.keys(listeners).forEach((e) => {
        bus.removeListener(e, listeners[e]);
      });
      componentWillUnmount && componentWillUnmount.call(this);
    };
    return newComponent;
  };
}
window.busState = busState;
export default bus;
