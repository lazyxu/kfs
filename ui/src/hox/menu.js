import { useState } from 'react';
import { createGlobalStore } from 'hox';

function useFunc() {
  const [menu, setMenu] = useState('');
  return {
    menu,
    setMenu,
  };
}

const [useMenu] = createGlobalStore(useFunc);

export default useMenu;
