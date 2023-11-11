import { createGlobalStore } from 'hox';
import { useState } from 'react';

const initialState = {
  drivers: [],
  // driverName: "",
  filePath: [],
  dirItems: null,
  content: null,
};

export const [useResourceManager, getResourceManager] = createGlobalStore(() => useState(initialState));

export default useResourceManager;
