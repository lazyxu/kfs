import { getSysConfig } from "@kfs/common/hox/sysConfig";
import { useEffect, useRef, useState } from "react";
import { fetchEventSource } from "@kfs/common/api/webServer/fetchEventSource";
import ThumbnailList from "@kfs/common/components/ThumbnailList/ThumbnailListRN";
import { listDCIM } from "@kfs/common/api/webServer/eventsource";
import { Text } from "react-native-paper";

export default function ({ metadataList, getTag, elementsPerLine }) {
    const sysConfig = getSysConfig();
    const [metadataTagList, setMetadataTagList] = useState([]);
    const [list, setList] = useState([]);
    const [dcimList, setDcimList] = useState([]);
    const [metadataListTotal, setMetadataListTotal] = useState(0);
    const abortControllerRef = useRef(null);
    
    // useEffect(() => {
    //     refresh();
    // }, [metadataList]);
    
    const fetchDCIMData = () => {
        if (abortControllerRef.current) {
            abortControllerRef.current.abort();
        }
        abortControllerRef.current = new AbortController();

        listDCIM(abortControllerRef.current.signal, () => {
            setDcimList([]);
            setMetadataListTotal(0);
            setList([]);
            setMetadataTagList([]);
        }, json => {
            if (json.dcimList) {
                setDcimList(prev => [...prev, ...json.dcimList]);
            } else {
                setMetadataListTotal(json.n);
            }
        }, (tag, err) => {
            window.noteError(tag+": " + err);
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
                thumbnailUrl: `${sysConfig.webServer}/thumbnail?id=${m.id}&size=${m.size}&md5=${m.md5}&sha256=${m.sha256}&hw=256&thumbnail=true`,
                uri: `${sysConfig.webServer}/api/v1/image?id=${m.id}&size=${m.size}&md5=${m.md5}&sha256=${m.sha256}`,
                id: m.id,
                size: m.size,
                md5: m.md5,
                sha256: m.sha256,
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
            }
            mtList.push({ index, hash: m.md5 + '_' + m.sha256 });
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
        console.log('dcimList:', dcimList);
        console.log('metadataList:', metadataList);
        // 当 dcimList 更新时重新处理数据
        if (dcimList.length > 0) {
            console.log('dcimList[0]:', dcimList[0]);
            processData(dcimList);
        } else if (metadataList) {
            console.log('metadataList[0]:', metadataList[0]);
            processData(metadataList);
        }
    }, [dcimList, metadataList]);
    
    return (
        <>
            <Text>({list.length}/{metadataListTotal})</Text>
            <ThumbnailList metadataTagList={metadataTagList} list={list} refresh={refresh} elementsPerLine={elementsPerLine} />
        </>
    );
}
