import { useState } from 'react';
import { createGlobalStore } from 'hox';

function useFunc() {
  const [menu, setMenu] = useState('文件');
  return {
    menu,
    setMenu,
  };
}

const [useMenu] = createGlobalStore(useFunc);

export default useMenu;
