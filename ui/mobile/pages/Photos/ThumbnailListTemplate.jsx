import { getSysConfig } from "@kfs/common/hox/sysConfig";
import { useEffect, useState } from "react";
import LongListTest from "./LongListTest";

export default function ({ metadataList, listDCIMMetadataTime, getTag, elementsPerLine }) {
    const [metadataTagList, setMetadataTagList] = useState([]);
    const sysConfig = getSysConfig();
    const [list, setList] = useState([]);
    // useEffect(() => {
    //     refresh();
    // }, [metadataList]);
    const refresh = async () => {
        console.log("refresh")
        let originlist;
        if (metadataList) {
            originlist = metadataList;
        } else {
            originlist = await listDCIMMetadataTime();
        }
        let tag = "";
        let mtList = [];
        let lineList;
        const l = [];
        // let backList = originlist;
        // originlist = [];
        // for (let i = 0; i < backList.length; i++) {
        //     for (let j = 0; j < 100000; j++) {
        //         originlist.push(backList[i]);
        //     }
        // }
        // originlist = originlist.slice(0, 100);
        for (let index = 0; index < originlist.length; index++) {
            const m = originlist[index];
            l.push({
                url: `${sysConfig.webServer}/api/v1/image?hash=${m.hash}`,
                hash: m.hash,
                type: m.fileType.type,
                duration: m.duration,
                height: m.heightWidth.height,
                width: m.heightWidth.width,
            });
            const curTag = getTag(m);
            if (tag !== curTag) {
                tag = curTag;
                mtList.push(tag);
                mtList.push({ index, hash: m.hash });
            } else {
                mtList.push({ index, hash: m.hash });
            }
        }
        setList(l);
        setMetadataTagList(mtList);
    }
    useEffect(() => {
        refresh();
    }, []);
    return (
        // <ThumbnailList metadataTagList={metadataTagList} list={list} refresh={listDCIMMetadataTime ? refresh : undefined} elementsPerLine={elementsPerLine} />
        <LongListTest metadataTagList={metadataTagList} list={list} refresh={listDCIMMetadataTime ? refresh : undefined} elementsPerLine={elementsPerLine} />
    );
}
