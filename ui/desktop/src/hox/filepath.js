import { useState, useEffect } from 'react';
import { createModel } from 'hox';

function useFilePath() {
  const [filepath, setFilepath] = useState('');
  return [
    filepath,
    setFilepath,
  ];
}

export default createModel(useFilePath);
