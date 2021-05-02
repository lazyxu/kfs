import React from 'react';
import Store from './store';
import './fs';

export * from './global';

export const { ctxInState } = Store;
export const StoreContext = React.createContext();
export default Store;
