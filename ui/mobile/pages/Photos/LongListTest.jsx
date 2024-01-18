import React, { useState } from "react";
import { View } from "react-native";
import { Text, } from "react-native-paper";
import LongList from "./LongList";
import Thumbnail from "./Thumbnail";

export default function ({ metadataTagList, elementsPerLine = 5, list, refresh }) {
    const t0 = Date.now();
    const navigation = window.kfsNavigation;
    const [widthThumbnail, setWidthThumbnail] = useState(0);

    const itemHeightList = [];
    for (let i = 0; i < metadataTagList.length; i++) {
        const data = metadataTagList[i];
        if (typeof data === 'object') {
            itemHeightList[i] = widthThumbnail;
        } else {
            itemHeightList[i] = 16 * 2;
        }
    }

    const rowRenderer = (_, index) => {
        const data = metadataTagList[index];
        return typeof data === 'object' ?
            <View style={{
                display: "flex",
                width: "100%",
                flexDirection: 'row',
                flexWrap: "wrap",
                alignContent: "flex-start"
            }}>
                {data.map(({ hash, index }) =>
                    <Thumbnail key={index} width={widthThumbnail} navigation={navigation} list={list} index={index} inView={true} />
                )}
            </View> :
            <View style={{ margin: 6 }}><Text style={{ fontSize: 16 }}>{data}</Text></View>
    }

    console.log("LongListTest.1", Date.now() - t0);
    return <View
        style={{ flex: 1 }}
        onLayout={e => {
            const { layout } = e.nativeEvent;
            if (layout.width) {
                const w = layout.width / elementsPerLine;
                setWidthThumbnail(w);
            }
        }}>
        <LongList
            itemHeightList={itemHeightList}
            dataList={metadataTagList}
            renderItem={rowRenderer}
        />
    </View>;
}
