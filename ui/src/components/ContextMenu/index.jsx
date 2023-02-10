import React from 'react';
import styles from './index.module.scss';

const onFinishList = {};
let id = 0;

document.addEventListener('click', function () {
    Object.values(onFinishList).forEach(onFinish => {
        onFinish?.();
    });
}, false);

export default ({
                    left, top, options, onFinish,
                }) => {
    React.useEffect(() => {
        id++;
        onFinishList[id] = onFinish;
        return () => {
            delete onFinishList[id];
        };
    }, []);
    return (
        <div className={styles.contextMenu}
             style={{left: `${left}px`, top: `${top}px`,}}
             onMouseDown={(e) => e.stopPropagation()}
             onClick={e => {
                 e.stopPropagation();
             }}
        >
            {Object.keys(options).map((o) => {
                const option = options[o];
                let fn = option;
                if (option.fn) {
                    fn = option.fn;
                    if (!option.enabled) {
                        return (
                            <div className={styles.disable} key={o}>{o}</div>
                        );
                    }
                }
                if (option?.type?.name === 'component') {
                    return (
                        <div className={styles.option} key={o}>{option}</div>
                    );
                }
                return (
                    <div className={styles.option}
                         key={o}
                         onMouseDown={(e) => {
                             fn(e);
                             onFinish && onFinish();
                             e.stopPropagation();
                         }}
                    >
                        {o}
                    </div>
                );
            })}
        </div>
    );
};
