import { useState } from 'react';
import { createGlobalStore } from 'hox';

const initialState = {
  drivers: [],
  // driverName: "",
  filePath: [],
  dirItems: null,
  content: null,
};

const [useResourceManager] = createGlobalStore(() => useState(initialState));

export default useResourceManager;
