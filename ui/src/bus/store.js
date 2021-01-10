import React from 'react';
import { EventEmitter } from 'events';

window.stores = [];
export default class Store {
  constructor(state = {}) {
    this.bus = new EventEmitter();
    this.state = state;
    this.triggers = {};
    window.stores.push(this);
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
      } else if (!Object.prototype.hasOwnProperty.call(this.state, k) || v !== this.state[k]) {
        this.state[k] = v;
      }
      this.bus.emit(`state-${k}`, this.state[k]);
      window.b2 = this.bus;
      console.log(`state-${k}`, this.bus.listeners(`state-${k}`));
      const functions = this.triggers[k];
      functions && functions.forEach((fn) => fn(this.state[k]));
    });
  }

  static ctxInState(contextType, ...bindStates) {
    return function (component) {
      console.log('inState', bindStates);
      const listeners = {};
      const newCtxComponent = class extends component {
        static contextType = contextType;

        constructor(props, ctx) {
          super(props, ctx);
          if (!this.state) {
            this.state = {};
          }
          console.log('this.context', this.context);
          bindStates.forEach((k) => {
            if (Object.prototype.hasOwnProperty.call(this.context.state, k)) {
              this.state[k] = this.context.state[k];
            }
            const key = `state-${k}`;
            listeners[key] = (v) => {
              console.log('---ctxInState this.setState---', k, v);
              this.setState({ [k]: v });
            };
            this.context.bus.addListener(key, listeners[key]);
            window.b1 = this.context.bus;
            console.log('---addListener---', key, this.context.bus.listeners(key));
          });
          console.log('---ctxInState constructor---', this.state);
        }
      };
      const { componentWillUnmount } = component.prototype;
      component.prototype.componentWillUnmount = function () {
        Object.keys(listeners).forEach((e) => {
          this.context.bus.removeListener(e, listeners[e]);
        });
        componentWillUnmount && componentWillUnmount.call(this);
      };
      return newCtxComponent;
    };
  }

  inState(...bindStates) {
    const that = this;
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
