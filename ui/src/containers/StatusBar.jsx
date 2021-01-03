import React from 'react';

import styled from 'styled-components';

import Icon from 'components/Icon';
import Informations from 'components/Informations';
import Notifications from 'containers/Notifications';
import { inState } from 'bus/bus';
import moment from 'moment';
import 'moment/locale/zh-cn';

const StatusBar = styled.div`
  position: relative;
  height: 100%;
  width: calc(100% - 0.5em);
  background-color: var(--footer-color);
  display: flex;
  padding: 0 0.5em 0 0;
`;
const Empty = styled.div`
  flex: 1;
`;
const Notification = styled.div`
  height: 100%;
  :hover {
    background-color: var(--footer-hover-color);
    cursor: pointer;
  }
`;
const Text = styled.div`
  line-height: 1.5em;
  height: 100%;
  padding: 0 0.5em;
  user-select: none;
  font-size: 1em;
`;
const ClickableText = styled(Text)`
  :hover {
    background-color: var(--footer-hover-color);
    cursor: pointer;
  }
`;
const Triangle = styled.div`
  position: fixed;
  right: 1em;
  bottom: 1.5em;
  width: 0;
  height: 0;
  border: 0.5em solid transparent;
  border-bottom: 0.5em solid var(--footer-color);
`;

function toDecimal2NoZero(x) {
  const f = Math.round(x * 100) / 100;
  const s = f.toString();
  return s;
}

function formatDate(ms) {
  return moment(ms).locale('zh-cn').format('llll');
}

@inState('showNotifications', 'showDateInformations', 'files', 'fileSize', 'atimems', 'mtimems', 'ctimems', 'birthtimems',
  'chosen', 'boxChosen')
class component extends React.Component {
  getFileSize() {
    let { fileSize } = this.state;
    if (typeof fileSize === 'number') {
      if (fileSize < 1024) {
        return `${toDecimal2NoZero(fileSize)} B`;
      }
      fileSize /= 1024;
      if (fileSize < 1024) {
        return `${toDecimal2NoZero(fileSize)} K`;
      }
      fileSize /= 1024;
      if (fileSize < 1024) {
        return `${toDecimal2NoZero(fileSize)} M`;
      }
      fileSize /= 1024;
      if (fileSize < 1024) {
        return `${toDecimal2NoZero(fileSize)} G`;
      }
      fileSize /= 1024;
      if (fileSize < 1024) {
        return `${toDecimal2NoZero(fileSize)} T`;
      }
    }
    return '';
  }

  getCnt() {
    const { chosen, boxChosen } = this.state;
    const result = {};
    Object.keys(chosen).forEach((k) => chosen[k] && (result[k] = chosen[k]));
    Object.keys(boxChosen).forEach((k) => boxChosen[k] && (result[k] = boxChosen[k]));
    return Object.values(result).filter((cnt) => cnt > 0).reduce((a, b) => a + b, 0);
  }

  render() {
    const {
      atimems, mtimems, ctimems, birthtimems, showDateInformations, showNotifications,
    } = this.state;
    const cnt = this.getCnt();
    return (
      <StatusBar>
        <Text>
          {` ${this.state.files.length} 个对象`}
        </Text>
        {cnt ? (
          <Text>
            {`选中 ${cnt} 个项目`}
          </Text>
        ) : null}
        <Text>{this.getFileSize()}</Text>
        <Empty />
        {mtimems ? showDateInformations && (
          <Informations informations={[
            `最近访问：${formatDate(atimems)}`,
            `最近写入：${formatDate(mtimems)}`,
            `最近修改：${formatDate(ctimems)}`,
            `创建时间：${formatDate(birthtimems)}`,
          ]}
          />
        ) : null}
        {mtimems
          ? (
            <ClickableText
              style={showDateInformations ? {
                backgroundColor: 'var(--footer-hover-color)',
                cursor: 'pointer',
              } : {}}
              onClick={() => {
                this.setState((prevState) => ({ showDateInformations: !prevState.showDateInformations }));
              }}
            >
              {formatDate(this.state.mtimems)}
            </ClickableText>
          ) : null}
        <Notification
          style={showNotifications ? {
            backgroundColor: 'var(--footer-hover-color)',
            cursor: 'pointer',
          } : {}}
        >
          <Icon
            icon="notice"
            color="white"
            size="1em"
            padding="0.2em 0.5em 0 0.5em"
            onClick={() => this.setState(
              (prevState) => ({ showNotifications: !prevState.showNotifications }),
            )}
          />
        </Notification>
        {showNotifications && <Triangle />}
        {showNotifications && <Notifications />}
      </StatusBar>
    );
  }
}

export default component;
