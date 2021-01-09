import React from 'react';
import { EventEmitter } from 'events';

export default class Store {
  constructor(state = {}) {
    this.bus = new EventEmitter();
    this.state = state;
    this.triggers = {};
  }

  addTrigger(stateName, fn) {
    if (this.triggers[stateName]) {
      this.triggers[stateName].push(fn);
    } else {
      this.triggers[stateName] = [fn];
    }
    console.log('this.triggers', this.triggers);
  }

  setState(state) {
    console.log('---setState---', state);
    Object.keys(state).forEach((k) => {
      const v = state[k];
      if (typeof v === 'function') {
        console.log('---setState before---', this.state[k]);
        v(this.state[k]);
        console.log('---setState after---', this.state[k]);
        this.bus.emit(`state-${k}`, this.state[k]);
      } else if (!Object.prototype.hasOwnProperty.call(this.state, k) || v !== this.state[k]) {
        this.state[k] = v;
        this.bus.emit(`state-${k}`, this.state[k]);
      }
      const functions = this.triggers[k];
      functions && functions.forEach((fn) => fn(this.state[k]));
    });
  }

  inState(...bindStates) {
    const that = this;
    console.log('that', that);
    return function (component) {
      console.log('inState', bindStates);
      const listeners = {};
      const newComponent = class extends component {
        constructor(props, ctx) {
          super(props, ctx);
          if (!this.state) {
            this.state = {};
          }
          console.log('that', that);
          bindStates.forEach((k) => {
            if (Object.prototype.hasOwnProperty.call(that.state, k)) {
              this.state[k] = that.state[k];
            }
            const key = `state-${k}`;
            listeners[key] = (v) => {
              console.log('---inState this.setState---', k, v);
              this.setState({ [k]: v });
            };
            that.bus.addListener(key, listeners[key]);
          });
          console.log('---inState constructor---', this.state);
        }
      };
      const { componentWillUnmount } = component.prototype;
      component.prototype.componentWillUnmount = function () {
        Object.keys(listeners).forEach((e) => {
          that.bus.removeListener(e, listeners[e]);
        });
        componentWillUnmount && componentWillUnmount.call(this);
      };
      return newComponent;
    };
  }
}
