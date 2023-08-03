import { useState } from 'react';
import { createGlobalStore } from 'hox';

function useFunc() {
  const [menu, setMenu] = useState('备份任务');
  return {
    menu,
    setMenu,
  };
}

const [useMenu] = createGlobalStore(useFunc);

export default useMenu;
