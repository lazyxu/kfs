import { Skeleton } from "@mui/material";
import { useCallback, useEffect, useRef, useState } from "react";

const useGetState = (initiateState) => {
    const [state, setState] = useState(initiateState);
    const stateRef = useRef(state);
    stateRef.current = state;
    const getState = useCallback(() => stateRef.current, []);
    return [state, setState, getState];
};

export default ({ src, inView, onClick, title }) => {
    const [url, setUrl] = useState();
    const [loaded, setLoaded, getLoaded] = useGetState(0);
    const controller = useRef(new AbortController());
    useEffect(() => {
        const l = getLoaded();
        // console.log(src, inView, l);
        if (inView && l === 0) {
            setLoaded(1);
            // console.log(src, "fetch");
            fetch(src, {
                method: 'get',
                signal: controller.current.signal,
            }).then(response => {
                setLoaded(2);
                return response.blob();
            }).then(blob => {
                setUrl(URL.createObjectURL(blob));
            });
        }
        if (!inView && l === 1) {
            // console.log(src, "abort", controller.current);
            controller.current.abort();
            controller.current = new AbortController();
            setLoaded(0);
        }
    }, [inView]);
    if (url) {
        return (
            <img title={title} src={url} loading="lazy" onClick={onClick} />
        );
    } else {
        return (
            <Skeleton title={title} variant="rectangular" animation={false} width="100%" height="100%" />
        )
    }
};
