import React from 'react';
import styled from 'styled-components';
import { Icon } from 'kfs-components';
import { getBranchList } from 'bus/grpcweb';
import { StoreContext } from 'bus/bus';

export default function () {
  const [options, setOptions] = React.useState([]);
  const Branch = styled.div`
  `;
  const context = React.useContext(StoreContext);
  const branches = React.useRef();
  if (!branches.current) {
    // componentWillMount
    branches.current = true;
    getBranchList().then(list => {
      setOptions(list);
      const branch = list[0];
      context.setState({ branch });
      context.cd(branch, context.state.pwd);
    });
  }
  return (
    <Branch>
      <Icon
        icon="git"
        color="#cccccc"
        size="1.5em"
        hoverColor="white"
        hoverCursor="pointer"
      />
      <select
        onChange={e => {
          const branch = e.target.value;
          const pwd = '/';
          context.setState({ branch, pwd });
          context.cd(branch, pwd);
        }}
      >
        {options.map(o => <option key={o} value={o}>{o}</option>)}
      </select>
    </Branch>
  );
}
