import Store from './store';

export const globalStore = new Store({
  pwd: '/',
  files: [],
  showNotifications: false,
  notifications: [],
  chosen: {},
  fileSize: null,
  filesPositions: [],
  boxChosen: {},
  windows: {},
});
window.globalStore = globalStore;

export const busValue = {};
export const busState = globalStore.state;
export const addTrigger = globalStore.addTrigger.bind(globalStore);
export const setState = globalStore.setState.bind(globalStore);
export const inState = globalStore.inState.bind(globalStore);
