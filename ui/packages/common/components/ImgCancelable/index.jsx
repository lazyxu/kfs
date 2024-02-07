import { useCallback, useEffect, useRef, useState } from "react";

const useGetState = (initiateState) => {
    const [state, setState] = useState(initiateState);
    const stateRef = useRef(state);
    stateRef.current = state;
    const getState = useCallback(() => stateRef.current, []);
    return [state, setState, getState];
};

export default ({ isNative, src, inView = true, renderImg, renderSkeleton, onLoaded, tag }) => {
    const [url, setUrl, getUrl] = useGetState();
    const [loaded, setLoaded, getLoaded] = useGetState(0);
    const controller = useRef(new AbortController());
    const aborted = useRef(false);
    useEffect(() => {
        return () => {
            console.log("ImgCancelable.unmount", tag);
            controller.current.abort();
            aborted.current = true;
        }
    }, []);
    useEffect(() => {
        const l = getLoaded();
        console.log("ImgCancelable.update", inView, l, tag);
        // console.log(src, inView, l);
        if (inView && l === 0) {
            setLoaded(1);
            // console.log("mount", src);
            aborted.current = false;
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
                if (!aborted.current) {
                    console.log("下载文件缩略图失败", tag, src, typeof e.response?.data === 'string' ? e.response?.data : e.message);
                    window.noteError("下载文件缩略图失败：" + (typeof e.response?.data === 'string' ? e.response?.data : e.message));
                }
            });
        }
        if (!inView && l === 1) {
            console.log("abort", src);
            controller.current.abort();
            aborted.current = true;
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
