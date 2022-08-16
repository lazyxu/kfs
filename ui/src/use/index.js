import React from 'react';

export function useClick(onClick, onDoubleClick) {
    // const [clicked, setClicked] = React.useState();
    const handle = React.useRef();
    React.useEffect(() => {
        // console.log('--00--', clicked, handle);
        return () => {
            // console.log('--55--', clicked, handle);
            clearTimeout(handle.current);
        };
    }, []);
    return e => {
        if (e.button === 2) {
            return;
        }
        if (handle.current) {
            // console.log('--44--', clicked, handle);
            onDoubleClick && onDoubleClick(e);
            // setClicked();
            clearTimeout(handle.current);
            handle.current = undefined;
            return;
        }
        // console.log('--11--', clicked);
        // setClicked(true);
        handle.current = setTimeout(() => {
            // console.log('--22--', clicked);
            // setClicked();
            clearTimeout(handle.current);
            handle.current = undefined;
        }, 200);
        // console.log('--33--', clicked);
        onClick(e);
    };
}
