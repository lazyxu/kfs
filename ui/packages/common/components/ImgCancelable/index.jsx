import { useCallback, useEffect, useRef, useState } from "react";

const useGetState = (initiateState) => {
    const [state, setState] = useState(initiateState);
    const stateRef = useRef(state);
    stateRef.current = state;
    const getState = useCallback(() => stateRef.current, []);
    return [state, setState, getState];
};

export default ({ isNative, src, inView = true, renderImg, renderSkeleton, onLoaded }) => {
    const [url, setUrl, getUrl] = useGetState();
    const [loaded, setLoaded, getLoaded] = useGetState(0);
    const controller = useRef(new AbortController());
    useEffect(() => {
        return () => {
            // console.log("unmount", src);
            controller.current.abort();
        }
    }, []);
    useEffect(() => {
        const l = getLoaded();
        // console.log(src, inView, l);
        if (inView && l === 0) {
            setLoaded(1);
            // console.log("mount", src);
            fetch(src, {
                method: 'get',
                signal: controller.current.signal,
            }).then(response => {
                setLoaded(2);
                return response.blob();
            }).then(blob => {
                if (isNative) {
                    const fileReaderInstance = new FileReader();
                    fileReaderInstance.readAsDataURL(blob);
                    fileReaderInstance.onload = () => {
                        base64data = fileReaderInstance.result;
                        setUrl(base64data);
                    }
                } else {
                    setUrl(URL.createObjectURL(blob));
                }
                onLoaded?.();
            }).catch(e => {
                if (e.name !== 'AbortError') {
                    window.noteError("下载文件缩略图失败：" + (typeof e.response?.data === 'string' ? e.response?.data : e.message));
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
