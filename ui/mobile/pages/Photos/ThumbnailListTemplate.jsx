import { getSysConfig } from "@kfs/common/hox/sysConfig";
import { EventStreamContentType, fetchEventSource } from "@microsoft/fetch-event-source";
import { useEffect, useRef, useState } from "react";
import LongListTest from "./LongListTest";

export default function ({ metadataList, getTag, elementsPerLine }) {
    const [metadataTagList, setMetadataTagList] = useState([]);
    const sysConfig = getSysConfig();
    const [list, setList] = useState([]);
    const [dcimList, setDcimList] = useState([]);
    const abortControllerRef = useRef(null);
    
    // useEffect(() => {
    //     refresh();
    // }, [metadataList]);
    
    const fetchDCIMData = () => {
        if (abortControllerRef.current) {
            abortControllerRef.current.abort();
        }
        abortControllerRef.current = new AbortController();
        
        fetchEventSource(`${getSysConfig().webServer}/api/v1/event/listDCIM`, {
            signal: abortControllerRef.current.signal,
            openWhenHidden: true,
            async onopen(response) {
                if (response.ok && response.headers.get('content-type').includes(EventStreamContentType)) {
                    return; // everything's good
                }
                console.error(response);
                window.noteError("event/listDCIM.onopen: " + response.status);
            },
            onmessage(msg) {
                if (msg.event === 'FatalError') {
                    console.error(msg);
                    window.noteError("event/listDCIM.onmessage: " + msg);
                    return;
                }
                let info = JSON.parse(msg.data);
                console.log(info);
                if (info.errMsg) {
                    window.noteError(info.errMsg);
                }
                if (info.dcimList) {
                    setDcimList(prev => [...prev, ...info.dcimList]);
                }
            },
            onclose() {
                // connection closed
            },
            onerror(err) {
                console.error(err);
            }
        });
    };
    
    const processData = (originlist) => {
        let mtList = [];
        const l = [];
        let tagObj = {};
        for (let index = 0; index < originlist.length; index++) {
            const m = originlist[index];
            
            // 如果是从 DCIM 获取的数据，需要解析时间字段
            if (m.time && !m.year) {
                const timeDate = new Date(m.time);
                m.year = timeDate.getFullYear();
                m.month = timeDate.getMonth() + 1; // getMonth() 返回 0-11
                m.day = timeDate.getDate();
            }
            
            l.push({
                url: `${sysConfig.webServer}/api/v1/image?id=${m.id}&size=${m.size}&md5=${m.md5}&sha256=${m.sha256}`,
                hash: m.md5 + '_' + m.sha256,
                type: m.type,
                duration: m.duration,
                height: m.height && m.height.Int64 ? m.height.Int64 : m.height,
                width: m.width && m.width.Int64 ? m.width.Int64 : m.width,
            });
            const tag = getTag(m);
            if (tagObj.tag !== tag) {
                tagObj.end = index - 1;
                tagObj = { tag, start: index };
                mtList.push(tagObj);
                mtList.push({ index, hash: m.md5 + '_' + m.sha256 });
            } else {
                mtList.push({ index, hash: m.md5 + '_' + m.sha256 });
            }
        }
        tagObj.end = originlist.length - 1;
        setList(l);
        setMetadataTagList(mtList);
        window.noteInfo("刷新照片成功");
    };

    const refresh = async () => {
        try {
            window.noteInfo("开始刷新照片...");
            console.log("refresh")
            if (metadataList) {
                processData(metadataList);
            } else {
                // 重新获取 DCIM 数据
                setDcimList([]);
                fetchDCIMData();
                return; // 等待新数据通过 useEffect 触发渲染
            }
        } catch (e) {
            window.noteError("刷新照片失败：" + (typeof e.response?.data === 'string' ? e.response?.data : e.message));
        }
    }
    
    useEffect(() => {
        if (!metadataList) {
            // 只有当没有传入 metadataList 时才获取 DCIM 数据
            setDcimList([]);
            fetchDCIMData();
        }
        return () => {
            if (abortControllerRef.current) {
                abortControllerRef.current.abort();
            }
        }
    }, []);
    
    useEffect(() => {
        // 当 dcimList 更新时重新处理数据
        if (dcimList.length > 0) {
            processData(dcimList);
        } else if (metadataList) {
            processData(metadataList);
        }
    }, [dcimList, metadataList]);
    
    return (
        // <ThumbnailList metadataTagList={metadataTagList} list={list} refresh={refresh} elementsPerLine={elementsPerLine} />
        <LongListTest metadataTagList={metadataTagList} list={list} refresh={refresh} elementsPerLine={elementsPerLine} />
    );
}
