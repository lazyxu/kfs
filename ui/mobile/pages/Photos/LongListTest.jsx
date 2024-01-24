import React, { useState } from "react";
import { View } from "react-native";
import { Text, } from "react-native-paper";
import LongList from "./LongList";
import Thumbnail from "./Thumbnail";

export default function ({ metadataTagList, elementsPerLine = 5, list, refresh }) {
    const t0 = Date.now();
    const navigation = window.kfsNavigation;
    const [width, setWidth] = useState(0);
    // console.log("metadataTagList", metadataTagList, width)
    const itemHeightWidthList = [];
    let elementWidth = 0;
    if (width !== 0) {
        elementWidth = (width - 1) / elementsPerLine;
        for (let i = 0; i < metadataTagList.length; i++) {
            const data = metadataTagList[i];
            if (typeof data === 'object') {
                itemHeightWidthList[i] = { height: elementWidth, width: elementWidth };
            } else {
                itemHeightWidthList[i] = { height: 16 * 2, width: width };
            }
        }
    }

    const renderItem = function (_, index, cacheItem) {
        const data = metadataTagList[index];
        if (typeof data !== 'object') {
            cacheItem?.();
            return <View style={{ margin: 6 }}><Text style={{ fontSize: 16 }}>{data}</Text></View>
        }
        return <Thumbnail key={data.index} width={elementWidth} navigation={navigation} list={list} index={data.index} inView={true} onLoaded={cacheItem} />
    }

    console.log("LongListTest.1", Date.now() - t0);
    return <View
        style={{ flex: 1 }}
        onLayout={e => {
            const { layout } = e.nativeEvent;
            if (layout.width) {
                setWidth(layout.width);
            }
        }}>
        <LongList
            itemHeightWidthList={itemHeightWidthList}
            width={width}
            dataList={metadataTagList}
            renderItem={renderItem}
        />
    </View>;
}
