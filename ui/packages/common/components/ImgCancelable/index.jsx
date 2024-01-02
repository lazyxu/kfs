import { useCallback, useEffect, useRef, useState } from "react";
import { Platform } from 'react-native';

const useGetState = (initiateState) => {
    const [state, setState] = useState(initiateState);
    const stateRef = useRef(state);
    stateRef.current = state;
    const getState = useCallback(() => stateRef.current, []);
    return [state, setState, getState];
};

export default ({ src, inView, renderImg, renderSkeleton }) => {
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
                if (Platform.OS === "web") {
                    setUrl(URL.createObjectURL(blob));
                } else {
                    const fileReaderInstance = new FileReader();
                    fileReaderInstance.readAsDataURL(blob);
                    fileReaderInstance.onload = () => {
                        base64data = fileReaderInstance.result;
                        setUrl(base64data);
                    }
                }
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
        return renderImg(url);
    } else {
        return renderSkeleton();
    }
};
