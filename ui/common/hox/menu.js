import { useState } from 'react';
import { createModel } from 'hox';

function useMenu() {
  const [menu, setMenu] = useState('');
  return {
    menu,
    setMenu,
  };
}

export default createModel(useMenu);
