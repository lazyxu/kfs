import React from 'react';
import styled from 'styled-components';
import Icon from 'components/Icon';
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
      context.setState({ branch: list[0] });
      context.cd(context.state.pwd);
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
          context.setState({ branch: e.target.value, pwd: '/' });
          context.cd(context.state.pwd);
        }}
      >
        {options.map(o => <option key={o} value={o}>{o}</option>)}
      </select>
    </Branch>
  );
}
