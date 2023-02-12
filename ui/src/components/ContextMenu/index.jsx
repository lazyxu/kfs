import React from 'react';
import styles from './index.module.scss';
import {ListItemText, MenuItem, MenuList, Paper} from "@mui/material";

const onFinishList = {};
let id = 0;

document.addEventListener('click', function () {
    Object.values(onFinishList).forEach(onFinish => {
        onFinish?.();
    });
}, false);

export default ({left, top, right, bottom, maxWidth, maxHeight, options, onFinish}) => {
    React.useEffect(() => {
        id++;
        onFinishList[id] = onFinish;
        return () => {
            delete onFinishList[id];
        };
    }, []);
    if (left + maxWidth > right) {
        left = right - maxWidth;
    }
    if (top + maxHeight > bottom) {
        top = bottom - maxHeight;
    }
    return (
        <Paper className={styles.contextMenu}
             style={{left: `${left}px`, top: `${top}px`,}}
             onMouseDown={(e) => e.stopPropagation()}
             onClick={e => {
                 e.stopPropagation();
             }}
        >   <MenuList>
            {Object.keys(options).map((o) => {
                const option = options[o];
                let fn = option;
                if (option.fn) {
                    fn = option.fn;
                    if (!option.enabled) {
                        return (
                            <MenuItem key={o} disabled={true}>
                                <ListItemText>{o}</ListItemText>
                            </MenuItem>
                        );
                    }
                }
                if (option?.type?.name === 'component') {
                    return (
                        <MenuItem key={o}>{option}</MenuItem>
                    );
                }
                return (
                    <MenuItem
                         key={o}
                         onMouseDown={(e) => {
                             e.stopPropagation();
                             fn(e);
                             onFinish?.();
                         }}
                    >
                        <ListItemText>{o}</ListItemText>
                    </MenuItem>
                );
            })}
        </MenuList>
        </Paper>
    );
};
