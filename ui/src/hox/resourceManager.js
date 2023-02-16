import { useState } from 'react';
import { createGlobalStore } from 'hox';

const initialState = {
  branchName: '默认文件夹',
  filePath: [],
  dirItems: [],
  content: null,
};

const [useResourceManager] = createGlobalStore(() => useState(initialState));

export default useResourceManager;
